BEGIN;
DROP TABLE IF EXISTS financial_group CASCADE;
DROP TABLE IF EXISTS financial_account CASCADE;
DROP TABLE IF EXISTS financial_journal_entry CASCADE;
DROP TABLE IF EXISTS financial_transaction CASCADE;

DROP TYPE IF EXISTS financial_transaction_types CASCADE;
DROP TYPE IF EXISTS financial_group_types CASCADE;
DROP TYPE IF EXISTS financial_account_types CASCADE;

DROP INDEX IF EXISTS idx_financial_transaction_journal_entry;
DROP INDEX IF EXISTS idx_financial_account_type;
DROP INDEX IF EXISTS idx_financial_account_group;
DROP INDEX IF EXISTS idx_financial_group_type;
DROP INDEX IF EXISTS idx_financial_journal_entry_updated_at;
DROP INDEX IF EXISTS idx_financial_transaction_updated_at;
COMMIT;