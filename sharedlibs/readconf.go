package sharedlibs

import (
	"bufio"
	"cmp"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"slices"
)

// ErrMalformedConfFile will be returned if the conf file can't be parsed
var ErrMalformedConfFile = errors.New("malformed conf file")

// ReadConfFile will parse a config file like /etc/ld.so.conf, in which each
// line of the file looks like one of the following:
//
//	whitespace-only: ignored
//	begins with optional whitespace and #: ignored
//	an absolute path to a directory: copied
//	has 'include <glob>': each file matching the glob will be processed recursively
//	with this function
//
// A list of unique paths will be returned, with only the first instance of a path
// included
func ReadConfFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var paths []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		switch {
		case line == "", strings.HasPrefix(line, "#"):
			continue
		case strings.HasPrefix(line, "/"):
			st, err := os.Stat(line)
			if err != nil && !errors.Is(err, fs.ErrNotExist) {
				return nil, fmt.Errorf("%w: %w", ErrMalformedConfFile, err)
			}

			if !st.IsDir() {
				return nil, fmt.Errorf("%w: %s is not a directory", ErrMalformedConfFile, line)
			}

			paths = appendNoDuplicates(paths, line)
		case strings.HasPrefix(strings.ToLower(line), "include "):
			if strings.EqualFold(line, "include ") {
				return nil, fmt.Errorf("%w: include must include a glob", ErrMalformedConfFile)
			}

			line = strings.TrimSpace(line[len("include "):])
			absPath, err := filepath.Abs(filePath)
			if err != nil {
				return nil, err
			}

			dir, _ := filepath.Split(absPath)
			if !filepath.IsAbs(line) {
				line = filepath.Join(dir, line)
			}
			
			matches, err := filepath.Glob(line)
			if err != nil {
				return nil, fmt.Errorf("%w: %w", ErrMalformedConfFile, err)
			}

			for _, m := range matches {
				subPaths, e := ReadConfFile(m)
				if e != nil {
					return nil, e
				}

				paths = appendNoDuplicates(paths, subPaths...)
			}
		}
	}

	return paths, nil
}

func appendNoDuplicates[T cmp.Ordered](s []T, a ...T) []T {
	y := make([]T, len(s), len(s)+len(a))
	copy(y, s)
	for _, x := range a {
		if !slices.Contains(y, x) {
			y = append(y, x)
		}
	}

	return y
}
