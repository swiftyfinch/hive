package main

import (
	"flag"
	"fmt"
	"log"
	"main/internal/commands/check"
	"main/internal/commands/tidy"
	"os"
)

const Config_Path = ".devtools/hive"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("missing subcommand")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "check":
		checkFlagSet := flag.NewFlagSet("check", flag.ExitOnError)
		checkFlagSet.Parse(os.Args[2:])
		if err := check.Check(Config_Path); err != nil {
			log.Fatal(err)
		}
	case "tidy":
		tidyFlagSet := flag.NewFlagSet("tidy", flag.ExitOnError)
		tidyFlagSet.Parse(os.Args[2:])
		if err := tidy.Tidy(Config_Path); err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Println("unexpected subcommand")
		os.Exit(1)
	}
}
