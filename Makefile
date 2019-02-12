all: test      \
	 benchmark \
     coverage

format: 
	gofmt -w=true src/uhppote/*.go
	gofmt -w=true src/uhppote-cli/*.go
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

search: build
	./bin/uhppote-cli -debug search

get-time: build
	./bin/uhppote-cli -debug get-time 423187757

set-time: build
	./bin/uhppote-cli -debug set-time 423187757 '2019-01-08 12:34:56'

set-address: build
	./bin/uhppote-cli -debug set-ip-address 423187757 '192.168.1.150' '255.255.254.0' '0.0.0.0'

get-auth-rec: build
	./bin/uhppote-cli -debug get-auth-rec 423187757

authorise: build
	./bin/uhppote-cli -debug add-auth 423187757 12345 2019-01-01 2019-12-31 1,4

simulator: build
	./bin/uhppote-simulator
