// Take well-formed json from either stdin or an input file and create an elasticsearch document to be used to
// generate user specific dashboards or highly contextual alerts.
//
// LICENSE:
//   Copyright 2015 Yieldbot. <devops@yieldbot.com>
//   Released under the MIT License; see LICENSE
//   for details.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/olivere/elastic"
	"github.com/yieldbot/dracky/src"
	"io/ioutil"
	"os"
	"time"
)

func main() {

	// set commandline flags
	esIndexPtr := flag.String("index", dracky.StatusEsIndex, "the elasticsearch index to use")
	esHostPtr := flag.String("host", dracky.DefaultEsHost, "the elasticsearch host")
	esPortPtr := flag.String("port", dracky.DefaultEsPort, "the elasticsearch port")
	stdinPtr := flag.Bool("read-stdin", true, "read input from stdin")
	//timePtr := flag.string("t-format", "", "time format to suffix on the index name")
	inputFilePtr := flag.String("input-file", "", "file to read json in from, check docs for proper format")

	flag.Parse()
	esIndex := *esIndexPtr
	esType := dracky.DefaultEsType
	esHost := *esHostPtr
	esPort := *esPortPtr
	rdStdin := *stdinPtr
	inputFile := *inputFilePtr

	// I don't want to call these if they are not needed
	sensuEvent := new(dracky.SensuEvent)
	userEvent := new(dracky.UserEvent)
	//t_format := *timePtr

	sensuEnv := dracky.SetSensuEnv()

	//   if t_format != "" {
	//     // get the format of the time
	//     esIndex = esIndex + t_format
	//   }

	if (rdStdin == false) && (inputFile != "") {
		userInput, err := ioutil.ReadFile(inputFile)
		if err != nil {
			dracky.Check(err)
		}
		err = json.Unmarshal(userInput, &userEvent)
		if err != nil {
			dracky.Check(err)
		}
		esType = "user"
	} else if (rdStdin == false) && (inputFile == "") {
		fmt.Printf("Please enter a file to read from")
		os.Exit(1)
	} else {
		sensuEvent = sensuEvent.AcquireSensuEvent()
	}

	// Create a client
	client, err := elastic.NewClient(
		elastic.SetURL("http://" + esHost + ":" + esPort),
	)
	if err != nil {
		dracky.Check(err)
	}

	// Check to see if the index exists and if not create it
	if client.IndexExists == nil { // need to test to make sure this does what I want
		_, err = client.CreateIndex(esIndex).Do()
		if err != nil {
			dracky.Check(err)
		}
	}

	// Create an Elasticsearch document. The document type will define the mapping used for the document.
	doc := make(map[string]string)
	var docID string
	switch esType {
	case "sensu":
		docID = dracky.EventName(sensuEvent.Client.Name, sensuEvent.Check.Name)
		doc["monitored_instance"] = sensuEvent.AcquireMonitoredInstance()
		doc["sensu_client"] = sensuEvent.Client.Name
		doc["incident_timestamp"] = time.Unix(sensuEvent.Check.Issued, 0).Format(time.RFC3339)
		doc["check_name"] = dracky.CreateCheckName(sensuEvent.Check.Name)
		doc["check_state"] = dracky.DefineStatus(sensuEvent.Check.Status)
		doc["sensuEnv"] = dracky.DefineSensuEnv(sensuEnv.Sensu.Environment)
		doc["instance_address"] = sensuEvent.Client.Address
		doc["check_state_duration"] = dracky.DefineCheckStateDuration()
	case "user":
		doc["product"] = userEvent.Product
		doc["data"] = userEvent.Data
		doc["timestamp"] = time.Unix(sensuEvent.Check.Issued, 0).Format(time.RFC3339) // dracky.Set_time(userEvent.Timestamp)
	default:
		fmt.Printf("Type is not correctly set")
		os.Exit(2)
	}

	// Add a document to the Elasticsearch index
	_, err = client.Index().
		Index(esIndex).
		Type(esType).
		Id(docID).
		BodyJson(doc).
		Do()
	if err != nil {
		dracky.Check(err)
	}

	// Log a successful document push to stdout. I don't add the id here as some id's are fixed but
	// the user has the ability to autogenerate an id if they don't want to provide one.
	fmt.Printf("Record added to ES\n")
}
