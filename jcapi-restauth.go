package jcapi

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
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

func (e JCRestAuth) marshalJSON() (jsonData []byte) {
	var builder []string

	builder = append(builder, buildJSONKeyValuePair("username", e.Username))
	builder = append(builder, buildJSONKeyValuePair("password", e.Password))
	builder = append(builder, buildJSONKeyValuePair("tag", e.Tag))

	jsonData = []byte("{" + strings.Join(builder, ",") + "}")

	return
}

func (jc JCAPI) AuthUser(username, password, tag string) (userAuthenticated bool, err error) {
	userAuthenticated = false

	auth := JCRestAuth{
		Username: username,
		Password: password,
		Tag:      tag,
	}

	data := auth.marshalJSON()

	fullUrl := jc.UrlBase + "/authenticate"

	client := &http.Client{}

	req, err := http.NewRequest("POST", fullUrl, bytes.NewReader(data))
	if err != nil {
		err = fmt.Errorf("ERROR: Could not build POST request: '%s'", err)
		return
	}

	jc.setHeader(req)

	resp, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("ERROR: client.Do() failed, err='%s'", err)
		return
	}

	defer resp.Body.Close()

	if resp.Status == "200 OK" {
		userAuthenticated = true
	}

	return
}
