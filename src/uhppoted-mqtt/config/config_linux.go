package config

const mqttBrokerCertificate string = "/etc/uhppoted/mqtt/broker.cert"
const mqttClientCertificate string = "/etc/uhppoted/mqtt/client.cert"
const mqttClientKey string = "/etc/uhppoted/mqtt/client.key"
const users string = "/etc/uhppoted/mqtt.permissions.users"
const groups string = "/etc/uhppoted/mqtt.permissions.groups"
const hotpSecrets string = "/etc/uhppoted/mqtt.hotp.secrets"
const rsaPrivateKey string = "/etc/uhppoted/mqtt/rsa/mttd.key"
const rsaClientKeys string = "/etc/uhppoted/mqtt/rsa/clients"

const eventIDs string = "/var/uhppoted/mqtt.events.retrieved"
const hotpCounters string = "/var/uhppoted/mqtt.hotp.counters"
const rsaCounters string = "/var/uhppoted/mqtt.rsa.counters"
