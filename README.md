# Bring
[![Build Status](https://github.com/deluan/bring/workflows/CI/badge.svg)](https://github.com/deluan/bring/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/deluan/bring)](https://goreportcard.com/report/github.com/deluan/bring)
[![Documentation](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat)](https://godoc.org/github.com/deluan/bring) 

Go Client for [Apache Guacamole](http://guacamole.apache.org) Protocol.

## Why?
Apache Guacamole was created with the goal of making a dedicated client unnecessary. So why create a client?!

The idea is that if you need to control a remote machine from your Go code, you can leverage the Guacamole protocol and the `guacd` server as a bridge. This way you can use any protocol supported by Guacamole (currently RDP and VNC, with X11 coming in the future) to do screen capture and remote control of networked servers/desktop machines.

My use case was to automate some tasks in a VirtualBox VM, but there was no functional support for the VirtualBox XPCOM API on Macs (my host platform), nor a working RDP implementation.

## Documentation

The API is provided by the [Client](client.go) struct. The [documentation](https://godoc.org/github.com/deluan/bring) is a work in progress, but the API is very simple and you can take a look at all features available in the [sample app](app) provided. Here are the steps to run the app:

1) You'll need a working `guacd` server in your machine. The easiest way is using docker and docker-compose. Just call `docker-compose up -d` in the root of this project. It starts the `guacd` server and a sample headless linux with a VNC server

2) Run the sample app with `make run`. It will connect to the linux container started by docker.

Take a look at the Makefile to learn how to run it in different scenarios. Keep in mind that this sample has a hardcoded resolution of 1024x768

## References:
- [The Guacamole protocol](http://guacamole.apache.org/doc/gug/guacamole-protocol.html)
- [Guacamole protocol reference](http://guacamole.apache.org/doc/gug/protocol-reference.html#rect-instruction)
- [Apache Guacamole Client implementation](https://github.com/apache/guacamole-client/tree/master/guacamole-common-js)
