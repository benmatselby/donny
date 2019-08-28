# Donny

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=donny&metric=alert_status)](https://sonarcloud.io/dashboard?id=donny)
[![Go Report Card](https://goreportcard.com/badge/github.com/benmatselby/donny?style=flat-square)](https://goreportcard.com/report/github.com/benmatselby/donny)

_Forget it, Donny, you're out of your element!_

CLI application for getting information out of Azure DevOps. It's based on [Trello CLI](https://github.com/benmatselby/trello-cli) so the aims are the same:

```shell
COMMANDS:
     help, h  Shows a list of commands or help for one command
   build:
     build:list, bl      List all the builds
     build:overview, bo  Show build overview for build definitions in a given path
   code:
     code:branches, cb  Show branch information for a repo
   iteration:
     iteration:burndown, ib  Show column based data for the iteration
     iteration:items, ii     List all the work items in a given iteration
     iteration:list, il      List all the iterations
     iteration:people, ip    Show people based data for the iteration
   plans:
     plan:list, pll      List all the delivery plans
     plan:timeline, plt  Show the timeline for the delivery plan
   pull requests:
     pr:list, pul  List all the pull requests
```

## Requirements

If you are wanting to build and develop this, you will need the following items installed. If, however, you just want to run the application I recommend using the docker container (See below)

- Go version 1.12+

## Configuration

You will need the following environment variables defining:

```bash
export AZURE_DEVOPS_ACCOUNT=""
export AZURE_DEVOPS_PROJECT=""
export AZURE_DEVOPS_TEAM=""
export AZURE_DEVOPS_TOKEN=""
```

## Installation via Docker

Other than requiring [docker](http://docker.com) to be installed, there are no other requirements to run the application this way. This is the preferred method of running the `donny`. The image is [here](https://hub.docker.com/r/benmatselby/donny/).

```bash
$ docker run \
    --rm \
    -t \
    -eAZURE_DEVOPS_ACCOUNT \
    -eAZURE_DEVOPS_PROJECT \
    -eAZURE_DEVOPS_TEAM \
    -eAZURE_DEVOPS_TOKEN \
    benmatselby/donny "$@"
```

## Installation via Git

```bash
git clone git@github.com:benmatselby/donny.git
cd donny
make all
./donny iteration:list
```

You can also install into your `$GOPATH/bin` by `go install`
