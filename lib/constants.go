// Library for all constants used by the elasticsearch sensu packages
//
// LICENSE:
//   Copyright 2015 Yieldbot. <devops@yieldbot.com>
//   Released under the MIT License; see LICENSE
//   for details.

package lib

// Default values for connecting with and indexing Elasticsearch.
const (
	DefaultEsType string = "sensu"
	DefaultEsPort string = "9200"
	StatusEsIndex string = "monitoring-status"
	DefaultEsHost string = "localhost"
)
