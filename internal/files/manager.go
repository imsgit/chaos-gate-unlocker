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
	"runtime"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"github.com/goccy/go-json"
)

var (
	sep = []byte("\r\n")

	ErrWrongSaveFileFormat = errors.New("\n\nWrong save file format.")
	ErrSaveFile            = errors.New("\n\nCannot save file.")
)

const (
	minVersion = 1170

	onShipLocation = "COMMON_Baleful_Edict"

	saveDir   = "AppData/LocalLow/Complex Games Inc_/GreyKnights/SaveGames/Campaign"
	protonDir = "/1611910/pfx"
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

func (m *Manager) OnShip() bool {
	return m.header != nil && m.header.Location == onShipLocation
}

func (m *Manager) GetCurrentPath() string {
	pathBinding := binding.BindPreferenceString("path", fyne.CurrentApp().Preferences())
	currentPath, _ := pathBinding.Get()
	dir := filepath.Dir(currentPath)
	if !dirNotExists(dir) {
		return dir
	}

	dir, _ = os.UserHomeDir()

	switch runtime.GOOS {
	case "linux":
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

	if dirNotExists(dir) {
		dir, _ = os.Getwd()
	}

	return dir
}

func searchDir(root, searchPath string) string {
	var result string

	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || !info.IsDir() || path == root || !strings.HasSuffix(path, searchPath) {
			return nil
		}

		result = path
		return fs.SkipAll
	})

	return result
}

func dirNotExists(path string) bool {
	_, err := os.Stat(path)
	return err != nil && os.IsNotExist(err)
}

func (m *Manager) Load(reader fyne.URIReadCloser) error {
	m.filePath = reader.URI().Path()

	file, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	_ = reader.Close()

	chunks := bytes.Split(file, sep)
	if len(chunks) < 3 {
		return ErrWrongSaveFileFormat
	}

	headerBytes, stateBytes, combatStateBytes := chunks[0], chunks[1], chunks[2]

	if err = m.loadHeader(headerBytes, &m.header); err != nil {
		return err
	}

	if err = m.loadState(stateBytes, &m.state); err != nil {
		return err
	}

	m.combatStateBytes = combatStateBytes

	for i := range m.onLoadState {
		m.onLoadState[i](m.state)
	}

	pathBinding := binding.BindPreferenceString("path", fyne.CurrentApp().Preferences())
	_ = pathBinding.Set(m.filePath)

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

func (m *Manager) loadHeader(headerBytes []byte, header **internal.Header) error {
	var newHeader internal.Header
	if err := json.Unmarshal(headerBytes, &newHeader); err != nil {
		return ErrWrongSaveFileFormat
	}
	*header = &newHeader

	ver, _ := strconv.Atoi(newHeader.Version)
	if ver < minVersion {
		return ErrWrongSaveFileFormat
	}

	return nil
}

func (m *Manager) loadState(stateBytes []byte, state **internal.State) error {
	if len(stateBytes) == 0 || stateBytes[0] != 194 {
		return ErrWrongSaveFileFormat
	}

	decodedState := decode(stateBytes)

	var newState internal.State
	if err := json.Unmarshal(decodedState, &newState); err != nil {
		return ErrWrongSaveFileFormat
	}
	*state = &newState

	return nil
}

func (m *Manager) Save() error {
	headerBytes, err := json.Marshal(m.header)
	if err != nil {
		return ErrSaveFile
	}

	stateBytes, err := json.Marshal(m.state)
	if err != nil {
		return ErrSaveFile
	}
	stateBytes = encode(stateBytes)

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

	err = os.WriteFile(m.filePath, file, 0644)
	if err != nil {
		return ErrSaveFile
	}

	return nil
}
