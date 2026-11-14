package ssh

import (
	"fmt"
	"net"
	"net/netip"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

type ConnectionOptions struct {
	Host     string
	Port     int
	Username string
	Password string

	approvedHostKey ssh.PublicKey
}

func connectionAddress(host string, port int) (string, error) {
	host = strings.TrimSpace(host)
	if host == "" || port < 1 || port > 65535 {
		return "", fmt.Errorf("enter a host and a port between 1 and 65535")
	}

	if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
		host = host[1 : len(host)-1]
	}

	_, ipErr := netip.ParseAddr(host)
	if host == "" || (strings.ContainsAny(host, "[]:/") && ipErr != nil) {
		return "", fmt.Errorf("enter a hostname or IP address in the host field")
	}

	return net.JoinHostPort(host, strconv.Itoa(port)), nil
}

func Connect(opts ConnectionOptions) (*ssh.Client, error) {
	address, err := connectionAddress(opts.Host, opts.Port)
	if err != nil {
		return nil, err
	}

	hostKeyCallback, err := loadHostKeys()
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: opts.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(opts.Password),
		},
		HostKeyCallback: verifyHostKey(hostKeyCallback, opts.approvedHostKey),
		Timeout:         10 * time.Second,
	}

	return ssh.Dial("tcp", address, config)
}

func TestConnection(opts ConnectionOptions) error {
	client, err := Connect(opts)
	if err != nil {
		fmt.Printf("Error1: %s\n", err.Error())
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		fmt.Printf("Error2: %s\n", err.Error())
		return err
	}
	defer session.Close()

	output, err := session.CombinedOutput("dir")
	fmt.Printf("%s\n", string(output))
	if err != nil {
		fmt.Printf("Error3: %s\n", err.Error())
		return err
	}

	return nil
}
