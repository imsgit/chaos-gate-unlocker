//go:build !js && embedwasm

package main

import (
	"embed"
	"io/fs"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
)

//go:embed wasm
var wasmBundle embed.FS

const browserSupported = true

var (
	browserOnce sync.Once
	browserURL  string
	browserErr  error
)

func openInBrowser() error {
	browserOnce.Do(func() {
		sub, err := fs.Sub(wasmBundle, "wasm")
		if err != nil {
			browserErr = err
			return
		}

		ln, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			browserErr = err
			return
		}

		files := http.FileServer(http.FS(sub))
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasSuffix(r.URL.Path, ".wasm") {
				w.Header().Set("Content-Type", "application/wasm")
				if gz, err := fs.ReadFile(sub, strings.TrimPrefix(r.URL.Path, "/")+".gz"); err == nil {
					w.Header().Set("Content-Encoding", "gzip")
					w.Header().Set("Content-Length", strconv.Itoa(len(gz)))
					w.Write(gz)
					return
				}
			}
			files.ServeHTTP(w, r)
		})

		browserURL = "http://" + ln.Addr().String() + "/"
		go http.Serve(ln, handler)
	})

	if browserErr != nil {
		return browserErr
	}

	u, err := url.Parse(browserURL)
	if err != nil {
		return err
	}

	return fyne.CurrentApp().OpenURL(u)
}
