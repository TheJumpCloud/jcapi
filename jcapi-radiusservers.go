package jcapi

import (
	"encoding/json"
	"fmt"
)

type JCRadiusServer struct {
	Id              string   `json:"_id,omitempty"`
	Name            string   `json:"name,omitempty"`
	NetworkSourceIP string   `json:"networkSourceIp,omitempty"`
	SharedSecret    string   `json:"sharedSecret,omitempty"`
	TagList         []string `json:"tags,omitempty"`
}

func (e JCRadiusServer) ToString() string {
	return fmt.Sprintf("radiusserver: id='%s' - name='%s' - IP='%s' - Secret='%s' - TagList=%s",
		e.Id, e.Name, e.NetworkSourceIP, e.SharedSecret, e.TagList)
}

func getJCRadiusServerFieldsFromInterface(fields map[string]interface{}, e *JCRadiusServer) {
	e.Id = fields["_id"].(string)

	e.Name = fields["name"].(string)

	if _, exists := fields["networkSourceIp"]; exists {
		e.NetworkSourceIP = fields["networkSourceIp"].(string)
	}
	if _, exists := fields["sharedSecret"]; exists {
		e.SharedSecret = fields["sharedSecret"].(string)
	}
	if _, exists := fields["tags"]; exists {
		for i, _ := range fields["tags"].([]interface{}) {
			e.TagList = append(e.TagList, (fields["tags"].([]interface{}))[i].(string))
		}
	}
}

func getJCRadiusServerFromInterface(radiusServer interface{}) []JCRadiusServer {

	var returnVal []JCRadiusServer

	recMap := radiusServer.(map[string]interface{})

	results := recMap["results"].([]interface{})

	returnVal = make([]JCRadiusServer, len(results))

	for idx, result := range results {
		getJCRadiusServerFieldsFromInterface(result.(map[string]interface{}), &returnVal[idx])
	}

	return returnVal
}

func (jc JCAPI) GetAllRadiusServers() (radiusServers []JCRadiusServer, err JCError) {

	result, err := jc.Get("/radiusservers")
	if err != nil {
		err = fmt.Errorf("ERROR: Could not list RADIUS servers, err='%s'", err)
		return
	}

	radiusServers = getJCRadiusServerFromInterface(result)

	return
}

func FindRadiusServerById(radiusServers []JCRadiusServer, id string) (radiusServer *JCRadiusServer) {
	index := FindObject(GetInterfaceArrayFromJCRadiusServer(radiusServers), "Id", id)
	if index >= 0 {
		radiusServer = &radiusServers[index]
	}

	return
}

func GetInterfaceArrayFromJCRadiusServer(radiusServers []JCRadiusServer) (interfaceArray []interface{}) {
	interfaceArray = make([]interface{}, len(radiusServers), len(radiusServers))

	for i := range radiusServers {
		interfaceArray[i] = radiusServers[i]
	}

	return
}

//
// Add or Update a radiusserver in place on JumpCloud
//
func (jc JCAPI) AddUpdateRadiusServer(op JCOp, radiusServer JCRadiusServer) (id string, err JCError) {
	data, err := json.Marshal(radiusServer)
	if err != nil {
		return "", fmt.Errorf("ERROR: Could not marshal JCRadius object, err='%s'", err)
	}

	url := "/radiusservers"
	if op == Update {
		url += "/" + radiusServer.Id
	}

	radiusServerRec, err := jc.Do(MapJCOpToHTTP(op), url, data)
	if err != nil {
		return "", fmt.Errorf("ERROR: Could not post new JCIDSource object, err='%s'", err)
	}

	var resultES JCRadiusServer
	getJCRadiusServerFieldsFromInterface(radiusServerRec.(map[string]interface{}), &resultES)

	if resultES.Name != radiusServer.Name {
		return "", fmt.Errorf("ERROR: JumpCloud did not return the same ID source name - this should never happen!")
	}

	return resultES.Id, nil
}

func (jc JCAPI) DeleteRadiusServer(radiusServer JCRadiusServer) JCError {
	_, err := jc.Delete(fmt.Sprintf("/radiusservers/%s", radiusServer.Id))
	if err != nil {
		return fmt.Errorf("ERROR: Could not delete ID source ID '%s': err='%s'", radiusServer.Id, err)
	}

	return nil
}
