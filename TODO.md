## v0.40

## IN PROGRESS

### uhppoted

- [x] use flag.FlagSet for commands
- [ ] linux.Daemonize: implement custom flag for uid:gid and commonalize Parse across all commands
- [ ] include missing/unexpected devices in monitoring
- [x] restructure 'run' command
- [x] keep-alive/health-check
- [x] watchdog
- [x] add --debug to command line flags
- [x] MacOS launchd SIG_TERM
- [x] MacOS launchd --daemonize
- [x] MacOS launchd reinstall
- [x] MacOS launchd --undaemonize
- [x] MacOS launchd XML plist file marshalling/unmarshalling
- [x] MacOS launchd newsyslog log rotate
- [ ] MacOS launchd socket handoff
- [x] Linux systemd service
- [x] Linux systemd logrotate
- [x] Linux systemd replace chown with uid+gid
- [x] Linux systemd uhppoted daemonize --user uhppoted:uhppoted
- [x] Linux systemd create initial /etc/uhppoted/uhppote.conf file
- [x] Linux systemd add note for UDP UFW rules for broadcast
- [ ] Linux systemd socket handoff
- [x] Windows service
- [x] Windows service: use %PROGRAMDATA% folder for conf files
- [x] Windows service: log to Event Log
- [x] Windows service: rethink --console option

- [ ] MQTT 
- [ ] websocket command interface
- [ ] GraphQL
- [ ] watchdog

#### REST API
- [x] TLS with mutual authentication
- [x] Swagger UI
- [x] gzip response
- [x] Add HTTP method to dispatch matching
- [x] Log internal error information
- [x] Include logging in context
- [x] Get devices
- [x] Get device
- [x] Get device status
- [x] Get device time
- [x] Set device time
- [x] Get door delay
- [x] Set door delay
- [x] Get cards
- [x] Get card
- [x] Delete cards
- [x] Delete card
- [x] Get events
- [x] Get event

## TODO

### uhppote
- [ ] concurrent requests
- [ ] update tests with Errorf + return to Fatalf
- [ ] commonalise ACL implementation

### CLI
- [ ] Rework grant/revoke for individual doors (labelled)
- [ ] get-acl
- [ ] Human readable output for e.g. get-status
- [ ] JSON formatted output for e.g. get-status
- [ ] Interactive shell (https://drewdevault.com/2019/09/02/Interactive-SSH-programs.html)

### simulator
- [ ] concurrent requests
- [ ] simulator-cli
- [ ] HTML
- [ ] Rework simulator.run to use rx channels
- [ ] Reload simulator on device file change
- [ ] Swagger UI
- [ ] Autodetect gzipped files (https://stackoverflow.com/questions/28309988/how-to-read-from-either-gzip-or-plain-text-reader-in-golang)

### Documentation

- [ ] TeX protocol description
- [ ] ASN.1 protocol specification
- [ ] godoc
- [ ] build documentation
- [ ] install documentation
- [ ] user manuals

### Other

1.  Rework uhppote to use bidirectional channel to serialize requests
2.  Consistently include device serial number in output e.g. of get-time
3.  Look into UDP multicast
4.  Look into ARP for set-address
5.  Rework error handling to use Wrap/Frame
6.  Integration tests
7.  conf file encoder/decoder with reflection and/or JSON
8.  Verify fields in listen events/status replies:
    - battery status can be (at least) 0x00, 0x01 and 0x04
9.  Update to use modules
10. Mojave/HomeKit
11. Phoenix UI
12. step-ca (https://smallstep.com/blog/private-acme-server)
13. fuse
14. EventLogger 
    - MacOS: use [system logging](https://developer.apple.com/documentation/os/logging)
    - Windows: event logging
16. [Streamsheets](https://github.com/cedalo/streamsheets)
17. PDL + go generate
    - [lipPDL](http://nmedit.sourceforge.net/subprojects/libpdl.html)
    - [Diva](http://www.diva-portal.org/smash/get/diva2:407713/FULLTEXT01.pdf)
    - [PADS/ML](https://pads.cs.tufts.edu/papers/tfp07.pdf)
    - [PADS](https://www.cs.princeton.edu/~dpw/papers/700popl06.pdf)
    - [DataScript](https://www.researchgate.net/publication/221108676_DataScript-_A_Specification_and_Scripting_Language_for_Binary_Data)
    - [PADS/ML](https://www.cs.princeton.edu/~dpw/papers/padsml06.pdf)
    - [PADS Project](http://www.padsproj.org/)
    - [Mozilla IPDL](https://developer.mozilla.org/en-US/docs/Mozilla/IPDL/Tutorial)
    - [PDL: Failure Semanics](https://www.researchgate.net/publication/2784726_A_Protocol_Description_Language_for_Customizing_Failure_Semantics)
    - https://en.wikipedia.org/wiki/Abstract_Syntax_Notation_One

