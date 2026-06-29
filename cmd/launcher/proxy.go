package main

import (
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type asset struct {
	body            []byte
	contentType     string
	contentEncoding string
	etag            string
	lastModified    string
}

type siteProxy struct {
	origin   string
	client   *http.Client
	cacheDir string
}

func newSiteProxy(site string) (*siteProxy, string, error) {
	u, err := url.Parse(site)
	if err != nil {
		return nil, "", err
	}
	path := u.Path
	if path == "" {
		path = "/"
	}

	cacheRoot, err := os.UserCacheDir()
	if err != nil {
		cacheRoot = os.TempDir()
	}
	cacheDir := filepath.Join(cacheRoot, "chaos-gate-unlocker", "site")
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return nil, "", err
	}

	transport := &http.Transport{
		DisableCompression:    true,
		DialContext:           (&net.Dialer{Timeout: 5 * time.Second}).DialContext,
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,
	}

	return &siteProxy{
		origin:   u.Scheme + "://" + u.Host,
		client:   &http.Client{Transport: transport},
		cacheDir: cacheDir,
	}, path, nil
}

func (p *siteProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqPath := r.URL.Path
	q := r.URL.Query()
	q.Del("t")

	if cached, ok := p.load(reqPath); ok {
		serveAsset(w, cached)
		go p.fetch(reqPath, q, cached.etag, cached.lastModified)
		return
	}

	if fresh, ok := p.fetch(reqPath, q, "", ""); ok {
		serveAsset(w, fresh)
		return
	}

	http.Error(w, "site unavailable and not cached", http.StatusBadGateway)
}

func (p *siteProxy) fetch(reqPath string, q url.Values, etag, lastMod string) (asset, bool) {
	upstream := p.origin + reqPath
	if enc := q.Encode(); enc != "" {
		upstream += "?" + enc
	}
	req, err := http.NewRequest(http.MethodGet, upstream, nil)
	if err != nil {
		return asset{}, false
	}
	if etag != "" {
		req.Header.Set("If-None-Match", etag)
	}
	if lastMod != "" {
		req.Header.Set("If-Modified-Since", lastMod)
	}
	resp, err := p.client.Do(req)
	if err != nil {
		return asset{}, false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return asset{}, false
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return asset{}, false
	}
	a := asset{
		body:            body,
		contentType:     resp.Header.Get("Content-Type"),
		contentEncoding: resp.Header.Get("Content-Encoding"),
		etag:            resp.Header.Get("ETag"),
		lastModified:    resp.Header.Get("Last-Modified"),
	}
	p.store(reqPath, a)
	return a, true
}

func serveAsset(w http.ResponseWriter, a asset) {
	if a.contentType != "" {
		w.Header().Set("Content-Type", a.contentType)
	}
	if a.contentEncoding != "" {
		w.Header().Set("Content-Encoding", a.contentEncoding)
	}
	w.Write(a.body)
}

func (p *siteProxy) cachePath(reqPath string) string {
	clean := filepath.FromSlash(strings.TrimPrefix(reqPath, "/"))
	if clean == "" || strings.HasSuffix(reqPath, "/") {
		clean = filepath.Join(clean, "index.html")
	}
	full := filepath.Join(p.cacheDir, clean)
	if rel, err := filepath.Rel(p.cacheDir, full); err != nil || strings.HasPrefix(rel, "..") {
		return ""
	}
	return full
}

func (p *siteProxy) store(reqPath string, a asset) {
	dst := p.cachePath(reqPath)
	if dst == "" {
		return
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return
	}
	_ = os.WriteFile(dst, a.body, 0o644)
	hdr := strings.Join([]string{a.contentType, a.contentEncoding, a.etag, a.lastModified}, "\n")
	_ = os.WriteFile(dst+".hdr", []byte(hdr), 0o644)
}

func (p *siteProxy) load(reqPath string) (asset, bool) {
	src := p.cachePath(reqPath)
	if src == "" {
		return asset{}, false
	}
	body, err := os.ReadFile(src)
	if err != nil {
		return asset{}, false
	}
	hdr, err := os.ReadFile(src + ".hdr")
	if err != nil {
		return asset{}, false
	}
	f := strings.SplitN(string(hdr), "\n", 4)
	for len(f) < 4 {
		f = append(f, "")
	}
	return asset{body: body, contentType: f[0], contentEncoding: f[1], etag: f[2], lastModified: f[3]}, true
}
