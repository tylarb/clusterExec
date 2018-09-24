package clusterExec

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/crypto/ssh"
)

const USER = "root"

func TestCreateNode(t *testing.T) {

	currentDir, err := getDir()
	if err != nil {
		t.Error(err)
	}
	config := &ssh.ClientConfig{
		User: USER,
		Auth: []ssh.AuthMethod{
			ssh.Password("password"),
		},
	}
	node1, err := CreateNode(USER, "172.22.0.10", NodeOptionConfig(config))
	if err != nil {
		t.Error(err)
	}
	if node1.Config != config {
		t.Log("Failed to set config for node1")
		t.Fail()
	}
	if node1.Hostname != "172.22.0.10" {
		t.Log("Failed to set hostname for node1")
		t.Fail()
	}
	if node1.User != USER {
		t.Log("Failed to set user for node1")
		t.Fail()
	}

	node2, err := CreateNode(USER, "172.25.0.10", NodeOptionPort(25), NodeOptionAuthMethod(ssh.Password("password")), NodeOptionKnownHostsFile(currentDir+"/known_hosts"))
	if err != nil {
		t.Error(err)
	}
	if node2.Port != 25 {
		t.Log("Failed to set port for node2")
		t.Fail()
	}
	if node2.KnownHostsFile != currentDir+"known_hosts" {
		t.Log("Failed to set known_hosts for node2")
		t.Fail()
	}

}

func TestGetConfig(t *testing.T) {

	currentDir, err := getDir()
	if err != nil {
		t.Error(err)
	}

	node, err := CreateNode(USER, "172.25.0.10", NodeOptionPort(25), NodeOptionAuthMethod(ssh.Password("password")), NodeOptionKnownHostsFile(currentDir+"/known_hosts"))
	if err != nil {
		t.Error(err)
	}
	if err := node.GetConfig(); err != nil {
		t.Log(err)
		t.Fail()
	}
	config := &ssh.ClientConfig{
		User: USER,
		Auth: []ssh.AuthMethod{
			ssh.Password("password"),
		},
	}
	config.HostKeyCallback, err = parseHostKeys("172.25.0.10", currentDir+"/known_hosts")
	if node.Config != config {
		t.Log("Failed to set config correctly")
		t.Fail()
	}
}

func getDir() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	currentDir := filepath.Dir(ex)
	return currentDir, nil
}
