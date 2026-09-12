> **Language:** English | [Español](licensing_ES.md)

# Licensing

SV Print uses a **freemium** model with offline, Ed25519-signed licenses.

## Modes

| Mode | Activation | Watermark | Daily quota | Pro features |
|---|---|---|---|---|
| **Trial** | none (default) | yes, at the start and end of every receipt | 50 prints/day | no |
| **Beta** | signed license with `expiry` | no | unlimited | depends on `features` |
| **Full** | signed license | no | unlimited | depends on `features` |

The trial watermark is `*** SV PRINT — LICENCIA DE PRUEBA ***`.

## Gated features

| Feature | What it unlocks |
|---|---|
| `raw_print` | `POST /api/v1/print` (raw ESC/POS payloads) |
| `websocket` | `GET /api/v1/events` (WebSocket event stream) |

A license must list a feature in its `features` array to enable it. Without a
license (trial mode), only the structured receipt endpoint
(`POST /api/v1/print/receipt`) and the printer/test endpoints are available.

## How licenses work

1. A **signing keypair** is generated once. The **public key** is embedded in the
   agent binary (`internal/license/public.key`); the **private key** is kept
   secret and is never committed to the repository.
2. Each license is a JSON envelope:
   ```json
   {
     "payload": "<base64-encoded claims>",
     "signature": "<base64 Ed25519 signature over the payload>"
   }
   ```
   The claims are `license_id`, `product`, `customer`, `tier`, `features`,
   `expiry` (optional RFC3339), and `fingerprint` (optional).
3. The agent verifies the signature against the embedded public key at startup.
   A missing, tampered, expired, or wrong-product license falls back to trial mode.

## Issuing a license

```bash
# 1. Generate the keypair (once).
go run ./cmd/sv-license keygen
#   -> private key (hex) printed to stderr: keep it in 1Password / a secret manager
#   -> public key (hex) printed to stdout: place it in internal/license/public.key

# 2. Sign a license.
SV_LICENSE_KEY=<private-key-hex> go run ./cmd/sv-license sign \
  -customer "ACME" \
  -tier full \
  -features raw_print,websocket \
  > license.key

# With an expiry (e.g. a beta):
SV_LICENSE_KEY=<private-key-hex> go run ./cmd/sv-license sign \
  -customer "ACME" \
  -tier beta \
  -expiry 2027-01-01T00:00:00Z \
  -features raw_print,websocket \
  > license.key
```

## Installing a license on the agent

The agent loads its license from, in order:

1. the `-license` flag (`sv-print -license ./license.key`),
2. the `SV_PRINT_LICENSE` environment variable,
3. `license.key` next to the config file (default location).

The current license state is exposed by `GET /api/v1/info` via the `tier` and
`licensed` fields.

## Device binding

A license can be bound to a single machine by setting the `fingerprint` claim.
When a fingerprint is present, the agent verifies it against the machine's
device ID (via [`machineid`](https://github.com/denisbrodbeck/machineid)) and
falls back to trial mode if it does not match.

Flow:

1. The customer runs `sv-print device-id` and sends you the fingerprint.
2. You sign the license with `-fingerprint <value>`.
3. The license only activates on that machine.

```bash
# Customer:
sv-print device-id
# -> e.g. 3f8c2a...

# Vendor:
go run ./cmd/sv-license sign \
  -customer "ACME" \
  -tier full \
  -features raw_print,websocket \
  -fingerprint "3f8c2a..." \
  > license.key
```

## Security notes

- The **private key must never** be committed or shared. Anyone with it can sign
  arbitrary licenses.
- For a public repository, this is the core control: even with the full source,
  a license cannot be forged without the private key.
- **Device binding** prevents sharing a license across machines, but it is
  software-enforced and can be defeated by an attacker who patches the binary.
- Offline-only licensing cannot revoke a license or prevent clock rollback.
  If revocation or hard metering becomes a requirement, add a phone-home
  validation step (see the roadmap).
