package save

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

const (
	appID      = "1611910"
	dir        = "AppData/LocalLow/Complex Games Inc_/GreyKnights/SaveGames/Campaign"
	protonDir  = "/1611910/pfx"
	protonUser = "pfx/drive_c/users/steamuser"
)

func Dir() string { return dir }

func Discover(currentPath string) string {
	d := filepath.Dir(currentPath)
	if currentPath != "" && dirExists(d) {
		return d
	}

	d, _ = os.UserHomeDir()

	switch runtime.GOOS {
	case "linux":
		if found := steamSaveDir(d); found != "" {
			return found
		}

		dirSteam := searchDir(filepath.Join(d, ".steam"), protonDir)
		if dirSteam != "" {
			d = dirSteam
		} else {
			d = searchDir(d, protonDir)
		}

		if d == "" {
			for _, path := range []string{"/run/media", "/media", "/mnt"} {
				d = searchDir(path, protonDir)
				if d != "" {
					break
				}
			}
		}

		if d != "" {
			d = searchDir(d, dir)
		}
	case "windows":
		d = filepath.Join(d, dir)
	}

	if !dirExists(d) {
		d, _ = os.Getwd()
	}

	return d
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
