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

