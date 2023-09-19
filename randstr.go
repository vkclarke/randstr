package main

import (
	"bytes"
	"crypto/rand"
	_ "embed"
	"fmt"
	"os"
	"unicode"
)

//go:embed usage.txt
var usage string

var random = make(chan byte)

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
		out := bytes.NewBuffer(make([]byte, 0, os.Getpagesize()))
		for i := range args {
			switch a := args[i].(type) {
			case string:
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
					// handle unrecognized
					fmt.Fprintf(os.Stderr, "randstr: invalid argument: %q\n", a)
					os.Stderr.WriteString(usage)
					return 1
				}
			case int:
				randstr(out, a, numbers, upper, lower, hex)
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
