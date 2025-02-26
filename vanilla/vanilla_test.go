package vanilla

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
	assert.Nil(t, config.versionManifest)
}

func TestLoadVersionManifest_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := versionManifest{
			Versions: []struct {
				ID  string `json:"id"`
				URL string `json:"url"`
			}{
				{ID: "1.18.2", URL: "https://example.com/1.18.2.json"},
				{ID: "1.17.1", URL: "https://example.com/1.17.1.json"},
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	versionManifestURL = server.URL

	config := New("1.18.2")
	err := config.loadVersionManifest()

	assert.NoError(t, err)
	assert.NotNil(t, config.versionManifest)
	assert.Equal(t, 2, len(config.versionManifest.Versions))
}

func TestLoadVersionManifest_Failure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	versionManifestURL = server.URL

	config := New("1.18.2")
	err := config.loadVersionManifest()

	assert.Error(t, err)
	assert.Nil(t, config.versionManifest)
}

func TestMirror_Success(t *testing.T) {
	manifestServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := versionManifest{
			Versions: []struct {
				ID  string `json:"id"`
				URL string `json:"url"`
			}{
				{ID: "1.18.2", URL: "https://example.com/1.18.2.json"},
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer manifestServer.Close()

	detailsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]any{
			"downloads": map[string]any{
				"server": map[string]any{
					"url": "https://example.com/server.jar",
				},
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer detailsServer.Close()

	versionManifestURL = manifestServer.URL
	config := New("1.18.2")
	config.versionManifest = &versionManifest{
		Versions: []struct {
			ID  string `json:"id"`
			URL string `json:"url"`
		}{
			{ID: "1.18.2", URL: detailsServer.URL},
		},
	}

	mirrorURL, err := config.Mirror()

	assert.NoError(t, err)
	assert.Equal(t, "https://example.com/server.jar", mirrorURL)
}

func TestMirror_InvalidVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := versionManifest{
			Versions: []struct {
				ID  string `json:"id"`
				URL string `json:"url"`
			}{
				{ID: "1.18.2", URL: "https://example.com/1.18.2.json"},
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	versionManifestURL = server.URL

	config := New("invalid-version")
	_, err := config.Mirror()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid version")
}

func TestMirror_VersionDetailsFailure(t *testing.T) {
	manifestServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := versionManifest{
			Versions: []struct {
				ID  string `json:"id"`
				URL string `json:"url"`
			}{
				{ID: "1.18.2", URL: "https://example.com/1.18.2.json"},
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer manifestServer.Close()

	detailsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer detailsServer.Close()

	versionManifestURL = manifestServer.URL
	config := New("1.18.2")
	config.versionManifest = &versionManifest{
		Versions: []struct {
			ID  string `json:"id"`
			URL string `json:"url"`
		}{
			{ID: "1.18.2", URL: detailsServer.URL},
		},
	}

	_, err := config.Mirror()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to fetch version details: status 500")
}
