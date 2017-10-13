package main

import jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"

// the following constants are used for API v2 calls:
const (
	apiKeyEnvVariable  = "JUMPCLOUD_APIKEY"
	apiKeyHeader       = "x-api-key"
	contentType        = "application/json"
	accept             = "application/json"
	searchLimit        = 100
	searchSkipInterval = 100
)

// isGroupsOrg returns true if this org is groups enabled:
func isGroupsOrg(urlBase string, apiKey string) (bool, error) {
	// instantiate a new API object for User Groups:
	userGroupsAPI := jcapiv2.NewUserGroupsApiWithBasePath(urlBase + "/v2")
	userGroupsAPI.Configuration.APIKey[apiKeyHeader] = apiKey
	// in order to check for groups support, we just query for the list of User groups
	// (we just ask to retrieve 1) and check the response status code:
	_, res, err := userGroupsAPI.GroupsUserList(contentType, accept, "", "", 1, 0, "")

	// check if we're using the API v1:
	// we need to explicitly check for 404, since GroupsUserList will also return a json
	// unmarshalling error (err will not be nil) if we're running this endpoint against
	// a Tags org and we don't want to treat this case as an error:
	if res.Response != nil && res.Response.StatusCode == 404 {
		return false, nil
	}

	// if there was any kind of other error, return that:
	if err != nil {
		return false, err
	}

	// if we're using API v2, we're expecting a 200:
	if res.Response.StatusCode == 200 {
		return true, nil
	}

	return false, nil
}
