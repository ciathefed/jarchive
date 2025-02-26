package forge

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	config := New("1.18.2")
	assert.Equal(t, "1.18.2", config.Version)
	assert.Equal(t, "", config.ForgeVersion)
}

func TestMirror_Success(t *testing.T) {
	promotionsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]any{
			"promos": map[string]string{
				"1.18.2-latest": "40.1.0",
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer promotionsServer.Close()

	mavenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer mavenServer.Close()

	promotionsSlimURL = promotionsServer.URL
	baseURL = mavenServer.URL + "/net/minecraftforge/forge"

	config := New("1.18.2")
	mirrorURL, err := config.Mirror()

	assert.NoError(t, err)
	expectedURL := mavenServer.URL + "/net/minecraftforge/forge/1.18.2-40.1.0/forge-1.18.2-40.1.0-installer.jar"
	assert.Equal(t, expectedURL, mirrorURL)
}

func TestMirror_WithForgeVersion(t *testing.T) {
	mavenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer mavenServer.Close()

	baseURL = mavenServer.URL + "/net/minecraftforge/forge"

	config := New("1.18.2")
	config.ForgeVersion = "40.1.0"
	mirrorURL, err := config.Mirror()

	assert.NoError(t, err)
	expectedURL := mavenServer.URL + "/net/minecraftforge/forge/1.18.2-40.1.0/forge-1.18.2-40.1.0-installer.jar"
	assert.Equal(t, expectedURL, mirrorURL)
}

func TestMirror_InvalidMinecraftVersion(t *testing.T) {
	promotionsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]any{
			"promos": map[string]string{
				"1.18.2-latest": "40.1.0",
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer promotionsServer.Close()

	promotionsSlimURL = promotionsServer.URL

	config := New("invalid-version")
	_, err := config.Mirror()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no Forge version found for Minecraft version")
}

func TestMirror_InvalidMavenURL(t *testing.T) {
	promotionsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]any{
			"promos": map[string]string{
				"1.18.2-latest": "40.1.0",
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer promotionsServer.Close()

	mavenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer mavenServer.Close()

	promotionsSlimURL = promotionsServer.URL
	baseURL = mavenServer.URL + "/net/minecraftforge/forge"

	config := New("1.18.2")
	_, err := config.Mirror()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid URL: status code 404")
}

func TestGetLatestForgeVersion_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]any{
			"promos": map[string]string{
				"1.18.2-latest": "40.1.0",
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	promotionsSlimURL = server.URL

	forgeVersion, err := getLatestForgeVersion("1.18.2")

	assert.NoError(t, err)
	assert.Equal(t, "40.1.0", forgeVersion)
}

func TestGetLatestForgeVersion_InvalidMinecraftVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]any{
			"promos": map[string]string{
				"1.18.2-latest": "40.1.0",
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	promotionsSlimURL = server.URL

	_, err := getLatestForgeVersion("invalid-version")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no Forge version found for Minecraft version")
}

func TestGetLatestForgeVersion_InvalidResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	promotionsSlimURL = server.URL

	_, err := getLatestForgeVersion("1.18.2")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid response: status code 500")
}
