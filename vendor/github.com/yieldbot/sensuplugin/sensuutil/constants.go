// Library for all constants used by the Yieldbot Infrastructure
// teams in sensu
//
// LICENSE:
//   Copyright 2015 Yieldbot. <devops@yieldbot.com>
//   Released under the MIT License; see LICENSE
//   for details.

package sensuutil
import  "github.com/op/go-logging"

// EnvironmentFile contains environmental details generated during the Chef run by Oahi.
const EnvironmentFile string = "/etc/sensu/conf.d/monitoring_infra.json"

type Password string
var SyslogFormat = logging.MustStringFormatter(
     `%{time:15:04:05.000} ▶ %{level} %{message}`,
)
var StderrFormat = logging.MustStringFormatter(
  `%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level} %{id:03x}%{color:reset} %{message}`,

)
