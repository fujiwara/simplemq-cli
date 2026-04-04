# Changelog

## [v0.9.2](https://github.com/fujiwara/simplemq-cli/compare/v0.9.1...v0.9.2) - 2026-04-04
- Recommend SAKURA_ENDPOINTS_SIMPLE_MQ_MESSAGE over SIMPLEMQ_MESSAGE_API_URL by @fujiwara in https://github.com/fujiwara/simplemq-cli/pull/66

## [v0.9.1](https://github.com/fujiwara/simplemq-cli/compare/v0.9.0...v0.9.1) - 2026-04-03
- Add --timeout option and localserver --latency option by @fujiwara in https://github.com/fujiwara/simplemq-cli/pull/64

## [v0.9.0](https://github.com/fujiwara/simplemq-cli/compare/v0.8.0...v0.9.0) - 2026-04-03
- Add --file option to message send command by @fujiwara in https://github.com/fujiwara/simplemq-cli/pull/62
- Bump github.com/sacloud/saclient-go from 0.3.1 to 0.3.5 by @dependabot[bot] in https://github.com/fujiwara/simplemq-cli/pull/61
- Bump golang.org/x/sys from 0.41.0 to 0.42.0 by @dependabot[bot] in https://github.com/fujiwara/simplemq-cli/pull/59
- Bump actions/setup-go from 6.3.0 to 6.4.0 by @dependabot[bot] in https://github.com/fujiwara/simplemq-cli/pull/58
- Bump modernc.org/sqlite from 1.46.1 to 1.48.0 by @dependabot[bot] in https://github.com/fujiwara/simplemq-cli/pull/60

## [v0.8.0](https://github.com/fujiwara/simplemq-cli/compare/v0.7.1...v0.8.0) - 2026-03-12
- Add --stdin, --each-line, --each-json options to message send by @fujiwara in https://github.com/fujiwara/simplemq-cli/pull/56

## [v0.7.1](https://github.com/fujiwara/simplemq-cli/compare/v0.7.0...v0.7.1) - 2026-03-09
- Update simplemq-api-go to v0.5.0 by @fujiwara in https://github.com/fujiwara/simplemq-cli/pull/54

## [v0.7.0](https://github.com/fujiwara/simplemq-cli/compare/v0.6.1...v0.7.0) - 2026-03-07
- Switch SQLite driver to modernc.org/sqlite (pure Go, no CGO) by @fujiwara in https://github.com/fujiwara/simplemq-cli/pull/52

## [v0.6.1](https://github.com/fujiwara/simplemq-cli/compare/v0.6.0...v0.6.1) - 2026-03-07
- Address review comments from PR #47 by @fujiwara in https://github.com/fujiwara/simplemq-cli/pull/48
- Improve README structure and add coding guidelines to CLAUDE.md by @fujiwara in https://github.com/fujiwara/simplemq-cli/pull/50

## [v0.6.1](https://github.com/fujiwara/simplemq-cli/compare/v0.6.0...v0.6.1) - 2026-03-07
- Address review comments from PR #47 by @fujiwara in https://github.com/fujiwara/simplemq-cli/pull/48
- Improve README structure and add coding guidelines to CLAUDE.md by @fujiwara in https://github.com/fujiwara/simplemq-cli/pull/50

## [v0.6.0](https://github.com/fujiwara/simplemq-cli/compare/v0.5.0...v0.6.0) - 2026-03-06
- Add API request rate limit documentation by @fujiwara in https://github.com/fujiwara/simplemq-cli/pull/44
- Add --queue-id option to queue commands by @fujiwara in https://github.com/fujiwara/simplemq-cli/pull/46
- Add SQLite storage backend for localserver by @fujiwara in https://github.com/fujiwara/simplemq-cli/pull/47

## [v0.5.0](https://github.com/fujiwara/simplemq-cli/compare/v0.4.0...v0.5.0) - 2026-03-05
- Add request logging to localserver by @fujiwara in https://github.com/fujiwara/simplemq-cli/pull/41
- Add simplemq-localserver to goreleaser and Makefile by @fujiwara in https://github.com/fujiwara/simplemq-cli/pull/43

## [v0.4.0](https://github.com/fujiwara/simplemq-cli/compare/v0.3.0...v0.4.0) - 2026-03-04
- Add Config struct to localserver, migrate from flag to kong by @fujiwara in https://github.com/fujiwara/simplemq-cli/pull/38
- Move Local Server section after Options in README by @fujiwara in https://github.com/fujiwara/simplemq-cli/pull/40

