PRAGMA journal_mode=WAL;

CREATE TABLE IF NOT EXISTS goals (
  id INTEGER PRIMARY KEY,
  title TEXT NOT NULL,
  priority INTEGER DEFAULT 100,
  status TEXT CHECK(status IN ('active','paused','done')) DEFAULT 'active',
  notes TEXT,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS adrs (
  id TEXT PRIMARY KEY,          -- e.g. ADR-0001
  title TEXT NOT NULL,
  path TEXT NOT NULL,           -- filesystem path, e.g. ADR/0001-queue-invariants.md
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS ci_runs (
  id INTEGER PRIMARY KEY,
  scope TEXT,                   -- e.g. ./cmd/jukebox
  status TEXT CHECK(status IN ('pass','fail','error')) NOT NULL,
  started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  finished_at DATETIME
);

CREATE TABLE IF NOT EXISTS markdown_templates (
  id TEXT PRIMARY KEY,          -- e.g. README, ADR, CHANGELOG
  name TEXT NOT NULL,
  description TEXT,
  category TEXT,                -- e.g. documentation, project-management, development
  content TEXT NOT NULL,        -- template content with variables
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS template_variables (
  id INTEGER PRIMARY KEY,
  template_id TEXT NOT NULL,
  name TEXT NOT NULL,
  type TEXT CHECK(type IN ('string','date','list','number','boolean')) DEFAULT 'string',
  required BOOLEAN DEFAULT FALSE,
  default_value TEXT,
  description TEXT,
  FOREIGN KEY (template_id) REFERENCES markdown_templates(id) ON DELETE CASCADE
);
