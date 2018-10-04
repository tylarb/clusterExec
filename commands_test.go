package clusterExec

import (
	"testing"
	"time"
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
	t.Log("creating cluster command complete")
}

func TestRunRemoteCommand(t *testing.T) {

	command := "hostname"
	args := []string{"-f"}

	clusterCmd := CreateClusterCommand(command, args)

	stdOut, stdErr, err := sshNode.runRemoteCommand(clusterCmd)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	if stdOut != "ssh0.cluster22\n" || stdErr != "" {
		t.Logf("Incorrect output - got %s, %s; expected %s, %s", stdOut, stdErr, "ssh0.cluster22\n", "")
		t.Fail()
	} else {
		t.Logf("ran command %s with output %s", command, stdOut)
	}
}

func TestRunRemoteCommandTimeout(t *testing.T) {

	timeout := time.Second * 5

	command := "cat"
	args := []string{"/dev/random"}
	clusterCmd := CreateClusterCommand(command, args, ClusterCmdOptionTimeout(timeout))
	stdOut, stdErr, err := sshNode.runRemoteCommand(clusterCmd)
	if T, ok := err.(*CommandTimeoutError); !ok {
		t.Logf("expected timeout error, but instead received type %v", T)
		t.Fail()
	} else {
		t.Logf("ran command %s with timeout of %s", composeCmd(command, args), timeout)
		t.Logf("stdout: %s, stderr: %s", stdOut, stdErr)
	}
}

func TestRemoteCommandFailedStart(t *testing.T) {
	command := "NotACommand"
	args := []string{"not", "an", "arg"}

	clusterCmd := CreateClusterCommand(command, args)

	_, _, err := sshNode.runRemoteCommand(clusterCmd)
	if T, ok := err.(*CommandExecutionError); !ok {
		t.Logf("Expected execution error, instead recieved type %v", T)
		t.Fail()
	} else {
		t.Logf("ran command %s and correctly received err %s", composeCmd(command, args), err)
	}
}

func TestLocalCommand(t *testing.T) {
	command := "echo"
	args := []string{"This is a test"}

	clusterCmd := CreateClusterCommand(command, args)

	stdOut, stdErr, err := localhostNode.runLocalCommand(clusterCmd)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	if stdOut != "This is a test\n" || stdErr != "" {
		t.Logf("Incorrect output - got %s, %s; expected %s, %s", stdOut, stdErr, "This is a test\n", "")
		t.Fail()
	} else {
		t.Logf("ran command %s with output %s", command, stdOut)
	}

}

func TestLocalCommandTimeout(t *testing.T) {

	timeout := time.Second * 5

	command := "cat"
	args := []string{"/dev/random"}
	clusterCmd := CreateClusterCommand(command, args, ClusterCmdOptionTimeout(timeout))
	stdOut, stdErr, err := localhostNode.runLocalCommand(clusterCmd)
	if T, ok := err.(*CommandTimeoutError); !ok {
		t.Logf("expected timeout error, but instead received type %v", T)
		t.Fail()
	} else {
		t.Logf("ran command %s with timeout of %s", composeCmd(command, args), timeout)
		t.Logf("stdout: %s, stderr: %s", stdOut, stdErr)
	}
}

func TestLocalCommandFailedStart(t *testing.T) {
	command := "NotACommand"
	args := []string{"not", "an", "arg"}

	clusterCmd := CreateClusterCommand(command, args)

	_, _, err := localhostNode.runLocalCommand(clusterCmd)
	if T, ok := err.(*CommandExecutionError); !ok {
		t.Logf("Expected execution error, instead recieved type %v", T)
		t.Fail()
	} else {
		t.Logf("ran command %s and correctly received err %s", composeCmd(command, args), err)
	}

}
