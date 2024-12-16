BEGIN;
CREATE TYPE inventory_movement_types AS ENUM ('increase', 'decrease');
CREATE TYPE storehouse_types AS ENUM ('owned', 'not_owned');

CREATE TABLE item(
    id INT GENERATED ALWAYS AS IDENTITY,
    item_name VARCHAR(300) NOT NULL UNIQUE,
    description VARCHAR(5000),
    inventory_unit VARCHAR(10),
    image_link VARCHAR(500),
    count NUMERIC, -- allowed to be null in case it's a digital product
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE TABLE inventory_journal_entry(
    id INT GENERATED ALWAYS AS IDENTITY,
    description VARCHAR(5000),
    entry_date timestamptz NOT NULL DEFAULT NOW(),
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE TABLE category(
    id INT GENERATED ALWAYS AS IDENTITY,
    category_name VARCHAR(300) UNIQUE,
    description VARCHAR(5000),
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE TABLE item_category(
    item INT NOT NULL REFERENCES item(id) ON UPDATE CASCADE ON DELETE CASCADE,
    category INT NOT NULL REFERENCES category(id) ON UPDATE CASCADE ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    UNIQUE (item, category),
    PRIMARY KEY (item, category)
);

CREATE TABLE storehouse(
    id INT GENERATED ALWAYS AS IDENTITY,
    storehouse_name VARCHAR(300) NOT NULL UNIQUE,
    description VARCHAR(5000),
    storehouse_type storehouse_types NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE TABLE storehouse_item_quantity(
    quantity NUMERIC NOT NULL,
    storehouse INT NOT NULL REFERENCES storehouse(id) ON UPDATE CASCADE ON DELETE CASCADE,
    item INT NOT NULL REFERENCES item(id) ON UPDATE CASCADE ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    UNIQUE (item, storehouse),
    PRIMARY KEY (item, storehouse)
);

CREATE TABLE storehouse_group(
    id INT GENERATED ALWAYS AS IDENTITY,
    storehouse_group_name VARCHAR(300) NOT NULL UNIQUE,
    description VARCHAR(5000),
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE TABLE storehouse_storehouse_group(
    storehouse INT NOT NULL REFERENCES storehouse(id) ON UPDATE CASCADE ON DELETE CASCADE,
    storehouse_group INT NOT NULL REFERENCES storehouse_group(id) ON UPDATE CASCADE ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    UNIQUE (storehouse, storehouse_group),
    PRIMARY KEY (storehouse, storehouse_group)
);

CREATE TABLE inventory_movement(
    id INT GENERATED ALWAYS AS IDENTITY,
    storehouse INT NOT NULL REFERENCES storehouse(id) ON UPDATE CASCADE ON DELETE CASCADE,
    item INT NOT NULL REFERENCES item(id) ON UPDATE CASCADE ON DELETE CASCADE,
    quantity NUMERIC, -- allowed to be null in case it's a digital product
    inventory_movement_type inventory_movement_types NOT NULL,
    inventory_journal_entry INT REFERENCES inventory_journal_entry(id) ON UPDATE CASCADE ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE INDEX idx_inventory_movement_invenotry_journal_entry
ON inventory_movement(inventory_journal_entry);

CREATE INDEX idx_inventory_movement_storehouse
ON inventory_movement(storehouse);

CREATE INDEX idx_storehouse_storehouse_group_storehouse
ON storehouse_storehouse_group(storehouse);

CREATE INDEX idx_storehouse_storehouse_group_storehouse_group
ON storehouse_storehouse_group(storehouse_group);

CREATE INDEX idx_storehouse_item_quantity_storehouse
ON storehouse_item_quantity(storehouse);

CREATE INDEX idx_storehouse_item_quantity_item
ON storehouse_item_quantity(item);

COMMIT;