# uhppoted

`uhppoted` implements a set of cross-platform building blocks for access control systems based on the 
*UHPPOTE UT0311-L0x* TCP/IP Wiegand access control boards. Currently available:

- device API
- external application API
- CLI for scripting and system administration
- REST service for integration with HTTP servers and mobile clients
- MQTT endpoint for integration with IOT systems

Supported operating systems:
- Linux
- MacOS
- Windows

This project is a fork of the original [Go CLI](https://github.com/twystd/uhppote-go) project which had outgrown
its initial scope and was relocated to [uhppoted](https://github.com/uhppoted) to simplify future development.

## Raison d'être

The components supplement the manufacturer supplied application which is 'Windows-only' and provides limited support 
for integration with other systems. 

The components are intended to simplify the integration of access control into systems based on:
- standard REST architectecture
- [AWS IoT](https://aws.amazon.com/iot)
- [Google Cloud IoT](https://cloud.google.com/solutions/iot)
- [IBM Watson IoT Platform](https://internetofthings.ibmcloud.com)

## Releases

- v0.5.1: Initial release following restructuring into standalone Go *modules* and *git submodules*
- v0.5.0: Add MQTT endpoint for remote access to UT0311-L0x controllers
- v0.4.2: Reworked `GetDevice` REST API to use directed broadcast and added get-device to CLI
- v0.4.1: Get/set door control state functionality added to simulator, CLI and REST API
- v0.4.0: REST API service
- v0.3.1: Functional simulator with minimal command API
- v0.2.0: Load access control list from TSV file
- v0.1.0: Bare-bones but functional CLI

## Modules

| *Module*               | *Description*                                                                 |
| ---------------------- | ----------------------------------------------------------------------------- |
| [uhppote-core][1]      | core library, implements the UDP interface to UT0311-L0x controllers          |
| [uhppoted-api][2]      | common API for external applications                                          |
| [uhppote-simulator][3] | UT0311-L04 simulator for development use                                      |
| [uhppote-cli][4]       | command line interface                                                        |
| [uhppoted-rest][5]     | daemon/service with REST API for remote access to UT0311-L0x controllers      |
| [uhppoted-mqtt][6]     | daemon/service with MQTT endpoint for remote access to UT0311-L0x controllers |

## Installation

Binaries for Linux, Windows, MacOS and Raspbian/ARM7 are distributed in the tarball for each release. To install
the binaries, download and extract the tarball to a directory of your choice.

### Building from source

*uhppoted* is the parent project for the individual components which are referenced as git submodules -
to clone the entire source tree:

```
git clone --recurse-submodules https://github.com/uhppoted/uhppoted.git

```

The supplied `Makefile` has targets to build binaries for all the supported operating systems:
```
make build
```
or 
```
make release
```

To pull upstream changes for all submodules:

```
git submodule update --remote
```

#### Dependencies

| *Dependency*                             | *Description*                                          |
| ---------------------------------------- | ------------------------------------------------------ |
| com.github/uhppoted/uhppote-core[1]      | Device level API implementation                        |
| com.github/uhppoted/uhppoted-api[2]      | External API implementation                            |
| com.github/uhppoted/uhppote-cli[4]       | CLI user application                                   |
| com.github/uhppoted/uhppoted-rest[5]     | REST API                                               |
| com.github/uhppoted/uhppoted-mqtt[6]     | MQTT endpoint                                          |
| com.github/uhppoted/uhppote-simulator[3] | Device simulator for development use                   |
| golang.org/x/sys/windows                 | Support for Windows services                           |
| golang.org/x/lint/golint                 | Additional *lint* check for release builds             |
| github.com/eclipse/paho.mqtt.golang      | Eclipse Paho MQTT client                               |
| github.com/gorilla/websocket             | paho.mqtt.golang dependency                            |

## References and Related Projects

1. [carbonsphere/UHPPOTE](https://github.com/carbonsphere/UHPPOTE) `PHP`
2. [carbonsphere/DoorControl](https://github.com/carbonsphere/DoorControl) `PHP`
2. [andrewvaughan/uhppote-rfid](https://github.com/andrewvaughan/uhppote-rfid) `Python`
3. [tachoknight/uhppote-tools](https://github.com/tachoknight/uhppote-tools): `Go`
4. [jjhuff/uhppote-go](https://github.com/jjhuff/uhppote-go): `Go`

[1]: https://github.com/uhppoted/uhppote-core
[2]: https://github.com/uhppoted/uhppoted-api
[3]: https://github.com/uhppoted/uhppote-simulator
[4]: https://github.com/uhppoted/uhppote-cli
[5]: https://github.com/uhppoted/uhppoted-rest
[6]: https://github.com/uhppoted/uhppoted-mqtt




