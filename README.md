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

Run PostgreSQL separately, set `DATABASE_URL='postgres://vanta:password@localhost:5432/vanta?sslmode=disable'`, then:

```sh
cd backend && go run ./cmd/api
cd frontend && npm install && npm run dev
```

The API expects `FRONTEND_ORIGIN` (defaults to `http://localhost:5173`). The frontend supports `VITE_API_URL` (defaults to `http://localhost:8080/api`).

## Architecture and API

See [architecture](docs/ARCHITECTURE.md) and [API reference](docs/API.md). Stack: React/TypeScript/Vite, Go, PostgreSQL, and Docker Compose.

## Status and limitations

This is a presentation-ready first phase. It has no authentication, ZAP importer, reporting exports, automatic finding correlation, risk scoring, or complete user/audit workflow. Nmap observations are intentionally separate from confirmed findings. See [ROADMAP.md](ROADMAP.md).

## Troubleshooting

- If port 5173/8080 is occupied, change `FRONTEND_PORT`/`API_PORT` in `.env`.
- A malformed XML report returns a safe 400 response. The server accepts `.xml` uploads up to 5 MB.
- Rebuild after code changes with `docker compose up --build`.
