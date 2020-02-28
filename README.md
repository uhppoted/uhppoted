# uhppote-go

A set of cross-platform building blocks for access control systems based on the *UHPPOTE UT0311-L0x* TCP/IP 
Wiegand access control boards. Currently available:

- low level communication library
- CLI for scripting and system administration
- REST service for integration with HTTP servers and mobile clients
- MQTT endpoint for integration with IOT systems*

Supported operating systems:
- Linux
- MacOS
- Windows

## Raison d'être

The manufacturer supplied application is 'Windows-only' and provides limited support for integration with other
systems.

## Releases

- v0.5.0: Add MQTT endpoint for remote access to UT0311-L0x controllers
- v0.4.2: Reworked `GetDevice` REST API to use directed broadcast and added get-device to CLI
- v0.4.1: Get/set door control state functionality added to simulator, CLI and REST API
- v0.4.0: REST API service
- v0.3.1: Functional simulator with minimal command API
- v0.2.0: Load access control list from TSV file
- v0.1.0: Bare-bones but functional CLI

## Modules

| *Module*          | *Description*                                                                 |
| ----------------- | ----------------------------------------------------------------------------- |
| uhppote           | core library, implements the UDP interface to UT0311-L0x controllers          |
| uhppote-simulator | UT0311-L04 simulator for development use                                      |
| uhppote-cli       | command line interface                                                        |
| uhppoted-rest     | daemon/service with REST API for remote access to UT0311-L0x controllers      |
| uhppoted-mqtt     | daemon/service with MQTT endpoint for remote access to UT0311-L0x controllers |

## Installation

### Building from source

#### Dependencies

| *Dependency*                          | *Description*                                          |
| ------------------------------------- | ------------------------------------------------------ |
| com.github/uhppoted/uhppote-core      | Device level API implementation                        |
| com.github/uhppoted/uhppoted-api      | External API implementation                            |
| com.github/uhppoted/uhppote-cli       | CLI user application                                   |
| com.github/uhppoted/uhppoted-rest     | REST API                                               |
| com.github/uhppoted/uhppoted-mqtt     | MQTT endpoint                                          |
| com.github/uhppoted/uhppote-simulator | Device simulator for development use                   |
| golang.org/x/sys/windows              | Support for Windows services                           |
| golang.org/x/lint/golint              | Additional *lint* check for release builds             |
| github.com/eclipse/paho.mqtt.golang   | Eclipse Paho MQTT client                               |
| github.com/gorilla/websocket          | paho.mqtt.golang dependency                            |

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

Supported 'run' options:
- --console
- --debug

## uhppoted-mqtt

Usage: *uhppoted-mqtt \<command\> \<options\>*

Defaults to 'run' unless one of the commands below is specified: 

- daemonize
- undaemonize
- config
- help
- version

Supported 'run' options:
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






