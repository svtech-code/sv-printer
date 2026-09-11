# SV Print

> **Canonical version:** This document is the canonical specification. The Spanish
> version lives in [`spect_ES.md`](./spect_ES.md) and is a mirror for reference only.

## 1. Overview

**SV Print** is a cross-platform local agent developed in **Go**, whose goal is to allow web applications to communicate with thermal printers connected to the user's machine.

The agent acts as a bridge between:

```text
Web Application
      │
      │ HTTP / WebSocket
      ▼
SV Print
      │
      ├── USB
      ├── Serial
      ├── TCP/IP
      └── System printer
      │
      ▼
Thermal printer
```

The main objective is to enable thermal printing from web applications without relying on the browser's print dialog.

The agent must run on:

- Windows
- macOS
- Linux

The primary protocol for the MVP will be **ESC/POS**.

---

# 2. Objectives

## 2.1 Primary objectives

The system must:

- Detect printers available locally.
- Identify the connection type of each printer.
- Allow configuring printers manually.
- Expose a local API for web applications.
- Print documents using ESC/POS.
- Support USB printers.
- Support serial printers.
- Support network printers.
- Query printer status when the hardware/protocol allows it.
- Run on Windows, macOS and Linux.
- Run as a background process or service.
- Provide a CLI for diagnostics and administration.
- Protect the API with authentication.
- Maintain an architecture independent of any web framework.

---

## 2.2 Secondary objectives

The project must be prepared to incorporate later:

- Advanced automatic discovery.
- Persistent print queue.
- Automatic retries.
- Printer monitoring.
- WebSocket events.
- Multiple printers.
- Printer aliases.
- Persistent configuration.
- Job history.
- CUPS integration.
- Windows Print Spooler integration.
- Graphical interface.
- System Tray.
- Remote administration.

---

# 3. MVP Scope

The MVP must focus on the essential features.

### Operating systems

```text
Windows
macOS
Linux
```

### Connections

```text
USB
Serial
TCP/IP
```

### Features

```text
✓ Printer detection
✓ Manual configuration
✓ Printer listing
✓ Print test
✓ ESC/POS printing
✓ Print queue
✓ Job status
✓ Local HTTP API
✓ Authentication
✓ CORS restriction
✓ CLI
✓ Structured logs
✓ Health check
```

---

# 4. Technologies

## 4.1 Language

The project must be developed in:

```text
Go
```

Use the stable Go version available when development starts.

---

## 4.2 Technical principles

The project must prioritize:

- Simplicity.
- Portability.
- Low resource consumption.
- Security.
- Testability.
- Separation of concerns.
- Small interfaces.
- Idiomatic Go code.
- Minimal dependencies.
- Avoid CGO when technically possible.

---

# 5. Architecture

A modular architecture inspired by Clean Architecture will be used.

Proposed structure:

```text
sv-print/
├── cmd/
│   └── sv-print/
│       └── main.go
│
├── internal/
│   ├── application/
│   │   ├── discovery/
│   │   ├── printing/
│   │   └── health/
│   │
│   ├── domain/
│   │   ├── printer/
│   │   ├── job/
│   │   └── errors/
│   │
│   ├── infrastructure/
│   │   ├── usb/
│   │   ├── serial/
│   │   ├── network/
│   │   ├── system/
│   │   ├── storage/
│   │   └── logging/
│   │
│   └── interfaces/
│       ├── http/
│       ├── websocket/
│       └── cli/
│
├── pkg/
│   └── escpos/
│
├── docs/
├── scripts/
├── test/
│
├── go.mod
└── README.md
```

The structure may be modified during development if there is a technical reason, but the separation between domain, application, infrastructure and interfaces must be maintained.

---

# 6. High-level architecture

```text
                         WEB APPLICATION
                               │
                               │ HTTP
                               ▼
                        LOCAL SV PRINT API
                               │
                     ┌─────────┴─────────┐
                     │                   │
                    HTTP             WebSocket
                     │                   │
                     └─────────┬─────────┘
                               │
                       APPLICATION LAYER
                               │
                 ┌─────────────┴─────────────┐
                 │                           │
           Printer                     Print queue
          discovery                        │
                 │                    Print Worker
                 │                           │
                 └─────────────┬─────────────┘
                               │
                     Printer abstraction
                               │
             ┌─────────────────┼─────────────────┐
             │                 │                 │
            USB             Serial              TCP
             │                 │                 │
             └─────────────────┼─────────────────┘
                               │
                        THERMAL PRINTER
```

---

# 7. Cross-platform support

## 7.1 Windows

Target versions:

```text
Windows 10+
Windows 11
```

