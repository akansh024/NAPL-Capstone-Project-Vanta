CREATE TABLE IF NOT EXISTS assets (
  id BIGSERIAL PRIMARY KEY, address TEXT NOT NULL UNIQUE, hostname TEXT, source TEXT NOT NULL DEFAULT 'imported', created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS findings (
  id BIGSERIAL PRIMARY KEY, title TEXT NOT NULL, description TEXT, severity TEXT NOT NULL CHECK (severity IN ('Critical','High','Medium','Low','Info')), status TEXT NOT NULL DEFAULT 'Open' CHECK (status IN ('Open','In Progress','Resolved')), asset_id BIGINT REFERENCES assets(id), evidence TEXT, remediation TEXT, source TEXT NOT NULL DEFAULT 'development seed', last_seen_at TIMESTAMPTZ NOT NULL DEFAULT now(), created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS report_imports (
  id BIGSERIAL PRIMARY KEY, filename TEXT NOT NULL, scanner TEXT NOT NULL, status TEXT NOT NULL, asset_count INT NOT NULL DEFAULT 0, observation_count INT NOT NULL DEFAULT 0, error_message TEXT, imported_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS scanner_observations (
  id BIGSERIAL PRIMARY KEY, import_id BIGINT NOT NULL REFERENCES report_imports(id) ON DELETE CASCADE, asset_id BIGINT NOT NULL REFERENCES assets(id), port INT NOT NULL, protocol TEXT NOT NULL, service_name TEXT, state TEXT NOT NULL, script_id TEXT, evidence TEXT, observed_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS remediation_history (
  id BIGSERIAL PRIMARY KEY, finding_id BIGINT NOT NULL REFERENCES findings(id) ON DELETE CASCADE, old_status TEXT, new_status TEXT NOT NULL, changed_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
INSERT INTO assets(address, hostname, source) VALUES ('192.0.2.10','demo-web-01','development seed'),('192.0.2.20','demo-db-01','development seed') ON CONFLICT(address) DO NOTHING;
INSERT INTO findings(title,description,severity,status,asset_id,evidence,remediation,source,last_seen_at)
SELECT 'Outdated TLS configuration','The development sample asset allows a legacy TLS configuration. This illustrative record is not a scan result.','High','Open',id,'Development seed only','Disable obsolete TLS protocols and retest.','development seed',now() FROM assets WHERE address='192.0.2.10' AND NOT EXISTS (SELECT 1 FROM findings WHERE title='Outdated TLS configuration');
INSERT INTO findings(title,description,severity,status,asset_id,evidence,remediation,source,last_seen_at)
SELECT 'Default database administrator account','Illustrative seed finding for workflow demonstration.','Critical','In Progress',id,'Development seed only','Rotate credentials, disable default account, and document ownership.','development seed',now() FROM assets WHERE address='192.0.2.20' AND NOT EXISTS (SELECT 1 FROM findings WHERE title='Default database administrator account');
INSERT INTO findings(title,description,severity,status,asset_id,evidence,remediation,source,last_seen_at)
SELECT 'Missing security headers','Illustrative seed finding for presentation only.','Medium','Resolved',id,'Development seed only','Set appropriate security headers at the web server.','development seed',now() FROM assets WHERE address='192.0.2.10' AND NOT EXISTS (SELECT 1 FROM findings WHERE title='Missing security headers');
