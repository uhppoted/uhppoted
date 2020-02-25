## v0.60

## TODO

### uhppote
- [ ] concurrent requests
- [ ] update tests with Errorf + return to Fatalf
- [ ] commonalise ACL
- [ ] commonalise configuration
- [ ] make types consistent across API
- [ ] Genericize message unit tests
- [ ] Convert to 1.13 error handling
- [ ] Rework UHPPOTE response messages to use factory
- [ ] Fix golint errors
- [ ] Invert conf Unmarshal so that it iterates struct rather than file (simplifies e.g. DeviceMap)
- [ ] Rework plist encoder/decoder to be only for launchd (and remove 'parse' from daemonize/undaemonize)
- [ ] Unify event buffer operations

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
- [ ] Support for multiple brokers
- [ ] NACL/tweetnacl
- [ ] Report system events for e.g. listen bound/not bound

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
    - event buffer logic
17. GUI
    - [Muon](https://github.com/ImVexed/muon) 
    - [webview](https://github.com/zserge/webview)
    - [fyne](https://github.com/fyne-io/fyne)
    - https://instadeq.com/blog/posts/things-end-users-care-about-but-programmers-dont
    - [Naked Objects](https://en.wikipedia.org/wiki/Naked_objects)
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
30. Implement a lightweight end-to-end encryption protocol 
     - [MLS](https://mrosenberg.pub/cryptography/2019/07/10/molasses.html)
     - NACL/tweetnacl
31.  Consider moving to event bus architecture (?)
32. [Open Policy Agent](https://github.com/open-policy-agent) - for permissions
33. CouchBase (for JSON DB)
34. [AnyLog](https://blog.acolyer.org/2020/02/24/anylog)
35. [gRPC](https://www.programmableweb.com/news/how-to-build-streaming-api-using-grpc/how-to/2020/02/21)
