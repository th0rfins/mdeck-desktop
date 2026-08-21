package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type WindowState struct {
	X           int  `json:"x"`
	Y           int  `json:"y"`
	Width       int  `json:"width"`
	Height      int  `json:"height"`
	IsMaximized bool `json:"is_maximized"`
}

func DefaultWindowState() WindowState {
	return WindowState{X: -1, Y: -1, Width: 1280, Height: 800}
}

func getConfigPath(name string) (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	base := filepath.Join(dir, "MDeck", "config")
	if err := os.MkdirAll(base, 0755); err != nil {
		return "", fmt.Errorf("mkdir %s: %w", base, err)
	}
	return filepath.Join(base, name), nil
}

func LoadWindowState() WindowState {
	def := DefaultWindowState()
	p, err := getConfigPath("window.json")
	if err != nil {
		return def
	}
	b, err := os.ReadFile(p)
	if err != nil {
		return def
	}
	var s WindowState
	if err := json.Unmarshal(b, &s); err != nil {
		return def
	}
	if s.Width < 600 {
		s.Width = def.Width
	}
	if s.Height < 400 {
		s.Height = def.Height
	}
	return s
}

func SaveWindowState(s WindowState) error {
	p, err := getConfigPath("window.json")
	if err != nil {
		return err
	}
	b, _ := json.MarshalIndent(s, "", "  ")
	return os.WriteFile(p, b, 0644)
}
