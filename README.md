# uhppote-go

Go CLI implementation for the UHPPOTE UT0311-L04 TCP/IP Wiegand Access Control Board. The current incarnation is essentially a rework in Go of the [carbonsphere/UHPPOTE](https://github.com/carbonsphere/UHPPOTE) PHP implementation.

## Raison d'être

The manufacturer supplied software for the UHPPOTE UT0311-L04 TCP/IP Wiegand Access Control Board is 'Windows only' and is also not suitable for server use or integration with other applications.

## Status

*Under development*

## Modules

- uhppote-cli:       CLI for use with bash scripts
- uhppote-simualtor: UHPPOTE simulator for development use

## Installation

## uhppote-cli

Usage: *uhppote-cli \<command\> \<command options\>*

Supported commands:
- search
- get-time \<serial number\>
- set-time \<serial number\> --now

## uhppote-simulator

Usage: *uhppote-simulator*

Supported options:
- --help
- --version
- --debug








