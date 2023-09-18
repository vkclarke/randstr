package main

import (
	"bytes"
	"crypto/rand"
	_ "embed"
	"fmt"
	"os"
	"strconv"
	"unicode"
)

//go:embed usage.txt
var usage string

var args []string
var random = make(chan byte)

func init() {
	args = make([]string, 0, len(os.Args[1:]))
	for _, a := range os.Args[1:] {
		if a[0] == '-' && len(a) > 1 {
			if a[1] != '-' {
				for _, c := range a[1:] {
					args = append(args, string([]rune{'-', c}))
				}
				continue
			}
		}
		args = append(args, a)
	}
}

func main() {
	go func() {
		buf := make([]byte, os.Getpagesize())
		for {
			rand.Read(buf)
			for _, b := range buf {
				random <- b
			}
		}
	}()
	os.Exit(func() (ret int) {
		var num int
		var upper, lower, numbers, hex bool
		var out = bytes.NewBuffer(make([]byte, 0, os.Getpagesize()))
		for _, a := range args {
			switch a {
			case "-h", "--help":
				os.Stderr.WriteString(usage)
				return 1
			case "-l", "--lower":
				lower = true
			case "-L", "--upper":
				upper = true
			case "-n", "--numbers":
				numbers = true
			case "-x", "--hex":
				hex = true
			default:
				n, err := strconv.Atoi(a)
				if err != nil {
					fmt.Fprintln(os.Stderr, "randstr:", err)
					return 1
				}
				randstr(out, n, numbers, upper, lower, hex)
				num++
			}
		}
		if num < 1 {
			randstr(out, 8, numbers, upper, lower, hex)
		}
		out.WriteTo(os.Stdout)
		return
	}())
}

func randstr(out *bytes.Buffer, n int, numbers, upper, lower, hex bool) {
	if !numbers && !upper && !lower {
		numbers, upper, lower = true, true, true
	}
	ranges := &unicode.RangeTable{
		R16: func() (r []unicode.Range16) {
			if numbers {
				r = append(r, unicode.Range16{0x0030, 0x0039, 1})
			}
			if upper {
				if hex {
					r = append(r, unicode.Range16{0x0041, 0x0046, 1})
				} else {
					r = append(r, unicode.Range16{0x0041, 0x005a, 1})
				}
			}
			if lower {
				if hex {
					r = append(r, unicode.Range16{0x0061, 0x0066, 1})
				} else {
					r = append(r, unicode.Range16{0x0061, 0x007a, 1})
				}
			}
			return r
		}(),
		LatinOffset: 3,
	}
	var prev byte
	for i := 0; i < n; i++ {
		for {
			b := <-random
			if unicode.In(rune(b), ranges) && b != prev {
				out.WriteByte(b)
				prev = b
				break
			}
		}
	}
	out.WriteByte(byte('\n'))
}
