package utils

import (
	"fmt"
	"net/url"
	"path"
)

func URLJoin(u string, elem ...string) (string, error) {
	t, err := url.Parse(u)
	if err != nil || t.Scheme == "" || t.Host == "" {
		return "", fmt.Errorf("invalid URL: %s", u)
	}
	t.Path = path.Join(append([]string{t.Path}, elem...)...)
	return t.String(), nil
}
