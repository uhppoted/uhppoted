VERSION ?= 0.9.1
NODERED ?= 1.1.14
RELEASE ?= 
BUMP    ?= 
DEBUG   ?= --debug
DIST    ?= uhppoted_v${VERSION}

LINUX        = env GOWORK=off GOOS=linux   GOARCH=amd64       
DARWIN_X64   = env GOWORK=off GOOS=darwin  GOARCH=amd64        
DARWIN_ARM64 = env GOWORK=off GOOS=darwin  GOARCH=arm64        
WINDOWS      = env GOWORK=off GOOS=windows GOARCH=amd64         
ARM          = env GOWORK=off GOOS=linux   GOARCH=arm64     
ARM7         = env GOWORK=off GOOS=linux   GOARCH=arm GOARM=7
ARM6         = env GOWORK=off GOOS=linux   GOARCH=arm GOARM=6

.PHONY: debug
.PHONY: docker
.PHONY: simulator
.PHONY: integration-tests

all: test      \
	 benchmark \
     coverage

clean:
	go clean
	rm -rf bin

status-all:
	git -C uhppote-core              status
	git -C uhppoted-lib              status
	git -C uhppote-simulator         status
	git -C uhppoted-cli               status
	git -C uhppoted-rest             status
	git -C uhppoted-mqtt             status
	git -C uhppoted-httpd            status
	git -C uhppoted-tunnel           status
	git -C uhppoted-dll              status
	git -C uhppoted-codegen          status
	git -C uhppoted-lib-go           status
	git -C uhppoted-lib-nodejs       status
	git -C uhppoted-lib-python       status
	git -C uhppoted-app-s3           status
	git -C uhppoted-app-sheets       status
	git -C uhppoted-app-wild-apricot status
	git -C uhppoted-app-db           status
	git -C uhppoted-wiegand          status
	git                              status

update:
	@echo "Nothing to do"

update-release:
	@echo "Nothing to do"

format: 
	@echo "Nothing to do"

build: format
	@echo "Nothing to do"

test: build
	@echo "Nothing to do"

vet: build
	echo "Nothing to do"

lint: build
	echo "Nothing to do"

benchmark: build
	echo "Nothing to do"

coverage: build
	echo "Nothing to do"

integration-tests: 
	@echo "Nothing to do"

build-all: test vet
	@echo "Nothing to do"

release: update-release build-all docker
	echo "Nothing to do"

publish: release
	echo "Releasing version $(VERSION)"
# 	gh release create "$(VERSION)" ./dist/*.tar.gz ./dist/*.zip --draft --prerelease --title "$(VERSION)-beta" --notes-file release-notes.md

# e.g.
# make release-all
# make release-all RELEASE=--release
# make release-all BUMP=--bump
release-all: 
	echo "Releasing version $(VERSION)"
# 	. .venv/bin/activate; yapf -ri ./internal
# 	. .venv/bin/activate; python3 ./internal/release.py --versions=.versions $(RELEASE) $(BUMP)
# 	python ./internal/release.py --version=$(VERSION) --node-red=$(NODERED) $(RELEASE) $(BUMP)

build-github: 
	echo "Nothing to do"

debug: 
	. .venv/bin/activate; yapf -ri ./internal
	. .venv/bin/activate; python ./internal/debug.py

swagger: 
	docker run --detach --publish 80:8080 --name swagger --rm swaggerapi/swagger-editor 
	sleep 1
	open http://127.0.0.1:80

docker:
	cd uhppote-simulator && make docker-dev
	cd uhppoted-rest     && make docker-dev
	cd uhppoted-mqtt     && make docker-dev
	cd ./docker/hivemq   && docker build --no-cache -f Dockerfile -t hivemq/uhppoted    .

docker-clean:
	docker image     prune -f
	docker container prune -f

docker-build-all: docker-clean docker

docker-simulator:
	docker run --detach --publish 8000:8000 --publish 60000:60000 --publish 60000:60000/udp --name simulator --rm uhppoted/simulator-dev
	sleep 1
	./uhppoted-cli/bin/uhppoted-cli --debug set-listener 405419896 192.168.1.125:60001
	./uhppoted-cli/bin/uhppoted-cli --debug set-listener 303986753 192.168.1.125:60001
	./uhppoted-cli/bin/uhppoted-cli --debug set-listener 201020304 192.168.1.125:60001

docker-simulator-tunnel:
	docker run --detach --publish 8000:8000 --publish 60005:60000 --publish 60005:60000/udp --name simulator --rm uhppoted/simulator-dev

docker-rest:
	docker run --detach --publish 8080:8080 --name restd --rm uhppoted/rest

docker-mqtt:
	docker run --detach --name mqttd --rm uhppoted/uhppoted-mqtt-dev

docker-hivemq:
	docker run --detach --publish 8081:8080 --publish 1883:1883 --publish 8883:8883 --name hivemq --rm hivemq/uhppoted

docker-sql-server:
#	docker run -d --name sqld -e 'ACCEPT_EULA=Y' -e 'SA_PASSWORD=UBxNxrQiKWsjncow7mMx' -p 1433:1433 mcr.microsoft.com/mssql/server:2022-latest
	docker start sqld

docker-sql-server-cli:
	mssql-cli -U sa -P UBxNxrQiKWsjncow7mMx 

docker-mysql:
#	docker run --name mysqld -e MYSQL_ROOT_PASSWORD=password --publish 3306:3306 -d mysql:latest
	docker container start mysqld

docker-mysql-cli:
	docker exec -it mysqld bash

docker-stop:
	docker stop $$(docker container ls -q)

hivemq-listen:
	# mqtt subscribe --topic 'uhppoted/reply/#' | jq '.' 
	mqtt subscribe --topic 'uhppoted/#' | jq '.' 

hivemq-listen-events:
	mqtt subscribe --topic 'uhppoted/gateway/events' | jq '.' 

hivemq-listen-live-events:
	mqtt subscribe --topic 'uhppoted/gateway/events/live' | jq '.' 

hivemq-listen-all-events:
	mqtt subscribe --topic 'uhppoted/gateway/events/#' | jq '.' 