Must support:

- USB detection.
- Serial detection.
- Network printers.
- Future Windows Print Spooler integration.
- Running as a Windows Service.

Binary:

```text
sv-print.exe
```

Future installer:

```text
sv-print-installer.exe
```

---

## 7.2 macOS

Initial minimum version:

```text
macOS 12+
```

Must support:

- USB.
- Serial.
- TCP/IP.
- CUPS.
- Running in the background.

Architectures:

```text
amd64
arm64
```

Future distribution:

```text
SV Print.app
.dmg
.pkg
```

Production distribution must consider:

- Code signing.
- Notarization.
- System permissions.

---

## 7.3 Linux

Initial distributions:

```text
Ubuntu
Debian
```

Architectures:

```text
amd64
arm64
```

Must support:

- USB.
- Serial.
- TCP/IP.
- CUPS.
- systemd.

Future package:

```text
.deb
```

Possible later formats:

```text
.rpm
AppImage
```

---

# 8. Printer detection

Printer detection is a fundamental feature.

The agent must provide a discovery service:

```text
Printer Discovery
       │
       ├── USB
       ├── Serial
       ├── Network
       └── System
```

Each mechanism must implement a common interface.

Example:

```go
type PrinterDiscovery interface {
    Discover(ctx context.Context) ([]Printer, error)
}
```

The application layer must not depend directly on Windows-, macOS- or Linux-specific APIs.

---

# 9. USB detection

The agent must detect printers connected via USB.

When possible it must obtain:

- Vendor ID.
- Product ID.
- Manufacturer.
- Model.
- Serial number.
- Device path.
- Interface information.

Example:

```json
{
  "id": "usb-1234-5678",
  "name": "XPrinter XP-Q200",
  "connection": "usb",
  "manufacturer": "XPrinter",
  "model": "XP-Q200"
}
```

The implementation must use a maintained Go library compatible with the three operating systems.

The final library must be selected considering:

- Cross-platform compatibility.
- Maintenance.
- License.
- CGO requirement.
- Thermal device compatibility.

---

# 10. Serial detection

Must detect serial devices.

### Windows

```text
COM3
COM4
```

### Linux

```text
/dev/ttyUSB0
/dev/ttyACM0
```

### macOS

```text
/dev/cu.usbserial-XXXX
/dev/cu.usbmodemXXXX
```

Configuration:

```text
Baud rate
Data bits
Stop bits
Parity
Flow control
```

Recommended initial configuration:

```text
9600
8
N
1
```

Must be configurable.

---

# 11. Network printers

Raw printing over TCP must be supported.

Initial port:

```text
9100
```

Example:

```text
192.168.1.100:9100
```

Configuration:

```json
{
  "name": "Caja 1",
  "connection": "network",
  "address": "192.168.1.100:9100",
  "protocol": "escpos"
}
```

Advanced automatic discovery via:

```text
mDNS
DNS-SD
SNMP
```

is initially a future feature.

Indiscriminate network scanning must not be performed by default.

---

# 12. Printer model

All printers must be represented by a common model.

Example:

```go
type Printer struct {
    ID           string
    Name         string
    Manufacturer string
    Model        string
    Connection   ConnectionType
    Address      string
    Status       PrinterStatus
    Protocol     PrinterProtocol
    IsDefault    bool
}
```

Connection types:

```text
usb
serial
network
system
unknown
```

Protocols:

```text
escpos
cups
windows
raw_tcp
unknown
```

States:

```text
unknown
ready
printing
offline
paper_out
error
busy
```

The system must not assume that all printers provide detailed status information.

---

# 13. Printer identification

IDs must be as stable as possible.

Priority:

1. Device serial number.
2. Vendor ID + Product ID + path.
3. Network address.
4. Operating system identifier.

Random IDs must not be generated on every discovery run if a stable identification is available.

---

# 14. Transport abstraction

The application must not depend directly on USB, TCP or Serial.

Proposed interface:

```go
type PrinterTransport interface {
    Open(ctx context.Context) error
    Write(ctx context.Context, data []byte) error
    Close() error
}
```

Implementations:

```text
USBTransport
SerialTransport
TCPTransport
SystemPrinterTransport
```

This will allow adding new connection mechanisms without modifying the core printing logic.

---

# 15. ESC/POS

ESC/POS will be the primary MVP protocol.

There must be an independent package:

```text
pkg/escpos
```

It must allow building print jobs.

Example:

```go
receipt := escpos.NewReceipt()

receipt.
    Center().
    Bold().
    Text("SV TECH").
    LineFeed().
    ResetStyle().
    Text("------------------------------").
    LineFeed().
    Text("Servicio soporte     $25.000").
    LineFeed().
    Text("------------------------------").
    Bold().
    Text("TOTAL                $25.000").
    LineFeed().
    Cut()
```

