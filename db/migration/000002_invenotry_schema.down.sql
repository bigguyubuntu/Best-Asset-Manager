BEGIN;

DROP TABLE IF EXISTS item CASCADE;
DROP TABLE IF EXISTS category CASCADE;
DROP TABLE IF EXISTS item_category CASCADE;
DROP TABLE IF EXISTS storehouse CASCADE;
DROP TABLE IF EXISTS storehouse_item_quantity CASCADE;
DROP TABLE IF EXISTS storehouse_group CASCADE;
DROP TABLE IF EXISTS storehouse_storehouse_group CASCADE;
DROP TABLE IF EXISTS inventory_movement CASCADE;
DROP TABLE IF EXISTS inventory_journal_entry CASCADE;

DROP TYPE IF EXISTS inventory_movement_types CASCADE;
DROP TYPE IF EXISTS storehouse_types CASCADE;

DROP INDEX IF EXISTS idx_inventory_movement_invenotry_journal_entry;
DROP INDEX IF EXISTS idx_inventory_movement_storehouse;
DROP INDEX IF EXISTS idx_storehouse_storehouse_group_storehouse;
DROP INDEX IF EXISTS idx_storehouse_storehouse_group_storehouse_group;
DROP INDEX IF EXISTS idx_storehouse_item_quantity_storehouse;
DROP INDEX IF EXISTS idx_storehouse_item_quantity_item;

COMMIT;
