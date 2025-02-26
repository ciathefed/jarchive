package vanilla

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var versionManifestURL = "https://launchermeta.mojang.com/mc/game/version_manifest.json"

type versionManifest struct {
	Versions []struct {
		ID  string `json:"id"`
		URL string `json:"url"`
	} `json:"versions"`
}

type Config struct {
	Version         string
	versionManifest *versionManifest
}

func New(version string) *Config {
	return &Config{
		Version:         version,
		versionManifest: nil,
	}
}

func (s *Config) loadVersionManifest() error {
	if s.versionManifest != nil {
		return nil
	}

	manifest := new(versionManifest)

	resp, err := http.Get(versionManifestURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, manifest); err != nil {
		return err
	}

	s.versionManifest = manifest
	return nil
}

func (c *Config) Mirror() (string, error) {
	if err := c.loadVersionManifest(); err != nil {
		return "", fmt.Errorf("failed to get version manifest: %v", err)
	}

	for _, v := range c.versionManifest.Versions {
		if v.ID == c.Version {
			resp, err := http.Get(v.URL)
			if err != nil {
				return "", fmt.Errorf("failed to fetch version details: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode > 399 {
				return "", fmt.Errorf("failed to fetch version details: status %d", resp.StatusCode)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return "", fmt.Errorf("failed to read version details: %v", err)
			}

			var details struct {
				Downloads struct {
					Server struct {
						URL string `json:"url"`
					} `json:"server"`
				} `json:"downloads"`
			}

			if err := json.Unmarshal(body, &details); err != nil {
				return "", fmt.Errorf("failed to decode version details: %v", err)
			}

			return details.Downloads.Server.URL, nil
		}
	}

	return "", fmt.Errorf("invalid version")
}
