package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/TheJumpCloud/jcapi"
)

// DefaultURLBase is the production api endpoint.
const DefaultURLBase string = "https://console.jumpcloud.com/api"

type userAttribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type userAttributes struct {
	Attributes []userAttribute `json:"attributes"`
}

func getUserByEmail(jc jcapi.JCAPI, email string) ([]jcapi.JCUser, error) {
	jcUsers, jcErr := jc.GetSystemUserByEmail(email, false)

	if jcErr != nil {
		err := fmt.Errorf("Error retrieving user for email %s: %s", email, jcErr.Error())
		return jcUsers, err
	}
	return jcUsers, nil
}

func buildAttributes(userRecord []string, attributeNames []string) userAttributes {

	recordLen := len(userRecord)
	attributeArray := make([]userAttribute, len(attributeNames))

	for i, attributeName := range attributeNames {
		attributeArray[i] = userAttribute{Name: attributeName}
		// acount for empty attributes at end of record
		if recordLen > (i + 1) {
			attributeArray[i].Value = userRecord[i+1]
		}
	}

	return userAttributes{attributeArray}
}

func importUserAttributes(jc jcapi.JCAPI, user jcapi.JCUser, attributes userAttributes) error {

	b, err := json.Marshal(attributes)
	if err != nil {
		return fmt.Errorf("Error converting attributes to JSON: %s", err.Error())
	}

	url := "/systemusers/" + user.Id
	_, jcErr := jc.Put(url, b)
	if jcErr != nil {
		return fmt.Errorf("Error setting attribute(s) on user %s: %s", user.Email, jcErr.Error())
	}
	return nil

}

func main() {

	// input parameters
	apiKey := flag.String("api-key", "", "Your JumpCloud Administrator API Key")
	inputFilePath := flag.String("inputFile", "", "CSV file containing user identifier and attributes")
	baseURL := flag.String("url", DefaultURLBase, "Base API Url override")

	flag.Parse()

	if *apiKey == "" || *inputFilePath == "" {
		flag.Usage()
		return
	}

	// Attach to JumpCloud API
	jc := jcapi.NewJCAPI(*apiKey, *baseURL)

	// Setup access to input/output files
	inputFile, err := os.Open(*inputFilePath)
	if err != nil {
		fmt.Printf("Error opening input file %s: %s\n", *inputFilePath, err)
		return
	}
	defer inputFile.Close()

	// Read input file and process users one at a time
	reader := csv.NewReader(inputFile)
	reader.FieldsPerRecord = -1 // indicates records have optional fields

	// Read header row (userIdentifier, attributeName(s)...)
	headerRecord, err := reader.Read()
	if err != nil {
		fmt.Printf("Error reading header row: %s\n", err)
		return
	}

	if len(headerRecord) < 2 {
		fmt.Printf("Invalid header row: File must contain at least 2 columns\n")
		return
	}
	attributeNames := headerRecord[1:]
	// TODO: validate attribute names?

	// Read each row (user identifier + attribute values)
	userCount := 0
	importedUserCount := 0
	var unknownUsers []string
	var errorsByUser = make(map[string]error)

	for {

		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		userCount++
		if err != nil {
			fmt.Printf("Error reading record on line %d: %s\n", userCount+1, err)
			return
		}

		// Fetch user by email
		email := record[0]
		users, err := getUserByEmail(jc, email)
		if err != nil {
			errorsByUser[email] = err
			continue
		}

		if len(users) == 0 {
			unknownUsers = append(unknownUsers, email)
			continue
		}

		user := users[0]
		attributes := buildAttributes(record, attributeNames)
		err = importUserAttributes(jc, user, attributes)
		if err != nil {
			errorsByUser[email] = err
		} else {
			importedUserCount++
		}

	}

	fmt.Println("\nImport complete:")
	fmt.Printf("  %d users processed\n", userCount)
	fmt.Printf("  %d users imported\n", importedUserCount)
	fmt.Printf("  %d users not found\n", len(unknownUsers))
	fmt.Printf("  %d errors processing users\n", len(errorsByUser))

	if len(unknownUsers) > 0 {
		fmt.Println("\nUnknown Users:")
		for _, userEmail := range unknownUsers {
			fmt.Printf("  %s\n", userEmail)
		}
	}

	if len(errorsByUser) > 0 {
		fmt.Println("\nUser Errors:")
		for email, err := range errorsByUser {
			fmt.Printf("  %s: %s\n", email, err.Error())
		}
	}
	fmt.Println("")

	return

}
