package stopwatchapp

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	astilectron "github.com/asticode/go-astilectron"
)

// Bootstrap unpacks required assets and creates the application GUI window
func (a *App) Bootstrap() error {
	var devTools = DevTools()
	var err error
	var hadCrash bool

	// If asset dir isn't set, try to determine it.
	if a.assetDir == "" {
		if UseTemp() {
			a.assetDir, err = TmpDataDir()
		} else {
			a.assetDir, err = PersistentDataDir()
		}
		if err != nil {
			return err
		}
	}

	// Unpack assets
	if err = UnpackEmbeddedAssets(a.assetDir, a.assetRestore); err != nil {
		return err
	}

	// Initialize astilectron
	a.Renderer, err = astilectron.New(astilectron.Options{
		AppName:            Name(),
		AppIconDefaultPath: path.Join(a.assetDir, EmbeddedIconPath()),
		AppIconDarwinPath:  path.Join(a.assetDir, EmbeddedIconPath()),
		BaseDirectoryPath:  a.assetDir,
	})
	if err != nil {
		return err
	}

	// Set provisioning to load from binary data.
	a.Renderer.SetProvisioner(
		astilectron.NewDisembedderProvisioner(
			a.assetData,
			EmbeddedAstilectronPath(),
			EmbeddedElectronPath(),
		),
	)

	// Set handling for signals and capture crashes
	a.Renderer.HandleSignals()
	a.Renderer.On(astilectron.EventNameAppCrash, func(e astilectron.Event) (deleteListener bool) {
		log.Printf("[EXIT] GUI has exited")
		log.Printf("[EXIT] event: %s, message: %v", e.Name, e.Message)
		hadCrash = true
		return true
	})

	// Start astilectron
	if err = a.Renderer.Start(); err != nil {
		log.Fatal(err)
	}

	if hadCrash {
		log.Fatal("Crashed during init.")
	}

	// Create main window
	wd := WindowSize()
	if a.Window, err = a.Renderer.NewWindow(path.Join(a.assetDir, EmbeddedUIMountPoint()), &astilectron.WindowOptions{
		Center:    astilectron.PtrBool(true),
		Height:    astilectron.PtrInt(wd.Height),
		Width:     astilectron.PtrInt(wd.Width),
		MinHeight: astilectron.PtrInt(wd.Height),
		MinWidth:  astilectron.PtrInt(wd.Width),
		Title:     astilectron.PtrStr(Name()),
		Icon:      astilectron.PtrStr(path.Join(a.assetDir, EmbeddedIconPath())),
		WebPreferences: &astilectron.WebPreferences{
			DevTools:        &devTools,
			DefaultEncoding: astilectron.PtrStr("utf-8"),
			Webaudio:        astilectron.PtrBool(false),
		},
	}); err != nil {
		log.Fatal(err)
	}

	// Setup queue for message sending
	a.msgQueue = make(chan Message, 50)

	// Run routine for handling sending
	go a.messageQueueFlusher()

	// Setup handler for GUI messages
	a.Window.OnMessage(a.onWindowMessage())

	// Actually create the window to make it appear.
	a.Window.Create()

	// Open dev tools
	if devTools {
		a.Window.OpenDevTools()
	}

	// Clean vendor directory of unnecessary zip-files.
	if !UseTemp() {
		files, err := ioutil.ReadDir(filepath.Join(a.assetDir, "vendor/"))
		if err == nil {
			for _, finfo := range files {
				if !finfo.IsDir() && strings.HasSuffix(finfo.Name(), ".zip") {
					log.Printf("Removing unnecessary resource: %s", filepath.Join(a.assetDir, "vendor/", finfo.Name()))
					os.Remove(filepath.Join(a.assetDir, "vendor/", finfo.Name()))
				}
			}
		} else {
			log.Printf("error reading dir: %s", err)
		}
	}
	return nil
}
