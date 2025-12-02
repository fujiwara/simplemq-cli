# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

simplemq-cli is a CLI tool for interacting with SAKURA Cloud SimpleMQ service. It provides commands to send and receive messages from SimpleMQ queues.

## Build Commands

```bash
# Build the binary
make simplemq-cli

# Run tests
make test

# Install to $GOPATH/bin
make install

# Build distribution binaries
make dist

# Clean build artifacts
make clean
```

## Architecture

The project uses a standard Go CLI structure with [kong](https://github.com/alecthomas/kong) for argument parsing:

- `cmd/simplemq-cli/main.go` - Entry point with signal handling (platform-specific via build tags)
- `cli.go` - CLI struct definitions with kong tags for flags and environment variables
- `main.go` - Command routing via `Run()` function
- `send.go` - Send message command implementation
- `receive.go` - Receive message command implementation with polling support
- `version.go` - Version string (set during build)

The CLI interacts with SimpleMQ via `github.com/sacloud/simplemq-api-go` client library.

## CLI Usage

Commands follow the pattern: `simplemq-cli --queue-name <name> message <subcommand>`

- `message send <content>` - Send a message to the queue
- `message receive` - Receive messages from the queue (supports polling, auto-delete)

Configuration can be via flags or environment variables (SIMPLEMQ_QUEUE_NAME, SIMPLEMQ_API_KEY, etc.).
