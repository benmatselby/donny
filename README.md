# Donny

![Go](https://github.com/benmatselby/donny/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/benmatselby/donny?style=flat-square)](https://goreportcard.com/report/github.com/benmatselby/donny)

_Forget it, Donny, you're out of your element!_

CLI application for getting information out of Azure DevOps. It's based on [Trello CLI](https://github.com/benmatselby/trello-cli) so the aims are the same:

```shell
CLI application for retrieving data from Azure DevOps

Usage:
  donny [command]

Available Commands:
  branch             Provide a list of pull requests
  builds             Provide a list of builds
  help               Help about any command
  iteration-burndown Provide a burndown of the iterations
  iteration-items    Provide a list of items in an iterations
  iteration-people   Provide a person breakdown for the iteration
  iterations         Provide a list of iterations
  plan               Get information about a delivery plan
  plans              Provide a list of delivery plans
  prs                Provide a list of pull requests
  teams              Provide a list of teams (Defaults to teams you are in)

Flags:
      --config string   config file (default is $HOME/.benmatselby/donny.yaml)
  -h, --help            help for donny

Use "donny [command] --help" for more information about a command.
```

## Requirements

If you are wanting to build and develop this, you will need the following items installed. If, however, you just want to run the application I recommend using the docker container (See below)

- Go version 1.16+

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
./donny builds
```

You can also install into your `$GOPATH/bin` by `go install`
