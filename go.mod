module github.com/uhppoted/uhppoted

go 1.14

require (
	github.com/uhppoted/uhppoted-api v0.0.0-20200302181311-56c5fea77afc
	github.com/uhppoted/uhppoted-mqtt v0.0.0-20200302184633-fe8e9b264884
	github.com/uhppoted/uhppoted-rest v0.0.0-20200302184005-f30d02a22101
	golang.org/x/net v0.0.0-20200301022130-244492dfa37a // indirect
)

replace (
	github.com/uhppoted/uhppote-cli => ./uhppote-cli
	github.com/uhppoted/uhppote-core => ./uhppote-core
	github.com/uhppoted/uhppoted-api => ./uhppoted-api
	github.com/uhppoted/uhppoted-mqtt => ./uhppoted-mqtt
	github.com/uhppoted/uhppoted-rest => ./uhppoted-rest
)
