package clusterExec

import (
	"fmt"
	"io/ioutil"

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
	Addr           string
	Auth           []ssh.AuthMethod
	Config         *ssh.ClientConfig
	HostKeyCheck   bool
	KnownHostsFile string
	Client         *ssh.Client
}

// CreateNode returns a single node for the cluster
func CreateNode(user, hostname string, options ...NodeOption) (*ClusterNode, error) {
	node := ClusterNode{User: user, Hostname: hostname, Port: PORT, HostKeyCheck: true}

	for _, opt := range options {
		opt(&node)
	}

	if node.Config == nil && !node.Localhost { // config isn't provided as a NodeOption, we'll compose it from other options
		if err := node.GetConfig(); err != nil {
			return nil, err
		}
	}
	node.Addr = fmt.Sprintf("%s:%d", node.Hostname, node.Port)

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

// Dial creates an ssh connection on a cluster node. node.Close() must be called to close the networking connection
// Repeated calls to Dial return silently if there is alreay a connection present
func (node *ClusterNode) Dial() error {
	if node.Client != nil || node.Localhost { // FIXME: is there some ssh ping to ensure a valid connection?
		return nil
	}
	var err error
	node.Client, err = ssh.Dial("tcp", node.Addr, node.Config)
	return err
}

// Close closes the networking conection to a node. Multiple calls to Close return silently.
func (node *ClusterNode) Close() error {
	if node.Client == nil || node.Localhost {
		return nil
	}
	err := node.Client.Close()
	node.Client = nil
	return err

}

// GetPrivateKeyAuth returns an ssh PublicKeys auth method provided a unencrypted private key file name
func GetPrivateKeyAuth(file string) (ssh.AuthMethod, error) {
	var auth ssh.AuthMethod
	buff, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	key, err := ssh.ParsePrivateKey(buff)
	if err != nil {
		return nil, err
	}
	auth = ssh.PublicKeys(key)
	return auth, nil
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

// NodeOptionCompose composes several options into one
func NodeOptionCompose(options ...NodeOption) NodeOption {
	return func(node *ClusterNode) {
		for _, opt := range options {
			opt(node)
		}
	}
}
