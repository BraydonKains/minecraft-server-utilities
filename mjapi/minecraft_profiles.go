package mjapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

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
