CREATE TABLE IF NOT EXISTS jobs (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    vendor text NOT NULL,
    country text NOT NULL,
    amount integer NOT NULL,
    url text NOT NULL,
    );