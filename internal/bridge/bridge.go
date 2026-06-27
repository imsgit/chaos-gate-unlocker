package bridge

import (
	"crypto/subtle"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"chaos-gate-unlocker/internal/saveinfo"
)

const maxSaveBytes = 64 << 20

type Handler struct {
	token string
	dir   func() string
}

func New(token string, dir func() string) *Handler {
	return &Handler{token: token, dir: dir}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/api/list", h.list)
	mux.HandleFunc("/api/file", h.file)
	mux.HandleFunc("/api/open", h.open)
	mux.HandleFunc("/api/openurl", h.openURL)
}

func (h *Handler) authed(r *http.Request) bool {
	got := r.URL.Query().Get("t")
	return h.token != "" && subtle.ConstantTimeCompare([]byte(got), []byte(h.token)) == 1
}

func (h *Handler) resolve(name string) string {
	if name == "" || filepath.Base(name) != name || !strings.HasSuffix(name, ".gksave") {
		return ""
	}
	return filepath.Join(h.dir(), name)
}

type entry struct {
	Name    string `json:"name"`
	ModTime int64  `json:"modTime"`
	Title   string `json:"title"`
	Detail  string `json:"detail"`
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	if !h.authed(r) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	ents, err := os.ReadDir(h.dir())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out := make([]entry, 0, len(ents))
	for _, e := range ents {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".gksave") {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		si := saveinfo.ParseFile(filepath.Join(h.dir(), e.Name()))
		out = append(out, entry{
			Name: e.Name(), ModTime: info.ModTime().Unix(),
			Title: si.Title, Detail: si.Detail,
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ModTime > out[j].ModTime })

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

func (h *Handler) open(w http.ResponseWriter, r *http.Request) {
	if !h.authed(r) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if err := openDir(h.dir()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func openDir(dir string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", dir)
	case "darwin":
		cmd = exec.Command("open", dir)
	default:
		cmd = exec.Command("xdg-open", dir)
	}
	return cmd.Start()
}

func (h *Handler) openURL(w http.ResponseWriter, r *http.Request) {
	if !h.authed(r) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	u, err := url.Parse(r.URL.Query().Get("url"))
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		http.Error(w, "bad url", http.StatusBadRequest)
		return
	}
	if err := openInBrowser(u.String()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func openInBrowser(rawURL string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", rawURL)
	case "darwin":
		cmd = exec.Command("open", rawURL)
	default:
		cmd = exec.Command("xdg-open", rawURL)
	}
	return cmd.Start()
}

func (h *Handler) file(w http.ResponseWriter, r *http.Request) {
	if !h.authed(r) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	path := h.resolve(r.URL.Query().Get("name"))
	if path == "" {
		http.Error(w, "bad name", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		data, err := os.ReadFile(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(data)
	case http.MethodPost:
		data, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxSaveBytes))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := os.WriteFile(path, data, 0600); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
