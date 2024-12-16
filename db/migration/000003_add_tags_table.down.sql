BEGIN;
DROP TABLE IF EXISTS tag CASCADE;
DROP TABLE IF EXISTS tag_financial_journal_entry CASCADE;
DROP TABLE IF EXISTS tag_inventory_journal_entry CASCADE;

DROP INDEX IF EXISTS idx_tag_financial_journal_entry;
DROP INDEX IF EXISTS idx_financial_journal_entry_tag;
DROP INDEX IF EXISTS idx_tag_inventory_journal_entry;
DROP INDEX IF EXISTS idx_inventory_journal_entry_tag;
COMMIT;