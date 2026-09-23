//go:build !js

package main

import (
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"chaos-gate-unlocker/internal/display"
	"chaos-gate-unlocker/internal/save"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func openFile(w fyne.Window, beginLoad func(), onData func(name string, data []byte, err error)) {
	go func() {
		dir := filesManager.GetCurrentPath()
		entries, err := os.ReadDir(dir)

		infos := map[string]save.Info{}
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".gksave") {
				infos[e.Name()] = save.ParseFile(filepath.Join(dir, e.Name()))
			}
		}

		fyne.Do(func() {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if len(infos) == 0 {
				dialog.ShowError(errors.New("\n\n\nNo .gksave files found in the save folder.\n\n"), w)
				return
			}

			showSavePicker(w, infos, func(name string) {
				beginLoad()
				path := filepath.Join(dir, name)
				go func() {
					data, err := os.ReadFile(path)
					onData(path, data, err)
				}()
			}, func() {
				_ = fyne.CurrentApp().OpenURL(fileURL(dir))
			})
		})
	}()
}

func fileURL(dir string) *url.URL {
	p := filepath.ToSlash(dir)
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	return &url.URL{Scheme: "file", Path: p}
}

func saveFile() error { return filesManager.Save() }

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
