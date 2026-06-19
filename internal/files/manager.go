package files

import (
	"chaos-gate-unlocker/internal"

	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"github.com/goccy/go-json"
)

var (
	sep = []byte("\r\n")

	ErrWrongSaveFileFormat = errors.New("\n\n\nError. Wrong save file format.\n\n")
	ErrSaveFile            = errors.New("\n\n\nError. Cannot save file.\n\n")
)

const (
	minVersion = 1170

	appID      = "1611910"
	saveDir    = "AppData/LocalLow/Complex Games Inc_/GreyKnights/SaveGames/Campaign"
	protonDir  = "/1611910/pfx"
	protonUser = "pfx/drive_c/users/steamuser"
)

type Manager struct {
	filePath         string
	header           *internal.Header
	state            *internal.State
	combatStateBytes []byte
	onLoadState      []func(*internal.State)
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) OnLoadState(fn func(*internal.State)) {
	m.onLoadState = append(m.onLoadState, fn)
}

func (m *Manager) SaveDir() string {
	return saveDir
}

func (m *Manager) DefaultLocationHint() string {
	home, _ := os.UserHomeDir()

	switch runtime.GOOS {
	case "windows":
		if home == "" {
			return `%USERPROFILE%\` + filepath.FromSlash(saveDir)
		}
		return filepath.Join(home, filepath.FromSlash(saveDir))
	case "linux":
		if dir := steamSaveDir(home); dir != "" {
			return dir
		}
		base := filepath.Join(home, ".steam", "steam")
		return filepath.Join(base, "steamapps", "compatdata", appID, protonUser, saveDir)
	default:
		return saveDir
	}
}

func (m *Manager) GetCurrentPath() string {
	currentPath := fyne.CurrentApp().Preferences().String("path")
	dir := filepath.Dir(currentPath)
	if currentPath != "" && dirExists(dir) {
		return dir
	}

	dir, _ = os.UserHomeDir()

	switch runtime.GOOS {
	case "linux":
		if found := steamSaveDir(dir); found != "" {
			return found
		}

		dirSteam := searchDir(filepath.Join(dir, ".steam"), protonDir)
		if dirSteam != "" {
			dir = dirSteam
		} else {
			dir = searchDir(dir, protonDir)
		}

		if dir == "" {
			for _, path := range []string{"/run/media", "/media", "/mnt"} {
				dir = searchDir(path, protonDir)
				if dir != "" {
					break
				}
			}
		}

		if dir != "" {
			dir = searchDir(dir, saveDir)
		}
	case "windows":
		dir = filepath.Join(dir, saveDir)
	}

	if !dirExists(dir) {
		dir, _ = os.Getwd()
	}

	return dir
}

func steamSaveDir(home string) string {
	for _, lib := range steamLibraries(home) {
		dir := filepath.Join(lib, "steamapps", "compatdata", appID, protonUser, saveDir)
		if dirExists(dir) {
			return dir
		}
	}
	return ""
}

func steamLibraries(home string) []string {
	bases := []string{
		filepath.Join(home, ".steam", "steam"),
		filepath.Join(home, ".steam", "root"),
		filepath.Join(home, ".local", "share", "Steam"),
		filepath.Join(home, ".var", "app", "com.valvesoftware.Steam", ".local", "share", "Steam"),
	}

	seen := map[string]bool{}
	var libs []string
	add := func(path string) {
		if resolved, err := filepath.EvalSymlinks(path); err == nil {
			path = resolved
		}
		if path != "" && !seen[path] && dirExists(path) {
			seen[path] = true
			libs = append(libs, path)
		}
	}

	for _, base := range bases {
		add(base)
		for _, lib := range parseLibraryFolders(filepath.Join(base, "steamapps", "libraryfolders.vdf")) {
			add(lib)
		}
	}
	return libs
}

var libraryPathRe = regexp.MustCompile(`"path"\s+"([^"]+)"`)

func parseLibraryFolders(vdfPath string) []string {
	data, err := os.ReadFile(vdfPath)
	if err != nil {
		return nil
	}

	var paths []string
	for _, match := range libraryPathRe.FindAllStringSubmatch(string(data), -1) {
		paths = append(paths, strings.ReplaceAll(match[1], `\\`, `/`))
	}
	return paths
}

func searchDir(root, searchPath string) string {
	var result string

	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || !d.IsDir() || path == root || !strings.HasSuffix(path, searchPath) {
			return nil
		}

		result = path
		return filepath.SkipAll
	})

	return result
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func (m *Manager) Load(reader fyne.URIReadCloser) error {
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	return m.LoadBytes(reader.URI().Path(), data)
}

