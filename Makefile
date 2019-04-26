DEBUG = "--debug"
LOCAL = "192.168.1.100:51234"
CARD = 6154412
SERIALNO = 423187757
DOOR = 3

all: test      \
	 benchmark \
     coverage

format: 
	gofmt -w=true src/uhppote/*.go
	gofmt -w=true src/uhppote/types/*.go
	gofmt -w=true src/uhppote/encoding/UTO311-L0x/*.go
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
	go test -count=1 src/uhppote/encoding/UTO311-L0x/*.go
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
	./bin/uhppote-cli --bind $(LOCAL) --debug help get-status

help: build
	./bin/uhppote-cli --bind $(LOCAL) help

version: build
	./bin/uhppote-cli version

run: build
	./bin/uhppote-cli --bind $(LOCAL) $(DEBUG) get-devices
	./bin/uhppote-cli --bind $(LOCAL) $(DEBUG) set-address $(SERIALNO) '192.168.1.125' '255.255.255.0' '0.0.0.0'
	./bin/uhppote-cli --bind $(LOCAL) $(DEBUG) get-cards   $(SERIALNO)

get-devices: build
	./bin/uhppote-cli --bind $(LOCAL) $(DEBUG) get-devices

set-address: build
	./bin/uhppote-cli -bind $(LOCAL) $(DEBUG) set-address $(SERIALNO) '192.168.1.125' '255.255.255.0' '0.0.0.0'

get-status: build
	./bin/uhppote-cli --bind $(LOCAL) --debug get-status $(SERIALNO)

get-time: build
	./bin/uhppote-cli --bind $(LOCAL) --debug get-time $(SERIALNO)

set-time: build
	# ./bin/uhppote-cli -debug set-time 423187757 '2019-01-08 12:34:56'
	./bin/uhppote-cli --bind $(LOCAL) $(DEBUG) set-time $(SERIALNO)

get-door-delay: build
	./bin/uhppote-cli --bind $(LOCAL) $(DEBUG) get-door-delay $(SERIALNO) $(DOOR)

set-door-delay: build
	./bin/uhppote-cli --bind $(LOCAL) $(DEBUG) set-door-delay $(SERIALNO) $(DOOR) 5

get-listener: build
	./bin/uhppote-cli --bind $(LOCAL) $(DEBUG) get-listener $(SERIALNO)

set-listener: build
	./bin/uhppote-cli --bind $(LOCAL) $(DEBUG) set-listener $(SERIALNO) 192.168.1.100:40000

get-cards: build
	./bin/uhppote-cli --bind $(LOCAL) $(DEBUG) get-cards $(SERIALNO)

get-card: build
	./bin/uhppote-cli --bind $(LOCAL) $(DEBUG) get-card $(SERIALNO) $(CARD)

grant: build
	./bin/uhppote-cli --bind $(LOCAL) --debug grant $(SERIALNO) $(CARD) 2019-01-01 2019-12-31 1

revoke: build
	./bin/uhppote-cli --bind $(LOCAL) --debug revoke $(SERIALNO) $(CARD)

revoke-all: build
	./bin/uhppote-cli --bind $(LOCAL) --debug revoke-all $(SERIALNO)

get-events: build
	./bin/uhppote-cli --bind $(LOCAL) get-events $(SERIALNO)

get-event-index: build
	./bin/uhppote-cli --bind $(LOCAL) $(DEBUG) get-event-index $(SERIALNO)

set-events-index: build
	./bin/uhppote-cli --bind $(LOCAL) $(DEBUG) set-events-index $(SERIALNO) 23

open: build
	./bin/uhppote-cli --bind $(LOCAL) --debug open $(SERIALNO) 1

simulator: build
	./bin/uhppote-simulator --debug 
