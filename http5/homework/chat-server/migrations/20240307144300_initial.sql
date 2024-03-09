-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id INT PRIMARY KEY
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE private_messages (
    id INT PRIMARY KEY
    receiver VARCHAR(255) NOT NULL,
    sender VARCHAR(255) NOT NULL,
    message TEXT
);

CREATE TABLE public_messages (
    id INT PRIMARY KEY
    sender VARCHAR(255) NOT NULL,
    message TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE private_messages;
DROP TABLE public_messages;
DROP TABLE users;
-- +goose StatementEnd
