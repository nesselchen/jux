package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/nesselchen/jux/jux"
)

func parse() (*jux.Args, error) {
	limit := flag.Int("limit", 0, "Specify how many differences are tracked. Abort comparison after limit is surpassed.")
	flag.Parse()
	var left, right string
	if rest := flag.Args(); len(rest) == 2 {
		// ignore trailing slash from terminal autocompleter
		left, _ = strings.CutSuffix(rest[0], "/")
		right, _ = strings.CutSuffix(rest[1], "/")
	}
	args := &jux.Args{
		Left:  left,
		Right: right,
		Limit: *limit,
	}
	if err := args.Validate(); err != nil {
		return nil, err
	}
	return args, nil
}

func main() {
	if err := run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func run() error {
	args, err := parse()
	if err != nil {
		flag.PrintDefaults()
		return err
	}

	leftInfo, err := os.Stat(args.Left)
	if err != nil {
		return err
	}
	rightInfo, err := os.Stat(args.Right)
	if err != nil {
		return err
	}
	if os.SameFile(leftInfo, rightInfo) {
		fmt.Println("Warning: left and right tree point to the same file")
		return nil
	}

	leftFiles, err := fileMap(args.Left)
	if err != nil {
		return err
	}
	rightFiles, err := fileMap(args.Right)
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
		fmt.Printf("Only in %s: %s\n", args.Left, key)
	}
	for key := range rightFiles {
		fmt.Printf("Only in %s: %s\n", args.Right, key)
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
