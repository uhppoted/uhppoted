# uhppote-go

Go CLI and daemon/service implementation for the UHPPOTE UT0311-L04 TCP/IP Wiegand Access Control Board. Based largely on the [carbonsphere/UHPPOTE](https://github.com/carbonsphere/UHPPOTE) PHP implementation.

## Raison d'être

Provide a cross-platform base for access control systems based on the UHPPOTE UT0311-L04 TCP/IP Wiegand Access Control Board. The manufacturer supplied software is a 'Windows only' application which is not ideal for server use or integration with other systems.

## Releases

- v0.4.0: REST API service
- v0.3.1: functional simulator with minimal command API
- v0.2.0: load access control list from TSV file
- v0.1.0: bare-bones but functional CLI

## Modules

| Module            | Description                                                              |
| ----------------- | ------------------------------------------------------------------------ |
| uhppote           | core library, implements the UDP interface to UT0311-L0x controllers     |
| uhppote-cli       | command line interface                                                   |
| uhppoted-rest     | daemon/service with REST API for remote access to UT0311-L0x controllers |
| uhppote-simulator | UT0311-L04 simulator for development use                                 |

## Installation

### Building from source

#### Dependencies

- golang.org.x.sys (for uhppoted Windows service)

### Binaries

## uhppote

Supported functions:
- FindDevices
- SetAddress
- GetStatus
- GetTime
- SetTime
- GetDoorDelay
- SetDoorDelay
- GetListener
- SetListener
- GetCards
- GetCardByIndex
- GetCardById
- PutCard
- DeleteCard
- GetEvent
- GetEventIndex
- SetEventIndex
- OpenDoor
- Listen

## uhppote-cli

Usage: *uhppote-cli [--bind <address:port>] [--debug] \<command\> \<arguments\>*

Supported commands:

- help
- version
- get-devices
- set-address
- get-status
- get-time
- set-time
- get-door-delay
- set-door-delay
- get-listener
- set-listener
- get-cards
- get-card
- get-events
- get-swipe-index
- set-event-index
- open
- grant
- revoke
- revoke-all
- load-acl
- listen

## uhppoted-rest

Usage: *uhppoted-rest \<command\> \<options\>*

Defaults to 'run' unless one of the commands below is specified: 

- daemonize
- undaemonize
- help
- version

## uhppote-simulator

Usage: *uhppote-simulator --devices=\<dir\>*

Supported options:
- --help
- --version
- --debug








