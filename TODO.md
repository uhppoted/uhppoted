## v0.03

*FIX interface{} nil assignment*

1. Simulator
   - get-event
   - add door state/buttons to current state

# TODO

1.  Human readable output for e.g. get-status
2.  JSON formatted output for e.g. get-status
3.  Consistently include device serial number in output e.g. of get-time
4.  Document protocol
    - ASN.1
5.  fuse
6.  Look into ARP for set-address
7.  Rework error handling to use Wrap/Frame
8. godoc
9. Integration tests
10. Rework grant/revoke for individual doors (labelled)
11. Dig into simulator not receiving broadcast unless listening on 0.0.0.0:60000
    (Ref. https://groups.google.com/forum/#!topic/golang-nuts/nbmYWwHCgPc)
12. Autodetect gzipped files
    (Ref. ttps://stackoverflow.com/questions/28309988/how-to-read-from-either-gzip-or-plain-text-reader-in-golang)
13. Reload simulator on device file change
14. Verify fields in listen events/status replies:
    - battery status can be (at least) 0x00, 0x01 and 0x04
15. websocket command interface
16. get-acl
