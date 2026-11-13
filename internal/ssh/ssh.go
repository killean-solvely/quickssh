package ssh

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"
)

type SSHKey struct {
	Name           string
	PublicKeyPath  string
	PrivateKeyPath *string
	Algorithm      string
	Fingerprint    string
	Comment        *string
	HasPrivateKey  bool
	InAgent        bool
	Permissions    *string
}

type SSHService struct{}

func NewSSHService() *SSHService {
	return &SSHService{}
}

func (s SSHService) GetExistingKeys() ([]SSHKey, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	sshFolderPath := filepath.Join(home, ".ssh")

	entries, err := os.ReadDir(sshFolderPath)
	if err != nil {
		return nil, err
	}

	agentKeys, err := s.GetAgentFingerprints()
	if err != nil {
		agentKeys = make(map[string]bool)
	}

	keys := []SSHKey{}
	for _, entry := range entries {
		// Skip directories
		if entry.IsDir() {
			continue
		}

		// Skip non public key files
		if !isPubKey(entry.Name()) {
			continue
		}

		baseName := strings.TrimSuffix(entry.Name(), ".pub")
		fullPath := filepath.Join(sshFolderPath, entry.Name())

		data, err := os.ReadFile(fullPath)
		if err != nil {
			return nil, err
		}

		// Get algo and fingerprint
		pubKey, comment, _, _, err := ssh.ParseAuthorizedKey(data)
		algorithm := pubKey.Type()
		fingerprint := ssh.FingerprintSHA256(pubKey)

		// Check if the private form exists
		privPath := filepath.Join(sshFolderPath, baseName)
		privExists, err := fileExists(privPath)
		if err != nil {
			return nil, err
		}

		k := SSHKey{
			Name:          baseName,
			PublicKeyPath: fullPath,
			Algorithm:     algorithm,
			Fingerprint:   fingerprint,
			Comment:       &comment,
			HasPrivateKey: privExists,
			InAgent:       agentKeys[fingerprint],
		}

		// Get file permissions
		info, err := os.Stat(privPath)
		if err == nil {
			v := fmt.Sprintf("%04o", info.Mode().Perm())
			k.Permissions = &v
		}

		if privExists {
			k.PrivateKeyPath = &privPath
		}

		keys = append(keys, k)
	}

	return keys, nil
}

func (s SSHService) CreateKey(
	keyName string,
	keyPassphrase string,
	keyType string,
	comment string,
) error {
	// ssh-keygen -t ed25519 -f "$HOME/.ssh/my_new_key" -N "" -C "comment"
	/*
		-t ed25519: The type of key to create (e.g., ed25519, rsa, ecdsa).
		-f /path/to/key: The output file path. Providing this bypasses the "Enter file in which to save the key" prompt.
		-N "": The new passphrase. Providing an empty string "" bypasses the passphrase prompt entirely. If your user does want a passphrase via your GUI, you just pass it here: -N "my_super_secret_pass".
		-C "comment": The comment (usually an email address) attached to the public key.
		Optional: -b 4096 (If you are generating an RSA key and want to specify the bit size).
	*/

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Failed to get homedir\n")
		return err
	}

	_, err = executeCmd("ssh-keygen", []string{
		"-t", keyType,
		"-f", filepath.Join(homeDir, ".ssh", keyName),
		"-N", keyPassphrase,
		"-C", comment,
	}...)
	if err != nil {
		fmt.Printf("Failed to get execute: %s\n", err.Error())
		return err
	}

	return nil
}

func (s SSHService) DeleteKey(key SSHKey) error {
	if err := os.Remove(key.PublicKeyPath); err != nil {
		return err
	}
	if key.HasPrivateKey && key.PrivateKeyPath != nil {
		if err := os.Remove(*key.PrivateKeyPath); err != nil {
			return err
		}
	}
	return nil
}

func (s SSHService) GetAgentFingerprints() (map[string]bool, error) {
	output, err := executeCmd("ssh-add", "-l", "-E", "sha256")
	if err != nil {
		return nil, err
	}

	fingerprints := map[string]bool{}
	for _, line := range strings.Split(output, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		fingerprints[fields[1]] = true
	}
	return fingerprints, nil
}

func (s SSHService) GetSSHConfig(host string) (any, error) {
	panic("TODO")
	return nil, nil
}

func (s SSHService) DeployKey() {
	panic("TODO")
}
