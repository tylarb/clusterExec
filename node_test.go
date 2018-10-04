package clusterExec

import (
	"fmt"
	"testing"

	"golang.org/x/crypto/ssh"
)

func TestCreateNode(t *testing.T) {

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(cluster22.password),
		},
	}
	node22, err := CreateNode(user, cluster22.node0, NodeOptionConfig(config))
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
	if node22.User != user {
		t.Log("Failed to set user for node1")
		t.Fail()
	}

	node25, err := CreateNode(user, cluster25.node0, NodeOptionPort(cluster25.port), NodeOptionAuthMethod(ssh.Password(cluster25.password)), NodeOptionKnownHostsFile(dir+"/known_hosts"))
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

}

func TestGetConfig(t *testing.T) {

	node, err := CreateNode(user, cluster25.node0, NodeOptionPort(cluster25.port), NodeOptionAuthMethod(ssh.Password(cluster25.password)), NodeOptionKnownHostsFile(dir+"/known_hosts"))
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
