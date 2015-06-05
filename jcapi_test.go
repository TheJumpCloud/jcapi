package jcapi

import (
	"fmt"
	"os"
	"testing"
)

const (
	testUrlBase string = "https://console.jumpcloud.com/api"
	authUrlBase string = "https://auth.jumpcloud.com"
)

var testAPIKey string = os.Getenv("JUMPCLOUD_APIKEY")

func MakeTestUser() (user JCUser) {
	user = JCUser{
		UserName:          "testuser",
		FirstName:         "Test",
		LastName:          "User",
		Email:             "testuser@jumpcloud.com",
		Password:          "test!@#$ADSF",
		Activated:         true,
		Sudo:              true,
		Uid:               "2244",
		Gid:               "2244",
		EnableManagedUid:  true,
		TagList:           make([]string, 0),
		ExternallyManaged: false,
	}

	return
}

//
// Note: This test requires at least one system to be installed on the
// JumpCloud account referenced by the API key.
//
func TestSystems(t *testing.T) {
	jcapi := NewJCAPI(testAPIKey, testUrlBase)

	systems, err := jcapi.GetSystems(true)
	if err != nil {
		t.Fatalf("couldn't get systems, err='%s'", err.Error())
	}
	if len(systems) == 0 {
		t.Fatalf("no systems found")
	}

	t.Logf("%d Systems found\n", len(systems))

	testSystem := systems[0]
	sysByID, err := jcapi.GetSystemById(testSystem.Id, true)
	if testSystem.Id != sysByID.Id {
		t.Fatalf("Got ID='%s', expected '%s'", sysByID.Id, testSystem.Id)
	}
	t.Logf("TestSystem: '%s'", testSystem.ToString())
	var foundSystem int = -1
	sysByHostname, err := jcapi.GetSystemByHostName(testSystem.Hostname, true)
	for i, sys := range sysByHostname {
		if sys.Id == testSystem.Id {
			foundSystem = i
		}
	}
	if foundSystem == -1 {
		t.Fatalf("Didn't find test system '%s', foundSystem=%d", testSystem.Id, foundSystem)
	}
	t.Logf("TestSystem: '%s'", testSystem.ToString())
	tagsBefore := testSystem.Tags
	if len(tagsBefore) == 0 {
		t.Fatalf("no tags in test system :-(")
	}
	allTags, err := jcapi.GetAllTags()
	if err != nil {
		t.Fatalf("couldn't get the tags")
	}
	tagList := make([]string, len(allTags))
	for i, tag := range allTags {
		tagList[i] = tag.Name
	}
	testSystem.TagList = tagList

	updatedSystemId, err := jcapi.UpdateSystem(testSystem)
	if err != nil {
		t.Fatalf("Couldn't update system, err='%s'", err)
	}

	updatedSystem, err := jcapi.GetSystemById(updatedSystemId, true)
	if err != nil {
		t.Fatalf("error getting system")
	}

	tagsAfter := updatedSystem.Tags
	if len(tagsAfter) < len(allTags) {
		t.Fatalf("not enough tags!")
	}
	beforeTagList := make([]string, len(tagsBefore))
	for i, tag := range tagsBefore {
		beforeTagList[i] = tag.Name
	}

	updatedSystem.TagList = beforeTagList
	backToNormalId, err := jcapi.UpdateSystem(updatedSystem)
	if err != nil {
		t.Fatalf("Couldn't update system, err='%s'", err)
	}
	backToNormal, err := jcapi.GetSystemById(backToNormalId, true)
	if err != nil {
		t.Fatalf("error getting system")
	}

	// TODO: compare Tags contents
	if len(backToNormal.Tags) != len(tagsBefore) {
		t.Fatalf("Tags don't match, backToNormal.Tags=%d - tagsBefore=%d", len(backToNormal.Tags), len(tagsBefore))
	}

}

