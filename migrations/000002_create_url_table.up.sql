CREATE TABLE url (
    short_url VARCHAR(7) PRIMARY KEY,
    addr TEXT NOT NULL,
    redirect_count INTEGER NOT NULL,
    creation_date TIMESTAMPTZ NOT NULL
);