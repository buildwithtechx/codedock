package utils

import (
	"context"

	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"

	"codedock.run/codedock/internal/config"
)

func GetDataDir() string {
	return config.Get().Server.DataDir
}

func GetRuntimeNetwork() string {
	return config.Get().Docker.RuntimeNetwork
}

func EnsureCodedockNetwork(ctx context.Context, cli *client.Client) error {
	if cli == nil {
		return nil
	}

	netName := GetRuntimeNetwork()

	networks, err := cli.NetworkList(ctx, network.ListOptions{})
	if err != nil {
		return err
	}
	for _, net := range networks {
		if net.Name == netName {
			return nil
		}
	}
	_, err = cli.NetworkCreate(ctx, netName, network.CreateOptions{
		Driver: "bridge",
	})
	return err
}
