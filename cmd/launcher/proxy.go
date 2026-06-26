package main

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

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

	return &siteProxy{
		origin:   u.Scheme + "://" + u.Host,
		client:   &http.Client{Transport: &http.Transport{DisableCompression: true}},
		cacheDir: cacheDir,
	}, path, nil
}

func (p *siteProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	q.Del("t")
	upstream := p.origin + r.URL.Path
	if enc := q.Encode(); enc != "" {
		upstream += "?" + enc
	}

	resp, err := p.client.Get(upstream)
	if err == nil && resp.StatusCode == http.StatusOK {
		body, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr == nil {
			ct := resp.Header.Get("Content-Type")
			ce := resp.Header.Get("Content-Encoding")
			p.store(r.URL.Path, body, ct, ce)
			serveAsset(w, body, ct, ce)
			return
		}
	} else if resp != nil {
		resp.Body.Close()
	}

	if body, ct, ce, ok := p.load(r.URL.Path); ok {
		serveAsset(w, body, ct, ce)
		return
	}

	http.Error(w, "site unavailable and not cached", http.StatusBadGateway)
}

func serveAsset(w http.ResponseWriter, body []byte, contentType, contentEncoding string) {
	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}
	if contentEncoding != "" {
		w.Header().Set("Content-Encoding", contentEncoding)
	}
	w.Write(body)
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

func (p *siteProxy) store(reqPath string, body []byte, contentType, contentEncoding string) {
	dst := p.cachePath(reqPath)
	if dst == "" {
		return
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return
	}
	_ = os.WriteFile(dst, body, 0o644)
	_ = os.WriteFile(dst+".hdr", []byte(contentType+"\n"+contentEncoding), 0o644)
}

func (p *siteProxy) load(reqPath string) (body []byte, contentType, contentEncoding string, ok bool) {
	src := p.cachePath(reqPath)
	if src == "" {
		return nil, "", "", false
	}
	body, err := os.ReadFile(src)
	if err != nil {
		return nil, "", "", false
	}
	hdr, _ := os.ReadFile(src + ".hdr")
	ct, ce, _ := strings.Cut(string(hdr), "\n")
	return body, ct, ce, true
}
