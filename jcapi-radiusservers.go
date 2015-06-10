package jcapi

import (
	"encoding/json"
	"fmt"
)

const (
	RADIUS_SERVERS_PATH string = "/radiusservers"
)

type JCRadiusServerResults struct {
	Results []JCRadiusServer `json:"results"`
}

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

func (jc JCAPI) GetAllRadiusServers() (radiusServers []JCRadiusServer, err JCError) {
	result, err := jc.DoBytes(MapJCOpToHTTP(Read), RADIUS_SERVERS_PATH, []byte{})
	if err != nil {
		err = fmt.Errorf("ERROR: Could not list RADIUS servers, err='%s'", err)
		return
	}

	radiusResults := JCRadiusServerResults{}

	err = json.Unmarshal(result, &radiusResults)
	if err != nil {
		err = fmt.Errorf("ERROR: Could not unmarshal result buffer '%s', err='%s'", result, err.Error())
		return
	}

	radiusServers = radiusResults.Results
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
	interfaceArray = make([]interface{}, len(radiusServers))

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

	buffer, err := jc.DoBytes(MapJCOpToHTTP(op), url, data)
	if err != nil {
		return "", fmt.Errorf("ERROR: Could not post new JCIDSource object, err='%s'", err)
	}

	var resultES JCRadiusServer

	err = json.Unmarshal(buffer, &resultES)
	if err != nil {
		return "", fmt.Errorf("ERROR: Could not unmarshal buffer '%s', err='%s'", buffer, err.Error())
	}

	if resultES.Name != radiusServer.Name {
		return "", fmt.Errorf("ERROR: JumpCloud did not return the same ID source name - this should never happen!")
	}

	return resultES.Id, nil
}

func (jc JCAPI) DeleteRadiusServer(radiusServer JCRadiusServer) JCError {
	_, err := jc.Delete(fmt.Sprintf("%s/%s", RADIUS_SERVERS_PATH, radiusServer.Id))
	if err != nil {
		return fmt.Errorf("ERROR: Could not delete ID source ID '%s': err='%s'", radiusServer.Id, err)
	}

	return nil
}
