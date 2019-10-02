# Bring
[![Build Status](https://github.com/deluan/bring/workflows/CI/badge.svg)](https://github.com/deluan/bring/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/deluan/bring)](https://goreportcard.com/report/github.com/deluan/bring)
[![Documentation](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat)](https://godoc.org/github.com/deluan/bring) 

Go client library for [Apache Guacamole](http://guacamole.apache.org) Protocol.

## Quick start (tl;dr)

1. Install the library in your project: 
        
       go get github.com/deluan/bring

2. Create a [Session](https://godoc.org/github.com/deluan/bring#Session) with the `NewSession()` function.
3. Create a [Client](https://godoc.org/github.com/deluan/bring#Client) with the `NewClient()` function.
4. Start the client with `go client.Start()`
5. Get screen updates with `client.Screen()`
5. Send keystrokes with `client.SendKey()`
6. Send mouse updates with `client.SendMouse()`  

See the [sample app](sample/main.go) for a working example

## Documentation

The API is provided by the [Session](https://godoc.org/github.com/deluan/bring#Session) 
and the [Client](https://godoc.org/github.com/deluan/bring#Client) structs. 
The [documentation](https://godoc.org/github.com/deluan/bring) is a work in progress, 
but the API is very simple and you can take a look at all features available in the 
[sample app](sample) provided. Here are the steps to run the app:

1) You'll need a working `guacd` server in your machine. The easiest way is using docker 
and docker-compose. Just call `docker-compose up -d` in the root of this project. It 
starts the `guacd` server and a sample headless linux with a VNC server

2) Run the sample app with `make run`. It will connect to the linux container started by docker.

Take a look at the Makefile to learn how to run it in different scenarios.

## Why?

Apache Guacamole was created with the goal of making a dedicated client unnecessary. 
So why create a client?!

The idea is that if you need to control a remote machine from your Go code, you can 
leverage the Guacamole protocol and the `guacd` server as a bridge. This way you can 
use any protocol supported by Guacamole (currently RDP and VNC, with X11 coming in 
the future) to do screen capture and remote control of networked servers/desktop 
machines from within your Go app.

My use case was to automate some tasks in a VirtualBox VM, but there was no Go support 
for the VirtualBox XPCOM API on Macs (my host platform), nor a working RDP client 
implementation in Go. Instead of writing a new RDP client, why not leverage the awesome 
Guacamole project and get support for multiple protocols?

## References:
- [The Guacamole protocol](http://guacamole.apache.org/doc/gug/guacamole-protocol.html)
- [Guacamole protocol reference](http://guacamole.apache.org/doc/gug/protocol-reference.html#rect-instruction)
- [Apache Guacamole Client implementation](https://github.com/apache/guacamole-client/tree/master/guacamole-common-js)
