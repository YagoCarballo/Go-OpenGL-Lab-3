## Build the Project

##### Requeriments

- [Gigo](https://github.com/LyricalSecurity/gigo)

##### How to Install

###### Install Gigo

```bash

## Clone Gigo's Repo
$ git clone https://github.com/LyricalSecurity/gigo.git && cd gigo

## Set Go's Path
$ export GOPATH=`pwd`

## Get Gigo's Dependencies
$ go get github.com/LyricalSecurity/gigo/actions

## Build Gigo
$ go build -o dist/gigo main.go

## Move Gigo to Go's Bin Directory
$ mv dist/gigo /usr/local/go/bin/gigo

## Add Go's Bin directory to the PATH (If not added already)


```