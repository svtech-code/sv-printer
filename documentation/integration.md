> **Language:** English | [Español](integration_ES.md)

# Integration Guide

How to run SV Printer, configure it, and integrate it with your web application
using the local HTTP API.

## Requirements

- Go (see [`go.mod`](../go.mod))
- A thermal printer connected via USB, serial, or network (TCP)

## Running the agent

```bash
# Build
go build -o sv-printer ./cmd/sv-printer

# Run (auto-discovers printers, generates a token on first run)
./sv-printer

# Or configure manually
./sv-printer -token mysecret -port 9876 -printer "Caja 1@192.168.1.100:9100"
```

On first run without `-token`, SV Printer generates one and prints it to stdout.
Save it — you'll need it for every API call.

### Configuration

| Flag | Env var | Default | Description |
|---|---|---|---|
| `-host` | `SV_PRINT_HOST` | `127.0.0.1` | Bind address |
| `-port` | `SV_PRINT_PORT` | `9876` | Bind port |
| `-token` | `SV_PRINT_TOKEN` | (generated) | Auth token |
| `-max-payload` | `SV_PRINT_MAX_PAYLOAD` | `5242880` (5 MB) | Max request body |
| `-origin` | `SV_PRINT_ALLOWED_ORIGINS` | (none) | Allowed CORS origin (repeatable) |
| `-printer` | — | (none) | Manual printer as `name@address` (repeatable) |
| `-serial` | — | (none) | Manual serial printer as `name@port[@baud]` (repeatable) |
| `-config` | `SV_PRINT_CONFIG` | platform path | Config file path |
| `-license` | `SV_PRINT_LICENSE` | `license.key` next to config | License file path |

The config file is JSON (`sv-printer config` shows its path).

### Subcommands

| Command | Description |
|---|---|
| `sv-printer version` | Print version |
| `sv-printer status` | Check agent status |
| `sv-printer printers` | List configured printers |
| `sv-printer config` | Show current config |
| `sv-printer discover` | Re-run printer discovery |
| `sv-printer test <printer-id>` | Send a test receipt |
| `sv-printer print <printer-id> <file.json>` | Print a structured receipt from a JSON file |
| `sv-printer logs` | Tail log file |
| `sv-printer doctor` | Run diagnostics |
| `sv-printer device-id` | Show device fingerprint (for license binding) |
| `sv-printer service <install\|uninstall\|status>` | Manage auto-start at login |

## Security model

- The agent binds to **`127.0.0.1` only** (never `0.0.0.0` by default).
- All `/api/*` endpoints require `Authorization: Bearer <token>`.
- CORS: only origins listed in `-origin` / `allowed_origins` can call the API
  from a browser. Configure this to your web app's origin (e.g. `https://svtech.cl`).
- WebSocket auth: the native browser WebSocket API cannot set headers, so
  `?token=<token>` is accepted as an alternative for `GET /api/v1/events`.

## Quickstart (curl)

```bash
TOKEN="your-token-here"
BASE="http://127.0.0.1:9876"

# Health (no auth)
curl http://127.0.0.1:9876/health

# Info
curl -H "Authorization: Bearer $TOKEN" $BASE/api/v1/info

# List printers
curl -H "Authorization: Bearer $TOKEN" $BASE/api/v1/printers

# Test receipt (trial: adds watermark, counts toward quota)
curl -X POST -H "Authorization: Bearer $TOKEN" \
  $BASE/api/v1/printers/<printer-id>/test
```

## Printing a structured receipt

`POST /api/v1/print/receipt` accepts a JSON body describing the receipt layout.

### Request schema

```json
{
  "printer_id": "net-192.168.1.100:9100",
  "cut": true,
  "lines": [
    { "text": "MY STORE", "style": { "bold": true, "align": "center" } },
    { "text": "123 Main St" },
    { "text": "-------------------" },
    { "text": "Item 1", "style": { "align": "left" } },
    { "text": "$10.00", "style": { "align": "right" } },
    { "text": "TOTAL", "style": { "bold": true, "align": "center", "size": [2, 2] } }
  ]
}
```

### Fields

| Field | Type | Description |
|---|---|---|
| `printer_id` | string | **Required.** Printer ID (from `/api/v1/printers`) |
| `cut` | bool | Feed and cut paper after printing |
| `lines` | array | **Required.** Array of line objects |

### Line object

| Field | Type | Description |
|---|---|---|
| `text` | string | **Required.** Text to print |
| `style.bold` | bool | Bold text |
| `style.underline` | bool | Underlined text |
| `style.italic` | bool | Italic text |
| `style.reverse` | bool | Inverted (white on black) |
| `style.align` | string | `left` (default), `center`, or `right` |
| `style.size` | `[w, h]` | Character size multiplier `[width, height]` (1–8 each) |

### Example (curl)

