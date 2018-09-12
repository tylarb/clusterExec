package clusterExec

import (
	"time"
)

/* clusterExec takes hostnames and bash commands as input.
From there, it forms ssh connections using golang's ssh package,
then executes bash commands using concurrency.

*/

// SSHCluster is the base cluster struct.
// It forms a grouping of a host or hosts on which a list of commands will be executed
type SSHCluster struct {
	Nodes          []*ClusterNode
	GlobalTimeout  time.Duration
	CommandTimeout time.Duration // Option - filepath to make for easy scanning of keys?
	//	EnsureExecute       bool  //WIP may be added later
	//	ExecuteConcurrently bool  // WIP may be added later
}

// CreateCluster generates a pointer to a new SSH cluster
func CreateCluster(nodes []*ClusterNode, options ...ClusterOption) (*SSHCluster, error) {

	cluster := SSHCluster{Nodes: nodes}
	for _, opt := range options {
		opt(&cluster)
	}
	return &cluster, nil
}

// Exec runs bash commands on the SSH cluster
func (cluster *SSHCluster) Exec(commands []ClusterCmd) error {
	return nil
}
