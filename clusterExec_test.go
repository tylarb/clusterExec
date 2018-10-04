package clusterExec

import (
	"os"
	"testing"
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

var dir string

func TestMain(t *testing.M) {
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

}
