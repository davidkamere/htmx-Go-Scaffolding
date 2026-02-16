# Go + HTMX + Tailwind Scaffolding

Minimal starter for server-rendered Go apps with HTMX and Tailwind CSS.

## What this scaffold includes
- Go HTTP server with `gorilla/mux`
- HTML templates in `web/templates/*.gohtmx`
- Static asset serving from `web/assets`
- Tailwind CSS build/watch scripts
- Optional live-reload with Air

## Requirements
- Go 1.20+
- Node.js 18+

## Quick start
1. Install frontend dependencies:
   ```bash
   npm install
   ```
2. Build Tailwind CSS once:
   ```bash
   npm run build:css
   ```
3. Run the Go app:
   ```bash
   go run .
   ```
4. Open [http://localhost:8080](http://localhost:8080)

## Development workflow
Run these in separate terminals:

1. Watch and rebuild CSS on template changes:
   ```bash
   npm run watch:css
   ```
2. Run Go live reload (requires Air):
   ```bash
   air
   ```

If you do not use Air, run:
```bash
go run .
```

## Structure
- `main.go`: router, handlers, template rendering
- `web/templates/base.gohtmx`: base layout + HTMX script + CSS include
- `web/templates/index.gohtmx`: sample page
- `web/assets/input.css`: Tailwind input file
- `web/assets/build.css`: generated stylesheet

## Notes
- HTMX is loaded from CDN in `base.gohtmx`.
- Generated CSS is served at `/assets/build.css`.
