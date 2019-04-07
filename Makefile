DEBUG = "--debug"
LOCAL = "192.168.1.100:51234"
CARD = "6154412"
SERIALNO = "423187757"

all: test      \
	 benchmark \
     coverage

format: 
	gofmt -w=true src/uhppote/*.go
	gofmt -w=true src/uhppote/types/*.go
	gofmt -w=true src/uhppote/messages/*.go
	gofmt -w=true src/uhppote/encoding/*.go
	gofmt -w=true src/uhppote-cli/*.go
	gofmt -w=true src/uhppote-cli/commands/*.go
	gofmt -w=true src/uhppote-simulator/*.go
	gofmt -w=true src/uhppote-simulator/simulator/*.go
	gofmt -w=true src/encoding/bcd/*.go

dist: format
	mkdir -p dist/windows
	mkdir -p dist/macosx
	mkdir -p dist/linux
	env GOOS=windows GOARCH=amd64  go build uhppote-cli;       mv uhppote-cli.exe dist/windows
	env GOOS=darwin  GOARCH=amd64  go build uhppote-cli;       mv uhppote-cli dist/macosx
	env GOOS=linux   GOARCH=amd64  go build uhppote-cli;       mv uhppote-cli dist/linux
	env GOOS=windows GOARCH=amd64  go build uhppote-simulator; mv uhppote-simulator.exe dist/windows
	env GOOS=darwin  GOARCH=amd64  go build uhppote-simulator; mv uhppote-simulator dist/macosx
	env GOOS=linux   GOARCH=amd64  go build uhppote-simulator; mv uhppote-simulator dist/linux

build: format
	go install uhppote-cli
	go install uhppote-simulator

test: build
	go clean -testcache
	go test -count=1 src/uhppote/*.go
	go test -count=1 src/uhppote/messages/*.go
	go test -count=1 src/uhppote/encoding/*.go
	go test -count=1 src/encoding/bcd/*.go

benchmark: build
	go test src/encoding/bcd/*.go -bench .

coverage: build
	go test -cover .

clean:
	go clean
	rm -rf bin

usage: build
	./bin/uhppote-cli

debug: build
	./bin/uhppote-cli --bind $(LOCAL) --debug help get-devices

help: build
	./bin/uhppote-cli --bind $(LOCAL) help

version: build
	./bin/uhppote-cli version

get-devices: build
	./bin/uhppote-cli --bind $(LOCAL) --debug get-devices

get-status: build
	./bin/uhppote-cli --bind $(LOCAL) --debug get-status 423187757

get-time: build
	./bin/uhppote-cli --bind $(LOCAL) --debug get-time 423187757

get-cards: build
	./bin/uhppote-cli --bind $(LOCAL) $(DEBUG) get-cards 423187757

get-card: build
	./bin/uhppote-cli --bind $(LOCAL) $(DEBUG) get-card $(SERIALNO) $(CARD)

get-swipes: build
	./bin/uhppote-cli --bind $(LOCAL) --debug get-swipes 423187757 1

set-time: build
	# ./bin/uhppote-cli -debug set-time 423187757 '2019-01-08 12:34:56'
	./bin/uhppote-cli -debug set-time 423187757

set-address: build
	./bin/uhppote-cli -debug set-ip-address 423187757 '192.168.1.125' '255.255.255.0' '0.0.0.0'

grant: build
	./bin/uhppote-cli --bind $(LOCAL) --debug grant 423187757 12345 2019-01-01 2019-12-31 1,4

revoke: build
	./bin/uhppote-cli --bind $(LOCAL) --debug revoke 423187757 615441

revoke-all: build
	./bin/uhppote-cli --bind $(LOCAL) --debug revoke-all 423187757

open: build
	./bin/uhppote-cli --bind $(LOCAL) --debug open 423187757 4

simulator: build
	./bin/uhppote-simulator --debug 
