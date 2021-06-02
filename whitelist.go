package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	WhitelistUserByUsername(os.Args[1])
}

func getWhitelistCommand() (*flag.FlagSet, map[string]*string) {
	whitelistCmd := flag.NewFlagSet("whitelist", flag.ExitOnError)
	whitelistFlags := map[string]*string{
		"username": whitelistCmd.String("username", "", "The username to whitelist by."),
		"uuid":     whitelistCmd.String("uuid", "", "The the uuid to whitelist by."),
	}
	return whitelistCmd, whitelistFlags
}

type MinecraftUser struct {
	UUID     string `json:"id"`
	Username string `json:"name"`
}

func UserProfile(user string) []MinecraftUser {
	requestBody, err := json.Marshal([1]string{user})
	if err != nil {
		log.Fatal(err)
	}

	response, err := http.Post(
		"https://api.mojang.com/profiles/minecraft",
		"application/json",
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	userDataBytes := responseBody(response)
	return readUsers(userDataBytes)
}

func responseBody(response *http.Response) []byte {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(body)
		log.Fatal(err)
	}
	return body
}

func readUsers(userData []byte) []MinecraftUser {
	var users []MinecraftUser
	if err := json.Unmarshal(userData, &users); err != nil {
		log.Fatal(err)
	}
	return users
}

type whitelistUser MinecraftUser

type WhitelistEntry struct {
	UUID                string `json:"uuid"`
	Username            string `json:"name"`
	Level               int    `json:"level"`
	BypassesPlayerLimit bool   `json:"bypassesPlayerLimit"`
}

func (user *whitelistUser) makeWhitelistEntry() WhitelistEntry {
	return WhitelistEntry{
		UUID:                user.UUID,
		Username:            user.Username,
		Level:               4,
		BypassesPlayerLimit: false,
	}
}

type Whitelist []WhitelistEntry

func (whitelist *Whitelist) addUser(user whitelistUser) {
	*whitelist = append(*whitelist, user.makeWhitelistEntry())
}

func (whitelist Whitelist) writeToFile() {
	whitelistBytes, err := json.MarshalIndent(whitelist, "", "")
	if err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile("whitelist.json", whitelistBytes, 0644); err != nil {
		log.Fatal(err)
	}
}

func WhitelistUserByUsername(username string) {
	user := whitelistUser(UserProfile(username)[0])
	whitelist := readWhitelistFile()
	whitelist.addUser(user)
	whitelist.writeToFile()
}

func readWhitelistFile() Whitelist {
	whitelistFile, err := os.Open("whitelist.json")
	if err != nil {
		log.Fatal(err)
	}
	defer whitelistFile.Close()

	whitelistBytes, err := ioutil.ReadAll(whitelistFile)
	if err != nil {
		log.Fatal(err)
	}

	var whitelist []WhitelistEntry
	if err := json.Unmarshal(whitelistBytes, &whitelist); err != nil {
		log.Fatal(err)
	}

	return whitelist
}
