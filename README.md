# Linko

A small, open-source URL shortener built as a realistic target for logging,
metrics, tracing, and other observability work.

Linko is intentionally small enough to understand in one sitting, but it has the
pieces that make real services interesting: HTTP routing, Basic auth, file-backed
persistence, outbound network calls, in-memory state, and a few deliberate
performance traps. Use it to instrument a working service instead of a
hello-world.

> Linko is a lab, not a production URL shortener.

## What you can practice

- Structured logging and request-scoped context
- Request metrics: latency, throughput, status codes, and error rates
- Trace spans across handlers, storage, and outbound calls
- Dashboards, alerts, SLOs, and load testing
- Profiling and tracking down slow paths

Linko currently logs lifecycle events and every HTTP request. It does not emit
metrics or traces yet; adding them is the point of the project.

## Features

- Create short links and redirect visitors
- List recent links and view in-memory redirect stats
- Basic auth on API routes
- File-per-short-code storage
- Request logging with optional file output
- Graceful shutdown

## Quick start

Requires Go 1.27 or newer.

```sh
git clone https://github.com/iamadarsh0309/Linko-Observability.git
cd Linko-Observability
go build -o linko .
./linko
```

Open <http://localhost:8899>. By default, Linko listens on port `8899` and stores
short codes under `./data`.

To also write logs to a file:

```sh
LINKO_LOG_FILE=linko.access.log ./linko
```

File logging uses an 8 KiB buffered writer; the buffer is flushed when the
process exits.

### Configuration

| Flag / variable | Default | Purpose |
| --- | --- | --- |
| `-port` | `8899` | HTTP port |
| `-data` | `./data` | Directory for short-code files |
| `LINKO_LOG_FILE` | unset | Append logs to this file as well as stderr |
| `ENV` | unset | Set to `production` to disable `POST /admin/shutdown` |

## API

| Method | Path | Auth | Description |
| --- | --- | --- | --- |
| `GET` | `/` | No | Web UI |
| `POST` | `/api/login` | Yes | Validate credentials |
| `POST` | `/api/shorten` | Yes | Create a short link from form field `url` |
| `GET` | `/api/urls` | Yes | List recent links as JSON |
| `GET` | `/api/stats` | Yes | Redirect stats since the last restart |
| `GET` | `/{shortCode}` | No | Redirect to the original URL |
| `POST` | `/admin/shutdown` | No | Gracefully stop the server (development only) |

Demo users:

- `frodo` / `ofTheNineFingers`
- `samwise` / `theStrong`

The `saruman` user has an intentionally invalid password hash, which makes it
useful for exercising authentication error logs.

Create a short link:

```sh
curl -u frodo:ofTheNineFingers \
  -X POST \
  -d 'url=https://example.com' \
  http://localhost:8899/api/shorten
```

The response is the six-character short code. Visit
`http://localhost:8899/<code>` to follow the redirect.

## Development

```sh
go test ./...
go vet ./...
golangci-lint run
```

The CI workflow runs formatting checks, vet, tests, the build, and
`golangci-lint` on pushes to `main` and pull requests. Linting is configured in
`.golangci.yml`; `fmt.Print*` is forbidden because application logging should go
through the configured logger.

### Project layout

```text
main.go                    Entrypoint, flags, signal handling
internal/server/           HTTP server, routes, handlers, request logging
internal/auth/             Basic auth middleware and user context
internal/destination/      Outbound URL validation
internal/logging/          Configurable logger construction
internal/store/            File-backed short-code storage
```

## Design notes

Some of Linko's rough edges are intentional and useful for observability work:

- The redirect path performs a live outbound HTTP check and generates a bcrypt
  hash on every request.
- Redirect history is held in an unbounded in-memory slice, so stats reset when
  the process restarts.
- The link list is capped at 10 entries.
- The store performs blocking file I/O directly in request handlers.

Measure before changing these. They exist to create interesting latency, error,
and resource-usage signals to investigate.

Linko is not production-ready: credentials are hardcoded, there is no TLS or
rate limiting, and the stats store is in memory.

## License

MIT. See [LICENSE](LICENSE).
