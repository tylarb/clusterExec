package clusterexec

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/sync/errgroup"
)

/* clusterexec takes hostnames and bash commands as input.
From there, it forms ssh connections using golang's ssh package,
then executes bash commands using concurrency.

*/

// SSHCluster is the base cluster struct.
// It forms a grouping of a host or hosts on which a list of commands will be executed
type SSHCluster struct {
	Nodes          []*ClusterNode
	GlobalTimeout  time.Duration
	CommandTimeout time.Duration // Option - filepath to make for easy scanning of keys?
	Port           int
	User           string
	KnownHostsFile string
	HostKeyCheck   bool
	Auth           []ssh.AuthMethod
	wg             *sync.WaitGroup
	//	EnsureExecute       bool  //WIP may be added later
	//	ExecuteConcurrently bool  // WIP may be added later
}

//ClusterCmdOut is the output of a command run on a cluster
type ClusterCmdOut struct {
	Node           *ClusterNode
	Err            error
	Stdout, Stderr string
}

// CreateCluster generates a pointer to a new SSH cluster
func CreateCluster(user string, hosts []string, options ...ClusterOption) (*SSHCluster, error) {
	cluster := SSHCluster{Port: PORT, wg: &sync.WaitGroup{}}
	for _, opt := range options {
		opt(&cluster)
	}

	nodeOptions := composeNodeOptions(&cluster)

	nodes := make([]*ClusterNode, len(hosts))

	var err error
	for i, host := range hosts { //TODO handle this error better
		nodes[i], err = CreateNode(user, host, nodeOptions)
		if err != nil {
			return nil, err
		}
	}
	cluster.Nodes = append(cluster.Nodes, nodes...)
	return &cluster, nil
}

func composeNodeOptions(cluster *SSHCluster) NodeOption {
	return NodeOptionCompose(NodeOptionAuthMethod(cluster.Auth...), NodeOptionHostKeyCheck(cluster.HostKeyCheck), NodeOptionKnownHostsFile(cluster.KnownHostsFile), NodeOptionPort(cluster.Port))
}

// Dial connects to all nodes on an ssh cluster
func (cluster *SSHCluster) Dial() error {
	var g errgroup.Group
	for _, node := range cluster.Nodes {
		node := node
		f := func() error {
			err := node.Dial()
			return err
		}
		g.Go(f)
	}
	return g.Wait()
}

// Close closes connections on an ssh cluster
func (cluster *SSHCluster) Close() error {
	var g errgroup.Group
	for _, node := range cluster.Nodes {
		node := node
		f := func() error {
			err := node.Close()
			return err
		}
		g.Go(f)
	}
	return g.Wait()
}

// Run executes bash commands on the SSH cluster   // TDOD - clean up errors, add an errors channel to get specific nodes w/ issue?
func (cluster *SSHCluster) Run(command *ClusterCmd) ([]ClusterCmdOut, error) {
	n := len(cluster.Nodes)
	cluster.wg.Add(n)
	var numErrors int
	clusterOut := make([]ClusterCmdOut, n)
	nodeOut := make(chan ClusterCmdOut)
	for _, node := range cluster.Nodes {
		go func(node *ClusterNode) {
			defer cluster.wg.Done()
			out := ClusterCmdOut{Node: node}
			out.Stdout, out.Stderr, out.Err = node.Run(command)
			nodeOut <- out
		}(node)
	}
	go func() {
		var i int
		for out := range nodeOut {
			if out.Err != nil {
				numErrors++
			}
			clusterOut[i] = out
			i++
		}
	}()
	cluster.wg.Wait()
	var err error
	if numErrors > 0 {
		err = fmt.Errorf("clusterExec: %d errors executing commands", numErrors)
	}
	return clusterOut, err
}
