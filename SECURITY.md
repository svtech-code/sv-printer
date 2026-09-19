> **Language:** English | [Español](SECURITY_ES.md)

# Security Policy

We take the security of **SV Printer** seriously. Since this agent handles local network access, printer communication, and license verification, safeguarding its deployment is our top priority.

---

## 1. Supported Versions

We actively support and patch security issues for the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 0.1.x   | :white_check_mark: |
| < 0.1.0 | :x:                |

> Policy: **latest minor only.** Security fixes land on the current minor release. Older minors are expected to upgrade.

---

## 2. Reporting a Vulnerability

**Please do not open a public issue on GitHub for security vulnerabilities.**

If you discover a security vulnerability, please report it privately:

1. **GitHub Security Advisory:** Go to the **Security** tab of the repository and select **Advisories → New draft advisory**. This allows us to discuss and patch the issue in private.
2. **Email:** Alternatively, contact the SVTech team at `security@svtech.software`.

We will acknowledge your report within **48 hours** and provide a detailed timeline for a patch.

---

## 3. Security Model

SV Printer follows these security principles:

- **Local binding only:** The agent binds to `127.0.0.1` by default and never to `0.0.0.0`, preventing remote access unless explicitly configured.
- **Token authentication:** All `/api/*` endpoints require `Authorization: Bearer <token>`. The token is generated on first run and stored in the config file.
- **CORS restriction:** Only origins explicitly listed in `-origin` / `allowed_origins` can call the API from a browser.
- **Ed25519 license verification:** Licenses are cryptographically signed and verified offline. The private signing key is never committed to the repository.
- **Device binding:** Licenses can be bound to a specific machine via fingerprint verification.

---

## 4. Secret Hygiene

- **Private keys:** The Ed25519 signing private key (`*.priv`, `.sv-license.priv`) is **gitignored** and must never be committed or shared.
- **Config file:** The auth token is stored in the JSON config file. Ensure the config directory has appropriate file permissions.
- **No telemetry:** SV Printer does not phone home or transmit any data externally. All license verification is offline.

---

## 5. Scope

This security policy covers the SV Printer agent binary, its HTTP API, and the license verification system. It does not cover:

- Third-party dependencies (report upstream).
- Physical printer security.
- Network configurations beyond the agent's binding.

---

## 6. Updates

Security fixes are released as patch versions (e.g., `0.1.1`) and announced in the [CHANGELOG](./CHANGELOG.md). We recommend always running the latest version.
