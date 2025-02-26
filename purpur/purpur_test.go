package purpur

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
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/purpur/1.18.2":
			response := map[string]any{
				"builds": map[string]any{
					"latest": "123",
					"all":    []string{"120", "121", "122", "123"},
				},
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		case "/v2/purpur/1.18.2/123/download":
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	baseURL = server.URL + "/v2/purpur"

	config := New("1.18.2")
	mirrorURL, err := config.Mirror()

	assert.NoError(t, err)
	assert.Equal(t, server.URL+"/v2/purpur/1.18.2/123/download", mirrorURL)
}

func TestMirror_InvalidVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	baseURL = server.URL + "/v2/purpur"

	config := New("invalid-version")
	_, err := config.Mirror()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid version")
}

func TestGetLatestBuild_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]any{
			"builds": map[string]any{
				"latest": "123",
				"all":    []string{"120", "121", "122", "123"},
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	baseURL = server.URL + "/v2/purpur"

	latestBuild, err := getLatestBuild("1.18.2")

	assert.NoError(t, err)
	assert.Equal(t, "123", latestBuild)
}

func TestGetLatestBuild_NoLatestBuild(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]any{
			"builds": map[string]any{
				"latest": "",
				"all":    []string{"120", "121", "122", "123"},
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	baseURL = server.URL + "/v2/purpur"

	latestBuild, err := getLatestBuild("1.18.2")

	assert.NoError(t, err)
	assert.Equal(t, "123", latestBuild)
}

func TestGetLatestBuild_InvalidVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	baseURL = server.URL + "/v2/purpur"

	_, err := getLatestBuild("invalid-version")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid version")
}
