package jcapi

import "fmt"

type JCSystem struct {
	Os                             string  `json:os`
	TemplateName                   string  `json:templateName`
	AllowSshRootLogin              bool    `json:allowSsgRootLogin`
	Id                             string  `json:id`
	LastContact                    string  `json:lastContact`
	RemoteIP                       string  `json:remoteIP`
	Active                         bool    `json:active`
	SshRootEnabled                 bool    `json:sshRootEnabled`
	AmazonInstanceID               string  `json:amazonInstanceID,omitempty`
	SshPassEnabled                 bool    `json:sshPassEnabled`
	Version                        string  `json:version`
	AgentVersion                   string  `json:agentVersion`
	AllowPublicKeyAuth             bool    `json:allowPublicKeyAuthentication`
	Organization                   string  `json:organization`
	Created                        string  `json:created`
	Arch                           string  `json:arch`
	SystemTimezone                 float64 `json:systemTimeZone`
	AllowSshPasswordAuthentication bool    `json:allowSshPasswordAuthentication`
	DisplayName                    string  `json:displayName`
	ModifySSHDConfig               bool    `json:modifySSHDConfig`
	AllowMultiFactorAuthentication bool    `json:allowMultiFactorAuthentication`
	Hostname                       string  `json:hostname`

	Patches               []string `json:patches`
	SshParamList          []string `json:sshParams`
	PatchAlarmList        []string `json:patchAlarms`
	NetworkInterfaceList  []string `json:networkInterfaces`
	ConnectionHistoryList []string `json:connectionHistory`

	Tags              []JCTag
	SshdParams        []JCSSHDParam
	NetworkInterfaces []JCNetworkInterface
}

type JCSSHDParam struct {
	Name  string `json:name`
	Value string `json:value`
}

type JCNetworkInterface struct {
	Name     string `json:name`
	Internal bool   `json:internal`
	Family   string `json:family`
	Address  string `json:address`
}

func SystemsToString(systems []JCSystem) string {
	returnVal := ""

	for _, system := range systems {
		returnVal += system.ToString()
	}

	return returnVal
}

func (jcsystem JCSystem) ToString() string {
	//
	// WARNING: Never output password via this method, it could be logged in clear text
	//
	//	returnVal := fmt.Sprintf("JCUSER: Id=[%s] - UserName=[%s] - FName/LName=[%s/%s] - Email=[%s] - sudo=[%t] - Uid=%s - Gid=%s - enableManagedUid=%t\n",
	//		jcuser.Id, jcuser.UserName, jcuser.FirstName, jcuser.LastName,
	//		jcuser.Email, jcuser.Sudo, jcuser.Uid, jcuser.Gid, jcuser.EnableManagedUid)
	//
	//	returnVal += fmt.Sprintf("JCUSER: ExternallyManaged=[%t] - ExternalDN=[%s] - ExternalSourceType=[%s]\n",
	//		jcuser.ExternallyManaged, jcuser.ExternalDN, jcuser.ExternalSourceType)
	//
	//	for _, tag := range jcuser.Tags {
	//		returnVal += fmt.Sprintf("\t%s\n", tag.ToString())
	//	}
	returnVal := "I'm a return val"
	return returnVal
}

func getJCSSHDParamFieldsFromInterface(fields map[string]interface{}, params *JCSSHDParam) {
	params.Name = fields["name"].(string)
	params.Value = fields["value"].(string)
}

func getJCSSHDParamFromArray(paramArray []interface{}) []JCSSHDParam {
	var returnVal []JCSSHDParam
	returnVal = make([]JCSSHDParam, len(paramArray))

	for idx, rec := range paramArray {
		getJCSSHDParamFieldsFromInterface(rec.(map[string]interface{}), &returnVal[idx])
	}

	return returnVal

}

func getJCNetworkInterfaceFieldsFromInterface(fields map[string]interface{}, nic *JCNetworkInterface) {
	nic.Address = fields["address"].(string)
	nic.Family = fields["family"].(string)
	nic.Internal = fields["internal"].(bool)
	nic.Name = fields["name"].(string)
}

func getJCNetworkInterfacesFromArray(nicArray []interface{}) []JCNetworkInterface {
	var returnVal []JCNetworkInterface
	returnVal = make([]JCNetworkInterface, len(nicArray))

	for idx, rec := range nicArray {
		getJCNetworkInterfaceFieldsFromInterface(rec.(map[string]interface{}), &returnVal[idx])
	}

	return returnVal
}

func getJCSystemFieldsFromInterface(fields map[string]interface{}, system *JCSystem) {
	system.Os = fields["os"].(string)
	system.TemplateName = fields["templateName"].(string)
	system.AllowSshRootLogin = fields["allowSshRootLogin"].(bool)
	system.Id = fields["id"].(string)
	system.LastContact = fields["lastContact"].(string)
	system.RemoteIP = fields["remoteIP"].(string)
	system.Active = fields["active"].(bool)
	system.SshRootEnabled = fields["sshRootEnabled"].(bool)
	system.SshPassEnabled = fields["sshPassEnabled"].(bool)
	system.Version = fields["version"].(string)
	system.AgentVersion = fields["agentVersion"].(string)
	system.AllowPublicKeyAuth = fields["allowPublicKeyAuthentication"].(bool)
	system.Organization = fields["organization"].(string)
	system.Created = fields["created"].(string)
	system.Arch = fields["arch"].(string)
	system.SystemTimezone = fields["systemTimezone"].(float64)
	system.AllowSshPasswordAuthentication = fields["allowSshPasswordAuthentication"].(bool)
	system.DisplayName = fields["displayName"].(string)
	system.ModifySSHDConfig = fields["modifySSHDConfig"].(bool)
	system.AllowMultiFactorAuthentication = fields["allowMultiFactorAuthentication"].(bool)
	system.Hostname = fields["hostname"].(string)

	system.SshdParams = getJCSSHDParamFromArray(fields["sshdParams"].([]interface{}))
	system.NetworkInterfaces = getJCNetworkInterfacesFromArray(fields["networkInterfaces"].([]interface{}))

	if _, exists := fields["amazonInstanceID"]; exists {
		system.AmazonInstanceID = fields["amazonInstanceID"].(string)
	}
}

