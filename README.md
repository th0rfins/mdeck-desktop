# MDeck Desktop

Wrapper Go + Wails v2 untuk MDeck server (multi-machine hub). Mirip Gotion: frameless traffic-light titlebar + persistent WebView profile + bookmark hub.

## Fitur
- **Frameless** `--wails-draggable` + traffic lights `#ff5f56/#ffbd2e/#27c93f` + 8 resize handles (injected via `internal/script`)
- **Hub sebelum startup**: grid bookmark `URL+label+warna` (persist `~/.config/MDeck/config/bookmarks.json`), default `hk1/hk2.projectpop.xyz`
- **Connect** = `window.location.href = bookmark.url` (load MDeck web di WebView yang sama, login pakai halaman MDeck sendiri)
- **Multi-window** seperti terminal: `New Window` → spawn proses baru `mdeck-desktop --url https://hk2...` (SingleInstanceLock di-disable per-window)
- **Tab browser-style di bawah titlebar**: polish `TabBar.tsx` jadi compact browser-tab (bukan di dalam titlebar)
- Platform: **Linux (WebKitGTK) + Windows (WebView2)**

## Struktur
```
mdeck-desktop/
  main.go, app.go, wails.json, go.mod
  internal/config, profile, navigation, bookmarks, script
  frontend/src (hub launcher)
  build/windows, build/linux
```

## Dev
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
wails dev
wails build           # native
wails build -nsis     # windows installer (needs NSIS)
```

## Release (GitHub Actions)
- `release.yml` build untuk linux/amd64 + windows/amd64 + checksums

## Server vs Client
- **Server**: `mdeck` (Node, `mdeck-*.zip`, systemd `mdeck.service`, `hk1.projectpop.xyz` via cloudflared)
- **Client**: `mdeck-desktop` (Go binary, `mdeck-desktop-windows-amd64.exe` / `mdeck-desktop-linux-amd64`)
