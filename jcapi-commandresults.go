package jcapi

import (
	"encoding/json"
	"fmt"
	"net/url"
)

const (
	COMMAND_RESULTS_PATH string = "/commandresults"
)

type JCCommandResultResults struct {
	Results []JCCommandResult `json:"results"`
}

type JCCommandResult struct {
	Id                 string     `json:"_id,omitempty"`                // unique database ID
	Name               string     `json:"name"`                         // a title for display in the UI
	Command            string     `json:"command"`                      // the actual command string to execute
	RequestTime        string     `json:"requestTime,omitempty"`        // The time the command started
	ResponseTime       string     `json:"responseTime,omitempty"`       // The time the command exited
	Organization       string     `json:"organization,omitempty"`       // organization ID for this command (auto-populated)
	Sudo               bool       `json:"sudo"`                         // Indicates whether the command should be run with sudo
	System             string     `json:"system,omitempty"`             // The hostname of the system from which this result came
	WorkflowId         string     `json:"workflowId,omitempty"`         // The ID of the workflow of which this command was a part
	WorkflowInstanceId string     `json:"workflowInstanceId,omitempty"` // The instance ID of the workflow of which this command was a part
	Response           JCResponse `json:"response,omitempty"`           // Response data, including command output, and exit code
	Files              []string   `json:"files,omitempty"`              // Names of files uploaded for this command to use during execution
}

type JCResponse struct {
	Id    string `json:"id,omitempty"`
	Data  JCData `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

type JCData struct {
	Output   string `json:"output,omitempty"`
	ExitCode int    `json:"exitCode"`
}

func (e JCCommandResult) ToString() string {
	return fmt.Sprintf("CommandResult: %v", e)
}

func getJCCommandResultsFromResults(result []byte) (commands []JCCommandResult, err JCError) {
	commandResultResults := JCCommandResultResults{}

	err = json.Unmarshal(result, &commandResultResults)
	if err != nil {
		err = fmt.Errorf("Could not unmarshal result '%s', err='%s'", string(result), err.Error())
		return
	}

	commands = commandResultResults.Results

	return
}

func (jc JCAPI) GetCommandResultDetailsById(id string) (commandResult JCCommandResult, err JCError) {
	buffer, err := jc.DoBytes(MapJCOpToHTTP(Read), COMMAND_RESULTS_PATH+"/"+id, []byte{})
	if err != nil {
		err = fmt.Errorf("Could not get command result details for ID '%s', err='%s'", id, err.Error())
	}

	err = json.Unmarshal(buffer, &commandResult)
	if err != nil {
		err = fmt.Errorf("Could not unmarshal buffer '%s', err='%s'", buffer, err.Error())
	}

	return
}

func (jc JCAPI) GetCommandResultsByName(name string) (commandResultList []JCCommandResult, err JCError) {
	searchString1 := "search[fields][]"
	searchString2 := "=name" // can't escape the = here, or we'll get a failure
	searchString3 := "search[searchTerm]"

	if name == "" {
		return nil, fmt.Errorf("ERROR: Name is a required search field and cannot be \"\"")
	}

	for skip := 0; skip == 0 || len(commandResultList) == searchLimit; skip += searchSkipInterval {
		urlQuery := fmt.Sprintf("%s?skip=%d&limit=%d&sort=-requestTime&%s%s&%s=%s", COMMAND_RESULTS_PATH, skip, searchLimit,
			url.QueryEscape(searchString1), searchString2, url.QueryEscape(searchString3), url.QueryEscape(name))

		buffer, err2 := jc.DoBytes(MapJCOpToHTTP(Read), urlQuery, []byte{})
		if err2 != nil {
			return nil, fmt.Errorf("ERROR: Get CommandResults to JumpCloud failed, err='%s'", err2)
		}

		resultsBlock, err2 := getJCCommandResultsFromResults(buffer)
		if err2 != nil {
			err = fmt.Errorf("Could not get resultsBlock data, err='%s'", err2.Error())
			return
		}

		for _, result := range resultsBlock {
			if result.Id != "" {
				commandResultList = append(commandResultList, result)
			}
		}
	}

	return
}

func (jc JCAPI) GetCommandResultsBySavedCommandID(id string) (commandResults []JCCommandResult, err JCError) {
	url := fmt.Sprintf("%s/%s/results", COMMAND_PATH, id)
	body, err := jc.DoBytes(MapJCOpToHTTP(Read), url, nil)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &commandResults); err != nil {
		return nil, err
	}
	return commandResults, err
}

func FindCommandResultById(commandResults []JCCommandResult, id string) (result *JCCommandResult, index int) {
	index = FindObject(GetInterfaceArrayFromJCCommandResults(commandResults), "Id", id)
	if index >= 0 {
		result = &commandResults[index]
	}

	return
}

func GetInterfaceArrayFromJCCommandResults(commandResults []JCCommandResult) (interfaceArray []interface{}) {
	interfaceArray = make([]interface{}, len(commandResults), len(commandResults))

	for i := range commandResults {
		interfaceArray[i] = commandResults[i]
	}

	return
}

func (jc JCAPI) DeleteCommandResult(id string) (err JCError) {
	url := fmt.Sprintf("%s/%s", COMMAND_RESULTS_PATH, id)

	_, err2 := jc.DoBytes(MapJCOpToHTTP(Delete), url, []byte{})
	if err2 != nil {
		return fmt.Errorf("ERROR: DELETE CommandResults failed, err='%s'", err2)
	}

	return
}
