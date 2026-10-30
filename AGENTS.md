# AGENTS.md

## Repo context

- Boot.dev "Linko" starter for the Logging & Telemetry course: a deliberately small,
  messy Go URL shortener used to practice adding logs, metrics, and traces.
- Do not "clean up" intentional inefficiencies unless asked: bcrypt hashing on every
  redirect, an unbounded in-memory `redirects` slice, and live outbound HTTP checks
  in `destination.go` are there for the exercises.

## Commands

- Go 1.27 (see `go.mod`), module `boot.dev/linko`.
- Build: `go build -o linko .`
- Run: `go run .` (or `./linko`); defaults to port `8899` and data dir `./data`.
  Flags: `-port`, `-data`.
- Verify: `go vet ./...` and `go fmt ./...`; `go test ./...` reports no test files
  (this repo has no tests, CI, or lint config).
- `index.html` is embedded with `go:embed`; rebuild/restart after editing it.

## Layout and runtime quirks

- Single `main` package at the root plus `internal/store`. Routes are wired in
  `server.go:newServer`; add new handlers there.
- Uses Go 1.22+ `ServeMux` method patterns. `GET /` is a catch-all for unmatched GET
  paths, and `GET /{shortCode}` handles redirects.
- The store is file-per-short-code under `-data`. `Store.Lookup` uppercases codes;
  `Store.List` returns at most 10 entries.
- All `/api/*` routes require Basic auth. Seeded users: `frodo` / `ofTheNineFingers`
  and `samwise` / `theStrong`. `saruman` has an invalid hash on purpose and returns
  500, which is useful for exercising error logs.
- `POST /admin/shutdown` gracefully stops a manually launched server (404 when
  `ENV=production`).
- Shorten and redirect requests perform a live outbound GET, so manual testing needs
  network access and a reachable destination URL.
- Logging uses the stdlib `log` package (stderr). Keep existing log statements unless
  the task is specifically about changing telemetry.
- `data/`, the `linko` binary, and `*.log` (e.g. `linko.out.log`) are gitignored.
  Never commit runtime data or logs.

## Git

- `origin` is `https://github.com/iamadarsh0309/Linko-Observability.git`, the
  published project repo. Do not push or rewrite history unless asked.