```bash
curl -X POST -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "printer_id": "net-192.168.1.100:9100",
    "cut": true,
    "lines": [
      { "text": "ORDER #42", "style": { "bold": true, "align": "center" } },
      { "text": "Coffee x2", "style": { "align": "left" } },
      { "text": "$8.00", "style": { "align": "right" } },
      { "text": "THANK YOU", "style": { "align": "center" } }
    ]
  }' \
  $BASE/api/v1/print/receipt
```

Response (201 Created):

```json
{ "job_id": "job_f3a9c1d2e4b5061789a0b2c3d4e5f607", "status": "queued" }
```

## Raw ESC/POS printing (pro)

`POST /api/v1/print` sends raw ESC/POS bytes. Requires a license (pro feature).

```json
{
  "printer_id": "net-192.168.1.100:9100",
  "payload": "AQIDBA==",
  "format": "escpos"
}
```

`payload` is a **base64-encoded** string of raw ESC/POS commands.

## Job status

`GET /api/v1/jobs/{id}` returns:

```json
{
  "id": "job_f3a9c1d2e4b5061789a0b2c3d4e5f607",
  "printer_id": "net-192.168.1.100:9100",
  "created_at": "2026-09-11T12:00:00Z",
  "status": "completed"
}
```

Status values: `queued` → `printing` → `completed` | `failed` | `cancelled`.

## WebSocket events (pro)

`GET /api/v1/events` opens a WebSocket connection. Requires a license (pro feature).

### Authentication

Browsers cannot set custom headers on WebSocket, so use a query parameter:

```
ws://127.0.0.1:9876/api/v1/events?token=your-token
```

Node.js / non-browser clients can use the `Authorization` header:

```
ws://127.0.0.1:9876/api/v1/events
Authorization: Bearer your-token
```

### Event types

| Event | Data |
|---|---|
| `agent.status` | `{ "status": "started" }` / `{ "status": "stopped" }` |
| `print.started` | `{ "job_id": "...", "printer_id": "..." }` |
| `print.completed` | `{ "job_id": "...", "printer_id": "..." }` |
| `print.failed` | `{ "job_id": "...", "printer_id": "...", "error": "..." }` |
| `printer.discovered` | `{ "id": "...", "name": "..." }` |
| `printer.removed` | `{ "id": "..." }` |
| `printer.status_changed` | `{ "id": "...", "status": "..." }` |

### Example (JavaScript)

```js
const ws = new WebSocket(
  "ws://127.0.0.1:9876/api/v1/events?token=your-token"
);

ws.onmessage = (e) => {
  const event = JSON.parse(e.data);
  console.log("Event:", event.event, event);
};
```

## Error responses

All errors follow a consistent JSON envelope:

```json
{
  "error": {
    "code": "PRINTER_NOT_FOUND",
    "message": "Printer not found."
  }
}
```

### Error codes

| HTTP | Code | Meaning |
|---|---|---|
| 400 | `INVALID_PAYLOAD` | Malformed request body |
| 400 | `UNSUPPORTED_PROTOCOL` | Unknown `format` in raw print |
| 401 | `UNAUTHORIZED` | Missing or invalid token |
| 403 | `LICENSE_REQUIRED` | Pro feature without a license |
| 403 | `ORIGIN_NOT_ALLOWED` | CORS origin not in allowed list |
| 404 | `PRINTER_NOT_FOUND` | Unknown printer ID |
| 404 | `JOB_NOT_FOUND` | Unknown job ID |
| 409 | `PRINTER_BUSY` | Printer is busy |
| 409 | `PRINTER_PAPER_OUT` | Printer is out of paper |
| 409 | `JOB_ALREADY_EXISTS` | Duplicate job ID |
| 413 | `PAYLOAD_TOO_LARGE` | Request body exceeds limit |
| 429 | `LICENSE_QUOTA_EXCEEDED` | Daily trial limit reached (50/day) |
| 502 | `PRINTER_CONNECTION_FAILED` | Cannot connect to printer |
| 503 | `PRINTER_OFFLINE` | Printer is offline |
| 503 | `EVENTS_UNAVAILABLE` | Event bus not configured |
| 500 | `PRINT_FAILED` | Print job failed |
| 500 | `INTERNAL_ERROR` | Unexpected server error |

## Licensing / trial notes

Without a license (trial mode):

- Receipts are prefixed and suffixed with
  `*** SV PRINTER — LICENCIA DE PRUEBA ***` and `www.svtech.cl`.
- You can print up to **50 receipts per day** (quota resets at midnight).
- `POST /api/v1/print` (raw ESC/POS) and `GET /api/v1/events` (WebSocket)
  are **blocked** (403 `LICENSE_REQUIRED`).

To install a license:

```bash
./sv-printer -license ./license.key
```

## Troubleshooting

```bash
# Run diagnostics
sv-printer doctor

# Check logs
sv-printer logs

# Show device fingerprint (for license binding)
sv-printer device-id
```
