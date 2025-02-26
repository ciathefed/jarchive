package paper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ciathefed/jarchive/internal/utils"
)

var baseURL = "https://api.papermc.io/v2/projects/paper"

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
		"versions",
		c.Version,
		"builds",
		strconv.Itoa(latestVersion),
		"downloads",
		fmt.Sprintf("paper-%s-%d.jar", c.Version, latestVersion),
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

func getLatestBuild(version string) (int, error) {
	url, err := utils.URLJoin(baseURL, "versions", version)
	if err != nil {
		return 0, err
	}

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 399 {
		return 0, fmt.Errorf("invalid version")
	}

	var data struct {
		Builds []int `json:"builds"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	if len(data.Builds) == 0 {
		return 0, fmt.Errorf("no builds found for version %s", version)
	}

	return data.Builds[len(data.Builds)-1], nil
}
