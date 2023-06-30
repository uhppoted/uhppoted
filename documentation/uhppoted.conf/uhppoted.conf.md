# _uhppoted.conf_

`uhppoted.conf` is the communal configuration file shared by all the `uhppoted` project modules, an example of which
can be found [here](uhppoted.conf).

It comprises the following sections:

1. `system`
2. `controllers`

## `system`

The `system` section defines the configurable parameters that are common to all the projects that use _uhppote-core_:

| Parameter                       | Default value   | Description                                                            |
|---------------------------------|-----------------|------------------------------------------------------------------------|
| bind.address                    | 0.0.0.0         | IPv4 UDP address to bind to when creating a connection to a controller |
| broadcast.address               | 255.255.255.255 | IPv4 UDP broadcast address                                             |
| listen.address                  | 0.0.0.0:60001   | IPv4 address on which to listen for controller events                  |
| timeout                         | 2.5s            | Time to wait for a controller response before returning an error       |
| monitoring.healthcheck.interval | 15s             | Interval at which to check controller status                           |
| monitoring.healthcheck.idle     | 1m0s            | Time after which an unreachable controller is marked 'idle'            |
| monitoring.healthcheck.ignore   | 5m0s            | Time after which an unreachable controller is marked 'ignore'          |
| monitoring.watchdog.interval    | 5s              | Interval at which the health-check subsystem is checked                |
| card.format                     | Wiegand-26      | Card format to use for validation (_any_ or _Wiegand-26_)              |

## `controllers`
