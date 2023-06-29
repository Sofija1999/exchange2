CREATE TABLE IF NOT EXISTS egw.user
(
    id         UUID DEFAULT uuid_generate_v4() PRIMARY KEY,

    email      VARCHAR(255) NOT NULL UNIQUE,

    first_name  VARCHAR(255) NOT NULL,
    surname     VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,

    updated_at TIMESTAMP    NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP    NOT NULL DEFAULT NOW()
);