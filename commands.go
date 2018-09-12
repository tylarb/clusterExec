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
	CmdOut      ClusterOut
}

// ClusterOut is the response from ONE host for ONE command
type ClusterOut struct {
	done chan bool
	Stdin  *bytes.Buffer
	Stdout *bytes.Buffer
	Stderr *bytes.Buffer
	Err    error
}

// CreateClusterCommand creates a new cluster command
func CreateClusterCommand(cmd string, args []string, options ...ClusterCmdOption) *ClusterCmd {
	clusterCommand := ClusterCmd{Cmd: cmd, Args: args, KillCommand: make(chan bool)}
	for _, opt := range options {
		opt(&clusterCommand)
	}
	return &clusterCommand
}

// Run executes a cluster comman
func (clusterCommand *ClusterCmd) Run(wg *sync.WaitGroup, done chan bool) {
	if wg != nil { // perhaps this should be if wg == nil {return}
		defer wg.Done()
		cmd := exec.Command(clusterCommand.Cmd, clusterCommand.Args...)
		cmd.Stdout = clusterCommand.CmdOut.Stdout
		cmd.Stderr = clusterCommand.CmdOut.Stderr
		cmd.Stdin = clusterCommand.CmdOut.Stdin

		clusterCommand.CmdOut.Err = cmd.Run()

	}
}

// CommandsRun runs an array of cluster commands via goroutines
func CommandsRun(clusterCmds []ClusterCmd) ([]ClusterCmd, error) {
	var wg *sync.WaitGroup
	wg.Add(len(commands))
	go func(){
		defer wg.Wait()
	}
	for _,clusterCmd := range commands {
		go clusterCmd.Run(wg)
	}
	select {
		
	}
}
