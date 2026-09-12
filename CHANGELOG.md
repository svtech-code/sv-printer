# Changelog

All notable changes follow [Conventional Commits](https://www.conventionalcommits.org/).
Releases are tagged `vX.Y.Z`.

## [Unreleased]

## [0.1.0] - 2026-09-12

### Added

- **Initial release of SV Print**: cross-platform local agent (Go) that lets web applications print to thermal printers (USB, serial, TCP) via a local HTTP API.
- **ESC/POS protocol**: full command set for thermal receipt printers (`pkg/escpos/`).
- **Structured receipt API**: `POST /api/v1/print/receipt` accepts JSON with styled lines and generates ESC/POS.
- **Raw ESC/POS API**: `POST /api/v1/print` for direct byte payloads (pro feature).
- **WebSocket events**: `GET /api/v1/events` streams `print.*`, `agent.status`, and `printer.*` events in real time (pro feature).
- **Browser WebSocket auth**: accept `?token=<token>` query parameter for `GET /api/v1/events` since native WebSocket API cannot set headers.
- **Freemium licensing**: trial (watermark + 50/day), beta (signed with expiry), full; Ed25519-signed offline verification.
- **Device binding**: licenses can be bound to a specific machine via fingerprint.
- **`sv-license` CLI tool**: `keygen` and `sign` commands for issuing licenses.
- **Multi-platform CI/CD**: GitHub Actions matrix (Linux, macOS, Windows) with `gofmt`, `go vet`, `go test`.
- **Cross-platform releases**: goreleaser v2 config for binary distribution (Linux/macOS/Windows, amd64/arm64).
- **Integration guide**: full API reference, receipt schema, WebSocket events, error codes, and client examples (JavaScript, Python, PHP).
- **E2E smoke test**: end-to-end test with simulated TCP printer.

### Fixed

- **WebSocket race condition**: `EventsHandler` now subscribes to the bus before accepting the WebSocket handshake, preventing lost events on CI and under load.
- **Windows CI CRLF**: forced LF via `.gitattributes` to fix `gofmt` failures on Windows checkout.
- **CI actions bumped to v7**: `actions/checkout@v7` and `actions/setup-go@v7` (node24, eliminates deprecation warnings).

### Docs

- **README redesigned**: hero with badges, Mermaid architecture diagram, key features table, bilingual (EN/ES).
- **Spanish mirrors**: `README_ES.md`, `docs/integration_ES.md`, `docs/licensing_ES.md`.
- **Repository standards**: `CONTRIBUTING`, `SECURITY`, `CODE_OF_CONDUCT` (EN/ES).
- **License**: Business Source License 1.1 (converts to Apache 2.0 on 2030-09-11).

### Changed

- **Product rename**: SV Print Agent → SV Print (identifier: `sv-print`).
- **Config**: corrected `allowed-origins` example to `svtech.cl`.
