# VANTA

VANTA (Vulnerability Analysis, Normalization, Tracking and Assessment) is an early academic vulnerability-management platform. It imports authorized Nmap XML reports as scanner observations and provides a persisted workflow for reviewing tracked findings. It does **not** scan targets, fabricate CVEs, exploit systems, or automatically remediate issues.

## Working capabilities

- React dashboard backed by Go and PostgreSQL.
- Search/filter finding records; open detail; update `Open`, `In Progress`, or `Resolved` status.
- Nmap XML upload with safe size/format checks, host/service/script observation persistence, and import history.
- Clearly labelled deterministic development seed findings for a presentation.

## Quick start

```sh
cp .env.example .env
docker compose up --build
```

Open http://localhost:5173. The API is at http://localhost:8080/api/health. Docker creates the schema and seeds development records on API startup; PostgreSQL persists in `postgres_data`.

To demonstrate import, open **Report imports**, select `sample-data/nmap-sample.xml`, and click **Import Nmap XML**. The history should show 2 assets and 3 observations. They are deliberately not converted into findings.

## Local development

For the supported local setup, copy `.env.example` to the ignored `.env` file and start the stack with Docker Compose. Compose starts its own PostgreSQL container and supplies the API with a Docker-network `DATABASE_URL`; it does not use a production database URL.

```sh
cp .env.example .env
docker compose up --build
```

For a non-Compose local backend, set `DATABASE_URL` to a local PostgreSQL connection URL, then run `go run ./cmd/api` from `backend`. The API expects `FRONTEND_ORIGIN` (defaults to `http://localhost:5173`). The frontend supports `VITE_API_URL` (defaults to `http://localhost:8080/api`).

## Render PostgreSQL production

Deploy the backend image from `backend/Dockerfile` to the application host. The image includes `/app/migrations`, and the API applies the ordered SQL migrations before it accepts HTTP requests.

Set these backend environment variables in the deployment dashboard, never in source control:

- `DATABASE_URL`: the Render PostgreSQL connection URL supplied by Render.
- `FRONTEND_ORIGIN`: the exact public origin of the deployed frontend.
- `PORT`: the port supplied by the application host. The API listens on all interfaces at that port.

`DATABASE_URL` is required. There is no local connection-string fallback in the Go application. On startup, the API opens and pings the PostgreSQL pool, runs the idempotent migrations in `backend/migrations`, then begins serving requests. A failed connection or migration prevents startup. The API closes the HTTP server and pool on `SIGTERM`/`SIGINT`.

The migration currently creates schema and deterministic development seed findings. Review whether those seed records are appropriate before using a new production database.

## Architecture and API

See [architecture](docs/ARCHITECTURE.md) and [API reference](docs/API.md). Stack: React/TypeScript/Vite, Go, PostgreSQL, and Docker Compose.

## Status and limitations

This is a presentation-ready first phase. It has no authentication, ZAP importer, reporting exports, automatic finding correlation, risk scoring, or complete user/audit workflow. Nmap observations are intentionally separate from confirmed findings. See [ROADMAP.md](ROADMAP.md).

## Troubleshooting

- If port 5173/8080 is occupied, change `FRONTEND_PORT`/`API_PORT` in `.env`.
- A malformed XML report returns a safe 400 response. The server accepts `.xml` uploads up to 5 MB.
- Rebuild after code changes with `docker compose up --build`.
