package nomadhelper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	nomad "github.com/hashicorp/nomad/api"
	"github.com/ryanuber/columnize"
)

// NomadHelper provides custodian helper functions
type NomadHelper struct {
	Client *nomad.Client
}

// ScaleType specifies scaling in or out
type ScaleType int

// ScaleType specifies scaling in or out
const (
	ScaleIn ScaleType = iota
	ScaleOut
)

func (d ScaleType) String() string {
	return [...]string{"In", "Out"}[d]
}

// Init will initialize the NomadHelper object
func (n *NomadHelper) Init() {
	var err error
	nomadConfig := nomad.DefaultConfig()
	n.Client, err = nomad.NewClient(nomadConfig)
	if err != nil {
		fmt.Printf("Error encountered:%s\n", err)
	}
}

func check(e error) {
	if e != nil {
		// panic(e)
		fmt.Printf("%s\n", e)
	}
}

// DisplayJobDiff prints the simplified diff between job versions
func DisplayJobDiff(diff nomad.JobDiff) {
	var output []string

	// Display job plan diff
	output = nil
	output = append(output, "|What's Changing|From|To")
	fields := diff.Fields
	for _, field := range fields {
		output = append(output, "|"+field.Name+"|"+field.Old+"|"+field.New)
	}
	for _, taskGroup := range diff.TaskGroups {
		for _, taskField := range taskGroup.Fields {
			output = append(output, "|"+taskField.Name+"|"+taskField.Old+"|"+taskField.New)
		}
	}
	for _, object := range diff.Objects {
		fmt.Printf(object.Name, "\n")
	}
	result := columnize.SimpleFormat(output)
	fmt.Printf("%s\n\n", result)
}

// ScaleInJobs scales all jobs in to count=1
func (n *NomadHelper) ScaleInJobs(force bool, verbose bool) {
	var output []string
	var jobsSkipped []string
	var wg sync.WaitGroup

	if !force && verbose {
		fmt.Printf("Running a plan scale down action.\n")
	}

	jobs := n.Client.Jobs()
	jobStubList, _, err := jobs.List(nil)
	if err != nil {
		fmt.Printf("Error encountered: %s \n", err)
	}

	if verbose {
		fmt.Printf("Number of jobs running: %d\n", len(jobStubList))
	}

	for _, jobStub := range jobStubList {
		// Get the jobs object
		jobInfo, _, err := jobs.Info(jobStub.ID, nil)
		if err != nil {
			fmt.Printf("Error encountered: %s", err)
		}

		custodianIgnore, err := strconv.ParseBool(jobInfo.Meta["custodian-ignore"])
		if err != nil {
			if jobInfo.Meta["custodian-ignore"] != "" {
				fmt.Printf("Error encountered:%s\n", err)
			}
		}

		alreadyScaledIn := jobInfo.Meta["custodian-action"] == "scaled-in"
		jobIsRunning := *jobInfo.Status == "running"
		criteriaToScaleIn := !alreadyScaledIn && !custodianIgnore && jobIsRunning

		if criteriaToScaleIn {
			// Update job count
			scaledDownJobCount := new(int)
			*scaledDownJobCount = 1
			for _, taskGroup := range jobInfo.TaskGroups {
				key := fmt.Sprintf("custodian-%s-count", *taskGroup.Name)
				jobInfo.SetMeta(key, fmt.Sprint(*taskGroup.Count))
				taskGroup.Count = scaledDownJobCount
			}
			// Update meta kv
			jobInfo.SetMeta("custodian-action", "scaled-in")
			jobInfo.SetMeta("custodian-revert-version", fmt.Sprint(*jobInfo.Version))

		} else {
			jobsSkipped = append(jobsSkipped, fmt.Sprintf("%s|%s|%t", *jobInfo.Name,
				jobInfo.Meta["custodian-action"], custodianIgnore))
			continue
		}

		// Plan the change and get the response/diff
		jobPlanResponse, _, err := jobs.Plan(jobInfo, true, nil)
		if err != nil {
			fmt.Printf("Error encountered: %s", err)
		}
		diff := *jobPlanResponse.Diff
		fmt.Printf("Job: %s, %s\n", *jobInfo.Name, *jobInfo.Status)
		DisplayJobDiff(diff)

		if force {
			wg.Add(1)
			go n.ApplyChanges(jobInfo, &wg)
		}
	}

	output = append(output, "Jobs Skipped|Scale Status|Ignore")
	if len(jobsSkipped) == 0 {
		output = append(output, "None")
	} else {
		output = append(output, jobsSkipped...)
	}
	result := columnize.SimpleFormat(output)
	fmt.Printf("%s\n", result)

	wg.Wait()
}

