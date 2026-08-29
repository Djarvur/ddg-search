# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Mandatory check after every code change

After **any** change to the code, run all three, in this order, and fix everything they report before considering the work done:

```bash
mise lint     # golangci-lint run --fix ./...
mise test     # go test -cover -race ./...
go build ./...
```

This mirrors the CI pipeline (`.github/workflows/ci.yml`: lint → test → build, build gated on the other two). Do not report a change as finished while any of the three fails. `mise lint` runs with `--fix`, so it may rewrite files — re-read anything you edited after running it.

## Commands

```bash
# Single test / subset
go test -run TestSearchIntegrationBasic ./internal/search/
go test -race -v ./internal/perplexity/

# Coverage the way CI measures it (threshold: total >= 50%, CI fails below)
go test -coverprofile=coverage.out -covermode=atomic -race ./...
go tool cover -func=coverage.out | grep total

# Build the three binaries individually
go build -o bin/ddg-search ./cmd/ddg-search
go build -o bin/page-dump ./cmd/page-dump
go build -o bin/perplexity-search ./cmd/perplexity-search

# Vulnerability scan (weekly in CI)
govulncheck ./...
```

Tooling versions come from `.mise.toml` — golangci-lint is pinned to `GOLANGCI_VERSION` (v2.9.0) and installed via `go install`, isolated from any system-wide golangci-lint. Keep that pin and the version in `.github/workflows/ci.yml` in sync.

Note: `README.md` says `make build`, but there is no Makefile — use the `go build` commands above.

### Tests that talk to the network

`internal/perplexity/integration_test.go` hits the real Perplexity API and self-skips unless **both** `PERPLEXITY_API_KEY` and `INTEGRATION_TEST=true` are set. `internal/search/integration_test.go`, despite the name, is hermetic: it serves `internal/search/testdata/*.html` fixtures from an `httptest` server. There are no build tags anywhere — plain `go test ./...` is always safe.

## Architecture

Three independent CLI binaries under `cmd/`, each a thin `urfave/cli/v3` wrapper (flag parsing, env lookup, `Fprintln(os.Stdout, …)`) over one `internal/` package. All logic lives in `internal/`; keep `main.go` files free of it. All three emit Markdown on stdout and debug logs on stderr, because the output is meant to be piped into an LLM.

| Binary | Package | Transport |
|---|---|---|
| `ddg-search` | `internal/search` (+ `internal/config`) | scrapes `html.duckduckgo.com/html` with goquery |
| `page-dump` | `internal/dump` | fetches a URL, converts via `html-to-markdown/v2` |
| `perplexity-search` | `internal/perplexity` | Perplexity chat API, `PERPLEXITY_API_KEY` required |

### The two retry stacks are deliberately separate

`internal/search` and `internal/perplexity` each define their **own** `RetryOptions` type, `DefaultRetryOptions()`, `ErrRateLimited`/`ErrMaxRetries`, `isRateLimited()`, and `calculateDelay()`. `internal/config` serves only the DuckDuckGo side. This duplication is intentional — the two services have different retryable-status semantics (DDG treats HTTP 202 as a rate-limit signal; Perplexity distinguishes 401/400/402 as non-retryable and breaks out of the loop). Do not unify them without a spec change, and when touching retry behaviour, check whether the change belongs to one service or both.

Jitter also differs on purpose: `search` uses `math/rand` (with a `//nolint:gosec`), `perplexity` uses `crypto/rand`.

### Retry nesting on the DuckDuckGo path (a real gotcha)

`Searcher.Search` ([search.go:41](internal/search/search.go#L41)) loops up to `MaxRetries+1` times, and each attempt calls `Client.Do` ([client.go:76](internal/search/client.go#L76)), which itself loops up to `MaxRetries+1` times. The retry counts multiply: with the default `--max-retries 3`, a persistently rate-limited query can issue up to 16 HTTP requests. `Client.Do` signals exhaustion with `ErrMaxRetries`, and the outer loop treats only that error as retryable — every other error returns immediately. Keep that contract if you touch either loop.

### Rate-limit detection in HTML is switched off

`Parser.FindRateLimitIndicator` ([parser.go:109](internal/search/parser.go#L109)) ignores its argument and always returns `""`, disabling the outer loop's HTML-body check. This is deliberate: the indicator words ("captcha", "blocked", "anomaly", …) produced false positives when they appeared in legitimate result snippets. `Parser.IsRateLimitPage` and `Parser.IsEmptyResults` still hold that logic but are not wired into the search flow. Re-enabling this needs a snippet-safe heuristic, not just un-stubbing the function.

DuckDuckGo results are read from `.result__a` anchors, and hrefs are unwrapped from DDG's `//duckduckgo.com/l/?uddg=…` redirect form by `Parser.extractURL`.

### Skills

`skills/ddg-search/` and `skills/perplexity-search/` package the binaries as Claude Code skills. The `skill.yaml` `command:` fields invoke the binaries **by bare name**, so they assume the binaries are on `PATH` (`go install github.com/Djarvur/ddg-search/cmd/<name>@latest`). Changing a CLI flag name or the stdout format is a breaking change for these skills — update `skill.yaml` and `SKILL.md` in the same change.

## Spec-driven workflow (openspec)

This repo uses `openspec` (`openspec/config.yaml`, `schema: spec-driven`). `openspec/specs/<capability>/spec.md` are the living requirements, written as `### Requirement:` blocks with `#### Scenario:` WHEN/THEN clauses; `openspec/changes/archive/<date>-<slug>/` holds completed change proposals (`proposal.md`, `design.md`, `tasks.md`, plus the spec deltas that were merged).

The project context in `openspec/config.yaml` states the working rules: **write tests first**, aim for 50% coverage, leave no linter errors, and pass `go build` / `mise lint` / `mise test` after every change. When changing behaviour, check the matching spec under `openspec/specs/` and update it alongside the code.

## Lint conventions

`.golangci.yml` enables **all** linters, then disables a short list with reasons inline. Practical consequences when writing code here:

- `funlen` caps functions at 80 lines / 50 statements and `gocognit` at complexity 15 — the existing code splits helpers out (`waitForRetry`, `searchAttempt`, `debugRateLimitReason`) specifically to stay under these. Prefer extracting a helper over adding a `//nolint`.
- `wsl_v5` is active: blank line required before `return` and around blocks. This is the most common lint failure; `mise lint --fix` handles it.
- Every exported identifier needs a doc comment, and every package needs a package comment.
- Errors are sentinel `var Err… = errors.New(…)` values wrapped with `%w`; never return bare `errors.New`/`fmt.Errorf` strings from a new code path where a sentinel would let callers use `errors.Is`.
- Existing `//nolint` directives all carry a justification comment — match that style if you must add one.
