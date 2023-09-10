package main

import (
	"flag"
	"fmt"
	"hive/packages/check"
	"hive/packages/tidy"
	"log"
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
		if err := tidy.Tidy(); err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Println("unexpected subcommand")
		os.Exit(1)
	}
}
