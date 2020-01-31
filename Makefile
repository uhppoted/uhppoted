CLI = ./bin/uhppote-cli
SIMULATOR = ./bin/uhppote-simulator
DEBUG ?= --debug
LOCAL ?= 192.168.1.100:51234
CARD ?= 6154410
SERIALNO ?= 423187757
DOOR ?= 3
DIST ?= development
DATETIME = $(shell date "+%Y-%m-%d %H:%M:%S")
VERSION = v0.5.x
LDFLAGS = -ldflags "-X uhppote.VERSION=$(VERSION)" 

all: test      \
	 benchmark \
     coverage

clean:
	go clean
	rm -rf bin

format: 
	go fmt uhppote...

build: format
	go install $(LDFLAGS) uhppote-cli
	go install $(LDFLAGS) uhppoted-rest
	go install $(LDFLAGS) cmd/uhppoted-mqtt
	go install $(LDFLAGS) uhppote-simulator

test: build
	go test uhppote...

vet: build
	go vet uhppote...

lint: build
	golint uhppote...

benchmark: build
	go test -bench uhppote...

coverage: build
	go test -cover uhppote...

integration-tests: build
	go fmt src/integration-tests/cli/*.go
	go fmt src/integration-tests/mqttd/*.go
#	go test src/integration-tests/cli/*.go
	go clean -testcache && go test -count=1 src/integration-tests/mqttd/*.go

release: test vet
	mkdir -p dist/$(DIST)/windows
	mkdir -p dist/$(DIST)/darwin
	mkdir -p dist/$(DIST)/linux
	mkdir -p dist/$(DIST)/openapi
	env GOOS=linux   GOARCH=amd64  go build -o dist/$(DIST)/linux/uhppote-cli             uhppote-cli
	env GOOS=darwin  GOARCH=amd64  go build -o dist/$(DIST)/darwin/uhppote-cli            uhppote-cli
	env GOOS=windows GOARCH=amd64  go build -o dist/$(DIST)/windows/uhppote-cli.exe       uhppote-cli
	env GOOS=linux   GOARCH=amd64  go build -o dist/$(DIST)/linux/uhppoted-rest           uhppoted-rest
	env GOOS=darwin  GOARCH=amd64  go build -o dist/$(DIST)/darwin/uhppoted-rest          uhppoted-rest
	env GOOS=windows GOARCH=amd64  go build -o dist/$(DIST)/windows/uhppoted-rest.exe     uhppoted-rest
	env GOOS=linux   GOARCH=amd64  go build -o dist/$(DIST)/linux/uhppoted-mqtt           cmd/uhppoted-mqtt
	env GOOS=darwin  GOARCH=amd64  go build -o dist/$(DIST)/darwin/uhppoted-mqtt          cmd/uhppoted-mqtt
	env GOOS=windows GOARCH=amd64  go build -o dist/$(DIST)/windows/uhppoted-mqtt.exe     cmd/uhppoted-mqtt
	env GOOS=linux   GOARCH=amd64  go build -o dist/$(DIST)/linux/uhppote-simulator       uhppote-simulator
	env GOOS=darwin  GOARCH=amd64  go build -o dist/$(DIST)/darwin/uhppote-simulator      uhppote-simulator
	env GOOS=windows GOARCH=amd64  go build -o dist/$(DIST)/windows/uhppote-simulator.exe uhppote-simulator
	cp -r install/openapi/* dist/$(DIST)/openapi/

release-tar: release
	tar --directory=dist --exclude=".DS_Store" -cvzf dist/$(DIST).tar.gz $(DIST)

debug: build
#	sudo ./bin/uhppoted-mqtt daemonize
	go test -v src/uhppote/encoding/conf/*.go

simulator: build
	./bin/uhppote-simulator --debug --bind 192.168.1.100:54321 --rest 192.168.1.100:8008 --devices "./runtime/simulation/devices"

simulator-device: build
	./bin/uhppote-simulator --debug --devices "runtime/simulation/devices" new-device 678

uhppoted-rest: build
	./bin/uhppoted-rest --console --debug 

uhppoted-mqtt: build
	./bin/uhppoted-mqtt --console

swagger: 
	docker run --detach --publish 80:8080 --rm swaggerapi/swagger-editor 
	open http://127.0.0.1:80

docker: build
	env GOOS=linux GOARCH=amd64 go build -o docker/simulator/uhppote-simulator     uhppote-simulator
	env GOOS=linux GOARCH=amd64 go build -o docker/uhppoted-rest/uhppote-simulator uhppote-simulator
	env GOOS=linux GOARCH=amd64 go build -o docker/uhppoted-rest/uhppoted-rest     uhppoted-rest
	env GOOS=linux GOARCH=amd64 go build -o docker/integration-tests/simulator/uhppote-simulator uhppote-simulator
	docker image     prune -f
	docker container prune -f
	docker build -f ./docker/simulator/Dockerfile     -t simulator       . 
	docker build -f ./docker/uhppoted-rest/Dockerfile -t uhppoted        . 
	docker build -f ./docker/hivemq/Dockerfile        -t hivemq/uhppoted . 
	docker build -f ./docker/integration-tests/simulator/Dockerfile -t integration-tests/simulator . 

docker-simulator:
	docker run --detach --publish 8000:8000 --publish 60000:60000/udp --name simulator --rm simulator
	./bin/uhppote-cli --debug set-listener 405419896 192.168.1.100:60001
	./bin/uhppote-cli --debug set-listener 303986753 192.168.1.100:60001

docker-hivemq:
	docker run --detach --publish 8081:8080 --publish 1883:1883 --publish 8883:8883 --name hivemq --rm hivemq/uhppoted

docker-rest:
	docker run --detach --publish 8080:8080 --rm uhppoted

docker-stop:
	docker stop simulator
	docker stop hivemq

docker-integration-tests:
	docker run --detach --publish 8000:8000 --publish 60000:60000/udp --name qwerty --rm integration-tests/simulator

hivemq-listen:
#	mqtt subscribe --topic 'twystd/uhppoted/#'
	open runtime/mqtt-spy-0.5.4-jar-with-dependencies.jar



