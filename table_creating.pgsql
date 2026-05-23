CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20),
    nickname VARCHAR(255),
    city VARCHAR(100),
    language VARCHAR(10) DEFAULT 'ru',
    current_theme VARCHAR(10) DEFAULT 'light',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS cars (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    is_main BOOLEAN DEFAULT false,
    brand VARCHAR(50) NOT NULL,
    model VARCHAR(50) NOT NULL,
    year INTEGER,
    color VARCHAR(30),
    body_type VARCHAR(30),
    engine_type VARCHAR(30),
    engine_volume NUMERIC(3,1),
    transmission VARCHAR(30),
    mileage INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS expenses (
    id SERIAL PRIMARY KEY,
    car_id INTEGER NOT NULL REFERENCES cars(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    amount NUMERIC(12,2) NOT NULL,
    date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS diagnostics (
    id SERIAL PRIMARY KEY,
    car_id INTEGER NOT NULL REFERENCES cars(id) ON DELETE CASCADE,
    date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    error_codes TEXT,
    severity VARCHAR(20),
    ai_response_structured JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS ai_history (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    car_id INTEGER REFERENCES cars(id) ON DELETE SET NULL,
    message TEXT NOT NULL,
    response TEXT NOT NULL,
    response_structured JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE maintenance_logs (
    id SERIAL PRIMARY KEY,
    car_id INT NOT NULL,
    work_date DATE NOT NULL DEFAULT CURRENT_DATE,
    description TEXT NOT NULL,
    mileage INT, -- Пробег на момент ремонта/ТО
    cost NUMERIC(10, 2) DEFAULT 0.00,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Связываем лог с машиной. При удалении машины удалится и история
    CONSTRAINT fk_car FOREIGN KEY (car_id) REFERENCES cars(id) ON DELETE CASCADE
);

-- Создаем индекс, чтобы Go-бэкенд моментально доставал историю для конкретной машины
CREATE INDEX idx_maintenance_logs_car_id ON maintenance_logs(car_id);