package purpur

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ciathefed/jarchive/internal/utils"
)

var baseURL = "https://api.purpurmc.org/v2/purpur"

type Config struct {
	Version string
}

func New(version string) *Config {
	return &Config{
		Version: version,
	}
}

func (c *Config) Mirror() (string, error) {
	latestVersion, err := getLatestBuild(c.Version)
	if err != nil {
		return "", err
	}

	url, err := utils.URLJoin(
		baseURL,
		c.Version,
		latestVersion,
		"download",
	)
	if err != nil {
		return "", err
	}

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

func getLatestBuild(version string) (string, error) {
	url, err := utils.URLJoin(baseURL, version)
	if err != nil {
		return "", err
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 399 {
		return "", fmt.Errorf("invalid version")
	}

	var data struct {
		Builds struct {
			Latest string   `json:"string"`
			All    []string `json:"all"`
		} `json:"builds"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	if data.Builds.Latest == "" {
		return data.Builds.All[len(data.Builds.All)-1], nil
	}

	return data.Builds.Latest, nil
}
