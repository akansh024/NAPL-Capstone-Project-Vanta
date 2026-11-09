package findings

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vanta/vanta/backend/internal/models"
)

type Store struct{ DB *pgxpool.Pool }

const cols = `f.id,f.title,coalesce(f.description,''),f.severity,f.status,coalesce(a.address,''),coalesce(f.evidence,''),coalesce(f.remediation,''),f.source,f.last_seen_at`

func scan(row interface{ Scan(...any) error }) (models.Finding, error) {
	var f models.Finding
	err := row.Scan(&f.ID, &f.Title, &f.Description, &f.Severity, &f.Status, &f.Asset, &f.Evidence, &f.Remediation, &f.Source, &f.LastSeen)
	return f, err
}
func (s Store) List(ctx context.Context, q, severity, status string) ([]models.Finding, error) {
	sql := `SELECT ` + cols + ` FROM findings f LEFT JOIN assets a ON a.id=f.asset_id WHERE ($1='' OR f.title ILIKE '%'||$1||'%' OR a.address ILIKE '%'||$1||'%') AND ($2='' OR f.severity=$2) AND ($3='' OR f.status=$3) ORDER BY f.last_seen_at DESC`
	rows, e := s.DB.Query(ctx, sql, q, severity, status)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []models.Finding{}
	for rows.Next() {
		f, e := scan(rows)
		if e != nil {
			return nil, e
		}
		out = append(out, f)
	}
	return out, rows.Err()
}
func (s Store) Get(ctx context.Context, id int64) (models.Finding, error) {
	return scan(s.DB.QueryRow(ctx, `SELECT `+cols+` FROM findings f LEFT JOIN assets a ON a.id=f.asset_id WHERE f.id=$1`, id))
}
func (s Store) UpdateStatus(ctx context.Context, id int64, status string) (models.Finding, error) {
	if status != "Open" && status != "In Progress" && status != "Resolved" {
		return models.Finding{}, fmt.Errorf("status must be Open, In Progress, or Resolved")
	}
	tx, e := s.DB.Begin(ctx)
	if e != nil {
		return models.Finding{}, e
	}
	defer tx.Rollback(ctx)
	var old string
	if e = tx.QueryRow(ctx, "SELECT status FROM findings WHERE id=$1", id).Scan(&old); e != nil {
		return models.Finding{}, e
	}
	tag, e := tx.Exec(ctx, "UPDATE findings SET status=$1,updated_at=now() WHERE id=$2", status, id)
	if e != nil || tag.RowsAffected() == 0 {
		return models.Finding{}, fmt.Errorf("finding not found")
	}
	_, e = tx.Exec(ctx, "INSERT INTO remediation_history(finding_id,old_status,new_status) VALUES($1,$2,$3)", id, old, status)
	if e != nil {
		return models.Finding{}, e
	}
	if e = tx.Commit(ctx); e != nil {
		return models.Finding{}, e
	}
	return s.Get(ctx, id)
}
