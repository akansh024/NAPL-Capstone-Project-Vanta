# Architecture

The browser serves the Vite-built React application. It communicates with the Go REST API using JSON over HTTP. The API owns validation, Nmap parsing, SQL access, and migrations; PostgreSQL is the persistent source of truth.

```text
React UI → Go API → PostgreSQL volume
               ↑
         Nmap XML upload
```

On startup, the API reads the required `DATABASE_URL`, opens and pings a pgx pool, applies ordered SQL migrations, and starts HTTP on `PORT`. The listener address is `:<PORT>`, which binds all interfaces for container hosts such as Render. Migrations create deterministic development seed findings once. The API restricts CORS to `FRONTEND_ORIGIN`. On `SIGTERM` or `SIGINT`, the API stops accepting new requests and closes its PostgreSQL pool.

Nmap uploads are limited to 5 MB and `.xml`, parsed with Go's strict XML decoder. Hosts become assets and open ports/services/scripts become scanner observations associated with an import. An observation is not a vulnerability finding. Docker Compose supplies PostgreSQL health checks, an API dependency, and persistent `postgres_data` storage for local development. Production supplies `DATABASE_URL` from the deployment environment; no production credential is stored in this repository.
