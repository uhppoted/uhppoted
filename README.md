![build](https://github.com/uhppoted/uhppoted/workflows/build/badge.svg)

# uhppoted [(juː hɒp ɪt eɪd)](/internal/uhppoted.mp3)

`uhppoted` implements a set of cross-platform building blocks for access control systems based on the 
*UHPPOTE UT0311-L0x* TCP/IP Wiegand access control boards. Currently available components include:

- controller API
- CLI for scripting and system administration
- REST service for integration with HTTP servers and mobile clients
- MQTT endpoint for integration with IOT systems e.g.
  - [Eclipse Mosquitto](https://mosquitto.org)
  - [HiveMQ](https://www.hivemq.com/)
  - [AWS IoT](https://aws.amazon.com/iot)
  - [Google Cloud IoT](https://cloud.google.com/solutions/iot)
  - [IBM Watson IoT Platform](https://internetofthings.ibmcloud.com)
- AWS S3 integration for file managed access control lists
- Google Sheets integration for spreadsheet managed access control lists
- Integration with the [Wild Apricot](https://www.wildapricot.com) member management system
- [uhppoted-app-db](https://github.com/uhppoted/uhppoted-app-db) SQL DB integration for access control lists managed in a database
- [uhppoted-lib-nodejs](https://github.com/uhppoted/uhppoted-lib-nodejs) NodeJS package
- [uhppoted-lib-python](https://github.com/uhppoted/uhppoted-lib-python) Python PyPI package
- [uhppoted-lib-dotnet](https://github.com/uhppoted/uhppoted-lib-dotnet) .NET NuGet package
- [uhppoted-lib-go](https://github.com/uhppoted/uhppoted-lib-go) standalone Go library for external use
- [Node-RED](https://nodered.org) low code environment integration
- [Home Assistant](https://www.home-assistant.io/) custom component

The components supplement the manufacturer supplied application which is 'Windows-only' and provides limited support 
for integration with other systems. 

Hardware:
- [uhppoted-breakout](https://github.com/uhppoted/uhppoted-breakout), a customisable/hackable replacement CPU board
- [uhppoted-wiegand](https://github.com/uhppoted/uhppoted-wiegand) an M5Stack form factor Wiegand reader/emulator

_Note_:

FYA, after much discussion it has been decided that _uhppoted_ is pronounced _u-hop-id-ay-d_ (juː hɒp ɪt eɪd).

### Operating systems

Supported operating systems:
- Linux
- MacOS
- Windows
- RaspberryPi (ARM/ARM7/ARM6)

### Language Bindings

The primary development language for this repository is Go, but code generated native bindings for the following languages
also exist:
- [NodeJS](https://github.com/uhppoted/uhppoted-nodejs)
- [Python](https://github.com/uhppoted/uhppoted-lib-python)
- [Rust](https://github.com/uhppoted/uhppoted-codegen/tree/main/bindings/rust)
- [Zig](https://github.com/uhppoted/uhppoted-codegen/tree/main/bindings/zig)
- [Erlang](https://github.com/uhppoted/uhppoted-codegen/tree/main/bindings/erlang)
- [PHP](https://github.com/uhppoted/uhppoted-codegen/tree/main/bindings/php)
- [Lua](https://github.com/uhppoted/uhppoted-codegen/tree/main/bindings/lua)


### shared-lib/dylib/DLL

[uhppoted-dll](https://github.com/uhppoted/uhppoted-dll/blob/master/doc/C.md) implements a shared-lib/DLL for interop with languages other than Go. The implementation includes
bindings to:

- C
- C++
- C#
- Python
- Clozure Common Lisp

### 3rd party integrations

- [ioBroker](https://www.iobroker.net) - [kBrausew/ioBroker.wiegand-tcpip](https://github.com/kBrausew/ioBroker.wiegand-tcpip)
- [Andre Becker](https://github.com/Andre1Becker) has written a short guide to getting up and running with
  [_uhppoted-mqtt_](https://github.com/uhppoted/uhppoted/discussions/46#discussioncomment-7886167) and _OpenHAB_ using
  _Docker_.

### Compatible Hardware

As per this issue [[Question] Compatible Hardware](https://github.com/uhppoted/uhppote-core/issues/1), 
**UHPPOTE** appears to be a specific branding (or distributor) for the access control boards manufactured
by [Shenzhen Wiegand Industrial Co., Ltd](http://www.wiegand.com.cn/english).

The software in this repository has been tested and is known to work with these specific boards:

| Source | Item |
| ------ | ---- |
| Amazon | [UHPPOTE Professional Wiegand 26-40 Bit TCP IP Network Access Control Board with Software For 4 Door 4 Reader](https://www.amazon.com/UHPPOTE-Professional-Wiegand-Network-Controller/dp/B00UX02JWE) |
| AliExpress | [TCP/IP RFID ACCESS CONTROL SYSTEM Wiegand 26](https://de.aliexpress.com/item/4000781912427.html) |

#### Firmware versions

| Version | Notes |
| ------- | ----- |
| 6.56    | Minimum firmware version (cf. [carbonsphere/UHPPOTE](https://github.com/carbonsphere/UHPPOTE) |
| 6.62    | Lowest firmware version in use |
| 8.92    | Latest tested firmware version |

_Notes:_

1. Firmware v6.62 sends anomalous _listen events_ with `0x19` as the start of message identifier. This appears to have
been fixed in later firmware versions but patches to support these events are included in:

- [`uhppote-core`](https://github.com/uhppoted/uhppote-core/blob/master/encoding/UTO311-L0x/UT0311-L0x.go#L234-L237)
- [`node-red-contrib-uhppoted`](https://github.com/uhppoted/node-red-contrib-uhppoted/blob/main/nodes/codec.js#L107-L109)
- [`uhppoted-nodejs`](https://github.com/uhppoted/uhppoted-nodejs/blob/master/src/codec.js#L114-L121)
- [`uhppoted-codegen`](https://github.com/uhppoted/uhppoted-codegen/blob/main/bindings/go/uhppote/decode.go#L36-L40)
- [`uhppoted-lib-python`](https://github.com/uhppoted/uhppoted-lib-python/blob/main/src/uhppoted/decode.py#L63-L65)
- [`uhppoted-lib-dotnet`](https://github.com/uhppoted/uhppoted-lib-dotnet/blob/main/uhppoted/uhppoted/decode.fs#L466-L467)
- [`uhppoted-lib-go`](https://github.com/uhppoted/uhppoted-lib-go)

#### Readers

Almost any reader with a Wiegand-26 interface should probably work (there have been reports of offbrand readers that don't)
but the readers below are in active use:

| Reader | Notes       |
| ------ | ----------- |
| [HID ProxPoint Plus 6500](https://www.hidglobal.com/products/readers/hid-proximity/6005) | Old stock and/or refurbished readers are often available on ebay |
| [IKeys barcode and QR code scanner](https://www.i-keys.de/de/zutrittskontrollsysteme/qr-code-zugangskontrolle-mit-rfid-reader-iso14443.html?c=91) | _Ref. [Put Card - Date with Time](https://github.com/uhppoted/uhppoted/issues/2)_ |
| [Newland FM430L barcode and QR code scanner](https://www.ico.de/newland-fm430l-u-barracuda-2d-fixmount-usb-kabel-wechselbar--bcpf6u?referer=froogle&gclid=Cj0KCQjwkIGKBhCxARIsAINMioJ4fxcqv9W3dINBdnEwEUjRJIBHhXoQfJMO5wL4uKlf902Gm8--FBkaAi3vEALw_wcB) | _Ref. [Put Card - Date with Time](https://github.com/uhppoted/uhppoted/issues/2)_ |

## Release Notes

### Current Release

**[v0.8.11](https://github.com/uhppoted/uhppoted-wiegand/releases/tag/v0.8.11) - 2025-07-01**

1. Added support `get/set-anti-passback` API function.
2. Added _decorated events_ and caching to _uhppoted-app-home-assistant.
3. Added M5 stack Wiegand emulator (in progress).
4. Renamed _upppoted-nodejs_ repository to _uhppoted-lib-nodejs_.
5. Added check to prevent UDP bind address from using broadcast port.
6. Updated to Go 1.24 throughout.


## Modules

| *Module*                                                                               | *Description*                                                                   |
| ---------------------------------------------------------------------------------------| --------------------------------------------------------------------------------|
| [uhppote-core](https://github.com/uhppoted/uhppote-core)                               | core library, implements the UP interface to UT0311-L0x controllers |
| [uhppoted-lib](https://github.com/uhppoted/uhppoted-lib)                               | common utility function library |
| [uhppote-simulator](https://github.com/uhppoted/uhppote-simulator)                     | UT0311-L04 simulator for development use |
| [uhppote-cli](https://github.com/uhppoted/uhppote-cli)                                 | command line interface  |
| [uhppoted-rest](https://github.com/uhppoted/uhppoted-rest)                             | daemon/service with REST API for remote access to UT0311-L0x controllers |
| [uhppoted-mqtt](https://github.com/uhppoted/uhppoted-mqtt)                             | daemon/service with MQTT endpoint for remote access to UT0311-L0x controllers |
| [uhppoted-app-s3](https://github.com/uhppoted/uhppoted-app-s3)                         | cron'able utility to download/upload access control lists from/to AWS S3 |
| [uhppoted-app-sheets](https://github.com/uhppoted/uhppoted-app-sheets)                 | cron'able utility to download/upload access control lists from/to Google Sheets|
| [uhppoted-app-db](https://github.com/uhppoted/uhppoted-app-db)                         | cron'able utility to download/upload access control lists from/to a SQL DB |
| [uhppoted-app-wild-apricot](https://github.com/uhppoted/uhppoted-app-wild-apricot)     | cron'able utility to manage access control lists from Wild Apricot |
| [uhppoted-app-home-assistant](https://github.com/uhppoted/uhppoted-app-home-assistant) | Home Assistant custom component |
| [node-red-contrib-uhppoted](https://github.com/uhppoted/node-red-contrib-uhppoted)     | NodeJS nodes for [Node-RED](https://nodered.org) low code environment |
| [uhppoted-nodejs](https://github.com/uhppoted/uhppoted-nodejs)                         | Standalone NodeJS module |
| [uhppoted-lib-python](https://github.com/uhppoted/uhppoted-lib-python)                         | Standalone Python package for PyPI |
| [uhppoted-lib-dotnet](https://github.com/uhppoted/uhppoted-lib-dotnet)                 | .NET package for NuGet |
| [uhppoted-lib-go](https://github.com/uhppoted/uhppoted-lib-go)                         | Standalone Go library (for external use) |
| [uhppoted-dll](https://github.com/uhppoted/uhppoted-dll)                               | shared-lib/DLL for cross-language interop |
| [uhppoted-codegen](https://github.com/uhppoted/uhppoted-codegen)                       | Model based code generator for native language bindings |
| [uhppoted-tunnel](https://github.com/uhppoted/uhppoted-tunnel)                         | UDP tunnel for remote access |

## Integrations

| *Module*                                                                               | *Description*                                                                   |
| -------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------- |
| [uhppoted-app-db](https://github.com/uhppoted/uhppoted-app-db)                         | cron'able utility to download/upload access control lists from/to a SQL DB |
| [uhppoted-app-s3](https://github.com/uhppoted/uhppoted-app-s3)                         | cron'able utility to download/upload access control lists from/to AWS S3 |
| [uhppoted-app-sheets](https://github.com/uhppoted/uhppoted-app-sheets)                 | cron'able utility to download/upload access control lists from/to Google Sheets |
| [uhppoted-app-wild-apricot](https://github.com/uhppoted/uhppoted-app-wild-apricot)     | cron'able utility to manage access control lists from/t Wild Apricot  |
| [uhppoted-app-home-assistant](https://github.com/uhppoted/uhppoted-app-home-assistant) | Home Assistant custom component                                       |
| [node-red-contrib-uhppoted](https://github.com/uhppoted/node-red-contrib-uhppoted)     | NodeJS nodes for [Node-RED](https://nodered.org) low code environment |
| [kBrausew/ioBroker.wiegand-tcpip](https://github.com/kBrausew/ioBroker.wiegand-tcpip)  | [ioBroker](https://www.iobroker.net)  |


## UI

| *Module*                                                               | *Description*            |
| ---------------------------------------------------------------------- | -------------------------|
| [uhppote-cli](https://github.com/uhppoted/uhppote-cli)                 | Command line interface   |
| [uhppoted-httpd](https://github.com/uhppoted/uhppote-httpd)            | HTML browser interface   |

## Hardware

| *Module*                                                          | *Description*                               |
| ----------------------------------------------------------------- | --------------------------------------------|
| [uhppoted-wiegand](https://github.com/uhppoted/uhppoted-wiegand)  | M5Stack form factor Wiegand reader/emulator |
| [uhppoted-breakout](https://github.com/uhppoted/uhppote-breakout) | Customisable replacement CPU board          |


## Installation

Binaries for Linux, Windows, MacOS and Raspbian/ARM7 are distributed in the tarball for each release. To install
the binaries, download and extract the tarball to a directory of your choice.

The NodeRED, NodeJS and Python packages are installable from the public repositories:
- [node-red-contrib-uhppoted](https://www.npmjs.com/package/node-red-contrib-uhppoted)
- [uhppoted-nodejs](https://www.npmjs.com/package/uhppoted)
- [uhppoted](https://pypi.org/project/uhppoted)

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

## Donations

If you like the project and you really want to :slightly_smiling_face::

[<img height="41" width="174" src="https://cdn.buymeacoffee.com/buttons/default-orange.png">](https://www.buymeacoffee.com/uhppoted)


## References and related projects

1.  [carbonsphere/UHPPOTE](https://github.com/carbonsphere/UHPPOTE) `PHP`
2.  [carbonsphere/DoorControl](https://github.com/carbonsphere/DoorControl) `PHP`
2.  [andrewvaughan/uhppote-rfid](https://github.com/andrewvaughan/uhppote-rfid) `Python`
3.  [tachoknight/uhppote-tools](https://github.com/tachoknight/uhppote-tools): `Go`
4.  [jjhuff/uhppote-go](https://github.com/jjhuff/uhppote-go): `Go`
5.  [pawl/Chinese-RFID-Access-Control-Library](https://github.com/pawl/Chinese-RFID-Access-Control-Library)
6.  [Dallas Makerspace:Reverse Engineering RFID Reader](https://dallasmakerspace.org/wiki/ReverseEngineeringRFIDReader)
7.  [wemakerspace/wiegand-daemon](https://github.com/wemakerspace/wiegand-daemon)
8.  [wemakerspace/weigeng-js](https://github.com/wemakerspace/weigeng-js)
9.  [realashleybailey/DoorControl](https://github.com/realashleybailey/DoorControl)
10. [kBrausew/ioBroker.wiegand-tcpip](https://github.com/kBrausew/ioBroker.wiegand-tcpip)
11. [TCP/IP Wiegand Access Controller (Zutrittskontrolle)](https://ingenier.wordpress.com/zutrittskontrolle)
12. [YouTube: UHPPOTE 2 door basic set-up](https://www.youtube.com/watch?v=P8mxOY_IF4I&t=197s)
13. [wouterdebie/uhppote-rs](https://github.com/wouterdebie/uhppote-rs)

