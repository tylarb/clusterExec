package clusterexec

import (
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

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

// ClusterOptionAddNodes adds one or more fully composed nodes to a cluster
func ClusterOptionAddNodes(nodes ...*ClusterNode) ClusterOption {
	return func(cluster *SSHCluster) {
		for _, node := range nodes {
			cluster.Nodes = append(cluster.Nodes, node)
		}
	}
}

// ClusterOptionPort sets a common ssh port for all nodes in cluster
func ClusterOptionPort(port int) ClusterOption {
	return func(cluster *SSHCluster) {
		cluster.Port = port
	}
}

// ClusterOptionUser sets a common user for all nodes in cluster
func ClusterOptionUser(user string) ClusterOption {
	return func(cluster *SSHCluster) {
		cluster.User = user
	}
}

// ClusterOptionIncludeLocalhost adds localhost to the cluster (without using ssh)
func ClusterOptionIncludeLocalhost() ClusterOption {
	return func(cluster *SSHCluster) {
		h, _ := os.Hostname()
		cluster.Nodes = append(cluster.Nodes, &ClusterNode{Localhost: true, Hostname: h})
	}
}

// ClusterOptionKnownHostsFile adds a common known hosts file to the cluster
func ClusterOptionKnownHostsFile(file string) ClusterOption {
	return func(cluster *SSHCluster) {
		cluster.KnownHostsFile = file
	}
}

// ClusterOptionHostKeyCheck sets if a host check will be used on the cluster to verify ssh connection
func ClusterOptionHostKeyCheck(check bool) ClusterOption {
	return func(cluster *SSHCluster) {
		cluster.HostKeyCheck = check
	}
}

// ClusterOptionAuthMethod adds a common auth method to the clsuter
func ClusterOptionAuthMethod(auth ssh.AuthMethod) ClusterOption {
	return func(cluster *SSHCluster) {
		cluster.Auth = append(cluster.Auth, auth)
	}
}

// ClusterOptionCompose condenses several connect options into one
func ClusterOptionCompose(options ...ClusterOption) ClusterOption {
	return func(cluster *SSHCluster) {
		for _, opt := range options {
			opt(cluster)
		}
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
