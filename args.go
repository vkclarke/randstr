package main

import (
	"os"
	"strconv"
)

var args []any

func init() {
	args = make([]any, 0, len(os.Args[1:]))
	for _, s := range os.Args[1:] {
		if s[0] == '-' && len(s) > 1 {
			if s[1] != '-' {
				for _, c := range s[1:] {
					args = append(args, string([]rune{'-', c}))
				}
				continue
			}
		}
		if n, err := strconv.Atoi(s); err == nil {
			args = append(args, n)
			continue
		}
		args = append(args, s)
	}
}
