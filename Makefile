.PHONY: clean test

simplemq-cli: go.* *.go
	go build -o $@ ./cmd/simplemq-cli

clean:
	rm -rf simplemq-cli dist/

test:
	go test -v ./...

install:
	go install github.com/fujiwara/simplemq-cli/cmd/simplemq-cli

dist:
	goreleaser build --snapshot --clean
