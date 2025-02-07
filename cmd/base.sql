CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date DATE,
    title VARCHAR(128) NOT NULL DEFAULT '',
    comment VARCHAR(256) NOT NULL DEFAULT '',
    repeat VARCHAR(128) NOT NULL DEFAULT ''
);

CREATE INDEX date_sort ON scheduler (date);