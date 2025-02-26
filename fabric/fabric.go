package fabric

import (
	"fmt"
	"net/http"
)

var (
	downloadURLFormat       = "https://meta.fabricmc.net/v2/versions/loader/%s/%s/%s/server/jar"
	defaultLoaderVersion    = "0.16.10"
	defaultInstallerVersion = "1.0.1"
)

type Config struct {
	Version          string
	LoaderVersion    string
	InstallerVersion string
}

func New(version string) *Config {
	return &Config{
		Version:          version,
		LoaderVersion:    defaultLoaderVersion,
		InstallerVersion: defaultInstallerVersion,
	}
}

func (c *Config) Mirror() (string, error) {
	url := fmt.Sprintf(downloadURLFormat, c.Version, c.LoaderVersion, c.InstallerVersion)

	resp, err := http.Head(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 399 {
		return "", fmt.Errorf("invalid version")
	}

	return url, nil
}
