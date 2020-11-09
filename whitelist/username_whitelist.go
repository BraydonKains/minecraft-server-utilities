package whitelist

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"mcutil/mjapi"
	"os"
)

func WhitelistUserByUsername(username string) {
	user := whitelistUser(mjapi.UserProfile(username)[0])
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
