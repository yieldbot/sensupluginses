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
