SET statement_timeout = 0;

--bun:split

CREATE SCHEMA IF NOT EXISTS rbac;

--bun:split

CREATE TABLE IF NOT EXISTS rbac.roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    pretty_name VARCHAR(255) NOT NULL,
    description TEXT,
    permissions BIGINT[] DEFAULT '{}'::BIGINT[],
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
