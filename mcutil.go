package main

import (
	"mcutil/whitelist"

	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	whitelistCmd, whitelistFlags := getWhitelistCommand()

	if len(os.Args) < 2 {
		log.Fatal("What do I do?")
	}

	switch os.Args[1] {
	case "whitelist":
		whitelistCmd.Parse(os.Args[2:])
	}

	if whitelistCmd.Parsed() {
		if *whitelistFlags["username"] == "" && *whitelistFlags["uuid"] == "" {
			log.Fatal("Need a username or uuid")
		} else if *whitelistFlags["username"] != "" && *whitelistFlags["uuid"] != "" {
			log.Fatal("one at a time plz")
		}

		if *whitelistFlags["username"] != "" {
			whitelist.WhitelistUserByUsername(*whitelistFlags["username"])
		}

		if *whitelistFlags["uuid"] != "" {
			fmt.Println("This isn't supported yet!")
		}
	}
}

func getWhitelistCommand() (*flag.FlagSet, map[string]*string) {
	whitelistCmd := flag.NewFlagSet("whitelist", flag.ExitOnError)
	whitelistFlags := map[string]*string{
		"username": whitelistCmd.String("username", "", "The username to whitelist by."),
		"uuid":     whitelistCmd.String("uuid", "", "The the uuid to whitelist by."),
	}
	return whitelistCmd, whitelistFlags
}
