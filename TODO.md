# TODO


2.  Human readable output for e.g. get-status
3.  JSON formatted output for e.g. get-status
4.  Consistently include device serial number in output e.g. of get-time
5.  Document protocol
    - ASN.1
6.  Simulator
7.  fuse
8.  Look into ARP for set-address
9.  Rework error handling to use Wrap/Frame
10. godoc
11. Integration tests
14. Rework grant/revoke for individual doors (labelled)
15. Dig into simulator not receiving broadcast unless listening on 0.0.0.0:60000
    (Ref. https://groups.google.com/forum/#!topic/golang-nuts/nbmYWwHCgPc)
16. Autodetect gzipped files
    (Ref. ttps://stackoverflow.com/questions/28309988/how-to-read-from-either-gzip-or-plain-text-reader-in-golang)
17. Reload simulator on device file change
18. get-acl
19. Verify fields in listen events/status replies:
    - battery status can be (at least) 0x00, 0x01 and 0x04
