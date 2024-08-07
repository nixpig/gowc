package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

type config struct {
	c bool
	l bool
	w bool
	m bool
	t bool
}

func main() {

	c := flag.Bool("c", false, "print the byte counts")
	l := flag.Bool("l", false, "print the newline counts")
	w := flag.Bool("w", false, "print the word counts")
	m := flag.Bool("m", false, "print the character counts")
	t := flag.Bool("total", false, "print the total")

	flag.Parse()

	cfg := config{
		c: *c,
		l: *l,
		w: *w,
		m: *m,
		t: *t,
	}

	if len(flag.Args()) == 0 {
		contents, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Printf("error reading from stdin: %s", err)
			os.Exit(1)
		}

		fmt.Println(count(string(contents), cfg))
		return
	}

	for _, file := range flag.Args() {
		contents, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("error reading file '%s': %s\n", file, err)
			os.Exit(1)
		}

		fmt.Println(count(string(contents), cfg))
	}
}

func count(contents string, cfg config) string {
	counted := strings.Builder{}

	if !cfg.c && !cfg.l && !cfg.w && !cfg.m {
		cfg.c = true
		cfg.l = true
		cfg.w = true
	}

	if cfg.l {
		counted.WriteString(fmt.Sprintf("%d ", len(strings.Split(string(contents), "\n"))))
	}

	if cfg.w {
		counted.WriteString(fmt.Sprintf("%d ", len(strings.Fields(string(contents)))))
	}

	if cfg.c {
		counted.WriteString(fmt.Sprintf("%d ", len(contents)))
	}

	if cfg.m {
		counted.WriteString(fmt.Sprintf("%d ", utf8.RuneCount([]byte(contents))))
	}

	counted.WriteString("\n")

	return counted.String()
}
