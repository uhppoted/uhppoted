# TODO

-- Dig into simulator not receiving broadcast unless listening on 0.0.0.0:60000
   Ref. https://groups.google.com/forum/#!topic/golang-nuts/nbmYWwHCgPc

**Clean up 'card not found' handling in uhppote**
**Autodetect gzipped files**
(https://stackoverflow.com/questions/28309988/how-to-read-from-either-gzip-or-plain-text-reader-in-golang)

0. simulator: cards
1. Load cards from TSV file
2. Human readable output for e.g. get-status
3. JSON formatted output for e.g. get-status
4. Consistently include device serial number in output e.g. of get-time
5. Document protocol
6. Simulator
7. fuse
8. Look into ARP for set-address
9. Rework error handling to use Wrap/Frame
10. godoc
11. Integration tests
12. --config:none option
