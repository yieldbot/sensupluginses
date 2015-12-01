// Take well-formed json from stdin and create an elasticsearch document to be used to
// generate user specific dashboards or highly contextual alerts.
//
// LICENSE:
//   Copyright 2015 Yieldbot. <devops@yieldbot.com>
//   Released under the MIT License; see LICENSE
//   for details.

package main

import (
	"flag"
	"fmt"
	"github.com/olivere/elastic"
	dracky "github.com/yieldbot/sensu-yieldbot-library/src"
	"time"
)

func main() {

	// set commandline flags
	esIndexPtr := flag.String("index", dracky.StatusEsIndex, "the elasticsearch index to use")
	esHostPtr := flag.String("host", dracky.DefaultEsHost, "the elasticsearch host")
	esPortPtr := flag.String("port", dracky.DefaultEsPort, "the elasticsearch port")

	flag.Parse()
	esIndex := *esIndexPtr
	esType := dracky.DefaultEsType
	esHost := *esHostPtr
	esPort := *esPortPtr

	sensuEvent := new(dracky.SensuEvent)

	sensuEnv := dracky.SetSensuEnv()
	sensuEvent = sensuEvent.AcquireSensuEvent()

	// Create a client
	client, err := elastic.NewClient(
		elastic.SetURL("http://" + esHost + ":" + esPort),
	)
	if err != nil {
		dracky.Check(err)
	}

	// Check to see if the index exists and if not create it
	if client.IndexExists(esIndex) == nil { // need to test to make sure this does what I want
		_, err = client.CreateIndex(esIndex).Do()
		if err != nil {
			dracky.Check(err)
		}
	}

	// Create an Elasticsearch document. The document type will define the mapping used for the document.
	doc := make(map[string]interface{})
	var docID string
	docID = dracky.EventName(sensuEvent.Client.Name, sensuEvent.Check.Name)
	doc["monitored_instance"] = sensuEvent.AcquireMonitoredInstance()
	doc["sensu_client"] = sensuEvent.Client.Name
	doc["incident_timestamp"] = time.Unix(sensuEvent.Check.Issued, 0).Format(time.RFC3339)
	doc["check_name"] = dracky.CreateCheckName(sensuEvent.Check.Name)
	doc["check_state"] = dracky.DefineStatus(sensuEvent.Check.Status)
	doc["sensuEnv"] = dracky.DefineSensuEnv(sensuEnv.Sensu.Environment)
	doc["tags"] = sensuEvent.Check.Tags
	doc["instance_address"] = sensuEvent.Client.Address
	doc["check_state_duration"] = dracky.DefineCheckStateDuration()

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
