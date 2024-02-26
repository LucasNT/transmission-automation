# Description

This program is to solve two of my pains, which is to:

- save and rename the torrent downloaded by transmission
- learn Golang


# Tech Stack

This programn is written in Golang, at the moment the programn has 2 builds one with a cli interface and another with a web interface. both use the default library of golang.

# Features

Currentely the program needs to be refactored, the core is to rigid to allow a dynamic UI.

# How to run

Just download one version from the releases, and execute it.

# How to build

To build the programn with web interface run:
```bash
go build ./cmd/web/web.go
```

To build the programn with cli run:
```bash
go build ./cmd/main/transmission-list-automation.go
```
