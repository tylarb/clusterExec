package clusterExec

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

// ClusterCmd is a single command which can be passed to an individual node in the cluster
type ClusterCmd struct {
	Cmd            string
	Args           []string
	Timeout        time.Duration
	Stdout, Stderr *bytes.Buffer
}

// CreateClusterCommand creates a new cluster command
func CreateClusterCommand(cmd string, args []string, options ...ClusterCmdOption) *ClusterCmd {
	clusterCommand := ClusterCmd{Cmd: cmd, Args: args, Stdout: new(bytes.Buffer), Stderr: new(bytes.Buffer)}
	for _, opt := range options {
		opt(&clusterCommand)
	}
	return &clusterCommand
}

// Run executes a command on a cluster node.
func (node *ClusterNode) Run(command *ClusterCmd) error {
	var err error
	if node.Localhost {
		err = node.runLocalCommand(command)
	} else {
		err = node.runRemoteCommand(command)
	}
	return err
}

// runs command locally if localhost
func (node *ClusterNode) runLocalCommand(command *ClusterCmd) error {

	return nil
}

// runs commands remotely over ssh
func (node *ClusterNode) runRemoteCommand(command *ClusterCmd) error {

	if node.Client == nil {
		return &NodeConnectionError{"Existing ssh connection", node}
	}
	session, err := node.Client.NewSession()
	if err != nil {
		return err
	}
	session.Stdout = command.Stdout
	session.Stderr = command.Stderr
	cmdString := composeCmd(command.Cmd, command.Args)
	err = session.Start(cmdString)
	if err != nil {
		return err
	}
	go func() {
		session.Wait()
	}()

	return nil
}
func composeCmd(cmd string, args []string) string {
	command := cmd + " " + strings.Join(args, " ")
	return command
}

// NodeConnectionError is returned if a node does not have an existing client connection when
// a command is attempted to run, or has some other networking error
type NodeConnectionError struct {
	err  string
	node *ClusterNode
}

func (n NodeConnectionError) Error() string {
	return fmt.Sprintf("clusterExec: node address %s: %s", n.node.Addr, n.err)
}

// CommandsRun runs an array of cluster commands via goroutines
/*
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


*/
