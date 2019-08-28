## v0.04

## IN PROGRESS

### uhppoted
- [ ] REST API (debug)
      - [ ] Add HTTP method to dispatch matching
      - [ ] Include logging in context
- [ ] websocket command interface
- [ ] MQTT 
- [ ] watchdog

## TODO

### CLI
- [ ] Rework grant/revoke for individual doors (labelled)
- [ ] get-acl
- [ ] Human readable output for e.g. get-status
- [ ] JSON formatted output for e.g. get-status

### simulator
- [ ] simulator-cli
- [ ] HTML
- [ ] Rework simulator.run to use rx channels
- [ ] Reload simulator on device file change

### Other

1.  Rework uhppote to use bidirectional channel to serialize requests
2.  Consistently include device serial number in output e.g. of get-time
3.  Document protocol
    - TeX
    - ASN.1
4.  fuse
5.  Look into ARP for set-address
6.  Rework error handling to use Wrap/Frame
7.  godoc
8.  Integration tests
9. Dig into simulator not receiving broadcast unless listening on 0.0.0.0:60000
    (Ref. https://groups.google.com/forum/#!topic/golang-nuts/nbmYWwHCgPc)
10. Autodetect gzipped files
    (Ref. ttps://stackoverflow.com/questions/28309988/how-to-read-from-either-gzip-or-plain-text-reader-in-golang)
11. Verify fields in listen events/status replies:
    - battery status can be (at least) 0x00, 0x01 and 0x04
12. Update to use modules

