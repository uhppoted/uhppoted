## v0.7.0

- [x] Docker image for `uhppoted-mqtt`

## TODO

- [ ] [TOML](https://toml.io) files
- [ ] Docker buildkit (Ref. https://pythonspeed.com/articles/docker-buildkit)
- [ ] Docker compose (Ref. https://jvns.ca/blog/2021/01/04/docker-compose-is-nice)

- [ ] Integration tests
      - https://stackoverflow.com/questions/57549439/how-do-i-use-docker-with-github-actions
- [ ] Default bind and listen addresses to 0.0.0.0
- [ ][Architectural Decision Records](https://github.com/joelparkerhenderson/architecture_decision_record)
- [ ] Create release builds on github to avoid embedded references to local directories
- [ ] Use REST call to wait for container initialisation to complete
- [ ] Logo (https://math.stackexchange.com/questions/3742825/why-is-the-penrose-triangle-impossible)

### uhppote

- [ ] Erlang OTP-like supervision trees
- [ ] concurrent requests
- [ ] update tests with Errorf + return to Fatalf
- [ ] commonalise configuration
- [ ] make types consistent across API
- [ ] Genericize message unit tests
- [ ] Convert to 1.13 error handling throughout
- [ ] Rework UHPPOTE response messages to use factory
- [ ] Fix golint errors
- [ ] Invert conf Unmarshal so that it iterates struct rather than file (simplifies e.g. DeviceMap)
- [ ] Rework plist encoder/decoder to be only for launchd (and remove 'parse' from daemonize/undaemonize)
- [ ] Unify event buffer operations
- [ ] Convert configuration files to [TOML](https://github.com/toml-lang/toml)
- [ ] UQL (Ref. https://github.com/liljencrantz/crush)
- [ ] Use /var/lock for lockfiles
- [ ] https://tailscale.com/blog/netaddr-new-ip-type-for-go/

### uhppoted-api
- [ ] websocket + GraphQL (?)
      - [Hasura](https://hasura.io)
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
- [x] uhppoted-rest: PUT card
- [x] uhppoted-rest: DELETE card
- [ ] uhppoted-rest: get-events date/id range
- [x] commonalise functionality with uhppoted-mqttd

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
- [ ] https://www.asyncapi.com/docs/getting-started/coming-from-openapi

### CLI
- [x] Rework grant/revoke for individual doors (labelled)
- [x] get-acl
- [ ] Human readable output for e.g. get-status
- [ ] JSON formatted output for e.g. get-status
- [ ] Interactive shell (https://drewdevault.com/2019/09/02/Interactive-SSH-programs.html)
- [ ] use flag.FlagSet for commands
- [x] Default to commmon config file
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
      - [Taste](https://taste.tools)
- [ ] godoc
- [ ] build documentation
- [ ] install documentation
- [ ] user manuals
- [ ] man/info page
- [ ] cookbook for e.g. logging to AWS, VPN's etc
      - https://www.snia.org/about/corporate_info/logos/tutorial_graphics

### Other

1.  Rework uhppote to use bidirectional channel to serialize requests
2.  Consistently include device serial number in output e.g. of get-time
3.  Look into UDP multicast
4.  Look into ARP for set-address
6.  Integration tests
7.  Verify fields in listen events/status replies against SDK:
    - battery status can be (at least) 0x00, 0x01 and 0x04
8.  Mojave/HomeKit
9.  Phoenix UI
10. Update build system to [CMake or Meson](http://anadoxin.org/blog/is-it-worth-using-make.html)
11. step-ca (https://smallstep.com/blog/private-acme-server)
12. fuse
13. EventLogger 
    - MacOS: use [system logging](https://developer.apple.com/documentation/os/logging)
    - Windows: event logging
14. [Streamsheets](https://github.com/cedalo/streamsheets)
15. TLA+/Alloy models:
    - watchdog/health-check
    - concurrent connections
    - HOTP counter update
    - key-value stores
    - event buffer logic
16. GUI
    - [Muon](https://github.com/ImVexed/muon) 
    - [webview](https://github.com/zserge/webview)
    - [fyne](https://github.com/fyne-io/fyne)
    - https://instadeq.com/blog/posts/things-end-users-care-about-but-programmers-dont
    - [Naked Objects](https://en.wikipedia.org/wiki/Naked_objects)
    - [GIO](https://gioui.org)
    - [Revery](https://www.outrunlabs.com/revery)
    - [Nuklear](https://github.com/Immediate-Mode-UI/Nuklear)
    - [tview](https://github.com/rivo/tview)
    - [eDex-UI](https://github.com/GitSquared/edex-ui)
    - [Oracle APEX](https://apex.oracle.com/en)
    - [Proton Native](https://proton-native.js.org)
    
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
18. Update file watchers to fsnotify when that is merged into the standard library (1.4 ?)
    - https://github.com/golang/go/issues/4068
19. [Ballerina](https://ballerina.io)
20. [Eclipse Kura](https://www.eclipse.org/kura)
21. [Klee](https://klee.github.io)
22. [Fernet](https://github.com/fernet/spec/blob/master/Spec.md)
     - [cryptography.io:fernet](https://cryptography.io/en/latest/fernet)
23. [AsyncAPI](https://www.asyncapi.coms)
     - https://modeling-languages.com/asyncapi-modeling-editor-code-generator
24. go-fuzz
25. [OPA](https://github.com/open-policy-agent/opa) for permissions (?)
26. [Node-RED](https://hackaday.com/2020/01/15/automate-your-life-with-node-red-plus-a-dash-of-mqtt)
27. [Datomic ?](https://stackoverflow.com/questions/21245555/when-should-i-use-datomic)
28. [OCF-Over-Thread](https://www.infoq.com/news/2016/07/ocf-thread/)
29. Implement a lightweight end-to-end encryption protocol 
     - [MLS](https://mrosenberg.pub/cryptography/2019/07/10/molasses.html)
     - NACL/tweetnacl
30.  Consider moving to event bus architecture (?)
31. [Open Policy Agent](https://github.com/open-policy-agent) - for permissions
32. CouchBase (for JSON DB)
33. [AnyLog](https://blog.acolyer.org/2020/02/24/anylog)
34. [gRPC](https://www.programmableweb.com/news/how-to-build-streaming-api-using-grpc/how-to/2020/02/21)
35. [NodeRed/GreenGrass](https://iot.stackexchange.com/questions/2646/deploy-scripts-to-aws-greengrass-without-aws-lambda)
36. [dolt](https://github.com/liquidata-inc/dolt)
37. [Go/OpenAPI](https://stackoverflow.com/questions/48752908/openapi-spec-validation-in-golang)
38. [Datasette](https://datasette.readthedocs.io/en/stable)
39. [IDP: KeyClock](https://www.keycloak.org)
40. [IDP: shibboleth](https://www.shibboleth.net/products/identity-provider)
41. [Observability](https://charity.wtf/2020/03/03/observability-is-a-many-splendored-thing)
    -  auditor/monitor subsystem
42. [keys.pub](https://keys.pub/#what-is-it)
43. [Google Sheets](https://github.com/DecentLabs/officeAir)
44. [Chromium Embedded Framework](https://en.wikipedia.org/wiki/Chromium_Embedded_Framework)
45. [rsync.net](https://www.rsync.net)
46. [Rasbian/qemu](https://beta7.io/posts/running-raspbian-in-qemu.html)
47. [WebRTC](https://github.com/pion/webrtc)
48. [UDP tunnelling: ssh/nc](https://superuser.com/questions/53103/udp-traffic-through-ssh-tunnel)
49. [UDP tunnelling: socat](http://www.morch.com/2011/07/05/forwarding-snmp-ports-over-ssh-using-socat/)
50. [OpenMTC](https://www.openmtc.org)
51. [RSocket](https://rsocket.io)
52. For demo system: https://dev.to/github/10-standout-github-profile-readmes-h2o
53. [immudb](https://immudb.io)
54. [Oracle Apex](https://apex.oracle.com/en/learn/getting-started)
55. [Airtable](https://airtable.com)
56. [Notion](https://www.notion.so)
57. [Coda](https://coda.io)
58. [AnywhereDoor](https://github.com/iayanpahwa/anywheredoor)
59. [db](https://euandre.org/2020/08/31/the-database-i-wish-i-had.html)
60. [GNS3](https://www.gns3.com/software)
61. [kore](https://kore.io)
63. [Wireshark protocol](https://networkengineering.stackexchange.com/questions/67586/how-do-i-use-the-wireshark-i-o-graph-to-plot-the-value-of-an-arbitrary-bit-in-th)
64. [spiped](https://www.tarsnap.com/spiped.html)
65. [Tink](https://github.com/google/tink)
66. [Grafana](https://grafana.com)
67. [IOT monitoring tools](https://iot.stackexchange.com/questions/5295/monitoring-tools-for-iot-systems)
68. https://etcd.io
69. https://rollbar.com/pricing/
70. [IoT: balena](https://www.balena.io)
71. [bleve](https://blevesearch.com)
72. [ZUI](https://zircleui.github.io/docs/examples/home.html)
73. [Gitea](https://gitea.io/en-us)
74. [FHEM](https://fhem.de)
75. [ioBroker](https://iobroker.net)
76. [OpenHAB](https://www.openhab.org)
77. [Matrix](https://matrix.org/docs/spec)
78. [Thrift](https://thrift.apache.org)