func TestSystemUsersByOne(t *testing.T) {
	jcapi := NewJCAPI(testAPIKey, testUrlBase)

	newUser := MakeTestUser()

	userId, err := jcapi.AddUpdateUser(Insert, newUser)
	if err != nil {
		t.Fatalf("Could not add new user ('%s'), err='%s'", newUser.ToString(), err)
	}

	t.Logf("Returned userId=%s", userId)

	retrievedUser, err := jcapi.GetSystemUserById(userId, true)
	if err != nil {
		t.Fatalf("Could not get the system user I just added, err='%s'", err)
	}

	if userId != retrievedUser.Id {
		t.Fatalf("Got back a different user ID than expected, this shouldn't happen! Initial userId='%s' - returned object: '%s'",
			userId, retrievedUser.ToString())
	}

	retrievedUser.Email = "newtestemail@jumpcloud.com"

	// We have to do the following because of defect #96322248
	retrievedUser.Uid = "2244"
	retrievedUser.Gid = "2244"

	newUserId, err := jcapi.AddUpdateUser(Update, retrievedUser)
	if err != nil {
		t.Fatalf("Could not modify email on the just-added user ('%s'), err='%s'", retrievedUser.ToString(), err)
	}

	if userId != newUserId {
		t.Fatalf("The user ID of the updated user changed across updates, this should never happen!")
	}

	err = jcapi.DeleteUser(retrievedUser)
	if err != nil {
		t.Fatalf("Could not delete user ('%s'), err='%s'", retrievedUser.ToString(), err)
	}

	return
}

func TestSystemUsers(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping potentially long test in short mode")
	}

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

	// We have to do the following because of defect #96322248
	allUsers[foundUser].Uid = "2244"
	allUsers[foundUser].Gid = "2244"

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
	return mockEmptyTag("Test tag 1", "testtag1")
}

