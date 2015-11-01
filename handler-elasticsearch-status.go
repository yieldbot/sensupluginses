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
	"github.com/yieldbot/dhuran"
	"github.com/yieldbot/dracky"
	"github.com/olivere/elastic"
	"io/ioutil"
	"os"
	"time"
)

func main() {

	// set commandline flags
	es_indexPtr := flag.String("index", dracky.STATUS_ES_INDEX, "the elasticsearch index to use")
	es_hostPtr := flag.String("host", dracky.DEFAULT_ES_HOST, "the elasticsearch host")
	es_portPtr := flag.String("port", dracky.DEFAULT_ES_PORT, "the elasticsearch port")
	stdinPtr := flag.Bool("read-stdin", true, "read input from stdin")
	//timePtr := flag.string("t-format", "", "time format to suffix on the index name")
	input_filePtr := flag.String("input-file", "", "file to read json in from, check docs for proper format")

	flag.Parse()
	es_index := *es_indexPtr
	es_type := dracky.DEFAULT_ES_TYPE
	es_host := *es_hostPtr
	es_port := *es_portPtr
	rd_stdin := *stdinPtr
	input_file := *input_filePtr

	// I don't want to call these if they are not needed
	sensu_event := new(dracky.Sensu_Event)
	user_event := new(dracky.User_Event)
	//t_format := *timePtr

	sensu_env := dracky.Set_sensu_env()

	//   if t_format != "" {
	//     // get the format of the time
	//     es_index = es_index + t_format
	//   }

	if (rd_stdin == false) && (input_file != "") {
		user_input, err := ioutil.ReadFile(input_file)
		if err != nil {
			dhuran.Check(err)
		}
		err = json.Unmarshal(user_input, &user_event)
		if err != nil {
			dhuran.Check(err)
		}
		es_type = "user"
	} else if (rd_stdin == false) && (input_file == "") {
		fmt.Printf("Please enter a file to read from")
		os.Exit(1)
	} else {
		sensu_event = sensu_event.Acquire_sensu_event()
	}

	// Create a client
	client, err := elastic.NewClient(
		elastic.SetURL("http://" + es_host + ":" + es_port),
	)
	if err != nil {
		dhuran.Check(err)
	}

	// Check to see if the index exists and if not create it
	if client.IndexExists == nil { // need to test to make sure this does what I want
		_, err = client.CreateIndex(es_index).Do()
		if err != nil {
			dhuran.Check(err)
		}
	}

	// Create an Elasticsearch document. The document type will define the mapping used for the document.
	doc := make(map[string]string)
	var doc_id string
	switch es_type {
	case "sensu":
		doc_id = dracky.Event_name(sensu_event.Client.Name, sensu_event.Check.Name)
		doc["monitored_instance"] = sensu_event.Acquire_monitored_instance()
		doc["sensu_client"] = sensu_event.Client.Name
		doc["incident_timestamp"] = time.Unix(sensu_event.Check.Issued, 0).Format(time.RFC822Z)
		doc["check_name"] = dracky.Create_check_name(sensu_event.Check.Name)
		doc["check_state"] = dracky.Define_status(sensu_event.Check.Status)
		doc["sensu_env"] = dracky.Define_sensu_env(sensu_env.Sensu.Environment)
		doc["instance_address"] = sensu_event.Client.Address
		doc["check_state_duration"] = dracky.Define_check_state_duration()
	case "user":
		doc["product"] = user_event.Product
		doc["data"] = user_event.Data
		doc["timestamp"] = time.Unix(sensu_event.Check.Issued, 0).Format(time.RFC822Z) // dracky.Set_time(user_event.Timestamp)
	default:
		fmt.Printf("Type is not correctly set")
		os.Exit(2)
	}

	// Add a document to the Elasticsearch index
	_, err = client.Index().
		Index(es_index).
		Type(es_type).
		Id(doc_id).
		BodyJson(doc).
		Do()
	if err != nil {
		dhuran.Check(err)
	}

	// Log a successful document push to stdout. I don't add the id here as some id's are fixed but
	// the user has the ability to autogenerate an id if they don't want to provide one.
	fmt.Printf("Record added to ES\n")
}
