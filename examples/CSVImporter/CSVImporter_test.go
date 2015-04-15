package main

import (
	"os"
	"testing"
	"github.com/TheJumpCloud/jcapi"
)

const (
	testUrlBase string = "https://console.jumpcloud.com/api"
	authUrlBase string = "https://auth.jumpcloud.com"
)

var testAPIKey string = os.Getenv("JUMPCLOUD_APIKEY")
var testSystemID string = os.Getenv("JUMPCLOUD_SYSTEMID")


func TestCSVImporter(t *testing.T) {
	// Attach to JumpCloud
	jc := jcapi.NewJCAPI(testAPIKey, testUrlBase)

	// Fetch all users in JumpCloud
	userList, err := jc.GetSystemUsers(false)

	if err != nil {
		t.Fatalf("Could not read system users, err='%s'\n", err)
		return
	}

	// Fetch our system from JumpCloud
	system, err := jc.GetSystemById(testSystemID, true)

	if err != nil {
		t.Fatalf("Could not read system info for ID='%s', err='%s'\n", testSystemID, err)
		return
	}

	if system.Hostname == "" {
		t.Fatalf("Could not read system info for ID='%s', err='%s'\n", testSystemID, err)
		return
	}

	// Create a CSV record to add a test user
	csvrec := []string{"Joe", "Smith", "js", "TheMan@jumpcloud.com", "", "", "T", "", ""}

	// Process this request record
	ProcessCSVRecord(jc, userList, csvrec)

	// Fetch our freshly minted user
	ourUserList, err := jc.GetSystemUserByEmail("TheMan@jumpcloud.com", true)

	if err != nil {
		t.Fatalf("Could not read system user, err='%s'\n", err)
		return
	}

	tempUserId := GetUserIdFromUserName(ourUserList, "js")

	if tempUserId == "" {
		t.Fatalf("Could not read system user, err='%s'\n", err)
		return
	}

	tempUser, err := jc.GetSystemUserById(tempUserId, true)

	if err != nil {
		t.Fatalf("Could not read system user, err='%s'\n", err)
		return
	}

	// Ensure the user has no associated tags
	if len(tempUser.Tags) > 0 {
		t.Fatalf("Unexpectedly found tags associated with user\n")
		return
	}

	// Re-fetch all users in JumpCloud
	userList, err = jc.GetSystemUsers(false)

	if err != nil {
		t.Fatalf("Could not re-read system users, err='%s'\n", err)
		return
	}

	// Update our user to add a tag
	csvrec = []string{"Joe", "Smith", "js", "TheMan@jumpcloud.com", "", "", "T", "", system.Hostname}

	ProcessCSVRecord(jc, userList, csvrec)

	// Refetch our user...they should now have a tag associated with the host and tag name we provided
	tempUser, err = jc.GetSystemUserById(tempUserId, true)

	if err != nil {
		t.Fatalf("Could not read system user, err='%s'\n", err)
		return
	}

	if len(tempUser.Tags) <= 0 {
		t.Fatalf("No tags associated with user\n")
		return
	}

	foundIt := false
	testTagName := system.Hostname + " - Joe Smith"

	for _, checkTag := range tempUser.Tags {
		if checkTag.Name == testTagName {
			for _, checkTagHost := range checkTag.Systems {
				if checkTagHost == testSystemID {
					foundIt = true
				}
			}
		}
	}

	if !foundIt {
		t.Fatalf("Did not find expected tag associated with user\n")
		return
	}

	return
}
