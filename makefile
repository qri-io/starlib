GOFILES = $(shell find . -name '*.go' -not -path './vendor/*')

default: test

lint:
	golint ./...

test:
	go test ./... -v --coverprofile=coverage.txt --covermode=atomic

test-all-coverage:
	./.circleci/cover.test.sh

update-changelog:
	conventional-changelog -p angular -i CHANGELOG.md -s

docs:
	outline template $(find . | grep doc.go | tr '\n' ' ') > ../website/content/docs/transforms/starlib.md

list-deps:
	go list -f '{{.Deps}}' | tr "[" " " | tr "]" " " | xargs go list -f '{{if not .Standard}}{{.ImportPath}}{{end}}'

run-circleci-tests-locally:
	circleci local execute .
