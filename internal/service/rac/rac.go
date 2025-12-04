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

func (rac *RemoteAdministrativeClient) LockInfobase(ctx context.Context, infobaseID string, secret string) error {
	return nil
}

func (rac *RemoteAdministrativeClient) UnlockInfobase(ctx context.Context, infobaseID string) error {
	return nil
}

func (rac *RemoteAdministrativeClient) TerminateAllSessions(ctx context.Context, infobaseID string) error {
	return nil
}

func (rac *RemoteAdministrativeClient) TerminateSession(ctx context.Context, sessionID string) error {
	return nil
}

func (rac *RemoteAdministrativeClient) execRACCommand(ctx context.Context, args ...string) ([]byte, error) {
	return nil, nil
}
