all: lint test

.PHONY: lint
lint:
	go vet -v ./...

.PHONY: test
test:
	go test -v ./...

.PHONY: gendoc
gendoc:
	gomarkdoc --output README.md --footer "![test workflow](https://github.com/chg1f/errorx/actions/workflows/test.yml/badge.svg?branch=master)" ./...
