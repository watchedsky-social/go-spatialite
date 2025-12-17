package sharedlibs

import (
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/jghiloni/go-commonutils/v3/values"
)

const pathOverrideEnvVar = "GO_SPATIALITE_SEARCH_PATH"

var (
	pathListSeparator = string([]rune{os.PathListSeparator})
	homeDir           = values.Must(os.UserHomeDir())
)

// GetSharedLibrarySearchPath will return a slice of strings representing a set
// of directories will be searched for shared library files. If the env variable
// GO_SPATIALITE_SEARCH_PATH is set, it will be used. Panics if the current OS is
// darwin and the current working directory cannot be determined
func GetSharedLibrarySearchPath() []string {
	envPath, set := os.LookupEnv(pathOverrideEnvVar)
	if set {
		return strings.Split(envPath, pathListSeparator)
	}

	prependHomebrew := false
	var paths []string
	switch runtime.GOOS {
	case "darwin":
		prependHomebrew = true
		exePath, err := os.Executable()
		if err == nil {
			dir, _ := filepath.Split(exePath)
			if strings.HasSuffix(dir, "/Contents/bin") {
				paths = append(paths, filepath.Dir(dir))
			}
		}

		_, err = os.Stat("/Applications/Xcode.app/Contents")
		if err == nil {
			paths = append(paths, "/Applications/Xcode.app/Contents")
		}
	case "linux":
		var err error
		prependHomebrew = true
		if paths, err = ReadConfFile("/etc/ld.so.conf"); err != nil {
			return nil
		}
	case "windows":
		paths = strings.Split(os.Getenv("PATH"), pathListSeparator)
	}

	if prependHomebrew {
		hbDir, isSet := os.LookupEnv("HOMEBREW_PREFIX")
		if isSet {
			paths = appendNoDuplicates([]string{filepath.Join(hbDir, "lib")}, paths...)
		}
	}

	return paths
}

// FindSharedLibraryFiles will search the library search path returned by
// [GetSharedLibrarySearchPath] and their subdirectories for the requested
// shared library files. There are a few OS-specific behaviors that should be
// noted:
//
//   - For linux, it is assumed that the file's extension is .so*
//   - For darwin, it is assumed that the file's extension is .dylib
//   - For windows, it is assumed that the file's extension is .dll
//
// As the linux example shows, wildcards are allowed as long as they conform
// to the globbing patterns accepted by [path/filepath.Glob]. A slice of files
// is returned; an empty slice is not an error condition, and an error will only
// be returned if an underlying error occurs. Because files are globbed, it is
// possible for findFirstOnly to be set to true and still return multiple files,
// if the first glob returns multiple candidates
func FindSharedLibraryFiles(filePattern string, findFirstOnly bool) ([]string, error) {
	switch runtime.GOOS {
	case "linux":
		filePattern += ".so.*"
	case "darwin":
		filePattern += ".dylib"
	case "windows":
		filePattern += ".dll"
	}

	var matches []string
	for _, root := range GetSharedLibrarySearchPath() {
		files := findFileRecursively(root, filePattern)
		if len(files) > 0 && findFirstOnly {
			return files[:1], nil
		}

		matches = append(matches, files...)
	}

	return matches, nil
}

func findFileRecursively(absRoot string, filePattern string) []string {
	root := os.DirFS(absRoot)

	var matches []string
	fs.WalkDir(root, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if matched, e := filepath.Match(filePattern, filepath.Base(path)); e == nil && matched {
			matches = append(matches, filepath.Join(absRoot, path))
		}

		return nil
	})

	return matches
}
