package ssh

import (
	"errors"
	"io/fs"
	"os"
	"os/exec"
)

func executeCmd(cmd string, args ...string) (string, error) {
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

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}
