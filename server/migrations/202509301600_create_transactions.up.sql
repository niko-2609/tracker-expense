CREATE TABLE transactions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    amount NUMERIC(12,2) NOT NULL,
    txn_type VARCHAR(10) NOT NULL CHECK (txn_type IN ('income','expense')),
    frequency VARCHAR(20) NOT NULL CHECK (frequency IN ('daily','weekly','monthly','quarterly','yearly')),
    category_id BIGINT REFERENCES categories(id),
    txn_date DATE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);
