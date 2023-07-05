# _uhppoted.conf_

`uhppoted.conf` is the communal configuration file shared by all the `uhppoted` project modules, an example of which
can be found [here](uhppoted.conf).

It comprises the following sections:

1. `system`
2. `controllers`
3. `REST`

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

The `controllers` section declares the access controllers known to the system. It is not generally required except for:
- ACL commands, which use the door names
- remote controllers i.e. not on the same LAN segment as the host and therefore not accessible via UDP broadcast
- controllers located in different timezones

A controller entry comprises:

| Parameter   | Description                                                                                          |
|-------------|------------------------------------------------------------------------------------------------------|
| name        | Alphanumeric controller name (optional). CLI commands can use the name instead of the serial number. |
| address     | IPv4 UDP address in the format _address:port_. Defaults to port 60000.                               |
| door 1      | Unique alphanumeric door name. Used to resolve ACL access permissions.                               |
| door 2      | Unique alphanumeric door name. Used to resolve ACL access permissions.                               |
| door 3      | Unique alphanumeric door name. Used to resolve ACL access permissions.                               |
| door 4      | Unique alphanumeric door name. Used to resolve ACL access permissions.                               |
| timezone    | Optional timezone. Used for converting between host time and controller time.                        |

Example controller entry for the controller with serial number 405419896:
```
UT0311-L0x.405419896.name = Alpha
UT0311-L0x.405419896.address = 192.168.1.100:60000
UT0311-L0x.405419896.door.1 = Gryffindor
UT0311-L0x.405419896.door.2 = Hufflepuff
UT0311-L0x.405419896.door.3 = Ravenclaw
UT0311-L0x.405419896.door.4 = Slytherin
UT0311-L0x.405419896.timezone = PDT
```

## `REST`

The `REST` section defines the configuration for the _uhppoted-rest_ REST gateway server.

| Parameter                    | Default                    | Description                                             |
|------------------------------|----------------------------|---------------------------------------------------------|
| rest.http.enabled            | false                      | Enables the unsecured HTTP server                       |
| rest.http.port               | 8080                       | Unsecured HTTP server port                              |
| rest.https.enabled           | true                       | Enable the secure HTTPS server                          |
| rest.https.port              | 8443                       | Secured HTTPS server port                               |
| rest.tls.key                 | \<etc\>/rest/uhppoted.key  | HTTPS server key (PEM)                                  |
| rest.tls.certificate         | \<etc\>/rest/uhppoted.cert | HTTPS server certificate (PEM)                          |
| rest.tls.ca                  | \<etc\>/rest/ca.cert       | HTTPS CA certifcate for client verification (PEM)       |
| rest.tls.client.certificates | true                       | Enables TLS mutual authentication                       |
| rest.CORS.enabled            | true                       | Enables CORS                                            |
| rest.auth.enabled            | false                      | Enables user/group authorisation                        |
| rest.auth.users              | \<etc\>/rest/users         | List of authorised users (JSON)                         |
| rest.auth.groups             | \<etc\>/rest/groups        | List of authorised user groups (JSON)                   |
| rest.auth.hotp.range         | 8                          | HOTP authentication: valid counter range                |
| rest.auth.hotp.secrets       |                            | HOTP authentication: user secrets                       |
| rest.auth.hotp.counters      | \<etc\>/rest/counters      | HOTP authentication: user counters                      |


