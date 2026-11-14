package ssh

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"strings"
	"unicode/utf16"

	"golang.org/x/crypto/ssh"
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

func powershellEncodedCommand(script string) string {
	chars := utf16.Encode([]rune(script))

	data := make([]byte, len(chars)*2)
	for i, char := range chars {
		binary.LittleEndian.PutUint16(
			data[i*2:],
			char,
		)
	}

	encoded := base64.StdEncoding.EncodeToString(data)

	return "powershell.exe -NoProfile -NonInteractive " +
		"-EncodedCommand " + encoded
}

func normalizePublicKey(data []byte) (string, error) {
	key, comment, options, rest, err := ssh.ParseAuthorizedKey(data)
	if err != nil || len(bytes.TrimSpace(rest)) != 0 {
		return "", fmt.Errorf("invalid public key file")
	}

	if len(options) != 0 {
		return "", fmt.Errorf("authorized-key options are not supported for deployment")
	}

	normalized := strings.TrimPrefix(
		string(ssh.MarshalAuthorizedKey(key)),
		"\n",
	)

	if comment != "" {
		normalized += " " + comment
	}

	normalized += "\n"

	return normalized, nil
}
