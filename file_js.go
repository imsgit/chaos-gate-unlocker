//go:build js

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"syscall/js"

	"chaos-gate-unlocker/internal/files"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var (
	bridgeFile string
	bridgeWin  fyne.Window
)

func bridgeToken() string {
	search := js.Global().Get("location").Get("search").String()
	q, _ := url.ParseQuery(strings.TrimPrefix(search, "?"))
	return q.Get("t")
}

func showTryOnline() bool { return bridgeToken() != "" }

func bridgeBase() string {
	return js.Global().Get("location").Get("origin").String()
}

func openFile(w fyne.Window, fm *files.Manager, onData func(name string, data []byte)) {
	bridgeWin = w
	if tok := bridgeToken(); tok != "" {
		go bridgePick(w, tok, onData)
		return
	}

	doc := js.Global().Get("document")
	input := doc.Call("createElement", "input")
	input.Set("type", "file")
	input.Set("accept", ".gksave")
	input.Get("style").Set("display", "none")

	var onChange js.Func
	onChange = js.FuncOf(func(_ js.Value, _ []js.Value) any {
		list := input.Get("files")
		if list.Length() == 0 {
			input.Call("remove")
			onChange.Release()
			return nil
		}

		file := list.Index(0)
		name := file.Get("name").String()
		reader := js.Global().Get("FileReader").New()

		var onLoad js.Func
		onLoad = js.FuncOf(func(_ js.Value, _ []js.Value) any {
			buf := js.Global().Get("Uint8Array").New(reader.Get("result"))
			data := make([]byte, buf.Length())
			js.CopyBytesToGo(data, buf)

			input.Call("remove")
			onChange.Release()
			onLoad.Release()

			fyne.Do(func() { onData(name, data) })
			return nil
		})

		reader.Set("onload", onLoad)
		reader.Call("readAsArrayBuffer", file)
		return nil
	})

	input.Set("onchange", onChange)
	doc.Get("body").Call("appendChild", input)
	input.Call("click")
}

func bridgeGet(u string) ([]byte, error) {
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("\n\n\nLauncher bridge error (%s).\n\n", resp.Status)
	}
	return body, nil
}

func bridgePick(w fyne.Window, tok string, onData func(name string, data []byte)) {
	fail := func(err error) { fyne.Do(func() { dialog.ShowError(err, w) }) }

	body, err := bridgeGet(bridgeBase() + "/api/list?t=" + url.QueryEscape(tok))
	if err != nil {
		fail(err)
		return
	}
	var list []struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(body, &list); err != nil {
		fail(err)
		return
	}
	if len(list) == 0 {
		fail(errors.New("\n\n\nNo .gksave files found in the save folder.\n\n"))
		return
	}

	names := make([]string, len(list))
	for i, e := range list {
		names[i] = e.Name
	}
	fyne.Do(func() { showBridgePicker(w, tok, names, onData) })
}

func showBridgePicker(w fyne.Window, tok string, names []string, onData func(name string, data []byte)) {
	listw := widget.NewList(
		func() int { return len(names) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i widget.ListItemID, o fyne.CanvasObject) { o.(*widget.Label).SetText(names[i]) },
	)

	d := dialog.NewCustom("Open game save file", "Cancel", listw, w)
	d.Resize(fyne.NewSize(420, 420))

	listw.OnSelected = func(i widget.ListItemID) {
		name := names[i]
		d.Hide()
		go func() {
			data, err := bridgeGet(bridgeBase() + "/api/file?t=" + url.QueryEscape(tok) + "&name=" + url.QueryEscape(name))
			if err != nil {
				fyne.Do(func() { dialog.ShowError(err, w) })
				return
			}
			bridgeFile = name
			fyne.Do(func() { onData(name, data) })
		}()
	}

	d.Show()
}

func saveFile(fm *files.Manager) error {
	if tok := bridgeToken(); tok != "" {
		return bridgeSave(tok, fm)
	}

	data, err := fm.Encode()
	if err != nil {
		return err
	}
	download(fm.Name(), data)
	return nil
}

func bridgeSave(tok string, fm *files.Manager) error {
	data, err := fm.Encode()
	if err != nil {
		return err
	}
	name := bridgeFile
	if name == "" {
		name = fm.Name()
	}

	go func() {
		u := bridgeBase() + "/api/file?t=" + url.QueryEscape(tok) + "&name=" + url.QueryEscape(name)
		resp, err := http.Post(u, "application/octet-stream", bytes.NewReader(data))
		if err == nil {
			if resp.StatusCode >= 300 {
				err = fmt.Errorf("\n\n\nError. Cannot save file (%s).\n\n", resp.Status)
			}
			resp.Body.Close()
		}
		if err != nil {
			fyne.Do(func() { dialog.ShowError(err, bridgeWin) })
		}
	}()

	return nil
}

func confirmSave(w fyne.Window, do func()) {
	if bridgeToken() == "" {
		do()
		return
	}

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

func download(name string, data []byte) {
	buf := js.Global().Get("Uint8Array").New(len(data))
	js.CopyBytesToJS(buf, data)

	parts := js.Global().Get("Array").New()
	parts.Call("push", buf)
	blob := js.Global().Get("Blob").New(parts, map[string]any{"type": "application/octet-stream"})

	objURL := js.Global().Get("URL").Call("createObjectURL", blob)
	defer js.Global().Get("URL").Call("revokeObjectURL", objURL)

	doc := js.Global().Get("document")
	a := doc.Call("createElement", "a")
	a.Set("href", objURL)
	a.Set("download", name)
	doc.Get("body").Call("appendChild", a)
	a.Call("click")
	a.Call("remove")
}

func validateScale() {}
