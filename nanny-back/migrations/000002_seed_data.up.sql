INSERT INTO users (full_name, email, phone, password_hash, role) VALUES
('Aruzhan Akhmetova', 'aruzhan@example.com', '+77010000001', 'hash1', 'owner'),
('Nazerke Alpyssova', 'nazerke@example.com', '+77010000002', 'hash2', 'sitter'),
('Anara Armankyzy', 'anara@example.com', '+77010000003', 'hash3', 'owner'),
('Meyrim Sultan', 'meyrim@example.com', '+77010000004', 'hash4', 'sitter'),
('Admin User', 'admin@nanny.kz', '+77010000005', 'hash5', 'admin');

INSERT INTO pets (owner_id, name, type, age, notes) VALUES
(1, 'Mila', 'cat', 2, 'Very calm and fluffy'),
(3, 'Bobby', 'dog', 4, 'Needs daily walk'),
(1, 'Luna', 'rodent', 1, 'Hamster'),
(3, 'Sharik', 'dog', 3, 'Friendly with kids'),
(1, 'Simba', 'cat', 5, 'Prefers dry food');

INSERT INTO sitters (sitter_id, experience_years, certificates, preferences, location, status) VALUES
(2, 3, 'Pet Care Certificate 2022', 'Loves dogs and cats', 'Almaty', 'approved'),
(4, 5, 'Veterinary Basics 2021', 'Can handle rodents', 'Astana', 'approved');
