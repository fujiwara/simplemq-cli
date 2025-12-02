# simplemq-cli

A CLI tool for [SAKURA Cloud SimpleMQ](https://manual.sakura.ad.jp/cloud/manual-simplemq.html) service.

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
```

### Send a message

```bash
simplemq-cli message --queue-name myqueue --api-key <api-key> --base64 send "Hello, World!"
```

### Receive messages

```bash
# Receive a single message
simplemq-cli message --queue-name myqueue --api-key <api-key> --base64 receive

# Receive with polling (wait for messages)
simplemq-cli message --queue-name myqueue --api-key <api-key> --base64 receive --polling --count 10

# Auto-delete messages after receiving
simplemq-cli message --queue-name myqueue --api-key <api-key> --base64 receive --auto-delete
```

**Note:** SimpleMQ API only accepts message content matching `^[0-9a-zA-Z+/=]*$`. Use `--base64` flag to automatically encode/decode message content.

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

## Options

### Global Options

| Option | Environment Variable | Description |
|--------|---------------------|-------------|
| `--queue-name` | `SIMPLEMQ_QUEUE_NAME` | Queue name (required) |

### Message Options

| Option | Environment Variable | Description |
|--------|---------------------|-------------|
| `--api-key` | `SIMPLEMQ_API_KEY` | API Key (required) |
| `--base64` | `SIMPLEMQ_BASE64` | Use Base64 encoding for message content (default: false) |

### Receive Options

| Option | Environment Variable | Default | Description |
|--------|---------------------|---------|-------------|
| `--polling` | `SIMPLEMQ_POLLING` | false | Enable polling to receive messages |
| `--count` | `SIMPLEMQ_RECEIVE_COUNT` | 1 | Number of messages to receive |
| `--auto-delete` | `SIMPLEMQ_AUTO_DELETE` | false | Automatically delete messages after receiving |
| `--interval` | `SIMPLEMQ_POLLING_INTERVAL` | 1s | Polling interval |

## License

MIT

## Author

fujiwara
