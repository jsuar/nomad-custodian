# Nomad Custodian

Inspired by [Cloud Custodian](https://github.com/cloud-custodian/cloud-custodian), this simple CLI will help Nomad administrators manage job resources with cost optimization and maintenance in mind.

## Features
* Scale in all job task group counts to `count=1` during off business hours
* Scale out all jobs to original counts
* Delete all jobs
* Backup all jobs as JSON files

# How to use

## `list`
Running `nomad-custodian list` will list the meta tags and current task group counts for each job.

```
$ nomad-custodian list
Number of jobs running: 4

+  Job: couchbase    Status: running
   Field             Value
   Count             2
+  Job: demo-webapp  Status: running
   Field             Value
   Count             3
+  Job: example      Status: running
   Field             Value
   Count             2
+  Job: nginx        Status: pending
   Field             Value
   Count             2
   custodian-ignore  1
```
## `scale-in`
Excluding `--force` or `-f` with the `scale-in` and `scale-out` commands will provide a preview of what will change. For example, running `nomad-custodian scale-in` will provide the below output. 

```
$ nomad-custodian scale-in
Job: couchbase, running
  What's Changing                  From  To
  Meta[custodian-action]                 scaled-in
  Meta[custodian-couchbase-count]        2
  Meta[custodian-revert-version]         1
  Count                            2     1

Job: demo-webapp, running
  What's Changing                 From  To
  Meta[custodian-action]                scaled-in
  Meta[custodian-demo-count]            3
  Meta[custodian-revert-version]        2
  Count                           3     1

Job: example, running
  What's Changing                 From  To
  Meta[custodian-action]                scaled-in
  Meta[custodian-cache-count]           2
  Meta[custodian-revert-version]        2
  Count                           2     1

Jobs Skipped  Scale Status  Ignore
nginx                       true
```

Including the `--force` flag will produce similar output as the plan but the changes will take place.

```
$ nomad-custodian scale-in --force
Job: couchbase, running
  What's Changing                  From  To
  Meta[custodian-action]                 scaled-in
  Meta[custodian-couchbase-count]        2
  Meta[custodian-revert-version]         1
  Count                            2     1

Job: demo-webapp, running
  What's Changing                 From  To
  Meta[custodian-action]                scaled-in
  Meta[custodian-demo-count]            3
  Meta[custodian-revert-version]        2
  Count                           3     1

Job: example, running
  What's Changing                 From  To
  Meta[custodian-action]                scaled-in
  Meta[custodian-cache-count]           2
  Meta[custodian-revert-version]        2
  Count                           2     1

Jobs Skipped  Scale Status  Ignore
nginx                       true
```

## `scale-out`
The `scale-out` command is similar to the `scale-in` command in terms of output.

```
$ nomad-custodian scale-out -f
Job: couchbase, running
  What's Changing                  From       To
  Meta[custodian-action]           scaled-in
  Meta[custodian-couchbase-count]  2
  Meta[custodian-revert-version]   1
  Count                            1          2

Job: demo-webapp, running
  What's Changing                 From       To
  Meta[custodian-action]          scaled-in
  Meta[custodian-demo-count]      3
  Meta[custodian-revert-version]  2
  Count                           1          3

Job: example, running
  What's Changing                 From       To
  Meta[custodian-action]          scaled-in
  Meta[custodian-cache-count]     2
  Meta[custodian-revert-version]  2
  Count                           1          2

Jobs Skipped  Scale Status  Ignore
nginx                       true
```

## `backup-jobs`

The `backup-jobs` command provides an easy way to locally backup all the jobs registered in Nomad as JSON files. A new time stamped directory is created each time the command is executed.

```
$ nomad-custodian backup-jobs
mkdir jobs-backup: file exists
Job couchbase written to couchbase.json
Job demo-webapp written to demo-webapp.json
Job example written to example.json
Job nginx written to nginx.json

$ ls jobs-backup/1578492852
couchbase.json   demo-webapp.json example.json     nginx.json
```

## `delete-all-jobs`

The `delete-all-jobs` helps make bulk deregistering of jobs (and purging if `--purge` or `-p` is included) from Nomad.

```
nomad-custodian delete-all-jobs -f -p
Are you sure you want to continue? (y/N): y
Job couchbase deregister response: ece44f6c-e518-bbbe-7f06-41ee4f3b61c8Action: Deregister, Job: couchbase
Job demo-webapp deregister response: 97f82a9d-ddd1-dc31-1be6-e5e81440b00fAction: Deregister, Job: demo-webapp
Job example deregister response: b8c9885e-d87c-9d2c-fbdc-2b1f42a57422Action: Deregister, Job: example

Jobs Skipped  Scale Status  Ignore
nginx                       true
```

## Safety Controls

Prevent any custodian actions:

```
Nothing available at the moment
```

Prevent changes on specific jobs:
```
job "nginx" {
  datacenters = ["dc1"]

  meta {
    nomad-custodian-ignore = true
  }
  ...
```

# Development

To build the binary:

```
make build
```

# Improvement / Feature Ideas
* Service and UI components
  * Service could be deployed to the same cluster or a management cluster
  * UI would provide same functionality as the CLI
* Filtering capabilities
  * Namespaces
  * Job names
  * Time of day
* Globally prevent custodian changes
  * Enforce with Consul KV check
  * Enforce some other way with Nomad