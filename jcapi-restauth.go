package jcapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	AUTHENTICATE_PATH string = "/authenticate"
)

type JCRestAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Tag      string `json:"tag"`
}

func (e JCRestAuth) ToString() string {
	return fmt.Sprintf("jcRestAuth: username='%s' - password='<hidden>' - tag='%s'\n",
		e.Username, e.Tag)
}

func (jc JCAPI) AuthUser(username, password, tag string) (userAuthenticated bool, err error) {
	userAuthenticated = false

	auth := JCRestAuth{
		Username: username,
		Password: password,
		Tag:      tag,
	}

	data, err := json.Marshal(auth)
	if err != nil {
		return false, fmt.Errorf("ERROR: Could not marshal the authentication request, err='%s'", err.Error())
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", jc.UrlBase+AUTHENTICATE_PATH, bytes.NewReader(data))
	if err != nil {
		err = fmt.Errorf("ERROR: Could not build POST request: '%s'", err.Error())
		return
	}

	jc.setHeader(req)

	resp, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("ERROR: client.Do() failed, err='%s'", err.Error())
		return
	}

	defer resp.Body.Close()

	if resp.Status == "200 OK" {
		userAuthenticated = true
	}

	return
}
