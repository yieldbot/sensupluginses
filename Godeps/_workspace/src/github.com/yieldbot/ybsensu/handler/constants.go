// Library for all handler related constants used by the Yieldbot Infrastructure
// teams in sensu
//
// LICENSE:
//   Copyright 2015 Yieldbot. <devops@yieldbot.com>
//   Released under the MIT License; see LICENSE
//   for details.

package handler

// EnvironmentFile contains environmental details generated during the Chef run by Oahi.
const EnvironmentFile string = "/etc/sensu/conf.d/monitoring_infra.json"

// Error codes for applications.
// Please use the below codes instead of random non-zero so that monitoring can
// utilize existing maps for alerting and help avoid unnecessary noise.
const (
	GeneralGolangError int = 129
	ConfigError        int = 127
	PermissionError    int = 126
	RuntimeError       int = 42
)

// create a map for sensu environments and read from that

// Map of all known exit codes accepted by the monitoring stack
// STATUS_map = make(map[int]string)
// STATUS_map[0] = "OK"
// STATUS_map[1] = "WARNING"
// STATUS_map[2] = "CRITICAL"
// STATUS_map[3] = "UNKNOWN"
// STATUS_map[126] = "PERMISSION DENIED"
// STATUS_map[127] = "CONFIG ERROR"

// // Map of all known monitoring environments
// ENV_map := make(map[string]string)
// ENV_map["prd"] = "Prod "
// ENV_map["dev"] = "Dev "
// ENV_map["stg"] = "Stg "
// ENV_map["vagrant"] = "Vagrant "
// ENV_map["default"] = "Test "
