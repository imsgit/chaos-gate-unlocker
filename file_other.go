//go:build !js

package main

import (
	"errors"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"chaos-gate-unlocker/internal/display"
	"chaos-gate-unlocker/internal/files"
	"chaos-gate-unlocker/internal/saveinfo"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func openFile(w fyne.Window, fm *files.Manager, onData func(name string, data []byte)) {
	dir := fm.GetCurrentPath()

	entries, err := os.ReadDir(dir)
	if err != nil {
		dialog.ShowError(err, w)
		return
	}

	var names []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".gksave") {
			names = append(names, e.Name())
		}
	}
	if len(names) == 0 {
		dialog.ShowError(errors.New("\n\n\nNo .gksave files found in the save folder.\n\n"), w)
		return
	}

	infoCache := map[string]saveinfo.Info{}
	info := func(name string) saveinfo.Info {
		if v, ok := infoCache[name]; ok {
			return v
		}
		v := saveinfo.ParseFile(filepath.Join(dir, name))
		infoCache[name] = v
		return v
	}

	showSavePicker(w, names, info, func(name string) {
		path := filepath.Join(dir, name)
		data, err := os.ReadFile(path)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		onData(path, data)
	}, func() {
		openSaveDir(dir)
	})
}

func openSaveDir(dir string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", dir)
	case "darwin":
		cmd = exec.Command("open", dir)
	default:
		cmd = exec.Command("xdg-open", dir)
	}
	_ = cmd.Start()
}

func saveFile(fm *files.Manager) error {
	return fm.Save()
}

func showTryOnline() bool { return true }

func openWebsite(u *url.URL) { _ = fyne.CurrentApp().OpenURL(u) }

func confirmSave(w fyne.Window, do func()) {
	showSaveConfirm(w, do)
}

func validateScale() {
	if runtime.GOOS == "windows" {
		return
	}
	if display.IsHiDPI() {
		os.Setenv("FYNE_SCALE", "2.0")
	}
}
