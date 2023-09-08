package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/Season-04/appos/core/manifest"
)

func (cfg *Config) requiredManifests() ([]manifest.Manifest, error) {
	coreManifest := manifest.Manifest{}
	bytes, err := os.ReadFile("appOS.manifest.json")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, &coreManifest)
	if err != nil {
		return nil, err
	}

	return []manifest.Manifest{
		coreManifest,
		{
			ID:    "appos.core-frontend",
			Name:  "AppOS Core Frontend",
			Image: "ghcr.io/season-04/appos.core-frontend",
			Routes: map[string]manifest.Route{
				"/": {
					Port: 80,
				},
			},
		},
		{
			ID:    "appos.ui",
			Name:  "AppOS UI",
			Image: "ghcr.io/season-04/appos.ui",
			Routes: map[string]manifest.Route{
				"/appos-ui": {
					Port: 3000,
				},
			},
		},
	}, nil
}

func (cfg *Config) ensureRequestAppsAreEnabled() error {
	required, err := cfg.requiredManifests()
	if err != nil {
		return err
	}

	for _, r := range required {
		installed := false
		for _, m := range cfg.manifests {
			if m.ID == r.ID {
				installed = true
				break
			}
		}
		if installed {
			continue
		}

		log.Println("Installing required application", r.ID)

		bytes, err := json.MarshalIndent(r, "", "  ")
		if err != nil {
			return err
		}
		fileName := filepath.Join(cfg.Directory, "apps.enabled/"+r.ID+".json")
		err = os.WriteFile(fileName, bytes, 0777)
		if err != nil {
			return err
		}

		cfg.manifests = append(cfg.manifests, r)
	}

	return nil
}