Initial functions:

- Initialize printer.
- Text.
- Bold.
- Italic when supported.
- Underline.
- Alignment.
- Font size.
- Spacing.
- Line feeds.
- Tables.
- Images.
- QR.
- Barcodes.
- Paper cut.
- Cash drawer pulse.
- Reset.

---

# 16. Print job

Every print must be represented as a job.

Example:

```go
type PrintJob struct {
    ID        string
    PrinterID string
    Payload   []byte
    CreatedAt time.Time
    Status    PrintJobStatus
}
```

States:

```text
queued
printing
completed
failed
cancelled
```

Every job must receive a unique ID.

Example:

```text
job_01KXYZ
```

---

# 17. Print queue

The MVP must use an in-memory queue.

```text
HTTP Request
      │
      ▼
Print queue
      │
      ▼
Print Worker
      │
      ▼
Printer
```

For a physical printer, sending several jobs simultaneously must be avoided.

Later versions may incorporate:

- Persistent queue.
- Retries.
- Priorities.
- Cancellation.
- History.
- Recovery after restarting the agent.

---

# 18. Local API

The agent must expose a local HTTP API.

Default address:

```text
127.0.0.1:9876
```

The port must be configurable.

For security, the agent must listen only on localhost by default.

It must not be exposed to:

```text
0.0.0.0
```

nor to the Internet except by explicit configuration.

---

# 19. Endpoints

## Health

```http
GET /health
```

Response:

```json
{
  "status": "ok",
  "version": "0.1.0"
}
```

---

## Agent information

```http
GET /api/v1/info
```

Response:

```json
{
  "name": "sv-print",
  "version": "0.1.0",
  "platform": "darwin",
  "architecture": "arm64"
}
```

---

## List printers

```http
GET /api/v1/printers
```

Response:

```json
{
  "printers": [
    {
      "id": "usb-1234-5678",
      "name": "POS-58",
      "connection": "usb",
      "protocol": "escpos",
      "status": "ready"
    }
  ]
}
```

---

## Discover printers

```http
POST /api/v1/printers/discover
```

Must run the discovery process again.

---

## Get printer

```http
GET /api/v1/printers/{printerID}
```

---

## Print test

```http
POST /api/v1/printers/{printerID}/test
```

Must generate a standard test receipt.

---

## Create print job

```http
POST /api/v1/print
```

Example:

```json
{
  "printer_id": "usb-1234-5678",
  "format": "escpos",
  "payload": "..."
}
```

Response:

```json
{
  "job_id": "job_01KXYZ",
  "status": "queued"
}
```

---

## Query job

```http
GET /api/v1/jobs/{jobID}
```

Response:

```json
{
  "id": "job_01KXYZ",
  "status": "completed"
}
```

---

# 20. WebSocket

Support for WebSocket events must be left prepared.

Endpoint:

```text
ws://127.0.0.1:9876/api/v1/events
```

Events:

```text
printer.discovered
printer.removed
printer.status_changed
print.started
print.completed
print.failed
agent.status
```

Example:

```json
{
  "event": "printer.status_changed",
  "printer_id": "usb-1234-5678",
  "status": "paper_out"
}
```

It may be considered experimental during the MVP.

---

# 21. Authentication

The local API must use authentication.

During the first initialization a token must be generated.

Example:

```text
SV Print initialized.

Agent ID:
agent_xxxxxxxxx

Token:
xxxxxxxxxxxxxxxxxxxxxxxx
```

Requests will use:

```http
Authorization: Bearer <token>
```

Tokens must not appear in logs.

When possible, credentials must be stored using secure operating system mechanisms.

---

# 22. CORS

The API must implement CORS restrictions.

It must not use:

```text
Access-Control-Allow-Origin: *
```

as the default configuration.

There must be a list of authorized origins.

Example:

```text
https://svtech.cl
https://localhost:3000
```

The list must be configurable.

---

# 23. Security

Mandatory principles:

1. Listen on localhost by default.
2. Use authentication.
3. Validate CORS origins.
4. Validate all payloads.
5. Limit job size.
6. Limit queue size.
7. Do not allow arbitrary filesystem access.
8. Do not execute system commands sent from the web.
9. Do not allow arbitrary network connections.
10. Do not expose the agent directly to the Internet.
11. Do not log credentials.
12. Do not log sensitive customer information.
13. Validate printer IDs.
14. Avoid SSRF through controlled network configurations.

---

# 24. Limits

