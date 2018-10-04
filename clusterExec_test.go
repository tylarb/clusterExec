package clusterExec

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"golang.org/x/crypto/ssh"
)

var user string
var cluster22 struct {
	node0, node1, node2, node3 string
	port                       int
	password                   string
}

var cluster25 struct {
	node0, node1, node2, node3 string
	port                       int
	password                   string
}

var clusterNode *ClusterNode

var dir string

func TestMain(m *testing.M) {
	user = os.Getenv("TEST_USER")
	if user == "" {
		user = "root"
	}
	dir = os.Getenv("TRAVIS_BUILD_DIR")
	if dir == "" {
		dir = os.Getenv("HOME") + "/.ssh"
	}

	cluster22.node0 = "172.22.0.10"
	cluster22.node1 = "172.22.0.11"
	cluster22.node2 = "172.22.0.12"
	cluster22.node3 = "172.22.0.13"
	cluster22.port = 22

	cluster25.node0 = "172.25.0.10"
	cluster25.node1 = "172.25.0.11"
	cluster25.node2 = "172.25.0.12"
	cluster25.node3 = "172.25.0.13"
	cluster25.port = 25

	cluster22.password = os.Getenv("22PASSWORD")
	if cluster22.password == "" {
		cluster22.password = "password"
	}
	cluster25.password = os.Getenv("25PASSWORD")
	if cluster25.password == "" {
		cluster25.password = "password"
	}
	var err error
	clusterNode, err = CreateNode(user, cluster22.node0, NodeOptionAuthMethod(ssh.Password(cluster22.password)), NodeOptionHostKeyCheck(false))
	if err != nil {
		fmt.Printf("Could not connect to ssh node with err %v", err)
		os.Exit(2)
	}
	if err := clusterNode.Dial(); err != nil {
		fmt.Printf("Could not connect to ssh node with err %v", err)
		os.Exit(2)
	}

	flag.Parse()
	exitCode := m.Run()

	clusterNode.Close()
	os.Exit(exitCode)
}
