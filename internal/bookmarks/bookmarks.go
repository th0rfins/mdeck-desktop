package bookmarks

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type Bookmark struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	URL   string `json:"url"`
	Color string `json:"color"` // mantine color: violet, teal, blue, etc
}

func path() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	base := filepath.Join(dir, "MDeck", "config")
	_ = os.MkdirAll(base, 0755)
	return filepath.Join(base, "bookmarks.json"), nil
}

func Load() []Bookmark {
	p, err := path()
	if err != nil {
		return defaultBookmarks()
	}
	b, err := os.ReadFile(p)
	if err != nil {
		return defaultBookmarks()
	}
	var list []Bookmark
	if err := json.Unmarshal(b, &list); err != nil || len(list) == 0 {
		return defaultBookmarks()
	}
	return list
}

func Save(list []Bookmark) error {
	p, err := path()
	if err != nil {
		return err
	}
	b, _ := json.MarshalIndent(list, "", "  ")
	return os.WriteFile(p, b, 0644)
}

func Upsert(b Bookmark) ([]Bookmark, error) {
	list := Load()
	b.URL = strings.TrimSpace(b.URL)
	b.Label = strings.TrimSpace(b.Label)
	if b.URL == "" || b.Label == "" {
		return list, nil
	}
	if b.ID == "" {
		b.ID = newID()
	}
	if b.Color == "" {
		b.Color = "violet"
	}
	found := false
	for i, it := range list {
		if it.ID == b.ID {
			list[i] = b
			found = true
			break
		}
	}
	if !found {
		list = append(list, b)
	}
	return list, Save(list)
}

func Delete(id string) ([]Bookmark, error) {
	list := Load()
	out := []Bookmark{}
	for _, it := range list {
		if it.ID != id {
			out = append(out, it)
		}
	}
	return out, Save(out)
}

func newID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func defaultBookmarks() []Bookmark {
	return []Bookmark{
		{ID: "hk1", Label: "HK1 Production", URL: "https://hk1.projectpop.xyz", Color: "violet"},
		{ID: "hk2", Label: "HK2 Staging", URL: "https://hk2.projectpop.xyz", Color: "teal"},
	}
}
