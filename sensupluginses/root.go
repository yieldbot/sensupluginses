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
	"fmt"
	"log/syslog"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/hooks/syslog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yieldbot/sensuplugin/sensuutil"
	//"github.com/yieldbot/sensupluginses/version"
)

// Configuration via Viper
var cfgFile string

// Hostname for logging
var host string

// Create a logging instance.
var syslogLog = logrus.New()

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "sensupluginses",
	//Short: fmt.Sprintf("An elasticsearch handler for Sensu - (%s)", version.AppVersion()),
	Long: `This plugin currently contains a single handler that will drop
Sensu check results into Elasticsearch.`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Setup logging for the package. Doing it here is much eaiser than in each
	// binary. If you want to overwrite it in a specific binary then feel free.
	hook, err := logrus_syslog.NewSyslogHook("", "", syslog.LOG_INFO, "")
	if err != nil {
		panic(err)
	}
	syslogLog.Hooks.Add(hook)
	syslogLog.Formatter = new(logrus.JSONFormatter)

	// Set the hostname for use in logging within the package. Doing it here is
	// cleaner than in each binary but if you want to use some other method just
	// override the variable in the specific binary.
	host, err = os.Hostname()
	if err != nil {
		syslogLog.WithFields(logrus.Fields{
			"check":   "sensupluginses",
			"client":  "unknown",
			//"version": version.AppVersion(),
			"error":   err,
		}).Error(`Could not determine the hostname of this machine as reported by the kernel.`)
		sensuutil.Exit("GENERALGOLANGERROR")
	}

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.sensupluginses.yaml)")
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("sensupluginses")
		viper.AddConfigPath("/etc/sensuplugins/conf.d")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
