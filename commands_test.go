package clusterExec

import (
	"testing"
	"time"

	"golang.org/x/crypto/ssh"
)

func TestCreateClusterCommand(t *testing.T) {
	timeout := 5 * time.Second
	command := "hostname"
	args := []string{"-f"}

	clusterCmd := CreateClusterCommand(command, args, ClusterCmdOptionTimeout(timeout))

	if clusterCmd.Cmd != command {
		t.Log("Command not set")
		t.Fail()
	}
	for i, arg := range clusterCmd.Args {
		if arg != args[i] {
			t.Log("args not set correctly")
			t.Fail()
		}
	}

	if clusterCmd.Timeout != timeout {
		t.Log("Timeout not set")
		t.Fail()
	}
}

func TestRunRemoteCommand(t *testing.T) {
	node, err := CreateNode(user, cluster22.node0, NodeOptionAuthMethod(ssh.Password(cluster22.password)), NodeOptionHostKeyCheck(false))
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	command := "hostname"
	args := []string{"-f"}

	clusterCmd := CreateClusterCommand(command, args)

	stdOut, stdErr, err := node.runRemoteCommand(clusterCmd)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	if stdOut != "ssh0.cluster22" || stdErr != "" {
		t.Logf("Incorrect output - got %s, %s; expected %s, %s", stdOut, stdErr, "ssh0.cluster22", "")
		t.Fail()
	}

}

func TestRunRemoteCommandTimeout(t *testing.T) {
	node, err := CreateNode(user, cluster22.node0, NodeOptionAuthMethod(ssh.Password(cluster22.password)), NodeOptionHostKeyCheck(false))
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	timeout := time.Second * 5

	command := "cat"
	args := []string{"/dev/random"}
	clusterCmd := CreateClusterCommand(command, args, ClusterCmdOptionTimeout(timeout))
	_, _, err = node.runRemoteCommand(clusterCmd)
	if T, ok := err.(CommandTimeoutError); !ok {
		t.Logf("expected timeout error, but instead received type %v", T)
		t.Fail()
	}

}
