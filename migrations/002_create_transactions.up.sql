CREATE TABLE IF NOT EXISTS transactions (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID         NOT NULL REFERENCES users(id),
    amount      NUMERIC(15,2) NOT NULL CHECK (amount > 0),
    type        VARCHAR(10)  NOT NULL CHECK (type IN ('income', 'expense')),
    category    VARCHAR(100) NOT NULL,
    date        DATE         NOT NULL,
    notes       TEXT,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ  NULL
);

CREATE INDEX IF NOT EXISTS idx_txn_user_id   ON transactions (user_id)  WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_txn_date      ON transactions (date)     WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_txn_category  ON transactions (category) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_txn_type      ON transactions (type)     WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_txn_user_date ON transactions (user_id, date) WHERE deleted_at IS NULL;
