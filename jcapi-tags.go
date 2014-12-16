package jcapi

import (
	"fmt"
	"strings"
)

type JCTag struct {
	Id                 string   `json:"_id,omitempty"`
	Name               string   `json:"name"`
	GroupName          string   `json:"groupname"`
	Systems            []string `json:"systems"`
	SystemUsers        []string `json:"systemusers"`
	RegularExpressions []string `json:"regularExpressions"`
	ExpirationTime     string   `json:"expirationTime"`
	Expired            bool     `json:"expired"`
	Selected           bool     `json:"selected"`

	//
	// For identification as an external user directory source
	//
	ExternallyManaged  bool   `json:"externallyManaged"`
	ExternalDN         string `json:"externalDN,omitempty"`
	ExternalSourceType string `json:"externalSourceType,omitempty"`

	applyToJumpCloud bool
}

func (tag JCTag) ToString() string {
	return fmt.Sprintf("tag id=%s - name='%s' - groupName='%s' - expires='%s' - systems='%s' - systemusers='%s' - applyToJC='%t' - externally_managed='%t' (%s)",
		tag.Id, tag.Name, tag.GroupName, tag.ExpirationTime, strings.Join(tag.Systems, ","),
		strings.Join(tag.SystemUsers, ","), tag.applyToJumpCloud, tag.ExternallyManaged, tag.ExternalDN)
}

func GetTagNames(tags []JCTag) []string {
	var returnVal []string

	for _, tag := range tags {
		returnVal = append(returnVal, tag.Name)
	}

	return returnVal
}

func (tag *JCTag) MarshalJSON() ([]byte, error) {

	var builder []string

	if tag.Id != "" {
		builder = append(builder, buildJSONKeyValuePair("_id", tag.Id))
	}

	builder = append(builder, buildJSONKeyValuePair("name", tag.Name))
	builder = append(builder, buildJSONKeyValuePair("groupname", tag.GroupName))
	builder = append(builder, buildJSONStringArray("systems", tag.Systems))
	builder = append(builder, buildJSONStringArray("systemusers", tag.SystemUsers))
	builder = append(builder, buildJSONStringArray("regularExpressions", tag.RegularExpressions))
	builder = append(builder, buildJSONKeyValuePair("expirationTime", tag.ExpirationTime))
	builder = append(builder, buildJSONKeyValueBoolPair("expired", tag.Expired))
	builder = append(builder, buildJSONKeyValueBoolPair("selected", tag.Selected))

	builder = append(builder, buildJSONKeyValueBoolPair("externallyManaged", tag.ExternallyManaged))
	builder = append(builder, buildJSONKeyValuePair("externalDN", tag.ExternalDN))
	builder = append(builder, buildJSONKeyValuePair("externalSourceType", tag.ExternalSourceType))

	return []byte("{" + strings.Join(builder, ",") + "}"), nil
}

func (jc JCAPI) getTagFieldsFromInterface(tagData map[string]interface{}, tag *JCTag) {
	tag.Id = tagData["_id"].(string)
	tag.Name = tagData["name"].(string)

	if tagData["groupName"] != nil {
		tag.GroupName = getStringOrNil(tagData["groupName"].(interface{}))
	}

	if tagData["expirationTime"] != nil {
		tag.ExpirationTime = getStringOrNil(tagData["expirationTime"].(interface{}))
	}

	tag.Systems = jc.extractStringArray(tagData["systems"].([]interface{}))
	tag.SystemUsers = jc.extractStringArray(tagData["systemusers"].([]interface{}))
	tag.RegularExpressions = jc.extractStringArray(tagData["regularExpressions"].([]interface{}))

	tag.ExternallyManaged = tagData["externallyManaged"].(bool)

	if tagData["externalDN"] != nil {
		tag.ExternalDN = getStringOrNil(tagData["externalDN"].(interface{}))
	}

	if tagData["externalSourceType"] != nil {
		tag.ExternalSourceType = getStringOrNil(tagData["externalSourceType"].(interface{}))
	}
}

func (jc JCAPI) GetAllTags() (tagList []JCTag, err JCError) {
	var returnVal []JCTag

	for skip := 0; skip == 0 || len(returnVal) == searchLimit; skip += searchSkipInterval {
		url := fmt.Sprintf("/tags?sort=username&skip=%d&limit=%d", skip, searchLimit)

		result, err := jc.Get(url)
		if err != nil {
			return nil, fmt.Errorf("ERROR: Get tags from JumpCloud failed, err='%s'", err)
		}

		recMap := result.(map[string]interface{})

		resultsMap := recMap["results"].([]interface{})

		returnVal := make([]JCTag, len(resultsMap))

		for idx, tagData := range resultsMap {
			jc.getTagFieldsFromInterface(tagData.(map[string]interface{}), &returnVal[idx])
		}

		for i, _ := range returnVal {
			tagList = append(tagList, returnVal[i])
		}
	}

	return
}

//
// Add or Update a tag in place on JumpCloud
//
func (jc JCAPI) AddUpdateTag(op JCOp, tag JCTag) (tagId string, err JCError) {
	data, err := tag.MarshalJSON()
	if err != nil {
		return "", fmt.Errorf("ERROR: Could not marshal JCTag object, err='%s'", err)
	}

	url := "/tags"
	if op == Update {
		url += "/" + tag.Id
	}

	jcTagRec, err := jc.Do(MapJCOpToHTTP(op), url, data)
	if err != nil {
		return "", fmt.Errorf("ERROR: Could not post new JCTag object, err='%s'", err)
	}

	var resultTag JCTag
	jc.getTagFieldsFromInterface(jcTagRec.(map[string]interface{}), &resultTag)

	if resultTag.Name != tag.Name {
		return "", fmt.Errorf("ERROR: JumpCloud did not return the same tag name - this should never happen!")
	}

	tagId = resultTag.Id

	return
}

func (jc JCAPI) DeleteTag(tag JCTag) JCError {
	_, err := jc.Delete(fmt.Sprintf("/tags/%s", tag.Id))
	if err != nil {
		return fmt.Errorf("ERROR: Could not delete tag ID '%s': err='%s'", tag.Id, err)
	}

	return nil
}
