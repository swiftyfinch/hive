package main

import (
	"flag"
	"fmt"
	"hive/check"
	"hive/tidy"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("missing subcommand")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "check":
		checkFlagSet := flag.NewFlagSet("check", flag.ExitOnError)
		checkFlagSet.Parse(os.Args[2:])
		check.Check()
	case "tidy":
		tidyFlagSet := flag.NewFlagSet("tidy", flag.ExitOnError)
		tidyFlagSet.Parse(os.Args[2:])
		tidy.Tidy()
	default:
		fmt.Println("unexpected subcommand")
		os.Exit(1)
	}
}
