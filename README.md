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

## uhppote-cli

Usage: *uhppote-cli [--bind <address:port>] [--debug] \<command\> \<command options\>*

Supported commands:
- help
- version
- get-devices
- get-status
- get-time
- get-cards
- get-swipes
- set-time
- set-ip-address
- grant
- revoke
- revoke-all
- open-door

## uhppote-simulator

Usage: *uhppote-simulator*

Supported options:
- --help
- --version
- --debug








