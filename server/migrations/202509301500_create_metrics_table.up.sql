CREATE TABLE user_dashboard_metrics (
    user_id BIGINT PRIMARY KEY,
    total_income NUMERIC(12,2) DEFAULT 0,
    total_expense NUMERIC(12,2) DEFAULT 0,
    net_savings NUMERIC(12,2) DEFAULT 0,
    monthly_totals JSONB DEFAULT '{}',       -- e.g., {"2025-01": 1200, "2025-02": 250}
    top_expense_categories JSONB DEFAULT '[]', -- e.g., [{"category":"Food","amount":600}]
    updated_at TIMESTAMP DEFAULT NOW()
);
