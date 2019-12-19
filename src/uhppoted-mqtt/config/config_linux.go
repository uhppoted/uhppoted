package config

const secrets string = "/etc/uhppoted/mqtt.hotp.secrets"
const users string = "/etc/uhppoted/mqtt.permissions.users"
const groups string = "/etc/uhppoted/mqtt.permissions.groups"
const brokerCertificate string = "/etc/uhppoted/mqtt.broker.pem"
const clientCertificate string = "/etc/uhppoted/mqtt-client.cert"
const clientKey string = "/etc/uhppoted/mqtt-client.key"

const counters string = "/var/uhppoted/mqtt.hotp.counters"
const eventIDs string = "/var/uhppoted/mqtt.events.retrieved"
