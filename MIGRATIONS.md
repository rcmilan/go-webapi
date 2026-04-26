# Database Migrations

This project uses [Ent](https://entgo.io) as the ORM and [Atlas](https://atlasgo.io) to manage versioned SQL migrations.

## How it works

```
ent/schema/book.go          ← you edit this (source of truth)
        │
        ▼
go generate ./ent/...       ← regenerates ent/ Go code
        │
        ▼
go run ./cmd/migrate        ← syncs schema into bookstore_dev
        │
        ▼
atlas migrate diff <name>   ← compares bookstore_dev vs migrations dir → writes .sql file
        │
        ▼
atlas migrate apply         ← applies pending .sql files to bookstore_db
```

Three databases are used:

| Database | Purpose |
|---|---|
| `bookstore_db` | Your app database. Migrations are applied here. |
| `bookstore_dev` | Desired state. `cmd/migrate` syncs the Ent schema here so Atlas can read it. |
| `bookstore_scratch` | Atlas's internal scratch space for computing diffs. Never touched manually. |

---

## Prerequisites

- Docker running: `docker compose up -d`
- Atlas CLI installed at `~/go/bin/atlas.exe` (already set up)

---

## Creating a new migration

Do this every time you change a schema file under `ent/schema/`.

```bash
# 1. Edit the schema
#    Example: add a field to ent/schema/book.go

# 2. Regenerate the Ent Go code
go generate ./ent/...

# 3. Sync the new schema into bookstore_dev
go run ./cmd/migrate

# 4. Generate the versioned SQL migration file
atlas migrate diff <name> --env local
#    Example: atlas migrate diff add_author_field --env local

# 5. Review the generated file in ent/migrate/migrations/
#    Always read the SQL before applying it.

# 6. Apply to your local bookstore_db
atlas migrate apply --env local
```

Steps 2 and 3 are also available as VS Code tasks (`ent: generate` and `migrate: sync dev db`).

---

## Applying migrations (after pulling new code)

If a teammate added a migration, apply it after pulling:

```bash
# Locally
atlas migrate apply --env local

# Via Docker (no local Atlas CLI needed)
docker compose run --rm atlas migrate apply --env docker
```

Atlas is idempotent — it tracks which migrations have been applied and skips them.

---

## Checking migration status

```bash
atlas migrate status --env local

# or via Docker
docker compose run --rm atlas migrate status --env docker
```

This shows which migrations are applied and which are pending.

---

## Reverting a migration

Atlas does not auto-generate rollback SQL. The correct approach is to write a new forward migration that undoes the previous change.

```bash
# 1. Create an empty migration file
atlas migrate new revert_<name> --env local
#    Example: atlas migrate new revert_add_author_field --env local

# 2. Open the generated .sql file in ent/migrate/migrations/ and write the reverse DDL
#    Examples:
#      DROP COLUMN `author`;
#      ALTER TABLE `books` MODIFY `price` float NOT NULL;

# 3. Apply it
atlas migrate apply --env local
```

For local development, it is often faster to just reset the database and re-apply everything from scratch:

```bash
docker exec bookstore_mysql mysql -uroot -proot_password \
  -e "DROP DATABASE bookstore_db; CREATE DATABASE bookstore_db;"

atlas migrate apply --env local
```

---

## Validating migrations (CI)

Run this in CI to detect tampered or missing migration files:

```bash
atlas migrate validate --env local
```

Atlas maintains an `atlas.sum` checksum file alongside the SQL files. This command verifies that all files are consistent with it.

---

## Schema file reference

The Ent schema lives in `ent/schema/`. Each file defines one entity (database table).

```go
// ent/schema/book.go
func (Book) Fields() []ent.Field {
    return []ent.Field{
        field.Uint32("id").SchemaType(map[string]string{dialect.MySQL: "int unsigned"}).Immutable(),
        field.String("title").MaxLen(255).NotEmpty(),
        field.String("isbn").MaxLen(13).Unique().NotEmpty(),
        field.Float("price").Positive(),
        field.Int64("created_at").Immutable().DefaultFunc(...),
    }
}
```

Full Ent field reference: https://entgo.io/docs/schema-fields
