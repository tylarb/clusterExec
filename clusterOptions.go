package clusterExec

import "time"

// ConnectOption allows for functional API options to be added to an SSH Cluster
type ConnectOption func(*SSHCluster)

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
func ConnectOptionGlobalTimeout(timeout time.Duration) ConnectOption {
	return func(cluster *SSHCluster) {
		cluster.GlobalTimeout = timeout
	}
}

// ConnectOptionCommandTimeout sets the timeout for any individual command.
// Global timeout superceeds this setting
func ConnectOptionCommandTimeout(timeout time.Duration) ConnectOption {
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

//ConnectOptionUsername changes the connecting username (default is os.CurrentUser)
func ConnectOptionUsername(name string) ConnectOption {
	return func(cluster *SSHCluster) {
		cluster.Username = name
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
