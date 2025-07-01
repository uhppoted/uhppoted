# TODO

- [x] Update to Go 1.24
- [x] anti-passback (cf. https://github.com/uhppoted/uhppoted/issues/60)

- [ ] Format/lint Makefile (https://github.com/EbodShojaei/bake)
- [ ] Format README like https://github.com/BelfrySCAD/BOSL2/wiki

- [ ] Release script
      - [x] npm: `... uhppoted-lib-nodejs is waiting for release on npm`
      - [ ] README: doesn't check uhppoted-lib-python
      - [ ] README: copy CHANGELOG to clipboard
      - [ ] Pre-release builds
      - [ ] Build checksums
      - [ ] Bump version in uhppote-core and Makefiles
      - [ ] https://pydoit.org/
      - [ ] https://github.com/bitfield/script

- [ ] ghcr
      - [ ] Build versioned images before 'latest' so that the install instruction refers to 'latest'
      - [ ] Extract version fromm release tag
            - https://stackoverflow.com/questions/58177786/get-the-current-pushed-tag-in-github-actions
            - https://github.com/orgs/community/discussions/26625
            - https://stackoverflow.com/questions/67231657/how-to-replace-string-in-expression-with-github-actions

- [ ] Some config ideas
      - https://chshersh.com/blog/2025-01-06-the-most-elegant-configuration-language.html
- [ ] Make `uhppoted` a service user (`sudo adduser --system uhppoted`)
- [ ] https://forum.universal-devices.com/topic/38892-wiegand-wg26-four-door-access-controller-aka-wgaccess-uhppote-wg2004/

## TODO

- [ ] Replace DeviceID with generic type
- [ ] [TAK](https://hackaday.com/2022/09/08/the-tak-ecosystem-military-coordination-goes-open-source)
- [ ] [Home Assistant](https://developers.home-assistant.io)
- [ ] [Koyeb](https://www.koyeb.com)
- [ ] Take a look at [CloudFlare Spectrum](https://www.cloudflare.com/products/cloudflare-spectrum) - proxies UDP!
- [ ] Switch to libsodium
      - https://blog.trailofbits.com/2019/07/08/fuck-rsa/
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
- [ ] uhppoted-rest: get-events date/id range

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
- [ ] Human readable output for e.g. get-status
- [ ] JSON formatted output for e.g. get-status
- [ ] Interactive shell (https://drewdevault.com/2019/09/02/Interactive-SSH-programs.html)
- [ ] use flag.FlagSet for commands
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

- https://www.mkdocs.org/
- [ ] TeX protocol description
- [ ] ASN.1 protocol specification
      - [Taste](https://taste.tools)
      - https://medium.com/erlang-battleground/erlang-asn-1-abstract-syntax-notation-one-deeb8300f479
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
3.  ~~Look into UDP multicast~~
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
    - [Flutter](https://blog.whidev.com/native-looking-desktop-app-with-flutter)
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
    - [imtui](https://github.com/ggerganov/imtui)
    
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
26. [Datomic ?](https://stackoverflow.com/questions/21245555/when-should-i-use-datomic)
27. [OCF-Over-Thread](https://www.infoq.com/news/2016/07/ocf-thread/)
28. Implement a lightweight end-to-end encryption protocol 
     - [MLS](https://mrosenberg.pub/cryptography/2019/07/10/molasses.html)
     - NACL/tweetnacl
29.  Consider moving to event bus architecture (?)
30. [Open Policy Agent](https://github.com/open-policy-agent) - for permissions
31. [zk-SNARK](https://www.entropy1729.com/the-hunting-of-the-zk-snark)
32. CouchBase (for JSON DB)
33. [AnyLog](https://blog.acolyer.org/2020/02/24/anylog)
34. [gRPC](https://www.programmableweb.com/news/how-to-build-streaming-api-using-grpc/how-to/2020/02/21)
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
79. [LipGloss](https://charm.sh)
80. [Svelte/Tauri](https://css-tricks.com/how-i-built-a-cross-platform-desktop-application-with-svelte-redis-and-rust)
81. https://blog.marcua.net/2022/02/20/data-diffs-algorithms-for-explaining-what-changed-in-a-dataset.html
82. [Charm](https://charm.sh)
83. [Riffle][https://riffle.systems/essays/prelude]
84. http://anachronauts.club/~voidstar/log/2022-03-24-openapi-for-binfmt.gmi
85. https://siraben.dev/2022/03/22/tree-sitter-linter.html
86. https://maori.geek.nz/golang-desktop-app-with-webview-lorca-wasm-and-bazel-3283813bf89
87. https://gokrazy.org/
88. https://www.fstar-lang.org/papers/EverParse3D.pdf
89. [Hesiod](https://en.wikipedia.org/wiki/Hesiod_(name_service))
90. [Grist](https://www.getgrist.com/)
91. ~~[HomeAssistant](https://www.home-assistant.io)~~
92. [FPGA](https://github.com/JulianKemmerer/PipelineC-Graphics/blob/main/doc/Sphery-vs-Shapes.pdf)
93. [Charm::Bubbletea](https://dlvhdr.me/posts/the-renaissance-of-the-command-line)
94. [Textualize](https://www.textualize.io)
95. [OpenTelemetry](https://github.com/open-telemetry/opentelemetry-demo/tree/v1.0.0)
96. [Reowolf](https://reowolf.net/)
97. [OSDP](https://www.getkisi.com/guides/osdp)
98. [wails](https://wails.io/)
99. [Sandstorm](https://sandstorm.io/)
100. [Spline](https://spline.design)
101. [Womp](https://www.womp.com)
102. [Voice](https://arstechnica.com/gadgets/2022/12/with-voice-assistants-in-trouble-home-assistant-starts-a-local-alternative/)
103. [OPC UA](https://stackoverflow.com/questions/52074597/when-it-is-justified-to-use-ua-opc-and-ua-opc-architectures-over-mqtt)
104. [Homebridge](https://homebridge.io)
105. [serverless DB](https://leerob.io/blog/backend)
106. [Tigris](https://www.tigrisdata.com/docs/quickstarts)
107. [Gitbook](https://developer.gitbook.com/gitbook-api/reference)
108. [Matter](https://interrupt.memfault.com/blog/memfault-matter)
109. [marimo](https://docs.marimo.io)
110. https://github.com/huginn/huginn
111. [slint](https://slint.dev)
112. https://github.com/francoismichel/ssh3
113. https://steampipe.io/
114. https://nats.io/
115. https://www.jolie-lang.org/index.html
116. https://arxiv.org/pdf/2503.04084v1