package config

import (
	"golang.org/x/sys/windows"
	"path/filepath"
)

var mqttBrokerCertificate string = filepath.Join(workdir(), "mqtt", "broker.cert")
var mqttClientCertificate string = filepath.Join(workdir(), "mqtt", "client.cert")
var mqttClientKey string = filepath.Join(workdir(), "mqtt", "client.key")
var users string = filepath.Join(workdir(), "mqtt.permissions.users")
var groups string = filepath.Join(workdir(), "mqtt.permissions.groups")
var hotpSecrets string = filepath.Join(workdir(), "mqtt.hotp.secrets")
var rsaKeyDir string = filepath.Join(workdir(), "mqtt", "rsa")

var eventIDs string = filepath.Join(workdir(), "mqtt.events.retrieved")
var hotpCounters string = filepath.Join(workdir(), "mqtt.hotp.counters")
var nonceServer string = filepath.Join(workdir(), "mqtt.nonce")
var nonceCounters string = filepath.Join(workdir(), "mqtt.nonce.counters")

func workdir() string {
	programData, err := windows.KnownFolderPath(windows.FOLDERID_ProgramData, windows.KF_FLAG_DEFAULT)
	if err != nil {
		return `C:\uhppoted`
	}

	return filepath.Join(programData, "uhppoted")
}
