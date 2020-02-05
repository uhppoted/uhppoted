## v0.50

## IN PROGRESS

### uhppoted-mqtt

- [x] get-devices
- [x] get-device
- [x] get-status
- [x] get-time
- [x] set-time
- [x] get-door-delay
- [x] set-door-delay
- [x] get-door-control
- [x] set-door-control
- [x] get-cards
- [x] delete-cards
- [x] get-card
- [x] put-card
- [x] delete-card
- [x] get-events
- [x] get-event
- [x] get-events: date/ID range
- [x] listen/events
- [x] command protocol: reply topic
- [x] command protocol: request ID
- [x] command protocol: authentication/HOTP
- [x] command protocol: authorisation
- [x] command protocol: add 'client-id' to response meta-info (?)
- [x] file watch on HOTP, user and permissions files
- [x] listen/events: retrieve and send actual events
- [x] clean up 'Request' implementation
- [x] command protocol: rework response JSON marshaling
- [x] command protocol: add 'operation' to response meta-info
- [x] wrap request handling in go routine
- [x] rework GetDevices to also find 'known' devices
- [x] sign & encrypt
- [x] move UHPPOTE from context to UHPPOTED
- [x] move incoming requests to /requests subtopic
- [x] health check
- [x] Make health check interval configurable 
- [x] Send health-check event to system topic
- [x] watchdog
- [x] Make watchdog check interval configurable 
- [x] Include get-listener in health check
- [x] daemonize/darwin
- [x] configurable healthcheck parameters
- [ ] daemonize/linux
- [ ] daemonize/windows
- [ ] 'dump config' command
- [ ] Replace all "path" with "filepath"
- [ ] Rework uhppoted API functions to use errors.Is(..) rather than returning status
  - https://blog.golang.org/go1.13-errors
- [ ] Fix weird error on restart after failure (might be related to filesystem permissions after daemonize)
      '''
      2020/02/05 10:12:13 INFO  publishing events to twystd/uhppoted/gateway/events
      2020/02/05 10:12:13 WARN  0            listen               Failed to open UDP socket [listen udp 192.168.1.100:60001: bind: address already in use]
      2020/02/05 10:12:13 [net]      incoming stopped with error EOF
      2020/02/05 10:12:13 [net]      error triggered, stopping
      2020/02/05 10:12:13 [net]      logic stopped
      2020/02/05 10:12:13 [net]      incoming stopped with error EOF
      2020/02/05 10:12:13 [net]      error triggered, stopping
      2020/02/05 10:12:13 [net]      logic stopped
      2020/02/05 10:12:13 [net]      incoming stopped with error EOF
      2020/02/05 10:12:13 [net]      error triggered, stopping
      '''

- [x] subscribe
- [x] error handling
- [x] 'reply'
- [x] TLS connection
- [x] TLS connection: client authentication
- [ ] Exit on error when multiple mqttd's are started (or generate unique client ID??? or .. lockfile)
- [ ] Implement retry + backoff for connection to broker
- [ ] Implement retry + backoff for 'listen'
- [ ] Rework listen logic to handle errors robustly
- [x] Identify UTO311-L01..L04 based on serial number prefix
- [ ] Rework conf.Marshal to use Range

