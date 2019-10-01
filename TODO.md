## v0.04

## IN PROGRESS

### uhppoted

- [ ] Daemon/service interface
- [ ] websocket command interface
- [ ] MQTT 
- [ ] GraphQL
- [ ] watchdog

#### REST API
- [ ] TLS with pinned certificates
- [ ] Swagger UI
- [ ] gzip response
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
3.  fuse
4.  Look into ARP for set-address
5.  Rework error handling to use Wrap/Frame
6.  Integration tests
7.  Dig into simulator not receiving broadcast unless listening on 0.0.0.0:60000
    (Ref. https://groups.google.com/forum/#!topic/golang-nuts/nbmYWwHCgPc)
8.  Verify fields in listen events/status replies:
    - battery status can be (at least) 0x00, 0x01 and 0x04
9.  Update to use modules
10. Mojave/HomeKit
11. Phoenix UI
12. step-ca (https://smallstep.com/blog/private-acme-server)

