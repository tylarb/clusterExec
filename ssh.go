package clusterExec

import "golang.org/x/crypto/ssh"

// ClusterNode contains information to collect to a single ssh node
type ClusterNode struct {
	localhost bool // true if host is local, to avoid using ssh
	hostname  string
	port      int
	config    *ssh.ClientConfig
}
