package fabric

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	config := New("1.18.2")
	assert.Equal(t, "1.18.2", config.Version)
	assert.Equal(t, "0.16.10", config.LoaderVersion)
	assert.Equal(t, "1.0.1", config.InstallerVersion)
}

func TestMirror_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	downloadURLFormat = server.URL + "/v2/versions/loader/%s/%s/%s/server/jar"

	config := New("1.18.2")
	mirrorURL, err := config.Mirror()

	assert.NoError(t, err)
	expectedURL := server.URL + "/v2/versions/loader/1.18.2/0.16.10/1.0.1/server/jar"
	assert.Equal(t, expectedURL, mirrorURL)
}

func TestMirror_CustomLoaderAndInstallerVersions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	downloadURLFormat = server.URL + "/v2/versions/loader/%s/%s/%s/server/jar"

	config := New("1.18.2")
	config.LoaderVersion = "0.15.0"
	config.InstallerVersion = "0.9.0"
	mirrorURL, err := config.Mirror()

	assert.NoError(t, err)
	expectedURL := server.URL + "/v2/versions/loader/1.18.2/0.15.0/0.9.0/server/jar"
	assert.Equal(t, expectedURL, mirrorURL)
}

func TestMirror_InvalidVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	downloadURLFormat = server.URL + "/v2/versions/loader/%s/%s/%s/server/jar"

	config := New("invalid-version")
	_, err := config.Mirror()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid version")
}

func TestMirror_NetworkError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	server.Close()

	downloadURLFormat = server.URL + "/v2/versions/loader/%s/%s/%s/server/jar"

	config := New("1.18.2")
	_, err := config.Mirror()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "connection refused")
}
