package main

import (
	"embed"
	"log"
	"os"
	"runtime/debug"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/src
var assets embed.FS

//go:embed build/appicon.png
var appIcon []byte

func main() {
	// disk cache like Gotion: 300MB cache, 100MB media
	browserArgs := "--disk-cache-size=314572800 --media-cache-size=104857600 --disable-features=Translate,OptimizationHints,MediaRouter"
	_ = os.Setenv("WEBVIEW2_ADDITIONAL_BROWSER_ARGUMENTS", browserArgs)
	_ = os.Setenv("WEBKIT_DISABLE_DMABUF_RENDERER", "1")
	_ = os.Setenv("WEBKIT_DISABLE_COMPOSITING_MODE", "1")

	debug.SetGCPercent(30)

	app := NewApp()

	windowsOptions := &windows.Options{
		WebviewUserDataPath:               app.GetUserDataDir(),
		Theme:                             windows.SystemDefault,
		BackdropType:                      windows.Auto,
		IsZoomControlEnabled:              true,
		DisablePinchZoom:                  false,
		DisableFramelessWindowDecorations: false,
		CustomTheme: &windows.ThemeSettings{
			DarkModeTitleBar:          windows.RGB(14, 16, 21),
			DarkModeTitleBarInactive:  windows.RGB(18, 20, 26),
			DarkModeTitleText:         windows.RGB(240, 240, 240),
			DarkModeTitleTextInactive: windows.RGB(140, 140, 140),
			DarkModeBorder:            windows.RGB(35, 35, 35),
			DarkModeBorderInactive:    windows.RGB(30, 30, 30),
		},
	}

	linuxOptions := &linux.Options{
		Icon:             appIcon,
		WebviewGpuPolicy: linux.WebviewGpuPolicyNever,
		ProgramName:      "mdeck-desktop",
	}

	err := wails.Run(&options.App{
		Title:            "MDeck",
		Width:            app.windowState.Width,
		Height:           app.windowState.Height,
		MinWidth:         800,
		MinHeight:        560,
		Frameless:        true,
		CSSDragProperty:  "--wails-draggable",
		CSSDragValue:     "drag",
		DisableResize:    false,
		EnableDefaultContextMenu: false,
		BackgroundColour: &options.RGBA{R: 9, G: 10, B: 13, A: 255},
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup:     app.startup,
		OnDomReady:    app.domReady,
		OnBeforeClose: app.beforeClose,
		OnShutdown:    app.shutdown,
		Windows:       windowsOptions,
		Linux:         linuxOptions,
		Bind: []interface{}{app},
	})
	if err != nil {
		log.Fatalf("[MDeck] %v", err)
	}
}
