package config

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/traefik/traefik/v2/pkg/config/dynamic"
	"gopkg.in/yaml.v3"
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
			err = cfg.pullImage(ctx, dockerClient, manifest.Image)
			if err != nil {
				return err
			}

			log.Println("Creating container for", manifest.ID)
			_, err = dockerClient.ContainerCreate(
				ctx,
				&container.Config{
					Image:    manifest.Image,
					Hostname: manifest.ID,
				},
				&container.HostConfig{
					NetworkMode: "appos",
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

		stats, err := dockerClient.ContainerInspect(ctx, manifest.ID)
		if err != nil {
			return err
		}

		// if manifest.ID != "appos.core" {
		if !stats.State.Running {
			log.Println("(Re)starting container for", manifest.ID)
			err = dockerClient.ContainerRestart(ctx, manifest.ID, container.StopOptions{})
			if err != nil {
				return err
			}
		}
	}

	return cfg.updateTraefik(ctx, dockerClient)
}

func (cfg *Config) pullImage(ctx context.Context, dockerClient DockerClient, image string) error {
	images, err := dockerClient.ImageList(ctx, types.ImageListOptions{
		Filters: filters.NewArgs(filters.Arg("reference", image)),
	})
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

func (cfg *Config) updateTraefik(ctx context.Context, dockerClient DockerClient) error {
	traefikConfig := &dynamic.Configuration{
		HTTP: &dynamic.HTTPConfiguration{
			Services: make(map[string]*dynamic.Service),
			Routers:  make(map[string]*dynamic.Router),
			Middlewares: map[string]*dynamic.Middleware{
				"appos-auth": {
					ForwardAuth: &dynamic.ForwardAuth{
						Address:             "http://appos.core:3000/auth/check",
						AuthResponseHeaders: []string{"X-AppOS-User"},
					},
				},
			},
		},
	}

	for _, manifest := range cfg.Manifests() {
		for path, route := range manifest.Routes {
			serviceName := fmt.Sprintf("%s-%v", manifest.ID, route.Port)
			traefikConfig.HTTP.Services[serviceName] = &dynamic.Service{
				LoadBalancer: &dynamic.ServersLoadBalancer{
					Servers: []dynamic.Server{
						{
							URL: fmt.Sprintf("http://%s:%v", manifest.ID, route.Port),
						},
					},
				},
			}

			traefikConfig.HTTP.Routers[path] = &dynamic.Router{
				EntryPoints: []string{"web"},
				Rule:        fmt.Sprintf("PathPrefix(`%s`)", path),
				Service:     serviceName,
				Middlewares: []string{"appos-auth"},
			}
		}
	}

	traefikConfigBytes, err := yaml.Marshal(traefikConfig)
	if err != nil {
		return err
	}
	traefikConfigFilename := filepath.Join(cfg.Directory, "traefik.yaml")
	err = os.WriteFile(traefikConfigFilename, traefikConfigBytes, 0777)
	if err != nil {
		return err
	}

	containers, err := dockerClient.ContainerList(ctx, types.ContainerListOptions{
		All:     true,
		Filters: filters.NewArgs(filters.Arg("name", "appos.traefik")),
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
					"--providers.file=true",
					"--providers.file.filename=/appos-config/traefik.yaml",
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
				NetworkMode: "appos",
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
						Source: cfg.hostConfigDir(),
						Target: "/appos-config",
					},
				},
			},
			&network.NetworkingConfig{},
			nil,
			"appos.traefik",
		)
		if err != nil {
			return err
		}
	}

	return dockerClient.ContainerRestart(ctx, "appos.traefik", container.StopOptions{})
}
