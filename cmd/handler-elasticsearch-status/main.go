// Take well-formed json from a sensu check result and create an elasticsearch document to be used to
// generate user specific dashboards or highly contextual alerts.
//
// LICENSE:
//   Copyright 2015 Yieldbot. <devops@yieldbot.com>
//   Released under the MIT License; see LICENSE
//   for details.

package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/olivere/elastic"
	"github.com/yieldbot/sensues/lib"
	"github.com/yieldbot/sensuplugin/sensuhandler"
	"github.com/yieldbot/sensuplugin/sensuutil"
	"os"
	"time"
)

func main() {

	var esIndex string
	esType := lib.DefaultEsType
	var esHost string
	var esPort string
	var debug bool

	app := cli.NewApp()
	app.Name = "handler-elasticsearch-status"
	app.Usage = "Send updated status notifications to elasticsearch"
	app.Action = func(c *cli.Context) {

		sensuEvent := new(sensuhandler.SensuEvent)

		sensuEnv := sensuhandler.SetSensuEnv()
		sensuEvent = sensuEvent.AcquireSensuEvent()

		// Create a client
		client, err := elastic.NewClient(
			elastic.SetURL("http://" + esHost + ":" + esPort),
		)
		fmt.Printf("http://" + esHost + ":" + esPort)
		if err != nil {
			sensuutil.EHndlr(err)
		}

		// Check to see if the index exists and if not create it
		if client.IndexExists(esIndex) == nil { // need to test to make sure this does what I want
			_, err = client.CreateIndex(esIndex).Do()
			if err != nil {
				sensuutil.EHndlr(err)
			}
		}

		// Create an Elasticsearch document. The document type will define the mapping used for the document.
		doc := make(map[string]interface{})
		var docID string
		docID = sensuhandler.EventName(sensuEvent.Client.Name, sensuEvent.Check.Name)
		doc["monitored_instance"] = sensuEvent.AcquireMonitoredInstance()
		doc["sensu_client"] = sensuEvent.Client.Name
		doc["incident_timestamp"] = time.Unix(sensuEvent.Check.Issued, 0).Format(time.RFC3339)
		doc["check_name"] = sensuhandler.CreateCheckName(sensuEvent.Check.Name)
		doc["check_state"] = sensuhandler.DefineStatus(sensuEvent.Check.Status)
		doc["sensuEnv"] = sensuhandler.DefineSensuEnv(sensuEnv.Sensu.Environment)
		doc["tags"] = sensuEvent.Check.Tags
		doc["instance_address"] = sensuEvent.Client.Address
		doc["check_state_duration"] = sensuhandler.DefineCheckStateDuration()

		// Add a document to the Elasticsearch index
		_, err = client.Index().
			Index(esIndex).
			Type(esType).
			Id(docID).
			BodyJson(doc).
			Do()
		if err != nil {
			sensuutil.EHndlr(err)
		}

		// Log a successful document push to stdout. I don't add the id here as some id's are fixed but
		// the user has the ability to autogenerate an id if they don't want to provide one.
		fmt.Printf("Record added to ES\n")
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "index, i",
			Value:       lib.StatusEsIndex,
			Usage:       "the elasticsearch index to use",
			EnvVar:      "ES_STATUS_INDEX",
			Destination: &esIndex,
		},
		cli.StringFlag{
			Name:        "host",
			Value:       lib.DefaultEsHost,
			Usage:       "the elasticsearch host",
			EnvVar:      "ES_HOST",
			Destination: &esHost,
		},
		cli.StringFlag{
			Name:        "port, p",
			Value:       lib.DefaultEsPort,
			Usage:       "the elasticsearch port",
			Destination: &esPort,
		},
		cli.BoolFlag{
			Name:        "debug, d",
			Usage:       "Set this to print debugging information. No notifications will be sent",
			Destination: &debug,
		},
	}
	app.Run(os.Args)
}
