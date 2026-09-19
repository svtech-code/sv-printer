# Changelog

All notable changes follow [Conventional Commits](https://www.conventionalcommits.org/).
Releases are tagged `vX.Y.Z`.

## [Unreleased]

## [0.1.1] - 2026-09-19

### Fixed

- **WebSocket origin verification**: the events endpoint no longer sets `InsecureSkipVerify`. Origins are now validated with `OriginPatterns` derived from `allowed_origins` (wildcard `*` still opts into permissive behavior).

### Changed

- **Atomic releases**: the release workflow now creates the release as a draft, attaches the native installers, and only then publishes it, preventing incomplete or duplicate releases on the same tag.
- **Reproducible builds**: Go binaries are built with `-trimpath` and `mod_timestamp`, so the same commit yields identical checksums.
- **Universal macOS app**: releases now include `SV_Printer_darwin_universal.app.zip` (amd64 + arm64 via `lipo`).
- **Asset labels**: release assets are labelled as GUI application vs headless CLI.

### Docs

- **Download guide**: README (EN/ES) and `documentation/README.md` now explain which release asset to download per OS/architecture (GUI `.app.zip` vs headless CLI `.tar.gz`), how to verify checksums, and the macOS Gatekeeper step.
- **Product name standardised**: "SV Printer" is now used consistently across code, docs, spec and license.
- **Badges fixed**: CI/release badge URLs corrected (`sv-printerer` → `sv-printer`).

## [0.1.0] - 2026-09-19

### Added

- **Native Installers**: automated generation of `SV_Printer_Setup_windows_amd64.exe` (NSIS) for Windows and `SV_Printer_darwin_<arch>.app.zip` for macOS directly from GitHub Releases. Native installers are distributed as compressed archives; the macOS `.app` is unsigned, so first launch requires **right-click → Open** (or `xattr -dr com.apple.quarantine "SV Printer.app"`).
- **System Tray GUI**: native status menu on macOS and Windows to easily copy the token, device ID, install licenses, and manage OS auto-start.
- **Dynamic License UI**: the system tray dynamically adapts, hiding pro features ("Instalar Licencia", "Copiar ID") for users who already have an active PRO license.
- **Native OS Auto-Start**: built-in registry injection on Windows and LaunchAgents on macOS to run the agent silently on boot.
- **Windows System Printers**: native integration with Windows Spooler API (without CGO) to automatically discover and print to system-installed printers.
- **CORS Support**: explicitly configurable `allowed_origins` in `config.json` (including `*` wildcard support) to secure web printing.
- **Interactive License Generator**: developer script (`scripts/generate_license.sh`) to easily prompt for customer data and output `.key` files.
- **Initial release of SV Printer**: cross-platform local agent (Go) that lets web applications print to thermal printers (USB, serial, TCP) via a local HTTP API.
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
- **Spanish mirrors**: `README_ES.md`, `documentation/integration_ES.md`, `documentation/licensing_ES.md`.
- **Repository standards**: `CONTRIBUTING`, `SECURITY`, `CODE_OF_CONDUCT` (EN/ES).
- **License**: Business Source License 1.1 (converts to Apache 2.0 on 2030-09-11).

### Changed

- **Product rename**: SV Printer Agent → SV Printer (identifier: `sv-printer`).
- **Config**: corrected `allowed-origins` example to `svtech.cl`.
