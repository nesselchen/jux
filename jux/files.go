package jux

import (
	"io/fs"
	"iter"
	"os"
	"strings"
)

type items struct {
	Prefix string
	Items  map[string]bool
}

func (m items) all() iter.Seq[string] {
	return func(yield func(string) bool) {
		for k := range m.Items {
			if !yield(k) {
				return
			}
		}
	}
}

func sameFile(left, right string) (bool, error) {
	leftInfo, err := os.Stat(left)
	if err != nil {
		return false, err
	}
	rightInfo, err := os.Stat(right)
	if err != nil {
		return false, err
	}
	return os.SameFile(leftInfo, rightInfo), nil
}

func accumulate(dir string) (items, error) {
	im := items{
		Prefix: dir,
		Items:  make(map[string]bool),
	}
	err := fs.WalkDir(os.DirFS("."), dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		subpath, _ := strings.CutPrefix(path, dir+"/")
		im.Items[subpath] = true
		return nil
	})
	if err != nil {
		return items{}, err
	}

	return im, nil
}