## [v0.3.0](https://github.com/fujiwara/simplemq-cli/compare/v0.2.2...v0.3.0) - 2026-03-03
- Add optional API key validation to localserver by @fujiwara in https://github.com/fujiwara/simplemq-cli/pull/35
- Bump golang.org/x/sys from 0.40.0 to 0.41.0 by @dependabot[bot] in https://github.com/fujiwara/simplemq-cli/pull/34
- Bump Songmu/tagpr from 1.15.0 to 1.17.1 by @dependabot[bot] in https://github.com/fujiwara/simplemq-cli/pull/33
- Bump github.com/sacloud/saclient-go from 0.2.6 to 0.3.1 by @dependabot[bot] in https://github.com/fujiwara/simplemq-cli/pull/32
- Bump goreleaser/goreleaser-action from 6.4.0 to 7.0.0 by @dependabot[bot] in https://github.com/fujiwara/simplemq-cli/pull/30
- Bump actions/setup-go from 6.2.0 to 6.3.0 by @dependabot[bot] in https://github.com/fujiwara/simplemq-cli/pull/29
- Bump github.com/alecthomas/kong from 1.13.0 to 1.14.0 by @dependabot[bot] in https://github.com/fujiwara/simplemq-cli/pull/31
- Support Go 1.25 and 1.26 in CI, build release with 1.26 by @fujiwara in https://github.com/fujiwara/simplemq-cli/pull/37

## [v0.2.2](https://github.com/fujiwara/simplemq-cli/compare/v0.2.1...v0.2.2) - 2026-02-27
- Fix data race in localserver queue operations by @fujiwara in https://github.com/fujiwara/simplemq-cli/pull/27

## [v0.2.1](https://github.com/fujiwara/simplemq-cli/compare/v0.2.0...v0.2.1) - 2026-02-08
- Add localserver and message integration tests by @fujiwara in https://github.com/fujiwara/simplemq-cli/pull/24
- Replace map[string]any with typed structs for localserver responses by @fujiwara in https://github.com/fujiwara/simplemq-cli/pull/26
- Bump golang.org/x/sys from 0.38.0 to 0.40.0 by @dependabot[bot] in https://github.com/fujiwara/simplemq-cli/pull/22

## [v0.2.0](https://github.com/fujiwara/simplemq-cli/compare/v0.1.0...v0.2.0) - 2026-02-06
- update SDKs by @fujiwara in https://github.com/fujiwara/simplemq-cli/pull/17
- Update simplemq-api-go to v0.4.0 by @fujiwara in https://github.com/fujiwara/simplemq-cli/pull/23
- Bump actions/checkout from 5.0.1 to 6.0.2 by @dependabot[bot] in https://github.com/fujiwara/simplemq-cli/pull/20
- Bump Songmu/tagpr from 1.8.4 to 1.15.0 by @dependabot[bot] in https://github.com/fujiwara/simplemq-cli/pull/19
- Bump actions/setup-go from 6.1.0 to 6.2.0 by @dependabot[bot] in https://github.com/fujiwara/simplemq-cli/pull/18

## [v0.1.0](https://github.com/fujiwara/simplemq-cli/compare/v0.0.4...v0.1.0) - 2025-12-03
- send - accepts from stdin by @fujiwara in https://github.com/fujiwara/simplemq-cli/pull/9
- Add queue management commands by @fujiwara in https://github.com/fujiwara/simplemq-cli/pull/11
- add confirmation for queue rotate-api-key command by @fujiwara in https://github.com/fujiwara/simplemq-cli/pull/12

## [v0.0.4](https://github.com/fujiwara/simplemq-cli/compare/v0.0.3...v0.0.4) - 2025-12-03
- fix receive --count by @fujiwara in https://github.com/fujiwara/simplemq-cli/pull/7

## [v0.0.3](https://github.com/fujiwara/simplemq-cli/compare/v0.0.2...v0.0.3) - 2025-12-03
- Add delete command and replace --base64 with --raw option by @fujiwara in https://github.com/fujiwara/simplemq-cli/pull/6

## [v0.0.2](https://github.com/fujiwara/simplemq-cli/compare/v0.0.1...v0.0.2) - 2025-12-02

## [v0.0.2](https://github.com/fujiwara/simplemq-cli/compare/v0.0.1...v0.0.2) - 2025-12-02

## [v0.0.1](https://github.com/fujiwara/simplemq-cli/commits/v0.0.1) - 2025-12-02
