package config

const mqttBrokerCertificate string = "/etc/uhppoted/mqtt/broker.cert"
const mqttClientCertificate string = "/etc/uhppoted/mqtt/client.cert"
const mqttClientKey string = "/etc/uhppoted/mqtt/client.key"
const users string = "/etc/uhppoted/mqtt.permissions.users"
const groups string = "/etc/uhppoted/mqtt.permissions.groups"
const hotpSecrets string = "/etc/uhppoted/mqtt.hotp.secrets"
const rsaKeyDir string = "/etc/uhppoted/mqtt/rsa"

const eventIDs string = "/var/uhppoted/mqtt.events.retrieved"
const hotpCounters string = "/var/uhppoted/mqtt.hotp.counters"
const nonceCounters string = "/var/uhppoted/mqtt.nonce.counters"
