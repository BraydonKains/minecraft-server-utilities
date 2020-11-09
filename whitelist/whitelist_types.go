package whitelist

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"mcutil/mjapi"
)

type whitelistUser mjapi.MinecraftUser

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
