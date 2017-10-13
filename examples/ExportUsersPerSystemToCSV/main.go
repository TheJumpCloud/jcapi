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

// getUsersBoundToSystemV1 returns the list of users associated with the given system
// for a Tags org using the /systems/<system_id>/users endpoint:
// This endpoint will return all the system-user bindings including those made
// via tags and via direct system-user binding
func getUsersBoundToSystemV1(jcapiv1 *jcapi.JCAPI, systemId string) (userIds []string, err error) {

	systemUserBindings, err := jcapiv1.GetSystemUserBindingsById(systemId)
	if err != nil {
		return userIds, fmt.Errorf("Could not get system user bindings for system %s, err='%s'\n", systemId, err)
	}
	// add the retrieved user Ids to our userIds list:
	for _, systemUserBinding := range systemUserBindings {
		userIds = append(userIds, systemUserBinding.UserId)
	}
	return
}

// getUsersBoundToSystemV2 returns the list of users associated with the given system
// for a Groups org using the /v2/systems/<system_id>/users endpoint:
func getUsersBoundToSystemV2(systemsAPIv2 *jcapiv2.SystemsApi, systemId string) (userIds []string, err error) {
	var graphs []jcapiv2.GraphObjectWithPaths
	for skip := 0; skip == 0 || len(graphs) == searchLimit; skip += searchSkipInterval {
		graphs, _, err := systemsAPIv2.GraphSystemTraverseUser(systemId, contentType, accept, int32(searchLimit), int32(skip))
		if err != nil {
			return userIds, fmt.Errorf("Could not retrieve users for system %s, err='%s'\n", systemId, err)
		}
		// add the retrieved user Ids to our userIds list:
		for _, graph := range graphs {
			userIds = append(userIds, graph.Id)
		}
	}
	return
}

func main() {
	var apiKey string
	var apiUrl string

	// Obtain the input parameters: api key and url (if we want to override the default url)
	flag.StringVar(&apiKey, "key", "", "-key=<API-key-value>")
	flag.StringVar(&apiUrl, "url", apiUrlDefault, "-url=<jumpcloud-api-url>")
	flag.Parse()

	// if the api key isn't specified, try to obtain it through environment variable:
	if apiKey == "" {
		apiKey = os.Getenv(apiKeyEnvVariable)
	}

	if apiKey == "" {
		fmt.Println("Usage:")
		fmt.Println("  -key=\"\": -key=<API-key-value>")
		fmt.Println("  -url=\"\": -url=<jumpcloud-api-url> (optional)")
		fmt.Println("You can also set the API key via the JUMPCLOUD_APIKEY environment variable:")
		fmt.Println("Run: export JUMPCLOUD_APIKEY=<your-JumpCloud-API-key>")
		return
	}

	// instantiate a new API v1 object for all v1 endpoints:
	jcapiv1 := jcapi.NewJCAPI(apiKey, apiUrl)

	// check if this org is on Groups or Tags:
	isGroups, err := isGroupsOrg(apiUrl, apiKey)
	if err != nil {
		log.Fatalf("Could not determine your org type, err='%s'\n", err)
	}
	// if we're on a groups org, instantiate the systems API v2
	// which we'll need  to list the users associated to a system:
	var systemsAPIv2 *jcapiv2.SystemsApi
	if isGroups {
		systemsAPIv2 = jcapiv2.NewSystemsApiWithBasePath(apiUrl + "/v2")
		systemsAPIv2.Configuration.APIKey[apiKeyHeader] = apiKey
	}

	csvWriter := csv.NewWriter(os.Stdout)
	defer csvWriter.Flush()

	headers := []string{"SystemId", "DisplayName", "HostName", "Active", "Instance ID", "OS", "OSVersion",
		"AgentVersion", "CreatedDate", "LastContactDate", "Users"}

	csvWriter.Write(headers)

	// retrieve all the systems (note this is a v1 endpoint):
	systems, err := jcapiv1.GetSystems(false)
	if err != nil {
		log.Fatalf("Could not get systems from your JumpCloud account, err='%s'\n", err)
	}

	for _, system := range systems {

		outLine := []string{system.Id, system.DisplayName, system.Hostname, fmt.Sprintf("%t", system.Active),
			system.AmazonInstanceID, system.Os, system.Version, system.AgentVersion, system.Created,
			system.LastContact}

		var userIds []string

		if isGroups {
			userIds, err = getUsersBoundToSystemV2(systemsAPIv2, system.Id)
		} else {
			userIds, err = getUsersBoundToSystemV1(&jcapiv1, system.Id)
		}

		if err != nil {
			// if we fail to retrieve users for the current system, log a msg:
			log.Printf("Failed to retrieve system user bindings: err='%s'\n", err)
			// make sure we still write the system details before skipping:
			csvWriter.Write(outLine)
			continue
		}

		// get details for each bound user and append it to the current system:
		for _, userId := range userIds {
			user, err := jcapiv1.GetSystemUserById(userId, false)
			if err != nil {
				log.Printf("Could not retrieve system user for ID '%s', err='%s'\n", userId, err)
			} else {
				outLine = append(outLine, fmt.Sprintf("%s (%s)", user.UserName, user.Email))
			}
		}
		csvWriter.Write(outLine)
	}
}
