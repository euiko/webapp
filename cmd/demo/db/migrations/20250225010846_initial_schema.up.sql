SET statement_timeout = 0;

--bun:split
CREATE TABLE companies (
    id serial PRIMARY KEY,
    name text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz
)

--bun:split

CREATE TABLE employees (
    id serial PRIMARY KEY,
    name text NOT NULL,
    company_id integer NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz
)
