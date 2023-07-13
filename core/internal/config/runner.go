package config

import (
	"context"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

func (cfg *Config) Run(ctx context.Context, dockerClient client.ContainerAPIClient) error {
	for _, manifest := range cfg.Manifests() {
		containers, err := dockerClient.ContainerList(ctx, types.ContainerListOptions{
			All:     true,
			Filters: filters.NewArgs(filters.Arg("name", manifest.ID)),
		})
		if err != nil {
			return err
		}

		if len(containers) == 0 {
			log.Println("Creating container for", manifest.ID)
			_, err = dockerClient.ContainerCreate(
				ctx,
				&container.Config{
					Image: manifest.Image,
				},
				&container.HostConfig{
					NetworkMode: "app-os",
				},
				&network.NetworkingConfig{},
				nil,
				manifest.ID,
			)
			if err != nil {
				return err
			}
		}

		stats, err := dockerClient.ContainerInspect(ctx, manifest.ID)
		if err != nil {
			return err
		}

		if !stats.State.Running || manifest.ID != "appos.core" {
			log.Println("(Re)starting container for", manifest.ID)
			err = dockerClient.ContainerRestart(ctx, manifest.ID, container.StopOptions{})
			if err != nil {
				return err
			}
		}
	}

	return nil

}
