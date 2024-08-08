package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/spf13/pflag"
)

type config struct {
	file  string
	c     bool
	l     bool
	w     bool
	m     bool
	total string
	L     bool
}

type count struct {
	file string
	c    int
	l    int
	w    int
	m    int
	L    int
}

func main() {
	c := pflag.BoolP("bytes", "c", false, "print the byte counts")
	l := pflag.BoolP("lines", "l", false, "print the newline counts")
	w := pflag.BoolP("words", "w", false, "print the word counts")
	m := pflag.BoolP("chars", "m", false, "print the character counts")
	t := pflag.String("total", "auto", "when to print a line with total counts; WHEN can be: auto, always, only, never")
	L := pflag.BoolP("max-line-length", "L", false, "print the maximum display width")

	pflag.Parse()

	cfg := config{
		c:     *c,
		l:     *l,
		w:     *w,
		m:     *m,
		total: *t,
		L:     *L,
	}

	var counted []count

	files := pflag.Args()

	if len(files) == 0 || files[0] == "-" {
		contents, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Printf("error reading from stdin: %s", err)
			os.Exit(1)
		}

		counted = append(counted, wc(string(contents), cfg))
	} else {
		counted = make([]count, len(files))

		for i, file := range files {
			contents, err := os.ReadFile(file)
			if err != nil {
				fmt.Printf("error reading file '%s': %s\n", file, err)
				os.Exit(1)
			}

			cfg.file = file

			counted[i] = wc(string(contents), cfg)
		}
	}

	switch cfg.total {
	case "always":
		counted = append(counted, total(counted))

	case "only":
		counted = []count{total(counted)}

	case "never":
		break

	case "auto":
		fallthrough
	default:
		if len(counted) > 1 {
			counted = append(counted, total(counted))
		}

	}

	for _, line := range counted {
		b := strings.Builder{}

		if cfg.c {
			b.WriteString(fmt.Sprintf("%d ", line.c))
		}

		if cfg.l {
			b.WriteString(fmt.Sprintf("%d ", line.l))
		}

		if cfg.w {
			b.WriteString(fmt.Sprintf("%d ", line.w))
		}

		if cfg.m {
			b.WriteString(fmt.Sprintf("%d ", line.m))
		}

		if cfg.L {
			b.WriteString(fmt.Sprintf("%d ", line.L))
		}

		b.WriteString(fmt.Sprintf("%s", line.file))

		fmt.Println(b.String())
	}
}

func wc(contents string, cfg config) count {
	counted := count{file: cfg.file}

	if !cfg.c && !cfg.l && !cfg.w && !cfg.m && !cfg.L {
		cfg.c = true
		cfg.l = true
		cfg.w = true
	}

	if cfg.l {
		counted.l = len(strings.Split(string(contents), "\n"))
	}

	if cfg.w {
		counted.w = len(strings.Fields(string(contents)))
	}

	if cfg.c {
		counted.c = len(contents)
	}

	if cfg.m {
		counted.m = utf8.RuneCount([]byte(contents))
	}

	if cfg.L {
		for _, line := range strings.Split(contents, "\n") {
			length := len(line)
			if length > counted.L {
				counted.L = length
			}

		}
	}

	return counted
}

func total(counts []count) count {
	totalled := count{file: "total"}
	for _, count := range counts {
		totalled.c += count.c
		totalled.l += count.l
		totalled.w += count.w
		totalled.m += count.m
		totalled.L += count.L
	}

	return totalled
}
