# Cron Expression Parser

This is a command line application which parses a cron string and expands each field to show the times at which it will run.

The current implementation supports the cron format of 5 time fields in the order minute, hour, day of month, month and day of week) plus a command.

The cron string has to be passed as a single argument on a single line.

It is written in golang 1.13.

## Installation
The program requires golang 1.13 or greater. Once the repository is cloned to build the build run `make` command from the root of the cloned repository. The commands will generate an executable called `cep`.

Once it is installed you can invoke the command e.g.:
```
./cep "*/15 0 1,15 * 1-5 /bin/ls"

```

The program runs on OSX and linux.

## Docker installation
If you have installed docker it is possible to run the program in a container, please build the container with `make docker-build` and the you can run invoking the command via container e.g.:
```
 docker run -it --rm cronparserexpander  "*/15 0 1,15 * 1-5 /bin/ls"

```
