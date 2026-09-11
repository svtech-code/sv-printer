# SV Print

Cross-platform local agent written in **Go** that lets web applications print to
thermal printers connected to the user's machine — without relying on the
browser's print dialog.

```
Web Application  ── HTTP/WebSocket ──>  SV Print  ──>  Thermal printer
                                                             (USB / Serial / TCP)
```

The primary printing protocol for the MVP is **ESC/POS**.

> **Specification:** the canonical spec lives in
> [`documentation/spect.md`](./documentation/spect.md) (English). A Spanish
> mirror is available in [`documentation/spect_ES.md`](./documentation/spect_ES.md).

## Requirements

- Go (stable version; see [`go.mod`](./go.mod)).
- Dependencies: [`go.bug.st/serial`](https://github.com/bugst/go-serial) (serial
  transport) and [`coder/websocket`](https://github.com/coder/websocket)
  (WebSocket events). Everything else uses the Go standard library.

## Quickstart

```bash
# Build
go build ./cmd/sv-print

# Run the CLI
./sv-print version
./sv-print status
./sv-print printers
```

## Architecture

Clean Architecture with a clear separation of concerns:

```text
cmd/sv-print/            CLI entry point
internal/
  application/           Use cases (discovery, printing)
  domain/                Core models (printer, job, errors)
  infrastructure/        Transports (network/TCP) and discovery mechanisms
  interfaces/http/       Local HTTP API + security middleware
pkg/escpos/              ESC/POS command builder
```

The application layer does not depend on OS-specific APIs: printer detection is
abstracted behind `PrinterDiscovery` and I/O behind `PrinterTransport`.

## Local API

The agent exposes a local HTTP API on `127.0.0.1:9876` (configurable) and never
binds to `0.0.0.0` by default. Requests require a `Authorization: Bearer <token>`
header.

| Method | Endpoint | Description |
|---|---|---|
| GET | `/health` | Health check |
| GET | `/api/v1/info` | Agent info (name, version, platform, tier) |
| GET | `/api/v1/printers` | List detected printers |
| POST | `/api/v1/printers/discover` | Re-run discovery |
| GET | `/api/v1/printers/{id}` | Get a printer |
| POST | `/api/v1/printers/{id}/test` | Print a test receipt |
| POST | `/api/v1/print` | Create a raw print job (pro) |
| POST | `/api/v1/print/receipt` | Print a structured receipt |
| GET | `/api/v1/jobs/{id}` | Query a job |
| GET | `/api/v1/events` | WebSocket event stream (pro) |

Errors use structured codes and a JSON envelope: `{"error":{"code","message"}}`.

## Licensing

SV Print uses a freemium model:

| Mode | Activation | What you get |
|---|---|---|
| **Trial** (default) | none | Structured receipts with a watermark (`*** SV PRINT — LICENCIA DE PRUEBA ***` at the start and end) and a **50 prints/day** limit |
| **Beta** | signed license with an expiry | No watermark, no quota, until the expiry date |
| **Full** | signed license | All features, no limits |

**Pro features** (require a license): `POST /api/v1/print` (raw ESC/POS) and
`GET /api/v1/events` (WebSocket events).

### Issuing a license

Licenses are Ed25519-signed and verified offline by the agent against an
embedded public key. The signing (private) key is kept outside the repository.
See [`docs/licensing.md`](./docs/licensing.md) for the full guide.

```bash
# 1. Generate a keypair (once). Keep the private key secret.
go run ./cmd/sv-license keygen

# 2. Sign a license (e.g. a 90-day beta for a customer)
go run ./cmd/sv-license sign \
  -key "$SV_LICENSE_KEY" \
  -customer "ACME" \
  -tier beta \
  -expiry 2027-01-01T00:00:00Z \
  -features raw_print,websocket \
  > license.key

# 3. Install it next to the config (or pass -license / SV_PRINT_LICENSE)
./sv-print -license ./license.key
```

## Development

The workflow follows the project conventions recorded in memory:

- **TDD:** write tests after each task and keep them passing (`go test ./...`).
- **Standard library first:** avoid external packages when the Go stdlib suffices.
- **Conventional Commits:** commits use `type(scope): message` (English).
- **Spec-driven:** behavior/architecture changes go through the `sv-memory` spec
  flow before implementation.

Useful commands:

```bash
go fmt ./...
go vet ./...
go test ./...
go build ./...
```

## Status

Early MVP. The spec is the source of truth; see the Roadmap section
(§43) in [`documentation/spect.md`](./documentation/spect.md) for the phased plan.

## License

The agent is distributed under a freemium license (see [Licensing](#licensing)
above and [`docs/licensing.md`](./docs/licensing.md)). The source is proprietary;
see the repository owner for licensing and distribution terms.
