package main

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"chaos-gate-unlocker/internal/bridge"
	"chaos-gate-unlocker/internal/savedir"
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

	dir := *dirFlag
	if dir == "" {
		dir = savedir.Discover("")
	}

	token := newToken()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	proxy, sitePath, err := newSiteProxy(*site)
	if err != nil {
		log.Fatalf("proxy: %v", err)
	}

	boot := fmt.Sprintf(`<!doctype html><meta charset="utf-8"><style>html,body{margin:0;height:100%%;background:#151515}@media(prefers-color-scheme:light){html,body{background:#fff}}</style><script>location.replace(%q)</script>`, sitePath+"?t="+token)

	mux := http.NewServeMux()
	bridge.New(token, func() string { return dir }).Register(mux)
	mux.HandleFunc("/__boot", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(boot))
	})
	mux.Handle("/", proxy)

	go func() { log.Fatalf("serve: %v", http.Serve(ln, mux)) }()

	openWindow("Chaos Gate Unlocker", "http://"+ln.Addr().String()+"/__boot")
}

func newToken() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		log.Fatalf("token: %v", err)
	}
	return hex.EncodeToString(b)
}
