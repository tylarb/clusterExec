package clusterExec

import "time"

// ClusterOption allows for functional API options to be added to an SSH Cluster
type ClusterOption func(*SSHCluster)

// TODO: Add Options...
// different id_rsa, etc file
// Different user name

// ClusterOptionGlobalTimeout sets the timeout for executing all provided commands
func ClusterOptionGlobalTimeout(timeout time.Duration) ClusterOption {
	return func(cluster *SSHCluster) {
		cluster.GlobalTimeout = timeout
	}
}

// ClusterOptionCommandTimeout sets the timeout for any individual command.
// Global timeout superceeds this setting
func ClusterOptionCommandTimeout(timeout time.Duration) ClusterOption {
	return func(cluster *SSHCluster) {
		cluster.CommandTimeout = timeout
	}
}

// ConnectOptionIgnoreExecuteSuccess indiates that the command should be
// run without ensuring completion WIP maybe added later
/*
func ConnectOptionIgnoreExecuteSuccess() ConnectOption {
	return func(cluster *SSHCluster) {
		cluster.EnsureExecute = false
	}
}
*/
// ConnectOptionNotConcurrent indicates executed commands should not be run
// concurrently on the destination system - each command must complete one
// at a time. WIP may be added later
/*
func ConnectOptionNotConcurrent() ConnectOption {
	return func(cluster *SSHCluster) {
		cluster.ExecuteConcurrently = false
	}
}
*/

// ClusterOptionCompose condenses several connect options into one
func ClusterOptionCompose(options ...ClusterOption) ClusterOption {
	return func(cluster *SSHCluster) {
		for _, opt := range options {
			opt(cluster)
		}
	}
}
