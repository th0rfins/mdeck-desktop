package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"mdeck-desktop/internal/bookmarks"
	"mdeck-desktop/internal/config"
	"mdeck-desktop/internal/navigation"
	"mdeck-desktop/internal/profile"
	"mdeck-desktop/internal/script"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx         context.Context
	userDataDir string
	windowState config.WindowState
}

func NewApp() *App {
	dir, err := profile.GetUserDataDir()
	if err != nil {
		log.Printf("[MDeck] userDataDir err: %v", err)
	}
	ws := config.LoadWindowState()
	return &App{userDataDir: dir, windowState: ws}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	log.Printf("[MDeck] startup dir=%s", a.userDataDir)
	if a.windowState.X >= 0 && a.windowState.Y >= 0 {
		runtime.WindowSetPosition(ctx, a.windowState.X, a.windowState.Y)
		runtime.WindowSetSize(ctx, a.windowState.Width, a.windowState.Height)
	} else {
		runtime.WindowCenter(ctx)
	}
	if a.windowState.IsMaximized {
		runtime.WindowMaximise(ctx)
	}

	// events from frontend launcher
	runtime.EventsOn(ctx, "mdeck:connect", func(data ...interface{}) {
		if len(data) == 0 {
			return
		}
		url, _ := data[0].(string)
		if url != "" {
			log.Printf("[MDeck] connect %s", url)
			runtime.WindowExecJS(ctx, fmt.Sprintf("window.location.href=%q", url))
		}
	})
	runtime.EventsOn(ctx, "mdeck:new-window", func(data ...interface{}) {
		url := ""
		if len(data) > 0 {
			url, _ = data[0].(string)
		}
		_ = a.NewWindow(url)
	})

	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		js := script.GetInjectionScript()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				runtime.WindowExecJS(ctx, js)
			}
		}
	}()
}

func (a *App) domReady(ctx context.Context) {
	runtime.WindowExecJS(ctx, script.GetInjectionScript())
}

func (a *App) beforeClose(ctx context.Context) bool {
	a.SaveWindowState()
	return false
}

func (a *App) shutdown(ctx context.Context) { a.SaveWindowState() }

func (a *App) SaveWindowState() {
	if a.ctx == nil {
		return
	}
	w, h := runtime.WindowGetSize(a.ctx)
	if w <= 0 || h <= 0 {
		return
	}
	x, y := runtime.WindowGetPosition(a.ctx)
	st := config.WindowState{X: x, Y: y, Width: w, Height: h, IsMaximized: runtime.WindowIsMaximised(a.ctx)}
	_ = config.SaveWindowState(st)
}

// bindings

func (a *App) GetUserDataDir() string { return a.userDataDir }
func (a *App) GetBookmarks() []bookmarks.Bookmark { return bookmarks.Load() }
func (a *App) SaveBookmark(b bookmarks.Bookmark) ([]bookmarks.Bookmark, error) {
	list, err := bookmarks.Upsert(b)
	if err != nil {
		return nil, err
	}
	return list, nil
}
func (a *App) DeleteBookmark(id string) ([]bookmarks.Bookmark, error) { return bookmarks.Delete(id) }
func (a *App) GetWindowState() config.WindowState { return config.LoadWindowState() }
func (a *App) OpenExternal(url string) { runtime.BrowserOpenURL(a.ctx, url) }
func (a *App) IsInternalURL(url string) bool { return navigation.IsInternalURL(url) }
func (a *App) Navigate(url string) {
	if a.ctx != nil {
		runtime.WindowExecJS(a.ctx, fmt.Sprintf("window.location.href=%q", url))
	}
}

// Multi-window: spawn new process (each window = own OS window, own profile share but separate WebView)
// For Linux/Windows we duplicate current executable with --url flag.
func (a *App) NewWindow(url string) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	args := []string{}
	if url != "" {
		args = append(args, "--url", url)
	}
	cmd := exec.Command(exe, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Start()
}

func (a *App) Close() { if a.ctx != nil { runtime.Quit(a.ctx) } }
func (a *App) Minimise() { if a.ctx != nil { runtime.WindowMinimise(a.ctx) } }
func (a *App) ToggleMaximise() { if a.ctx != nil { runtime.WindowToggleMaximise(a.ctx) } }
func (a *App) Reload() { if a.ctx != nil { runtime.WindowReload(a.ctx) } }
func (a *App) GetVersion() string { return "1.0.0" }
