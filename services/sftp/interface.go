package sftp

import "golang.org/x/crypto/ssh"

type SFTP interface {
	dial() (*ssh.Client, error)
}

type SSH struct {
	User                    string // Username of remote ssh
	Password                string // Password of remote ssh
	Host                    string // Host URL of remote ssh
	Port                    int    // := ":22"
	PathToSSHKnownHostsFile string
}
