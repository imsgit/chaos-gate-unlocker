package bridge

import (
	"crypto/subtle"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
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
		out = append(out, entry{Name: e.Name(), ModTime: info.ModTime().Unix()})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ModTime > out[j].ModTime })

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
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
