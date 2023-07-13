package runner

import (
	"context"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/staugaard/app-os/core/internal/config"
)

func Run(ctx context.Context, cfg config.Config) error {
	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	for _, manifest := range cfg.Manifests {
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
