package jux

import (
	"errors"
	"io"
	"strings"
)

type Args struct {
	Left      string
	Right     string
	Limit     int
	ErrWriter io.Writer
}

func (a Args) Validate() error {
	if a.Left == "" || a.Right == "" {
		return errors.New("two few files/directories passed")
	}
	if strings.HasPrefix(a.Left, "/") {
		return errors.New("path may not end in a /")
	}
	if strings.HasPrefix(a.Right, "/") {
		return errors.New("path may not end in a /")
	}

	return nil
}
