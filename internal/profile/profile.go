package profile

import (
	"os"
	"path/filepath"
)

func GetUserDataDir() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	// like Gotion: persistent WebView profile
	p := filepath.Join(dir, "MDeck", "profile")
	if err := os.MkdirAll(p, 0755); err != nil {
		return "", err
	}
	return p, nil
}
