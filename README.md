# clusterExec [![Build Status](https://travis-ci.org/tylarb/clusterExec.svg?branch=master)](https://travis-ci.org/tylarb/clusterExec)  [![Go Report Card](https://goreportcard.com/badge/github.com/tylarb/clusterexec)](https://goreportcard.com/report/github.com/tylarb/clusterexec)
Execute bash across a cluster using Golang's ssh package


## Concept

clusterExec is an attempt to abstract command execution accross a cluster quickly and easily. It is intended to be a flexible API which can easily be adapted for cluster-centric utilities. 

## Usage

clusterExec is built around the concept of _clusterNodes_ and _clusterCommands_. ClusterCommands can be created and executed across nodes. 

Each node is one of a localhost or remote host, and appropriate ssh connection files must be provided to a node. This is made simple by included code, including a flexible API. Take some examples:

~~~
>> A node connecting on standard port 22 with private key auth, and no strict host checking:

AuthKey, err := clusterExec.GetPrivateAuthKey(<private_key_file>)
if err != nil {
   log.Fatal(err)
}
node, err := clusterExec.CreateNode(<user>, <hostname>, clusterExec.NodeOptionAuthMethod(AuthKey), clusterExec.NodeOptionHostKeyCheck(false))


>> A node connecting on port 25 with password auth:

node, err := clusterExec.CreateNode(<user>, <hostname>, clusterExec.NodeOptionPort(<port>), clusterExec.NodeOptionAuthMethod(ssh.Password(<password>)), clusterExec.NodeOptionKnownHostsFile(<known_hosts_file>))

>> Localhost node 

node, err := clusterExec.CreateNode(<user>, "", clusterExec.NodeOptionIsLocalhost())

~~~

After a node is created, connect to it - this allows ssh connections to be shared across commands:

~~~
if err := node.Dail(); err != nil {
  log.Error(err)
}
defer node.Close()
~~~

Once you have created a node, you can create a cluster command, with the option of a timeout, to execute accross the cluster:

~~~
clusterCmd := clusterExec.CreateClusterCommand("ls", []string{"-l"}, clusterExec.ClusterCmdOptionTimeout(5 * time.Second))

// On a single node:
stdOut, stdErr, err := node.Run(clusterCmd)
if err != nil {
   log.Error(err)
}

fmt.Printf("Output of command was %s\n", stdOut)

~~~

## Contributing

Feel free to fork and submit pull requests, or submit issues and feature requests.

