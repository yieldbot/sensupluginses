// Copyright © 2016 Yieldbot <devops@yieldbot.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package sensupluginses

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/olivere/elastic"
	"github.com/spf13/cobra"
	"github.com/yieldbot/sensuplugin/sensuhandler"
	"github.com/yieldbot/sensupluginses/version"
)

// elasticsearch endpoint configuration
var esHost string
var esIndex string
var esPort string
var esType = DefaultEsType

// Bring in the environmant details
var sensuEnv = new(sensuhandler.EnvDetails)

// main entry point
var handlerElasticsearchStatusCmd = &cobra.Command{
	Use:   "handlerElasticsearchStatus --index <index> --host <host> --port <port>",
	Short: "This will input a single record for each check result given, overwriting the currect record.",
	Long: `This will take a single check result and create a key based upon the host name and
  the check name. This key will remain consistent so that only the latest status will be available in
  the index. This is designed to allow the creation of current dashboards from Kibana or Dashing.`,

	Run: func(sensupluginses *cobra.Command, args []string) {

		// read in the event data from the sensu server
		sensuEvent := new(sensuhandler.SensuEvent)
		sensuEvent = sensuEvent.AcquireSensuEvent()

		// set the environment this is running in (prd, dev,stg)
		sensuEnv = sensuEnv.SetSensuEnv()

		// Create a client
		client, err := elastic.NewClient(
			elastic.SetURL("http://" + esHost + ":" + esPort),
		)
		if err != nil {
			syslogLog.WithFields(logrus.Fields{
				"check":   "sensupluginses",
				"client":  host,
				"version": version.AppVersion(),
				"error":   err,
				"esHost":  esHost,
				"esPort":  esPort,
			}).Error(`Could not create an elasticsearch client`)

		}

		// Check to see if the index exists and if not create it
		if client.IndexExists(esIndex) == nil { // need to test to make sure this does what I want
			_, err = client.CreateIndex(esIndex).Do()
			if err != nil {
				syslogLog.WithFields(logrus.Fields{
					"check":   "sensupluginses",
					"client":  host,
					"version": version.AppVersion(),
					"error":   err,
					"esIndex": esIndex,
				}).Error(`Could not create an elasticsearch index`)

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
			syslogLog.WithFields(logrus.Fields{
				"check":   "sensupluginses",
				"client":  host,
				"version": version.AppVersion(),
				"error":   err,
				"esHost":  esHost,
				"esPort":  esPort,
				"esIndex": esIndex,
			}).Error(`Could not post a document to elasticsearch`)

		}

		// Log a successful document push to stdout. I don't add the id here as some id's are fixed but
		// the user has the ability to autogenerate an id if they don't want to provide one.
		syslogLog.WithFields(logrus.Fields{
			"check":   "sensupluginses",
			"client":  host,
			"version": version.AppVersion(),
			"esHost":  esHost,
			"esPort":  esPort,
			"esIndex": esIndex,
		}).Info(`Document posted to elasticsearch`)
	},
}

func init() {
	RootCmd.AddCommand(handlerElasticsearchStatusCmd)

	// set commandline flags
	handlerElasticsearchStatusCmd.Flags().StringVarP(&esHost, "host", "", DefaultEsHost, "the elasticsearch host")
	handlerElasticsearchStatusCmd.Flags().StringVarP(&esIndex, "index", "", StatusEsIndex, "the es index to populate")
	handlerElasticsearchStatusCmd.Flags().StringVarP(&esPort, "port", "", DefaultEsPort, "the elasticsearch port")

}
