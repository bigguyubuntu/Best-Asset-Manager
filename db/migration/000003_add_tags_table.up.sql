BEGIN;
CREATE TABLE tag(
    id INT GENERATED ALWAYS AS IDENTITY,
    tag_name VARCHAR(1000) NOT NULL UNIQUE,
    description VARCHAR(10000),
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE TABLE tag_financial_journal_entry(
    tag INT NOT NULL REFERENCES tag(id) ON UPDATE CASCADE ON DELETE CASCADE,
    financial_journal_entry INT NOT NULL REFERENCES financial_journal_entry(id) ON UPDATE CASCADE ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    UNIQUE (tag, financial_journal_entry),
    PRIMARY KEY (tag, financial_journal_entry)
);

CREATE TABLE tag_inventory_journal_entry(
    tag INT NOT NULL REFERENCES tag(id) ON UPDATE CASCADE ON DELETE CASCADE,
    inventory_journal_entry INT NOT NULL REFERENCES inventory_journal_entry(id) ON UPDATE CASCADE ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    UNIQUE (tag, inventory_journal_entry),
    PRIMARY KEY (tag, inventory_journal_entry)
);

CREATE INDEX idx_tag_financial_journal_entry
ON tag_financial_journal_entry(tag);
CREATE INDEX idx_financial_journal_entry_tag
ON tag_financial_journal_entry(financial_journal_entry);

CREATE INDEX idx_tag_inventory_journal_entry
ON tag_inventory_journal_entry(tag);
CREATE INDEX idx_inventory_journal_entry_tag
ON tag_inventory_journal_entry(inventory_journal_entry);

COMMIT;