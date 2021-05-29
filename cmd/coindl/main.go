package main

import (
	"fmt"
	"os"

	"github.com/markdl-dev/coin-dl/internal/cli/coins/info"
	"github.com/markdl-dev/coin-dl/internal/cli/coins/market"
	"github.com/markdl-dev/coin-dl/internal/cli/exchangerate"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) > 0 {
		var err error
		switch args[0] {
		case exchangerate.CmdName:
			err = exchangerate.Cmd()
		case market.CmdName:
			err = market.Cmd()
		case info.CmdName:
			err = info.Cmd()
		default:
			fmt.Println("default")
		}
		return err
	}
	// TODO handle no command
	return nil
}
