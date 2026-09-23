package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"chaos-gate-unlocker/internal/bridge"
)

const (
	cacheMagic    = "cgu1"
	stallTimeout  = 10 * time.Second
	maxAssetBytes = 32 << 20
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
	inflight sync.Map
}

func newSiteProxy(site string) (*siteProxy, string, error) {
	u, err := url.Parse(site)
	if err != nil {
		return nil, "", err
	}
	if u.Scheme != "https" && u.Hostname() != "localhost" && u.Hostname() != "127.0.0.1" {
		return nil, "", errors.New("site URL must use https")
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

	if cached, ok := p.load(reqPath); ok {
		serveAsset(w, cached)
		p.revalidate(reqPath, cached.etag, cached.lastModified)
		return
	}

	if !p.fetch(reqPath, "", "", w) {
		http.Error(w, "site unavailable and not cached", http.StatusBadGateway)
	}
}

func (p *siteProxy) revalidate(reqPath, etag, lastMod string) {
	if _, busy := p.inflight.LoadOrStore(reqPath, struct{}{}); busy {
		return
	}
	go func() {
		defer p.inflight.Delete(reqPath)
		p.fetch(reqPath, etag, lastMod, nil)
	}()
}

func (p *siteProxy) fetch(reqPath, etag, lastMod string, w http.ResponseWriter) bool {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.origin+reqPath, nil)
	if err != nil {
		return false
	}
	req.Header.Set("Accept-Encoding", "gzip")
	if etag != "" {
		req.Header.Set("If-None-Match", etag)
	}
	if lastMod != "" {
		req.Header.Set("If-Modified-Since", lastMod)
	}
	resp, err := p.client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return false
	}
	a := asset{
		contentType:     resp.Header.Get("Content-Type"),
		contentEncoding: resp.Header.Get("Content-Encoding"),
		etag:            resp.Header.Get("ETag"),
		lastModified:    resp.Header.Get("Last-Modified"),
	}

	var buf bytes.Buffer
	if resp.ContentLength > 0 && resp.ContentLength <= maxAssetBytes {
		buf.Grow(int(resp.ContentLength))
	}
	var dst io.Writer = &buf
	if w != nil {
		setAssetHeaders(w, a)
		if resp.ContentLength >= 0 {
			w.Header().Set("Content-Length", strconv.FormatInt(resp.ContentLength, 10))
		}
		dst = io.MultiWriter(&buf, w)
	}

	stall := time.AfterFunc(stallTimeout, cancel)
	defer stall.Stop()
	n, err := io.Copy(dst, io.LimitReader(&stallReader{r: resp.Body, timer: stall}, maxAssetBytes+1))
	if err != nil || n > maxAssetBytes {
		return w != nil
	}
	if w != nil {
		http.NewResponseController(w).Flush()
	}
	a.body = buf.Bytes()
	p.store(reqPath, a)
	return true
}

type stallReader struct {
	r     io.Reader
	timer *time.Timer
}

func (s *stallReader) Read(p []byte) (int, error) {
	n, err := s.r.Read(p)
	if n > 0 {
		s.timer.Reset(stallTimeout)
	}
	return n, err
}

func setAssetHeaders(w http.ResponseWriter, a asset) {
	if a.contentType != "" {
		w.Header().Set("Content-Type", a.contentType)
	}
	if a.contentEncoding != "" {
		w.Header().Set("Content-Encoding", a.contentEncoding)
	}
	w.Header().Set("Cache-Control", "no-store")
}

func serveAsset(w http.ResponseWriter, a asset) {
	setAssetHeaders(w, a)
	w.Header().Set("Content-Length", strconv.Itoa(len(a.body)))
	w.Write(a.body)
}

func (p *siteProxy) cachePath(reqPath string) string {
	clean := filepath.FromSlash(strings.TrimPrefix(reqPath, "/"))
	if clean == "" || strings.HasSuffix(reqPath, "/") {
		clean = filepath.Join(clean, "index.html")
	}
	if !filepath.IsLocal(clean) {
		return ""
	}
	return filepath.Join(p.cacheDir, clean)
}

func (p *siteProxy) store(reqPath string, a asset) {
	dst := p.cachePath(reqPath)
	if dst == "" {
		return
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return
	}
	header := strings.Join([]string{cacheMagic, a.contentType, a.contentEncoding, a.etag, a.lastModified}, "\n") + "\n"
	_ = bridge.WriteFileAtomic(dst, []byte(header), a.body)
}

func (p *siteProxy) load(reqPath string) (asset, bool) {
	src := p.cachePath(reqPath)
	if src == "" {
		return asset{}, false
	}
	data, err := os.ReadFile(src)
	if err != nil {
		return asset{}, false
	}
	f := bytes.SplitN(data, []byte("\n"), 6)
	if len(f) < 6 || string(f[0]) != cacheMagic {
		return asset{}, false
	}
	return asset{
		contentType:     string(f[1]),
		contentEncoding: string(f[2]),
		etag:            string(f[3]),
		lastModified:    string(f[4]),
		body:            f[5],
	}, true
}
