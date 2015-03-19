package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/shavac/readline"
)

var histFile = "/tmp/$USER/rdlin_hist"

func main() {
	prompt := "foo >"
	histFn := os.ExpandEnv(histFile)
	err := readline.ReadHistoryFile(histFn)
	if err != nil {
		if err != syscall.ENOENT {
			fmt.Fprintf(os.Stderr, "could not read %s: %s\n", histFn, err)
			os.Exit(2)
		}

	}

	for {
		input := readline.ReadLine(&prompt)
		if *input != "" {
			readline.AddHistory(*input)
			fmt.Println(*input)
			readline.WriteHistoryFile(histFn)
		}
	}
}


