# Nomad Custodian

Inspired by [Cloud Custodian](https://github.com/cloud-custodian/cloud-custodian), this simple CLI will help Nomad administrators manage job resources with cost optimization and maintenance in mind.

## Features
* Scale in all job task group counts to `count=1` during off business hours
* Scale out all jobs to original counts
* Delete all jobs
* Backup all jobs as JSON files

# How to use

To list the meta tags and current task group counts for each job:

```
nomad-custodian list
```

To what changes will take place in a scale in all job task group counts to `1`:

```
nomad-custodian scale-in
```

To scale in all job task group counts to `1`:

```
nomad-custodian scale-in --force
```

To scale out all job task group counts to their original numbers:

```
nomad-custodian scale-out
```

To backup all jobs registered in Nomad:

```
./nomad-custodian backup-jobs
```

To delete and purge all jobs registered in Nomad:

```
nomad-custodian delete-all-jobs -f -p
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