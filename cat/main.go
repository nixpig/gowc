package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

type config struct {
	n bool
	b bool
}

func main() {
	cfg := config{}
	pflag.BoolVarP(&cfg.n, "number", "n", false, "number all output lines")
	pflag.BoolVarP(&cfg.b, "number-nonblank", "b", false, "number non-empty output lines, overrides -n")

	pflag.Parse()
	files := pflag.Args()

	concat := strings.Builder{}

	if len(files) == 0 || files[0] == "-" {
		contents, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println(fmt.Sprintf("unable to read from stdin: %s", err))
			os.Exit(1)
		}

		concat.WriteString(string(contents))
	} else {
		for _, file := range files {
			contents, err := os.ReadFile(file)
			if err != nil {
				fmt.Println(fmt.Sprintf("unable to read from file '%s': %s", file, err))
				os.Exit(1)
			}

			concat.WriteString(string(contents))
		}
	}

	scanner := bufio.NewScanner(strings.NewReader(concat.String()))

	i := 1
	for scanner.Scan() {
		t := scanner.Text()

		if cfg.b {
			fmt.Printf("%d  ", i)
			i++
		} else if cfg.n {
			if t != "" {
				fmt.Printf("%d  ", i)
				i++
			}
		}

		fmt.Printf("%s\n", t)
	}
}
