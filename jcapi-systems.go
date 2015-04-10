package jcapi

import (
	"fmt"
)

type JCSystem struct {
	Id               string `json:"_id,omitempty"`
	Active           bool   `json:"active"`
	HostName         string `json:"hostname,omitempty"`
	DisplayName      string `json:"displayName,omitempty"`

	Tags []JCTag
}

func SystemsToString(systems []JCSystem) string {
	returnVal := ""

	for _, sys := range systems {
		returnVal += sys.ToString()
	}

	return returnVal
}

func (jcsystem JCSystem) ToString() string {
	returnVal := fmt.Sprintf("JCSYSTEM: Id=[%s] - HostName=[%s] - DisplayName=[%s] - Active=[%t]\n",
		jcsystem.Id, jcsystem.HostName, jcsystem.DisplayName, jcsystem.Active)

	for _, tag := range jcsystem.Tags {
		returnVal += fmt.Sprintf("\t%s\n", tag.ToString())
	}

	return returnVal
}

func (jcsystem JCSystem) SystemHasTag(tagName string) (bool, string) {
	for _, tag := range jcsystem.Tags {
		if tag.Name == tagName {
			return true, tag.Id
		}
	}

	return false, ""
}

func getJCSystemFieldsFromInterface(fields map[string]interface{}, system *JCSystem) {
	system.Id = fields["_id"].(string)

	if _, exists := fields["displayName"]; exists {
		system.DisplayName = fields["displayName"].(string)
	}

	if _, exists := fields["hostname"]; exists {
		system.HostName = fields["hostname"].(string)
	}

	system.Active = fields["active"].(bool)
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

			system.AddSystemJCTags(tags)
		}
	}

	return
}

func (jc JCAPI) GetSystems(withTags bool) (systemList []JCSystem, err JCError) {
	var returnVal []JCSystem

	for skip := 0; skip == 0 || len(returnVal) == searchLimit; skip += searchSkipInterval {
		url := fmt.Sprintf("/systems?sort=hostname&skip=%d&limit=%d", skip, searchLimit)

		jcSysRec, err2 := jc.Get(url)
		if err != nil {
			return nil, fmt.Errorf("ERROR: Post to JumpCloud failed, err='%s'", err2)
		}

		if jcSysRec == nil {
			return nil, fmt.Errorf("ERROR: No systems found")
		}

		// We really only care about the ID for the following call...
		returnVal = getJCSystemsFromInterface(jcSysRec)

		for i, _ := range returnVal {
			if returnVal[i].Id != "" {

				//
				// Get the rest of the system record
				//
				// We'll get all the tags one time later, so don't get the tags on this call...
				//
				// See above about the compiler error that requires me to use err2 instead of err below...
				//
				detailedSystem, err2 := jc.GetSystemById(returnVal[i].Id, false)
				if err != nil {
					err = fmt.Errorf("ERROR: Could not get details for system ID '%s', err='%s'", returnVal[i].Id, err2)
					return
				}

				if detailedSystem.Id != "" {
					systemList = append(systemList, detailedSystem)
				}
			}
		}
	}

	if withTags {
		tags, err := jc.GetAllTags()
		if err != nil {
			return nil, fmt.Errorf("ERROR: Could not get tags, err='%s'", err)
		}

		for idx, _ := range systemList {
			systemList[idx].AddSystemJCTags(tags)
		}
	}

	return
}