The agent must have limits to avoid abuse or excessive resource consumption.

Suggested initial values:

```text
Maximum payload size: 5 MB
Maximum queue size: 100 jobs
Simultaneous jobs per printer: 1
```

These values must be configurable.

---

# 25. CLI

The MVP does not require a graphical interface.

It must provide an administrative CLI.

Examples:

```bash
sv-print status

sv-print printers

sv-print discover

sv-print test <printer-id>

sv-print print <printer-id> receipt.json

sv-print config

sv-print logs

sv-print doctor

sv-print version
```

Example:

```bash
sv-print printers
```

Result:

```text
SV Print

Printers
────────────────────────────────────────────

ID              NAME            CONNECTION

usb-1234        POS-58          USB
net-001         Epson TM20      TCP 192.168.1.50:9100

Status
────────────────────────────────────────────

POS-58          READY
Epson TM20      OFFLINE
```

---

# 26. Graphical UI or CLI?

The first version must use a **CLI and not a graphical UI**.

The web will be the main interface for the end user.

The CLI will be intended for:

- Installation.
- Diagnostics.
- Configuration.
- Testing.
- Technical support.
- Troubleshooting.

The agent's core must remain completely independent of the CLI.

---

# 27. Possible future UI

If the user experience shows that configuration becomes complex, a graphical interface may be incorporated.

The future architecture could be:

```text
              SV Print UI
                  │
               Tauri
                  │
                  ▼
            SV Print Core
                  │
        ┌─────────┼─────────┐
        │         │         │
       USB      Serial      TCP
```

The UI must not contain the core printing logic.

The Go core must continue working independently.

---

# 28. Configuration

Configuration must be stored outside the binary.

Example:

```yaml
server:
  host: 127.0.0.1
  port: 9876

security:
  allowed_origins:
    - https://svtech.cl

printers:
  - id: usb-1234-5678
    name: Caja 1
    protocol: escpos

logging:
  level: info
```

The file location must follow the conventions of each operating system.

---

# 29. Logs

Structured logs will be used.

Recommended:

```text
log/slog
```

Levels:

```text
debug
info
warn
error
```

Example:

```text
INFO printer discovered
    printer_id=usb-1234
    connection=usb
    model=POS-58
```

Must never log:

- Tokens.
- Passwords.
- Credentials.
- Full payloads that may contain sensitive information.
- Unnecessary personal data.

---

# 30. Diagnostics

There must be:

```bash
sv-print doctor
```

Example:

```text
SV Print Doctor

✓ Operating system detected
✓ Valid configuration
✓ Local API available
✓ USB subsystem available
✓ Network available
✓ Printer detected

Printers:

✓ POS-58
```

If there is a problem:

```text
✗ POS-58

Reason:
Could not establish connection.

Suggested actions:
- Check USB cable.
- Check power.
- Check permissions.
- Run detection again.
```

---

# 31. Installation as a service

The agent must start automatically with the operating system.

## Windows

Use:

```text
Windows Service
```

## Linux

Use:

```text
systemd
```

## macOS

Use:

```text
launchd
```

---

# 32. Lifecycle

The CLI must allow:

```bash
sv-print service start

sv-print service stop

sv-print service restart

sv-print service status
```

The specific implementation of each operating system must be isolated from the core logic.

---

# 33. Printer monitoring

The agent must be able to monitor configured printers.

Example:

```text
READY
  │
  ▼
PRINTING
  │
  ▼
READY
```

Error:

```text
READY
  │
  ▼
OFFLINE
  │
  ▼
READY
```

The monitoring interval must be configurable.

The availability of states such as `paper_out` will depend on the capabilities of each printer and protocol.

---

# 34. Error handling

Errors must use structured codes.

Examples:

```text
PRINTER_NOT_FOUND
PRINTER_OFFLINE
PRINTER_BUSY
PRINTER_PAPER_OUT
PRINTER_CONNECTION_FAILED
PRINT_FAILED
INVALID_PAYLOAD
UNAUTHORIZED
ORIGIN_NOT_ALLOWED
PAYLOAD_TOO_LARGE
UNSUPPORTED_PROTOCOL
```

Response:

```json
{
  "error": {
    "code": "PRINTER_OFFLINE",
    "message": "Printer is currently offline."
  }
}
```

The web application will be responsible for translating messages to the user's language via i18n.

---

# 35. Testing

## Unit

Must cover:

- Entities.
- Validations.
- ESC/POS.
- Queue.
- Jobs.
- Printer identification.
- Configuration.
- Authentication.
- Error handling.

## Integration

Must test:

- HTTP API.
- WebSocket.
- TCP.
- Serial when possible.
- Complete printing flow.

