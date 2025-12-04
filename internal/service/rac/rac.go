// Module implements work with 1C:Remote Administrative Client
package rac

import (
	"context"
	"fmt"

	"github.com/vilasle/updater1c/internal/util"
)

type RemoteAdministrativeClient struct {
	socket   string
	user     string
	password string
	executor util.Executor
}

func NewRemoteAdministrativeClient(host, port, password, user string, executor util.Executor) *RemoteAdministrativeClient {
	return &RemoteAdministrativeClient{
		socket:   fmt.Sprintf("%s:%s", host, port),
		user:     user,
		password: password,
		executor: executor,
	}
}

func (rac *RemoteAdministrativeClient) GetClusterId(ctx context.Context) (string, error) {
	select {
	default:
	case <-ctx.Done():
		return "", ctx.Err()
	}

	command := NewClusterCommand(clusterListCmd, rac.socket)

	result := rac.executor.Execute(command)

	if err := result.Error(); err != nil {
		return "", err
	}

	if result.Code() != 0 {
		return "", fmt.Errorf(
			"getting information about cluster failed. error info: %s; code: %d",
			string(result.Stderr()), result.Code())
	}

	clusterInfo, err := parseClusterInfoResponse(result.Stdout())

	return clusterInfo["cluster"], err
}

func (rac *RemoteAdministrativeClient) GetInfobaseId(ctx context.Context, infobaseName, clusterId string) (Infobase, error) {
	select {
	default:
	case <-ctx.Done():
		return Infobase{}, ctx.Err()
	}

	command := NewInfobaseCommand(infobaseListAllCommand, rac.socket)
	command.SetClusterId(clusterId)

	result := rac.executor.Execute(command)

	if err := result.Error(); err != nil {
		return Infobase{}, err
	}

	if result.Code() != 0 {
		return Infobase{}, fmt.Errorf(
			"getting list of information bases failed. error info: %s; code: %d",
			string(result.Stderr()), result.Code())
	}

	ibLs, err := parseInfobaseListResponse(result.Stdout())
	if err != nil {
		return Infobase{}, err
	}

	if ib, ok := ibLs[infobaseName]; ok {
		return ib, nil
	} else {
		err = fmt.Errorf("infobase %s not found in cluster %s", infobaseName, clusterId)
	}
	return Infobase{}, err
}

type LockInfobaseParams struct {
	InfobaseID       string
	AccessToken      string
	InfobaseUser     string
	InfobasePassword string
}

func (rac *RemoteAdministrativeClient) LockInfobase(ctx context.Context, clusterId string, lockParams LockInfobaseParams) error {
	select {
	default:
	case <-ctx.Done():
		return ctx.Err()
	}

	command := rac.lockingCommand(true, clusterId, lockParams)

	result := rac.executor.Execute(command)
	if err := result.Error(); err != nil {
		if err != nil {
			return err
		}
	}

	if result.Code() != 0 {
		return fmt.Errorf(
			"locking infobase failed. error info: %s; code: %d",
			string(result.Stderr()), result.Code())
	}
	return nil
}

func (rac *RemoteAdministrativeClient) UnlockInfobase(ctx context.Context, clusterId string, lockParams LockInfobaseParams) error {
	select {
	default:
	case <-ctx.Done():
		return ctx.Err()
	}
	command := rac.lockingCommand(false, clusterId, lockParams)

	result := rac.executor.Execute(command)
	if err := result.Error(); err != nil {
		if err != nil {
			return err
		}
	}

	if result.Code() != 0 {
		return fmt.Errorf(
			"unlocking infobase failed. error info: %s; code: %d",
			string(result.Stderr()), result.Code())
	}
	return nil
}

func (rac *RemoteAdministrativeClient) lockingCommand(lock bool, clusterId string, lockParams LockInfobaseParams) util.Commander {
	command := NewInfobaseCommand(infobaseUpdate, rac.socket)
	command.SetClusterId(clusterId)
	command.SetInfobaseId(lockParams.InfobaseID)
	command.SetLockSession(true, lockParams.AccessToken)
	command.SetAuth(lockParams.InfobaseUser, lockParams.InfobasePassword)

	return command
}

func (rac *RemoteAdministrativeClient) TerminateAllSessions(ctx context.Context, clusterID, infobaseID string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	allSessionCmd := NewSessionCommand(sessionListCommand, rac.socket)
	allSessionCmd.SetClusterId(clusterID)
	allSessionCmd.SetInfobaseId(infobaseID)

	result := rac.executor.Execute(allSessionCmd)

	if err := result.Error(); err != nil {
		return err
	}

	if result.Code() != 0 {
		return fmt.Errorf(
			"getting list of sessions failed. error info: %s; code: %d",
			string(result.Stderr()), result.Code())
	}

	ls, err := parseSessoinListResponse(result.Stdout())
	if err != nil {
		return err
	}

	for _, sessionID := range ls {
		if err != rac.TerminateSession(ctx, clusterID, sessionID) {
			if err != nil {
				//TODO: log
				continue
			}
		}
	}
	return nil
}

func (rac *RemoteAdministrativeClient) TerminateSession(ctx context.Context, clusterID, sessionID string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	terminateCmd := NewSessionCommand(sessionTerminateCommand, rac.socket)
	terminateCmd.SetClusterId(clusterID)
	terminateCmd.SetSessionId(sessionID)

	result := rac.executor.Execute(terminateCmd)

	if err := result.Error(); err != nil {
		return err
	}

	if result.Code() != 0 {
		return fmt.Errorf(
			"terminating session failed. error info: %s; code: %d",
			string(result.Stderr()), result.Code())
	}
	return nil
}
