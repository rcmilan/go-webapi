# env "local" — used when running atlas CLI locally (atlas migrate diff / apply)
#
# Workflow:
#   1. go run ./cmd/migrate          → syncs Ent schema into bookstore_dev
#   2. atlas migrate diff <name> --env local  → diffs bookstore_dev vs migrations dir
#   3. atlas migrate apply --env local        → applies pending migrations to bookstore_db
#
env "local" {
  src = "mysql://root:root_password@localhost:3306/bookstore_dev"
  dev = "mysql://root:root_password@localhost:3306/bookstore_scratch"
  url = "mysql://root:root_password@localhost:3306/bookstore_db"

  migration {
    dir    = "file://ent/migrate/migrations"
    format = atlas
  }
}

# env "docker" — used inside the atlas Docker service (apply only, no Go required)
#   docker compose run --rm atlas migrate apply --env docker
#   docker compose run --rm atlas migrate status --env docker
#
env "docker" {
  url = "mysql://root:root_password@mysql_db:3306/bookstore_db"

  migration {
    dir    = "file://ent/migrate/migrations"
    format = atlas
  }
}
