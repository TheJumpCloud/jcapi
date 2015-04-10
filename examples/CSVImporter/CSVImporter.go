package main

import (
    "fmt"
    "flag"
    "encoding/csv"
    "os"
    "io"
    "github.com/TheJumpCloud/jcapi"
)


//
// This program will process each line (record) of a CSV file as a user
// import request into JumpCloud.
//
// The CSV file must have each line formatted as:
// 
// first_name, last_name, USER_NAME, EMAIL, uid, gid, SUDO_FLAG, password, host_name, tag_name, admin1, admin2, ...
//
// Values shown in all lowercase are optional, while those in all uppercase
// are required.
//
// For each line of the CSV file, this program will:
//
// 1 - Insert the USER_NAME specified as a new JumpCloud user, if the name
//     does not already exist, otherwise this line will be treated as an
//     update to an existing user matching USER_NAME.
//
// 2 - If the password was specified, the USER_NAME will have that password
//     assigned to it immediately, otherwise, JumpCloud will automatically
//     generate an email message, delivered to the EMAIL specified, requesting
//     the user to provide an appropriate password for their account.
//
// 3 - The value of the SUDO_FLAG will always be applied to this user.  All
//     other optional values specified will be applied as appropriate.
//
// 4 - When both the host_name and tag_name are specified, a tag will be
//     created or updated for this USER_NAME associated with the host_name
//     provided.
//     (n.b. - if the USER_NAME already exists on the host_name in question,
//     specifying these values here will result in the account on host_name
//     being "taken over" by JumpCloud)
//
// 5 - Any invalid combination of values will generate an error and terminate
//     processing for that entry, but the program will attempt to process the
//     remainder of the file.
//
// 6 - A summary is printed at the conclusion of processing.
//


const (
    urlBase string = "https://console.jumpcloud.com/api"
)


//
// Returns the ID for the username specified if it is contained in the
// list of users provided.  (helper function)
//

func GetUserIdFromUserName(users []jcapi.JCUser, name string) string {
    returnVal := ""

    for _, user := range users {
        if user.UserName == name {
            returnVal = user.Id
            break
        }
    }

    return returnVal
}


//
// Main Entry Point...
//

func main() {
    // Input parameters
    var apiKey string
    var csvFile string

	// Obtain the input parameters
	flag.StringVar(&csvFile, "csv", "", "-csv=<filename>")
	flag.StringVar(&apiKey, "key", "", "-key=<API-key-value>")
	flag.Parse()

    // Attach to JumpCloud
    jc := jcapi.NewJCAPI(apiKey, urlBase)

	// Setup access the CSV file specified
	inFile, err := os.Open(csvFile)

    if err != nil {
        fmt.Println(err)
        return
    }

    defer inFile.Close()

    reader := csv.NewReader(inFile)
    reader.FieldsPerRecord = -1    // indicates records have optional fields

    // Fetch all systems in JumpCloud for lookups below
    systemList, err := jc.GetSystems(true)

    if err != nil {
        fmt.Printf("Could not read systems, err='%s'\n", err)
        return
    }

    // Process each user/request record found in the CSV file...
    recordCount := 0

    for {
        // Setup work variables
        var currentUser jcapi.JCUser
        var currentHost string
        var currentTag  string

        currentAdmins :=  make(map[string]string)  // "user name", "user id"

        var fieldMap = map[int]*string {
            0 : &currentUser.FirstName,
            1 : &currentUser.LastName,
            2 : &currentUser.UserName,
            3 : &currentUser.Email,
            4 : &currentUser.Uid,
            5 : &currentUser.Gid,
            // "Sudo" boolean will be handled separately, so no 6
            7 : &currentUser.Password,
            8 : &currentHost,
            9 : &currentTag,
        }

        // (Re)Fetch all users in JumpCloud (pickup any newly added users)
        userList, err := jc.GetSystemUsers(false)

        if err != nil {
            fmt.Printf("Could not read system users, err='%s'\n", err)
            return
        }

        // Read next record from CSV file
    	record, err := reader.Read()

        // Exit loop at the end of file or on error
    	if err == io.EOF {
    		fmt.Println("EOF")
    		break
    	} else if err != nil {
    		fmt.Println(err)
    		return
    	}

        recordCount = recordCount + 1

    	// Parse the record just read into our work vars
    	for i, element := range record {
            // Handle variable fields separately
            if i > 9 {
                break
            }

            // Special case for sole boolean to be parsed
            if i == 6 {
                currentUser.Sudo = jcapi.GetTrueOrFalse(element)
            } else {
                // Default case is to move the string into the var
                *fieldMap[i] = element
            }
    	}

        // The administrators list is optional, and variable.  Using a slice
        // will pick it up if it exists without causing errors otherwise.
        // Map any names found to their ID values for later use.
        adminsSlice := record[10:]

        for _, tempAdmin := range adminsSlice {
            currentAdmins[tempAdmin] = GetUserIdFromUserName(userList, tempAdmin)
        }

        // Determine operation to perform based on whether the current user
        // is already in JumpCloud...
        var opCode jcapi.JCOp
        currentUserId := GetUserIdFromUserName(userList, currentUser.UserName)

        if currentUserId != "" {
            opCode = jcapi.Update
            currentUser.Id = currentUserId
        } else {
            opCode = jcapi.Insert
        }

        // Perform the requested operation on the current user and report results
        currentUserId, err = jc.AddUpdateUser(opCode, currentUser)

        if err != nil {
            fmt.Printf("Could not process user '%s', err='%s'", currentUser.ToString(), err)
            continue
        } else {
            fmt.Printf("Processed user ID '%s'\n", currentUserId)
        }

        // Create/associate JumpCloud tags for the host and user...
        if currentHost != "" {
            // Determine if the host specified is defined in JumpCloud
            var currentJCSystem jcapi.JCSystem

            for _, testSys := range systemList {
                if testSys.Hostname == currentHost {
                    currentJCSystem = testSys
                    break
                }
            }

            if currentJCSystem.Id != "" {
                // Determine operation to perform based on whether the tag
                // is already in JumpCloud...
                var tempTag jcapi.JCTag

                tempTag.Name = currentHost + " - " + currentUser.FirstName + " " + currentUser.LastName

                hasTag, tagId := currentJCSystem.SystemHasTag(tempTag.Name)

                if hasTag {
                    opCode = jcapi.Update
                    tempTag.Id = tagId
                } else {
                    opCode = jcapi.Insert
                }

                // Build a suitable tag from the request's elements
                tempTag.ApplyToJumpCloud = true
                tempTag.Systems = append(tempTag.Systems, currentJCSystem.Id)
                tempTag.SystemUsers = append(tempTag.SystemUsers, currentUserId)

                for _, adminId := range currentAdmins {
                    tempTag.SystemUsers = append(tempTag.SystemUsers, adminId)
                }

                // Create or modify the tag in JumpCloud
                tempTag.Id, err = jc.AddUpdateTag(opCode, tempTag)

                if err != nil {
                    fmt.Printf("Could not process tag '%s', err='%s'", tempTag.ToString(), err)
                } else {
                    fmt.Printf("Processed tag ID '%s'\n", tempTag.Id)
                }
            }
        }
    }

    // Print run summary
    fmt.Printf("\n\nProcessed %d records from file %s \n", recordCount, csvFile)

	return
}
