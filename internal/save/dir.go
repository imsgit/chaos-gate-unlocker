package save

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

const (
	appID      = "1611910"
	dir        = "AppData/LocalLow/Complex Games Inc_/GreyKnights/SaveGames/Campaign"
	protonDir  = "/1611910/pfx"
	protonUser = "pfx/drive_c/users/steamuser"
)

var discovered = sync.OnceValue(func() string {
	if d := discover(); dirExists(d) {
		return d
	}
	return ""
})

func Discover(currentPath string) string {
	if d := filepath.Dir(currentPath); currentPath != "" && dirExists(d) {
		return d
	}
	if d := discovered(); dirExists(d) {
		return d
	}
	d, _ := os.Getwd()
	return d
}

func discover() string {
	home, _ := os.UserHomeDir()

	switch runtime.GOOS {
	case "linux":
		if found := steamSaveDir(home); found != "" {
			return found
		}
		for _, root := range []string{filepath.Join(home, ".steam"), home, "/run/media", "/media", "/mnt"} {
			if d := searchDir(root, protonDir); d != "" {
				return searchDir(d, dir)
			}
		}
		return ""
	case "windows":
		return filepath.Join(home, dir)
	}
	return home
}

func steamSaveDir(home string) string {
	for _, lib := range steamLibraries(home) {
		d := filepath.Join(lib, "steamapps", "compatdata", appID, protonUser, dir)
		if dirExists(d) {
			return d
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
		filepath.Join(home, "snap", "steam", "common", ".local", "share", "Steam"),
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

	filepath.WalkDir(root, func(path string, e fs.DirEntry, err error) error {
		if err != nil || !e.IsDir() || path == root || !strings.HasSuffix(path, searchPath) {
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
