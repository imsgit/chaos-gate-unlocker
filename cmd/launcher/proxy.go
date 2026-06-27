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

	ct, ce, etag, lastMod, haveCache := p.loadMeta(r.URL.Path)

	if req, err := http.NewRequest(http.MethodGet, upstream, nil); err == nil {
		if haveCache {
			if etag != "" {
				req.Header.Set("If-None-Match", etag)
			}
			if lastMod != "" {
				req.Header.Set("If-Modified-Since", lastMod)
			}
		}
		if resp, derr := p.client.Do(req); derr == nil {
			if resp.StatusCode == http.StatusOK {
				body, readErr := io.ReadAll(resp.Body)
				resp.Body.Close()
				if readErr == nil {
					nct := resp.Header.Get("Content-Type")
					nce := resp.Header.Get("Content-Encoding")
					p.store(r.URL.Path, body, nct, nce,
						resp.Header.Get("ETag"), resp.Header.Get("Last-Modified"))
					serveAsset(w, body, nct, nce)
					return
				}
			} else {
				notModified := resp.StatusCode == http.StatusNotModified
				resp.Body.Close()
				if notModified && haveCache {
					if body, ok := p.loadBody(r.URL.Path); ok {
						serveAsset(w, body, ct, ce)
						return
					}
				}
			}
		}
	}

	if haveCache {
		if body, ok := p.loadBody(r.URL.Path); ok {
			serveAsset(w, body, ct, ce)
			return
		}
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

func (p *siteProxy) store(reqPath string, body []byte, contentType, contentEncoding, etag, lastModified string) {
	dst := p.cachePath(reqPath)
	if dst == "" {
		return
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return
	}
	_ = os.WriteFile(dst, body, 0o644)
	hdr := strings.Join([]string{contentType, contentEncoding, etag, lastModified}, "\n")
	_ = os.WriteFile(dst+".hdr", []byte(hdr), 0o644)
}

func (p *siteProxy) loadMeta(reqPath string) (contentType, contentEncoding, etag, lastModified string, ok bool) {
	src := p.cachePath(reqPath)
	if src == "" {
		return "", "", "", "", false
	}
	if _, err := os.Stat(src); err != nil {
		return "", "", "", "", false
	}
	hdr, err := os.ReadFile(src + ".hdr")
	if err != nil {
		return "", "", "", "", false
	}
	f := strings.SplitN(string(hdr), "\n", 4)
	for len(f) < 4 {
		f = append(f, "")
	}
	return f[0], f[1], f[2], f[3], true
}

func (p *siteProxy) loadBody(reqPath string) (body []byte, ok bool) {
	src := p.cachePath(reqPath)
	if src == "" {
		return nil, false
	}
	body, err := os.ReadFile(src)
	if err != nil {
		return nil, false
	}
	return body, true
}
