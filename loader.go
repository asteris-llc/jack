package jack

import (
	"errors"
	"github.com/kardianos/osext"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var (
	ErrNotFound = errors.New("executable not found in any search path")
)

type Loader struct {
	paths []string
}

func NewLoader(paths []string) (*Loader, error) {
	// we're going to contruct the paths that we'll search for plugins from a
	// couple of different sources. These are ordered highest-to-least priority.
	// We'll start with the paths given, as those should be the highest priority.

	// add the local executable directory
	local, err := osext.ExecutableFolder()
	if err != nil {
		return nil, err
	}
	paths = append(paths, local)

	// finish with paths in $PATH
	path := os.Getenv("PATH")
	paths = append(paths, strings.Split(path, ":")...)

	return &Loader{paths}, nil
}

func (c *Loader) searchInPaths(name string) (string, error) {
	for _, dir := range c.paths {
		files, err := ioutil.ReadDir(dir)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return "", err
		}

		for _, file := range files {
			if file.Name() != name {
				continue
			}

			mode := file.Mode()
			if !mode.IsRegular() {
				continue
			}

			if mode.Perm()&0111 != 0 {
				return path.Join(dir, file.Name()), nil
			}
		}
	}

	return "", ErrNotFound
}

func (c *Loader) Load(name string) (*Client, error) {
	fullPath, err := c.searchInPaths(name)
	if err != nil {
		return nil, err
	}

	return NewClient(fullPath), nil
}
