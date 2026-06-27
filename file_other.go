//go:build !js

package main

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"chaos-gate-unlocker/internal/files"

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

	showSavePicker(w, names, func(name string) {
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

func confirmSave(w fyne.Window, do func()) {
	d := dialog.NewConfirm(
		"Save confirmation",
		"\n\n\nThis will override the existing save file. Are you sure?\nPlease make a backup if needed.",
		func(r bool) {
			if r {
				do()
			}
		}, w)
	d.SetConfirmText("Save")
	d.SetDismissText("Cancel")
	d.Show()
}

func validateScale() {
	if runtime.GOOS == "windows" {
		return
	}
	cmd := exec.Command("xdpyinfo")
	out, err := cmd.Output()
	if err != nil {
		return
	}
	re := regexp.MustCompile(`resolution:\s+(\d+)x`)
	match := re.FindStringSubmatch(string(out))
	if len(match) == 2 {
		if dpi, _ := strconv.Atoi(match[1]); dpi > 96 {
			os.Setenv("FYNE_SCALE", "2.0")
		}
	}
}
