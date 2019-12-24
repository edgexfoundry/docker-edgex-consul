/********************************************************************************
 *  Copyright 2019 Dell Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const (
	ServiceNameIdentifier = "ServiceName"
	StatusIdentifier      = "Status"
)

var cfgFile string
var verboseOutput bool
var consulAgentCheckURL string
var desiredState string

type HTTPGetFunc func(string) (*http.Response, error)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "health service-name [flags]",
	Short: "Check the status of a service via Consul",
	Long: `Check the status of a service via Consul by querying Consul 
and inspecting the JSON response with JSON path to verify 
the service is in a desired state.
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		result, err := CheckServiceStatus(http.Get, consulAgentCheckURL, args[0], desiredState)

		if err != nil {
			log(err.Error())
		}

		if err != nil || !result {
			os.Exit(2)
		}
	},
	Example: "health edgex-mongo --url \"http://localhost:8500/v1/agent/checks\" --verbose",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.health.yaml)")
	rootCmd.PersistentFlags().BoolVar(&verboseOutput, "verbose", false, "Increase logging for debugging")
	rootCmd.PersistentFlags().StringVar(&consulAgentCheckURL, "url", "http://edgex-core-consul:8500/v1/agent/checks", "Consul agent check URL")
	rootCmd.PersistentFlags().StringVar(&desiredState, "desired-state", "passing", "The status state which indicates a service is healthy")

	rootCmd.Flags()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".health" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".health")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// checkServiceStatus verifies the service's status matches the expected state.
func CheckServiceStatus(
	getter HTTPGetFunc,
	consulAgentChecksURL,
	serviceName,
	successfulStatusState string) (bool, error) {

	jsonResponse, err := getRawStatusFromConsul(getter, consulAgentChecksURL)
	if err != nil {
		return false, err
	}
	serviceStatus, err := extractResultFromResponse(jsonResponse.(map[string]interface{}), serviceName)
	if err != nil {
		return false, err
	}

	log(fmt.Sprintf(`
##########################
Service Status Evaluation
--------------------------
Service's Status: %s
Acceptable Status: %s 
##########################`, serviceStatus, successfulStatusState))

	return serviceStatus == successfulStatusState, nil
}

// getRawStatusFromConsul retrieves the agent statuses from Consul and unmarshals the result.
func getRawStatusFromConsul(getter HTTPGetFunc, consulAgentCheckURL string) (interface{}, error) {
	r, err := getter(consulAgentCheckURL)
	if err != nil || (r.StatusCode < 200 && r.StatusCode > 299) || r.Body == nil {
		log(err.Error())
		return nil, errors.New("unable to connect to Consul")
	}

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	log(fmt.Sprintf(`
#################################
Information from Consul Request
-------------------------------
GET Request URL: %s
Response Status Code: %d
Response Body: '%s'
##################################`, consulAgentCheckURL, r.StatusCode, string(b)))

	var jsonResponse interface{}
	err = json.Unmarshal(b, &jsonResponse)
	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("marshal response from Consul")
	}

	log("Successfully read and unmarshaled response body.")

	return jsonResponse, nil
}

// extractResultFromResponse retrieves the service's status from a jsonResponse using the provided JSON path query and
// service name.
func extractResultFromResponse(jsonResponse map[string]interface{}, serviceName string) (string, error) {

	for _, element := range jsonResponse {
		node := element.(map[string]interface{})
		if node[ServiceNameIdentifier] == serviceName {
			status := node[StatusIdentifier]
			fmt.Println(status)
			return status.(string), nil
		}
	}

	return "", errors.New("no match for service name: " + serviceName)
}

func log(message string) {
	if verboseOutput {
		fmt.Println(message)
	}
}
