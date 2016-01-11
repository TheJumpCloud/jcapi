package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/TheJumpCloud/jcapi"
)

const (
	apiUrl string = "https://console.jumpcloud.com/api"
)

type systemMapToUserMap map[string]map[string]struct{}

func main() {
	apiKey := os.Getenv("JUMPCLOUD_APIKEY")
	if apiKey == "" {
		log.Fatalf("%s: Please run:\n\n\texport JUMPCLOUD_APIKEY=<your-JumpCloud-API-key>\n", os.Args[0])
	}

	jc := jcapi.NewJCAPI(apiKey, apiUrl)

	tags, err := jc.GetAllTags()
	if err != nil {
		log.Fatalf("Could not get tags from your JumpCloud account, err='%s'", err)
	}

	systemUserMap := make(systemMapToUserMap)

	// Walk the tags, and map each system to a map of users
	for _, tag := range tags {
		for _, systemId := range tag.Systems {
			for _, userId := range tag.SystemUsers {
				if systemUserMap[systemId] == nil {
					systemUserMap[systemId] = make(map[string]struct{})
				}

				systemUserMap[systemId][userId] = struct{}{}
			}
		}
	}

	csvWriter := csv.NewWriter(os.Stdout)
	defer csvWriter.Flush()

	headers := []string{"SystemId", "DisplayName", "HostName", "Active", "Instance ID", "OS", "OSVersion",
		"AgentVersion", "CreatedDate", "LastContactDate", "Users"}

	csvWriter.Write(headers)

	for systemId, userMap := range systemUserMap {
		var system jcapi.JCSystem

		system, err = jc.GetSystemById(systemId, false)
		if err != nil {
			log.Fatalf("Could not retrieve system for ID '%s', err='%s'", system.Id, err)
		}

		outLine := []string{system.Id, system.DisplayName, system.Hostname, fmt.Sprintf("%t", system.Active),
			system.AmazonInstanceID, system.Os, system.Version, system.AgentVersion, system.Created,
			system.LastContact}

		for userId, _ := range userMap {
			user, err := jc.GetSystemUserById(userId, false)
			if err != nil {
				log.Fatalf("Could not retrieve system user for ID '%s', err='%s'", userId, err)
			}

			outLine = append(outLine, fmt.Sprintf("%s (%s)", user.UserName, user.Email))
		}

		csvWriter.Write(outLine)
	}
}
