package clusterExec

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"
)

const (
	// Default ssh port
	PORT = 22
)

// ClusterNode contains information to collect to a single ssh node
type ClusterNode struct {
	Localhost bool // true if host is local, to avoid using ssh
	User      string
	Hostname  string
	Port      int
	Auth      []ssh.AuthMethod
	Config    *ssh.ClientConfig
	Files     NodeFiles // connection files to access the node
}

// NodeFiles are the files used to connect to a cluster node - i.e. keys and known_hosts
type NodeFiles struct {
	PublicKeyFile  string
	KnownHostsFile string
}

// CreateNode returns a single node for the cluster
func CreateNode(user, hostname string, options ...NodeOption) (*ClusterNode, error) {
	node := ClusterNode{User: user, Hostname: hostname, Port: PORT}
	node.GetDefaultFiles()

	for _, opt := range options {
		opt(&node)
	}
	return &node, nil
}

func (node *ClusterNode) GetDefaultFiles() {
	node.Files.KnownHostsFile = filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts")
	node.Files.PublicKeyFile = filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa.pub") // FIXME other key types?
}

func (node *ClusterNode) GetConfig() error {
	var config ssh.ClientConfig
	// file, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	file, err := os.Open(node.Files.KnownHostsFile)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var hostKey ssh.PublicKey
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) != 3 {
			continue
		}
		if strings.Contains(fields[0], node.Hostname) {
			var err error
			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
			if err != nil {
				return err

			}
			break
		}
	}
	if hostKey == nil {
		return errors.New("key for this host is nil - make sure known_hosts is valid")
	}

	config.User = node.User
	config.HostKeyCallback = ssh.FixedHostKey(hostKey)
	config.Auth = node.Auth
	node.Config = &config
	return nil

}

// NodeOption is a function which changes settings on a cluster node
type NodeOption func(*ClusterNode)

// NodeOptionPort allows you to set the ssh port
func NodeOptionPort(port int) NodeOption {
	return func(node *ClusterNode) {
		node.Port = port
	}
}

func NodeOptionAuthMethod(auth ssh.AuthMethod) NodeOption {
	return func(node *ClusterNode) {
		node.Auth = append(node.Auth, auth)
	}
}

func NodeOptionKnownHostFile(file string) NodeOption {
	// if file == ""  return err if no file
	return func(node *ClusterNode) {
		node.Files.KnownHostsFile = file
	}
}

func NodeOptionPubKeyFile(file string) NodeOption {
	// if file == "" return err if no file
	return func(node *ClusterNode) {
		node.Files.PublicKeyFile = file
	}
}
