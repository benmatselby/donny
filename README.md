Donny
=====

[![Build Status](https://travis-ci.org/benmatselby/donny.png?branch=master)](https://travis-ci.org/benmatselby/donny)
[![Go Report Card](https://goreportcard.com/badge/github.com/benmatselby/donny?style=flat-square)](https://goreportcard.com/report/github.com/benmatselby/donny)

_Forget it, Donny, you're out of your element!_

CLI application for getting information out of Visual Studio Team Services. It's based on [Trello CLI](https://github.com/benmatselby/trello-cli) so the aims are the same:

```
COMMANDS:
     help, h  Shows a list of commands or help for one command
   build:
     build:list          List all the builds
   iteration:
     iteration:burndown  Show column based data for the iteration
     iteration:items     List the work items in a given iteration
     iteration:list      List all the iterations
     iteration:people    Show people based data for the iteration
   plans:
     plan:list           List all the delivery plans
   pull requests:
     pr:list             List all the pull requests
```


# Configuration

You will need the following environment variables defining:

```
$ export VSTS_ACCOUNT=""
$ export VSTS_PROJECT=""
$ export VSTS_TEAM=""
$ export VSTS_TOKEN=""
```

# Installation via Git

```
$ git clone git@github.com:benmatselby/donny.git
$ cd donny
$ make all
$ ./donny iteration:list
```

# Installation via Docker

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
