package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/TheJumpCloud/jcapi"
	jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"
)

const (
	apiUrlDefault string = "https://console.jumpcloud.com/api"
)

func main() {
	var apiKey string
	var apiUrl string

	// Obtain the input parameters: api key and url (if we want to override the default url)
	flag.StringVar(&apiKey, "key", "", "-key=<API-key-value>")
	flag.StringVar(&apiUrl, "url", apiUrlDefault, "-url=<jumpcloud-api-url>")
	flag.Parse()

	// if the api key isn't specified, try to obtain it through environment variable:
	if apiKey == "" {
		apiKey = os.Getenv("JUMPCLOUD_APIKEY")
	}

	if apiKey == "" {
		fmt.Println("Usage:")
		fmt.Println("  -key=\"\": -key=<API-key-value>")
		fmt.Println("  -url=\"\": -url=<jumpcloud-api-url> (optional)")
		fmt.Println("You can also set the API key via the JUMPCLOUD_APIKEY environment variable:")
		fmt.Println("Run: export JUMPCLOUD_APIKEY=<your-JumpCloud-API-key>")
		return
	}

	// check if this org is on Groups or Tags:
	isGroups, err := isGroupsOrg(apiUrl, apiKey)
	if err != nil {
		log.Fatalf("Could not determine your org type, err='%s'\n", err)
	}
	// if we're on a groups org, instantiate API v2 objects for systems and system groups
	// which we'll need  to list the parent system groups for a given system:
	var systemsAPIv2 *jcapiv2.SystemsApi
	var systemGroupsAPIv2 *jcapiv2.SystemGroupsApi
	if isGroups {
		systemsAPIv2 = jcapiv2.NewSystemsApiWithBasePath(apiUrl + "/v2")
		systemsAPIv2.Configuration.APIKey["x-api-key"] = apiKey
		systemGroupsAPIv2 = jcapiv2.NewSystemGroupsApiWithBasePath(apiUrl + "/v2")
		systemGroupsAPIv2.Configuration.APIKey["x-api-key"] = apiKey
	}

	// instantiate a jcapi v1 object for all v1 endpoints:
	jcapiv1 := jcapi.NewJCAPI(apiKey, apiUrl)

	// Grab all systems (with their tags for a Tags)
	systems, err := jcapiv1.GetSystems(!isGroups)
	if err != nil {
		log.Fatalf("Could not read systems, err='%s'\n", err)
	}

	csvWriter := csv.NewWriter(os.Stdout)
	defer csvWriter.Flush()

	headers := []string{"Id", "DisplayName", "HostName", "Active", "Instance ID", "OS", "OSVersion",
		"AgentVersion", "CreatedDate", "LastContactDate"}

	if isGroups {
		headers = append(headers, "SystemGroups")
	} else {
		headers = append(headers, "Tags")
	}

	csvWriter.Write(headers)

	for _, system := range systems {
		outLine := []string{system.Id, system.DisplayName, system.Hostname, fmt.Sprintf("%t", system.Active),
			system.AmazonInstanceID, system.Os, system.Version, system.AgentVersion, system.Created,
			system.LastContact}

		if isGroups {
			// for a Groups org, let's retrieve the system groups this system is a member of:
			// NOTE: there are more associations for a system in a Groups org we may want to list here as well:
			// Policies, direct Users associations, etc
			var graphs []jcapiv2.GraphObjectWithPaths
			for skip := 0; skip == 0 || len(graphs) == searchLimit; skip += searchSkipInterval {
				graphs, _, err := systemsAPIv2.GraphSystemMemberOf(system.Id, contentType, accept, int32(searchLimit), int32(skip))
				if err != nil {
					fmt.Printf("Could not retrieve parent groups for system %s, err='%s'\n", system.Id, err)
				} else {
					// add the retrieved system groups names to the list for the current system:
					for _, graph := range graphs {
						// get the details of the current system group:
						systemGroup, _, err := systemGroupsAPIv2.GroupsSystemGet(graph.Id, contentType, accept)
						if err != nil {
							// just log a message and skip the system group if there's an error retrieving details:
							fmt.Printf("Could not retrieve info for system group ID %s, err='%s'\n", graph.Id, err)
						} else {
							outLine = append(outLine, systemGroup.Name)
						}
					}
				}
			}
		} else {
			// for Tags orgs, we've already retrieved the list of tags in GetSystems:
			for _, tag := range system.Tags {
				outLine = append(outLine, tag.Name)
			}
		}

		csvWriter.Write(outLine)
	}
}
