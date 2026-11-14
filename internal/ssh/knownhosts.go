package ssh

import (
	"crypto/rand"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

type unknownHostKey struct {
	address string
	key     ssh.PublicKey
}

func (e *unknownHostKey) Error() string {
	return "unknown SSH host key for " + e.address
}

func verifyHostKey(check ssh.HostKeyCallback, approved ssh.PublicKey) ssh.HostKeyCallback {
	return func(address string, remote net.Addr, key ssh.PublicKey) error {
		if approved != nil {
			if err := ssh.FixedHostKey(approved)(address, remote, key); err != nil {
				return fmt.Errorf("server presented a different key on retry: %w", err)
			}
		}

		err := check(address, remote, key)
		var keyErr *knownhosts.KeyError
		if errors.As(err, &keyErr) && len(keyErr.Want) == 0 {
			if approved != nil {
				return nil // The key exactly matches this attempt's user approval
			}
			return &unknownHostKey{address: address, key: key}
		}
		return err // Includes mismatches and revoked keys
	}
}

type HostKeyPrompt struct {
	ApprovalID  string `json:"approvalId"`
	Address     string `json:"address"`
	Algorithm   string `json:"algorithm"`
	Fingerprint string `json:"fingerprint"`
}

type pendingHostKey struct {
	address   string
	key       ssh.PublicKey
	expiresAt time.Time
}

type hostKeyApprovals struct {
	mu      sync.Mutex
	pending map[string]pendingHostKey
}

func (a *hostKeyApprovals) issue(unknown *unknownHostKey) HostKeyPrompt {
	id := rand.Text()
	now := time.Now()

	a.mu.Lock()
	defer a.mu.Unlock()
	if a.pending == nil {
		a.pending = make(map[string]pendingHostKey)
	}
	for id, approval := range a.pending {
		if !now.Before(approval.expiresAt) {
			delete(a.pending, id)
		}
	}
	a.pending[id] = pendingHostKey{unknown.address, unknown.key, now.Add(2 * time.Minute)}

	return HostKeyPrompt{
		ApprovalID:  id,
		Address:     unknown.address,
		Algorithm:   unknown.key.Type(),
		Fingerprint: ssh.FingerprintSHA256(unknown.key),
	}
}

func (a *hostKeyApprovals) take(id, address string) (ssh.PublicKey, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	approval, ok := a.pending[id]
	delete(a.pending, id)
	if !ok || approval.address != address || !time.Now().Before(approval.expiresAt) {
		return nil, errors.New(
			"SSH host approval has expired or does not match. Start the deployment again",
		)
	}
	return approval.key, nil
}

var hostKeysWriteMu sync.Mutex

func loadHostKeys() (ssh.HostKeyCallback, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	appDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	dirs := []string{
		filepath.Join(home, ".ssh", "known_hosts"),
		filepath.Join(appDir, "quickssh", "known_hosts"),
	}

	var files []string
	for _, path := range dirs {
		if _, err := os.Stat(path); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return nil, err
		}
		files = append(files, path)
	}
	return knownhosts.New(files...)
}

func rememberHostKey(address string, remote net.Addr, key ssh.PublicKey) error {
	hostKeysWriteMu.Lock()
	defer hostKeysWriteMu.Unlock()

	check, err := loadHostKeys()
	if err != nil {
		return err
	}
	err = check(address, remote, key)
	if err == nil {
		return nil
	}
	var keyErr *knownhosts.KeyError
	if !errors.As(err, &keyErr) || len(keyErr.Want) != 0 {
		return err
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	dir := filepath.Join(configDir, "quickssh")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	path := filepath.Join(dir, "known_hosts")
	saved, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if len(saved) > 0 && saved[len(saved)-1] != '\n' {
		saved = append(saved, '\n')
	}
	host := knownhosts.HashHostname(knownhosts.Normalize(address))
	saved = append(saved, knownhosts.Line([]string{host}, key)+"\n"...)

	// Write to a temporary file so a failed write leaves the saved keys intact.
	tmp, err := os.CreateTemp(dir, ".known_hosts-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	_, writeErr := tmp.Write(saved)
	if err := errors.Join(writeErr, tmp.Close()); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}
