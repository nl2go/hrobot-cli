# hrobot-cli: Command-line interface for Hetzner Robot Webservice

![Build](https://gitlab.com/newsletter2go/hrobot-cli/badges/master/pipeline.svg) ![Coverage](https://gitlab.com/newsletter2go/hrobot-cli/badges/master/coverage.svg)

`hrobot-cli` is a command-line interface for interacting with the [Hetzner Robot API](https://robot.your-server.de/doc/webservice/en.html).

## Set environment variables for robot webservice credentials

`hrobot-cli` use environment variables for authenticating with the Hetzner Robot. It is
recommended to use an env file similar to `sample.env` which can be used in the local environment as
well as with the docker container. The expected variables are
* HROBOTCLI_USER
* HROBOTCLI_PASSWORD

## Build manually and run on local machine

If you have Go installed, you can build `hrobot-cli` with:

    go build

and run it with:

    export $(cat sample.env | xargs) # or set the expected environemnt variables in any other way
    ./hrobot-cli

## Build and run in docker

The project contains a Dockerfile which will build a container with the binary by running:

    docker build --rm -f "Dockerfile" -t hrobot-cli:latest .

Afterwards you can run it with:

    docker run -it --env-file=sample.env hrobot-cli:latest

## Run latest version in docker using pre-built docker image

Run latest version from docker registry:

    docker run -it --env-file=sample.env registry.github.com/nl2go/hrobot-cli:v0-1-1

## Update hrobot-go mocks

This project uses `gomock`. When the hrobot-go library is upgraded to a new version, the mocks used 
for the tests need to be updated. This can be done by running the following command:

    mockgen -package mock github.com/nl2go/hrobot-go RobotClient > test/mock/mock_robot_client.go

## Features & overview

Currently implemented commands are:

```
CLI application for the Hetzner Robot API - version 0.1.1

Usage:
  hrobot-cli [command]

Available Commands:
  help               Help about any command
  ip:list            Print list of IP's
  key:list           Print list of ssh keys
  rdns:get           Print single reverse DNS entry
  rdns:list          Print list of reverse DNS entries
  server:ansible-inv Generates ansible inventory from server list
  server:get         Print single server
  server:list        Print list of servers
  server:rescue      Activate rescue mode for single server
  server:reset       Reset single server (hardware reset)
  server:reverse     Revert single server order
  server:set-name    Sets name for selected servers
  version            Print the version number of hrobot-cli

Flags:
  -h, --help   help for hrobot-cli

Use "hrobot-cli [command] --help" for more information about a command.
```

For commands that are applied to a single server it is possible to select the server interactively 
by searching through the server list. For other commands (like server renaming) multiple items can
be selected for executing the respective command.
