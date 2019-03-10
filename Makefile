LOCAL = "192.168.1.100:51234"

all: test      \
	 benchmark \
     coverage

format: 
	gofmt -w=true src/uhppote/*.go
	gofmt -w=true src/uhppote-cli/*.go
	gofmt -w=true src/uhppote-cli/commands/*.go
	gofmt -w=true src/uhppote-simulator/*.go
	gofmt -w=true src/uhppote/types/*.go
	gofmt -w=true src/uhppote/messages/*.go
	gofmt -w=true src/encoding/bcd/*.go

build: format
	go install uhppote-cli
	go install uhppote-simulator

test: build
	go test src/uhppote/messages/*.go
	go test src/encoding/bcd/*.go

benchmark: build
	go test src/encoding/bcd/*.go -bench .

coverage: build
	go test -cover .

clean:
	go clean
	rm -rf bin

usage: build
	./bin/uhppote-cli

help: build
	./bin/uhppote-cli --bind $(LOCAL) help

version: build
	./bin/uhppote-cli version

list-devices: build
	./bin/uhppote-cli -debug list-devices

get-status: build
	./bin/uhppote-cli --debug --bind $(LOCAL) get-status 423187757

get-time: build
	./bin/uhppote-cli -debug get-time 423187757

set-time: build
	# ./bin/uhppote-cli -debug set-time 423187757 '2019-01-08 12:34:56'
	./bin/uhppote-cli -debug set-time 423187757

set-address: build
	./bin/uhppote-cli -debug set-ip-address 423187757 '192.168.1.125' '255.255.255.0' '0.0.0.0'

get-authorised: build
	./bin/uhppote-cli -debug get-authorised 423187757

authorise: build
	./bin/uhppote-cli -debug authorise 423187757 12345   2019-01-01 2019-12-31 1,4

get-swipe: build
	./bin/uhppote-cli --debug get-swipe 423187757 1

simulator: build
	./bin/uhppote-simulator
