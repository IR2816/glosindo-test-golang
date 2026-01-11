-- Drop tables if exists
DROP TABLE IF EXISTS presensi CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- Create users table
CREATE TABLE users (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    address TEXT,
    division VARCHAR(50),
    position VARCHAR(50),
    join_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);

-- Create presensi table
CREATE TABLE presensi (
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL,
    date DATE NOT NULL,
    check_in_time TIMESTAMP,
    check_out_time TIMESTAMP,
    check_in_lat DECIMAL(10, 8),
    check_in_lng DECIMAL(11, 8),
    check_out_lat DECIMAL(10, 8),
    check_out_lng DECIMAL(11, 8),
    check_in_address TEXT,
    check_out_address TEXT,
    status VARCHAR(20) DEFAULT 'hadir',
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_presensi_user_id ON presensi(user_id);
CREATE INDEX idx_presensi_date ON presensi(date);
CREATE INDEX idx_presensi_user_date ON presensi(user_id, date);

-- Insert dummy user (password: 123456)
INSERT INTO users (id, name, email, password, phone, address, division, position, join_date)
VALUES (
    'USER001',
    'John Doe',
    'user@glosindo.com',
    '$2a$10$YourHashedPasswordHere',
    '081234567890',
    'Jl. Raya Bogor No. 123, Depok, Jawa Barat',
    'IT Department',
    'Software Engineer',
    '2023-01-15'
);

-- Note: Hash password "123456" dengan Go nanti