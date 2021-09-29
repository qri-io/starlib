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
	@if [ ! -d ../website/content/docs ]; then \
		echo "website repo not found, should be a sibling to starlib/"; \
		exit 1; \
	fi
	@mkdir -p ../website/content/docs/reference/starlark-packages
	for sourcefile in $$(find . | grep doc.go) ; do \
		targetfile="`echo $${sourcefile} | sed 's/\/doc.go/.md/' | sed 's/encoding\///' `"; \
		outline template --template asset/doc_template.txt $${sourcefile} > ../website/content/docs/reference/starlark-packages/$${targetfile} ; \
	done

list-deps:
	go list -f '{{.Deps}}' | tr "[" " " | tr "]" " " | xargs go list -f '{{if not .Standard}}{{.ImportPath}}{{end}}'

run-circleci-tests-locally:
	circleci local execute .
