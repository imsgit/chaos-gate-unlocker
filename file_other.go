//go:build !js

package main

import (
	"io"

	"chaos-gate-unlocker/internal/files"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
)

func openFile(w fyne.Window, fm *files.Manager, onData func(name string, data []byte)) {
	fileDialog := dialog.NewFileOpen(func(rc fyne.URIReadCloser, _ error) {
		if rc == nil {
			return
		}

		name := rc.URI().Path()
		data, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		onData(name, data)
	}, w)

	l, _ := storage.ListerForURI(storage.NewFileURI(fm.GetCurrentPath()))
	fileDialog.SetTitleText("Open game save file ../" + fm.SaveDir())
	fileDialog.SetConfirmText("Open")
	fileDialog.SetDismissText("Cancel")
	fileDialog.SetLocation(l)
	fileDialog.Resize(fyne.NewSize(800, 600))
	fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".gksave"}))
	fileDialog.Show()
}

func saveFile(fm *files.Manager) error {
	return fm.Save()
}

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
