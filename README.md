# uhppote-go

Go CLI implementation for the UHPPOTE UT0311-L04 TCP/IP Wiegand Access Control Board. The current incarnation is essentially a rework in Go of the [carbonsphere/UHPPOTE](https://github.com/carbonsphere/UHPPOTE) PHP implementation.

## Raison d'être

The manufacturer supplied software for the UHPPOTE UT0311-L04 TCP/IP Wiegand Access Control Board is 'Windows only' and is also not suitable for server use or integration with other applications.

## Releases

- v0.03: functional simulator with minimal command API
- v0.02: load access control list from TSV file
- v0.01: bare-bones but functional CLI

## Modules

- uhppote:           UDP interface to UHPPOTE UTO0311-L0x controllers
- uhppote-cli:       CLI for use with bash scripts
- uhppote-simulator: UHPPOTE simulator for development use

## Installation

## uhppote

Supported functions:
- FindDevices
- SetAddress
- GetTime
- SetTime
- GetDoorDelay
- SetDoorDelay
- GetListener
- SetListener
- GetStatus
- GetCards
- GetCardByIndex
- GetCardById
- PutCard
- DeleteCard
- GetEvent
- GetEventIndex
- SetEventIndex
- Open
- Listen

## uhppote-cli

Usage: *uhppote-cli [--bind <address:port>] [--debug] \<command\> \<arguments\>*

Supported commands:

- help
- version
- get-devices
- set-address
- get-time
- set-time
- get-door-delay
- set-door-delay
- get-listener
- set-listener
- get-status
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

## uhppote-simulator

Usage: *uhppote-simulator --devices=<dir>*

Supported options:
- --help
- --version
- --debug








