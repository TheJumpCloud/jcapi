package main

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

func usersToString(users []JCUser) string {
	returnVal := ""

	for _, user := range users {
		returnVal += user.toString()
	}

	return returnVal
}

func (jcuser JCUser) toString() string {
	//
	// WARNING: Never output password via this method, it could be logged in clear text
	//
	returnVal := fmt.Sprintf("id=[%s] - userName=[%s] - email=[%s] - externally_managed=[%t] - sudo=[%t] - Uid=%d - Gid=%d - enableManagedUid=%t\n", jcuser.Id, jcuser.UserName,
		jcuser.Email, jcuser.ExternallyManaged, jcuser.Sudo, jcuser.Uid, jcuser.Gid, jcuser.EnableManagedUid)

	for _, tag := range jcuser.tags {
		returnVal += fmt.Sprintf("\t%s\n", tag.toString())
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

	if WillDebug(3) {
		dbg(3, "recMap[\"results\"]=%U\n\n------\n", recMap["results"])

		for key, value := range recMap {
			dbg(3, "recMap: key=[%s] - value=[%s]\n", key, value)
		}
	}

	results := recMap["results"].([]interface{})

	returnVal = make([]JCUser, len(results))

	for idx, result := range results {
		if WillDebug(3) {
			for key, value := range result.(map[string]interface{}) {
				dbg(3, "results: key=[%s] - value=[%s]\n", key, value)
			}
		}

		getJCUserFieldsFromInterface(result.(map[string]interface{}), &returnVal[idx])
	}

	return returnVal
}

// Executes a search by email via the JumpCloud API
func (jc JCAPI) getSystemUserByEmail(email string, withTags bool) ([]JCUser, JCError) {
	var returnVal []JCUser

	jcUserRec, err := jc.post("/search/systemusers", jc.emailFilter(email))
	if err != nil {
		return nil, fmt.Errorf("ERROR: Post to JumpCloud failed, err='%s'", err)
	}

	returnVal = getJCUsersFromInterface(jcUserRec)

	if withTags {
		tags, err := jc.getAllTags()
		if err != nil {
			return nil, fmt.Errorf("ERROR: Could not get tags, err='%s'", err)
		}

		for idx, _ := range returnVal {
			returnVal[idx].addTags(tags)
		}
	}

	return returnVal, nil
}

func (jc JCAPI) getSystemUsers(withTags bool) ([]JCUser, JCError) {
	var returnVal []JCUser

	jcUserRec, err := jc.post("/search/systemusers", []byte("{}"))
	if err != nil {
		return nil, fmt.Errorf("ERROR: Post to JumpCloud failed, err='%s'", err)
	}

	returnVal = getJCUsersFromInterface(jcUserRec)

	if withTags {
		tags, err := jc.getAllTags()
		if err != nil {
			return nil, fmt.Errorf("ERROR: Could not get tags, err='%s'", err)
		}

		for idx, _ := range returnVal {
			returnVal[idx].addTags(tags)
		}
	}

	return returnVal, nil
}

//
// Add or Update a new user to JumpCloud
//
func (jc JCAPI) addUpdateUser(op JCOp, user JCUser) (string, JCError) {
	if user.Password != "" {
		user.PasswordDate = getTimeString()
	}

	data, err := json.Marshal(user)
	if err != nil {
		return "", fmt.Errorf("ERROR: Could not marshal JCUser object, err='%s'", err)
	}

	url := "/systemusers"
	if op == update {
		url += "/" + user.Id
	}

	jcUserRec, err := jc.do(mapJCOpToHTTP(op), url, data)
	if err != nil {
		return "", fmt.Errorf("ERROR: Could not post new JCUser object, err='%s'", err)
	}

	var returnUser JCUser
	getJCUserFieldsFromInterface(jcUserRec.(map[string]interface{}), &returnUser)

	if returnUser.Email != user.Email {
		return "", fmt.Errorf("ERROR: JumpCloud did not return the same email - this should never happen!")
	}

	return returnUser.Id, nil
}

func (jc JCAPI) deleteUser(user JCUser) JCError {
	_, err := jc.delete(fmt.Sprintf("/systemusers/%s", user.Id))
	if err != nil {
		return fmt.Errorf("ERROR: Could not delete user '%s': err='%s'", user.Email, err)
	}

	return nil
}
