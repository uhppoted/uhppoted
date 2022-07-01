# CHANGELOG

## [0.8.0](https://github.com/uhppoted/uhppoted/releases/tag/v0.8.0) - 2022-07-01

### Added
1. [uhppoted-httpd](https://github.com/uhppoted/uhppoted-httpd) browser user interface for access control management
2. [uhppoted-tunnel](https://github.com/uhppoted/uhppoted-tunnel) UDP tunnel to connect access control systems and controllers
running on disparate networks


## [0.7.3](https://github.com/uhppoted/uhppoted/releases/tag/v0.7.3) - 2022-06-01

### Added
1. [uhppoted-dll](https://github.com/uhppoted/uhppoted-dll) DLL/shared-lib/dylib for cross-language
   support

### Changed
1. Included -trimpath option in all build paths to remove local machine information from executables


## [0.7.2](https://github.com/uhppoted/uhppoted/releases/tag/v0.7.2) - 2022-01-27

### Changed

1. Replaced event rollover throughout with handling for _nil_ and _overwritten_ events
2. Reworked [`uhppoted-nodejs`](https://github.com/uhppoted/uhppoted-nodejs) for compatibility with NodeJS v14.18.3
   (cf. https://github.com/uhppoted/uhppoted-nodejs/issues/5)
3. Reworked [`node-red-contrib-uhppoted`](https://github.com/uhppoted/node-red-contrib-uhppoted) for compatibility with NodeJS v14.18.3
   (cf. https://github.com/uhppoted/uhppoted-nodejs/issues/5)


