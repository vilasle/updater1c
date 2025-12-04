package rac

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

const (
	sessionMainCommand      = "session"
	sessionListCommand      = "list"
	sessionTerminateCommand = "terminate"
)

type SessionCommand struct {
	clusterId  string
	infobaseId string
	sessionId  string
	command    string
	args       []string
}

func NewSessionCommand(command string, socket string) *SessionCommand {
	return &SessionCommand{
		command: command,
		args:    []string{sessionListCommand, socket},
	}
}

func (cmd *SessionCommand) SetInfobaseId(id string) {
	cmd.infobaseId = id
}

func (cmd *SessionCommand) SetSessionId(id string) {
	cmd.sessionId = id
}

func (cmd *SessionCommand) SetClusterId(id string) {
	cmd.clusterId = id
}

func (cmd *SessionCommand) Command() []string {
	args := make([]string, 0)

	args = append(args, sessionMainCommand)
	args = append(args, cmd.command)
	args = append(args, fmt.Sprintf("--cluster=%s", cmd.clusterId))

	if cmd.command == sessionListCommand {
		args = append(args, fmt.Sprintf("--infobase=%s", cmd.infobaseId))
	}

	if cmd.command == sessionTerminateCommand {
		args = append(args, fmt.Sprintf("--session=%s", cmd.sessionId))
	}

	args = append(args, cmd.args...)

	return args
}

func parseSessoinListResponse(response []byte) ([]string, error) {
	rd := bytes.NewReader(response)
	sc := bufio.NewScanner(rd)

	sessionList := make([]string, 0)

	kSession := "session"
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

		if key == kSession && len(value) > 0 {
			sessionList = append(sessionList, value)
		}
	}
	return sessionList, sc.Err()
}
