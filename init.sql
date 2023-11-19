-- Create tables
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    password VARCHAR(255),
    email VARCHAR(255) UNIQUE
);

CREATE TABLE properties (
    id SERIAL PRIMARY KEY,
    address VARCHAR(255),
    contact_name VARCHAR(255),
    contact_phone VARCHAR(20),
    contact_email VARCHAR(255),
    tags TEXT,
    rating DECIMAL(2,1),
    type VARCHAR(50)
);

CREATE TABLE rooms (
    id SERIAL PRIMARY KEY,
    room_nr INT,
    hotel INT REFERENCES properties(id),
    description TEXT,
    count INT,
    price DECIMAL(10,2),
    availability_interval tsrange
);

-- Seed users
-- INSERT INTO users (name, password, email) VALUES 
-- ('u', 'pass', 'u@example.com'),

-- Seed properties (2 hotels, 1 house, 2 flats)
INSERT INTO properties (address, contact_name, contact_phone, contact_email, tags, rating, type) VALUES 
('101 Budapest Street, Budapest', 'John Doe', '123456789', 'john@example.com', '#dogfriendly, #nonsmoking', 8.5, 'hotel'),
('102 Budapest Street, Budapest', 'Jane Smith', '987654321', 'jane@example.com', '#family, #nonsmoking', 9.0, 'hotel'),
('103 Budapest Street, Budapest', 'Alice Johnson', '111222333', 'alice@example.com', '#luxury, #nonsmoking', 9.5, 'house'),
('104 Budapest Street, Budapest', 'Bob Brown', '444555666', 'bob@example.com', '#budget, #petfriendly', 7.5, 'flat'),
('105 Budapest Street, Budapest', 'Charlie Davis', '777888999', 'charlie@example.com', '#central, #modern', 8.0, 'flat');

-- Seed 7 rooms for these hotels
INSERT INTO rooms (room_nr, hotel, description, count, price, availability_interval) VALUES 
(101, 1, 'Single room with city view', 2, 50.00, '[2023-01-01, 2023-12-31]'),
(102, 1, 'Double room with garden view', 3, 70.00, '[2023-01-01, 2023-12-31]'),
(103, 1, 'Deluxe room with river view', 1, 90.00, '[2023-01-01, 2023-12-31]'),
(201, 2, 'Suite with balcony', 2, 120.00, '[2023-01-01, 2023-12-31]'),
(202, 3, 'Family room with kitchen', 3, 80.00, '[2023-01-01, 2023-12-31]'),
(203, 4, 'Executive suite with office space', 1, 150.00, '[2023-01-01, 2023-12-31]'),
(204, 5, 'Standard single room', 2, 60.00, '[2023-01-01, 2023-12-31]');