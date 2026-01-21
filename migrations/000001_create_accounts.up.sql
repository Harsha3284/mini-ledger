-- Accounts table for mini-ledger
CREATE TABLE IF NOT EXISTS accounts (
  id          UUID PRIMARY KEY,
  name        TEXT NOT NULL,
  type        TEXT NOT NULL,
  currency    TEXT NOT NULL,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),

  CONSTRAINT accounts_name_len CHECK (char_length(name) >= 1 AND char_length(name) <= 80),
  CONSTRAINT accounts_type_chk CHECK (type IN ('checking', 'savings', 'credit')),
  CONSTRAINT accounts_currency_chk CHECK (currency IN ('USD'))
);

CREATE INDEX IF NOT EXISTS idx_accounts_created_at ON accounts (created_at DESC);

-- Auto-update updated_at
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_accounts_set_updated_at ON accounts;
CREATE TRIGGER trg_accounts_set_updated_at
BEFORE UPDATE ON accounts
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();