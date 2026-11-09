# API

Base URL: `http://localhost:8080/api`. Errors use `{"error":"message"}`.

| Method | Route | Purpose |
|---|---|---|
| GET | `/health` | Liveness response |
| GET | `/ready` | PostgreSQL readiness |
| GET | `/dashboard` | Metrics, recent findings/imports |
| GET | `/findings?q=&severity=&status=` | Backend-filtered findings |
| GET | `/findings/{id}` | Finding details |
| PATCH | `/findings/{id}/status` | Persist a workflow status |
| POST | `/imports/nmap` | Multipart Nmap XML (`file`) |
| GET | `/imports` | Import history |
| GET | `/observations` | Persisted Nmap asset/service observations and evidence |

Status request:

```json
{"status":"In Progress"}
```

Allowed values are `Open`, `In Progress`, and `Resolved`. Invalid IDs, malformed JSON, invalid statuses, report formats, or XML yield safe `400`/`404` responses. A successful Nmap import returns an ID and persisted asset/observation counts, e.g. `{"id":1,"assets":2,"observations":3}`.
