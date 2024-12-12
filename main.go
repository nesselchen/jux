package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"
)

const usage = "usage: jux [file_root] [file_root]"

func main() {
	if err := run(); err != nil {
		fmt.Println("ERR:", err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) != 3 {
		return errors.New(usage)
	}

	lhs, rhs := os.Args[1], os.Args[2]

	// ignore trailing slash from terminal autocompleter
	lhs, _ = strings.CutSuffix(lhs, "/")
	rhs, _ = strings.CutSuffix(rhs, "/")

	leftInfo, err := os.Stat(lhs)
	if err != nil {
		return err
	}

	rightInfo, err := os.Stat(rhs)
	if err != nil {
		return err
	}

	if os.SameFile(leftInfo, rightInfo) {
		fmt.Println("WRN: left and right tree point to the same file")
		return nil
	}

	leftFiles, err := fileMap(lhs)
	if err != nil {
		return err
	}

	rightFiles, err := fileMap(rhs)
	if err != nil {
		return err
	}

	leftFiles, rightFiles, same := keyDiff(leftFiles, rightFiles)
	if same {
		return nil
	}

	// report differences
	for key := range leftFiles {
		if _, ok := rightFiles[key]; ok {
			fmt.Println("Mismatched bytes:", key)
			delete(rightFiles, key)
			continue
		}
		fmt.Printf("Only in %s: %s\n", lhs, key)
	}
	for key := range rightFiles {
		fmt.Printf("Only in %s: %s\n", rhs, key)
	}
	return nil
}

func fileMap(dir string) (map[string]string, error) {
	files := make(map[string]string)
	err := fs.WalkDir(os.DirFS("."), dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		key, _ := strings.CutPrefix(path, dir+"/")
		files[key] = path
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func keyDiff(m1, m2 map[string]string) (map[string]string, map[string]string, bool) {
	for key, v1 := range m1 {
		v2, ok := m2[key]
		if !ok {
			continue
		}
		// compare files
		content1, err := fs.ReadFile(os.DirFS("."), v1)
		if err != nil {
			return nil, nil, false
		}
		content2, err := fs.ReadFile(os.DirFS("."), v2)
		if err != nil {
			return nil, nil, false
		}
		if !bytes.Equal(content1, content2) {
			continue
		}
		delete(m1, key)
		delete(m2, key)
	}
	if len(m1) == 0 && len(m2) == 0 {
		return nil, nil, true
	}
	return m1, m2, false
}
