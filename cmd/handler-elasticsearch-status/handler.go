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
	"github.com/yieldbot/sensu-yieldbot-handler-slack/Godeps/_workspace/src/github.com/olivere/elastic"
	"github.com/yieldbot/sensu-yieldbot-handler-slack/Godeps/_workspace/src/github.com/yieldbot/ybsensu/handler"
	"time"
)

func main() {

	// set commandline flags
	esIndexPtr := flag.String("index", handler.StatusEsIndex, "the elasticsearch index to use")
	esHostPtr := flag.String("host", handler.DefaultEsHost, "the elasticsearch host")
	esPortPtr := flag.String("port", handler.DefaultEsPort, "the elasticsearch port")

	flag.Parse()
	esIndex := *esIndexPtr
	esType := handler.DefaultEsType
	esHost := *esHostPtr
	esPort := *esPortPtr

	sensuEvent := new(handler.SensuEvent)

	sensuEnv := handler.SetSensuEnv()
	sensuEvent = sensuEvent.AcquireSensuEvent()

	// Create a client
	client, err := elastic.NewClient(
		elastic.SetURL("http://" + esHost + ":" + esPort),
	)
	if err != nil {
		handler.Check(err)
	}

	// Check to see if the index exists and if not create it
	if client.IndexExists(esIndex) == nil { // need to test to make sure this does what I want
		_, err = client.CreateIndex(esIndex).Do()
		if err != nil {
			handler.Check(err)
		}
	}

	// Create an Elasticsearch document. The document type will define the mapping used for the document.
	doc := make(map[string]interface{})
	var docID string
	docID = handler.EventName(sensuEvent.Client.Name, sensuEvent.Check.Name)
	doc["monitored_instance"] = sensuEvent.AcquireMonitoredInstance()
	doc["sensu_client"] = sensuEvent.Client.Name
	doc["incident_timestamp"] = time.Unix(sensuEvent.Check.Issued, 0).Format(time.RFC3339)
	doc["check_name"] = handler.CreateCheckName(sensuEvent.Check.Name)
	doc["check_state"] = handler.DefineStatus(sensuEvent.Check.Status)
	doc["sensuEnv"] = handler.DefineSensuEnv(sensuEnv.Sensu.Environment)
	doc["tags"] = sensuEvent.Check.Tags
	doc["instance_address"] = sensuEvent.Client.Address
	doc["check_state_duration"] = handler.DefineCheckStateDuration()

	// Add a document to the Elasticsearch index
	_, err = client.Index().
		Index(esIndex).
		Type(esType).
		Id(docID).
		BodyJson(doc).
		Do()
	if err != nil {
		handler.Check(err)
	}

	// Log a successful document push to stdout. I don't add the id here as some id's are fixed but
	// the user has the ability to autogenerate an id if they don't want to provide one.
	fmt.Printf("Record added to ES\n")
}
