package jcapi

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	IDSOURCES_PATH string = "/idsources"
)

type JCIDSourceResults struct {
	Results []JCIDSource `json:"results"`
}

type JCIDSource struct {
	Id             string `json:"_id,omitempty"`
	Name           string `json:"name"`
	Organization   string `json:"organization,omitempty"`
	Type           string `json:type`
	Version        string `json:version`
	IpAddress      string `json:ipAddress`
	LastUpdateTime string `json:lastUpdateTime,omitempty`
	DN             string `json:dn`
	Active         bool   `json:active,omitempty`
}

func (e JCIDSource) ToString() string {
	return fmt.Sprintf("idsource: id='%s' - name='%s' - type='%s' - version='%s' - ipAddr='%s' - lastUpdate='%s' - DN='%s' - active='%t'\n",
		e.Id, e.Name, e.Type, e.Version, e.IpAddress, e.LastUpdateTime, e.DN, e.Active)
}

func (e JCIDSource) marshalJSON(writeActive bool) ([]byte, error) {

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

func (jc JCAPI) GetAllIDSources() (idSources []JCIDSource, err JCError) {
	result, err := jc.DoBytes(MapJCOpToHTTP(Read), IDSOURCES_PATH, []byte{})
	if err != nil {
		return idSources, fmt.Errorf("ERROR: Could not list ID sources, err='%s'", err)
	}

	idSourceResults := JCIDSourceResults{}

	err = json.Unmarshal(result, &idSourceResults)
	if err != nil {
		err = fmt.Errorf("Could not unmarshal result set, err='%s'", err.Error())
		return
	}

	idSources = idSourceResults.Results

	return
}

func (jc JCAPI) GetIDSourceByName(name string) (idSource JCIDSource, exists bool, err JCError) {
	e, err := jc.GetAllIDSources()
	if err != nil {
		return idSource, false, fmt.Errorf("ERROR: Could not gather all ID source objects, err='%s'", err)
	}

	for _, idSource = range e {
		if idSource.Name == name {
			exists = true
			return
		}
	}

	return
}

//
// Add or Update an ID source in place on JumpCloud
//
func (jc JCAPI) AddUpdateIDSource(op JCOp, idSource JCIDSource) (string, JCError) {
	data, err := idSource.marshalJSON(op == Insert)
	if err != nil {
		return "", fmt.Errorf("ERROR: Could not marshal JCIDSource object, err='%s'", err)
	}

	url := IDSOURCES_PATH
	if op == Update {
		url += "/" + idSource.Id
	}

	buffer, err := jc.DoBytes(MapJCOpToHTTP(op), url, data)
	if err != nil {
		return "", fmt.Errorf("ERROR: Could not post new JCIDSource object, err='%s'", err)
	}

	var resultES JCIDSource

	err = json.Unmarshal(buffer, &resultES)
	if err != nil {
		return "", fmt.Errorf("ERROR: Could not unmarshal result buffer '%s', err='%s'", buffer, err.Error())
	}

	if resultES.Name != idSource.Name {
		return "", fmt.Errorf("ERROR: JumpCloud did not return the same ID source name - this should never happen!")
	}

	return resultES.Id, nil
}

func (jc JCAPI) DeleteIDSource(idSource JCIDSource) JCError {
	_, err := jc.Delete(fmt.Sprintf("%s/%s", IDSOURCES_PATH, idSource.Id))
	if err != nil {
		return fmt.Errorf("ERROR: Could not delete ID source ID '%s': err='%s'", idSource.Id, err)
	}

	return nil
}
