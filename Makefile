.PHONY: clean test

simplemq-cli: go.* *.go cmd/simplemq-cli/*.go
	go build -o $@ ./cmd/simplemq-cli

simplemq-localserver: go.* *.go cmd/simplemq-localserver/*.go localserver/*.go
	go build -tags sqlite -o $@ ./cmd/simplemq-localserver

clean:
	rm -rf simplemq-cli simplemq-localserver dist/

test:
	go test -v ./...

install:
	go install github.com/fujiwara/simplemq-cli/cmd/simplemq-cli
	go install github.com/fujiwara/simplemq-cli/cmd/simplemq-localserver

dist:
	goreleaser build --snapshot --clean
