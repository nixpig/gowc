package main

import (
	"strings"
	"unicode/utf8"
)

type count struct {
	file string
	c    int
	l    int
	w    int
	m    int
	L    int
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
