package clusterExec

import (
	"bytes"
	"fmt"
	"os/exec"
	"time"

	"golang.org/x/crypto/ssh"
)

// ClusterCmd is a single command which can be passed to an individual node in the cluster
type ClusterCmd struct {
	Cmd     string
	Args    []string
	Timeout time.Duration
}

//ClusterCmdOut is the output of a command run on a cluster
/* type ClusterCmdOut struct {
	err            error
	Stdout, Stderr *bytes.Buffer
} */

// CreateClusterCommand creates a new cluster command
func CreateClusterCommand(cmd string, args []string, options ...ClusterCmdOption) *ClusterCmd {
	clusterCommand := ClusterCmd{Cmd: cmd, Args: args}
	for _, opt := range options {
		opt(&clusterCommand)
	}
	return &clusterCommand
}

// Run executes a command on a cluster node. // TODO: clarify errors returned
func (node *ClusterNode) Run(command *ClusterCmd) (stdOut, stdErr string, err error) {
	if node.Localhost {
		stdOut, stdErr, err = node.runLocalCommand(command)
	} else {
		stdOut, stdErr, err = node.runRemoteCommand(command)
	}
	return "", "", err
}

// runs command locally if localhost
func (node *ClusterNode) runLocalCommand(command *ClusterCmd) (stdOut, stdErr string, err error) {
	retError := make(chan error)
	timeout := make(chan bool)

	var stdOutBuff, stdErrBuff bytes.Buffer
	cmd := exec.Command(command.Cmd, command.Args...)
	cmd.Stdout = &stdOutBuff
	cmd.Stderr = &stdErrBuff
	if err := cmd.Start(); err != nil {
		return "", "", &CommandExecutionError{fmt.Sprintf("Command failed to start"), node, command, err}
	}
	go func() {
		retError <- cmd.Wait()
	}()
	if command.Timeout > 0 {
		go func(t time.Duration) {
			time.Sleep(t)
			timeout <- true
		}(command.Timeout)
	}

	select {
	case err := <-retError:

		if err != nil {
			return stdOutBuff.String(), stdErrBuff.String(), &CommandExecutionError{fmt.Sprintf("Command returned failed"), node, command, err}
		}
		return stdOutBuff.String(), stdErrBuff.String(), nil
	case <-timeout:
		cmd.Process.Kill()
		return "", "", &CommandTimeoutError{fmt.Sprintf("timeout after %s", command.Timeout), node, command}
	}
}

// runs commands remotely over ssh
func (node *ClusterNode) runRemoteCommand(command *ClusterCmd) (stdOut, stdErr string, err error) {
	retError := make(chan error)
	timeout := make(chan bool)
	if node.Client == nil {
		return "", "", &NodeConnectionError{"No existing ssh connection", node}
	}
	session, err := node.Client.NewSession()
	if err != nil {
		return "", "", err
	}
	defer session.Close()
	var stdOutBuff, stdErrBuff bytes.Buffer

	session.Stdout = &stdOutBuff
	session.Stderr = &stdErrBuff
	cmdString := composeCmd(command.Cmd, command.Args)
	if err = session.Start(cmdString); err != nil {
		return "", "", &CommandExecutionError{fmt.Sprintf("Command failed to start"), node, command, err} // Return err : command could not start with err
	}
	go func() {
		retError <- session.Wait()
	}()
	if command.Timeout > 0 {
		go func(t time.Duration) {
			time.Sleep(t)
			timeout <- true
		}(command.Timeout)
	}
	select {
	case err := <-retError:

		if err != nil {
			return stdOutBuff.String(), stdErrBuff.String(), &CommandExecutionError{fmt.Sprintf("Command returned failed"), node, command, err}
		}
		return stdOutBuff.String(), stdErrBuff.String(), nil
	case <-timeout:
		session.Signal(ssh.SIGHUP) // Terminate child process - this DOES NOT WORK
		// appars to be a limitation in ssh, see https://github.com/golang/go/issues/16597
		// The child process will still be reaped as the session will be closed, so long
		// as the child process responds to such things. `sleep` and `dd` don't
		return "", "", &CommandTimeoutError{fmt.Sprintf("timeout after %s", command.Timeout), node, command}
	}

}

// CommandTimeoutError is returned when an executing command times out rather than completing cleanly
// The relevant command will be terminated with SIGHUP
type CommandTimeoutError struct {
	err     string
	Node    *ClusterNode
	Command *ClusterCmd
}

func (t *CommandTimeoutError) Error() string {
	return fmt.Sprintf("clusterExec: node address %s: command %s: %s", t.Node.Addr, t.Command.Cmd, t.err)
}

// NodeConnectionError is returned if a node does not have an existing client connection when
// a command is attempted to run, or has some other networking error
type NodeConnectionError struct {
	err  string
	Node *ClusterNode
}

func (n *NodeConnectionError) Error() string {
	return fmt.Sprintf("clusterExec: node address %s: %s", n.Node.Addr, n.err)
}

// CommandExecutionError is returned when there is some problem starting or executing a command.
// executionError is the full error the command returns
type CommandExecutionError struct {
	err            string
	Node           *ClusterNode
	Command        *ClusterCmd
	ExecutionError error
}

func (c *CommandExecutionError) Error() string {
	return fmt.Sprintf("clusterExec: node address %s: Command: %s Command Execution failed %s", c.Node.Addr, c.Command.Cmd, c.err)
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
