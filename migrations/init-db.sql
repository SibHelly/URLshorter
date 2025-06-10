-- Таблица для хранения URL
CREATE TABLE IF NOT EXISTS urls (
    id SERIAL PRIMARY KEY,
    original_url TEXT NOT NULL,
    alias VARCHAR(50) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN DEFAULT TRUE NOT NULL,
    visit_count INT NOT NULL,
    title VARCHAR(255),
    description TEXT
);

-- Таблица для хранения статистики посещений
CREATE TABLE IF NOT EXISTS visits (
    id SERIAL PRIMARY KEY,
    url_id INT REFERENCES urls(id) ON DELETE CASCADE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Индексы для оптимизации
CREATE INDEX IF NOT EXISTS idx_urls_alias ON urls(alias);
CREATE INDEX IF NOT EXISTS idx_urls_expires_at ON urls(expires_at);
CREATE INDEX IF NOT EXISTS idx_visits_url_id ON visits(url_id);
CREATE INDEX IF NOT EXISTS idx_visits_created_at ON visits(created_at);