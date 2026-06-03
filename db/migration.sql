CREATE TABLE IF NOT EXISTS loans (
    merchant_id VARCHAR(50) PRIMARY KEY,
    shop_name VARCHAR(255) NOT NULL,
    principal_amount NUMERIC(15, 2) NOT NULL,
    current_balance NUMERIC(15, 2) NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);