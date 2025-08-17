CREATE TABLE IF NOT EXISTS managed_filepaths (
	file_id INTEGER PRIMARY KEY AUTOINCREMENT,
	filepath TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS tags_definition (
	tag_id INTEGER PRIMARY KEY AUTOINCREMENT,
	tag_name TEXT UNIQUE NOT NULL,
	tag_description TEXT NULL
);

CREATE TABLE IF NOT EXISTS file_tags (
	file_id INTEGER REFERENCES managed_filepaths(file_id),
	tag_id INTEGER REFERENCES tags_definition(tag_id)
);
