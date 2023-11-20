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

CREATE TABLE reserv (
    id SERIAL PRIMARY KEY,
    userId INT REFERENCES users(id),
    propertyId INT REFERENCES properties(id),
    roomId INT REFERENCES rooms(id),
    cost DECIMAL(10,2),
    reserv_interval tsrange,
    creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
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
('105 Budapest Street, Budapest', 'Charlie Davis', '777888999', 'charlie@example.com', '#central, #modern', 8.0, 'flat'),
('106 Budapest Street, Budapest', 'Eve Green', '111333444', 'eve@example.com', '#historic, #luxury', 9.0, 'hotel'),
('107 Budapest Street, Budapest', 'Sam Blue', '222444555', 'sam@example.com', '#contemporary, #citycenter', 8.0, 'flat'),
('108 Budapest Street, Budapest', 'Lily White', '333555666', 'lily@example.com', '#budget, #central', 7.0, 'hostel'),
('109 Budapest Street, Budapest', 'Max Black', '444666777', 'max@example.com', '#family, #nonsmoking', 8.5, 'hotel'),
('110 Budapest Street, Budapest', 'Nora Grey', '555777888', 'nora@example.com', '#petfriendly, #luxury', 9.5, 'villa');

-- Seed 7 rooms for these hotels
INSERT INTO rooms (room_nr, hotel, description, count, price, availability_interval) VALUES 
(101, 1, 'Single room with city view', 2, 50.00, '[2023-01-01, 2023-12-31]'),
(102, 1, 'Double room with garden view', 3, 70.00, '[2023-01-01, 2023-12-31]'),
(103, 1, 'Deluxe room with river view', 1, 90.00, '[2023-01-01, 2023-12-31]'),
(201, 2, 'Suite with balcony', 2, 120.00, '[2023-01-01, 2023-12-31]'),
(202, 3, 'Family room with kitchen', 3, 80.00, '[2023-01-01, 2023-12-31]'),
(203, 4, 'Executive suite with office space', 1, 150.00, '[2023-01-01, 2023-12-31]'),
(204, 5, 'Standard single room', 2, 60.00, '[2023-01-01, 2023-12-31]'),
(301, 6, 'Penthouse with cityscape view', 1, 150.00, '[2023-01-01, 2023-12-31]'),
(302, 6, 'Twin room with two beds', 2, 80.00, '[2023-01-01, 2023-12-31]'),
(401, 7, 'Studio apartment', 3, 100.00, '[2023-01-01, 2023-12-31]'),
(402, 8, 'Shared room with 4 beds', 4, 30.00, '[2023-01-01, 2023-12-31]'),
(501, 9, 'Queen suite with jacuzzi', 2, 120.00, '[2023-01-01, 2023-12-31]');