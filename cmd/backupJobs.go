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

// backupJobsCmd represents the backupJobs command
var backupJobsCmd = &cobra.Command{
	Use:     "backupJobs",
	Aliases: []string{"backup-jobs"},
	Short:   "Creates a backup of all jobs registered in Nomad",
	Long: `The backup-jobs command will created a new directory named with the current time
in seconds under the jobs-backup directory. All jobs register with Nomad will be written as 
a JSON file in the directory with the name of the job as the file name.`,
	Run: func(cmd *cobra.Command, args []string) {
		nh := new(nomadhelper.NomadHelper)
		nh.Init()
		nh.BackupJobs()
	},
}

func init() {
	rootCmd.AddCommand(backupJobsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// backupJobsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// backupJobsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
