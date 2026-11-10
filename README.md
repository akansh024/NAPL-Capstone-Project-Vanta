# VANTA
### Vulnerability Analysis, Normalization, Tracking and Assessment

**An Academic Vulnerability Management Platform**

**Project Team**

- **Akansh Pandey** | B.Sc. (Hons.) Cybersecurity
- **Aditya Chouhan** | B.Sc. (Hons.) Cybersecurity
- **Kartikey Drall** | B.Sc. (Hons.) Computer Science

---

## Overview

VANTA (Vulnerability Analysis, Normalization, Tracking and Assessment) is an academic vulnerability management platform designed to centralize, organize, and track security findings.

The platform imports authorized Nmap XML reports as scanner observations and provides a persistent workflow for reviewing and managing tracked findings.

VANTA does not directly scan targets, fabricate CVEs, exploit systems, or automatically remediate vulnerabilities.

## Working Capabilities

- **Dashboard:** React-based dashboard powered by a Go backend and PostgreSQL database.
- **Finding Management:** Search, filter, view, and update vulnerability finding records.
- **Status Tracking:** Manage findings through three workflow states: `Open`, `In Progress`, and `Resolved`.
- **Nmap XML Import:** Upload authorized Nmap XML reports with file-size and format validation.
- **Observation Persistence:** Store imported host, service, and script observations.
- **Import History:** Review previous report imports and their associated statistics.
- **Development Seed Data:** Clearly labelled, deterministic sample findings for demonstrations and presentations.

## Technology Stack

| Component | Technology |
|---|---|
| Frontend | React, TypeScript, Vite |
| Backend | Go |
| Database | PostgreSQL |
| Containerization | Docker, Docker Compose |

## Quick Start

### Prerequisites

- Docker
- Docker Compose

### Run the Application

Clone the repository and navigate to the project directory.

```bash
cp .env.example .env
docker compose up --build
```

Once the containers are running, access the application:

| Service | URL |
|---|---|
| Frontend Dashboard | http://localhost:5173 |
| Backend Health API | http://localhost:8080/api/health |

Docker Compose starts the frontend, backend, and PostgreSQL services.

The API automatically applies the database schema and initializes development seed records on startup. PostgreSQL data is persisted using the `postgres_data` volume.

### Demonstrating Nmap XML Import

To demonstrate the report-import functionality:

1. Open the VANTA dashboard.
2. Navigate to **Report Imports**.
3. Select `sample-data/nmap-sample.xml`.
4. Click **Import Nmap XML**.
5. Review the generated import history.

The sample report should produce:

- **2 assets**
- **3 observations**

These observations are intentionally stored separately and are not automatically converted into vulnerability findings.

## Local Development

The recommended local development environment uses Docker Compose.

Copy the example environment file:

```bash
cp .env.example .env
```

Start the application:

```bash
docker compose up --build
```

### Running the Backend Without Docker Compose

For standalone backend development, configure `DATABASE_URL` with a valid local PostgreSQL connection string.

Run the following command from the `backend` directory:

```bash
go run ./cmd/api
```

The backend expects the following environment variables:

| Variable | Description | Default |
|---|---|---|
| `DATABASE_URL` | PostgreSQL connection URL | Required |
| `FRONTEND_ORIGIN` | Frontend origin for CORS configuration | `http://localhost:5173` |

The frontend supports:

| Variable | Description | Default |
|---|---|---|
| `VITE_API_URL` | Backend API base URL | `http://localhost:8080/api` |

## Production Deployment (Render)

VANTA supports deployment using Render PostgreSQL and a backend container image.

### Backend Deployment

Deploy the backend using the Dockerfile located at:

```text
backend/Dockerfile
```

The Docker image includes the SQL migrations directory:

```text
/app/migrations
```

The API automatically applies ordered SQL migrations before accepting HTTP requests.

### Required Environment Variables

Configure the following variables in the Render service dashboard:

| Variable | Description |
|---|---|
| `DATABASE_URL` | PostgreSQL connection URL provided by Render |
| `FRONTEND_ORIGIN` | Exact public origin of the deployed frontend |
| `PORT` | Port assigned by the application host |

**Important:** Never commit production credentials or database connection strings to source control.

### Database Initialization

The Go application requires `DATABASE_URL`. There is no fallback to a local database connection string.

During startup, the API:

1. Establishes a PostgreSQL connection pool.
2. Connects to and verifies the database.
3. Applies the ordered, idempotent SQL migrations.
4. Starts accepting HTTP requests.

If the database connection or migration process fails, the API does not start serving requests.

The application also handles `SIGTERM` and `SIGINT` signals to shut down the HTTP server and database connection pool gracefully.

**Production consideration:** The current migrations include deterministic development seed findings. Review and adjust these records before deploying to a new production database.

## Architecture and API Documentation

For additional technical details, refer to:

- [System Architecture](docs/ARCHITECTURE.md)
- [API Reference](docs/API.md)
- [Project Roadmap](ROADMAP.md)

## Project Status and Limitations

VANTA currently represents a presentation-ready first development phase.

The following capabilities are not yet implemented:

- User authentication and authorization.
- OWASP ZAP report importer.
- Vulnerability reporting and export functionality.
- Automatic finding correlation and deduplication.
- Risk scoring and prioritization.
- Complete user and audit management workflows.

Nmap observations and confirmed vulnerability findings are intentionally maintained as separate records.

For planned improvements and future development, see [ROADMAP.md](ROADMAP.md).

## Troubleshooting

| Issue | Solution |
|---|---|
| Port `5173` is already occupied | Change `FRONTEND_PORT` in `.env`. |
| Port `8080` is already occupied | Change `API_PORT` in `.env`. |
| Invalid Nmap XML report | Ensure the uploaded file contains valid XML in a supported Nmap format. |
| Upload rejected | The server accepts `.xml` files up to 5 MB. |
| Changes are not reflected | Rebuild the application using `docker compose up --build`. |
| Database connection failure | Verify PostgreSQL availability and the configured `DATABASE_URL`. |

---

**VANTA | Vulnerability Analysis, Normalization, Tracking and Assessment**

*Academic Project | B.Sc. (Hons.) Cybersecurity & Computer Science*
