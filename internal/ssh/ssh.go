package ssh

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

		baseName, _, _ := strings.Cut(entry.Name(), ".")
		fullPath := filepath.Join(sshFolderPath, entry.Name())
		fingerprintFull, err := s.GetFingerprint(fullPath)
		if err != nil {
			return nil, err
		}

		fingerprint := strings.Split(fingerprintFull, " ")[1]

		// Break the comment and algo out of the fingerprint
		splits := strings.Split(fingerprintFull, " ")
		comment := strings.Join(splits[2:len(splits)-1], " ")
		algorithm := strings.Trim(
			strings.TrimSpace(strings.Join(splits[len(splits)-1:], " ")),
			"()",
		)

		// Check if the private form exists
		privPath := filepath.Join(sshFolderPath, baseName)
		privExists, err := fileExists(privPath)
		if err != nil {
			return nil, err
		}

		keys = append(keys, SSHKey{
			Name:           baseName,
			PublicKeyPath:  fullPath,
			Algorithm:      algorithm,
			Fingerprint:    fingerprint,
			Comment:        &comment,
			HasPrivateKey:  privExists,
			PrivateKeyPath: &privPath,
			// TODO: Check if in agent, and get permissions for file
		})
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

func (s SSHService) DeleteKey(path string) error {
	return os.Remove(path)
}

func (s SSHService) GetFingerprint(keyPath string) (string, error) {
	fingerprint, err := executeCmd("ssh-keygen", []string{
		"-lf", keyPath,
	}...)
	if err != nil {
		return "", err
	}

	return fingerprint, nil
}

func (s SSHService) GetSSHConfig(host string) (any, error) {
	panic("TODO")
	return nil, nil
}

func (s SSHService) DeployKey() {
	panic("TODO")
}
