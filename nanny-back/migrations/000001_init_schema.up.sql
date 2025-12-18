CREATE TABLE users (
                       user_id SERIAL PRIMARY KEY,
                       full_name VARCHAR(100) NOT NULL,
                       email VARCHAR(100) UNIQUE NOT NULL,
                       phone VARCHAR(20) NOT NULL,
                       password_hash VARCHAR(255) NOT NULL,
                       role VARCHAR(10) CHECK (role IN ('owner', 'sitter', 'admin')) NOT NULL,
                       created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE pets (
                      pet_id SERIAL PRIMARY KEY,
                      owner_id INT REFERENCES users(user_id) ON DELETE CASCADE,
                      name VARCHAR(50) NOT NULL,
                      type VARCHAR(20) CHECK (type IN ('cat', 'dog', 'rodent')) NOT NULL,
                      age INT,
                      notes TEXT
);

CREATE TABLE sitters (
                         sitter_id INT PRIMARY KEY REFERENCES users(user_id) ON DELETE CASCADE,
                         experience_years INT,
                         certificates TEXT,
                         preferences TEXT,
                         location VARCHAR(100),
                         status VARCHAR(10) CHECK (status IN ('pending', 'approved', 'rejected')) NOT NULL
);

CREATE TABLE services (
                          service_id SERIAL PRIMARY KEY,
                          sitter_id INT REFERENCES sitters(sitter_id) ON DELETE CASCADE,
                          type VARCHAR(20) CHECK (type IN ('walking', 'boarding', 'home-care')) NOT NULL,
                          price_per_hour DECIMAL(10,2) NOT NULL,
                          description TEXT
);

CREATE TABLE bookings (
                          booking_id SERIAL PRIMARY KEY,
                          owner_id INT REFERENCES users(user_id),
                          sitter_id INT REFERENCES sitters(sitter_id),
                          pet_id INT REFERENCES pets(pet_id),
                          service_id INT REFERENCES services(service_id),
                          start_time TIMESTAMP NOT NULL,
                          end_time TIMESTAMP NOT NULL,
                          status VARCHAR(15)
                              CHECK (status IN ('pending', 'confirmed', 'cancelled', 'completed'))
                              DEFAULT 'pending'
);

CREATE TABLE payments (
                          payment_id SERIAL PRIMARY KEY,
                          booking_id INT REFERENCES bookings(booking_id) ON DELETE CASCADE,
                          amount DECIMAL(10,2) NOT NULL,
                          method VARCHAR(50),
                          status VARCHAR(10)
                              CHECK (status IN ('paid', 'refunded', 'failed')) NOT NULL,
                          created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE reviews (
                         review_id SERIAL PRIMARY KEY,
                         booking_id INT REFERENCES bookings(booking_id) ON DELETE CASCADE,
                         owner_id INT REFERENCES users(user_id),
                         sitter_id INT REFERENCES sitters(sitter_id),
                         rating INT CHECK (rating BETWEEN 1 AND 5),
                         comment TEXT,
                         created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE chats (
                       chat_id SERIAL PRIMARY KEY,
                       booking_id INT REFERENCES bookings(booking_id) ON DELETE CASCADE,
                       created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE messages (
                          message_id SERIAL PRIMARY KEY,
                          chat_id INT REFERENCES chats(chat_id) ON DELETE CASCADE,
                          sender_id INT REFERENCES users(user_id),
                          content TEXT NOT NULL,
                          sent_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_pets_owner ON pets(owner_id);
CREATE INDEX idx_services_sitter ON services(sitter_id);
CREATE INDEX idx_bookings_owner ON bookings(owner_id);
CREATE INDEX idx_bookings_sitter ON bookings(sitter_id);
CREATE INDEX idx_bookings_status ON bookings(status);
CREATE INDEX idx_bookings_start_time ON bookings(start_time);
CREATE INDEX idx_reviews_sitter ON reviews(sitter_id);
CREATE INDEX idx_reviews_booking ON reviews(booking_id);
CREATE INDEX idx_chats_booking ON chats(booking_id);
CREATE INDEX idx_messages_chat ON messages(chat_id);