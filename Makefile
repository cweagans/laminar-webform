SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules
PROJECT_NAME = laminar-webform

.PHONY: bin
bin: bin/windows/$(PROJECT_NAME).exe bin/linux/$(PROJECT_NAME) bin/darwin/$(PROJECT_NAME)

bin/windows/$(PROJECT_NAME).exe: vendor main.go config.go
	env GOOS=windows GOARCH=amd64 go build -o bin/windows/$(PROJECT_NAME).exe

bin/linux/$(PROJECT_NAME): vendor main.go config.go
	env GOOS=linux GOARCH=amd64 go build -o bin/linux/$(PROJECT_NAME)

bin/darwin/$(PROJECT_NAME): vendor main.go config.go
	env GOOS=darwin GOARCH=amd64 go build -o bin/darwin/$(PROJECT_NAME)

.PHONY: release
release: bin release/$(PROJECT_NAME).win64.tar.gz release/$(PROJECT_NAME).linux64.tar.gz release/$(PROJECT_NAME).darwin64.tar.gz release/sha256sums.txt

release/sha256sums.txt: release/$(PROJECT_NAME).win64.tar.gz release/$(PROJECT_NAME).linux64.tar.gz release/$(PROJECT_NAME).darwin64.tar.gz
	shasum --algorithm 256 release/* > release/sha256sums.txt

release/$(PROJECT_NAME).win64.tar.gz: bin/windows/$(PROJECT_NAME).exe
	[ -d release ] || mkdir release
	tar czf release/$(PROJECT_NAME).win64.tar.gz -C bin/windows $(PROJECT_NAME).exe

release/$(PROJECT_NAME).linux64.tar.gz: bin/linux/$(PROJECT_NAME)
	[ -d release ] || mkdir release
	tar czf release/$(PROJECT_NAME).linux64.tar.gz -C bin/linux $(PROJECT_NAME)

release/$(PROJECT_NAME).darwin64.tar.gz: bin/darwin/$(PROJECT_NAME)
	[ -d release ] || mkdir release
	tar czf release/$(PROJECT_NAME).darwin64.tar.gz -C bin/darwin $(PROJECT_NAME)

vendor: go.sum
	go mod vendor

run: bin
	docker-compose run data-mover

go.sum:

docker-compose.yml:

clean:
	([ -d bin ] && rm -r bin) || true
	([ -d release ] && rm -r release) || true
