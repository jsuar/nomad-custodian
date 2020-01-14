// Package cmd handles all CLI calls
/*
Copyright © 2020 John Suarez jsuar@users.noreply.github.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/jsuar/nomad-custodian/pkg/nomadhelper"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all register jobs with relevant info",
	Long: `The list command is a quick view of all jobs registered with Nomad with 
count and meta data relevant to the nomad-custodian CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		nhelper := new(nomadhelper.NomadHelper)
		nhelper.Init()
		nhelper.ListJobs(true)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
