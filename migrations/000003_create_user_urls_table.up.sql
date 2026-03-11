CREATE TABLE user_urls (
    user_id INTEGER NOT NULL REFERENCES users(id),
    short_url VARCHAR(7) REFERENCES url(short_url),
    PRIMARY KEY (user_id, short_url)
);