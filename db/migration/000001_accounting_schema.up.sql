BEGIN;

CREATE TYPE financial_transaction_types AS ENUM ('debit', 'credit');

CREATE TYPE financial_group_types AS ENUM ('assets', 'liabilities', 'equity', 'expenses', 'income' );

CREATE TYPE financial_account_types AS ENUM ('debit_increased', 'credit_increased');

CREATE TABLE financial_group(
    id INT GENERATED ALWAYS AS IDENTITY,
    parent_id INT NOT NULL,
    group_name VARCHAR(1000) NOT NULL,
    description VARCHAR(10000),
    financial_group_type financial_group_types NOT NULL,
    balance BIGINT DEFAULT 0,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id),
    UNIQUE(group_name, parent_id)
);

CREATE TABLE financial_account(
    id INT GENERATED ALWAYS AS IDENTITY,
    account_name VARCHAR(1000) NOT NULL,
    description VARCHAR(10000),
    balance BIGINT DEFAULT 0,
    financial_account_type financial_account_types NOT NULL,
    financial_group INT NOT NULL REFERENCES financial_group(id) ON UPDATE CASCADE ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id),
    UNIQUE(account_name, financial_group)
);

CREATE TABLE financial_journal_entry(
    id INT GENERATED ALWAYS AS IDENTITY,
    entry_date timestamptz NOT NULL DEFAULT NOW(),
    description VARCHAR(5000),
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE TABLE financial_transaction(
    id INT GENERATED ALWAYS AS IDENTITY,
    amount BIGINT DEFAULT 0,
    financial_transaction_type financial_transaction_types NOT NULL,
    financial_journal_entry INT NOT NULL REFERENCES financial_journal_entry(id) ON UPDATE CASCADE ON DELETE CASCADE,
    financial_account INT NOT NULL REFERENCES financial_account(id) ON UPDATE CASCADE ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE INDEX idx_financial_transaction_journal_entry
ON financial_transaction(financial_journal_entry);

CREATE INDEX idx_financial_account_type
ON financial_account(financial_account_type);

CREATE INDEX idx_financial_account_group
ON financial_account(financial_group);

CREATE INDEX idx_financial_group_type
ON financial_group(financial_group_type);

CREATE INDEX idx_financial_transaction_updated_at
ON financial_transaction(updated_at);

CREATE INDEX idx_financial_journal_entry_updated_at
ON financial_journal_entry(updated_at);
