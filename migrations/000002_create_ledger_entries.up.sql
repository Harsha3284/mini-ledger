-- Ledger entries: atomic debit/credit records linked to an account
CREATE TABLE IF NOT EXISTS ledger_entries (
  id           UUID PRIMARY KEY,
  account_id   UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,

  direction    TEXT NOT NULL,
  amount       NUMERIC(14,2) NOT NULL,

  category     TEXT NULL,
  description  TEXT NULL,
  occurred_at  TIMESTAMPTZ NOT NULL,

  created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),

  CONSTRAINT ledger_direction_chk CHECK (direction IN ('debit','credit')),
  CONSTRAINT ledger_amount_chk CHECK (amount > 0)
);

CREATE INDEX IF NOT EXISTS idx_ledger_account_occurred ON ledger_entries (account_id, occurred_at DESC);
CREATE INDEX IF NOT EXISTS idx_ledger_occurred ON ledger_entries (occurred_at DESC);