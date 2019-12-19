package config

import (
	"golang.org/x/sys/windows"
	"path/filepath"
)

var secrets string = filepath.Join(workdir(), "mqtt.hotp.secrets")
var counters string = filepath.Join(workdir(), "mqtt.hotp.counters")
var users string = filepath.Join(workdir(), "mqtt.permissions.users")
var groups string = filepath.Join(workdir(), "mqtt.permissions.groups")
var eventIDs string = filepath.Join(workdir(), "mqtt.events.retrieved")
var brokerCertificate string = filepath.Join(workdir(), "mqtt.broker.pem")
var clientCertificate string = filepath.Join(workdir(), "mqtt-client.cert")
var clientKey string = filepath.Join(workdir(), "mqtt-client.key")

func workdir() string {
	programData, err := windows.KnownFolderPath(windows.FOLDERID_ProgramData, windows.KF_FLAG_DEFAULT)
	if err != nil {
		return `C:\uhppoted`
	}

	return filepath.Join(programData, "uhppoted")
}
