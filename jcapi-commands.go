package jcapi

import (
	"encoding/json"
	"fmt"
)

const (
	COMMAND_PATH     string = "/commands"
	RUN_COMMAND_PATH string = "/runCommand"
)

type JCCommandResults struct {
	Results []JCCommand `json:"results"`
}

type JCCommand struct {
	Id               string   `json:"_id,omitempty"`            // unique database ID
	Name             string   `json:"name"`                     // a title for display in the UI
	Command          string   `json:"command"`                  // the actual command string to execute
	CommandRunners   []string `json:"commandRunners,omitempty"` // Command Runner user IDs able to run this command
	CommandType      string   `json:"commandType"`              // linux/windows/mac
	User             string   `json:"user,omitempty"`           // user to run as (000000000000000000000000 for root)
	Files            []string `json:"files,omitempty"`          // list of files uploaded by the command
	Systems          []string `json:"systems,omitempty"`        // systems to run the command on
	Tags             []string `json:"tags,omitempty"`           // tags to run the command on (tags and systems are mutually exclusive)
	LaunchType       string   `json:"launchType"`               // manual/add-delete-user/repeated/scheduled
	ListensTo        string   `json:"listensTo"`                // AddUser/DeleteUser (when launchType is add-delete-user)
	Schedule         string   `json:"schedule,omitempty"`       // immediate/agentEvent (launchType=add-delete-user)/a crontab(5) time entry as in "0 0 2 * * 6"
	ScheduledRunDate string   `json:"scheduledRunDate"`         // when LaunchType='scheduled', set to the date on which to start the command
	ScheduledRunTime string   `json:"scheduledRunTime"`         // when LaunchType='scheduled', set to the time at which to start the command
	Trigger          string   `json:"trigger,omitempty"`        // generate trigger (No longer supported)
	Timeout          string   `json:"timeout"`                  // Command time out in seconds, after which it will be killed
	Organization     string   `json:"organization,omitempty"`   // organization ID for this command (auto-populated)
	Sudo             bool     `json:"sudo"`                     // Indicates whether the command should be run with sudo

	Skip  int `json:"skip"`  // Objects to skip on /search POST
	Limit int `json:"limit"` // Max objects to return on /search POST
}

func (e JCCommand) ToString() string {
	return fmt.Sprintf("command: %v", e)
}

// Necessary because json.Unmarshal() can't recognize the JCCommand struct as matching
// the return result.
func getJCCommandFieldsFromInterface(fields map[string]interface{}, e *JCCommand) {
	e.Id = getFieldString("_id", fields)
	e.Name = getFieldString("name", fields)

	e.Command = getFieldString("command", fields)
	e.CommandRunners = getFieldStringArray("commandRunners", fields)
	e.CommandType = getFieldString("commandType", fields)

	e.User = getFieldString("user", fields)

	e.Files = getFieldStringArray("files", fields)
	e.Systems = getFieldStringArray("systems", fields)
	e.Tags = getFieldStringArray("tags", fields)

	e.LaunchType = getFieldString("launchType", fields)
	e.ListensTo = getFieldString("listensTo", fields)
	e.Schedule = getFieldString("schedule", fields)
	e.Trigger = getFieldString("trigger", fields)
	e.Timeout = getFieldString("timeout", fields)
	e.Organization = getFieldString("organization", fields)
	e.Sudo = getFieldBool("sudo", fields)

	e.Skip = getFieldInt("skip", fields)
	e.Limit = getFieldInt("limit", fields)

	return
}

func getJCCommandResultsFromInterface(result interface{}) []JCCommand {

	var returnVal []JCCommand

	recMap := result.(map[string]interface{})

	results := recMap["results"].([]interface{})

	returnVal = make([]JCCommand, len(results), len(results))

	for idx, result := range results {
		getJCCommandFieldsFromInterface(result.(map[string]interface{}), &returnVal[idx])
	}

	return returnVal
}

func (jc JCAPI) GetAllCommands() (commandList []JCCommand, err JCError) {
	for skip := 0; skip == 0 || len(commandList) == searchLimit; skip += searchSkipInterval {
		url := fmt.Sprintf("%s?sort=hostname&skip=%d&limit=%d", COMMAND_PATH, skip, searchLimit)

		jcSysRec, err2 := jc.Get(url)

		if err2 != nil {
			return nil, fmt.Errorf("ERROR: Get commands to JumpCloud failed, err='%s'", err2)
		}

		if jcSysRec == nil {
			return nil, fmt.Errorf("ERROR: No commands found")
		}

		resultsBlock := getJCCommandResultsFromInterface(jcSysRec)

		for i, _ := range resultsBlock {
			if resultsBlock[i].Id != "" {
				commandList = append(commandList, resultsBlock[i])
			}
		}
	}

	return
}

func FindCommandById(commands []JCCommand, id string) (result *JCCommand, index int) {
	index = FindObject(GetInterfaceArrayFromJCCommand(commands), "Id", id)
	if index >= 0 {
		result = &commands[index]
	}

	return
}

func GetInterfaceArrayFromJCCommand(commands []JCCommand) (interfaceArray []interface{}) {
	interfaceArray = make([]interface{}, len(commands), len(commands))

	for i := range commands {
		interfaceArray[i] = commands[i]
	}

	return
}

//
// Add or Update a command in place on JumpCloud
//
func (jc JCAPI) AddUpdateCommand(op JCOp, command JCCommand) (id string, err JCError) {
	id, err = jc.HandleCommand(COMMAND_PATH, op, command)

	return
}

func (jc JCAPI) HandleCommand(path string, op JCOp, command JCCommand) (id string, err JCError) {
	data, err := json.Marshal(command)
	if err != nil {
		return "", fmt.Errorf("ERROR: Could not marshal JCCommand object, err='%s'", err.Error())
	}

	url := path
	if op == Update {
		url += "/" + command.Id
	}

	result, err := jc.Do(MapJCOpToHTTP(op), url, data)
	if err != nil {
		return "", fmt.Errorf("ERROR: Could not '%s' new JCCommand object, err='%s'", MapJCOpToHTTP(op), err.Error())
	}

	commandResult := JCCommand{}

	getJCCommandFieldsFromInterface(result.(map[string]interface{}), &commandResult)

	return commandResult.Id, nil
}

func (jc JCAPI) DeleteCommand(command JCCommand) JCError {
	_, err := jc.Delete(fmt.Sprintf("/%s/%s", COMMAND_PATH, command.Id))
	if err != nil {
		return fmt.Errorf("ERROR: Could not delete command ID '%s': err='%s'", command.Id, err.Error())
	}

	return nil
}

func (jc JCAPI) RunCommand(command JCCommand) JCError {
	_, err := jc.HandleCommand(RUN_COMMAND_PATH, Insert, command)

	return err
}
