package jcapi

import (
	"encoding/json"
	"fmt"
)

const (
	SYSTEMS_PATH string = "/systems"
)

type JCSystemResults struct {
	Results []JCSystem `json:"results"`
}

type JCSystem struct {
	Os                             string  `json:"os,omitempty"`
	TemplateName                   string  `json:"templateName,omitempty"`
	AllowSshRootLogin              bool    `json:"allowSshRootLogin"`
	Id                             string  `json:"_id"`
	LastContact                    string  `json:"lastContact,omitempty"`
	RemoteIP                       string  `json:"remoteIP,omitempty"`
	Active                         bool    `json:"active,omitempty"`
	SshRootEnabled                 bool    `json:"sshRootEnabled"`
	AmazonInstanceID               string  `json:"amazonInstanceID,omitempty"`
	SshPassEnabled                 bool    `json:"sshPassEnabled,omitempty"`
	Version                        string  `json:"version,omitempty"`
	AgentVersion                   string  `json:"agentVersion,omitempty"`
	AllowPublicKeyAuth             bool    `json:"allowPublicKeyAuthentication"`
	Organization                   string  `json:"organization,omitempty"`
	Created                        string  `json:"created,omitempty"`
	Arch                           string  `json:"arch,omitempty"`
	SystemTimezone                 float64 `json:"systemTimeZone,omitempty"`
	AllowSshPasswordAuthentication bool    `json:"allowSshPasswordAuthentication"`
	DisplayName                    string  `json:"displayName"`
	ModifySSHDConfig               bool    `json:"modifySSHDConfig"`
	AllowMultiFactorAuthentication bool    `json:"allowMultiFactorAuthentication"`
	Hostname                       string  `json:"hostname,omitempty"`

	ConnectionHistoryList []string             `json:"connectionHistory,omitempty"`
	SshdParams            []JCSSHDParam        `json:"sshdParams,omitempty"`
	NetworkInterfaces     []JCNetworkInterface `json:"networkInterfaces, omitempty"`

	// Derived by JCAPI
	TagList []string `json:"tags,omitempty"`
	Tags    []JCTag
}

// this struct describes a system user binding :
type SystemUserBinding struct {
	UserId   string   `json:"id,omitempty"`
	Username string   `json:"username,omitempty"`
	Tags     []string `json:"tags,omitempty"`
}

