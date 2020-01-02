package config

const secrets string = "/etc/uhppoted/mqtt.hotp.secrets"
const users string = "/etc/uhppoted/mqtt.permissions.users"
const groups string = "/etc/uhppoted/mqtt.permissions.groups"
const mqttBrokerCertificate string = "/etc/uhppoted/mqtt/broker.cert"
const mqttClientCertificate string = "/etc/uhppoted/mqtt/client.cert"
const mqttClientKey string = "/etc/uhppoted/mqtt/client.key"
const rsaClientKeys string = "/etc/uhppoted/mqtt/rsa/clients"

const counters string = "/var/uhppoted/mqtt.hotp.counters"
const eventIDs string = "/var/uhppoted/mqtt.events.retrieved"
