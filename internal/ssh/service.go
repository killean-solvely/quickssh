package ssh

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

type SSHService struct {
	hostKeyApprovals *hostKeyApprovals
}

func NewSSHService() *SSHService {
	return &SSHService{
		hostKeyApprovals: &hostKeyApprovals{},
	}
}
