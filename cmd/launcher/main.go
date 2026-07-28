package main

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"chaos-gate-unlocker/internal/bridge"
	"chaos-gate-unlocker/internal/save"
)

var version = "dev"

const siteURL = "https://imsgit.github.io/chaos-gate-unlocker/app.html"

func main() {
	site := flag.String("url", siteURL, "site URL to load the wasm from")
	dirFlag := flag.String("dir", "", "save directory (default: auto-detect Steam/Proton)")
	showVer := flag.Bool("version", false, "print launcher version and exit")
	flag.Parse()

	if *showVer {
		fmt.Println(version)
		return
	}

	dir := sync.OnceValue(func() string {
		if *dirFlag != "" {
			return *dirFlag
		}
		return save.Discover("")
	})

	token := newToken()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	proxy, sitePath, err := newSiteProxy(*site)
	if err != nil {
		log.Fatalf("proxy: %v", err)
	}

	host := ln.Addr().String()
	appURL := "http://" + host + sitePath + "?t=" + token
	boot := fmt.Sprintf(`<!doctype html><meta charset="utf-8"><style>html,body{margin:0;height:100%%;background:#151515}</style><script>requestAnimationFrame(function(){requestAnimationFrame(function(){location.replace(%q)})})</script>`, appURL)

	mux := http.NewServeMux()
	bridge.New(token, dir).Register(mux)
	var launched atomic.Bool
	mux.HandleFunc("/__launch", func(w http.ResponseWriter, r *http.Request) {
		if subtle.ConstantTimeCompare([]byte(r.URL.Query().Get("t")), []byte(token)) != 1 {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		if !launched.CompareAndSwap(false, true) {
			http.Error(w, "gone", http.StatusGone)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(boot))
	})
	mux.Handle("/", proxy)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Host != host {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		mux.ServeHTTP(w, r)
	})

	srv := &http.Server{Handler: handler, ReadHeaderTimeout: 10 * time.Second}
	go func() {
		if err := srv.Serve(ln); err != http.ErrServerClosed {
			log.Fatalf("serve: %v", err)
		}
	}()

	openWindow("Chaos Gate Unlocker", "http://"+host+"/__launch?t="+token)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}

func newToken() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		log.Fatalf("token: %v", err)
	}
	return hex.EncodeToString(b)
}
