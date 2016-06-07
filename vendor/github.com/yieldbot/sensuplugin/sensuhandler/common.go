// Library for all non-constant handler related variables used by the Yieldbot
// Infrastructure teams in sensu
//
// LICENSE:
//   Copyright 2015 Yieldbot. <devops@yieldbot.com>
//   Released under the MIT License; see LICENSE
//   for details.

package sensuhandler

// Notification Color is a set of colors used by all handlers to provide a rich
// notification.
var NotificationColor = map[string]string{
	"green":  "#33CC33",
	"orange": "#FFFF00",
	"red":    "#FF0000",
	"yellow": "#FF6600",
}