// ScaleOutJobs scales all jobs the original count
func (n *NomadHelper) ScaleOutJobs(force bool, verbose bool) {
	var output []string
	var jobsSkipped []string

	jobs := n.Client.Jobs()
	jobStubList, _, err := jobs.List(nil)
	if err != nil {
		fmt.Printf("Error encountered: %s", err)
	}

	if verbose {
		fmt.Printf("Number of jobs running: %d\n", len(jobStubList))
	}

	for _, jobStub := range jobStubList {
		// Get the jobs object
		jobInfo, _, err := jobs.Info(jobStub.ID, nil)
		if err != nil {
			fmt.Printf("Error encountered:%s\n", err)
		}

		custodianIgnore, err := strconv.ParseBool(jobInfo.Meta["custodian-ignore"])
		if err != nil {
			if jobInfo.Meta["custodian-ignore"] != "" {
				fmt.Printf("Error encountered:%s\n", err)
			}
		}

		alreadyScaledIn := jobInfo.Meta["custodian-action"] == "scaled-in"
		jobIsRunning := *jobInfo.Status == "running"
		criteriaToScaleOut := alreadyScaledIn && !custodianIgnore && jobIsRunning

		// Only proceed if job was scaled in using the tooling
		if criteriaToScaleOut {
			// Convert to uint64 for revert function
			previousVer, err := strconv.ParseUint(jobInfo.Meta["custodian-revert-version"], 10, 64)
			if err != nil {
				fmt.Printf("Error encountered:%s\n", err)
			}

			includeDiffs := false
			pastJobs, _, _, err := jobs.Versions(jobStub.ID, includeDiffs, nil)
			if err != nil {
				fmt.Printf("Error encountered:%s\n", err)
			}

			for _, pastJob := range pastJobs {
				if *pastJob.Version == previousVer {
					// Plan the change and get the response/diff
					jobPlanResponse, _, err := jobs.Plan(pastJob, true, nil)
					if err != nil {
						fmt.Printf("Error encountered:%s\n", err)
					}
					diff := *jobPlanResponse.Diff
					fmt.Printf("Job: %s, %s\n", *jobInfo.Name, *jobInfo.Status)
					DisplayJobDiff(diff)
					break
				}
			}

			if force {
				// Handle revert response
				jobRegisterResponse, _, err := jobs.Revert(*jobInfo.ID, previousVer, nil, nil, "")
				if err != nil {
					fmt.Printf("Error encountered:%s\n", err)
				}
				if jobRegisterResponse.Warnings != "" {
					fmt.Printf("Warnings: %s\n", jobRegisterResponse.Warnings)
				}
			}
		} else {
			jobsSkipped = append(jobsSkipped, fmt.Sprintf("%s|%s|%t", *jobInfo.Name,
				jobInfo.Meta["custodian-action"], custodianIgnore))
		}
	}

	output = append(output, "Jobs Skipped|Scale Status|Ignore")
	if len(jobsSkipped) == 0 {
		output = append(output, "None")
	} else {
		output = append(output, jobsSkipped...)
	}
	result := columnize.SimpleFormat(output)
	fmt.Printf("%s\n", result)
}

// ApplyChanges will register the job and any changes it has with Nomad
func (n *NomadHelper) ApplyChanges(job *nomad.Job, wg *sync.WaitGroup) {
	jobs := n.Client.Jobs()

	// fmt.Printf("\nApplying changes to job %s\n", *job.Name)
	jobRegisterResponse, _, err := jobs.Register(job, nil)
	if err != nil {
		fmt.Printf("Error encountered:%s\n", err)
	}
	if jobRegisterResponse.Warnings != "" {
		fmt.Printf("Warnings: %s\n", jobRegisterResponse.Warnings)
	}

	wg.Done()
}

