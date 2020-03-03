module github.com/uhppoted/uhppoted

go 1.14

replace (
	github.com/uhppoted/uhppote-cli => ./uhppote-cli
	github.com/uhppoted/uhppote-core => ./uhppote-core
	github.com/uhppoted/uhppoted-api => ./uhppoted-api
	github.com/uhppoted/uhppoted-mqtt => ./uhppoted-mqtt
	github.com/uhppoted/uhppoted-rest => ./uhppoted-rest
)
