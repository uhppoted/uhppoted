# uhppote-go

Go CLI implementation for the UHPPOTE UT0311-L04 TCP/IP Wiegand Access Control Board. The current incarnation is essentially a rework in Go of the [carbonsphere/UHPPOTE](https://github.com/carbonsphere/UHPPOTE) PHP implementation.

## Raison d'être

The manufacturer supplied software for the UHPPOTE UT0311-L04 TCP/IP Wiegand Access Control Board is 'Windows only' and is also not suitable for server use or integration with other applications.

## Status

*Under development*

## Modules

- uhppote-cli:       CLI for use with bash scripts
- uhppote-simulator: UHPPOTE simulator for development use

## Installation

## uhppote

Supported functions:
- FindDevices
- SetAddress
- GetStatus
- GetTime
- SetTime
- GetDoorDelay
- SetDoorDelay
- GetCards
- GetCardByIndex
- GetCardById
- GetEvent
- GetEventIndex
- SetEventIndex

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
- get-cards
- get-card
- grant
- revoke
- revoke-all
- get-events
- get-swipe-index
- set-event-index
- open             

## uhppote-simulator

Usage: *uhppote-simulator*

Supported options:
- --help
- --version
- --debug