- [x] Fix go vet errors
- [x] conf file decoder with reflection
- [x] conf file decoder: embedded structs
- [x] conf file encoder
- [x] docker: simulator
- [x] UT0311-L0x encoding: unmarshal arrays of structs (for broadcast)
- [x] Move version to [LDFLAGS](https://stackoverflow.com/questions/28459102/golang-compile-environment-variable-into-binary)

## TODO

### uhppote
- [ ] concurrent requests
- [ ] update tests with Errorf + return to Fatalf
- [ ] commonalise ACL
- [ ] commonalise configuration
- [ ] make types consistent across API
- [ ] Genericize message unit tests
- [ ] Add Rasbian/ARM7 target to build
- [ ] Convert to 1.13 error handling
- [ ] Rework UHPPOTE response messages to use factory
- [ ] Fix golint errors
- [ ] Invert conf Unmarshal so that it iterates struct rather than file (simplifies e.g. DeviceMap)
- [ ] Rework plist encoder/decoder to be only for launchd (and remove 'parse' from daemonize/undaemonize)

### uhppoted
- [ ] websocket + GraphQL (?)
- [ ] IFTTT
- [ ] Braid (?)
- [ ] MacOS launchd socket handoff
- [ ] Linux systemd socket handoff
- [ ] conf file decoder: JSON
- [ ] Rework plist encoder
- [ ] move ACL and events to separate API's
- [ ] Make events consistent across everything
- [ ] Rework uhppoted-xxx Run, etc to use [method expressions](https://talks.golang.org/2012/10things.slide#9)
- [ ] system API (for health-check, watchdog, configuration, etc)
- [ ] Parallel-ize health-check 

### uhppoted-rest
- [ ] Get events after XXX
- [ ] Client certificate revocation list
- [ ] Rework GetDevices to also find 'known' devices
- [ ] uhppoted-rest: PUT card
- [ ] uhppoted-rest: DELETE card
- [ ] uhppoted-rest: get-events date/id range
- [ ] commonalise functionality with uhppoted-mqttd

### uhppoted-mqtt
- [ ] last-will-and-testament (?)
- [ ] publish add/delete card, etc to event stream
- [ ] MQTT v5.0
- [ ] [JSON-RPC](https://en.wikipedia.org/wiki/JSON-RPC) (?)
- [ ] Add to CLI
- [ ] Non-ephemeral key transport:  https://tools.ietf.org/html/rfc5990#appendix-A
- [ ] user:open/get permissions require matching card number 
- [ ] [AEAD](http://alexander.holbreich.org/message-authentication)

### CLI
- [ ] Rework grant/revoke for individual doors (labelled)
- [ ] get-acl
- [ ] Human readable output for e.g. get-status
- [ ] JSON formatted output for e.g. get-status
- [ ] Interactive shell (https://drewdevault.com/2019/09/02/Interactive-SSH-programs.html)
- [ ] use flag.FlagSet for commands
- [ ] Default to commmon config file
- [ ] Use (loadable) text/template for output formats
- [ ] Rework GetDevices to also find 'known' devices
- [ ] events: retrieve and show actual events
- [ ] Generate OTP secret + QR code

### simulator
- [ ] concurrent requests
- [ ] simulator-cli
- [ ] HTML
- [ ] Rework simulator.run to use rx channels
- [ ] Reload simulator on device file change
- [ ] Implement JSON unmarshal to initialise default values
- [ ] Swagger UI
- [ ] Autodetect gzipped files (https://stackoverflow.com/questions/28309988/how-to-read-from-either-gzip-or-plain-text-reader-in-golang)

### Documentation

- [ ] TeX protocol description
- [ ] ASN.1 protocol specification
- [ ] godoc
- [ ] build documentation
- [ ] install documentation
- [ ] user manuals
- [ ] man/info page

### Other

1.  Update to use modules
    - Refactor uhppoted as uhppoted-api
    - Rename and restructure repo
2.  Rework uhppote to use bidirectional channel to serialize requests
3.  Consistently include device serial number in output e.g. of get-time
4.  Look into UDP multicast
5.  Look into ARP for set-address
6.  github project page
7.  Integration tests
8.  Verify fields in listen events/status replies against SDK:
    - battery status can be (at least) 0x00, 0x01 and 0x04
9.  Mojave/HomeKit
10. Phoenix UI
11. Update build system to [CMake or Meson](http://anadoxin.org/blog/is-it-worth-using-make.html)
12. step-ca (https://smallstep.com/blog/private-acme-server)
13. fuse
14. EventLogger 
    - MacOS: use [system logging](https://developer.apple.com/documentation/os/logging)
    - Windows: event logging
15. [Streamsheets](https://github.com/cedalo/streamsheets)
16. TLA+/Alloy models:
    - watchdog/health-check
    - concurrent connections
    - HOTP counter update
    - key-value stores
17. GUI
    - [Muon](https://github.com/ImVexed/muon) 
    - [webview](https://github.com/zserge/webview)
    - [fyne](https://github.com/fyne-io/fyne)
    - https://instadeq.com/blog/posts/things-end-users-care-about-but-programmers-dont
18. PDL + go generate
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
19. Update file watchers to fsnotify when that is merged into the standard library (1.4 ?)
    - https://github.com/golang/go/issues/4068
20. [Ballerina](https://ballerina.io)
21. [Eclipse Kura](https://www.eclipse.org/kura)
22. [Klee](https://klee.github.io)
23. [Fernet](https://github.com/fernet/spec/blob/master/Spec.md)
     - [cryptography.io:fernet](https://cryptography.io/en/latest/fernet)
24. [AsyncAPI](https://www.asyncapi.coms)
     - https://modeling-languages.com/asyncapi-modeling-editor-code-generator
25. go-fuzz
26. [OPA](https://github.com/open-policy-agent/opa) for permissions (?)
27. [Node-RED](https://hackaday.com/2020/01/15/automate-your-life-with-node-red-plus-a-dash-of-mqtt)
28. [Datomic ?](https://stackoverflow.com/questions/21245555/when-should-i-use-datomic)
29. [OCF-Over-Thread](https://www.infoq.com/news/2016/07/ocf-thread/)
