Donny
=====

[![Build Status](https://travis-ci.org/benmatselby/donny.png?branch=master)](https://travis-ci.org/benmatselby/donny)

_Forget it, Donny, you're out of your element!_

CLI application for getting information out of Visual Studio Team Services. It's based on [Trello CLI](https://github.com/benmatselby/trello-cli) so the aims are the same:

* Provide a list of open iterations/boards
* Provide details about which user stories are in the sprint


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
