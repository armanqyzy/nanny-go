CREATE DATABASE nanny_db;

CREATE TABLE users (
                       user_id SERIAL PRIMARY KEY,
                       full_name VARCHAR(100),
                       email VARCHAR(100) UNIQUE,
                       phone VARCHAR(20),
                       password_hash VARCHAR(255),
                       role VARCHAR(10) CHECK (role IN ('owner', 'sitter', 'admin')),
                       created_at TIMESTAMP DEFAULT NOW()
);


CREATE TABLE pets (
                      pet_id SERIAL PRIMARY KEY,
                      owner_id INT REFERENCES users(user_id),
                      name VARCHAR(50),
                      type VARCHAR(20) CHECK (type IN ('cat', 'dog', 'rodent')),
                      age INT,
                      notes TEXT
);


CREATE TABLE sitters (
                         sitter_id INT PRIMARY KEY REFERENCES users(user_id),
                         experience_years INT,
                         certificates TEXT,
                         preferences TEXT,
                         location VARCHAR(100),
                         status VARCHAR(10) CHECK (status IN ('pending', 'approved', 'rejected'))
);


CREATE TABLE services (
                          service_id SERIAL PRIMARY KEY,
                          sitter_id INT REFERENCES sitters(sitter_id),
                          type VARCHAR(20) CHECK (type IN ('walking', 'boarding', 'home-care')),
                          price_per_hour DECIMAL(10,2),
                          description TEXT
);


CREATE TABLE bookings (
                          booking_id SERIAL PRIMARY KEY,
                          owner_id INT REFERENCES users(user_id),
                          sitter_id INT REFERENCES sitters(sitter_id),
                          pet_id INT REFERENCES pets(pet_id),
                          service_id INT REFERENCES services(service_id),
                          start_time TIMESTAMP,
                          end_time TIMESTAMP,
                          status VARCHAR(15) CHECK (status IN ('pending', 'confirmed', 'cancelled', 'completed'))
);


CREATE TABLE payments (
                          payment_id SERIAL PRIMARY KEY,
                          booking_id INT REFERENCES bookings(booking_id),
                          amount DECIMAL(10,2),
                          method VARCHAR(50),
                          status VARCHAR(10) CHECK (status IN ('paid', 'refunded', 'failed')),
                          created_at TIMESTAMP DEFAULT NOW()
);


CREATE TABLE reviews (
                         review_id SERIAL PRIMARY KEY,
                         booking_id INT REFERENCES bookings(booking_id),
                         owner_id INT REFERENCES users(user_id),
                         sitter_id INT REFERENCES sitters(sitter_id),
                         rating INT CHECK (rating BETWEEN 1 AND 5),
                         comment TEXT,
                         created_at TIMESTAMP DEFAULT NOW()
);


CREATE TABLE chats (
                       chat_id SERIAL PRIMARY KEY,
                       booking_id INT REFERENCES bookings(booking_id),
                       created_at TIMESTAMP DEFAULT NOW()
);


CREATE TABLE messages (
                          message_id SERIAL PRIMARY KEY,
                          chat_id INT REFERENCES chats(chat_id),
                          sender_id INT REFERENCES users(user_id),
                          content TEXT,
                          sent_at TIMESTAMP DEFAULT NOW()
);



INSERT INTO users (full_name, email, phone, password_hash, role) VALUES
                                                                     ('Aruzhan Akhmetova', 'aruzhan@example.com', '+77010000001', 'hash1', 'owner'),
                                                                     ('Nazerke Alpyssova', 'nazerke@example.com', '+77010000002', 'hash2', 'sitter'),
                                                                     ('Anara Armankyzy', 'anara@example.com', '+77010000003', 'hash3', 'owner'),
                                                                     ('Meyrim Sultan', 'meyrim@example.com', '+77010000004', 'hash4', 'sitter'),
                                                                     ('Admin User', 'admin@nanny.kz', '+77010000005', 'hash5', 'admin');


INSERT INTO pets (owner_id, name, type, age, notes) VALUES
                                                        (1, 'Mila', 'cat', 2, 'Very calm and fluffy'),
                                                        (3, 'Bobby', 'dog', 4, 'Needs daily walk'),
                                                        (1, 'Luna', 'rodent', 1, 'Hamster, likes sunflower seeds'),
                                                        (3, 'Sharik', 'dog', 3, 'Friendly with kids'),
                                                        (1, 'Simba', 'cat', 5, 'Prefers dry food');


INSERT INTO sitters (sitter_id, experience_years, certificates, preferences, location, status) VALUES
                                                                                                   (2, 3, 'Pet Care Certificate 2022', 'Loves dogs and cats', 'Almaty', 'approved'),
                                                                                                   (4, 5, 'Veterinary Basics 2021', 'Can handle rodents', 'Astana', 'approved');


INSERT INTO services (sitter_id, type, price_per_hour, description) VALUES
                                                                        (2, 'walking', 2500.00, '1-hour walk with your dog in the park'),
                                                                        (2, 'home-care', 4000.00, 'Visits home twice a day to feed your pet'),
                                                                        (4, 'boarding', 7000.00, 'Pet stays at sitter’s place overnight'),
                                                                        (4, 'walking', 2000.00, 'Evening walks near the river'),
                                                                        (2, 'boarding', 8000.00, 'Comfortable stay for cats and small dogs');


INSERT INTO bookings (owner_id, sitter_id, pet_id, service_id, start_time, end_time, status) VALUES
                                                                                                 (1, 2, 1, 1, '2025-10-15 10:00', '2025-10-15 11:00', 'completed'),
                                                                                                 (3, 4, 2, 3, '2025-10-16 09:00', '2025-10-16 18:00', 'confirmed'),
                                                                                                 (1, 2, 5, 5, '2025-10-20 08:00', '2025-10-21 08:00', 'pending'),
                                                                                                 (3, 4, 4, 4, '2025-10-22 18:00', '2025-10-22 19:00', 'cancelled'),
                                                                                                 (1, 2, 3, 2, '2025-10-25 09:00', '2025-10-25 10:00', 'confirmed');


INSERT INTO payments (booking_id, amount, method, status) VALUES
                                                              (1, 2500.00, 'card', 'paid'),
                                                              (2, 7000.00, 'card', 'paid'),
                                                              (3, 8000.00, 'card', 'failed'), -- ← заменили pending на failed
                                                              (4, 2000.00, 'cash', 'refunded'),
                                                              (5, 4000.00, 'card', 'paid');



INSERT INTO reviews (booking_id, owner_id, sitter_id, rating, comment) VALUES
                                                                           (1, 1, 2, 5, 'Great experience, sitter was kind!'),
                                                                           (2, 3, 4, 4, 'Dog came home happy.'),
                                                                           (5, 1, 2, 5, 'Very punctual sitter.'),
                                                                           (4, 3, 4, 3, 'Cancelled late, but polite.'),
                                                                           (3, 1, 2, 4, 'Nice person, clean environment.');


INSERT INTO chats (booking_id) VALUES
                                   (1), (2), (3), (4), (5);


INSERT INTO messages (chat_id, sender_id, content) VALUES
                                                       (1, 1, 'Hello, is the time okay for tomorrow?'),
                                                       (1, 2, 'Yes, I’ll be there at 10.'),
                                                       (2, 3, 'Can you take Bobby at 9am?'),
                                                       (2, 4, 'Sure, no problem.'),
                                                       (3, 1, 'Please send photo updates during boarding.');


select * from pets;
select * from messages;
select * from users;
select * from sitters;
select * from bookings;
