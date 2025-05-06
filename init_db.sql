-- Таблица пользователей
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    email VARCHAR(255) UNIQUE,
    password_hash VARCHAR(100),
    age INTEGER,
    description TEXT,
    city VARCHAR(100),
    coordinates POINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица предпочтений
CREATE TABLE preferences (
    user_id INTEGER PRIMARY KEY,
    gender VARCHAR(10),
    age_from INTEGER,
    age_to INTEGER,
    radius INTEGER,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Таблица фотографий
CREATE TABLE photos (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    photo_url TEXT,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Таблица совпадений
CREATE TABLE matches (
    id SERIAL PRIMARY KEY,
    user_id_1 INTEGER,
    user_id_2 INTEGER,
    FOREIGN KEY (user_id_1) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id_2) REFERENCES users(id) ON DELETE CASCADE
);

INSERT INTO users (name, email, password_hash, age, description, city, coordinates) VALUES
('Алексей', 'aleksei@example.com', '$2a$10$hash1', 25, 'Люблю путешествия и активный отдых', 'Москва', '(37.6176, 55.7558)'),
('Мария', 'maria@example.com', '$2a$10$hash2', 22, 'Ищу серьезные отношения', 'Санкт-Петербург', '(30.3159, 59.9391)'),
('Иван', 'ivan@example.com', '$2a$10$hash3', 30, 'Программист, люблю спорт', 'Новосибирск', '(82.9204, 55.0084)'),
('Елена', 'elena@example.com', '$2a$10$hash4', 28, 'Творческая личность', 'Казань', '(49.1088, 55.7963)'),
('Дмитрий', 'dmitry@example.com', '$2a$10$hash5', 35, 'Бизнесмен, автомобилист', 'Сочи', '(39.7230, 43.5855)');

INSERT INTO preferences (user_id, gender, age_from, age_to, radius) VALUES
(1, 'Женский', 20, 30, 50),
(2, 'Мужской', 25, 35, 30),
(3, 'Женский', 22, 32, 100),
(4, 'Мужской', 28, 38, 80),
(5, 'Женский', 25, 35, 150);

INSERT INTO photos (user_id, photo_url) VALUES
(1, 'https://70127a1b-bd9a-4d77-b11f-3b28f501631e.selstorage.ru/23%2FIMG_3263.jpg'),
(1, 'https://70127a1b-bd9a-4d77-b11f-3b28f501631e.selstorage.ru/23%2FIMG_3263.jpg'),
(2, 'https://70127a1b-bd9a-4d77-b11f-3b28f501631e.selstorage.ru/23%2FIMG_3263.jpg'),
(3, 'https://70127a1b-bd9a-4d77-b11f-3b28f501631e.selstorage.ru/23%2FIMG_3263.jpg'),
(3, 'https://70127a1b-bd9a-4d77-b11f-3b28f501631e.selstorage.ru/23%2FIMG_3263.jpg'),
(3, 'https://70127a1b-bd9a-4d77-b11f-3b28f501631e.selstorage.ru/23%2FIMG_3263.jpg'),
(4, 'https://70127a1b-bd9a-4d77-b11f-3b28f501631e.selstorage.ru/23%2FIMG_3263.jpg'),
(5, 'https://70127a1b-bd9a-4d77-b11f-3b28f501631e.selstorage.ru/23%2FIMG_3263.jpg');

INSERT INTO matches (user_id_1, user_id_2) VALUES
(1, 2),
(1, 3),
(2, 4),
(3, 5),
(4, 5);