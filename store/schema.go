package store

const schemaSQL = `
CREATE TABLE IF NOT EXISTS meta(
	key TEXT PRIMARY KEY,
	value TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS envelopes(
	id TEXT PRIMARY KEY,
	sender TEXT NOT NULL,
	recipient TEXT NOT NULL,
	subject TEXT NOT NULL,
	body TEXT NOT NULL,
	status TEXT NOT NULL,
	box TEXT NOT NULL,
	tries INTEGER NOT NULL,
	defer_to TEXT NOT NULL DEFAULT '',
	reason TEXT NOT NULL DEFAULT '',
	created TEXT NOT NULL,
	updated TEXT NOT NULL
);
`
