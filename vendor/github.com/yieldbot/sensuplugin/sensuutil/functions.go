// Library for general purpose functions used by the Yieldbot Infrastructure
// teams in sensu.
//
// LICENSE:
//   Copyright 2015 Yieldbot. <devops@yieldbot.com>
//   Released under the MIT License; see LICENSE
//   for details.

package sensuutil

import (
	"fmt"
	"os"
	"strings"

	"github.com/op/go-logging"
)

func (p Password) Redacted() interface{} {
	return logging.Redact(string(p))
}

// ConfigError is a generic message called when your config is boned.
func ConfigError() {
	fmt.Printf("You are missing a required configuration parameter")
	fmt.Printf("If unsure consult the documentation for examples and requirements\n")
	os.Exit(MonitoringErrorCodes["CONFIG_ERROR"])
}

// Exit method for all sensu checks that will print the output and desired exit code
// To use, pass it in the state you want and an optional text you would want outputted with the check.
//Ex. sensuutil.Exit("ok", "Everything is fine")
//    sensuutil.Exit("critical", variable)
// A list of error codes currently supported can be found in common.go
func Exit(args ...interface{}) {
	// YELLOW need to make sure that condition exists
	// panic is bad, need to add logging
	var exitCode int
	output := ""

	if len(args) == 0 {
		panic("Not enough parameters.")
	}

	for i, p := range args {
		switch i {
		case 0: // name
			param, ok := p.(string)
			if !ok {
				panic("1st parameter not type string.")
			}

			for k := range MonitoringErrorCodes {
				if k == strings.ToUpper(param) {
					exitCode = MonitoringErrorCodes[k]
				}
			}

		case 1: // optional text
			param, ok := p.(string)
			if !ok {
				panic("2nd parameter not type string.")
			}
			output = param

		default:
			panic("Incorrect parameters")
		}
	}

	fmt.Printf("%v\n", output)
	os.Exit(exitCode)
}
