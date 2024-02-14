-- CREATE SCHEMA IF NOT EXISTS `art_gallery`;

-- CREATE TABLE IF NOT EXISTS `art_gallery`.`genres` (
-- `id` VARCHAR(36) NOT NULL,
-- `name` VARCHAR(1024) NOT NULL,
-- PRIMARY KEY (`id`));

-- CREATE TABLE IF NOT EXISTS `art_gallery`.`users` (
-- `id` VARCHAR(36) NOT NULL,
-- `first_name` VARCHAR(1024) NULL,
-- `last_name` VARCHAR(1024) NULL,
-- `date_of_registration` VARCHAR(1024) NOT NULL,
-- `email` VARCHAR(1024) NOT NULL,
-- `password` VARCHAR(1024) NULL,
-- PRIMARY KEY (`id`));

-- CREATE TABLE IF NOT EXISTS `art_gallery`.`paintings` (
-- `id` VARCHAR(36) NOT NULL,
-- `title` VARCHAR(1024) NULL,
-- `description` VARCHAR(1024) NULL,
-- `mime_type` VARCHAR(255),
-- `data` LONGBLOB,
-- `author` VARCHAR(1024) NOT NULL,
-- `date_of_publication` VARCHAR(1024) NULL,
-- `width`  INT NULL,
-- `height`  INT NULL,
-- `genre`  VARCHAR(36),
-- `price` DOUBLE,
-- PRIMARY KEY (`id`),
-- FOREIGN KEY (genre) REFERENCES genres(id));
