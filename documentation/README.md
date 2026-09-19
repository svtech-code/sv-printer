# Documentation

This directory holds supporting project documentation.

## Source of truth

- **Specification (canonical):** [`../documentation/spect.md`](../documentation/spect.md)
  (English) — the authoritative system specification.
- **Specification (Spanish mirror):** [`../documentation/spect_ES.md`](../documentation/spect_ES.md).

## Getting started

See [`../README.md`](../README.md) for build instructions, architecture overview,
and the local API contract. Also available in [Español](../README_ES.md).

## Which release to download

Releases publish several artifacts. Pick the one that matches your goal and OS:

| Goal | Platform | Asset |
|:---|:---|:---|
| Desktop app (tray + icon) | Windows | `SV_Printer_Setup_windows_amd64.exe` |
| Desktop app (tray + icon) | macOS Apple Silicon | `SV_Printer_darwin_arm64.app.zip` |
| Desktop app (tray + icon) | macOS Intel | `SV_Printer_darwin_amd64.app.zip` |
| Desktop app (tray + icon) | macOS (any) | `SV_Printer_darwin_universal.app.zip` |
| Headless CLI | Linux / macOS / Windows | `sv-printer_<version>_<os>_<arch>.tar.gz` (zip on Windows) |

- The `.app.zip` archives contain the full `SV Printer.app` bundle (GUI with icon).
- The `sv-printer_*.<os>_*` archives contain **only the headless CLI binary** (no icon).
- Check your Mac architecture with `uname -m` (`arm64` = Apple Silicon, `x86_64` = Intel).
- The macOS app is not code-signed; on first launch use **right-click → Open** or
  `xattr -dr com.apple.quarantine "SV Printer.app"`.
- Verify downloads with `checksums.txt` (`shasum -a 256 -c checksums.txt`).
- Always download from the **Latest** release, never from a draft.
- The universal macOS bundle is included since `v0.1.1`.

## Integration

See [`integration.md`](./integration.md) ([Español](./integration_ES.md)) for full
API documentation, receipt schema, WebSocket events, error codes, and client
examples (JavaScript, Python, PHP).

## Licensing

See [`licensing.md`](./licensing.md) ([Español](./licensing_ES.md)) for the
freemium model, gated features, and how to issue and install licenses.

## Memory

Design knowledge, decisions, and standards are captured in the project memory
store (`.sv-memory/`) and managed via the `sv-memory` toolchain.
