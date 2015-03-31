package main

import (
	"fmt"
	"github.com/TheJumpCloud/jcapi"
)

const (
	// Change this to your own API key, here
	apiKey  string = "<your-api-key-here>"
	urlBase string = "https://console.jumpcloud.com/api"
)

//
// This program will "release" all AD-owned user accounts that have been
// imported by the JumpCloud AD Bridge agent. This is useful if you'd like to
// use the AD Bridge agent to bring all your user accounts and their group memberships
// into JumpCloud, and then disable your AD server, and manage users from JumpCloud.
//
// This script will have no effect on your JumpCloud account if you have no AD-managed
// users in your account.
//
func main() {
	jc := jcapi.NewJCAPI(apiKey, urlBase)

	userList, err := jc.GetSystemUsers(false)
	if err != nil {
		fmt.Printf("Could not read system users, err='%s'\n", err)
		return
	}

	var updateCount = 0

	for i, _ := range userList {
		if userList[i].ExternallyManaged == true {
			userList[i].ExternallyManaged = false
			userList[i].ExternalDN = ""
			userList[i].ExternalSourceType = ""

			userId, err := jc.AddUpdateUser(3, userList[i])
			if err != nil {
				fmt.Printf("Could not update user '%s', err='%s'", userList[i].ToString(), err)
				return
			} else {
				fmt.Printf("Updated user ID '%s'\n", userId)
			}

			updateCount++
		}
	}

	fmt.Printf("%d users released\n", updateCount)

	return
}
