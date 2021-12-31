## TROUBLE SHOOTING

##### Is your controller firmware version supported?

- The minimum supported controller firmware version is `6.56`
- The lowest known version in actual use is `6.62`
- The code is developed and tested against `8.92`

Firmware version `v6.62` is known to send anomalous _listen events_ for which a patch has been added to:

- [`uhppote-core`](https://github.com/uhppoted/uhppote-core/blob/75a185a48184ecb68a07a09ebdd9ea1a8f96ba2c/encoding/UTO311-L0x/UT0311-L0x.go#L201-L204)
- [`uhppote-simulator`](https://github.com/uhppoted/uhppote-simulator/blob/f599512fb821c892a75786bbe4f35f6ebb4563d9/commands/run.go#L125-L134)
- [`node-red-contrib-uhppoted`](https://github.com/uhppoted/node-red-contrib-uhppoted/blob/74de32d62bee8097c03c9a1abc2bb45b0160f7b2/nodes/codec.js#L93-L100)

##### Not receiving _door open_, _door close_ or _door button_ events?

`door open`, `door close` and `door button` events need to be specifically enabled using the `record-special-events` function.


##### `uhppote-cli`: Microsoft Windows 7 crash with `unknown pc ...`

On Microsoft Windows 7, a `uhppote-cli` crash report with `unknown pc ...` (see example below) is typically associated
with an anti-virus program (principally _webroot_) and appears to be related to an incompatibiliy between the structure of
the EXE file created by the Go compiler and the anti-virus signature detection code:

- https://github.com/golang/go/issues/45153
- https://github.com/golang/go/issues/40878

As of 2021-07-01, the only known workaround is to disable _webroot_ - there appear to be no plans to address this in either
the Go compiler or in _webroot_.

```
Exception 0xc0000005 0x0 0x77776fff 0xb60000
PC=0xb60000

runtime: unknown pc 0xb60000
stack: frame={sp:0x22f220, fp:0x0} stack=[0x0,0x22ff30)
000000000022f120:  0000000000200000  000000000022f170
000000000022f130:  000000000022f170  000000000022f170
000000000022f140:  0000000000000020  00000000775e04c6
000000000022f150:  0000000000000001  0000000000001112
000000000022f160:  0000000000000000  000007fefaa60000
000000000022f170:  000000000022f348  000000007763be5f
...
...
rax     0x22f3e2
rbx     0x775c0000
rcx     0x8000000000000000
rdi     0x774a9bfa
rsi     0x1
rbp     0xbbfc00
rsp     0x22f220
r8      0x0
r9      0x0
r10     0x0
r11     0x22f3e2
r12     0x0
r13     0x776d614c
r14     0x776d4200
r15     0x776d8050
rip     0xb60000
rflags  0x10206
cs      0x33
fs      0x53
gs      0x2b
```

### MacOS firewall: Accept incoming connections

See [StackExchange: How to get rid of firewall "accept incoming connections" dialog?](https://apple.stackexchange.com/questions/3271/how-to-get-rid-of-firewall-accept-incoming-connections-dialog)


