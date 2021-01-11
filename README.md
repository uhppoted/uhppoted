![build](https://github.com/uhppoted/uhppoted/workflows/build/badge.svg)

# uhppoted

`uhppoted` implements a set of cross-platform building blocks for access control systems based on the 
*UHPPOTE UT0311-L0x* TCP/IP Wiegand access control boards. Currently available:

- device API
- external application API
- CLI for scripting and system administration
- REST service for integration with HTTP servers and mobile clients
- MQTT endpoint for integration with IOT systems
- AWS S3 integration for file managed access control lists
- Google Sheets integration for spreadsheet managed access control lists
- [Node-RED](https://nodered.org) low code environment integration

Supported operating systems:
- Linux
- MacOS
- Windows
- ARM7 _(e.g. RaspberryPi)_

This project is a fork of the original [Go CLI](https://github.com/twystd/uhppote-go) project which had outgrown
its initial scope and was relocated to [uhppoted](https://github.com/uhppoted) to simplify future development.

### Compatible Hardware

As per this issue [[Question] Compatible Hardware](https://github.com/uhppoted/uhppote-core/issues/1), **UHPPOTE** appears to be a specific 
branding (or distributor) for the access control boards manufactured by [Shenzhen Wiegand Industrial Co., Ltd](http://www.wiegand.com.cn/english).

The software in this repository has been tested and is known to work with these specific boards:

| Source | Item |
| ------ | ---- |
| Amazon | [UHPPOTE Professional Wiegand 26-40 Bit TCP IP Network Access Control Board with Software For 4 Door 4 Reader](https://www.amazon.com/UHPPOTE-Professional-Wiegand-Network-Controller/dp/B00UX02JWE) |
| AliExpress | [TCP/IP RFID ACCESS CONTROL SYSTEM Wiegand 26](https://de.aliexpress.com/item/4000781912427.html) |

#### Firmware versions

| Version | Notes |
| ------- | ---- |
| 6.56    | Minimum firmware version (cf. [carbonsphere/UHPPOTE](https://github.com/carbonsphere/UHPPOTE) |
| 6.62    | Lowest firmware version in use |
| 8.92    | Latest tested firmware version |

_Notes:_

1. Firmware v6.62 sends anomalous _listen events_ with `0x19` as the start of message identifier. This appears to have
been fixed in later firmware versions but patches to support these events are included in:

- [`uhppote-core`](https://github.com/uhppoted/uhppote-core/blob/75a185a48184ecb68a07a09ebdd9ea1a8f96ba2c/encoding/UTO311-L0x/UT0311-L0x.go#L201-L204)
- [`uhppote-simulator`](https://github.com/uhppoted/uhppote-simulator/blob/f599512fb821c892a75786bbe4f35f6ebb4563d9/commands/run.go#L125-L134)
- [`node-red-contrib-uhppoted`](https://github.com/uhppoted/node-red-contrib-uhppoted/blob/74de32d62bee8097c03c9a1abc2bb45b0160f7b2/nodes/codec.js#L93-L100)

## Raison d'être

The components supplement the manufacturer supplied application which is 'Windows-only' and provides limited support 
for integration with other systems. 

The components are intended to simplify the integration of access control into systems based on:
- standard REST architectecture
- [AWS IoT](https://aws.amazon.com/iot)
- [Google Cloud IoT](https://cloud.google.com/solutions/iot)
- [IBM Watson IoT Platform](https://internetofthings.ibmcloud.com)
- [Node-RED](https://nodered.org) low code environment integration

## Releases

- v0.6.8: Adds handling for v6.62 firmware _listen events_ to `node-red-contrib-uhppoted`
- v0.6.7: Implements `record-special-events` to enable/disable door events
- v0.6.5: `node-red-contrib-uhppoted` module for use with NodeRED low code environment
- v0.6.4: `uhppoted-app-sheets` Google Sheets integration module
- v0.6.3: Added access control list commands to `uhppoted-mqtt`
- v0.6.2: Added access control list commands to `uhppoted-rest`
- v0.6.1: Added access control list commands to `uhppote-cli`
- v0.6.0: `uhppoted-app-s3` AWS S3 integration module
- v0.5.1: Initial release following restructuring into standalone Go *modules* and *git submodules*
- v0.5.0: Add MQTT endpoint for remote access to UT0311-L0x controllers
- v0.4.2: Reworked `GetDevice` REST API to use directed broadcast and added get-device to CLI
- v0.4.1: Get/set door control state functionality added to simulator, CLI and REST API
- v0.4.0: REST API service
- v0.3.1: Functional simulator with minimal command API
- v0.2.0: Load access control list from TSV file
- v0.1.0: Bare-bones but functional CLI

## Modules

| *Module*                                                                           | *Description*                                                                   |
| ---------------------------------------------------------------------------------- | ------------------------------------------------------------------------------- |
| [uhppote-core](https://github.com/uhppoted/uhppote-core)                           | core library, implements the UP interface to UT0311-L0x controllers             |
| [uhppoted-api](https://github.com/uhppoted/uhppoted-api)                           | common API for external applications                                            |
| [uhppote-simulator](https://github.com/uhppoted/uhppote-simulator)                 | UT0311-L04 simulator for development use                                        |
| [uhppote-cli](https://github.com/uhppoted/uhppote-cli)                             | command line interface                                                          |
| [uhppoted-rest](https://github.com/uhppoted/uhppoted-rest)                         | daemon/service with REST API for remote access to UT0311-L0x controllers        |
| [uhppoted-mqtt](https://github.com/uhppoted/uhppoted-mqtt)                         | daemon/service with MQTT endpoint for remote access to UT0311-L0x controllers   |
| [uhppoted-app-s3](https://github.com/uhppoted/uhppoted-app-s3)                     | cron'able utility to download/upload access control lists from/to AWS S3        |
| [uhppoted-app-sheets](https://github.com/uhppoted/uhppoted-app-sheets)             | cron'able utility to download/upload access control lists from/to Google Sheets |
| [node-red-contrib-uhppoted](https://github.com/uhppoted/node-red-contrib-uhppoted) | NodeJS nodes for [Node-RED](https://nodered.org) low code environment           |

## Integrations

| *Module*                                                                           | *Description*                                                                   |
| ---------------------------------------------------------------------------------- | ------------------------------------------------------------------------------- |
| [uhppoted-app-s3](https://github.com/uhppoted/uhppoted-app-s3)                     | cron'able utility to download/upload access control lists from/to AWS S3        |
| [uhppoted-app-sheets](https://github.com/uhppoted/uhppoted-app-sheets)             | cron'able utility to download/upload access control lists from/to Google Sheets |
| [node-red-contrib-uhppoted](https://github.com/uhppoted/node-red-contrib-uhppoted) | NodeJS nodes for [Node-RED](https://nodered.org) low code environment           |

## UI

| *Module*                                                               | *Description*                                                                   |
| ---------------------------------------------------------------------- | ------------------------------------------------------------------------------- |
| [uhppote-cli](https://github.com/uhppoted/uhppote-cli)                 | Command line interface                                                          |

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

| *Dependency*                                                                           | *Description*                                           |
| -------------------------------------------------------------------------------------- | ------------------------------------------------------- |
| [com.github/uhppoted/uhppote-core](https://github.com/uhppoted/uhppote-core)           | Device level API implementation                         |
| [com.github/uhppoted/uhppoted-api](https://github.com/uhppoted/uhppoted-api)           | External API implementation                             |
| [com.github/uhppoted/uhppote-cli](https://github.com/uhppoted/uhppote-cli)             | CLI user application                                    |
| [com.github/uhppoted/uhppoted-rest](https://github.com/uhppoted/uhppoted-rest)         | REST API                                                |
| [com.github/uhppoted/uhppoted-mqtt](https://github.com/uhppoted/uhppoted-mqtt)         | MQTT endpoint                                           |
| [com.github/uhppoted/uhppoted-app-s3](https://github.com/uhppoted/uhppoted-app-s3)     | cron'able access control management using AWS S3        |
| [com.github/uhppoted/uhppoted-app-sheets](https://github.com/uhppoted/uhppoted-app-s3) | cron'able access control management using Google Sheets |
| [com.github/uhppoted/uhppote-simulator](https://github.com/uhppoted/uhppote-simulator) | Device simulator for development use                    |
| [node-red-contrib-uhppoted](https://github.com/uhppoted/node-red-contrib-uhppoted)     | NodeJS nodes for [Node-RED](https://nodered.org) low code environment           |
| golang.org/x/sys/windows                                                               | Support for Windows services                            |
| golang.org/x/lint/golint                                                               | Additional *lint* check for release builds              |
| github.com/eclipse/paho.mqtt.golang                                                    | Eclipse Paho MQTT client                                |
| github.com/gorilla/websocket                                                           | paho.mqtt.golang dependency                             |
| github.com/aws/aws-sdk-go                                                              | AWS SDK                                                 |
| google.golang.org/api                                                                  | Google SDK                                              |
| golang.org/x/net                                                                       | google.golang.org/api dependency                        |
| golang.org/x/oauth2                                                                    | google.golang.org/api dependency                        |
| golang.org/x/tools                                                                     | google.golang.org/api dependency                        |
| golang.org/x/xerrors                                                                   | google.golang.org/api dependency                        |
| github.com/golang/protobuf                                                             | google.golang.org/api dependency                        |
| google.golang.org/appengine                                                            | google.golang.org/api dependency                        |

## References and Related Projects

1. [carbonsphere/UHPPOTE](https://github.com/carbonsphere/UHPPOTE) `PHP`
2. [carbonsphere/DoorControl](https://github.com/carbonsphere/DoorControl) `PHP`
2. [andrewvaughan/uhppote-rfid](https://github.com/andrewvaughan/uhppote-rfid) `Python`
3. [tachoknight/uhppote-tools](https://github.com/tachoknight/uhppote-tools): `Go`
4. [jjhuff/uhppote-go](https://github.com/jjhuff/uhppote-go): `Go`
5. [pawl/Chinese-RFID-Access-Control-Library](https://github.com/pawl/Chinese-RFID-Access-Control-Library)
6. [Dallas Makerspace:Reverse Engineering RFID Reader](https://dallasmakerspace.org/wiki/ReverseEngineeringRFIDReader)
7. [wemakerspace/wiegand-daemon](https://github.com/wemakerspace/wiegand-daemon)
8. [wemakerspace/weigeng-js](https://github.com/wemakerspace/weigeng-js)
9. [realashleybailey/DoorControl](https://github.com/realashleybailey/DoorControl)


