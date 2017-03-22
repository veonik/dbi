package main

import (
	"github.com/peterh/liner"
	"gopkg.in/mattes/go-expand-tilde.v1"
	"flag"
	"fmt"
	"os"
	"strings"
)

var historyPathFlag = flag.String("history-file", "~/.dbi-history", "Specify history file.")

var commands = []string{"connect", "exit"}

func main() {
	flag.Parse()

	histFile, err := tilde.Expand(*historyPathFlag)
	if err != nil {
		fmt.Println("error: ", err.Error())
		return
	}
	cli := liner.NewLiner()
	if f, err := os.Open(histFile); err == nil {
		cli.ReadHistory(f)
		f.Close()
	}
	defer func() {
		f, err := os.Create(histFile)
		if err != nil {
			fmt.Println("error saving history:", err.Error())
			return
		}
		cli.WriteHistory(f)
		f.Close()
	}()

	cli.SetCompleter(func(line string) []string {
		res := []string{}
		for _, v := range commands {
			if strings.HasPrefix(v, line) {
				res = append(res, v)
			}
		}
		return res
	})

	for {
		cmd, err := cli.Prompt("> ")
		if err != nil {
			fmt.Println("error reading input:", err.Error())
			return
		}
		args := strings.Split(cmd, " ")
		if len(args) == 0 {
			continue
		}
		switch {
		case args[0] == "exit" || args[0] == "quit":
			return
		case args[0] == "connect":
			fmt.Println(cmd)
		case args[0] == "repl":
			// Run a query, print a result
		case args[0] == "paste":
			// Paste a query, print result
		default:
			continue
		}
		cli.AppendHistory(cmd)
	}
}