//go:build js

package main

import (
	"syscall/js"

	"chaos-gate-unlocker/internal/files"

	"fyne.io/fyne/v2"
)

func openFile(_ fyne.Window, _ *files.Manager, onData func(name string, data []byte)) {
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

func saveFile(fm *files.Manager) error {
	data, err := fm.Encode()
	if err != nil {
		return err
	}
	download(fm.Name(), data)
	return nil
}

func confirmSave(_ fyne.Window, do func()) { do() }

func download(name string, data []byte) {
	buf := js.Global().Get("Uint8Array").New(len(data))
	js.CopyBytesToJS(buf, data)

	parts := js.Global().Get("Array").New()
	parts.Call("push", buf)
	blob := js.Global().Get("Blob").New(parts, map[string]any{"type": "application/octet-stream"})

	url := js.Global().Get("URL").Call("createObjectURL", blob)
	defer js.Global().Get("URL").Call("revokeObjectURL", url)

	doc := js.Global().Get("document")
	a := doc.Call("createElement", "a")
	a.Set("href", url)
	a.Set("download", name)
	doc.Get("body").Call("appendChild", a)
	a.Call("click")
	a.Call("remove")
}
