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
)

// EHndlr is for generic error handling in all Yieldbot monitoring packages.
func EHndlr(e error) {
	if e != nil {
		fmt.Printf("ERROR: %v", e)
		os.Exit(129)
	}
}

// Exit method for all sensu checks that will print the output and desired
// exit code
func Exit(args ...interface{}) {
	// YELLOW need to make sure that condition exists
	var exitCode int
	output := ""

	if 1 > len(args) {
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

		case 1: // x
			param, ok := p.(string)
			if !ok {
				panic("2nd parameter not type string.")
			}
			output = param

		default:
			panic("Wrong parameter count.")
		}
	}

	fmt.Printf("%v", output)
	os.Exit(exitCode)
}
