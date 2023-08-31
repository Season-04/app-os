package config

import (
	"context"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type DockerClient interface {
	client.ContainerAPIClient
	client.ImageAPIClient
}

func (cfg *Config) hostConfigDir() string {
	configRoot := os.Getenv("CONFIG_DIR")
	if configRoot == "" {
		panic("Missing CONFIG_DIR envoronment variable")
	}
	return configRoot
}

func (cfg *Config) Run(ctx context.Context, dockerClient DockerClient) error {
	for _, manifest := range cfg.Manifests() {
		relativeAppDataPath := filepath.Join("app-data", manifest.ID)
		err := os.MkdirAll(filepath.Join(cfg.Directory, relativeAppDataPath), 0777)
		if err != nil {
			return err
		}

		containers, err := dockerClient.ContainerList(ctx, types.ContainerListOptions{
			Filters: filters.NewArgs(filters.Arg("name", "^/"+manifest.ID+"$")),
		})
		if err != nil {
			return err
		}

		if len(containers) == 0 {
			// err = cfg.pullImage(ctx, dockerClient, manifest.ID)
			// if err != nil {
			// 	return err
			// }

			labels := map[string]string{}
			for path, route := range manifest.Routes {
				labels["traefik.enable"] = "true"
				routeName := fmt.Sprintf("route%v", crc32.ChecksumIEEE([]byte(path)))
				labels[fmt.Sprintf("traefik.http.services.%s.loadbalancer.server.port", routeName)] = strconv.FormatUint(uint64(route.Port), 10)
				labels[fmt.Sprintf("traefik.http.routers.%s.entrypoints", routeName)] = "web"
				labels[fmt.Sprintf("traefik.http.routers.%s.service", routeName)] = routeName
				labels[fmt.Sprintf("traefik.http.routers.%s.rule", routeName)] = fmt.Sprintf("PathPrefix(`%s`)", path)
				labels[fmt.Sprintf("traefik.http.routers.%s.middlewares", routeName)] = "appos-auth"
				labels["traefik.http.middlewares.appos-auth.forwardauth.address"] = "http://appos.core:3000/auth/check"
				labels["traefik.http.middlewares.appos-auth.forwardauth.authResponseHeaders"] = "X-AppOS-User"
			}

			log.Println("Creating container for", manifest.ID)
			_, err = dockerClient.ContainerCreate(
				ctx,
				&container.Config{
					Image:    manifest.Image,
					Labels:   labels,
					Hostname: manifest.ID,
				},
				&container.HostConfig{
					NetworkMode: "app-os",
					Mounts: []mount.Mount{
						{
							Type:   mount.TypeBind,
							Source: filepath.Join(cfg.hostConfigDir(), relativeAppDataPath),
							Target: "/data",
						},
					},
				},
				&network.NetworkingConfig{},
				nil,
				manifest.ID,
			)
			if err != nil {
				return err
			}
		}

		// stats, err := dockerClient.ContainerInspect(ctx, manifest.ID)
		// if err != nil {
		// 	return err
		// }

		// if manifest.ID != "appos.core" {
		// if !stats.State.Running || manifest.ID != "appos.core" {
		log.Println("(Re)starting container for", manifest.ID)
		err = dockerClient.ContainerRestart(ctx, manifest.ID, container.StopOptions{})
		if err != nil {
			return err
		}
		// }
	}

	containers, err := dockerClient.ContainerList(ctx, types.ContainerListOptions{
		All:     true,
		Filters: filters.NewArgs(filters.Arg("name", "appos-traefik")),
	})
	if err != nil {
		return err
	}

	if len(containers) == 0 {
		err = cfg.pullImage(ctx, dockerClient, "traefik:v2.10")
		if err != nil {
			return err
		}

		log.Println("Creating container for traefik")
		_, err = dockerClient.ContainerCreate(
			ctx,
			&container.Config{
				Image: "traefik:v2.10",
				Cmd: []string{
					"--api.insecure=true",
					"--providers.docker=true",
					"--providers.docker.exposedByDefault=false",
					"--entrypoints.web.address=:80",
					"--log.level=DEBUG",
				},
				ExposedPorts: nat.PortSet{
					"80/tcp":   struct{}{},
					"443/tcp":  struct{}{},
					"8080/tcp": struct{}{},
				},
			},
			&container.HostConfig{
				NetworkMode: "app-os",
				PortBindings: nat.PortMap{
					"80/tcp": []nat.PortBinding{
						{
							HostIP:   "0.0.0.0",
							HostPort: "80",
						},
					},
					"443/tcp": []nat.PortBinding{
						{
							HostIP:   "0.0.0.0",
							HostPort: "443",
						},
					},
					"8080/tcp": []nat.PortBinding{
						{
							HostIP:   "0.0.0.0",
							HostPort: "8080",
						},
					},
				},
				Mounts: []mount.Mount{
					{
						Type:   mount.TypeBind,
						Source: "/var/run/docker.sock",
						Target: "/var/run/docker.sock",
					},
				},
			},
			&network.NetworkingConfig{},
			nil,
			"appos-traefik",
		)
		if err != nil {
			return err
		}
	}

	err = dockerClient.ContainerRestart(ctx, "appos-traefik", container.StopOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (cfg *Config) pullImage(ctx context.Context, dockerClient DockerClient, image string) error {
	images, err := dockerClient.ImageList(ctx, types.ImageListOptions{
		Filters: filters.NewArgs(filters.Arg("reference", image)),
	})
	log.Println("images", images, len(images))
	if err != nil {
		return err
	}

	if len(images) > 0 {
		log.Println("got image")
		return nil
	}

	log.Println("Pulling image", image)
	r, err := dockerClient.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	_, err = io.ReadAll(r)
	return err
}
