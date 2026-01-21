ALTER TABLE ledger_entries DROP CONSTRAINT IF EXISTS fk_ledger_entries_account_id;

ALTER TABLE ledger_entries
  ADD CONSTRAINT fk_ledger_entries_account_id
  FOREIGN KEY (account_id)
  REFERENCES accounts(id)
  ON DELETE CASCADE;