type JCSSHDParam struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type JCNetworkInterface struct {
	Name     string `json:"name"`
	Internal bool   `json:"internal"`
	Family   string `json:"family"`
	Address  string `json:"address"`
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

func GetInterfaceArrayFromJCSystems(systems []JCSystem) (interfaceArray []interface{}) {
	interfaceArray = make([]interface{}, len(systems), len(systems))

	for i := range systems {
		interfaceArray[i] = systems[i]
	}

	return
}

// Executes a search by hostname via the JumpCloud API
func (jc JCAPI) GetSystemByHostName(hostname string, withTags bool) ([]JCSystem, JCError) {
	var returnVal []JCSystem

	buffer, err := jc.DoBytes(MapJCOpToHTTP(Insert), "/search"+SYSTEMS_PATH, jc.hostnameFilter(hostname))

	if err != nil {
		return nil, fmt.Errorf("ERROR: Post to JumpCloud failed, err='%s'", err)
	}

	systemResults := JCSystemResults{}

	err = json.Unmarshal(buffer, &systemResults)
	if err != nil {
		return nil, fmt.Errorf("ERROR: Could not unmarshal buffer '%s', err='%s'", buffer, err.Error())
	}

	returnVal = systemResults.Results

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
	url := fmt.Sprintf("%s/%s", SYSTEMS_PATH, systemId)

	buffer, err := jc.DoBytes(MapJCOpToHTTP(Read), url, []byte{})
	if err != nil {
		return system, fmt.Errorf("ERROR: Could not get system by ID '%s', err='%s'", systemId, err.Error())
	}

	err = json.Unmarshal(buffer, &system)
	if err != nil {
		return system, fmt.Errorf("ERROR: Could not unmarshal buffer '%s', err='%s'", buffer, err.Error())
	}

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

	return
}

func (jc JCAPI) GetSystems(withTags bool) (systems []JCSystem, err JCError) {
	for skip := 0; skip == 0 || len(systems) == searchLimit; skip += searchSkipInterval {
		url := fmt.Sprintf("%s?sort=hostname&skip=%d&limit=%d", SYSTEMS_PATH, skip, searchLimit)

		buffer, err2 := jc.DoBytes(MapJCOpToHTTP(Read), url, []byte{})
		if err2 != nil {
			return nil, fmt.Errorf("ERROR: Get to JumpCloud failed, err='%s'", err2)
		}

		systemResults := JCSystemResults{}

		err = json.Unmarshal(buffer, &systemResults)
		if err != nil {
			return nil, fmt.Errorf("ERROR: Could not unmarshal buffer '%s', err='%s'", buffer, err.Error())
		}

		for _, system := range systemResults.Results {
			systems = append(systems, system)
		}

	}

	if withTags {
		tags, err := jc.GetAllTags()
		if err != nil {
			return nil, fmt.Errorf("ERROR: Could not get tags, err='%s'", err)
		}

		for idx, _ := range systems {
			systems[idx].AddJCTagsToSystem(tags)
		}
	}

	return
}

//
// Update a system
//
func (jc JCAPI) UpdateSystem(system JCSystem) (systemId string, err JCError) {
	data, err := json.Marshal(system)
	if err != nil {
		return "", fmt.Errorf("ERROR: Could not marshal JCSystem object, err='%s'", err)
	}

	buffer, err := jc.DoBytes(MapJCOpToHTTP(Update), SYSTEMS_PATH+"/"+system.Id, data)
	if err != nil {
		return "", fmt.Errorf("ERROR: Could not update JCSystem object, err='%s'", err)
	}

	var returnSystem JCSystem

	err = json.Unmarshal(buffer, &returnSystem)
	if err != nil {
		return "", fmt.Errorf("ERROR: Could not unmarshal result buffer '%s', err='%s'", buffer, err.Error())
	}

	if returnSystem.Id != system.Id {
		return "", fmt.Errorf("ERROR: JumpCloud did not return the same ID ('%s', should be '%s'), return buffer='%s'",
			returnSystem.Id, system.Id, buffer)
	}

	systemId = returnSystem.Id

	return systemId, nil
}

//
//               !!!!!!!!!!!!WARNING!!!!!!!!!!!!
//
// This will cause JumpCloud to uninstall the agent on this system.
//    You will lose control of the system after the call returns.
//
func (jc JCAPI) DeleteSystem(system JCSystem) JCError {
	_, err := jc.Delete(fmt.Sprintf("%s/%s", SYSTEMS_PATH, system.Id))
	if err != nil {
		return fmt.Errorf("ERROR: Could not delete system '%s': err='%s'", system.Hostname, err)
	}

	return nil
}

// GetSystemUserBindingsById returns all the user bindings for the given system Id:
// this includes the direct system-user bindings as well as the bindings made via tags
func (jc JCAPI) GetSystemUserBindingsById(systemId string) (systemUserBindings []SystemUserBinding, err JCError) {
	url := fmt.Sprintf("%s/%s/systemusers", SYSTEMS_PATH, systemId)

	buffer, err := jc.DoBytes(MapJCOpToHTTP(Read), url, []byte{})
	if err != nil {
		return systemUserBindings, fmt.Errorf("ERROR: Could not get system user bindings for system ID '%s', err='%s'", systemId, err.Error())
	}

	// The response from the /systems/<system_id>/users endpoint is a map of system-user bindings
	// keyed on the user ID, so we can't unmarshall it in one operation:
	// first parse the response into a map of generic json objects:
	var userBindingMap map[string]*json.RawMessage
	err = json.Unmarshal(buffer, &userBindingMap)

	if err != nil {
		return systemUserBindings, fmt.Errorf("ERROR: Could not unmarshal buffer '%s', err='%s'", buffer, err.Error())
	}

	// iterate over the map of user bindings:
	for userId, _ := range userBindingMap {
		var binding SystemUserBinding
		// unmarshall the current raw json into a system user binding:
		err = json.Unmarshal(*userBindingMap[userId], &binding)
		if err != nil {
			return systemUserBindings, fmt.Errorf("ERROR: Could not unmarshal buffer '%s', err='%s'", *userBindingMap[userId], err.Error())
		}
		// set the user id (obtained from the current key)
		// since it isn't populated in the json response:
		binding.UserId = userId
		// append the current user binding to the array:
		systemUserBindings = append(systemUserBindings, binding)
	}

	return
}
