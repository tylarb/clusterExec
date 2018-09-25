package clusterExec

import (
	"fmt"
	"os"
	"testing"

	"golang.org/x/crypto/ssh"
)

const USER = "root"

var dir = os.Getenv("HOME") + "/.ssh"

func TestCreateNode(t *testing.T) {

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

	node2, err := CreateNode(USER, "172.25.0.10", NodeOptionPort(25), NodeOptionAuthMethod(ssh.Password("password")), NodeOptionKnownHostsFile(dir+"/known_hosts"))
	if err != nil {
		t.Error(err)
	}
	if node2.Port != 25 {
		t.Log("Failed to set port for node2")
		t.Fail()
	}
	if node2.KnownHostsFile != dir+"/known_hosts" {
		t.Log("Failed to set known_hosts for node2")
		t.Fail()
	}

}

func TestGetConfig(t *testing.T) {

	node, err := CreateNode(USER, "172.25.0.10", NodeOptionPort(25), NodeOptionAuthMethod(ssh.Password("password")), NodeOptionKnownHostsFile(dir+"/known_hosts"))
	if err != nil {
		t.Error(err)
	}
	if err := node.GetConfig(); err != nil {
		t.Log(err)
		t.Fail()
	}

	hostaddress := fmt.Sprintf("%s:%d", node.Hostname, node.Port)
	client, err := ssh.Dial("tcp", hostaddress, node.Config)
	if err != nil {
		t.Log("Failed to generate valid config")
		t.Log(err)
		t.Fail()
	}
	client.Close()
}
