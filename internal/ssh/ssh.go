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

func (s SSHService) GetPublicSSHKey(sshKey SSHKey) (string, error) {
	data, err := os.ReadFile(sshKey.PublicKeyPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s SSHService) GetSSHConfig(host string) (any, error) {
	panic("TODO")
}

func (s SSHService) DeployKey() {
	panic("TODO")
}
