package ssh

import (
	"os"
	"os/exec"
)

type SSHKey struct {
	Name string
	Path string
	Pub  string
}

type SSHService struct{}

func NewSSHService() *SSHService {
	return &SSHService{}
}

func (s SSHService) GetExistingKeys() ([]SSHKey, error) {
	entries, err := os.ReadDir("~/.ssh")
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

		fullPath := "~/.ssh/" + entry.Name()

		contents, err := os.ReadFile(fullPath)
		if err != nil {
			return nil, err
		}

		keys = append(keys, SSHKey{
			Name: entry.Name(),
			Path: fullPath,
			Pub:  string(contents),
		})
	}

	return keys, nil
}

func (s SSHService) executeCmd(cmd string, args ...string) (string, error) {
	c := exec.Command(cmd, args...)
	out, err := c.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func isPubKey(name string) bool {
	return len(name) >= 4 && name[len(name)-4:] == ".pub"
}
