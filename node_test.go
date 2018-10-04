package clusterExec

import (
	"fmt"
	"testing"

	"golang.org/x/crypto/ssh"
)

func TestCreateNode(t *testing.T) {

	config := &ssh.ClientConfig{
		User: USER,
		Auth: []ssh.AuthMethod{
			ssh.Password(cluster22.password),
		},
	}
	node22, err := CreateNode(USER, cluster22.node0, NodeOptionConfig(config))
	if err != nil {
		t.Error(err)
	}
	if node22.Config != config {
		t.Log("Failed to set config for node1")
		t.Fail()
	}
	if node22.Hostname != cluster22.node0 {
		t.Log("Failed to set hostname for node1")
		t.Fail()
	}
	if node22.User != USER {
		t.Log("Failed to set USER for node1")
		t.Fail()
	}

	node25, err := CreateNode(USER, cluster25.node0, NodeOptionPort(cluster25.port), NodeOptionAuthMethod(ssh.Password(cluster25.password)), NodeOptionKnownHostsFile(dir+"/known_hosts"))
	if err != nil {
		t.Error(err)
	}
	if node25.Port != cluster25.port {
		t.Log("Failed to set port for node2")
		t.Fail()
	}
	if node25.KnownHostsFile != dir+"/known_hosts" {
		t.Log("Failed to set known_hosts for node2")
		t.Fail()
	}
	t.Log("Created node with options")
}

func TestGetConfig(t *testing.T) {

	node, err := CreateNode(USER, cluster25.node0, NodeOptionPort(cluster25.port), NodeOptionAuthMethod(ssh.Password(cluster25.password)), NodeOptionKnownHostsFile(dir+"/known_hosts"))
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
		t.Logf("hostaddress: %s", hostaddress)
		t.Log(err)
		t.Fail()
	} else {
		t.Log("Created a config")
	}
	client.Close()
}
