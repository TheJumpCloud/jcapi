package jcapi

import (
	"testing"
)

const (
	testAPIKey  string = "a7bb8f46cce718636f31e29a64db6200f0304701"
	testUrlBase string = "https://test-console.jumpcloud.com/api"
)

func MakeTestUser() (user JCUser) {
	user = JCUser{
		UserName:          "testuser",
		FirstName:         "Test",
		LastName:          "User",
		Email:             "testuser@jumpcloud.com",
		Password:          "test!@#$ADSF",
		Activated:         false,
		Sudo:              true,
		Uid:               "26753",
		Gid:               "26753",
		EnableManagedUid:  true,
		TagList:           make([]string, 0),
		ExternallyManaged: false,
	}

	return
}

func TestSystemUsers(t *testing.T) {
	jcapi := NewJCAPI(testAPIKey, testUrlBase)

	newUser := MakeTestUser()

	userId, err := jcapi.AddUpdateUser(Insert, newUser)
	if err != nil {
		t.Fatalf("Could not add new user ('%s'), err='%s'", newUser.ToString(), err)
	}

	t.Logf("Returned userId=%s", userId)

	allUsers, err := jcapi.GetSystemUsers(true)
	if err != nil {
		t.Fatalf("Could not get all system users, err='%s'", err)
	}

	t.Logf("GetSystemUsers() returned %d users", len(allUsers))

	var foundUser int = -1

	for i, user := range allUsers {
		if user.Id == userId {
			foundUser = i
			t.Logf("Matched user[%d]='%s'", i, user.ToString())
		}
	}

	if foundUser == -1 {
		t.Fatalf("Could not find the user ID just added '%s', foundUser=%d", userId, foundUser)
	}

	allUsers[foundUser].Email = "newtestemail@jumpcloud.com"

	newUserId, err := jcapi.AddUpdateUser(Update, allUsers[foundUser])
	if err != nil {
		t.Fatalf("Could not modify email on the just-added user ('%s'), err='%s'", allUsers[foundUser].ToString(), err)
	}

	if userId != newUserId {
		t.Fatalf("The user ID of the updated user changed across updates, this should never happen!")
	}

	err = jcapi.DeleteUser(allUsers[foundUser])
	if err != nil {
		t.Fatalf("Could not delete user ('%s'), err='%s'", allUsers[foundUser].ToString(), err)
	}

	return
}

func MakeTestTag() (tag JCTag) {
	tag = JCTag{
		Name:              "Test tag #1",
		GroupName:         "testtag1",
		Systems:           make([]string, 0),
		SystemUsers:       make([]string, 0),
		Expired:           false,
		Selected:          false,
		ExternallyManaged: false,
	}

	return
}

func TestTags(t *testing.T) {
	jcapi := NewJCAPI(testAPIKey, testUrlBase)

	newTag := MakeTestTag()

	tagId, err := jcapi.AddUpdateTag(Insert, newTag)
	if err != nil {
		t.Fatalf("Could not add new tag ('%s'), err='%s'", newTag.ToString(), err)
	}

	t.Logf("Returned tagId=%d", tagId)

	allTags, err := jcapi.GetAllTags()
	if err != nil {
		t.Fatalf("Could not GetAllTags, err='%s'", err)
	}

	var foundTag int

	for i, tag := range allTags {
		t.Logf("Tag[%d]='%s'", i, tag)
		if tag.Id == tagId {
			foundTag = i
		}
	}

	allTags[foundTag].Name = "Test tag #1 with a name change"

	newTagId, err := jcapi.AddUpdateTag(Update, allTags[foundTag])
	if err != nil {
		t.Fatalf("Could not change the test tag's name, err='%s'", err)
	}

	if tagId != newTagId {
		t.Fatalf("The ID of the tag changed during an update, this shouldn't happen.")
	}

	err = jcapi.DeleteTag(allTags[foundTag])
	if err != nil {
		t.Fatalf("Could not delete the tag I just added ('%s'), err='%s'", allTags[foundTag].ToString(), err)
	}
}

func MakeIDSource() JCIDSource {

	return JCIDSource{
		Name:           "Test Name",
		Type:           "Active Directory",
		Version:        "1.0.0",
		IpAddress:      "127.0.0.1",
		LastUpdateTime: "2014-10-14 23:34:33",
		DN:             "CN=JumpCloud;CN=Users;DC=jumpcloud;DC=com",
		Active:         true,
	}
}

func TestIDSources(t *testing.T) {

	jcapi := NewJCAPI(testAPIKey, testUrlBase)

	e := MakeIDSource()

	result, err := jcapi.AddUpdateIDSource(Insert, e)
	if err != nil {
		t.Fatalf("Could not post a new ID Source object, err='%s'", err)
	}

	t.Logf("Post to idsources API successful, result=%U", result)

	extSourceList, err := jcapi.GetAllIDSources()
	if err != nil {
		t.Fatalf("Could not list all external sources, err='%s'", err)
	}

	for idx, source := range extSourceList {
		t.Logf("Result %d: '%s'", idx, source.ToString())
	}

	eGet, exists, err := jcapi.GetIDSourceByName(e.Name)
	if err != nil {
		t.Fatalf("Could not get an external source by name '%s', err='%s'", e.Name, err)
	} else if exists && eGet.Name != e.Name {
		t.Fatalf("Received name is different ('%s') than what was sent ('%s')", eGet.Name, e.Name)
	} else if !exists {
		t.Fatalf("Could not find the record we just put in '%c'")
	}

	//
	// If there's more than one test object with this name, let's just
	// loop over and delete them until we find no more of them...
	//
	for exists, err = true, nil; exists; eGet, exists, err = jcapi.GetIDSourceByName(e.Name) {
		if err != nil {
			t.Fatalf("ERROR: getIDSourceByName() on '%s' failed, err='%s'", eGet.ToString(), err)
		}

		err = jcapi.DeleteIDSource(eGet)
		if err != nil {
			t.Fatalf("ERROR: Delete on '%s' failed, err='%s'", eGet.ToString(), err)
		}
	}
}
