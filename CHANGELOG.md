# CHANGELOG

## Unreleased

### Updated
1. Renamed _uhppoted-lib_ _master_ branch to _main_.


## [0.8.8](https://github.com/uhppoted/uhppoted/releases/tag/v0.8.8) - 2024-03-28

### Added
1. Added _uhppoted-app-home-assistant_ experimental Home Assistant integration.
2. `restore-default-parameters` function across all subprojects.
3. Added public Docker images for _uhppote-simulator_, uhppoted-rest_, uhppoted-mqtt_, and
   _uhppoted-httpd_ to ghcr.io.

### Updated
1. Bumped Go version to 1.22.
2. Reworked _uhppoted-app-wild-apricot_ member/group resolution logic.


## [0.8.7](https://github.com/uhppoted/uhppoted/releases/tag/v0.8.7) - 2023-12-01

### Added
1. `set-door-passcodes` function across all subprojects.
2. Added PostgreSQL bindings to _uhppoted-app-db_.
3. Added _Lua_ bindings to _uhppoted-codegen_.
4. Added Visual Studio C# examples to _uhppoted-dll_.
5. Added _live_ events to _uhppoted-mqtt_.
6. Added keypad emulation _uhppoted-wiegand_.

### Updated
1. Bumped Go version to 1.21.


## [0.8.6](https://github.com/uhppoted/uhppoted/releases/tag/v0.8.6) - 2023-08-30

### Added
1. Repackaged uhppoted-codegen Python bindings as _uhppoted-python_ for PyPI
2. Implemented `activate-keypads` function across all subprojects.
3. Added bindings to MySQL and Microsoft SQL server to _uhppoted-app-db_.
4. Added support for _tmpfs_ filesystems.
5. Added _Erlang_ bindings to _uhppoted-codegen_.
6. Preliminary documentation for _uhppoted.conf_ file.


## [0.8.5](https://github.com/uhppoted/uhppoted/releases/tag/v0.8.5) - 2023-06-14

### Added
1. Added _uhppoted-app-db_ project with initial sqlite3 support only.
2. Added support for PicoW to _uhppoted-wiegand_ project.
3. Implemented `set-interlock` function across all subprojects.
4. Added PHP bindings to _uhppoted-codegen_.
5. Added _tailscale_ integration to _uhppoted-tunnel_.
6. Repacked _uhppoted-codegen_ Python bindings as _uhppoted-python_ and uploaded to PyPI.


## [0.8.4](https://github.com/uhppoted/uhppoted/releases/tag/v0.8.4) - 2023-03-17

### Added
1. Included _uhppoted-wiegand_ project in submodules.

### Updated
1. Updated documentation and docker for card keypan PINs.


## [0.8.3](https://github.com/uhppoted/uhppoted/releases/tag/v0.8.3) - 2022-12-16

### Added
1. Added ARM64 to release build artifacts


## [0.8.2](https://github.com/uhppoted/uhppoted/releases/tag/v0.8.2) - 2022-10-14

### Added
1. _uhppoted-codegen_ interface generator

### Updated
1. Bumped Go version to 1.9

## [0.8.1](https://github.com/uhppoted/uhppoted/releases/tag/v0.8.1) - 2022-08-01

### Changed
1. Added support for human-readable event fields
2. Updated health-check to handle INADDR_ANY listen addresses correctly.
3. Minor bug fixes


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


