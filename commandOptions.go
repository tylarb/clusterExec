package clusterexec

import "time"

// ClusterCmdOption is a functional option on a cluster command
type ClusterCmdOption func(*ClusterCmd)

// ClusterCmdOptionTimeout sets the timeout of a command
func ClusterCmdOptionTimeout(t time.Duration) ClusterCmdOption {
	return func(cmd *ClusterCmd) {
		cmd.Timeout = t
	}
}
