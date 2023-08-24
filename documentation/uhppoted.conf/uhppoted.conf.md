# _uhppoted.conf_

`uhppoted.conf` is the communal configuration file shared by all the `uhppoted` project modules, an example of which
can be found [here](uhppoted.conf).

It comprises the following sections:

1. [`system`](#system)
2. [`controllers`](#controllers)
3. [`REST`](#rest)
4. [`MQTT`](#mqtt)
5. [`AWS`](#aws)
6. [`HTTPD`](#httpd)
7. [`Wild Apricot`](#wild-apricot)
7. [`OpenAPI`](#openapi)

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


## `MQTT`

The `MQTT` section defines the configuration for the _uhppoted-mqtt_ MQTT gateway server.

| Parameter                          | Default                         | Description                                             |
|------------------------------------|---------------------------------|---------------------------------------------------------|
| mqtt.server.ID                     | uhppoted                        | (future use)                                            |
| mqtt.connection.broker             | tcp://127.0.0.1:1883            | MQTT broker IP address and port                         |
| mqtt.connection.client.ID          | uhppoted-mqttd                  | Client ID for connection to MQTT broker                 |
| mqtt.connection.username           |                                 | User name for connection to MQTT broker                 |
| mqtt.connection.password           |                                 | Password for connection  to MQTT broker                 |
| mqtt.connection.broker.certificate | \<etc\>/mqtt/broker.cert        | CA certificate for MQTT broker                          |
| mqtt.connection.client.certificate | \<etc\>/mqtt/client.cert        | Client RSA certificate for TLS mutual authentication    |
| mqtt.connection.client.key         | \<etc\>/mqtt/client.key         | Client RSA key for TLS mutual authentication            |
| mqtt.connection.verify             |                                 | _allow-insecure_ disables TLS verification              |
| mqtt.topic.root                    | uhppoted/gateway                | Root MQTT topic for _uhppoted-mqtt_ messages            |
| mqtt.topic.requests                | ./requests                      | _request_ messages subtopic                             |
| mqtt.topic.replies                 | ./replies                       | _reply_ messages subtopic                               |
| mqtt.topic.events                  | ./events                        | _event_ messages subtopic                               |
| mqtt.topic.system                  | ./system                        | _system_ messages subtopic                              |
| mqtt.translation.locale            |                                 | Locale for internationalisation of messages             |
| mqtt.protocol.version              |                                 | MQTT (for future use)                                   |
| mqtt.alerts.qos                    | 1                               | MQTT _quality of service_                               |
| mqtt.alerts.retained               | true                            | Enable _retained messages_                              |
| mqtt.events.key                    | events                          | Key ID for encrypted event messages                     |
| mqtt.system.key                    | system                          | Key ID for encrypted system messages                    |
| mqtt.events.index.filepath         | \<var\>/mqtt.events.retrieved   | File for retrieved events                               |
| mqtt.permissions.enabled           | false                           | Enables/disables user permissions                       |
| mqtt.permissions.users             | \<etc\>/mqtt.permissions.users  | File containing authorised users                        |
| mqtt.permissions.groups            | \<etc\>/mqtt.permissions.groups | Permissions groups for authorised users                 |
| mqtt.cards                         | \<etc\>/mqtt/cards              | Authorised cards for remote door open                   |
| mqtt.security.HMAC.required        | false                           | Requires valid HMAC on received messages                |
| mqtt.security.HMAC.key             |                                 | HMAC key for message validation                         |
| mqtt.security.authentication       | NONE                            | Sets user authentication (none/HOTP/RSA/any)            |
| mqtt.security.hotp.range           | 8                               | Maximum discrepancy between HOTP counters               |
| mqtt.security.hotp.secrets         | \<etc\>/mqtt.hotp.secrets       | _secrets_ file for HOTP authentication                  |
| mqtt.security.hotp.counters        | \<var\>/mqtt.hotp.counters      | _counters_ file for HOTP authentication                 |
| mqtt.security.rsa.keys             | \<etc\>/mqtt/rsa                | _RSA keys_ file for RSA authentication                  |
| mqtt.security.nonce.required       | false                           | Validates message nonce field                           |
| mqtt.security.nonce.server         | \<var\>/mqtt.nonce              | _nonce_ file for server messages                        |
| mqtt.security.nonce.clients        | \<var\>/mqtt.nonce.counters     | _nonce_ file for client messages                        |
| mqtt.security.outgoing.sign        | false                           | Enables/disables outgoing message signatures            |
| mqtt.security.outgoing.encrypt     | false                           | Enables/disables outgoing message encryption            |
| mqtt.lockfile.remove               | false                           | Explicitly remove lockfile on termination               |
| mqtt.disconnects.enabled           | true                            | Reconnect if disconnected                               |
| mqtt.disconnects.interval          | 5m0s                            | Interval between reconnects if disconnected             |
| mqtt.disconnects.max               | 10                              | Maximum number of reconnects before terminating         |
| mqtt.acl.verify                    | RSA                             | ACL verification (none/any/RSA)                         |


## `AWS`

The `AWS` section defines the credentials for accessing AWS S3 for _uhppoted-app-s3_ and AWS Greengrass (_uhppoted-mqtt_).

| Parameter       | Default     | Description                                             |
|-----------------|-------------|---------------------------------------------------------|
| aws.credentials |             | AWS credentials file                                    |
| aws.profile     | default     | Profile in AWS credentials file                         |
| aws.region      | us-east-1   | AWS region                                              |


## `HTTPD`

The `HTTPD` section defines the configuration for the _uhppoted-httpd_ HTTP server.

| Parameter                              | Default                               | Description                                   |
|----------------------------------------|---------------------------------------|-----------------------------------------------|
| httpd.html                             |                                       | External folder for HTML files                |
| httpd.http.enabled                     | false                                 | Enables unsecured HTTP                        |
| httpd.http.port                        | 8080                                  | Port for unsecured HTTP                       |
| httpd.https.enabled                    | true                                  | Enables HTTPS                                 |
| httpd.https.port                       | 8443                                  | Port for HTTPS                                |
| httpd.tls.ca                           | \<etc\>/httpd/ca.cert                 | CA certificate (for client authentication)    |
| httpd.tls.certificate                  | \<etc\>/httpd/uhppoted.cert           | HTTPS server certficate                       |
| httpd.tls.key                          | \<etc\>/httpd/uhppoted.key            | HTTPS server key                              |
| httpd.tls.client.certificates.required | false                                 | Requires TLS client authentication            |
| httpd.security.auth                    | basic                                 | Enables client authorisation (_none_,_basic_) |
| httpd.security.local.db                | \<etc\>/httpd/auth.json               | JSON file contained authorised clients        |
| httpd.security.cookie.max-age          | 24                                    | Cookie Max-Age                                |
| httpd.security.login.expiry            | 1m                                    | Login cookie validity time                    |
| httpd.security.session.expiry          | 60m                                   | Session cookie validity time                  |
| httpd.security.otp.issuer              | uhppoted-httpd                        | Issuer ID for OTP application                 |
| httpd.security.otp.login               | allow                                 | Allow login with OTP                          |
| httpd.request.timeout                  | 5s                                    | Request processing time limit                 |
| httpd.system.interfaces                | \<var\>/httpd/system/interfaces.json  | JSON file for interfaces table                |
| httpd.system.controllers               | \<var\>/httpd/system/controllers.json | JSON file for controllers table               |
| httpd.system.doors                     | \<var\>/httpd/system/doors.json       | JSON file for doors table                     |
| httpd.system.groups                    | \<var\>/httpd/system/groups.json      | JSON file for groups table                    |
| httpd.system.cards                     | \<var\>/httpd/system/cards.json       | JSON file for cards table                     |
| httpd.system.events                    | \<var\>/httpd/system/events.json      | JSON file for events table                    |
| httpd.system.logs                      | \<var\>/httpd/system/logs.json        | JSON file for logs table                      |
| httpd.system.users                     | \<var\>/httpd/system/users.json       | JSON file for users table                     |
| httpd.system.history                   | \<var\>/httpd/system/history.json     | JSON file for history table                   |
| httpd.system.refresh                   | 30s                                   | System monitor refresh interval               |
| httpd.system.windows.ok                | 1m0s                                  | System monitor _healthy time window           |
| httpd.system.windows.uncertain         | 5m0s                                  | System monitor _uncertain_ window             |
| httpd.system.windows.systime           | 5m0s                                  | System monitor allowed _system time_ window   |
| httpd.system.windows.expires           | 2m0s                                  | System monitor _controller_ time window       |
| httpd.db.rules.acl                     | \<etc\>/httpd/acl.grl                 | Optional _grules_ file for ACL                |
| httpd.db.rules.interfaces              |                                       | Optional _grules_ file for interfaces table   |
| httpd.db.rules.controllers             |                                       | Optional _grules_ file for controllers table  |
| httpd.db.rules.cards                   |                                       | Optional _gruels_ file for cards table        |
| httpd.db.rules.doors                   |                                       | Optional _grules_ file for doors table        |
| httpd.db.rules.groups                  |                                       | Optional _grules_ file for groups table       |
| httpd.db.rules.events                  |                                       | Optional _grules_ file for events table       |
| httpd.db.rules.logs                    |                                       | Optional _grules_ file for logs table         |
| httpd.db.rules.users                   |                                       | Optional _grules_ file for users table        |
| httpd.audit.file                       | \<var\>/httpd/audit/audit.log         | Audit log file                                |
| httpd.retention                        | 6h0m0s                                | Retention time for deleted items              |
| httpd.timezones                        |                                       | Optional timezones definition file            |
| httpd.with-pin                         | false                                 | Includes card reader PIN in ACL               |


## `Wild Apricot`

The `Wild Apricot` section defines the configuration for _uhppoted-app-wild-apricot_.

| Parameter                              | Default         | Description                                               |
|----------------------------------------|-----------------|-----------------------------------------------------------|
| wild-apricot.http.client-timeout       | 10s             | Wild Apricot API request timeout                          |
| wild-apricot.http.retries              | 3               | Wild Apricot API request retry count                      |
| wild-apricot.http.retry-delay          | 5s              | Interval between API request retries                      |
| wild-apricot.fields.card-number        | Card Number     | Field name for card number                                |
| wild-apricot.display-order.groups      |                 | Optionally sets groups order for reports                  |
| wild-apricot.display-order.doors       |                 | Optionally sets doors order for reports                   |
| wild-apricot.facility-code             |                 | Default facility code to use for abbreviated card numbers |


## `OpenAPI`

The `OpenAPI` section defines the configuration for the optional REST _OpenAPI_ (_swagger_) server.

| Parameter         | Default    | Description                                               |
|-------------------|------------|-----------------------------------------------------------|
| openapi.enabled   | false      | Enables OpenAPI server                                    |
| openapi.directory | ./openapi  | Folder for OpenAPI YAML files                             |


