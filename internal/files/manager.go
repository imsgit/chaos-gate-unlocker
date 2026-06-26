package files

import (
	"chaos-gate-unlocker/internal"
	"chaos-gate-unlocker/internal/savedir"

	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"fyne.io/fyne/v2"
	"github.com/goccy/go-json"
)

var (
	sep = []byte("\r\n")

	ErrWrongSaveFileFormat = errors.New("\n\n\nError. Wrong save file format.\n\n")
	ErrSaveFile            = errors.New("\n\n\nError. Cannot save file.\n\n")
)

const minVersion = 1170

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

func (m *Manager) SaveDir() string { return savedir.SaveDir() }

func (m *Manager) DefaultLocationHint() string { return savedir.DefaultLocationHint() }

func (m *Manager) GetCurrentPath() string {
	return savedir.Discover(fyne.CurrentApp().Preferences().String("path"))
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
