package main

import (
	"fmt"
	"github.com/TheJumpCloud/jcapi"
	"os"
)

const (
	apiUrl string = "https://console.jumpcloud.com/api"
)

func outFirst(data string) {
	fmt.Printf("\"%s\"", data)
}

func out(data string) {
	fmt.Printf(",\"%s\"", data)
}

func endLine() {
	fmt.Printf("\n")
}

func main() {
	apiKey := os.Getenv("JUMPCLOUD_APIKEY")
	if apiKey == "" {
		fmt.Printf("%s: Please run: export JUMPCLOUD_APIKEY=<your-JumpCloud-API-key>\n")
		os.Exit(1)
	}

	jc := jcapi.NewJCAPI(apiKey, apiUrl)

	// Grab all system users with their tags
	userList, err := jc.GetSystemUsers(true)
	if err != nil {
		fmt.Printf("Could not read system users, err='%s'\n", err)
		return
	}

	outFirst("Username")
	out("FirstName")
	out("LastName")
	out("Email")
	out("UID")
	out("GID")
	out("Activated")
	out("PasswordExpired")
	out("Sudo")
	out("Tags")
	endLine()

	for _, user := range userList {
		outFirst(user.UserName)
		out(user.FirstName)
		out(user.LastName)
		out(user.Email)
		out(user.Uid)
		out(user.Gid)
		out(fmt.Sprintf("%t", user.Activated))
		out(fmt.Sprintf("%t", user.PasswordExpired))
		out(fmt.Sprintf("%t", user.Sudo))

		for _, tag := range user.Tags {
			out(tag.Name)
		}

		endLine()
	}
}
