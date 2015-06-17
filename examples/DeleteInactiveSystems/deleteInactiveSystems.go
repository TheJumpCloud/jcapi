package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/TheJumpCloud/jcapi"
)

const (
	URL_BASE string = "https://console.jumpcloud.com/api"
)

func dateBeforeNDays(date string, days int) (before bool, err error) {
	dateField, err := time.Parse(time.RFC3339, date)
	if err != nil {
		err = fmt.Errorf("Could not parse date value '%s', err='%s'", date, err.Error())
		return
	}

	before = dateField.Before(time.Now().Add(time.Duration(-days) * 24 * time.Hour))

	return
}

func main() {
	apiKey := flag.String("api-key", "", "Your JumpCloud Administrator API Key")
	daysSinceLastConnection := flag.Int("days-since-last-connect", 30,
		"Systems that have not connected in this many days or more, will be deleted from JumpCloud.")

	flag.Parse()

	if apiKey != nil && *apiKey == "" {
		log.Fatalf("%s: You must specify an API key value (--api-key=keyValue)", os.Args[0])
	}

	jc := jcapi.NewJCAPI(*apiKey, URL_BASE)

	// Get all the systems in the account
	systems, err := jc.GetSystems(false)
	if err != nil {
		log.Fatalf("Could not get all systems in the account, err='%s'", err.Error())
	}

	for _, system := range systems {
		if system.Active == false {
			okToDelete, err := dateBeforeNDays(system.LastContact, *daysSinceLastConnection)
			if err != nil {
				log.Fatalf("Could not compare date '%s' for system ID '%s' (%s), err='%s'", system.LastContact, system.Id, system.Hostname, err.Error())
			}

			if okToDelete {
				fmt.Printf("Deleting [%s] - ", system.ToString())

				err = jc.DeleteSystem(system)
				if err != nil {
					log.Fatalf("Delete failed, err='%s'\n", err.Error())
				}

				fmt.Printf("SUCCESS!\n")
			}
		}
	}
}
