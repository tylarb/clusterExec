package clusterExec

import (
	"bytes"
	"os/exec"
	"sync"
	"time"
)

// ClusterCmd is a single command which can be passed to an individual node in the cluster
type ClusterCmd struct {
	Cmd         string
	Args        []string
	Timeout     time.Duration
	KillCommand chan bool
	CmdOut      chan ClusterOut
}

// ClusterOut is the response from ONE host for ONE command
type ClusterOut struct {
	Stdin  *bytes.Buffer
	Stdout *bytes.Buffer
	Stderr *bytes.Buffer
	Err    error
}

// CreateClusterCommand creates a new cluster command
func CreateClusterCommand(cmd string, args []string, options ...ClusterCmdOption) *ClusterCmd {
	clusterCommand := ClusterCmd{Cmd: cmd, Args: args, KillCommand: make(chan bool), CmdOut: make(chan ClusterOut)}
	for _, opt := range options {
		opt(&clusterCommand)
	}
	return &clusterCommand
}

// Run executes a cluster comman
func (clusterCommand *ClusterCmd) Run(wg *sync.WaitGroup) {
	var output ClusterOut
	if wg != nil {
		cmd := exec.Command(clusterCommand.Cmd, clusterCommand.Args...)
		cmd.Stdout = output.Stdout
	}
}
