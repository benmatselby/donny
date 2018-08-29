# Donny

[![Build Status](https://travis-ci.org/benmatselby/donny.png?branch=master)](https://travis-ci.org/benmatselby/donny)
[![Go Report Card](https://goreportcard.com/badge/github.com/benmatselby/donny?style=flat-square)](https://goreportcard.com/report/github.com/benmatselby/donny)

_Forget it, Donny, you're out of your element!_

CLI application for getting information out of Visual Studio Team Services. It's based on [Trello CLI](https://github.com/benmatselby/trello-cli) so the aims are the same:

```
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

* Go version 1.10+
* [Dep installed](https://github.com/golang/dep)

## Configuration

You will need the following environment variables defining:

```
$ export VSTS_ACCOUNT=""
$ export VSTS_PROJECT=""
$ export VSTS_TEAM=""
$ export VSTS_TOKEN=""
```

## Installation via Docker

Other than requiring [docker](http://docker.com) to be installed, there are no other requirements to run the application this way. This is the preferred method of running the `donny`. The image is [here](https://hub.docker.com/r/benmatselby/donny/).

```
$ docker run \
    --rm \
    -t \
    -eVSTS_ACCOUNT \
    -eVSTS_PROJECT \
    -eVSTS_TEAM \
    -eVSTS_TOKEN \
    benmatselby/donny "$@"
```

## Installation via Git

```
$ git clone git@github.com:benmatselby/donny.git
$ cd donny
$ make all
$ ./donny iteration:list
```
