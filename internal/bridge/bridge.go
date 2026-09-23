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
	"strconv"
	"strings"

	"chaos-gate-unlocker/internal/save"
)

const maxSaveBytes = 16 << 20

type Handler struct {
	token string
	dir   func() string
}

func New(token string, dir func() string) *Handler {
	return &Handler{token: token, dir: dir}
}

func (h *Handler) Register(mux *http.ServeMux) {
	for path, f := range map[string]http.HandlerFunc{
		"/api/list":    h.list,
		"/api/file":    h.file,
		"/api/open":    h.open,
		"/api/openurl": h.openURL,
	} {
		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			if !h.authed(r) {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			f(w, r)
		})
	}
}

func (h *Handler) authed(r *http.Request) bool {
	got := r.URL.Query().Get("t")
	return h.token != "" && subtle.ConstantTimeCompare([]byte(got), []byte(h.token)) == 1
}

func (h *Handler) resolve(name string) string {
	if name == "" || filepath.Base(name) != name || strings.ContainsRune(name, ':') || !strings.HasSuffix(name, ".gksave") {
		return ""
	}
	return filepath.Join(h.dir(), name)
}

type entry struct {
	Name string `json:"name"`
	save.Info
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	dir := h.dir()
	ents, err := os.ReadDir(dir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out := make([]entry, 0, len(ents))
	for _, e := range ents {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".gksave") {
			continue
		}
		out = append(out, entry{e.Name(), save.ParseFile(filepath.Join(dir, e.Name()))})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

func (h *Handler) open(w http.ResponseWriter, r *http.Request) {
	if err := OpenDir(h.dir()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func openExternally(winName string, winArgs []string, target string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command(winName, append(winArgs, target)...)
	case "darwin":
		cmd = exec.Command("open", target)
	default:
		cmd = exec.Command("xdg-open", target)
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	go cmd.Wait()
	return nil
}

func OpenDir(dir string) error {
	return openExternally("explorer", nil, dir)
}

func (h *Handler) openURL(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.Query().Get("url"))
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" || u.User != nil {
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
	return openExternally("rundll32", []string{"url.dll,FileProtocolHandler"}, rawURL)
}

func WriteFileAtomic(path string, parts ...[]byte) error {
	mode := os.FileMode(0o644)
	if info, err := os.Stat(path); err == nil {
		mode = info.Mode().Perm()
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+".*.tmp")
	if err != nil {
		return err
	}
	name := tmp.Name()
	for _, part := range parts {
		if _, err = tmp.Write(part); err != nil {
			break
		}
	}
	if err == nil {
		err = tmp.Chmod(mode)
	}
	if err == nil {
		err = tmp.Sync()
	}
	if cerr := tmp.Close(); err == nil {
		err = cerr
	}
	if err == nil {
		err = os.Rename(name, path)
	}
	if err != nil {
		os.Remove(name)
		return err
	}
	return nil
}

func (h *Handler) file(w http.ResponseWriter, r *http.Request) {
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
		w.Header().Set("Content-Length", strconv.Itoa(len(data)))
		w.Write(data)
	case http.MethodPost:
		data, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxSaveBytes))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := WriteFileAtomic(path, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
