-- +goose Up
-- +goose StatementBegin
CREATE TABLE posts(
id UUID PRIMARY KEY,
title TEXT NOT NULL,
url TEXT UNIQUE NOT NULL,
description TEXT,
published_at TIMESTAMP NOT NULL,
created_at TIMESTAMP NOT NULL,
updated_at TIMESTAMP NOT NULL,
feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE posts;
-- +goose StatementEnd
