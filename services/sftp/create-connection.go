package sftp

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

func (s SSH) dial() (*ssh.Client, error) {
	if s.PathToSSHKnownHostsFile == "" {
		s.PathToSSHKnownHostsFile = filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts")
	}
	hostKeyCallback, err := knownhosts.New(s.PathToSSHKnownHostsFile)
	if err != nil {
		return nil, err
	}
	addr := fmt.Sprintf(s.Host+":%v", s.Port)
	// connect
	return ssh.Dial("tcp", addr, &ssh.ClientConfig{
		User: s.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.Password),
		},
		// HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		HostKeyCallback: hostKeyCallback,
	})
}
