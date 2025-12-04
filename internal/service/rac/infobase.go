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
	infobaseUpdate         = "update"
)

type Infobase struct {
	id   string
	name string
	desc string
}

type InfobaseCommand struct {
	clusterId   string
	infobaseId  string
	command     string
	args        []string
	lockSession bool
	lockCode    string
	auth        [2]string
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

func (cmd *InfobaseCommand) SetLockSession(lockSession bool, lockCode string) {
	cmd.lockSession = true
	cmd.lockCode = lockCode
}

func (cmd *InfobaseCommand) Command() []string {
	args := make([]string, 0)

	args = append(args, infobaseMainCommand)
	args = append(args, cmd.command)
	args = append(args, fmt.Sprintf("--cluster=%s", cmd.clusterId))

	if cmd.command == infobaseUpdate {
		args = append(args, cmd.infobaseUpdateArgs()...)
	}

	args = append(args, cmd.args...)

	return args
}

func (cmd *InfobaseCommand) infobaseUpdateArgs() []string {
	args := make([]string, 0)
	//add infobase id
	args = append(args,
		fmt.Sprintf("--infobase=%s", cmd.infobaseId),
	)
	//add auth infobase if exists
	if len(cmd.auth[0]) > 0 {
		args = append(args,
			fmt.Sprintf("--infobase-user=%s", cmd.auth[0]),
			fmt.Sprintf("--infobase-pwd=%s", cmd.auth[1]),
		)
	}
	//add access code if exists
	if len(cmd.lockCode) > 0 {
		args = append(args,
			fmt.Sprintf("--permission-code=%s", cmd.lockCode),
		)
	}

	//set lock session mode
	denySession := "off"
	if cmd.lockSession {
		denySession = "on"
	}
	args = append(args, fmt.Sprintf("--sessions-deny=%s", denySession))
	return args
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
