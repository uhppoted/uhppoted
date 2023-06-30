RELEASE = 0.8.4
DEBUG ?= --debug
DIST  ?= development

LINUX   = env GOWORK=off GOOS=linux   GOARCH=amd64       
DARWIN  = env GOWORK=off GOOS=darwin  GOARCH=amd64        
WINDOWS = env GOWORK=off GOOS=windows GOARCH=amd64         
ARM     = env GOWORK=off GOOS=linux   GOARCH=arm64     
ARM7    = env GOWORK=off GOOS=linux   GOARCH=arm GOARM=7

.PHONY: debug
.PHONY: docker
.PHONY: simulator
.PHONY: uhppoted-rest
.PHONY: uhppoted-mqtt
.PHONY: uhppoted-app-s3
.PHONY: uhppoted-app-sheets
.PHONY: uhppoted-app-wild-apricot
.PHONY: uhppoted-app-db
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
	git -C uhppote-cli               status
	git -C uhppoted-rest             status
	git -C uhppoted-mqtt             status
	git -C uhppoted-httpd            status
	git -C uhppoted-tunnel           status
	git -C uhppoted-dll              status
	git -C uhppoted-codegen          status
	git -C uhppoted-app-s3           status
	git -C uhppoted-app-sheets       status
	git -C uhppoted-app-wild-apricot status
	git -C uhppoted-app-wild-db      status
	git -C uhppoted-wiegand          status
	git                              status

push-all:
	git push --recurse-submodules=on-demand 
#	git -C uhppote-core              push
#	git -C uhppoted-lib              push
#	git -C uhppote-simulator         push
#	git -C uhppote-cli               push
#	git -C uhppoted-rest             push
#	git -C uhppoted-mqtt             push
#	git -C uhppoted-app-s3           push
#	git -C uhppoted-app-sheets       push
#	git -C uhppoted-app-wild-apricot push
#	git -C uhppoted-app-db           push
#	git -C uhppoted-httpd            push
#	git -C uhppoted-tunnel           push
#	git -C uhppoted-dll              push
#	git -C uhppoted-codegen          push
#	git -C uhppoted-wiegand          push
#	git                              push

update:
	go get -u github.com/uhppoted/uhppote-core@master
	go mod tidy

update-release:
	go get -u github.com/uhppoted/uhppote-core
	go mod tidy


format: 
	cd uhppote-core              && make format
	cd uhppoted-lib              && make format
	cd uhppote-simulator         && make format
	cd uhppote-cli               && make format
	cd uhppoted-rest             && make format
	cd uhppoted-mqtt             && make format
	cd uhppoted-httpd            && make format
	cd uhppoted-tunnel           && make format
	cd uhppoted-dll              && make format
	cd uhppoted-codegen          && make format
	cd uhppoted-app-s3           && make format
	cd uhppoted-app-sheets       && make format
	cd uhppoted-app-wild-apricot && make format
	cd uhppoted-app-db           && make format
	cd integration-tests         && go fmt ./...
	go fmt ./...

build: format
	mkdir -p bin
	cd uhppote-core              && go build -trimpath            ./...
	cd uhppoted-lib              && go build -trimpath            ./...
	cd uhppote-simulator         && go build -trimpath -o ../bin/ ./...
	cd uhppote-cli               && go build -trimpath -o ../bin/ ./...
	cd uhppoted-rest             && go build -trimpath -o ../bin/ ./...
	cd uhppoted-mqtt             && go build -trimpath -o ../bin/ ./...
	cd uhppoted-httpd            && go build -trimpath -o ../bin/ ./...
	cd uhppoted-tunnel           && go build -trimpath -o ../bin/ ./...
#	cd uhppoted-dll              && make build
	cd uhppoted-codegen          && go build -trimpath -o ../bin/ ./...
	cd uhppoted-app-s3           && go build -trimpath -o ../bin/ ./...
	cd uhppoted-app-sheets       && go build -trimpath -o ../bin/ ./...
	cd uhppoted-app-wild-apricot && go build -trimpath -o ../bin/ ./...
	cd uhppoted-app-db           && go build -trimpath -o ../bin/ ./...

