all: test      \
	 benchmark \
     coverage

format: 
	gofmt -w=true src/uhppote/*.go
	gofmt -w=true src/uhppote-cli/*.go
	gofmt -w=true src/uhppote/types/*.go
	gofmt -w=true src/uhppote/messages/*.go

build: format
	go install uhppote-cli

test: build
	go test src/uhppote/messages/*.go

benchmark: build
	go test -bench .

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
	./bin/uhppote-cli -debug set-time 423187757 now



