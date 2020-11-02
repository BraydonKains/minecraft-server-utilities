package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type MinecraftUser struct {
	UUID     string `json:"id"`
	Username string `json:"name"`
}

type WhitelistEntry struct {
	UUID                string `json:"uuid"`
	Username            string `json:"name"`
	Level               int    `json:"level"`
	BypassesPlayerLimit bool   `json:"bypassesPlayerLimit"`
}

func (user MinecraftUser) makeWhitelistEntry() WhitelistEntry {
	return WhitelistEntry{
		UUID:                user.UUID,
		Username:            user.Username,
		Level:               4,
		BypassesPlayerLimit: false,
	}
}

func main() {
	requestedUsername := os.Args[1]

	apiResponse := userAPICall(requestedUsername)
	user := getReturnedUser(apiResponse)
	whitelist := getWhitelist()

	whitelist = append(whitelist, user.makeWhitelistEntry())

	writeWhitelist(whitelist)
}

func userAPICall(user string) []byte {
	requestBody, err := json.Marshal([1]string{user})
	if err != nil {
		log.Fatal(err)
	}

	response, err := http.Post("https://api.mojang.com/profiles/minecraft", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(body)
		log.Fatal(err)
	}

	return body
}

func getReturnedUser(userData []byte) MinecraftUser {
	var users []MinecraftUser
	if err := json.Unmarshal(userData, &users); err != nil {
		log.Fatal(err)
	}

	return users[0]
}

func getWhitelist() []WhitelistEntry {
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

func writeWhitelist(whitelist []WhitelistEntry) {
	whitelistBytes, err := json.MarshalIndent(whitelist, "", "")
	if err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile("whitelist.json", whitelistBytes, 0644); err != nil {
		log.Fatal(err)
	}
}
