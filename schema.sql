CREATE TABLE IF NOT EXISTS code_execution_results (
    id SERIAL PRIMARY KEY,
    code TEXT NOT NULL,
    language TEXT NOT NULL,
    encoded_code TEXT NOT NULL,
    stdout TEXT,
    stderr TEXT,
    error TEXT,
    build_duration REAL,
    exec_duration REAL,
    created_at TIMESTAMPTZ DEFAULT NOW()
); 