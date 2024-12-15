package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/nesselchen/jux/jux"
)

func main() {
	err := run()
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, err.Error())
	if err, ok := err.(exitError); ok {
		os.Exit(err.code)
	}
	os.Exit(1)
}

func run() error {
	args, err := parse()
	if err != nil {
		flag.Usage()
		fmt.Println()
		return err
	}

	diff, err := jux.CompareTrees(args)
	if err != nil {
		return err
	}

	if diff != 0 {
		return exitError{
			msg:  fmt.Sprintf("%d differences", diff),
			code: 1,
		}
	}
	return nil
}

type exitError struct {
	code int
	msg  string
}

func (err exitError) Error() string {
	return err.msg
}

func parse() (jux.Args, error) {
	limit := flag.Int("limit", 0, "Specify how many differences are tracked. Abort comparison after limit is surpassed.")
	flag.Parse()
	var left, right string
	if rest := flag.Args(); len(rest) == 2 {
		// ignore trailing slash from terminal autocompleter
		left, _ = strings.CutSuffix(rest[0], "/")
		right, _ = strings.CutSuffix(rest[1], "/")
	}
	args := jux.Args{
		Left:      left,
		Right:     right,
		Limit:     *limit,
		ErrWriter: os.Stdout,
	}
	if err := args.Validate(); err != nil {
		return jux.Args{}, err
	}
	return args, nil
}
