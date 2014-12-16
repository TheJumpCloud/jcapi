package jcapi

import (
	"encoding/json"
	"fmt"
)

type JCUser struct {
	Id               string `json:"_id,omitempty"`
	UserName         string `json:"username,omitempty"`
	FirstName        string `json:"firstname,omitempty"`
	LastName         string `json:"lastname,omitempty"`
	Email            string `json:"email"`
	Password         string `json:"password,omitempty"`
	PasswordDate     string `json:"password_date,omitempty"`
	Activated        bool   `json:"activated"`
	ActivationKey    string `json:"activation_key"`
	ExpiredWarned    bool   `json:"expired_warned"`
	PasswordExpired  bool   `json:"password_expired"`
	Sudo             bool   `json:"sudo"`
	Uid              string `json:"unix_uid"`
	Gid              string `json:"unix_guid"`
	EnableManagedUid bool   `json:"enable_managed_uid"`

	TagList []string `json:"tags"`

	//
	// For identification as an external user directory source
	//
	ExternallyManaged  bool   `json:"externally_managed"`
	ExternalDN         string `json:"external_dn,omitempty"`
	ExternalSourceType string `json:"external_source_type,omitempty"`

	tags []JCTag
}

func UsersToString(users []JCUser) string {
	returnVal := ""

	for _, user := range users {
		returnVal += user.ToString()
	}

	return returnVal
}

func (jcuser JCUser) ToString() string {
	//
	// WARNING: Never output password via this method, it could be logged in clear text
	//
	returnVal := fmt.Sprintf("id=[%s] - userName=[%s] - email=[%s] - externally_managed=[%t] - sudo=[%t] - Uid=%s - Gid=%s - enableManagedUid=%t\n", jcuser.Id, jcuser.UserName,
		jcuser.Email, jcuser.ExternallyManaged, jcuser.Sudo, jcuser.Uid, jcuser.Gid, jcuser.EnableManagedUid)

	for _, tag := range jcuser.tags {
		returnVal += fmt.Sprintf("\t%s\n", tag.ToString())
	}

	return returnVal
}

func getJCUserFieldsFromInterface(fields map[string]interface{}, user *JCUser) {
	user.Email = fields["email"].(string)

	if _, exists := fields["firstname"]; exists {
		user.FirstName = fields["firstname"].(string)
	}
	if _, exists := fields["lastname"]; exists {
		user.LastName = fields["lastname"].(string)
	}

	user.UserName = fields["username"].(string)
	user.Id = fields["_id"].(string)

	if _, exists := fields["externally_managed"]; exists {
		user.ExternallyManaged = fields["externally_managed"].(bool)
	} else {
		user.ExternallyManaged = false
	}

	user.Sudo = fields["sudo"].(bool)

	if _, exists := fields["external_dn"]; exists {
		user.ExternalDN = fields["external_dn"].(string)
	}

	if _, exists := fields["external_source_type"]; exists {
		user.ExternalSourceType = fields["external_source_type"].(string)
	}

	if _, exists := fields["unix_uid"]; exists {
		user.Uid = getStringOrNil(fields["unix_uid"])
	}
	if _, exists := fields["unix_gid"]; exists {
		user.Gid = getStringOrNil(fields["unix_gid"])
	}
	if _, exists := fields["enable_managed_uid"]; exists {
		user.EnableManagedUid = fields["enable_managed_uid"].(bool)
	}
}

func getJCUsersFromInterface(userInt interface{}) []JCUser {

	var returnVal []JCUser

	recMap := userInt.(map[string]interface{})

	results := recMap["results"].([]interface{})

	returnVal = make([]JCUser, len(results))

	for idx, result := range results {
		getJCUserFieldsFromInterface(result.(map[string]interface{}), &returnVal[idx])
	}

	return returnVal
}

// Executes a search by email via the JumpCloud API
func (jc JCAPI) GetSystemUserByEmail(email string, withTags bool) ([]JCUser, JCError) {
	var returnVal []JCUser

	jcUserRec, err := jc.Post("/search/systemusers", jc.emailFilter(email))
	if err != nil {
		return nil, fmt.Errorf("ERROR: Post to JumpCloud failed, err='%s'", err)
	}

	returnVal = getJCUsersFromInterface(jcUserRec)

	if withTags {
		tags, err := jc.GetAllTags()
		if err != nil {
			return nil, fmt.Errorf("ERROR: Could not get tags, err='%s'", err)
		}

		for idx, _ := range returnVal {
			returnVal[idx].AddJCTags(tags)
		}
	}

	return returnVal, nil
}

func (jc JCAPI) GetSystemUsers(withTags bool) (userList []JCUser, err JCError) {
	var returnVal []JCUser

	for skip := 0; skip == 0 || len(returnVal) == searchLimit; skip += searchSkipInterval {
		url := fmt.Sprintf("/systemusers?sort=username&skip=%d&limit=%d", skip, searchLimit)

		jcUserRec, err := jc.Get(url)
		if err != nil {
			return nil, fmt.Errorf("ERROR: Post to JumpCloud failed, err='%s'", err)
		}

		returnVal = getJCUsersFromInterface(jcUserRec)

		for i, _ := range returnVal {
			userList = append(userList, returnVal[i])
		}
	}

	if withTags {
		tags, err := jc.GetAllTags()
		if err != nil {
			return nil, fmt.Errorf("ERROR: Could not get tags, err='%s'", err)
		}

		for idx, _ := range userList {
			userList[idx].AddJCTags(tags)
		}
	}

	return
}

//
// Add or Update a new user to JumpCloud
//
func (jc JCAPI) AddUpdateUser(op JCOp, user JCUser) (userId string, err JCError) {
	if user.Password != "" {
		user.PasswordDate = getTimeString()
	}

	data, err := json.Marshal(user)
	if err != nil {
		return "", fmt.Errorf("ERROR: Could not marshal JCUser object, err='%s'", err)
	}

	url := "/systemusers"
	if op == Update {
		url += "/" + user.Id
	}

	jcUserRec, err := jc.Do(MapJCOpToHTTP(op), url, data)
	if err != nil {
		return "", fmt.Errorf("ERROR: Could not post new JCUser object, err='%s'", err)
	}

	var returnUser JCUser
	getJCUserFieldsFromInterface(jcUserRec.(map[string]interface{}), &returnUser)

	if returnUser.Email != user.Email {
		return "", fmt.Errorf("ERROR: JumpCloud did not return the same email - this should never happen!")
	}

	userId = returnUser.Id

	return
}

func (jc JCAPI) DeleteUser(user JCUser) JCError {
	_, err := jc.Delete(fmt.Sprintf("/systemusers/%s", user.Id))
	if err != nil {
		return fmt.Errorf("ERROR: Could not delete user '%s': err='%s'", user.Email, err)
	}

	return nil
}
