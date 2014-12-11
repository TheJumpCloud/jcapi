package main

import (
	"fmt"
	"strings"
)

type JCIDSource struct {
	Id             string `json:"_id,omitempty"`
	Name           string `json:"name"`
	Organization   string `json:"organization,omitempty"`
	Type           string `json:type`
	Version        string `json:version`
	IpAddress      string `json:ipAddress`
	LastUpdateTime string `json:lastUpdateTime,omitempty`
	DN             string `json:dn`
	Active         bool   `json:active, omitempty`
}

func (e JCIDSource) toString() string {
	return fmt.Sprintf("idsource: id='%s' - name='%s' - type='%s' - version='%s' - ipAddr='%s' - lastUpdate='%s' - DN='%s' - active='%t'\n",
		e.Id, e.Name, e.Type, e.Version, e.IpAddress, e.LastUpdateTime, e.DN, e.Active)
}

func getIDSourceFieldsFromInterface(fields map[string]interface{}, e *JCIDSource) {
	e.Id = fields["_id"].(string)

	e.Name = fields["name"].(string)

	if _, exists := fields["organization"]; exists {
		e.Organization = fields["organization"].(string)
	}
	if _, exists := fields["type"]; exists {
		e.Type = fields["type"].(string)
	}
	if _, exists := fields["version"]; exists {
		e.Version = fields["version"].(string)
	}
	if _, exists := fields["ipAddress"]; exists {
		e.IpAddress = fields["ipAddress"].(string)
	}
	if _, exists := fields["lastUpdateTime"]; exists {
		e.LastUpdateTime = fields["lastUpdateTime"].(string)
	}
	if _, exists := fields["dn"]; exists {
		e.DN = fields["dn"].(string)
	}
	if _, exists := fields["active"]; exists {
		e.Active = fields["active"].(bool)
	}
}

func (e JCIDSource) MarshalJSON(writeActive bool) ([]byte, error) {

	var builder []string

	if e.Id != "" {
		builder = append(builder, buildJSONKeyValuePair("_id", e.Id))
	}

	builder = append(builder, buildJSONKeyValuePair("name", e.Name))
	builder = append(builder, buildJSONKeyValuePair("organization", e.Organization))
	builder = append(builder, buildJSONKeyValuePair("type", e.Type))
	builder = append(builder, buildJSONKeyValuePair("version", e.Version))
	builder = append(builder, buildJSONKeyValuePair("ipAddress", e.IpAddress))
	builder = append(builder, buildJSONKeyValuePair("lastUpdateTime", e.LastUpdateTime))
	builder = append(builder, buildJSONKeyValuePair("dn", e.DN))

	//
	// We never write 'active' out on a PUT, to prevent a race condition around
	// where the the user may change the setting between when we read the
	// object and write it back in to update the lastUpdateTime.
	//
	if writeActive {
		builder = append(builder, buildJSONKeyValueBoolPair("active", e.Active))
	}

	return []byte("{" + strings.Join(builder, ",") + "}"), nil
}

func getJCIDSourcesFromInterface(idSource interface{}) []JCIDSource {

	var returnVal []JCIDSource

	recMap := idSource.(map[string]interface{})

	if WillDebug(3) {
		dbg(3, "recMap[\"results\"]=%U\n\n------\n", recMap["results"])

		for key, value := range recMap {
			dbg(3, "recMap: key=[%s] - value=[%s]\n", key, value)
		}
	}

	results := recMap["results"].([]interface{})

	returnVal = make([]JCIDSource, len(results))

	for idx, result := range results {
		if WillDebug(3) {
			for key, value := range result.(map[string]interface{}) {
				dbg(3, "results: key=[%s] - value=[%s]\n", key, value)
			}
		}

		getIDSourceFieldsFromInterface(result.(map[string]interface{}), &returnVal[idx])
	}

	return returnVal
}

func (jc JCAPI) getAllIDSources() ([]JCIDSource, JCError) {
	var returnValue []JCIDSource

	result, err := jc.get("/idsources")
	if err != nil {
		return returnValue, fmt.Errorf("ERROR: Could not list ID sources, err='%s'", err)
	}

	returnValue = getJCIDSourcesFromInterface(result)

	return returnValue, nil
}

func (jc JCAPI) getIDSourceByName(name string) (JCIDSource, bool, JCError) {
	var returnValue JCIDSource

	e, err := jc.getAllIDSources()
	if err != nil {
		return returnValue, false, fmt.Errorf("ERROR: Could not gather all ID source objects, err='%s'", err)
	}

	for _, returnValue = range e {
		if returnValue.Name == name {
			return returnValue, true, nil
		}
	}

	return returnValue, false, nil
}

//
// Add or Update an ID source in place on JumpCloud
//
func (jc JCAPI) addUpdateIDSource(op JCOp, idSource JCIDSource) (string, JCError) {
	data, err := idSource.MarshalJSON(op == insert)
	if err != nil {
		return "", fmt.Errorf("ERROR: Could not marshal JCTag object, err='%s'", err)
	}

	url := "/idsources"
	if op == update {
		url += "/" + idSource.Id
	}

	idSourceRec, err := jc.do(mapJCOpToHTTP(op), url, data)
	if err != nil {
		return "", fmt.Errorf("ERROR: Could not post new JCIDSource object, err='%s'", err)
	}

	var resultES JCIDSource
	getIDSourceFieldsFromInterface(idSourceRec.(map[string]interface{}), &resultES)

	if resultES.Name != idSource.Name {
		return "", fmt.Errorf("ERROR: JumpCloud did not return the same ID source name - this should never happen!")
	}

	return resultES.Id, nil
}

func (jc JCAPI) deleteIDSource(idSource JCIDSource) JCError {
	dbg(3, "Deleting External Source ID=%d", idSource.Id)

	_, err := jc.delete(fmt.Sprintf("/idsources/%s", idSource.Id))
	if err != nil {
		return fmt.Errorf("ERROR: Could not delete ID source ID '%s': err='%s'", idSource.Id, err)
	}

	return nil
}
