package jcapi

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// If you add a field here make sure to add corresponding logic to getJCUserFieldsFromInterface
type JCUser struct {
	Id                     string    `json:"_id,omitempty"`
	UserName               string    `json:"username,omitempty"`
	FirstName              string    `json:"firstname,omitempty"`
	LastName               string    `json:"lastname,omitempty"`
	Email                  string    `json:"email"`
	Password               string    `json:"password,omitempty"`
	PasswordDate           string    `json:"password_date,omitempty"`
	Activated              bool      `json:"activated"`
	ActivationKey          string    `json:"activation_key"`
	ExpiredWarned          bool      `json:"expired_warned"`
	PasswordExpired        bool      `json:"password_expired"`
	PasswordExpirationDate time.Time `json:"password_expiration_date,omitempty"`
	PendingProvisioning    bool      `json:"pendingProvisioning,omitempty"`
	Sudo                   bool      `json:"sudo"`
	Uid                    string    `json:"unix_uid"`
	Gid                    string    `json:"unix_guid"`
	EnableManagedUid       bool      `json:"enable_managed_uid"`

	TagIds []string `json:"tags,omitempty"` // the list of tag IDs that this user should be put in

	//
	// For identification as an external user directory source
	//
	ExternallyManaged  bool   `json:"externally_managed"`
	ExternalDN         string `json:"external_dn,omitempty"`
	ExternalSourceType string `json:"external_source_type,omitempty"`

	Tags []JCTag // the list of actual tags the user is in
}

//
// Special request structure for sending activation emails
//
type JCUserEmailRequest struct {
	IsSelectAll bool     `json:"isSelectAll"`
	Models      []JCUser `json:"models"`
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
	returnVal := fmt.Sprintf("JCUSER: Id=[%s] - UserName=[%s] - FName/LName=[%s/%s] - Email=[%s] - sudo=[%t] - Uid=%s - Gid=%s - enableManagedUid=%t\n",
		jcuser.Id, jcuser.UserName, jcuser.FirstName, jcuser.LastName,
		jcuser.Email, jcuser.Sudo, jcuser.Uid, jcuser.Gid, jcuser.EnableManagedUid)

	returnVal += fmt.Sprintf("JCUSER: ExternallyManaged=[%t] - ExternalDN=[%s] - ExternalSourceType=[%s]\n",
		jcuser.ExternallyManaged, jcuser.ExternalDN, jcuser.ExternalSourceType)

	returnVal += fmt.Sprintf("JCUSER: PasswordExpired=[%t] - Active=[%t] - PendingProvisioning=[%t]\n", jcuser.PasswordExpired, jcuser.Activated,
		jcuser.PendingProvisioning)

	for _, tag := range jcuser.Tags {
		returnVal += fmt.Sprintf("\t%s\n", tag.ToString())
	}

	return returnVal
}

func setTagIds(user *JCUser) {
	for idx, _ := range user.Tags {
		user.TagIds = append(user.TagIds, user.Tags[idx].Id)
	}
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

	// Currently returned as float64, not string (though they are posted as string),
	// defect #96322248...
	if floatVal, ok := fields["unix_uid"].(float64); ok {
		user.Uid = strconv.FormatInt(int64(floatVal), 10)
	}
	if floatVal, ok := fields["unix_guid"].(float64); ok {
		user.Gid = strconv.FormatInt(int64(floatVal), 10)
	}

	if _, exists := fields["enable_managed_uid"]; exists {
		user.EnableManagedUid = fields["enable_managed_uid"].(bool)
	}

	if _, exists := fields["password_expired"]; exists {
		user.PasswordExpired = fields["password_expired"].(bool)
	}

	if _, exists := fields["activated"]; exists {
		user.Activated = fields["activated"].(bool)
	}

	if _, exists := fields["pendingProvisioning"]; exists {
		user.PendingProvisioning = fields["pendingProvisioning"].(bool)
	}

	if _, exists := fields["password_date"]; exists {
		user.PasswordDate = fields["password_date"].(string)
	}

	if _, exists := fields["password_expiration_date"]; exists {
		user.PasswordExpirationDate, _ = time.Parse(time.RFC3339, fields["password_expiration_date"].(string))
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

func (jc JCAPI) GetSystemUserById(userId string, withTags bool) (user JCUser, err JCError) {
	url := fmt.Sprintf("/systemusers/%s", userId)

	retVal, err := jc.Get(url)
	if err != nil {
		err = fmt.Errorf("ERROR: Could not get system user by ID '%s', err='%s'", userId, err)
	}

	if retVal != nil {
		getJCUserFieldsFromInterface(retVal.(map[string]interface{}), &user)

		if withTags {
			// I should be able to use err below as the err return value, but there's
			// a compiler bug here in that it thinks a := of err is shadowed here,
			// even though tags should be the only variable declared with the :=
			tags, err2 := jc.GetAllTags()
			if err != nil {
				err = fmt.Errorf("ERROR: Could not get tags, err='%s'", err2)
				return
			}

			user.AddJCTags(tags)
			setTagIds(&user)
		}
	}

	return
}

func (jc JCAPI) GetSystemUsers(withTags bool) (userList []JCUser, err JCError) {
	var returnVal []JCUser

	for skip := 0; skip == 0 || len(returnVal) == searchLimit; skip += searchSkipInterval {
		url := fmt.Sprintf("/systemusers?sort=username&skip=%d&limit=%d", skip, searchLimit)

		jcUserRec, err2 := jc.Get(url)
		if err != nil {
			return nil, fmt.Errorf("ERROR: Post to JumpCloud failed, err='%s'", err2)
		}

		if jcUserRec == nil {
			return nil, fmt.Errorf("ERROR: No users found")
		}

		// We really only care about the ID for the following call...
		returnVal = getJCUsersFromInterface(jcUserRec)

		for i, _ := range returnVal {
			if returnVal[i].Id != "" {

				//
				// Get the rest of the user record, which includes details like
				// the externalDN...
				//
				// We'll get all the tags one time later, so don't get the tags on this call...
				//
				// See above about the compiler error that requires me to use err2 instead of err below...
				//
				detailedUser, err2 := jc.GetSystemUserById(returnVal[i].Id, false)
				if err != nil {
					err = fmt.Errorf("ERROR: Could not get details for user ID '%s', err='%s'", returnVal[i].Id, err2)
					return
				}

				if detailedUser.Id != "" {
					userList = append(userList, detailedUser)
				}
			}
		}
	}

	if withTags {
		tags, err := jc.GetAllTags()
		if err != nil {
			return nil, fmt.Errorf("ERROR: Could not get tags, err='%s'", err)
		}

		for idx, _ := range userList {
			userList[idx].AddJCTags(tags)
			setTagIds(&userList[idx])
		}
	}

	return
}

//
// Resend user email
//
func (jc JCAPI) SendUserActivationEmail(userList []JCUser) (err JCError) {
	for _, user := range userList {
		if user.Id == "" {
			return fmt.Errorf("ERROR: Cannot resend user activation email without a systemuser Id on user %v", user)
		}
	}

	emailRequest := JCUserEmailRequest{
		IsSelectAll: false,
		Models:      userList,
	}

	data, err := json.Marshal(emailRequest)
	if err != nil {
		return fmt.Errorf("ERROR: Could not marshal JCUserEmailRequest object, err='%s'", err)
	}

	url := "/systemusers/reactivate"

	_, err = jc.Do(MapJCOpToHTTP(Insert), url, data)
	if err != nil {
		return fmt.Errorf("ERROR: Could not post resend email request object, err='%s'", err)
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
