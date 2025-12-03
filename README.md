# simplemq-cli

A CLI tool for [SAKURA Cloud SimpleMQ](https://manual.sakura.ad.jp/cloud/appliance/simplemq/index.html) service.

## Installation

```bash
go install github.com/fujiwara/simplemq-cli/cmd/simplemq-cli@latest
```

Or download from [Releases](https://github.com/fujiwara/simplemq-cli/releases).

## Usage

```
Usage: simplemq-cli <command> [flags]

Flags:
  -h, --help       Show context-sensitive help.
  -v, --version    Show version and exit.
      --debug      Enable debug mode ($SIMPLEMQ_DEBUG).

Commands:
  message send --queue-name=STRING --api-key=STRING <content> [flags]
    Send message to queue

  message receive --queue-name=STRING --api-key=STRING [flags]
    Receive message from queue

  message delete --queue-name=STRING --api-key=STRING <message-id> [flags]
    Delete message from queue
```

### Send a message

```bash
simplemq-cli message send --queue-name myqueue --api-key <api-key> "Hello, World!"
```

### Receive messages

```bash
# Receive a single message
# Do not delete the message after receiving
simplemq-cli message receive --queue-name myqueue --api-key <api-key>

# Receive with polling (wait for messages)
simplemq-cli message receive --queue-name myqueue --api-key <api-key> --polling --count 10

# Auto-delete messages after receiving
simplemq-cli message receive --queue-name myqueue --api-key <api-key> --auto-delete
```

### Delete a message

```bash
simplemq-cli message delete --queue-name myqueue --api-key <api-key> <message-id>
```

**Note:** SimpleMQ API only accepts message content matching `^[0-9a-zA-Z+/=]*$`. By default, this CLI automatically encodes/decodes message content using Base64. Use `--raw` flag to disable this behavior.

Received messages are output as JSON to stdout.

For example:
```json
{
  "id": "019adef1-91d4-7aac-ad68-2e1a5c881e95",
  "content": "こんにちは世界",
  "created_at": "2025-12-02T21:02:44.82+09:00",
  "updated_at": "2025-12-02T21:02:47.11+09:00",
  "expires_at": "2025-12-06T21:02:44.82+09:00",
  "acquired_at": "2025-12-02T21:02:47.11+09:00",
  "visibility_timeout_at": "2025-12-02T21:03:17.11+09:00"
}
```

If `--raw` flag is specified, the raw message from SimpleMQ API is output without Base64 decoding and time conversion.
```json
{"id":"019adf15-f115-7efd-942c-423a6b6a2250","content":"44GT44KT44Gr44Gh44Gv5LiW55WM","created_at":1764679348501,"updated_at":1764679355018,"expires_at":1765024948501,"acquired_at":1764679355018,"visibility_timeout_at":1764679385018}
```

## Options

### Message Options

| Option | Environment Variable | Description |
|--------|---------------------|-------------|
| `--queue-name` | `SIMPLEMQ_QUEUE_NAME` | Queue name (required) |
| `--api-key` | `SIMPLEMQ_API_KEY` | API Key (required) |
| `--raw` | `SIMPLEMQ_RAW` | Handle raw message without Base64 encoding/decoding (default: false) |

### Send Options

No additional options for sending messages.

### Receive Options

| Option | Environment Variable | Default | Description |
|--------|---------------------|---------|-------------|
| `--polling` | `SIMPLEMQ_POLLING` | false | Enable polling to receive messages |
| `--count` | `SIMPLEMQ_RECEIVE_COUNT` | 1 | Number of messages to receive |
| `--auto-delete` | `SIMPLEMQ_AUTO_DELETE` | false | Automatically delete messages after receiving |
| `--interval` | `SIMPLEMQ_POLLING_INTERVAL` | 1s | Polling interval |

### Delete Options

No additional options for deleting messages.

## License

MIT

## Author

fujiwara
