package jux

import "errors"

type Args struct {
	Left  string
	Right string
	Limit int
}

func (a Args) Validate() error {
	if a.Left == "" || a.Right == "" {
		return errors.New("two few files/directories passed")
	}

	return nil
}
