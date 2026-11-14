package ssh

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

type RemoteOS string

const (
	RemoteLinux   RemoteOS = "linux"
	RemoteMacOS   RemoteOS = "darwin"
	RemoteWindows RemoteOS = "windows"
	RemoteUnknown RemoteOS = "unknown"
)

type KeyDeployer interface {
	InstallPublicKey(
		client *ssh.Client,
		publicKey []byte,
	) error
}

type UnixKeyDeployer struct{}

func (d UnixKeyDeployer) InstallPublicKey(client *ssh.Client, publicKey []byte) error {
	normalized, err := normalizePublicKey(publicKey)
	if err != nil {
		return err
	}

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	session.Stdin = strings.NewReader(normalized)

	const command = `set -eu

umask 077
mkdir -p "$HOME/.ssh"
chmod 700 "$HOME/.ssh"

touch "$HOME/.ssh/authorized_keys"
chmod 600 "$HOME/.ssh/authorized_keys"

IFS= read -r incoming || [ -n "$incoming" ]

key_type=$(printf '%s\n' "$incoming" | awk '{print $1}')
key_data=$(printf '%s\n' "$incoming" | awk '{print $2}')

if ! awk -v type="$key_type" -v data="$key_data" '
    $1 == type && $2 == data { found = 1 }
    END { exit !found }
' "$HOME/.ssh/authorized_keys"; then
    printf '%s\n' "$incoming" >> "$HOME/.ssh/authorized_keys"
fi`

	output, err := session.CombinedOutput(command)
	if err != nil {
		return fmt.Errorf("install public key: %w: %s", err, output)
	}

	return nil
}

type WindowsKeyDeployer struct{}

func (d WindowsKeyDeployer) InstallPublicKey(client *ssh.Client, publicKey []byte) error {
	normalized, err := normalizePublicKey(publicKey)
	if err != nil {
		return err
	}

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	session.Stdin = strings.NewReader(normalized)

	const script = `
$ErrorActionPreference = 'Stop'

$key = [Console]::In.ReadLine()
if ([string]::IsNullOrWhiteSpace($key)) {
    throw 'No public key received'
}

$sshDir = Join-Path $env:USERPROFILE '.ssh'
$authorizedKeys = Join-Path $sshDir 'authorized_keys'

if (-not (Test-Path -LiteralPath $sshDir)) {
    New-Item -ItemType Directory -Path $sshDir | Out-Null
}

if (-not (Test-Path -LiteralPath $authorizedKeys)) {
    New-Item -ItemType File -Path $authorizedKeys | Out-Null
}

$parts = $key -split '\s+'
$keyType = $parts[0]
$keyData = $parts[1]

$exists = $false

foreach ($line in [System.IO.File]::ReadAllLines($authorizedKeys)) {
    $fields = $line.Trim() -split '\s+'

    for ($i = 0; $i -lt ($fields.Length - 1); $i++) {
        if ($fields[$i] -eq $keyType -and
            $fields[$i + 1] -eq $keyData) {
            $exists = $true
            break
        }
    }

    if ($exists) { break }
}

if (-not $exists) {
    [System.IO.File]::AppendAllText(
        $authorizedKeys,
        $key + [Environment]::NewLine,
        [System.Text.UTF8Encoding]::new($false)
    )
}

Write-Output 'Public key installed'
`

	// PowerShell's -EncodedCommand expects UTF-16LE.
	command := powershellEncodedCommand(script)

	output, err := session.CombinedOutput(command)
	if err != nil {
		return fmt.Errorf(
			"Windows key deployment failed: %w: %s",
			err,
			strings.TrimSpace(string(output)),
		)
	}

	return nil
}

type DeployRequest struct {
	Host              string
	Port              int
	Username          string
	Password          string
	PublicKeyPath     string
	RemoteMachineOS   RemoteOS
	HostKeyApprovalID string `json:"hostKeyApprovalID,omitempty"`
}

type DeployResult struct {
	HostKey *HostKeyPrompt `json:"hostKey,omitempty"`
}

func (s SSHService) DeployKey(req DeployRequest) (DeployResult, error) {
	address, err := connectionAddress(req.Host, req.Port)
	if err != nil {
		return DeployResult{}, err
	}

	opts := ConnectionOptions{
		Host:     req.Host,
		Port:     req.Port,
		Username: req.Username,
		Password: req.Password,
	}
	if err != nil {
		return DeployResult{}, err
	}

	if req.HostKeyApprovalID != "" {
		opts.approvedHostKey, err = s.hostKeyApprovals.take(req.HostKeyApprovalID, address)
		if err != nil {
			return DeployResult{}, err
		}
	}

	client, err := Connect(opts)
	var unknown *unknownHostKey
	if errors.As(err, &unknown) {
		prompt := s.hostKeyApprovals.issue(unknown)
		return DeployResult{HostKey: &prompt}, nil
	}
	if err != nil {
		return DeployResult{}, err
	}
	defer client.Close()

	if opts.approvedHostKey != nil {
		if err := rememberHostKey(address, client.RemoteAddr(), opts.approvedHostKey); err != nil {
			return DeployResult{}, err
		}
	}

	// Read local public keyfile
	data, err := os.ReadFile(req.PublicKeyPath)
	if err != nil {
		return DeployResult{}, err
	}

	switch req.RemoteMachineOS {
	case RemoteLinux, RemoteMacOS:
		return DeployResult{}, UnixKeyDeployer{}.InstallPublicKey(client, data)
	default:
		return DeployResult{}, fmt.Errorf("unknown/unsupported OS %s", req.RemoteMachineOS)
	}
}
