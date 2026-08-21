package navigation

import (
	"net/url"
	"strings"
)

// IsInternalURL: allow *.projectpop.xyz, localhost, and MDeck subdomains.
// External links -> open in system browser.
func IsInternalURL(raw string) bool {
	if raw == "" {
		return false
	}
	u, err := url.Parse(raw)
	if err != nil {
		return false
	}
	h := strings.ToLower(u.Hostname())
	if h == "localhost" || h == "127.0.0.1" {
		return true
	}
	// your cloudflared base domain + any mdeck hosts
	if h == "projectpop.xyz" || strings.HasSuffix(h, ".projectpop.xyz") {
		return true
	}
	// allow file:// and wails:// for local launcher
	if u.Scheme == "file" || u.Scheme == "wails" || u.Scheme == "http" && strings.Contains(h, "localhost") {
		return true
	}
	return false
}
