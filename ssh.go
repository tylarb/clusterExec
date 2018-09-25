package clusterExec

import (
	"bufio"
	"errors"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

const (
	// PORT is the default ssh port
	PORT = 22
)

// ClusterNode contains information to collect to a single ssh node
type ClusterNode struct {
	Localhost      bool // true if host is local, to avoid using ssh
	User           string
	Hostname       string
	Port           int
	Auth           []ssh.AuthMethod
	Config         *ssh.ClientConfig
	HostKeyCheck   bool
	KnownHostsFile string
}

// CreateNode returns a single node for the cluster
func CreateNode(user, hostname string, options ...NodeOption) (*ClusterNode, error) {
	node := ClusterNode{User: user, Hostname: hostname, Port: PORT, HostKeyCheck: true}

	for _, opt := range options {
		opt(&node)
	}

	if node.Config == nil { // config isn't provided as a NodeOption, we'll compose it from other options
		if err := node.GetConfig(); err != nil {
			return nil, err
		}
	}

	return &node, nil
}

// GetConfig generates an ssh config given an auth method, user, and (if applicable) known hosts file
func (node *ClusterNode) GetConfig() error {
	var config ssh.ClientConfig

	if len(node.Auth) == 0 {
		panic("programming error: no auth method provided, cannot compose ssh config")
	}

	config.User = node.User
	config.Auth = node.Auth

	if node.HostKeyCheck {
		if node.KnownHostsFile == "" {
			panic("programming error: no known hosts file name provided")
		}
		var err error
		config.HostKeyCallback, err = knownhosts.New(node.KnownHostsFile)
		if err != nil {
			return err
		}
	} else {
		config.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	}

	node.Config = &config
	return nil

}

func parseHostKeys(hostname, keyfile string) (ssh.HostKeyCallback, error) {
	file, err := os.Open(keyfile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var hostKey ssh.PublicKey
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) != 3 {
			continue
		}
		if strings.Contains(fields[0], hostname) {
			var err error
			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
			if err != nil {
				return nil, err
			}
			break
		}
	}
	if hostKey == nil {
		return nil, errors.New("clusterExec: No key found for this host - make sure known_hosts file is valid")
	}
	return ssh.FixedHostKey(hostKey), nil

}

// NodeOption is a function which changes settings on a cluster node
type NodeOption func(*ClusterNode)

// NodeOptionPort allows you to set the ssh port
func NodeOptionPort(port int) NodeOption {
	return func(node *ClusterNode) {
		node.Port = port
	}
}

// NodeOptionIsLocalhost indicates the node is a localhost (avoids using ssh)
func NodeOptionIsLocalhost() NodeOption {
	return func(node *ClusterNode) {
		node.Localhost = true
	}
}

// NodeOptionAuthMethod adds an ssh authentication method to a node
func NodeOptionAuthMethod(auth ssh.AuthMethod) NodeOption {
	return func(node *ClusterNode) {
		node.Auth = append(node.Auth, auth)
	}
}

// NodeOptionKnownHostsFile sets the filename for the known hosts file (will be ignored if "HostKeyCheck" is false
func NodeOptionKnownHostsFile(file string) NodeOption {
	// if file == ""  return err if no file
	return func(node *ClusterNode) {
		node.KnownHostsFile = file
	}
}

// NodeOptionConfig provdes a full ssh.ClientConfig to a node
func NodeOptionConfig(config *ssh.ClientConfig) NodeOption {
	return func(node *ClusterNode) {
		node.Config = config
	}
}

// NodeOptionHostKeyCheck sets if a known host key check will be to verify the ssh connection
func NodeOptionHostKeyCheck(check bool) NodeOption {
	return func(node *ClusterNode) {
		node.HostKeyCheck = check
	}
}
