CREATE TABLE IF NOT EXISTS cryptocurrencies (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(10) NOT NULL UNIQUE,      
    name VARCHAR(50) NOT NULL UNIQUE          
);

CREATE TABLE IF NOT EXISTS trackings (
    id SERIAL PRIMARY KEY,
    cryptocurrency_id INT NOT NULL UNIQUE REFERENCES cryptocurrencies(id) ON DELETE CASCADE,
    is_active BOOLEAN DEFAULT TRUE,
);

CREATE TABLE IF NOT EXISTS price_history (
    id BIGSERIAL PRIMARY KEY,
    cryptocurrency_id INT NOT NULL REFERENCES cryptocurrencies(id) ON DELETE CASCADE,
    price NUMERIC(20, 8) NOT NULL,
    timestamp TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_price_history ON price_history (cryptocurrency_id, timestamp DESC);