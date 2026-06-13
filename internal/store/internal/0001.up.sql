CREATE TABLE IF NOT EXISTS branches (
	name TEXT PRIMARY KEY,
	description TEXT NOT NULL,
	parent_ref TEXT NOT NULL,
	parent_sha BLOB NOT NULL,
	created_at INTEGER NOT NULL,
	updated_at INTEGER NOT NULL,
	last_accessed_at INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS aliases (
	alias TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	FOREIGN KEY (name) REFERENCES branches (name)
		ON DELETE CASCADE
);