test: build
	cd uhppote-core              && go test ./...
	cd uhppoted-lib              && go test ./...
	cd uhppote-simulator         && go test ./...
	cd uhppote-cli               && go test ./...
	cd uhppoted-rest             && go test ./...
	cd uhppoted-mqtt             && go test ./...
	cd uhppoted-httpd            && go test  -tags "tests" ./...
	cd uhppoted-tunnel           && go test ./...
#	cd uhppoted-dll              && make tests
	cd uhppoted-codegen          && go test ./...
	cd uhppoted-app-s3           && go test ./...
	cd uhppoted-app-sheets       && go test ./...
	cd uhppoted-app-wild-apricot && go test ./...
	cd uhppoted-app-db           && go test ./...

vet: build
	go vet ./...

lint: build
	golint ./...

benchmark: build
	go test -bench ./...

coverage: build
	go test -cover ./...

integration-tests: 
	cd integration-tests && go fmt ./...
#	go test integration-tests/cli/*.go
	go test -v integration-tests/mqttd/*.go
	# go clean -testcache && go test -count=1 integration-tests/simulator/*.go

build-all: test vet
	mkdir -p dist/linux/$(DIST)
	mkdir -p dist/darwin/$(DIST)
	mkdir -p dist/windows/$(DIST)
	mkdir -p dist/openapi/$(DIST)
	mkdir -p dist/arm/$(DIST)
	mkdir -p dist/arm7/$(DIST)

	cd uhppote-cli         && $(LINUX)   go build -trimpath -o ../dist/linux/$(DIST)   ./...
	cd uhppote-cli         && $(DARWIN)  go build -trimpath -o ../dist/darwin/$(DIST)  ./...
	cd uhppote-cli         && $(WINDOWS) go build -trimpath -o ../dist/windows/$(DIST) ./...
	cd uhppote-cli         && $(ARM)     go build -trimpath -o ../dist/arm/$(DIST)     ./...
	cd uhppote-cli         && $(ARM7)    go build -trimpath -o ../dist/arm7/$(DIST)    ./...

	cd uhppoted-rest       && $(LINUX)   go build -trimpath -o ../dist/linux/$(DIST)   ./...
	cd uhppoted-rest       && $(DARWIN)  go build -trimpath -o ../dist/darwin/$(DIST)  ./...
	cd uhppoted-rest       && $(WINDOWS) go build -trimpath -o ../dist/windows/$(DIST) ./...
	cd uhppoted-rest       && $(ARM)     go build -trimpath -o ../dist/arm/$(DIST)     ./...
	cd uhppoted-rest       && $(ARM7)    go build -trimpath -o ../dist/arm7/$(DIST)    ./...

	cd uhppoted-mqtt       && $(LINUX)   go build -trimpath -o ../dist/linux/$(DIST)   ./...
	cd uhppoted-mqtt       && $(DARWIN)  go build -trimpath -o ../dist/darwin/$(DIST)  ./...
	cd uhppoted-mqtt       && $(WINDOWS) go build -trimpath -o ../dist/windows/$(DIST) ./...	
	cd uhppoted-mqtt       && $(ARM)     go build -trimpath -o ../dist/arm/$(DIST)     ./...
	cd uhppoted-mqtt       && $(ARM7)    go build -trimpath -o ../dist/arm7/$(DIST)    ./...

	cd uhppoted-httpd    && $(LINUX)   go build -trimpath -o ../dist/linux/$(DIST)   ./...
	cd uhppoted-httpd    && $(DARWIN)  go build -trimpath -o ../dist/darwin/$(DIST)  ./...
	cd uhppoted-httpd    && $(WINDOWS) go build -trimpath -o ../dist/windows/$(DIST) ./...	
	cd uhppoted-httpd    && $(ARM)     go build -trimpath -o ../dist/arm/$(DIST)     ./...
	cd uhppoted-httpd    && $(ARM7)    go build -trimpath -o ../dist/arm7/$(DIST)    ./...

	cd uhppoted-tunnel   && $(LINUX)   go build -trimpath -o ../dist/linux/$(DIST)   ./...
	cd uhppoted-tunnel   && $(DARWIN)  go build -trimpath -o ../dist/darwin/$(DIST)  ./...
	cd uhppoted-tunnel   && $(WINDOWS) go build -trimpath -o ../dist/windows/$(DIST) ./...	
	cd uhppoted-tunnel   && $(ARM)     go build -trimpath -o ../dist/arm/$(DIST)     ./...
	cd uhppoted-tunnel   && $(ARM7)    go build -trimpath -o ../dist/arm7/$(DIST)    ./...

	cd uhppoted-codegen  && $(LINUX)   go build -trimpath -o ../dist/linux/$(DIST)   ./...
	cd uhppoted-codegen  && $(DARWIN)  go build -trimpath -o ../dist/darwin/$(DIST)  ./...
	cd uhppoted-codegen  && $(WINDOWS) go build -trimpath -o ../dist/windows/$(DIST) ./...	
	cd uhppoted-codegen  && $(ARM)     go build -trimpath -o ../dist/arm/$(DIST)     ./...
	cd uhppoted-codegen  && $(ARM7)    go build -trimpath -o ../dist/arm7/$(DIST)    ./...

	cd uhppoted-app-s3     && $(LINUX)   go build -trimpath -o ../dist/linux/$(DIST)   ./...
	cd uhppoted-app-s3     && $(DARWIN)  go build -trimpath -o ../dist/darwin/$(DIST)  ./...
	cd uhppoted-app-s3     && $(WINDOWS) go build -trimpath -o ../dist/windows/$(DIST) ./...	
	cd uhppoted-app-s3     && $(ARM)     go build -trimpath -o ../dist/arm/$(DIST)     ./...
	cd uhppoted-app-s3     && $(ARM7)    go build -trimpath -o ../dist/arm7/$(DIST)    ./...

	cd uhppoted-app-sheets && $(LINUX)   go build -trimpath -o ../dist/linux/$(DIST)   ./...
	cd uhppoted-app-sheets && $(DARWIN)  go build -trimpath -o ../dist/darwin/$(DIST)  ./...
	cd uhppoted-app-sheets && $(WINDOWS) go build -trimpath -o ../dist/windows/$(DIST) ./...	
	cd uhppoted-app-sheets && $(ARM)     go build -trimpath -o ../dist/arm/$(DIST)     ./...
	cd uhppoted-app-sheets && $(ARM7)    go build -trimpath -o ../dist/arm7/$(DIST)    ./...

	cd uhppoted-app-wild-apricot && $(LINUX)   go build -trimpath -o ../dist/linux/$(DIST)   ./...
	cd uhppoted-app-wild-apricot && $(DARWIN)  go build -trimpath -o ../dist/darwin/$(DIST)  ./...
	cd uhppoted-app-wild-apricot && $(WINDOWS) go build -trimpath -o ../dist/windows/$(DIST) ./...	
	cd uhppoted-app-wild-apricot && $(ARM)     go build -trimpath -o ../dist/arm/$(DIST)     ./...
	cd uhppoted-app-wild-apricot && $(ARM7)    go build -trimpath -o ../dist/arm7/$(DIST)    ./...

	cd uhppoted-app-db && $(LINUX)   go build -trimpath -o ../dist/linux/$(DIST)   ./...
	cd uhppoted-app-db && $(DARWIN)  go build -trimpath -o ../dist/darwin/$(DIST)  ./...
	cd uhppoted-app-db && $(WINDOWS) go build -trimpath -o ../dist/windows/$(DIST) ./...	
	cd uhppoted-app-db && $(ARM)     go build -trimpath -o ../dist/arm/$(DIST)     ./...
	cd uhppoted-app-db && $(ARM7)    go build -trimpath -o ../dist/arm7/$(DIST)    ./...

	cd uhppote-simulator && $(LINUX)   go build -trimpath -o ../dist/linux/$(DIST)   ./...
	cd uhppote-simulator && $(DARWIN)  go build -trimpath -o ../dist/darwin/$(DIST)  ./...
	cd uhppote-simulator && $(WINDOWS) go build -trimpath -o ../dist/windows/$(DIST) ./...
	cd uhppote-simulator && $(ARM)     go build -trimpath -o ../dist/arm/$(DIST)     ./...
	cd uhppote-simulator && $(ARM7)    go build -trimpath -o ../dist/arm7/$(DIST)    ./...

	cp uhppoted-rest/documentation/uhppoted-api.yaml      documentation/openapi/
	cp uhppote-simulator/documentation/simulator-api.yaml documentation/openapi/
	cp uhppoted-rest/documentation/uhppoted-api.yaml      install/openapi/
	cp uhppote-simulator/documentation/simulator-api.yaml install/openapi/
	cp -r install/openapi/* dist/openapi/$(DIST)/

	cp uhppoted-mqtt/documentation/TLS.md        cookbook/mqtt/
	cp uhppoted-mqtt/documentation/signatures.md cookbook/mqtt/
	cp uhppoted-app-s3/documentation/signing.md  cookbook/s3/

release: build-all docker
	find . -name ".DS_Store" -delete

	tar --directory=uhppoted-codegen --exclude=".DS_Store" -cvzf dist/linux/$(DIST)/uhppoted-codegen-bindings.tar.gz   bindings
	tar --directory=uhppoted-codegen --exclude=".DS_Store" -cvzf dist/darwin/$(DIST)/uhppoted-codegen-bindings.tar.gz  bindings
	tar --directory=uhppoted-codegen --exclude=".DS_Store" -cvzf dist/windows/$(DIST)/uhppoted-codegen-bindings.tar.gz bindings
	tar --directory=uhppoted-codegen --exclude=".DS_Store" -cvzf dist/arm/$(DIST)/uhppoted-codegen-bindings.tar.gz     bindings
	tar --directory=uhppoted-codegen --exclude=".DS_Store" -cvzf dist/arm7/$(DIST)/uhppoted-codegen-bindings.tar.gz    bindings
	tar --directory=uhppoted-codegen/generated --exclude=".DS_Store" --exclude="go/bin"               -cvzf dist/linux/$(DIST)/uhppoted-codegen-go.tar.gz       go
	tar --directory=uhppoted-codegen/generated --exclude=".DS_Store" --exclude="rust/uhppoted/target" -cvzf dist/linux/$(DIST)/uhppoted-codegen-rust.tar.gz     rust
	tar --directory=uhppoted-codegen/generated --exclude=".DS_Store" --exclude="python/__pycache__"   -cvzf dist/linux/$(DIST)/uhppoted-codegen-python.tar.gz   python
	tar --directory=uhppoted-codegen/generated --exclude=".DS_Store" --exclude="go/bin"               -cvzf dist/darwin/$(DIST)/uhppoted-codegen-go.tar.gz      go
	tar --directory=uhppoted-codegen/generated --exclude=".DS_Store" --exclude="rust/uhppoted/target" -cvzf dist/darwin/$(DIST)/uhppoted-codegen-rust.tar.gz    rust
	tar --directory=uhppoted-codegen/generated --exclude=".DS_Store" --exclude="python/__pycache__"   -cvzf dist/darwin/$(DIST)/uhppoted-codegen-python.tar.gz  python
	tar --directory=uhppoted-codegen/generated --exclude=".DS_Store" --exclude="go/bin"               -cvzf dist/windows/$(DIST)/uhppoted-codegen-go.tar.gz     go
	tar --directory=uhppoted-codegen/generated --exclude=".DS_Store" --exclude="rust/uhppoted/target" -cvzf dist/windows/$(DIST)/uhppoted-codegen-rust.tar.gz   rust
	tar --directory=uhppoted-codegen/generated --exclude=".DS_Store" --exclude="python/__pycache__"   -cvzf dist/windows/$(DIST)/uhppoted-codegen-python.tar.gz python
	tar --directory=uhppoted-codegen/generated --exclude=".DS_Store" --exclude="go/bin"               -cvzf dist/arm/$(DIST)/uhppoted-codegen-go.tar.gz         go
	tar --directory=uhppoted-codegen/generated --exclude=".DS_Store" --exclude="rust/uhppoted/target" -cvzf dist/arm/$(DIST)/uhppoted-codegen-rust.tar.gz       rust
	tar --directory=uhppoted-codegen/generated --exclude=".DS_Store" --exclude="python/__pycache__"   -cvzf dist/arm/$(DIST)/uhppoted-codegen-python.tar.gz     python
	tar --directory=uhppoted-codegen/generated --exclude=".DS_Store" --exclude="go/bin"               -cvzf dist/arm7/$(DIST)/uhppoted-codegen-go.tar.gz        go
	tar --directory=uhppoted-codegen/generated --exclude=".DS_Store" --exclude="rust/uhppoted/target" -cvzf dist/arm7/$(DIST)/uhppoted-codegen-rust.tar.gz      rust
	tar --directory=uhppoted-codegen/generated --exclude=".DS_Store" --exclude="python/__pycache__"   -cvzf dist/arm7/$(DIST)/uhppoted-codegen-python.tar.gz    python

	tar --directory=dist/linux  --exclude=".DS_Store" -cvzf dist/$(DIST)-linux.tar.gz $(DIST)
	tar --directory=dist/darwin --exclude=".DS_Store" -cvzf dist/$(DIST)-darwin.tar.gz $(DIST)
	tar --directory=dist/arm    --exclude=".DS_Store" -cvzf dist/$(DIST)-arm.tar.gz $(DIST)
	tar --directory=dist/arm7   --exclude=".DS_Store" -cvzf dist/$(DIST)-arm7.tar.gz $(DIST)
	cd dist/windows && zip --recurse-paths ../$(DIST)-windows.zip $(DIST)

release-v0.8.6: 
	yapf -ri ./internal
	python ./internal/release.py --version=v0.8.6 --node-red=1.1.5 --release

bump: 
	yapf -ri ./internal
	python ./internal/release.py --version=v0.8.6 --node-red=1.1.5 --bump

publish: release
	echo "Releasing version $(VERSION)"
	rm -f dist/development-arm.tar.gz
	rm -f dist/development-arm7.tar.gz
	rm -f dist/development-darwin.tar.gz
	rm -f dist/development-linux.tar.gz
	rm -f dist/development-windows.zip
	gh release create "$(VERSION)" ./dist/*.tar.gz ./dist/*.zip --draft --prerelease --title "$(VERSION)-beta" --notes-file release-notes.md

build-github: 
	cd uhppote-core              && go build -trimpath ./...
	cd uhppoted-lib              && go build -trimpath ./...
	cd uhppote-simulator         && go build -trimpath ./...
	cd uhppote-cli               && go build -trimpath ./...
	cd uhppoted-rest             && go build -trimpath ./...
	cd uhppoted-mqtt             && go build -trimpath ./...
	cd uhppoted-httpd            && go build -trimpath ./...
	cd uhppoted-tunnel           && go build -trimpath ./...
	cd uhppoted-codegen          && go build -trimpath ./...
#	make -C ./uhppoted-dll -f Makefile build-all
	cd uhppoted-app-s3           && go build -trimpath ./...
	cd uhppoted-app-sheets       && go build -trimpath ./...
	cd uhppoted-app-wild-apricot && go build -trimpath ./...
	cd uhppoted-app-db           && go build -trimpath ./...

debug: 
	cd uhppote-simulator && $(WINDOWS) go build -trimpath -o ../dist/windows/$(DIST) ./...

simulator: 
	./bin/uhppote-simulator --debug --bind 0.0.0.0:60000 --rest 0.0.0.0:8000 --devices "./runtime/simulation/devices"

uhppoted-rest:
	./bin/uhppoted-rest --console

uhppoted-mqtt: 
	./bin/uhppoted-mqtt --console

swagger: 
	docker run --detach --publish 80:8080 --name swagger --rm swaggerapi/swagger-editor 
	sleep 1
	open http://127.0.0.1:80

docker:
	cd uhppote-simulator && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -o ../docker/simulator     ./...
	cd uhppote-simulator && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -o ../docker/uhppoted-rest ./...
	cd uhppote-simulator && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -o ../docker/integration-tests/simulator ./...
	cd uhppoted-rest     && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -o ../docker/uhppoted-rest ./...
	cd uhppoted-mqtt     && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -o ../docker/uhppoted-mqtt ./...
	
	docker image     prune -f
	docker container prune -f
	cd ./docker/simulator     && docker build -f Dockerfile -t uhppoted/simulator . 
	cd ./docker/uhppoted-rest && docker build -f Dockerfile -t uhppoted/rest      . 
	cd ./docker/uhppoted-mqtt && docker build -f Dockerfile -t uhppoted/mqtt      . 
	cd ./docker/hivemq        && docker build -f Dockerfile -t hivemq/uhppoted    . 
	cd ./docker/integration-tests/simulator && docker build -f Dockerfile -t integration-tests/simulator . 
	cd ./docker/integration-tests/mqttd     && docker build -f Dockerfile -t integration-tests/mqttd     . 
	cd ./docker/integration-tests/hivemq    && docker build -f Dockerfile -t integration-tests/hivemq    . 

docker-simulator:
	docker run --detach --publish 8000:8000 --publish 60000:60000/udp --name simulator --rm uhppoted/simulator
	sleep 1
	./bin/uhppote-cli --debug set-listener 405419896 192.168.1.100:60001
	./bin/uhppote-cli --debug set-listener 303986753 192.168.1.100:60001
	./bin/uhppote-cli --debug set-listener 201020304 192.168.1.100:60001

docker-simulator-tunnel:
	docker run --detach --publish 8000:8000 --publish 60005:60000/udp --name simulator --rm uhppoted/simulator

docker-hivemq:
	docker run --detach --publish 8081:8080 --publish 1883:1883 --publish 8883:8883 --name hivemq --rm hivemq/uhppoted

docker-rest:
	docker run --detach --publish 8080:8080 --name restd --rm uhppoted/rest

docker-mqtt:
	docker run --detach --name mqttd --rm uhppoted/mqtt

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

docker-integration-tests:
	docker run --detach --publish 8000:8000 --publish 60000:60000/udp --name qwerty --rm integration-tests/simulator

hivemq-listen:
	# mqtt subscribe --topic 'uhppoted/reply/#' | jq '.' 
	mqtt subscribe --topic 'uhppoted/#' | jq '.' 
#	open runtime/mqtt-spy-0.5.4-jar-with-dependencies.jar

# fly.io:
# 	cd uhppote-simulator && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -o ../docker/fly.io ./...
# 	cd uhppoted-httpd    && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -o ../docker/fly.io ./...
# 	cp -R uhppoted-httpd/html         docker/fly.io/httpd
# 	cp -R uhppoted-httpd/translations docker/fly.io/httpd
# 	cd ./docker/fly.io   && docker build -f Dockerfile -t uhppoted/fly.io . 
# 	# docker run --publish 8080:8080 --name fly.io --rm uhppoted/fly.io
# 	docker run --detach --publish 8080:8080 --name fly.io --rm uhppoted/fly.io