// ListJobs scales all jobs in to count=1 or out to the jobs original count
func (n *NomadHelper) ListJobs(verbose bool) {
	var output []string

	nomadConfig := nomad.DefaultConfig()
	nomadClient, err := nomad.NewClient(nomadConfig)
	if err != nil {
		fmt.Printf("Error encountered:%s\n", err)
	}

	jobs := nomadClient.Jobs()
	jobStubList, _, err := jobs.List(nil)
	if err != nil {
		fmt.Printf("Error encountered:%s\n", err)
	}

	if verbose {
		fmt.Printf("Number of jobs running: %d\n", len(jobStubList))
	}

	for _, jobStub := range jobStubList {
		// Get the jobs object
		jobInfo, _, err := jobs.Info(jobStub.ID, nil)
		if err != nil {
			fmt.Printf("Error encountered:%s\n", err)
		}

		// Display job plan diff
		output = append(output, fmt.Sprintf("+|Job: %s|Status: %s|", *jobInfo.Name, *jobInfo.Status))
		output = append(output, "|Field|Value|")
		for _, taskGroup := range jobInfo.TaskGroups {
			count := strconv.Itoa(*taskGroup.Count)
			output = append(output, "|Count|"+count+"|")
		}
		for k, v := range jobInfo.Meta {
			output = append(output, fmt.Sprintf("|%s|%s|", k, v))
		}
	}

	result := columnize.SimpleFormat(output)
	fmt.Printf("\n%s\n", result)
}

// AskForConfirmation prompts the user for confirmation before proceeding
func AskForConfirmation() bool {
	var s string

	fmt.Printf("(y/N): ")
	_, err := fmt.Scan(&s)
	if err != nil {
		panic(err)
	}

	s = strings.TrimSpace(s)
	s = strings.ToLower(s)

	if s == "y" || s == "yes" {
		return true
	}
	return false
}

// DeleteAllJobs deregisters all jobs currently running in Nomad
func (n *NomadHelper) DeleteAllJobs(force bool, autoApprove bool, purge bool, verbose bool) {
	var output []string
	var jobsSkipped []string
	var userConfirmation bool

	jobs := n.Client.Jobs()
	jobStubList, _, err := jobs.List(nil)
	if err != nil {
		fmt.Printf("Error encountered: %s", err)
	}

	if force {
		if autoApprove {
			userConfirmation = true
		} else {
			userConfirmation = AskForConfirmation()
		}
	}

	for _, jobStub := range jobStubList {
		// Get the jobs object
		jobInfo, _, err := jobs.Info(jobStub.ID, nil)
		if err != nil {
			fmt.Printf("Error encountered:%s\n", err)
		}

		custodianIgnore, err := strconv.ParseBool(jobInfo.Meta["custodian-ignore"])
		if err != nil {
			if jobInfo.Meta["custodian-ignore"] != "" {
				fmt.Printf("Error encountered:%s\n", err)
			}
		}

		if custodianIgnore {
			jobsSkipped = append(jobsSkipped, fmt.Sprintf("%s|%s|%t", *jobInfo.Name,
				jobInfo.Meta["custodian-action"], custodianIgnore))
			continue

		} else {
			if userConfirmation {
				deregisterResponse, _, err := jobs.Deregister(jobStub.ID, purge, nil)
				if err != nil {
					fmt.Printf("Error encountered: %s", err)
				}
				fmt.Printf("Job %s deregister response: %s", jobStub.Name, deregisterResponse)

			}
			fmt.Printf("Action: Deregister, Job: %s\n", jobStub.Name)
		}
	}

	output = append(output, "Jobs Skipped|Scale Status|Ignore")
	if len(jobsSkipped) == 0 {
		output = append(output, "None")
	} else {
		output = append(output, jobsSkipped...)
	}
	result := columnize.SimpleFormat(output)
	fmt.Printf("\n%s\n", result)
}

// BackupJobs ...
func (n *NomadHelper) BackupJobs() {
	jobs := n.Client.Jobs()
	jobStubList, _, err := jobs.List(nil)
	if err != nil {
		fmt.Printf("Error encountered: %s", err)
	}

	dir := fmt.Sprintf("jobs-backup")
	err = os.Mkdir(dir, 0755)
	check(err)

	now := time.Now()
	secs := now.Unix()
	dir = fmt.Sprintf("jobs-backup/%d/", secs)
	err = os.Mkdir(dir, 0755)
	check(err)

	for _, jobStub := range jobStubList {
		// Get the jobs object
		jobInfo, _, err := jobs.Info(jobStub.ID, nil)
		if err != nil {
			fmt.Printf("Error encountered:%s\n", err)
		}

		jobJSON, _ := json.Marshal(jobInfo)
		filename := fmt.Sprintf("%s.json", *jobInfo.Name)
		fmt.Printf("Writing %s to a file\n", *jobInfo.Name)

		err = ioutil.WriteFile(filepath.Join(dir, filepath.Base(filename)), jobJSON, 0644)
		if err != nil {
			panic(err)
		}
	}
}
