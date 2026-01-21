DO $$
DECLARE conname text;
BEGIN
  SELECT c.conname INTO conname
  FROM pg_constraint c
  JOIN pg_class t ON t.oid = c.conrelid
  WHERE t.relname = 'ledger_entries'
    AND c.contype = 'f';

  IF conname IS NOT NULL THEN
    EXECUTE format('ALTER TABLE ledger_entries DROP CONSTRAINT %I', conname);
  END IF;
END $$;

ALTER TABLE ledger_entries
  ADD CONSTRAINT fk_ledger_entries_account_id
  FOREIGN KEY (account_id)
  REFERENCES accounts(id)
  ON DELETE RESTRICT;