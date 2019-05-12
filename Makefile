CLI = ./bin/uhppote-cli
DEBUG = --debug
LOCAL = 192.168.1.100:51234
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
	gofmt -w=true src/uhppote-cli/config/*.go
	gofmt -w=true src/uhppote-simulator/*.go
	gofmt -w=true src/uhppote-simulator/simulator/*.go
	gofmt -w=true src/encoding/bcd/*.go

release: format
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
	$(CLI)

debug: build
#	$(CLI) $(DEBUG) --bind 192.168.0.14:12345 --broadcast 192.168.0.255:60000 get-devices
#	$(CLI) $(DEBUG) --bind 0.0.0.0:12345                                      get-devices
#	$(CLI) $(DEBUG) --bind 0.0.0.0:12345 get-card $(SERIALNO) $(CARD)
	$(CLI) $(DEBUG) get-card $(SERIALNO) $(CARD)

help: build
	$(CLI) --bind $(LOCAL) help

version: build
	$(CLI) version

run: build
	$(CLI) --bind $(LOCAL) $(DEBUG) get-devices
#	$(CLI) --bind $(LOCAL) $(DEBUG) set-address    $(SERIALNO) '192.168.1.125' '255.255.255.0' '0.0.0.0'
	$(CLI) --bind $(LOCAL) $(DEBUG) get-cards      $(SERIALNO)
	$(CLI) --bind $(LOCAL) $(DEBUG) get-door-delay $(SERIALNO) $(DOOR)
	$(CLI) --bind $(LOCAL) $(DEBUG) set-time       $(SERIALNO)
	$(CLI) --bind $(LOCAL) $(DEBUG) revoke         $(SERIALNO) $(CARD)

get-devices: build
#	$(CLI) --bind 0.0.0.0:0 $(DEBUG) get-devices
	$(CLI) --bind $(LOCAL) $(DEBUG) get-devices

set-address: build
	$(CLI) -bind $(LOCAL) $(DEBUG) set-address $(SERIALNO) '192.168.1.125' '255.255.255.0' '0.0.0.0'

get-status: build
	$(CLI) --bind $(LOCAL) $(DEBUG) get-status $(SERIALNO)

get-time: build
	$(CLI) --bind $(LOCAL) $(DEBUG) get-time $(SERIALNO)

set-time: build
	# $(CLI) -debug set-time 423187757 '2019-01-08 12:34:56'
	$(CLI) --bind $(LOCAL) $(DEBUG) set-time $(SERIALNO)

get-door-delay: build
	$(CLI) --bind $(LOCAL) $(DEBUG) get-door-delay $(SERIALNO) $(DOOR)

set-door-delay: build
	$(CLI) --bind $(LOCAL) $(DEBUG) set-door-delay $(SERIALNO) $(DOOR) 5

get-listener: build
	$(CLI) --bind  $(DEBUG) get-listener $(SERIALNO)

set-listener: build
	$(CLI) --bind $(LOCAL) $(DEBUG) set-listener $(SERIALNO) 192.168.1.100:40000

get-cards: build
	$(CLI) --bind $(LOCAL) $(DEBUG) get-cards $(SERIALNO)

get-card: build
	$(CLI) --bind $(LOCAL) $(DEBUG) get-card $(SERIALNO) $(CARD)

grant: build
	$(CLI) --bind $(LOCAL) $(DEBUG) grant $(SERIALNO) $(CARD) 2019-01-01 2019-12-31 1

revoke: build
	$(CLI) --bind $(LOCAL) $(DEBUG) revoke $(SERIALNO) $(CARD)

revoke-all: build
	$(CLI) --bind $(LOCAL) $(DEBUG) revoke-all $(SERIALNO)

get-events: build
	$(CLI) --bind $(LOCAL) $(DEBUG) get-events $(SERIALNO)

get-event-index: build
	$(CLI) --bind $(LOCAL) $(DEBUG) get-event-index $(SERIALNO)

set-events-index: build
	$(CLI) --bind $(LOCAL) $(DEBUG) set-events-index $(SERIALNO) 23

open: build
	$(CLI) --bind $(LOCAL) $(DEBUG) open $(SERIALNO) 1

listen: build
	$(CLI) --bind 192.168.1.100:40000 $(DEBUG) listen 

simulator: build
	./bin/uhppote-simulator --debug 
