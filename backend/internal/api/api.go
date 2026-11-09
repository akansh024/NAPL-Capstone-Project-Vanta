package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vanta/vanta/backend/internal/findings"
	"github.com/vanta/vanta/backend/internal/imports"
	"github.com/vanta/vanta/backend/internal/models"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type API struct {
	DB    *pgxpool.Pool
	Store findings.Store
}

func (a API) Router(origin string) http.Handler {
	r := chi.NewRouter()
	r.Use(cors(origin))
	r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) { respond(w, 200, map[string]string{"status": "ok"}) })
	r.Get("/api/ready", a.ready)
	r.Get("/api/dashboard", a.dashboard)
	r.Get("/api/findings", a.list)
	r.Get("/api/findings/{id}", a.get)
	r.Patch("/api/findings/{id}/status", a.status)
	r.Post("/api/imports/nmap", a.nmap)
	r.Get("/api/imports", a.importList)
	r.Get("/api/observations", a.observations)
	return r
}
func cors(origin string) func(http.Handler) http.Handler {
	return func(n http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PATCH,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if r.Method == "OPTIONS" {
				w.WriteHeader(204)
				return
			}
			n.ServeHTTP(w, r)
		})
	}
}
func respond(w http.ResponseWriter, c int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(c)
	json.NewEncoder(w).Encode(v)
}
func fail(w http.ResponseWriter, c int, msg string) { respond(w, c, map[string]string{"error": msg}) }
func (a API) ready(w http.ResponseWriter, r *http.Request) {
	if e := a.DB.Ping(r.Context()); e != nil {
		fail(w, 503, "database unavailable")
		return
	}
	respond(w, 200, map[string]string{"status": "ready"})
}
func (a API) list(w http.ResponseWriter, r *http.Request) {
	f, e := a.Store.List(r.Context(), r.URL.Query().Get("q"), r.URL.Query().Get("severity"), r.URL.Query().Get("status"))
	if e != nil {
		fail(w, 500, "could not retrieve findings")
		return
	}
	respond(w, 200, map[string]any{"items": f, "total": len(f)})
}
func (a API) get(w http.ResponseWriter, r *http.Request) {
	id, e := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if e != nil {
		fail(w, 400, "invalid finding id")
		return
	}
	f, e := a.Store.Get(r.Context(), id)
	if errors.Is(e, pgx.ErrNoRows) {
		fail(w, 404, "finding not found")
		return
	}
	if e != nil {
		fail(w, 500, "could not retrieve finding")
		return
	}
	respond(w, 200, f)
}
func (a API) status(w http.ResponseWriter, r *http.Request) {
	id, e := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if e != nil {
		fail(w, 400, "invalid finding id")
		return
	}
	var b struct {
		Status string `json:"status"`
	}
	if json.NewDecoder(http.MaxBytesReader(w, r.Body, 1024)).Decode(&b) != nil {
		fail(w, 400, "invalid JSON body")
		return
	}
	f, e := a.Store.UpdateStatus(r.Context(), id, b.Status)
	if e != nil {
		if errors.Is(e, pgx.ErrNoRows) {
			fail(w, 404, "finding not found")
		} else {
			fail(w, 400, e.Error())
		}
		return
	}
	respond(w, 200, f)
}
func (a API) dashboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var total, critical, high, open, resolved int
	e := a.DB.QueryRow(ctx, `SELECT count(*),count(*) FILTER(WHERE severity='Critical'),count(*) FILTER(WHERE severity='High'),count(*) FILTER(WHERE status='Open'),count(*) FILTER(WHERE status='Resolved') FROM findings`).Scan(&total, &critical, &high, &open, &resolved)
	if e != nil {
		fail(w, 500, "could not load dashboard")
		return
	}
	fs, e := a.Store.List(ctx, "", "", "")
	if e != nil {
		fail(w, 500, "could not load dashboard")
		return
	}
	if len(fs) > 6 {
		fs = fs[:6]
	}
	ims, _ := a.getImports(ctx)
	respond(w, 200, map[string]any{"metrics": map[string]int{"total": total, "critical": critical, "high": high, "open": open, "resolved": resolved}, "findings": fs, "imports": ims})
}
func (a API) nmap(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 5<<20)
	if e := r.ParseMultipartForm(5 << 20); e != nil {
		fail(w, 400, "file upload must be below 5 MB")
		return
	}
	file, h, e := r.FormFile("file")
	if e != nil {
		fail(w, 400, "provide a file field named file")
		return
	}
	defer file.Close()
	if !strings.HasSuffix(strings.ToLower(h.Filename), ".xml") {
		fail(w, 400, "only .xml Nmap reports are supported")
		return
	}
	report, e := imports.Parse(file)
	if e != nil {
		fail(w, 400, e.Error())
		return
	}
	ctx := r.Context()
	tx, e := a.DB.Begin(ctx)
	if e != nil {
		fail(w, 500, "could not start import")
		return
	}
	defer tx.Rollback(ctx)
	var iid int64
	if e = tx.QueryRow(ctx, "INSERT INTO report_imports(filename,scanner,status) VALUES($1,'Nmap','Completed') RETURNING id", h.Filename).Scan(&iid); e != nil {
		fail(w, 500, "could not save import")
		return
	}
	assets, obs := 0, 0
	seen := map[string]bool{}
	for _, host := range report.Hosts {
		addr := host.Address()
		if addr == "" {
			continue
		}
		var aid int64
		e = tx.QueryRow(ctx, `INSERT INTO assets(address,hostname,source) VALUES($1,$2,'Nmap import') ON CONFLICT(address) DO UPDATE SET hostname=COALESCE(EXCLUDED.hostname,assets.hostname),updated_at=now() RETURNING id`, addr, host.Hostname()).Scan(&aid)
		if e != nil {
			fail(w, 500, "could not save asset")
			return
		}
		if !seen[addr] {
			seen[addr] = true
			assets++
		}
		for _, p := range host.Ports {
			if p.State.State != "open" {
				continue
			}
			evidence := imports.ScriptEvidence(p)
			_, e = tx.Exec(ctx, "INSERT INTO scanner_observations(import_id,asset_id,port,protocol,service_name,state,script_id,evidence) VALUES($1,$2,$3,$4,$5,$6,$7,$8)", iid, aid, p.ID, p.Protocol, p.Service.Name, p.State.State, "", evidence)
			if e != nil {
				fail(w, 500, "could not save observation")
				return
			}
			obs++
		}
	}
	_, e = tx.Exec(ctx, "UPDATE report_imports SET asset_count=$1,observation_count=$2 WHERE id=$3", assets, obs, iid)
	if e != nil {
		fail(w, 500, "could not update import")
		return
	}
	if e = tx.Commit(ctx); e != nil {
		fail(w, 500, "could not finish import")
		return
	}
	respond(w, 201, map[string]any{"id": iid, "filename": h.Filename, "assets": assets, "observations": obs, "message": "Nmap observations imported. No vulnerability findings were created."})
}
func (a API) getImports(ctx context.Context) ([]models.Import, error) {
	rows, e := a.DB.Query(ctx, "SELECT id,filename,scanner,status,asset_count,observation_count,imported_at FROM report_imports ORDER BY imported_at DESC LIMIT 20")
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []models.Import{}
	for rows.Next() {
		var x models.Import
		if e = rows.Scan(&x.ID, &x.Filename, &x.Scanner, &x.Status, &x.AssetCount, &x.ObservationCount, &x.ImportedAt); e != nil {
			return nil, e
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
func (a API) importList(w http.ResponseWriter, r *http.Request) {
	x, e := a.getImports(r.Context())
	if e != nil {
		fail(w, 500, "could not retrieve imports")
		return
	}
	respond(w, 200, map[string]any{"items": x})
}

func (a API) observations(w http.ResponseWriter, r *http.Request) {
	rows, err := a.DB.Query(r.Context(), `SELECT o.id, a.address, COALESCE(a.hostname,''), o.port, o.protocol, COALESCE(o.service_name,''), o.state, COALESCE(o.evidence,''), i.filename, o.observed_at FROM scanner_observations o JOIN assets a ON a.id=o.asset_id JOIN report_imports i ON i.id=o.import_id ORDER BY o.observed_at DESC, o.id DESC LIMIT 100`)
	if err != nil {
		fail(w, 500, "could not retrieve observations")
		return
	}
	defer rows.Close()
	items := []models.Observation{}
	for rows.Next() {
		var item models.Observation
		if err := rows.Scan(&item.ID, &item.Asset, &item.Hostname, &item.Port, &item.Protocol, &item.Service, &item.State, &item.Evidence, &item.Filename, &item.ObservedAt); err != nil {
			fail(w, 500, "could not read observations")
			return
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		fail(w, 500, "could not retrieve observations")
		return
	}
	respond(w, 200, map[string]any{"items": items})
}

var _ = fmt.Sprint
var _ = time.Now
