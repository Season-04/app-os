package config

import (
	"sync"

	"github.com/Season-04/app-os/core/manifest"
)

type Config struct {
	Directory string
	manifests []manifest.Manifest
	errors    []ManifestError
	mtx       *sync.RWMutex
}

type ManifestError struct {
	Path  string
	Error error
}

func NewConfig(directory string) *Config {
	return &Config{
		Directory: directory,
		manifests: []manifest.Manifest{},
		errors:    []ManifestError{},
		mtx:       &sync.RWMutex{},
	}
}

func (cfg *Config) Manifests() []manifest.Manifest {
	cfg.mtx.RLock()
	defer cfg.mtx.RUnlock()

	return cfg.manifests
}

func (cfg *Config) Errors() []ManifestError {
	cfg.mtx.RLock()
	defer cfg.mtx.RUnlock()

	return cfg.errors
}
