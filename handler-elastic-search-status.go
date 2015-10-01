// A Sensu handler for dropping the check result into an elasticsearch queue
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"gopkg.in/olivere/elastic.v2"
	"io/ioutil"
	"os"
	"time"
)

// Data structure for holding Sensu generated check results.
type Sensu_Event struct {
	Action      string
	Occurrences int
	Client      struct {
		Name          string
		Address       string
		Subscriptions []string
		Timestamp     int64
	}
	Check struct {
		Source      string
		Name        string
		Issued      int64
		Subscribers []string
		Interval    int
		Command     string
		Output      string
		Status      int
		Handler     string
		History     []string
	}
}

// Data structure for holding environment variables provided by Oahi via Chef.
type Env_Details struct {
	Sensu struct {
		Environment string `json:"environment"`
		FQDN        string `json:"fqdn"`
		Hostname    string `json:"hostname"`
	}
}

// Data structure for holding generic user data that is entered via the commandline.
type User_Event struct {
	Product   string
	Timestamp int64
	Data      string
}

// Generic error handling for all calling packages.
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Generate a simple id for use by ES and internal logging.
func event_name(client string, check string) string {
	return client + "_" + check
}

// Generate the name of the device being monitored.
// The monitored instance will map to the source field in the check result.
// It refers to the device that is actually being monitored, in many cases this
// will be the same device that the sensu-client is living on but in the case of
// SNMP traps or containers it may be different. More information can be found here:
// https://github.com/yieldbot/rhapthorne/blob/master/source/monitoring_terms.md#monitored-instance
//
// If the source field is not present then it will be set to the client name.
func (e Sensu_Event) acquire_monitored_instance() string {
	var monitored_instance string
	if e.Check.Source != "" {
		monitored_instance = e.Check.Source
	} else {
		monitored_instance = e.Client.Name
	}
	return monitored_instance
}

// Set the environment that the machine is running in based upon values
// dropped via Oahi during the Chef run.
func define_sensu_env(env string) string {
	switch env {
	case "prd":
		return "Prod "
	case "dev":
		return "Dev "
	case "stg":
		return "Stg "
	case "vagrant":
		return "Vagrant "
	default:
		return "Test "
	}
}

// Convert the check result status from an integer to a string.
// Feel free to add more cases here if you find yourself coming across them often.
func define_status(status int) string {
	switch status {
	case 0:
		return "OK"
	case 1:
		return "WARNING"
	case 2:
		return "CRITICAL"
	case 3:
		return "UNKNOWN"
	case 126:
		return "PERMISSION DENIED"
	case 127:
		return "CONFIG ERROR"
	default:
		return "ERROR"
	}
}

// Creates a check name that is easliy searchable in ES using different
// levels of granularity.
func create_check_name(check string) string {
	return check
}

// Calculate how long a check has been in a given state.
func define_check_state_duration() string {
	return " "
}

// Read in the environment details provided by Oahi and drop it into a struct.
func set_sensu_env() *Env_Details {
	env_file, err := ioutil.ReadFile("/etc/sensu/conf.d/monitoring_infra.json")
	if err != nil {
		check(err)
	}

	var env_details Env_Details
	err = json.Unmarshal(env_file, &env_details)
	if err != nil {
		check(err)
	}
	return &env_details
}

// Read in the check_result via stdin and drop it into a struct.
func (e Sensu_Event) acquire_sensu_event() *Sensu_Event {
	results, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		check(err)
	}
	err = json.Unmarshal(results, &e)
	if err != nil {
		check(err)
	}
	return &e
}

func main() {

	es_indexPtr := flag.String("index", "monitoring-status", "the elasticsearch index to use")
	es_hostPtr := flag.String("host", "elasticsearch.service.consul", "the elasticsearch host")
	es_portPtr := flag.String("port", "9200", "the elasticsearch port")
	stdinPtr := flag.Bool("read-stdin", true, "read input from stdin")
	//timePtr := flag.string("t-format", "", "time format to suffix on the index name")
	input_filePtr := flag.String("input-file", "", "file to read json in from, check docs for proper format")

	flag.Parse()
	es_index := *es_indexPtr
	es_type := "sensu"
	es_host := *es_hostPtr
	es_port := *es_portPtr
	rd_stdin := *stdinPtr
	input_file := *input_filePtr
	sensu_event := new(Sensu_Event)
	user_event := new(User_Event)
	//t_format := *timePtr

	sensu_env := set_sensu_env()

	//   if t_format != "" {
	//     // get the format of the time
	//     es_index = es_index + t_format
	//   }

	if (rd_stdin == false) && (input_file != "") {
		user_input, err := ioutil.ReadFile(input_file)
		if err != nil {
			check(err)
		}
		err = json.Unmarshal(user_input, &user_event)
		if err != nil {
			check(err)
		}
		es_type = "user"
	} else if (rd_stdin == false) && (input_file == "") {
		fmt.Printf("Please enter a file to read from")
		os.Exit(1)
	} else {
		sensu_event = sensu_event.acquire_sensu_event()
	}

	// Create a client
	client, err := elastic.NewClient(
		elastic.SetURL("http://" + es_host + ":" + es_port),
	)
	if err != nil {
		check(err)
	}

	// Check to see if the index exists and if not create it
	if client.IndexExists == nil {
		_, err = client.CreateIndex(es_index).Do()
		if err != nil {
			check(err)
		}
	}

	// Create an Elasticsearch document. The document type will define the mapping used for the document.
	doc := make(map[string]string)
	var doc_id string
	switch es_type {
	case "sensu":
		doc_id = event_name(sensu_event.Client.Name, sensu_event.Check.Name)
		doc["monitored_instance"] = sensu_event.acquire_monitored_instance()
		doc["sensu_client"] = sensu_event.Client.Name
		doc["incident_timestamp"] = time.Unix(sensu_event.Check.Issued, 0).Format(time.RFC822Z)
		doc["check_name"] = create_check_name(sensu_event.Check.Name)
		doc["check_state"] = define_status(sensu_event.Check.Status)
		doc["sensu_env"] = define_sensu_env(sensu_env.Sensu.Environment)
		doc["instance_address"] = sensu_event.Client.Address
		doc["check_state_duration"] = define_check_state_duration()
	case "user":
		doc["product"] = user_event.Product
		doc["data"] = user_event.Data
		doc["timestamp"] = time.Unix(user_event.Timestamp, 0).Format(time.RFC822Z)
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
		check(err)
	}

	// Log a successful document push to stdout. I don't add the id here as some id's are fixed but
	// the user has the ability to autogenerate an id if they don't want to provide one.
	fmt.Printf("Record added to ES\n")
}
