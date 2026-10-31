# AGENTS.md

## Repo context

- Open-source "Linko" URL shortener used as a realistic target for logging,
  metrics, tracing, and other observability work.
- Preserve deliberate rough edges unless asked: bcrypt hashing on every redirect,
  an unbounded in-memory `redirects` slice, and live outbound HTTP checks in
  `destination.go` create useful observability signals.

## Commands

- Go 1.27 (see `go.mod`), module `boot.dev/linko`.
- Build: `go build -o linko .`
- Run: `go run .` (or `./linko`); defaults to port `8899` and data dir `./data`.
  Flags: `-port`, `-data`.
- Verify: `go vet ./...`, `go fmt ./...`, `go test ./...` (no test files yet), and
  `golangci-lint run` (config in `.golangci.yml`; only `forbidigo` is enabled).
  Use a golangci-lint binary built with Go >= 1.27.
- CI: `.github/workflows/ci.yml` runs formatting, vet, tests, build, and lint on
  pushes to `main` and pull requests.
- `index.html` is embedded from `internal/server` with `go:embed`; rebuild/restart
  after editing it.

## Layout and runtime quirks

- Root `main.go` is the entrypoint. Application code lives under `internal/`:
  `server` (HTTP routing, handlers, middleware), `auth` (Basic auth), `destination`
  (outbound URL checks), `logging` (logger construction), and `store` (file storage).
  Routes are wired in `internal/server.New`; add new handlers there.
- Uses Go 1.22+ `ServeMux` method patterns. `GET /` is a catch-all for unmatched GET
  paths, and `GET /{shortCode}` handles redirects.
- `requestLogger(logger)` in `internal/server/middleware.go` wraps the mux and logs
  `Served request: METHOD PATH` for every request.
- The store is file-per-short-code under `-data`. `Store.Lookup` uppercases codes;
  `Store.List` returns at most 10 entries.
- All `/api/*` routes require Basic auth. Seeded users: `frodo` / `ofTheNineFingers`
  and `samwise` / `theStrong`. `saruman` has an invalid hash on purpose and returns
  500, which is useful for exercising error logs.
- `POST /admin/shutdown` gracefully stops a manually launched server (404 when
  `ENV=production`).
- Shorten and redirect requests perform a live outbound GET, so manual testing needs
  network access and a reachable destination URL.
- Logging is configured in `internal/logging` by `New`, which writes to stderr and,
  when `LINKO_LOG_FILE` is set, also appends to that file via `io.MultiWriter`.
  Loggers are injected into `server.New` and `store.New`.
- `fmt.Print*` is banned by `forbidigo` ("Use logger instead."). Keep existing log
  statements unless the task is specifically about changing telemetry.
- `data/`, the `linko` binary, and `*.log` (e.g. `linko.out.log`) are gitignored.
  Never commit runtime data or logs.

## Git

- `origin` is `https://github.com/iamadarsh0309/Linko-Observability.git`, the
  published project repo. Do not push or rewrite history unless asked.
