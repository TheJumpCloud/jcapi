package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"

	"github.com/TheJumpCloud/jcapi"
)

func main() {
	// Input parameters
	var apiKey string
	var csvFile string

	// Obtain the input parameters
	flag.StringVar(&csvFile, "output", "o", "-output=<filename>")
	flag.StringVar(&apiKey, "key", "k", "-key=<API-key-value>")
	flag.Parse()

	if csvFile == "" || apiKey == "" {
		fmt.Println("Usage of ./CSVImporter:")
		fmt.Println("  -output=\"\": -output=<filename>")
		fmt.Println("  -key=\"\": -key=<API-key-value>")
		return
	}

	// Attach to JumpCloud
	jc := jcapi.NewJCAPI(apiKey, jcapi.StdUrlBase)

	// Fetch all users who's password expires between given dates in
	userList, err := jc.GetSysemUsersByPasswordEpiryDate()

	if err != nil {
		fmt.Printf("Could not read system users, err='%s'\n", err)
		return
	}

	// Setup access the CSV file specified
	w := csv.NewWriter(csvFile)

	if err := w.Write([]string{"FIRSTNAME", "LASTNAME", "EMAIL", "PASSWORD EXPIRY DATE"}); err != nil {
		log.Fatalln("error writing record to csv:", err)
	}

	for _, record := range userList {
		if err := w.Write([]string{record.FirstName, record.LastName, record.Email, record.PasswordExpirationDate.String()}); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}
	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Finished")

	return
}
