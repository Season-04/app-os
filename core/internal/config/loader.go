package config

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/staugaard/app-os/core/manifest"
)

func Load(directory string) (*Config, error) {
	files, err := filepath.Glob(filepath.Join(directory, "apps.enabled/*.json"))
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Manifests: make([]manifest.Manifest, 0, len(files)),
		Errors:    make([]ManifestError, 0),
	}

	for _, file := range files {
		manifest := manifest.Manifest{}
		bytes, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(bytes, &manifest)
		if err != nil {
			cfg.Errors = append(cfg.Errors, ManifestError{Path: file, Error: err})
		} else {
			cfg.Manifests = append(cfg.Manifests, manifest)
		}
	}

	return cfg, nil
}