## Hardware

Real tests with printers must remain separate from unit tests and run on physical hardware.

---

# 36. CI/CD

The project must use GitHub Actions.

Build matrix:

```text
windows/amd64
windows/arm64

linux/amd64
linux/arm64

darwin/amd64
darwin/arm64
```

The pipeline must run:

```text
go fmt
go vet
go test
go build
```

May be added later:

- golangci-lint.
- security analysis.
- integration tests.
- automatic release generation.

---

# 37. Release artifacts

Each release must generate:

```text
sv-print-windows-amd64.exe
sv-print-windows-arm64.exe

sv-print-linux-amd64
sv-print-linux-arm64

sv-print-darwin-amd64
sv-print-darwin-arm64
```

Later:

```text
Windows Installer
.deb
.rpm
AppImage
.dmg
.pkg
```

Each release must generate checksums.

---

# 38. Versioning

Use:

```text
Semantic Versioning
```

Format:

```text
MAJOR.MINOR.PATCH
```

Example:

```text
0.1.0
```

During initial development the project will remain at version `0.x`.

---

# 39. API versioning

All public endpoints will use:

```text
/api/v1/
```

Example:

```text
/api/v1/printers
```

Incompatible changes will require:

```text
/api/v2/
```

Compatible changes must remain within the same API version.

---

# 40. Main flow

The expected flow for a user will be:

```text
1. Install SV Print
             │
             ▼
2. Agent starts automatically
             │
             ▼
3. Detects printers
             │
             ▼
4. Web application detects Agent
             │
             ▼
5. Application gets printers
             │
             ▼
6. User selects printer
             │
             ▼
7. Application generates job
             │
             ▼
8. Application sends job
             │
             ▼
9. SV Print processes queue
             │
             ▼
10. Agent sends ESC/POS
             │
             ▼
11. Printer prints
             │
             ▼
12. Agent reports result
```

---

# 41. Integration example

The web application queries:

```http
GET http://127.0.0.1:9876/health
```

Then:

```http
GET http://127.0.0.1:9876/api/v1/printers
```

It gets:

```json
{
  "printers": [
    {
      "id": "usb-1234",
      "name": "Caja 1",
      "status": "ready"
    }
  ]
}
```

Finally:

```http
POST http://127.0.0.1:9876/api/v1/print
```

The agent:

```text
Validate request
        │
        ▼
Authenticate
        │
        ▼
Validate printer
        │
        ▼
Create Print Job
        │
        ▼
Add to queue
        │
        ▼
Open connection
        │
        ▼
Send ESC/POS
        │
        ▼
Close connection
        │
        ▼
Update status
```

---

# 42. Compatibility considerations

Compatibility with a printer must not be determined solely by its brand.

Must consider:

```text
Brand
Model
Connection type
Protocol
Supported ESC/POS commands
Operating system
Required driver
```

The project must maintain a progressive compatibility strategy:

```text
Standard ESC/POS
       ↓
Common extensions
       ↓
Manufacturer-specific commands
```

Manufacturer-specific commands must remain isolated.

---

# 43. Roadmap

## Phase 1 — MVP

```text
Go
CLI
HTTP API
USB
TCP
ESC/POS
Windows
macOS
Linux
```

## Phase 2

```text
Serial
WebSocket
Monitoring
Persistent configuration
Installers
System services
```

## Phase 3

```text
CUPS
Windows Print Spooler
Advanced discovery
SNMP
mDNS / DNS-SD
```

## Phase 4

```text
Graphical UI
System Tray
Friendly configuration
```

## Phase 5

```text
SV Print Cloud
Device management
Remote configuration
Licensing
Enterprise features
```

---

# 44. Design principles

The project must respect:

- Idiomatic Go code.
- Small interfaces.
- Separation of concerns.
- Dependency injection when appropriate.
- Inverted dependency.
- Cross-platform code.
- Isolation of operating system-specific code.
- Explicit error handling.
- Use of `context.Context`.
- Concurrency-safe operations.
- Avoid mutable global state.
- Testable code.
- Secure-by-default configuration.
- Minimal external dependencies.

---

# 45. Expected result

The final result must be a small, reliable agent that allows any compatible web application to:

```text
Detect printer
       ↓
Select printer
       ↓
Send job
       ↓
Print
       ↓
Query result
```

without the application having to know the specific details of:

```text
Windows
macOS
Linux
USB
Serial
Ethernet
Wi-Fi
ESC/POS
```

The local API must constitute the primary integration contract.

The CLI must constitute the primary technical administration tool.

The graphical UI remains a possible future layer and must not be part of the project's core.
