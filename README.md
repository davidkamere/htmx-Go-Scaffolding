# HTMX Go Scaffolding (v2)

Production-leaning starter for server-rendered Go apps with HTMX partial updates and Tailwind CSS.

## Included by default
- Real HTMX CRUD flow (create/delete tasks with server-rendered partials)
- `cmd` + `internal` project layout
- Environment-based config (`PORT`, `APP_ENV`, `LOG_LEVEL`, `DB_PATH`)
- Persistent file-backed storage by default (`./data/tasks.json`)
- Middleware stack: security headers, panic recovery, request logging
- Graceful shutdown on `SIGINT`/`SIGTERM`
- Go tests and CI workflow
- Dockerfile + docker-compose
- Local vendored HTMX runtime (served from `/assets/vendor/htmx.min.js`)

## Project layout
- `cmd/server/main.go`: binary entrypoint
- `internal/app`: app startup and graceful shutdown
- `internal/server`: router and dependency wiring
- `internal/handlers`: HTTP handlers and template rendering
- `internal/tasks`: persistent task repository
- `internal/config`: env config loading
- `internal/middleware`: common middleware
- `internal/templates`: template discovery/parser
- `web/templates`: HTML templates and HTMX partials
- `web/assets`: Tailwind input/output and vendor JS

## Requirements
- Go 1.20+
- Node.js 18+

## Quick start
```bash
cp .env.example .env
npm install
npm run build:css
npm run build:vendor
go run ./cmd/server
```

Open: [http://localhost:8080](http://localhost:8080)

## Common commands
```bash
make deps       # npm install
make css        # build css + vendor htmx
make run        # run server
make test       # go test ./...
make vet        # go vet ./...
make fmt        # format Go files
make ci         # CI-equivalent local check
```

## Dev workflow
Terminal 1:
```bash
npm run watch:css
```

Terminal 2:
```bash
air
```

## Docker
```bash
docker compose up --build
```

Open: [http://localhost:8080](http://localhost:8080)

Data persists in `./data/tasks.json` via the mounted volume.

## Notes
- This scaffold intentionally keeps rendering on the server side.
- The data file and parent directory are auto-created at startup.
