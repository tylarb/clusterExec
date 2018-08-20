package clusterExec

import (
	"os/user"
	"time"

	"golang.org/x/crypto/ssh"
)

/* clusterExec takes hostnames and bash commands as input.
From there, it forms ssh connections using golang's ssh package,
then executes bash commands using concurrency.

*/

const (
	// Default ssh port
	PORT = 22
)

// SSHCluster is the base cluster struct.
// It forms a grouping of a host or hosts on which a list of commands will be executed
type SSHCluster struct {
	Hosts               []string
	Port                int
	SSHConfig           *ssh.ClientConfig //TODO unsure if this should be in the cluster config
	GlobalTimeout       time.Time
	CommandTimeout      time.Time
	EnsureExecute       bool
	ExecuteConcurrently bool
	Username            string
}

// ConnectOption allows for functional API options to be added to an SSH Cluster
type ConnectOption func(*SSHCluster)

// CreateCluster generates a pointer to a new SSH cluster
func CreateCluster(hosts []string, options ...ConnectOption) (*SSHCluster, error) {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	cluster := SSHCluster{Hosts: hosts, Port: PORT, EnsureExecute: true, ExecuteConcurrently: true, Username: user.Username}
	for _, opt := range options {
		opt(&cluster)
	}
	return &cluster, nil
}

// Exec runs bash commands on the SSH cluster
func (cluster *SSHCluster) Exec(commands []string) error {
	return nil
}
