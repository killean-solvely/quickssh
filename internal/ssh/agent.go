package ssh

import (
	"fmt"
	"strings"
)

func (s SSHService) AddKeyToAgent(key SSHKey) error {
	if !key.HasPrivateKey || key.PrivateKeyPath == nil {
		return fmt.Errorf("this key has no private key")
	}

	_, err := executeCmd(
		"ssh-add",
		*key.PrivateKeyPath,
	)

	return err
}

func (s SSHService) RemoveKeyFromAgent(key SSHKey) error {
	if !key.HasPrivateKey || key.PrivateKeyPath == nil {
		return fmt.Errorf("this key has no private key")
	}

	_, err := executeCmd(
		"ssh-add",
		"-d",
		*key.PrivateKeyPath,
	)

	return err
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
