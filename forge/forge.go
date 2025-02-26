package forge

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	promotionsSlimURL = "https://files.minecraftforge.net/net/minecraftforge/forge/promotions_slim.json"
	baseURL           = "https://maven.minecraftforge.net/net/minecraftforge/forge"
)

type Config struct {
	Version      string // Minecraft version
	ForgeVersion string // Forge version (optional)
}

func New(version string) *Config {
	return &Config{
		Version: version,
	}
}

// Mirror fetches the download URL for the Forge installer.
func (c *Config) Mirror() (string, error) {
	// If no Forge version is specified, fetch the latest one
	if c.ForgeVersion == "" {
		latestForgeVersion, err := getLatestForgeVersion(c.Version)
		if err != nil {
			return "", fmt.Errorf("failed to get latest Forge version: %w", err)
		}
		c.ForgeVersion = latestForgeVersion
	}

	// Construct the Maven URL for the Forge installer
	mavenURL := fmt.Sprintf(
		"%s/%s-%s/forge-%s-%s-installer.jar",
		baseURL,
		c.Version,
		c.ForgeVersion,
		c.Version,
		c.ForgeVersion,
	)

	// Verify the URL by making a HEAD request
	resp, err := http.Head(mavenURL)
	if err != nil {
		return "", fmt.Errorf("failed to verify URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 399 {
		return "", fmt.Errorf("invalid URL: status code %d", resp.StatusCode)
	}

	return mavenURL, nil
}

// getLatestForgeVersion fetches the latest Forge version for a specific Minecraft version.
func getLatestForgeVersion(mcVersion string) (string, error) {
	// Fetch the list of Forge versions for the specified Minecraft version
	resp, err := http.Get(promotionsSlimURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 399 {
		return "", fmt.Errorf("invalid response: status code %d", resp.StatusCode)
	}

	var promotions struct {
		Promos map[string]string `json:"promos"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&promotions); err != nil {
		return "", err
	}

	// Find the latest Forge version for the specified Minecraft version
	key := fmt.Sprintf("%s-latest", mcVersion)
	forgeVersion, ok := promotions.Promos[key]
	if !ok {
		return "", fmt.Errorf("no Forge version found for Minecraft version %s", mcVersion)
	}

	return forgeVersion, nil
}
