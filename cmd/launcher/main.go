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

	dir := *dirFlag
	if dir == "" {
		dir = save.Discover("")
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

	appURL := "http://" + ln.Addr().String() + sitePath + "?t=" + token
	boot := fmt.Sprintf(`<!doctype html><meta charset="utf-8"><style>html,body{margin:0;height:100%%;background:#151515}</style><script>requestAnimationFrame(function(){requestAnimationFrame(function(){location.replace(%q)})})</script>`, appURL)

	mux := http.NewServeMux()
	bridge.New(token, func() string { return dir }).Register(mux)
	mux.HandleFunc("/__launch", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(boot))
	})
	mux.Handle("/", proxy)

	go func() { log.Fatalf("serve: %v", http.Serve(ln, mux)) }()

	openWindow("Chaos Gate Unlocker", "http://"+ln.Addr().String()+"/__launch")
}

func newToken() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		log.Fatalf("token: %v", err)
	}
	return hex.EncodeToString(b)
}
