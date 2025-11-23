-- Usar la base de datos creada por Docker (users_db)
USE users_db;

-- Tabla de usuarios
CREATE TABLE IF NOT EXISTS users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    email_verified BOOLEAN DEFAULT FALSE NOT NULL,
    email_verification_token VARCHAR(255),
    password_reset_token VARCHAR(255),
    password_reset_expires TIMESTAMP NULL,
    name VARCHAR(100) NOT NULL,
    lastname VARCHAR(100) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role ENUM('user', 'admin') DEFAULT 'user' NOT NULL,
    phone VARCHAR(20) NOT NULL,
    street VARCHAR(255) NOT NULL,
    number INT NOT NULL,
    photo_url VARCHAR(255),
    sex ENUM('hombre', 'mujer', 'otro') NOT NULL,
    avg_driver_rating DECIMAL(3,2) DEFAULT 0.00,
    avg_passenger_rating DECIMAL(3,2) DEFAULT 0.00,
    total_trips_passenger INT DEFAULT 0,
    total_trips_driver INT DEFAULT 0,
    birthdate DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_email (email),
    INDEX idx_email_verification_token (email_verification_token),
    INDEX idx_password_reset_token (password_reset_token)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabla de calificaciones
CREATE TABLE IF NOT EXISTS ratings (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    rater_id BIGINT NOT NULL,
    rated_user_id BIGINT NOT NULL,
    trip_id VARCHAR(24) NOT NULL,
    role_rated ENUM('conductor', 'pasajero') NOT NULL,
    score TINYINT NOT NULL CHECK (score >= 1 AND score <= 5),
    comment TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_rating (rater_id, trip_id, rated_user_id),
    INDEX idx_rater_id (rater_id),
    INDEX idx_rated_user_id (rated_user_id),
    INDEX idx_trip_id (trip_id),
    CONSTRAINT fk_ratings_rater FOREIGN KEY (rater_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_ratings_rated_user FOREIGN KEY (rated_user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Insertar usuario administrador inicial
-- Email: admin@carpooling.com
-- Password: admin123! (CAMBIAR EN PRODUCCIÓN)
-- Hash generado con bcrypt cost 10
INSERT INTO users (
    email,
    email_verified,
    name,
    lastname,
    password_hash,
    role,
    phone,
    street,
    number,
    sex,
    birthdate
) VALUES (
    'admin@admin.com',
    TRUE,
    'Administrador',
    'Sistema',
    '$2a$10$hS/B/Jctte1MfYFfFEySjOI4V8UhDYyqp189mWHrHtEeh2RSQA6Te',  -- Password: admin123!
    'admin',
    '0000000000',
    'N/A',
    0,
    'otro',
    '1990-01-01'
) ON DUPLICATE KEY UPDATE email=email;  -- Evita error si ya existe