func (m *Manager) LoadBytes(path string, file []byte) error {
	m.filePath = path

	chunks := bytes.SplitN(file, sep, 3)
	if len(chunks) < 3 {
		return ErrWrongSaveFileFormat
	}

	headerBytes, stateBytes, combatStateBytes := chunks[0], chunks[1], chunks[2]

	if err := m.loadHeader(headerBytes); err != nil {
		return err
	}

	if err := m.loadState(stateBytes); err != nil {
		return err
	}

	m.combatStateBytes = bytes.Clone(combatStateBytes)

	for _, callback := range m.onLoadState {
		callback(m.state)
	}

	fyne.CurrentApp().Preferences().SetString("path", m.filePath)

	return nil
}

func (m *Manager) Name() string { return filepath.Base(m.filePath) }

func (m *Manager) loadHeader(headerBytes []byte) error {
	var header internal.Header
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return ErrWrongSaveFileFormat
	}

	if ver, _ := strconv.Atoi(header.Version); ver < minVersion {
		return ErrWrongSaveFileFormat
	}

	m.header = &header
	return nil
}

func (m *Manager) loadState(stateBytes []byte) error {
	if len(stateBytes) == 0 || stateBytes[0] != 194 {
		return ErrWrongSaveFileFormat
	}

	var state internal.State
	if err := json.Unmarshal(encodeDecode(stateBytes), &state); err != nil {
		return ErrWrongSaveFileFormat
	}

	m.state = &state
	return nil
}

func (m *Manager) Encode() ([]byte, error) {
	headerBytes, err := json.Marshal(m.header)
	if err != nil {
		return nil, ErrSaveFile
	}

	stateBytes, err := json.Marshal(m.state)
	if err != nil {
		return nil, ErrSaveFile
	}
	stateBytes = encodeDecode(stateBytes)

	fileLength := len(headerBytes) + len(sep) + len(stateBytes) + len(sep)

	if len(m.combatStateBytes) > 0 {
		fileLength += len(m.combatStateBytes) + len(sep)
	}

	file := make([]byte, 0, fileLength)
	file = append(file, headerBytes...)
	file = append(file, sep...)
	file = append(file, stateBytes...)
	file = append(file, sep...)

	if len(m.combatStateBytes) > 0 {
		file = append(file, m.combatStateBytes...)
		file = append(file, sep...)
	}

	return file, nil
}

func (m *Manager) Save() error {
	file, err := m.Encode()
	if err != nil {
		return err
	}

	if err := os.WriteFile(m.filePath, file, 0600); err != nil {
		return ErrSaveFile
	}

	return nil
}

func (m *Manager) Status() string {
	var difficulty string
	switch m.header.Difficulty {
	case 3:
		difficulty = "Legendary"
	case 2:
		difficulty = "Ruthless"
	case 1:
		difficulty = "Standard"
	case 0:
		difficulty = "Merciful"
	}
	if m.header.IronMan {
		difficulty += " Ironman"
	}

	timestamp := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		m.header.SavedTimeStamp.Years, m.header.SavedTimeStamp.Months, m.header.SavedTimeStamp.Days,
		m.header.SavedTimeStamp.Hours, m.header.SavedTimeStamp.Minutes, m.header.SavedTimeStamp.Seconds)

	return fmt.Sprintf("%s | %s | Days: %d | Difficulty: %s | %s",
		filepath.Base(m.filePath), m.header.SaveName, m.header.GameDays, difficulty, timestamp)
}
