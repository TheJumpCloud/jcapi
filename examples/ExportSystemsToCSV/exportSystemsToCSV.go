package main

import (
	"encoding/csv"
	"fmt"
	"github.com/TheJumpCloud/jcapi"
	"os"
)

const (
	apiUrl string = "https://console.jumpcloud.com/api"
)

func main() {
	apiKey := os.Getenv("JUMPCLOUD_APIKEY")
	if apiKey == "" {
		fmt.Printf("%s: Please run:\n\n\texport JUMPCLOUD_APIKEY=<your-JumpCloud-API-key>\n", os.Args[0])
		os.Exit(1)
	}

	jc := jcapi.NewJCAPI(apiKey, apiUrl)

	// Grab all systems with their tags
	systems, err := jc.GetSystems(true)
	if err != nil {
		fmt.Printf("Could not read systems, err='%s'\n", err)
		return
	}

	csvWriter := csv.NewWriter(os.Stdout)
	defer csvWriter.Flush()

	headers := []string{"Id", "DisplayName", "HostName", "Active", "Instance ID", "OS", "OSVersion",
		"AgentVersion", "CreatedDate", "LastContactDate", "Tags"}

	csvWriter.Write(headers)

	for _, system := range systems {
		outLine := []string{system.Id, system.DisplayName, system.Hostname, fmt.Sprintf("%t", system.Active),
			system.AmazonInstanceID, system.Os, system.Version, system.AgentVersion, system.Created,
			system.LastContact}

		for _, tag := range system.Tags {
			outLine = append(outLine, tag.Name)
		}

		csvWriter.Write(outLine)
	}
}
