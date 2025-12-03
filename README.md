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
  message send --queue=STRING --api-key=STRING <content> [flags]
    Send message to queue

  message receive --queue=STRING --api-key=STRING [flags]
    Receive message from queue

  message delete --queue=STRING --api-key=STRING <message-id> [flags]
    Delete message from queue

  queue create --queue=STRING [flags]
    Create a new queue

  queue list [flags]
    List queues

  queue get --queue=STRING [flags]
    Get queue details

  queue modify --queue=STRING [flags]
    Modify queue settings

  queue delete --queue=STRING [flags]
    Delete a queue

  queue message-count --queue=STRING [flags]
    Get message count in a queue

  queue rotate-api-key --queue=STRING [flags]
    Rotate API key for a queue

  queue purge --queue=STRING [flags]
    Purge all messages in a queue
```

### Send a message

```bash
simplemq-cli message send --queue myqueue --api-key <api-key> "Hello, World!"
```

### Receive messages

```bash
# Receive a single message
# Do not delete the message after receiving
simplemq-cli message receive --queue myqueue --api-key <api-key>

# Receive with polling (wait for messages)
simplemq-cli message receive --queue myqueue --api-key <api-key> --polling --count 10

# Auto-delete messages after receiving
simplemq-cli message receive --queue myqueue --api-key <api-key> --auto-delete
```

### Delete a message

```bash
simplemq-cli message delete --queue myqueue --api-key <api-key> <message-id>
```

### Queue Management

Queue management commands require `SAKURACLOUD_ACCESS_TOKEN` and `SAKURACLOUD_ACCESS_TOKEN_SECRET` environment variables for authentication.

#### Create a queue

```bash
simplemq-cli queue create --queue myqueue --description "My queue description"
```

#### List queues

```bash
simplemq-cli queue list
```

#### Get queue details

```bash
simplemq-cli queue get --queue myqueue
```

#### Modify queue settings

```bash
# Modify visibility timeout and message expiration time
simplemq-cli queue modify --queue myqueue --visibility-timeout-seconds 60 --expire-seconds 86400
```

#### Delete a queue

```bash
# Delete with confirmation prompt
simplemq-cli queue delete --queue myqueue

# Delete without confirmation
simplemq-cli queue delete --queue myqueue -f
```

#### Get message count

```bash
simplemq-cli queue message-count --queue myqueue
```

#### Rotate API key

```bash
simplemq-cli queue rotate-api-key --queue myqueue
```

#### Purge all messages

```bash
# Purge with confirmation prompt
simplemq-cli queue purge --queue myqueue

# Purge without confirmation
simplemq-cli queue purge --queue myqueue -f
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
| `--queue` | `SIMPLEMQ_QUEUE_NAME` | Queue name (required) |
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

### Queue Options

Queue commands require the following environment variables for authentication:

| Environment Variable | Description |
|---------------------|-------------|
| `SAKURACLOUD_ACCESS_TOKEN` | Access token for queue management (required) |
| `SAKURACLOUD_ACCESS_TOKEN_SECRET` | Access token secret for queue management (required) |

| Option | Environment Variable | Description |
|--------|---------------------|-------------|
| `--queue`, `-q` | `SIMPLEMQ_QUEUE_NAME` | Queue name (required for most commands) |

### Queue Create Options

| Option | Description |
|--------|-------------|
| `--description` | Description of the queue |

### Queue Modify Options

| Option | Description |
|--------|-------------|
| `--visibility-timeout-seconds` | Visibility timeout in seconds |
| `--expire-seconds` | Message expiration time in seconds |

### Queue Delete / Purge Options

| Option | Environment Variable | Default | Description |
|--------|---------------------|---------|-------------|
| `-f`, `--force` | `SIMPLEMQ_FORCE` | false | Force operation without confirmation prompt |

## License

MIT

## Author

fujiwara
