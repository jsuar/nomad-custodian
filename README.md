# Nomad Custodian

Inspired by [Cloud Custodian](https://github.com/cloud-custodian/cloud-custodian), this simple CLI helps Nomad administrators manage Nomad job resources with cost optimization in mind.

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
nomad-custodian scaleIn
```

To scale in all job task group counts to `1`:

```
nomad-custodian scaleIn --force
```

To scale out all job task group counts to their original numbers:

```
nomad-custodian scaleOut
```

## Safety Controls

Prevent any custodian actions:

```
?
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