package jcapi

import (
	"encoding/json"
	"fmt"
)

type JCSystem struct {
	Os                             string  `json:os`
	TemplateName                   string  `json:templateName`
	AllowSshRootLogin              bool    `json:allowSshRootLogin`
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

	TagList               []string `json:"tags"`
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
	returnVal := fmt.Sprintf("JCSystem: OS=[%s] - TemplateName=[%s] - ID=[%s] - RemoteIP=[%s] - LastContact=[%v] - Version=%s - DisplayName=%s - Hostname=%s - Arch=%s\n",
		jcsystem.Os, jcsystem.TemplateName, jcsystem.Id, jcsystem.RemoteIP, jcsystem.LastContact,
		jcsystem.Version, jcsystem.DisplayName, jcsystem.Hostname, jcsystem.Arch)

	for _, tag := range jcsystem.Tags {
		returnVal += fmt.Sprintf("\t%s\n", tag.ToString())
	}

	return returnVal
}

func (jcsystem JCSystem) SystemHasTag(tagName string) (hasTag bool, tagId string) {
	for _, tag := range jcsystem.Tags {
		if tag.Name == tagName {
			hasTag = true
			tagId = tag.Id
			return
		}
	}

	return
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
	if _, exists := fields["address"]; exists {
		nic.Address = fields["address"].(string)
	}
	if _, exists := fields["family"]; exists {
		nic.Family = fields["family"].(string)
	}
	if _, exists := fields["internal"]; exists {
		nic.Internal = fields["internal"].(bool)
	}
	if _, exists := fields["name"]; exists {
		nic.Name = fields["name"].(string)
	}
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
	// doing this b/c the jsonSysRec that's returned from the update only has a subset
	// of the fields
	if _, exists := fields["os"]; exists {
		system.Os = fields["os"].(string)
	}
	if _, exists := fields["templateName"]; exists {
		system.TemplateName = fields["templateName"].(string)
	}
	if _, exists := fields["allowSsgRootLogin"]; exists {
		system.AllowSshRootLogin = fields["allowSshRootLogin"].(bool)
	}
	if _, exists := fields["id"]; exists {
		system.Id = fields["id"].(string)
	} else if _, exists := fields["_id"]; exists {
		system.Id = fields["_id"].(string)

	}
	if _, exists := fields["lastContact"]; exists {
		system.LastContact = fields["lastContact"].(string)
	}
	if _, exists := fields["remoteIP"]; exists {
		system.RemoteIP = fields["remoteIP"].(string)
	}
	if _, exists := fields["active"]; exists {
		system.Active = fields["active"].(bool)
	}
	if _, exists := fields["sshRootEnabled"]; exists {
		system.SshRootEnabled = fields["sshRootEnabled"].(bool)
	}
	if _, exists := fields["sshPassEnabled"]; exists {
		system.SshPassEnabled = fields["sshPassEnabled"].(bool)
	}
	if _, exists := fields["version"]; exists {
		system.Version = fields["version"].(string)
	}
	if _, exists := fields["agentVersion"]; exists {
		system.AgentVersion = fields["agentVersion"].(string)
	}
	if _, exists := fields["allowPublicKeyAuthentication"]; exists {
		system.AllowPublicKeyAuth = fields["allowPublicKeyAuthentication"].(bool)
	}
	if _, exists := fields["organization"]; exists {
		system.Organization = fields["organization"].(string)
	}
	if _, exists := fields["created"]; exists {
		system.Created = fields["created"].(string)
	}
	if _, exists := fields["arch"]; exists {
		system.Arch = fields["arch"].(string)
	}
	if _, exists := fields["systemTimeZone"]; exists {
		system.SystemTimezone = fields["systemTimezone"].(float64)
	}
	if _, exists := fields["allowSshPasswordAuthentication"]; exists {
		system.AllowSshPasswordAuthentication = fields["allowSshPasswordAuthentication"].(bool)
	}
	if _, exists := fields["displayName"]; exists {
		system.DisplayName = fields["displayName"].(string)
	}
	if _, exists := fields["modifySSHDConfig"]; exists {
		system.ModifySSHDConfig = fields["modifySSHDConfig"].(bool)
	}
	if _, exists := fields["allowMultiFactorAuthentication"]; exists {
		system.AllowMultiFactorAuthentication = fields["allowMultiFactorAuthentication"].(bool)
	}
	if _, exists := fields["hostname"]; exists {
		system.Hostname = fields["hostname"].(string)
	}

	if _, exists := fields["sshdParams"]; exists {
		system.SshdParams = getJCSSHDParamFromArray(fields["sshdParams"].([]interface{}))
	}
	if _, exists := fields["networkInterfaces"]; exists {
		system.NetworkInterfaces = getJCNetworkInterfacesFromArray(fields["networkInterfaces"].([]interface{}))
	}

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
func (jc JCAPI) GetSystemByHostName(hostname string, withTags bool) ([]JCSystem, JCError) {
	var returnVal []JCSystem

	jcSystemRec, err := jc.Post("/search/systems", jc.hostnameFilter(hostname))

	if err != nil {
		return nil, fmt.Errorf("ERROR: Post to JumpCloud failed, err='%s'", err)
	}

	if jcSystemRec == nil {
		return nil, fmt.Errorf("ERROR: No systems found")
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

func (jc JCAPI) GetSystems(withTags bool) ([]JCSystem, JCError) {
	var returnVal []JCSystem

	for skip := 0; skip == 0 || len(returnVal) == searchLimit; skip += searchSkipInterval {
		url := fmt.Sprintf("/systems?sort=hostname&skip=%d&limit=%d", skip, searchLimit)

		jcSysRec, err2 := jc.Get(url)

		if err2 != nil {
			return nil, fmt.Errorf("ERROR: Get to JumpCloud failed, err='%s'", err2)
		}

		if jcSysRec == nil {
			return nil, fmt.Errorf("ERROR: No systems found")
		}

		returnVal = getJCSystemsFromInterface(jcSysRec)

		if withTags {
			tags, err := jc.GetAllTags()
			if err != nil {
				return nil, fmt.Errorf("ERROR: Could not get tags, err='%s'", err)
			}

			for idx, _ := range returnVal {
				returnVal[idx].AddJCTagsToSystem(tags)
			}
		}

	}
	return returnVal, nil
}

//
// Update a system
//
func (jc JCAPI) UpdateSystem(system JCSystem) (systemId string, err JCError) {
	data, err := json.Marshal(system)
	if err != nil {
		return "", fmt.Errorf("ERROR: Could not marshal JCSystem object, err='%s'", err)
	}
	url := "/systems/" + system.Id

	jcSysRec, err := jc.Do(MapJCOpToHTTP(Update), url, data)
	if err != nil {
		return "", fmt.Errorf("ERROR: Could not update JCSystem object, err='%s'", err)
	}
	var returnSystem JCSystem
	getJCSystemFieldsFromInterface(jcSysRec.(map[string]interface{}), &returnSystem)

	if returnSystem.Id != system.Id {
		return "", fmt.Errorf("ERROR: JumpCloud did not return the same ID - this should never happen!")
	}
	systemId = returnSystem.Id
	return systemId, nil
}

//!!!!!!!!!!!!WARNING!!!!!!!!!!!!
//This will cause JumpCloud to uninstall the agent on this system
//You will lose control of the system after the call returns
//Seriously, It'll be gone
func (jc JCAPI) DeleteSystem(system JCSystem) JCError {
	_, err := jc.Delete(fmt.Sprintf("/system/%s", system.Id))
	if err != nil {
		return fmt.Errorf("ERROR: Could not delete system '%s': err='%s'", system.Hostname, err)
	}

	return nil
}
