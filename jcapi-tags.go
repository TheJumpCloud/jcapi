package jcapi

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	TAGS_PATH string = "/tags"
)

type JCTagResults struct {
	Results []JCTag `json:"results"`
}

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

	ApplyToJumpCloud bool
}

func (tag JCTag) ToString() string {
	return fmt.Sprintf("tag id=%s - name='%s' - groupName='%s' - expires='%s' - systems='%s' - systemusers='%s' - applyToJC='%t' - externally_managed='%t' (%s)",
		tag.Id, tag.Name, tag.GroupName, tag.ExpirationTime, strings.Join(tag.Systems, ","),
		strings.Join(tag.SystemUsers, ","), tag.ApplyToJumpCloud, tag.ExternallyManaged, tag.ExternalDN)
}

func GetTagNames(tags []JCTag) []string {
	var returnVal []string

	for _, tag := range tags {
		returnVal = append(returnVal, tag.Name)
	}

	return returnVal
}

func getJCTagsFromResults(result []byte) (tags []JCTag, err JCError) {
	tag := JCTag{}

	// Try unmarshalling as a single tag first, and if that works, stop there
	err = json.Unmarshal(result, &tag)
	if err != nil || tag.Id == "" {
		err = nil

		tagResults := JCTagResults{}

		// Nope, must be a results array...
		err = json.Unmarshal(result, &tagResults)
		if err != nil {
			err = fmt.Errorf("Could not unmarshal result '%s', err='%s'", string(result), err.Error())
		}

		tags = tagResults.Results
	} else {
		tags = append(tags, tag)
	}

	return
}

func (jc JCAPI) GetTagsByUrl(urlPath string) (tagList []JCTag, err JCError) {

	result, err := jc.DoBytes(MapJCOpToHTTP(Read), urlPath, []byte{})
	if err != nil {
		return nil, fmt.Errorf("ERROR: Get tags from JumpCloud failed with urlPath='%s', err='%s'", urlPath, err.Error())
	}

	tagList, err = getJCTagsFromResults(result)
	if err != nil {
		return nil, fmt.Errorf("ERROR: Could not get tags from results, err='%s'", err.Error())
	}

	return
}

func (jc JCAPI) GetAllTags() (tagList []JCTag, err JCError) {

	var returnVal []JCTag

	for skip := 0; skip == 0 || len(returnVal) == searchLimit; skip += searchSkipInterval {
		url := fmt.Sprintf("%s?sort=username&skip=%d&limit=%d", TAGS_PATH, skip, searchLimit)

		tags, err := jc.GetTagsByUrl(url)
		if err != nil {
			return nil, fmt.Errorf("ERROR: Could not query tags, err='%s'", err)
		}

		for _, tag := range tags {
			tagList = append(tagList, tag)
		}
	}

	return
}

func (jc JCAPI) GetTagByName(tagName string) (tag JCTag, err JCError) {
	url := fmt.Sprintf("%s/%s", TAGS_PATH, tagName)

	tags, err := jc.GetTagsByUrl(url)
	if err != nil {
		err = fmt.Errorf("ERROR: Could not get tags by name for '%s', url='%s', err='%s'", tagName, url, err.Error())
		return
	}

	if len(tags) > 0 {
		tag = tags[0]
	}

	return
}

//
// Add or Update a tag in place on JumpCloud
//
func (jc JCAPI) AddUpdateTag(op JCOp, tag JCTag) (tagId string, err JCError) {
	data, err := json.Marshal(tag)
	if err != nil {
		return "", fmt.Errorf("ERROR: Could not marshal JCTag object, err='%s'", err)
	}

	url := TAGS_PATH
	if op == Update {
		url += "/" + tag.Id
	}

	result, err := jc.DoBytes(MapJCOpToHTTP(op), url, data)
	if err != nil {
		return "", fmt.Errorf("ERROR: Could not post new JCTag object, err='%s'", err)
	}

	tagList, err := getJCTagsFromResults(result)
	if err != nil {
		return "", fmt.Errorf("ERROR: Could not get tags from results, err='%s'", err.Error())
	}

	var resultTag JCTag

	if len(tagList) > 0 {
		resultTag = tagList[0]
	} else {
		return "", fmt.Errorf("ERROR: Got no result back from JumpCloud, cannot return object ID")
	}

	if resultTag.Name != tag.Name {
		return "", fmt.Errorf("ERROR: JumpCloud did not return the same tag name - this should never happen!")
	}

	tagId = resultTag.Id

	return
}

func (jc JCAPI) DeleteTag(tag JCTag) JCError {
	_, err := jc.Delete(fmt.Sprintf("%s/%s", TAGS_PATH, tag.Id))
	if err != nil {
		return fmt.Errorf("ERROR: Could not delete tag ID '%s': err='%s'", tag.Id, err)
	}

	return nil
}
