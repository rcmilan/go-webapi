-- bookstore_db is created automatically by MYSQL_DATABASE in docker-compose.
-- bookstore_dev holds the desired Ent schema (source of truth for atlas migrate diff).
-- bookstore_scratch is Atlas's internal scratch space for diff computation.
CREATE DATABASE IF NOT EXISTS bookstore_dev;
CREATE DATABASE IF NOT EXISTS bookstore_scratch;
