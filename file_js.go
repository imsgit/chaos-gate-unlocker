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
	"sync"
	"syscall/js"

	"chaos-gate-unlocker/internal/save"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

var bridgeEnv = sync.OnceValues(func() (token, base string) {
	loc := js.Global().Get("location")
	q, _ := url.ParseQuery(strings.TrimPrefix(loc.Get("search").String(), "?"))
	return q.Get("t"), loc.Get("origin").String()
})

func bridgeToken() string {
	tok, _ := bridgeEnv()
	return tok
}

func bridgeURL(path string, q url.Values) string {
	tok, base := bridgeEnv()
	if q == nil {
		q = url.Values{}
	}
	q.Set("t", tok)
	return base + path + "?" + q.Encode()
}

func showTryOnline() bool { return bridgeToken() != "" }

func openWebsite(u *url.URL) {
	if bridgeToken() != "" {
		go bridgeGet(bridgeURL("/api/openurl", url.Values{"url": {u.String()}}))
		return
	}
	_ = fyne.CurrentApp().OpenURL(u)
}

func openFile(w fyne.Window, beginLoad func(), onData func(name string, data []byte, err error)) {
	if bridgeToken() != "" {
		go bridgePick(w, beginLoad, onData)
		return
	}

	doc := js.Global().Get("document")
	input := doc.Call("createElement", "input")
	input.Set("type", "file")
	input.Set("accept", ".gksave")
	input.Get("style").Set("display", "none")

	var onChange, onCancel js.Func
	cleanup := func() {
		input.Call("remove")
		onChange.Release()
		onCancel.Release()
	}
	onChange = js.FuncOf(func(_ js.Value, _ []js.Value) any {
		list := input.Get("files")
		if list.Length() == 0 {
			cleanup()
			return nil
		}

		file := list.Index(0)
		name := file.Get("name").String()
		reader := js.Global().Get("FileReader").New()

		fyne.Do(beginLoad)

		var onLoad, onError js.Func
		done := func(data []byte, err error) {
			cleanup()
			onLoad.Release()
			onError.Release()
			go onData(name, data, err)
		}
		onLoad = js.FuncOf(func(_ js.Value, _ []js.Value) any {
			buf := js.Global().Get("Uint8Array").New(reader.Get("result"))
			data := make([]byte, buf.Length())
			js.CopyBytesToGo(data, buf)
			done(data, nil)
			return nil
		})
		onError = js.FuncOf(func(_ js.Value, _ []js.Value) any {
			done(nil, errors.New("\n\n\nError. Cannot read the selected file.\n\n"))
			return nil
		})

		reader.Set("onload", onLoad)
		reader.Set("onerror", onError)
		reader.Call("readAsArrayBuffer", file)
		return nil
	})
	onCancel = js.FuncOf(func(_ js.Value, _ []js.Value) any {
		cleanup()
		return nil
	})

	input.Set("onchange", onChange)
	input.Call("addEventListener", "cancel", onCancel)
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

func bridgePick(w fyne.Window, beginLoad func(), onData func(name string, data []byte, err error)) {
	fail := func(err error) { fyne.Do(func() { dialog.ShowError(err, w) }) }

	body, err := bridgeGet(bridgeURL("/api/list", nil))
	if err != nil {
		fail(err)
		return
	}
	var list []struct {
		Name string `json:"name"`
		save.Info
	}
	if err := json.Unmarshal(body, &list); err != nil {
		fail(err)
		return
	}
	if len(list) == 0 {
		fail(errors.New("\n\n\nNo .gksave files found in the save folder.\n\n"))
		return
	}

	infos := make(map[string]save.Info, len(list))
	for _, e := range list {
		infos[e.Name] = e.Info
	}
	fyne.Do(func() {
		showSavePicker(w, infos, func(name string) {
			beginLoad()
			go func() {
				data, err := bridgeGet(bridgeURL("/api/file", url.Values{"name": {name}}))
				onData(name, data, err)
			}()
		}, func() {
			go bridgeGet(bridgeURL("/api/open", nil))
		})
	})
}

func saveFile() error {
	data, err := filesManager.Encode()
	if err != nil {
		return err
	}
	if bridgeToken() == "" {
		download(filesManager.Name(), data)
		return nil
	}
	resp, err := http.Post(bridgeURL("/api/file", url.Values{"name": {filesManager.Name()}}), "application/octet-stream", bytes.NewReader(data))
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("\n\n\nError. Cannot save file (%s).\n\n", resp.Status)
	}
	return nil
}

func confirmSave(w fyne.Window, do func()) {
	if bridgeToken() == "" {
		do()
		return
	}
	showSaveConfirm(w, do)
}

func download(name string, data []byte) {
	buf := js.Global().Get("Uint8Array").New(len(data))
	js.CopyBytesToJS(buf, data)

	blob := js.Global().Get("Blob").New([]any{buf}, map[string]any{"type": "application/octet-stream"})

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
