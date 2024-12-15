package jux

import (
	"bytes"
	"fmt"
	"os"
)

func CompareTrees(args Args) (int, error) {
	same, err := sameFile(args.Left, args.Right)
	if err != nil {
		return 0, err
	} else if same {
		fmt.Fprintln(args.ErrWriter, "left and right file root point to the same file")
		return 0, nil
	}

	lf, err := accumulate(args.Left)
	if err != nil {
		return 0, err
	}
	rf, err := accumulate(args.Right)
	if err != nil {
		return 0, err
	}

	var c Comparison = comparisonFunc(compareBytes)

	var differences int
	for path := range lf.all() {
		ok := rf.Items[path]
		if !ok {
			continue
		}
		leftPath, rightPath := lf.Prefix+"/"+path, rf.Prefix+"/"+path
		report, matched, err := c.Compare(leftPath, rightPath)
		if err != nil {
			return 0, err
		}
		if !matched {
			differences++
			fmt.Fprintf(args.ErrWriter, "%s: %s\n", report, path)
		}
		if args.Limit > 0 && differences >= args.Limit {
			break
		}

		delete(lf.Items, path)
		delete(rf.Items, path)
	}

	for path := range lf.all() {
		fmt.Fprintf(args.ErrWriter, "only in %s: %s\n", lf.Prefix, path)
	}
	for path := range rf.all() {
		fmt.Fprintf(args.ErrWriter, "only in %s: %s\n", rf.Prefix, path)
	}

	differences += len(lf.Items) + len(rf.Items)
	if differences == 0 {
		return 0, nil
	}
	return differences, nil
}

type Comparison interface {
	Compare(left, right string) (string, bool, error)
}

type comparisonFunc func(left, right string) (string, bool, error)

func (f comparisonFunc) Compare(left, right string) (string, bool, error) {
	return f(left, right)
}

func compareBytes(left, right string) (string, bool, error) {
	lb, err := os.ReadFile(left)
	if err != nil {
		return "", false, err
	}
	rb, err := os.ReadFile(right)
	if err != nil {
		return "", false, err
	}
	if bytes.Equal(lb, rb) {
		return "", true, nil
	}
	return "mismatched bytes", false, nil
}
