package rac

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

const (
	infobaseMainCommand    = "infobase"
	infobaseListAllCommand = "summary list"
)

type Infobase struct {
	id   string
	name string
	desc string
}

type InfobaseCommand struct {
	clusterId  string
	infobaseId string
	command    string
	auth       [2]string
	args       []string
}

func NewInfobaseCommand(command string, socket string) *InfobaseCommand {
	return &InfobaseCommand{
		command: command,
		args:    []string{socket},
		auth:    [2]string{},
	}
}

func (cmd *InfobaseCommand) SetClusterId(id string) {
	cmd.clusterId = id
}

func (cmd *InfobaseCommand) SetInfobaseId(id string) {
	cmd.infobaseId = id
}

func (cmd *InfobaseCommand) SetAuth(user, password string) {
	cmd.auth[0] = user
	cmd.auth[1] = password
}

func (cmd *InfobaseCommand) Command() (string, []string) {
	args := make([]string, 0, 4+len(cmd.args)+len(cmd.auth))

	args = append(args, infobaseMainCommand)
	args = append(args, cmd.command)
	args = append(args, fmt.Sprintf("--cluster=%s", cmd.clusterId))
	args = append(args, cmd.args...)

	return infobaseMainCommand, args
}

func parseInfobaseListResponse(response []byte) (map[string]Infobase, error) {
	rd := bytes.NewReader(response)
	sc := bufio.NewScanner(rd)

	infobaseList := make(map[string]Infobase)

	kInfobase, kName, kDesc := "infobase", "name", "desc"
	infobase := Infobase{}
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

		switch key {
		case kInfobase:
			if len(infobase.id) > 0 {
				infobaseList[infobase.name] = infobase
			}
			infobase = Infobase{
				id: value,
			}
		case kName:
			infobase.name = value
		case kDesc:
			infobase.desc = value
		}
	}

	infobaseList[infobase.name] = infobase

	return infobaseList, sc.Err()
}