func TestTags(t *testing.T) {
	jcapi := NewJCAPI(testAPIKey, testUrlBase)

	newTag := MakeTestTag()

	tagId, err := jcapi.AddUpdateTag(Insert, newTag)
	if err != nil {
		t.Fatalf("Could not add new tag ('%s'), err='%s'", newTag.ToString(), err)
	}

	t.Logf("Returned tagId=%s", tagId)

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

	oneTag, err := jcapi.GetTagByName(allTags[foundTag].Name)
	if err != nil {
		t.Fatalf("Could not get tag by name, '%s', err='%s'", allTags[foundTag].Name, err)
	}
	if oneTag.Name != allTags[foundTag].Name {
		t.Fatalf("Tag names don't match, '%s' != '%s'", oneTag.Name, allTags[foundTag].Name)
	}

	allTags[foundTag].Name = "Test tag 1 with a name change"

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

func checkAuth(t *testing.T, expectedResult bool, username, password, tag string) {
	authjc := NewJCAPI(testAPIKey, authUrlBase)

	userAuth, err := authjc.AuthUser(username, password, tag)
	if err != nil {
		t.Fatalf("Could not authenticate the user '%s' with password '%s' and tag '%s' err='%s'", username, password, tag, err)
	}

	if userAuth != expectedResult {
		t.Fatalf("userAuth=%t, we expected %s for user='%s', pass='%s', tag='%s'", userAuth, expectedResult, username, password, tag)
	}
}

// NOTE: Requires a functional Auth Server for testing, so if your auth server is
// local, make sure it's running and happy.
func TestRestAuth(t *testing.T) {
	jcapi := NewJCAPI(testAPIKey, testUrlBase)

	newUser := MakeTestUser()

	userId, err := jcapi.AddUpdateUser(Insert, newUser)
	if err != nil {
		t.Fatalf("Could not add new user ('%s'), err='%s'", newUser.ToString(), err)
	}
	newUser.Id = userId
	defer jcapi.DeleteUser(newUser)

	t.Logf("Returned userId=%s", userId)

	checkAuth(t, true, newUser.UserName, newUser.Password, "")
	checkAuth(t, false, newUser.UserName, newUser.Password, "mytesttag")
	checkAuth(t, false, newUser.UserName, "a0938mbo", "")
	checkAuth(t, false, "2309vnotauser", newUser.Password, "")
	checkAuth(t, false, "", "", "")

	//
	// Now add a tag and put the user in it, and let's try all the tag checking stuff
	//
	newTag := MakeTestTag()

	newTag.SystemUsers = append(newTag.SystemUsers, userId)

	tagId, err := jcapi.AddUpdateTag(Insert, newTag)
	if err != nil {
		t.Fatalf("Could not add new tag ('%s'), err='%s'", newTag.ToString(), err)
	}
	newTag.Id = tagId
	defer jcapi.DeleteTag(newTag)

	t.Logf("Returned tagId=%d", tagId)

	checkAuth(t, true, newUser.UserName, newUser.Password, newTag.Name)
	checkAuth(t, false, newUser.UserName, newUser.Password, "not a real tag")
	checkAuth(t, true, newUser.UserName, newUser.Password, "")
}

func mockEmptyTag(name, groupName string) (tag JCTag) {
	tag = JCTag{
		Name:              name,
		GroupName:         groupName,
		Systems:           make([]string, 0),
		SystemUsers:       make([]string, 0),
		Expired:           false,
		Selected:          false,
		ExternallyManaged: false,
	}

	return
}

func mockTestRadiusServer(name, ip, secret string, tagList []string) (radiusServer *JCRadiusServer) {
	return &JCRadiusServer{
		Name:            name,
		NetworkSourceIP: ip,
		SharedSecret:    secret,
		TagList:         tagList,
	}
}

const (
	RADIUS_SERVER_COUNT int = 5
)

func testRadiusServerCalls(t *testing.T, jcapi JCAPI, op JCOp, rs *JCRadiusServer, expectedError error) (id string) {
	var err error

	switch op {
	case Insert:
		id, err = jcapi.AddUpdateRadiusServer(op, *rs)
	case Update:
		id, err = jcapi.AddUpdateRadiusServer(op, *rs)
	case Delete:
		err = jcapi.DeleteRadiusServer(*rs)
	}

	opText := MapJCOpToHTTP(op)

	if err != nil && expectedError == nil {
		t.Fatalf("Testing with '%s' of '%s', expected nil, but got '%s'", opText, rs.ToString(), err.Error())
	}
	if err == nil && expectedError != nil {
		t.Fatalf("Testing with '%s' of '%s', expected '%s', but got nil", opText, rs.ToString(), expectedError.Error())
	}
	if err != nil && expectedError != nil && err.Error() != expectedError.Error() {
		t.Fatalf("Testing with '%s' of '%s', expected '%s', but got '%s'", opText, rs.ToString(), expectedError, err.Error())
	}

	return
}

func TestRadiusServer(t *testing.T) {

	jcapi := NewJCAPI(testAPIKey, testUrlBase)

	tagIds := make([]string, RADIUS_SERVER_COUNT)

	//
	// Let's get a few tags added to our account
	//
	for i := 0; i < RADIUS_SERVER_COUNT; i++ {
		tag := mockEmptyTag(fmt.Sprintf(" RS test tag %d", i), "")

		id, err := jcapi.AddUpdateTag(Insert, tag)
		if err != nil {
			t.Fatalf("Could not insert a new test tag for Radius Server test, err='%s'", err.Error())
		}

		tag.Id = id
		defer jcapi.DeleteTag(tag)

		tagIds[i] = id
	}

	// Should be okay on POST/PUT
	rs1 := mockTestRadiusServer("Boulder Network", "12.13.14.15", "my-super-secret", tagIds)

	// should fail on POST/PUT because of space in the secret
	rs2 := mockTestRadiusServer("Denver Network", "34.42.53.22", "another secret", tagIds)

	// Should be okay on POST/PUT
	emptyTags := make([]string, 0)
	rs3 := mockTestRadiusServer("Boston Network", "55.66.23.43", "secret", emptyTags)

	rs1.Id = testRadiusServerCalls(t, jcapi, Insert, rs1, nil)

	testRadiusServerCalls(t, jcapi, Update, rs1, nil)

	testRadiusServerCalls(t, jcapi, Delete, rs1, nil)

	// Test with the other two objects...
	rs2.Id = testRadiusServerCalls(t, jcapi, Insert, rs2, fmt.Errorf("ERROR: Could not post new JCIDSource object, err='JumpCloud HTTP response status='400 Bad Request''"))

	rs3.Id = testRadiusServerCalls(t, jcapi, Insert, rs3, nil)

	testRadiusServerCalls(t, jcapi, Update, rs3, nil)

	testRadiusServerCalls(t, jcapi, Delete, rs3, nil)

	// Validate the get and find functions...
	rs1.Id = testRadiusServerCalls(t, jcapi, Insert, rs1, nil)
	defer testRadiusServerCalls(t, jcapi, Delete, rs1, nil)

	rs3.Id = testRadiusServerCalls(t, jcapi, Insert, rs3, nil)
	defer testRadiusServerCalls(t, jcapi, Delete, rs3, nil)

	radservers, err := jcapi.GetAllRadiusServers()
	if err != nil {
		t.Fatalf("Could not get all the RADIUS servers, err='%s'", err.Error())
	}

	for _, radServer := range radservers {
		switch radServer.Id {
		case rs1.Id:
			if radServer.ToString() != rs1.ToString() {
				t.Fatalf("radServer='%s' - rs1='%s' - string compare failed", radServer.ToString(), rs1.ToString())
			}
		case rs3.Id:
			if radServer.ToString() != rs3.ToString() {
				t.Fatalf("radServer='%s' - rs3='%s' - string compare failed", radServer.ToString(), rs3.ToString())
			}
		}
	}

	foundRs := FindRadiusServerById(radservers, rs1.Id)

	if foundRs == nil {
		t.Fatalf("Could not find expected ID '%s'", rs1.Id)
	}

	if foundRs.ToString() != rs1.ToString() {
		t.Fatalf("Find by ID %s failed, got '%s', but expected '%s'", rs1.Id, foundRs.ToString(), rs1.ToString())
	}

	foundRs = FindRadiusServerById(radservers, rs3.Id)

	if foundRs == nil {
		t.Fatalf("Could not find expected ID '%s'", rs3.Id)
	}

	if foundRs.ToString() != rs3.ToString() {
		t.Fatalf("Find by ID %s failed, got '%s', but expected '%s'", rs3.Id, foundRs.ToString(), rs3.ToString())
	}

	foundRs = FindRadiusServerById(radservers, "id not there")

	if foundRs != nil {
		t.Fatalf("Found an ID we didn't expect and got back '%s'", foundRs.ToString())
	}

	return
}

func mockCommand(name, command, commandType, user string) (cmd JCCommand) {
	cmd = JCCommand{
		Name:        name,
		Command:     command,
		CommandType: commandType,
		User:        user,
		LaunchType:  "manual",
		Schedule:    "immediate",
		Timeout:     "0", // No timeout
		ListensTo:   "",
		Trigger:     "",
		Sudo:        false,
		Skip:        0,
		Limit:       10,
	}

	return
}

// Not an ideal test... depends on the existence of at least one system in the database,
// but without direct DB access, it's not possible to simply add one...
func TestCommands(t *testing.T) {
	jc := NewJCAPI(testAPIKey, testUrlBase)

	c := mockCommand("AAA test command", "/bin/echo \"hello\"", "linux", "000000000000000000000000")

	// Get a Linux system to attach to the command, it requires at least one...
	systems, err := jc.GetSystems(false)
	if err != nil {
		t.Fatalf("Could not get a list of all systems, err='%s'")
	}

	// Find the first Linux host available (doesn't matter what it is)...
	systemIndex, err := FindObjectByStringRegex(GetInterfaceArrayFromJCSystems(systems), "Os", "CentOS|Ubuntu|Amazon|Debian")
	if err != nil {
		t.Fatalf("Could search a list of systems for OS type, err='%s'", err.Error())
	}

	if systemIndex >= 0 {
		c.Systems = append(c.Systems, systems[systemIndex].Id)
	} else {
		t.Skip("No applicable systems to test with, skipping this test")
	}

	id, err := jc.AddUpdateCommand(Insert, c)
	if err != nil {
		t.Fatalf("Could not insert a new command, err='%s'", err.Error())
	}

	if id == "" {
		t.Fatalf("Didn't get back an ID from AddUpdateCommand!")
	}

	c.Id = id

	t.Logf("Returned ID value is '%s'", c.Id)

	commandList, err := jc.GetAllCommands()
	if err != nil {
		t.Fatalf("Could not get all commands, err='%s'", err.Error())
	}

	for idx, cmd := range commandList {
		t.Logf("Command %d='%s'", idx, cmd.ToString())
	}

	foundCommand, index := FindCommandById(commandList, c.Id)
	if foundCommand == nil {
		t.Fatalf("FindCommandByID() returned no matching command for id='%s', index=%d", c.Id, index)
	}

	if foundCommand.Id != c.Id {
		t.Fatalf("FindCommandById() returned '%s', but was expecting '%s'", foundCommand.ToString(), c.ToString())
	}

	err = jc.RunCommand(c)
	if err != nil {
		t.Fatalf("RunCommand() failed on '%s', err='%s'", c.ToString(), err)
	}

	err = jc.DeleteCommand(c)
	if err != nil {
		t.Fatalf("Could not delete command '%s', err='%s'", c.ToString(), err.Error())
	}

	commandList, err = jc.GetAllCommands()
	if err != nil {
		t.Fatalf("Could not get all commands, err='%s'", err.Error())
	}

	for _, cmd := range commandList {
		if cmd.Id == c.Id {
			t.Fatalf("DeleteCommand failed to delete '%s', it's still in the system, found at '%s'", c.ToString(), cmd.ToString())
		}
	}
}
