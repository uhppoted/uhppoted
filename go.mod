module github.com/uhppoted/uhppoted

go 1.14

require (
	github.com/uhppoted/uhppoted-api v0.0.0-20200303183137-339d1456290b
	github.com/uhppoted/uhppoted-mqtt v0.0.0-20200302184633-fe8e9b264884
	golang.org/x/net v0.0.0-20200301022130-244492dfa37a // indirect
	golang.org/x/sys v0.0.0-20200302150141-5c8b2ff67527 // indirect
)

replace (
	github.com/uhppoted/uhppote-cli => ./uhppote-cli
	github.com/uhppoted/uhppote-core => ./uhppote-core
	github.com/uhppoted/uhppoted-api => ./uhppoted-api
	github.com/uhppoted/uhppoted-mqtt => ./uhppoted-mqtt
	github.com/uhppoted/uhppoted-rest => ./uhppoted-rest
)
