package clusterExec

import (
	"time"

	"golang.org/x/crypto/ssh"
)

/* clusterExec takes hostnames and bash commands as input.
From there, it forms ssh connections using golang's ssh package,
then executes bash commands using concurrency.

*/

// SSHCluster is the ssh cluster struct
type SSHCluster struct {
	Hosts               []string
	Port                int
	SSHConfig           *ssh.ClientConfig
	GlobalTimeout       time.Time
	CommandTimeout      time.Time
	EnsureExecute       bool
	ExecuteConcurrently bool
}

// ConnectOption allows for functional API options to be added to an SSH Cluster
type ConnectOption func(*SSHCluster)

// Connect generates a pointer to a new SSH cluster
func Connect(hosts []string, options ...ConnectOption) (*SSHCluster, error) {
	return nil, nil
}

// Run runs bash commands on the SSH cluster
func (cluster *SSHCluster) Run(commands []string) error {
	return nil
}

// TODO: Add Options...
// different id_rsa, etc file
// Different user name

// ConnectOptionPort allows you to set the ssh port to connect to
func ConnectOptionPort(port int) ConnectOption {
	return func(cluster *SSHCluster) {
		cluster.Port = port
	}
}

// ConnectOptionGlobalTimeout sets the timeout for executing all provided commands
func ConnectOptionGlobalTimeout(timeout time.Time) ConnectOption {
	return func(cluster *SSHCluster) {
		cluster.GlobalTimeout = timeout
	}
}

// ConnectOptionCommandTimeout sets the timeout for any individual command.
// Global timeout superceeds this setting
func ConnectOptionCommandTimeout(timeout time.Time) ConnectOption {
	return func(cluster *SSHCluster) {
		cluster.CommandTimeout = timeout
	}
}

// ConnectOptionIgnoreExecuteSuccess indiates that the command should be
// run without ensuring completion
func ConnectOptionIgnoreExecuteSuccess() ConnectOption {
	return func(cluster *SSHCluster) {
		cluster.EnsureExecute = false
	}
}

// ConnectOptionNotConcurrent indicates executed commands should not be run
// concurrently on the destination system - each command must complete one
// at a time.
func ConnectOptionNotConcurrent() ConnectOption {
	return func(cluster *SSHCluster) {
		cluster.ExecuteConcurrently = false
	}
}

// ConnectOptionCompose condenses several connect options into one
func ConnectOptionCompose(options ...ConnectOption) ConnectOption {
	return func(cluster *SSHCluster) {
		for _, opt := range options {
			opt(cluster)
		}
	}
}
