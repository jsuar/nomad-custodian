// Package cmd deleteAllJobs contains the deleteAllJobs functionality
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

// deleteAllJobsCmd represents the deleteAllJobs command
var deleteAllJobsCmd = &cobra.Command{
	Use:     "deleteAllJobs",
	Aliases: []string{"delete-all-jobs"},
	Short:   "Deletes all jobs currently registered with Nomad",
	Long: `The delete-all-jobs command will loop through all jobs currently registered
with Nomad and deregister them. If the purge flag is set, then purge=true will be
passed in the deregistration call.`,
	Run: func(cmd *cobra.Command, args []string) {
		force, _ := cmd.Flags().GetBool("force")
		purge, _ := cmd.Flags().GetBool("purge")
		autoApprove, _ := cmd.Flags().GetBool("auto-approve")
		verbose, _ := cmd.Flags().GetBool("verbose")

		nh := new(nomadhelper.NomadHelper)
		nh.Init()
		nh.DeleteAllJobs(force, autoApprove, purge, verbose)
	},
}

func init() {
	rootCmd.AddCommand(deleteAllJobsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteAllJobsCmd.PersistentFlags().String("foo", "", "A help for foo")
	deleteAllJobsCmd.PersistentFlags().BoolP("force", "f", false, "Force action")
	deleteAllJobsCmd.PersistentFlags().BoolP("auto-approve", "", false, "Skip user confirmation")
	deleteAllJobsCmd.PersistentFlags().BoolP("purge", "p", false, "Purge job data from Nomad after deregister")
	deleteAllJobsCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose output")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteAllJobsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
