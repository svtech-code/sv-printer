# Contributing to SV Print

Thank you for your interest in contributing to **SV Print**! This guide outlines the workflow and standards.

> **Just want to install SV Print?** See the [Integration Guide](./documentation/integration.md) for build instructions.

---

## 1. Branching Strategy & Workflow

We follow a simple feature-branch workflow:

1. **Branch Naming:** Create a branch from `main` using descriptive prefixes:
   - `feat/feature-name` for new features.
   - `fix/bug-name` for bug fixes.
   - `documentation/doc-name` for documentation updates.
   - `refactor/refactor-name` for code cleanup.
2. **Pull Requests:** Once your changes are ready and tested:
   - Push your branch to GitHub.
   - Open a Pull Request (PR) against `main`.
   - Ensure all tests pass before requesting review.

> **Note for external contributors:** All changes must come through a pull request. Continuous Integration (CI) runs `gofmt`, `go vet`, `go test`, and a build check on Linux, macOS, and Windows, and **must pass green** before your PR can be merged.

---

## 2. Commit Message Standards

This repository strictly enforces the **[Conventional Commits](https://www.conventionalcommits.org/)** specification. **All commit messages must be written in English.**

Format: `<type>(<scope>): <description>`

- **Allowed Types:** `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `build`, `ci`, `chore`.
- **Scope:** Optional but recommended (e.g., `api`, `escpos`, `config`, `licensing`).
- **Examples:**
  - `feat(api): add receipt structured endpoint`
  - `fix(ws): subscribe before Accept to eliminate EventsHandler race`
  - `docs(readme): redesign with hero, badges and Mermaid`

---

## 3. Development Environment

Ensure Go 1.26+ is installed (see [`go.mod`](./go.mod)).

### Key Commands

```bash
# Build
go build ./cmd/sv-print

# Run unit tests
go test ./...

# Format codebase
gofmt -l .

# Static analysis
go vet ./...
```

---

## 4. Pre-Commit Gate (Required)

Before presenting any commit message, you **must** run the full barrier and fix any failures:

```bash
gofmt -l .
go vet ./...
go test ./...
```

If any of these fail, fix the issues first, re-run, and only then present the commit message. This prevents wasted CI cycles and broken main branches.

---

## 5. Project Conventions

- **TDD:** Write tests after each task and keep them passing.
- **Standard library first:** Avoid external packages when the Go stdlib suffices.
- **Spec-driven:** Behavior and architecture changes go through the `sv-memory` spec flow before implementation.
- **No git automation:** Never run `git add`, `git commit`, or `git push` autonomously — the developer handles these manually.
- **Memory protocol:** Use `sv-memory` tools (search, save, graph) during the workflow.

---

## 6. License

SV Print is licensed under the [Business Source License 1.1](./LICENSE) (BSL 1.1). By contributing, you agree that your contributions will be licensed under the same terms.
