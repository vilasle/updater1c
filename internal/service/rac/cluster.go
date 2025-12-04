package rac

import (
	"bufio"
	"bytes"
	"strings"
)

const (
	clusterMainCmd = "cluster"
	clusterListCmd = "list"
)

type ClusterCommand struct {
	command string
	args    []string
}

func NewClusterCommand(command string, socket string) *ClusterCommand {
	return &ClusterCommand{
		command: command,
		args:    []string{socket},
	}
}

func (cmd *ClusterCommand) Command() []string {
	args := make([]string, 0)

	args = append(args, clusterMainCmd)
	args = append(args, clusterListCmd)
	args = append(args, cmd.command)
	args = append(args, cmd.args...)

	return args
}

func parseClusterInfoResponse(response []byte) (map[string]string, error) {
	rd := bytes.NewReader(response)
	sc := bufio.NewScanner(rd)

	clusterInfo := make(map[string]string)
	for sc.Scan() {
		kv := strings.Split(sc.Text(), ":")
		if len(kv) == 0 {
			continue
		}

		key, value := "", ""
		key = strings.TrimSpace(kv[0])

		if len(kv) > 1 {
			value = strings.TrimSpace(kv[1])
		}
		clusterInfo[key] = value
	}

	return clusterInfo, sc.Err()
}