func getJCSystemsFromInterface(systemInt interface{}) []JCSystem {

	var returnVal []JCSystem

	recMap := systemInt.(map[string]interface{})

	results := recMap["results"].([]interface{})

	returnVal = make([]JCSystem, len(results))

	for idx, result := range results {
		getJCSystemFieldsFromInterface(result.(map[string]interface{}), &returnVal[idx])
	}

	return returnVal
}

// Executes a search by hostname via the JumpCloud API
func (jc JCAPI) GetSystemByhostname(hostname string, withTags bool) ([]JCSystem, JCError) {
	var returnVal []JCSystem

	jcSystemRec, err := jc.Post("/search/systems", jc.hostnameFilter(hostname))
	if err != nil {
		return nil, fmt.Errorf("ERROR: Post to JumpCloud failed, err='%s'", err)
	}

	returnVal = getJCSystemsFromInterface(jcSystemRec)

	if withTags {
		tags, err := jc.GetAllTags()
		if err != nil {
			return nil, fmt.Errorf("ERROR: Could not get tags, err='%s'", err)
		}

		for idx, _ := range returnVal {
			returnVal[idx].AddJCTagsToSystem(tags)
		}
	}

	return returnVal, nil
}

func (jc JCAPI) GetSystemById(systemId string, withTags bool) (system JCSystem, err JCError) {
	url := fmt.Sprintf("/systems/%s", systemId)

	retVal, err := jc.Get(url)
	if err != nil {
		err = fmt.Errorf("ERROR: Could not get system by ID '%s', err='%s'", systemId, err)
	}

	if retVal != nil {
		getJCSystemFieldsFromInterface(retVal.(map[string]interface{}), &system)

		if withTags {
			// I should be able to use err below as the err return value, but there's
			// a compiler bug here in that it thinks a := of err is shadowed here,
			// even though tags should be the only variable declared with the :=
			tags, err2 := jc.GetAllTags()
			if err != nil {
				err = fmt.Errorf("ERROR: Could not get tags, err='%s'", err2)
				return
			}

			system.AddJCTagsToSystem(tags)
		}
	}

	return
}

//
//func (jc JCAPI) GetSystemUsers(withTags bool) (userList []JCUser, err JCError) {
//	var returnVal []JCUser
//
//	for skip := 0; skip == 0 || len(returnVal) == searchLimit; skip += searchSkipInterval {
//		url := fmt.Sprintf("/systemusers?sort=username&skip=%d&limit=%d", skip, searchLimit)
//
//		jcUserRec, err2 := jc.Get(url)
//		if err != nil {
//			return nil, fmt.Errorf("ERROR: Post to JumpCloud failed, err='%s'", err2)
//		}
//
//		// We really only care about the ID for the following call...
//		returnVal = getJCUsersFromInterface(jcUserRec)
//
//		for i, _ := range returnVal {
//			if returnVal[i].Id != "" {
//
//				//
//				// Get the rest of the user record, which includes details like
//				// the externalDN...
//				//
//				// We'll get all the tags one time later, so don't get the tags on this call...
//				//
//				// See above about the compiler error that requires me to use err2 instead of err below...
//				//
//				detailedUser, err2 := jc.GetSystemUserById(returnVal[i].Id, false)
//				if err != nil {
//					err = fmt.Errorf("ERROR: Could not get details for user ID '%s', err='%s'", returnVal[i].Id, err2)
//					return
//				}
//
//				if detailedUser.Id != "" {
//					userList = append(userList, detailedUser)
//				}
//			}
//		}
//	}
//
//	if withTags {
//		tags, err := jc.GetAllTags()
//		if err != nil {
//			return nil, fmt.Errorf("ERROR: Could not get tags, err='%s'", err)
//		}
//
//		for idx, _ := range userList {
//			userList[idx].AddJCTags(tags)
//		}
//	}
//
//	return
//}
//
////
//// Add or Update a new user to JumpCloud
////
//func (jc JCAPI) AddUpdateUser(op JCOp, user JCUser) (userId string, err JCError) {
//	if user.Password != "" {
//		user.PasswordDate = getTimeString()
//	}
//
//	data, err := json.Marshal(user)
//	if err != nil {
//		return "", fmt.Errorf("ERROR: Could not marshal JCUser object, err='%s'", err)
//	}
//
//	url := "/systemusers"
//	if op == Update {
//		url += "/" + user.Id
//	}
//
//	jcUserRec, err := jc.Do(MapJCOpToHTTP(op), url, data)
//	if err != nil {
//		return "", fmt.Errorf("ERROR: Could not post new JCUser object, err='%s'", err)
//	}
//
//	var returnUser JCUser
//	getJCUserFieldsFromInterface(jcUserRec.(map[string]interface{}), &returnUser)
//
//	if returnUser.Email != user.Email {
//		return "", fmt.Errorf("ERROR: JumpCloud did not return the same email - this should never happen!")
//	}
//
//	userId = returnUser.Id
//
//	return
//}

func (jc JCAPI) DeleteSystem(system JCSystem) JCError {
	_, err := jc.Delete(fmt.Sprintf("/system/%s", system.Id))
	if err != nil {
		return fmt.Errorf("ERROR: Could not delete system '%s': err='%s'", system.Hostname, err)
	}

	return nil
}
