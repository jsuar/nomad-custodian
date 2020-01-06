GOOS=linux
GOARCH=386

.PHONY: build

# GIT_COMMIT := $(shell git rev-list -1 HEAD)
GIT_COMMIT := 9jd3j90j3

.PHONY: build
build:
	# GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "-X main.gitCommit=$(GIT_COMMIT)" .
	go build -ldflags "-X 'github.com/jsuar/nomad-custodian/cmd.gitCommit=$(GIT_COMMIT)'" .

.PHONY: load-jobs
load-jobs:
	nomad run testing/demo-webapp.nomad
	nomad run testing/nginx.nomad
	nomad run testing/redis.nomad