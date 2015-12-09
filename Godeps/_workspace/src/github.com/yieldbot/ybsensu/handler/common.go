// Library for all non-constant variables used by the Yieldbot Infrastructure teams
//
// LICENSE:
//   Copyright 2015 Yieldbot. <devops@yieldbot.com>
//   Released under the MIT License; see LICENSE
//   for details.

package handler

// Debug  Do we print debug statements or not. This is set in each binary but is placed here
// to avoid the use of global variables
var Debug bool

// JavaApp  This is used to let the process -> pid function know how it will match the process name
var JavaApp bool

var SlackChannels = map[string]string {
"monitoring-test" : "xxx",
"devops-alerts" : "xxx",
}
