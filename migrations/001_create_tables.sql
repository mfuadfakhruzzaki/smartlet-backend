-- PostgreSQL database schema
-- Create database: CREATE DATABASE swiflet_db;

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    location VARCHAR(255),
    no_telp VARCHAR(20),
    password VARCHAR(255) NOT NULL,
    img_profile VARCHAR(255),
    status INTEGER DEFAULT 0,
    role INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tags table
CREATE TABLE IF NOT EXISTS tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

-- Articles table
CREATE TABLE IF NOT EXISTS articles (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    cover_image VARCHAR(255),
    tag_id INTEGER REFERENCES tags(id) ON DELETE SET NULL,
    status INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Comments table
CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY,
    article_id INTEGER NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- EBooks table
CREATE TABLE IF NOT EXISTS ebooks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    file_path VARCHAR(255) NOT NULL,
    thumbnail_path VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Videos table
CREATE TABLE IF NOT EXISTS videos (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    youtube_link VARCHAR(255) NOT NULL,
    thumbnail_path VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- SwifletHouse table
CREATE TABLE IF NOT EXISTS swiflet_houses (
    id SERIAL PRIMARY KEY,
    id_user INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    location VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- IoTDevice table
CREATE TABLE IF NOT EXISTS iot_devices (
    id SERIAL PRIMARY KEY,
    id_swiflet_house INTEGER NOT NULL REFERENCES swiflet_houses(id) ON DELETE CASCADE,
    floor INTEGER NOT NULL,
    install_code VARCHAR(255) NOT NULL UNIQUE,
    status INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- WeeklyPrice table
CREATE TABLE IF NOT EXISTS weekly_prices (
    id SERIAL PRIMARY KEY,
    province VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    week_start DATE NOT NULL,
    week_end DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Harvest table
CREATE TABLE IF NOT EXISTS harvests (
    id SERIAL PRIMARY KEY,
    id_user INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    id_swiflet_house INTEGER NOT NULL REFERENCES swiflet_houses(id) ON DELETE CASCADE,
    floor INTEGER NOT NULL,
    bowl_weight DECIMAL(10,2) DEFAULT 0,
    bowl_pieces INTEGER DEFAULT 0,
    oval_weight DECIMAL(10,2) DEFAULT 0,
    oval_pieces INTEGER DEFAULT 0,
    corner_weight DECIMAL(10,2) DEFAULT 0,
    corner_pieces INTEGER DEFAULT 0,
    broken_weight DECIMAL(10,2) DEFAULT 0,
    broken_pieces INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- HarvestSales table
CREATE TABLE IF NOT EXISTS harvest_sales (
    id SERIAL PRIMARY KEY,
    id_user INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    province VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    bowl_weight DECIMAL(10,2) DEFAULT 0,
    oval_weight DECIMAL(10,2) DEFAULT 0,
    corner_weight DECIMAL(10,2) DEFAULT 0,
    broken_weight DECIMAL(10,2) DEFAULT 0,
    appointment_date DATE NOT NULL,
    proof_photo VARCHAR(255),
    status INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- InstallationRequest table
CREATE TABLE IF NOT EXISTS installation_requests (
    id SERIAL PRIMARY KEY,
    id_swiflet_house INTEGER NOT NULL REFERENCES swiflet_houses(id) ON DELETE CASCADE,
    floors VARCHAR(255) NOT NULL,
    sensor_count INTEGER NOT NULL,
    appointment_date DATE NOT NULL,
    status INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- MaintenanceRequest table
CREATE TABLE IF NOT EXISTS maintenance_requests (
    id SERIAL PRIMARY KEY,
    id_device INTEGER NOT NULL REFERENCES iot_devices(id) ON DELETE CASCADE,
    reason TEXT NOT NULL,
    appointment_date DATE NOT NULL,
    status INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- UninstallationRequest table
CREATE TABLE IF NOT EXISTS uninstallation_requests (
    id SERIAL PRIMARY KEY,
    id_device INTEGER NOT NULL REFERENCES iot_devices(id) ON DELETE CASCADE,
    reason TEXT NOT NULL,
    appointment_date DATE NOT NULL,
    status INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Transaction table
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    order_id VARCHAR(255) NOT NULL UNIQUE,
    status INTEGER NOT NULL DEFAULT 0,
    amount DECIMAL(10,2) NOT NULL,
    payment_type VARCHAR(100) NOT NULL,
    transaction_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Membership table
CREATE TABLE IF NOT EXISTS memberships (
    id SERIAL PRIMARY KEY,
    id_user INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    join_date DATE NOT NULL,
    exp_date DATE NOT NULL,
    order_id VARCHAR(255) NOT NULL REFERENCES transactions(order_id),
    status INTEGER DEFAULT 0
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_articles_tag_id ON articles(tag_id);
CREATE INDEX IF NOT EXISTS idx_comments_article_id ON comments(article_id);
CREATE INDEX IF NOT EXISTS idx_comments_user_id ON comments(user_id);
CREATE INDEX IF NOT EXISTS idx_swiflet_houses_user_id ON swiflet_houses(id_user);
CREATE INDEX IF NOT EXISTS idx_iot_devices_swiflet_house_id ON iot_devices(id_swiflet_house);
CREATE INDEX IF NOT EXISTS idx_iot_devices_install_code ON iot_devices(install_code);
CREATE INDEX IF NOT EXISTS idx_harvests_user_id ON harvests(id_user);
CREATE INDEX IF NOT EXISTS idx_harvests_swiflet_house_id ON harvests(id_swiflet_house);
CREATE INDEX IF NOT EXISTS idx_harvest_sales_user_id ON harvest_sales(id_user);
CREATE INDEX IF NOT EXISTS idx_weekly_prices_province ON weekly_prices(province);
CREATE INDEX IF NOT EXISTS idx_weekly_prices_week_start ON weekly_prices(week_start);
CREATE INDEX IF NOT EXISTS idx_memberships_user_id ON memberships(id_user);