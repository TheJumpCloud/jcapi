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

	outFirst("Id")
	out("DisplayName")
	out("HostName")
	out("Active")
	out("Instance ID")
	out("OS")
	out("OSVersion")
	out("AgentVersion")
	out("CreatedDate")
	out("LastContactDate")
	out("Tags")
	endLine()

	for _, system := range systems {
		outFirst(system.Id)
		out(system.DisplayName)
		out(system.Hostname)
		out(fmt.Sprintf("%t", system.Active))
		out(system.AmazonInstanceID)
		out(system.Os)
		out(system.Version)
		out(system.AgentVersion)
		out(system.Created)
		out(system.LastContact)

		for _, tag := range system.Tags {
			out(tag.Name)
		}

		endLine()
	}
}
