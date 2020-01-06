// Package cmd scaleOut contains the scaleOut functionality
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

// scaleUpCmd represents the scaleIn command
var scaleInCmd = &cobra.Command{
	Use:     "scaleIn",
	Aliases: []string{"scale-in"},
	Short:   "Scales in job task groups count to 1",
	Long: `The scale-in command will loop through all
jobs registered in Nomad and set count=1 for all tasks
in the task group. Jobs will be skipped if they:
* Are already scaled in to count=1
* Have the custodian-ignore=false meta key value set
* Are not running`,
	Run: func(cmd *cobra.Command, args []string) {
		force, _ := cmd.Flags().GetBool("force")
		verbose, _ := cmd.Flags().GetBool("verbose")

		nhelper := new(nomadhelper.NomadHelper)
		nhelper.Init()
		nhelper.ScaleInJobs(force, verbose)
	},
}

func init() {
	rootCmd.AddCommand(scaleInCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	scaleInCmd.PersistentFlags().BoolP("force", "f", false, "Force action")
	scaleInCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose output")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scaleInCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
