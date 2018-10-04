// misc.go contains various useful formatting/other functions
// Released under MIT liscence, copyright 2018 Tyler Ramer

package clusterExec

import (
	"strings"
)

// SSH commands need to be a single string, rather than command + []args
func composeCmd(cmd string, args []string) string {
	command := cmd + " " + strings.Join(args, " ")
	return command
}
