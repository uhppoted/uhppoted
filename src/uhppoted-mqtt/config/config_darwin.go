package config

const mqttBrokerCertificate string = "/usr/local/etc/com.github.twystd.uhppoted/mqtt/broker.cert"
const mqttClientCertificate string = "/usr/local/etc/com.github.twystd.uhppoted/mqtt/client.cert"
const mqttClientKey string = "/usr/local/etc/com.github.twystd.uhppoted/mqtt/client.key"
const users string = "/usr/local/etc/com.github.twystd.uhppoted/mqtt.permissions.users"
const groups string = "/usr/local/etc/com.github.twystd.uhppoted/mqtt.permissions.groups"
const hotpSecrets string = "/usr/local/etc/com.github.twystd.uhppoted/mqtt.hotp.secrets"
const rsaKeyDir string = "/usr/local/etc/com.github.twystd.uhppoted/mqtt/rsa"

const eventIDs string = "/usr/local/var/com.github.twystd.uhppoted/mqtt.events.retrieved"
const hotpCounters string = "/usr/local/var/com.github.twystd.uhppoted/mqtt.hotp.counters"
const nonceCounters string = "/usr/local/var/com.github.twystd.uhppoted/mqtt.nonce.counters"
