package config

import "github.com/staugaard/app-os/core/manifest"

type Config struct {
	Manifests []manifest.Manifest
	Errors    []ManifestError
}

type ManifestError struct {
	Path  string
	Error error
}
