DROP TABLE IF EXISTS friends;
DROP TABLE IF EXISTS post;
DROP TABLE IF EXISTS user;
DROP TABLE IF EXISTS city;

CREATE TABLE city
(
    id   INT AUTO_INCREMENT NOT NULL,
    name VARCHAR(128)       NOT NULL,
    PRIMARY KEY (`id`)
) CHARACTER SET = utf8;

CREATE TABLE user
(
    id         INT AUTO_INCREMENT NOT NULL,
    first_name VARCHAR(128)       NOT NULL,
    last_name  VARCHAR(128)       NOT NULL,
    email      VARCHAR(128)       NOT NULL,
    password   VARCHAR(128)       NOT NULL,
    age        INT                NOT NULL default 0,
    city_id    INT                NULL,
    hobbies    TINYTEXT,
    PRIMARY KEY (`id`),
    UNIQUE INDEX (email),
    INDEX (city_id),
    FOREIGN KEY (city_id) REFERENCES city (id) ON DELETE set null
) CHARACTER SET = utf8mb4;

CREATE TABLE friends
(
    user_id   INT NOT NULL,
    friend_id INT NOT NULL,
    UNIQUE INDEX (user_id, friend_id),
    PRIMARY KEY (user_id, friend_id),
    FOREIGN KEY (user_id) REFERENCES user (id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (friend_id) REFERENCES user (id) ON DELETE CASCADE ON UPDATE CASCADE
) CHARACTER SET = utf8;

CREATE TABLE post
(
    id         INT AUTO_INCREMENT NOT NULL,
    text       TEXT,
    created_at TIMESTAMP          NOT NULL,
    user_id    INT                NOT NULL,
    PRIMARY KEY (`id`),
    INDEX (user_id),
    FOREIGN KEY (user_id) REFERENCES user (id) ON DELETE CASCADE ON UPDATE CASCADE
) CHARACTER SET = utf8mb4;

INSERT INTO city
(name)
VALUES ('Saint-Petersburg'),
       ('Moscow');

INSERT INTO user
(first_name, last_name, email, password, age, city_id)
VALUES ('John', 'Smith', 'j.smith@social.net', '', 30, 1),
       ('John', 'Doe', 'j.doe@social.net', '', 25, 2);

