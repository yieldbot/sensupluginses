// Library for all handler data structures used by the Yieldbot Infrastructure
// teams in sensu.
//
// LICENSE:
//   Copyright 2015 Yieldbot. <devops@yieldbot.com>
//   Released under the MIT License; see LICENSE
//   for details.

package sensuhandler

// SensuEvent holds Sensu generated check results.
type SensuEvent struct {
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
		Tags        []string
	}
}

// EnvDetails holds environment variables provided by Oahi dropped via Chef.
type EnvDetails struct {
	Sensu struct {
		Environment string `json:"environment"`
		FQDN        string `json:"fqdn"`
		Hostname    string `json:"hostname"`
	}
}

// // UserEvent holds generic user data that is entered via an input file declared on the commndline.
// type UserEvent struct {
// 	Product   string
// 	Timestamp int64
// 	Data      string
// }

// Data structure for holding product configuration. Each product will have its own configuration. Then we just
// call that product on the commandline. Every value here is also represented via a commandline flag that will take preceence.
// type Config_Details struct {
// 	Sensu struct {
// 	}
// 	Elasticsearch struct {
// 		index []string
// 	}
// 	Slack struct {
// 		Channel []string
// 		Token   string
// 	}
// 	Pagerduty struct {
// 	}
// 	Statsd struct {
// 	}
// 	Mail struct {
// 	}
// }
