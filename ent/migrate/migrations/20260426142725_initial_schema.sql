-- Create "books" table
CREATE TABLE `books` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL,
  `isbn` varchar(13) NOT NULL,
  `price` double NOT NULL,
  `created_at` bigint NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `isbn` (`isbn`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
