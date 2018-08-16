package clusterExec

import (
	"time"

	"golang.org/x/crypto/ssh"
)

/* clusterExec takes hostnames and bash commands as input.
From there, it forms ssh connections using golang's ssh package,
then executes bash commands using concurrency.

*/

type SSHCluster struct {
	Hosts               []string
	Port                int
	SSHConfig           *ssh.ClientConfig
	Timeout             time.Time
	EnsureExecute       bool
	ExecuteConcurrently bool
}

type ConnectOption func(*SSHCluster)

func Connect(hosts []string, options ...ConnectOption) (*SSHCluster, error) {
	return nil, nil
}

func (cluster *SSHCluster) Run(commands []string) error {
	return nil
}

// TODO: Add Options...
// different id_rsa, etc file
// Different user name

func ConnectOptionPort(port int) ConnectOption {
	return func(cluster *SSHCluster) {
		cluster.Port = port
	}
}

func ConnectOptionTimeout(timeout time.Time) ConnectOption {
	return func(cluster *SSHCluster) {
		cluster.Timeout = timeout
	}
}

func ConnectOptionIgnoreExecute() ConnectOption {
	return func(cluster *SSHCluster) {
		cluster.EnsureExecute = false
	}
}

func ConnectOptionNotConcurrent() ConnectOption {
	return func(cluster *SSHCluster) {
		cluster.ExecuteConcurrently = false
	}
}

func ConnectOptionCompose(options ...ConnectOption) ConnectOption {
	return func(cluster *SSHCluster) {
		for _, opt := range options {
			opt(cluster)
		}
	}
}
