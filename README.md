# uhppote-go

Go CLI and daemon/service implementation for the UHPPOTE UT0311-L04 TCP/IP Wiegand Access Control Board. 

## Raison d'être

Provides a set of cross-platform building blocks for access control systems based on the UHPPOTE UT0311-L04 TCP/IP 
Wiegand Access Control Board, to supplement the 'Windows-only' out-of-the-box application supplied with the boards.

## Releases

- v0.4.2: Reworked `GetDevice` REST API to use directed broadcast and added get-device to CLI
- v0.4.1: Get/set door control state functionality added to simulator, CLI and REST API
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
- FindDevice
- SetAddress
- GetStatus
- GetTime
- SetTime
- GetDoorControlState
- SetDoorControlState
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
- get-device
- set-address
- get-status
- get-time
- set-time
- get-door-delay
- set-door-delay
- get-door-control
- set-door-control
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

Supported options:
- --console
- --debug


## uhppote-simulator

Usage: *uhppote-simulator \<command\> --devices=\<dir\>*

Defaults to 'run' unless one of the commands below is specified: 

- help
- version

Supported options:
- --bind <IP address to bind to>
- --devices <directory path for device files>
- --debug

## References and Related Projects

1. [carbonsphere/UHPPOTE](https://github.com/carbonsphere/UHPPOTE) `PHP`
2. [carbonsphere/DoorControl](https://github.com/carbonsphere/DoorControl) `PHP`
2. [andrewvaughan/uhppote-rfid](https://github.com/andrewvaughan/uhppote-rfid) `Python`
3. [tachoknight/uhppote-tools](https://github.com/tachoknight/uhppote-tools): `Go`
4. [jjhuff/uhppote-go](https://github.com/jjhuff/uhppote-go): `Go`






