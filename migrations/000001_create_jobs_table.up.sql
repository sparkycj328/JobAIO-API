CREATE TABLE IF NOT EXISTS jobs (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    vendor text NOT NULL,
    country text NOT NULL,
    amount integer NOT NULL,
    url text NOT NULL,
    version integer NOT NULL DEFAULT 1
);

ALTER TABLE jobs ADD CONSTRAINT jobs_amount_check CHECK (amount >=0);

CREATE INDEX IF NOT EXISTS jobs_vendor_idx ON jobs USING GIN (to_tsvector('simple', vendor));
