# All Atlas commands run via Docker:
#   docker compose run --rm atlas migrate diff <name>
#   docker compose run --rm atlas migrate apply
#   docker compose run --rm atlas migrate status
#   docker compose run --rm atlas migrate validate

env "docker" {
  src = "mysql://root:root_password@mysql_db:3306/bookstore_dev"
  dev = "mysql://root:root_password@mysql_db:3306/bookstore_scratch"
  url = "mysql://root:root_password@mysql_db:3306/bookstore_db"

  migration {
    dir    = "file://ent/migrate/migrations"
    format = atlas
  }
}
