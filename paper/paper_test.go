package paper

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
}

func TestMirror_Success(t *testing.T) {
	buildsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]any{
			"builds": []int{100, 101, 102},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer buildsServer.Close()

	downloadServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer downloadServer.Close()

	baseURL = buildsServer.URL + "/v2/projects/paper"

	config := New("1.18.2")
	mirrorURL, err := config.Mirror()

	assert.NoError(t, err)
	expectedURL := buildsServer.URL + "/v2/projects/paper/versions/1.18.2/builds/102/downloads/paper-1.18.2-102.jar"
	assert.Equal(t, expectedURL, mirrorURL)
}

func TestMirror_InvalidVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	baseURL = server.URL + "/v2/projects/paper"

	config := New("invalid-version")
	_, err := config.Mirror()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid version")
}

func TestGetLatestBuild_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]any{
			"builds": []int{100, 101, 102},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	baseURL = server.URL + "/v2/projects/paper"

	latestBuild, err := getLatestBuild("1.18.2")

	assert.NoError(t, err)
	assert.Equal(t, 102, latestBuild)
}

func TestGetLatestBuild_NoBuilds(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]any{
			"builds": []int{},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	baseURL = server.URL + "/v2/projects/paper"

	_, err := getLatestBuild("1.18.2")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no builds found for version")
}

func TestGetLatestBuild_InvalidVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	baseURL = server.URL + "/v2/projects/paper"

	_, err := getLatestBuild("invalid-version")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid version")
}
