# Database Migrations

This project uses [Ent](https://entgo.io) as the ORM and [Atlas](https://atlasgo.io) to manage versioned SQL migrations.

Atlas runs exclusively via Docker — no local installation needed.

## How it works

```
ent/schema/book.go               ← you edit this (source of truth)
        │
        ▼
go generate ./ent/...            ← regenerates ent/ Go code
        │
        ▼
go run ./cmd/migrate             ← syncs schema into bookstore_dev
        │
        ▼
docker compose run --rm atlas \
  migrate diff <name>            ← compares bookstore_dev vs migrations dir → writes .sql file
        │
        ▼
docker compose run --rm atlas \
  migrate apply                  ← applies pending .sql files to bookstore_db
```

Three databases are used:

| Database | Purpose |
|---|---|
| `bookstore_db` | Your app database. Migrations are applied here. |
| `bookstore_dev` | Desired state. `cmd/migrate` syncs the Ent schema here so Atlas can read it. |
| `bookstore_scratch` | Atlas's internal scratch space for computing diffs. Never touched manually. |

---

## Prerequisites

Docker running:

```bash
docker compose up -d
```

---

## Creating a new migration

Do this every time you change a schema file under `ent/schema/`.

```bash
# 1. Edit the schema
#    Example: add a field to ent/schema/book.go

# 2. Regenerate the Ent Go code
go generate ./ent/...
```

> **Stop here if `go generate` fails.** A common failure is a `displaywidth` compile error
> caused by `go mod tidy` downgrading a dependency that is incompatible with Go 1.26.
> Fix it by running `go get github.com/clipperhouse/displaywidth@v0.11.0` and retrying.
> Do not proceed to the next steps with a broken generate — `cmd/migrate` will panic.

```bash
# 3. Sync the new schema into bookstore_dev
go run ./cmd/migrate

# 4. Generate the versioned SQL migration file
docker compose run --rm atlas migrate diff <name> --env docker
#    Example: docker compose run --rm atlas migrate diff add_author_field --env docker

# 5. Review the generated file in ent/migrate/migrations/
#    Always read the SQL before applying it.

# 6. Apply to your local bookstore_db
docker compose run --rm atlas migrate apply --env docker
```

Steps 2 and 3 are also available as VS Code tasks (`ent: generate` and `migrate: sync dev db`).

### Why `go mod tidy` can break `go generate`

`go mod tidy` removes dependencies it considers unnecessary. `displaywidth` is a transitive
dependency of the Ent code generator — it is only used at code-gen time, not at runtime.
When tidy runs, it can downgrade it from `v0.11.0` back to `v0.6.2`, which fails to compile
on Go 1.26.

`tools/tools.go` exists to prevent this. It imports `displaywidth` and `entgo.io/ent/entc`
under a `//go:build tools` tag — the file is never compiled into the binary, but `go mod tidy`
sees the imports and keeps both packages pinned at their correct versions.

If you ever see the downgrade happen again, run:

```bash
go get github.com/clipperhouse/displaywidth@v0.11.0
go mod tidy
```

---

## Applying migrations (after pulling new code)

```bash
docker compose run --rm atlas migrate apply --env docker
```

Atlas is idempotent — it tracks which migrations have been applied and skips them.

---

## Checking migration status

```bash
docker compose run --rm atlas migrate status --env docker
```

This shows which migrations are applied and which are pending.

---

## Reverting a migration

Atlas does not auto-generate rollback SQL. The correct approach is to write a new forward migration that undoes the previous change.

```bash
# 1. Create an empty migration file
docker compose run --rm atlas migrate new revert_<name> --env docker
#    Example: docker compose run --rm atlas migrate new revert_add_author_field --env docker

# 2. Open the generated .sql file in ent/migrate/migrations/ and write the reverse DDL
#    Examples:
#      DROP COLUMN `author`;
#      ALTER TABLE `books` MODIFY `price` float NOT NULL;

# 3. Apply it
docker compose run --rm atlas migrate apply --env docker
```

For local development, it is often faster to just reset the database and re-apply everything from scratch:

```bash
docker exec bookstore_mysql mysql -uroot -proot_password \
  -e "DROP DATABASE bookstore_db; CREATE DATABASE bookstore_db;"

docker compose run --rm atlas migrate apply --env docker
```

---

## Validating migrations (CI)

Run this in CI to detect tampered or missing migration files:

```bash
docker compose run --rm atlas migrate validate --env docker
```

Atlas maintains an `atlas.sum` checksum file alongside the SQL files. This command verifies that all files are consistent with it.

---

## Adding a new field (end-to-end example)

Changes always start from the domain and ripple outward. Never start from the database.

```
internal/domain/book/book.go    ← 1. add field to aggregate + validate in NewBook
ent/schema/book.go              ← 2. add field to Ent schema
go generate ./ent/...           ← 3. regenerate Ent code
ent_book_repository.go          ← 4. update Create (SetX) and toDomainBook (map it back)
internal/application/           ← 5. update command DTO and service
internal/interface/api/         ← 6. update API input/output structs
go run ./cmd/migrate            ← 7. sync schema to bookstore_dev
atlas migrate diff <name>       ← 8. generate the SQL migration file
atlas migrate apply             ← 9. apply to bookstore_db
```

**Example — adding `ReleaseYear`:**

**Step 1 — Domain** (`internal/domain/book/book.go`)
```go
type Book struct {
    ID          uint32
    Title       string
    ISBN        ISBN
    Price       float64
    ReleaseYear int        // ← new
}

func NewBook(title string, isbn ISBN, price float64, releaseYear int) (*Book, error) {
    if releaseYear < 1450 {
        return nil, errors.New("ano de lançamento inválido")
    }
    ...
}
```

**Step 2 — Ent schema** (`ent/schema/book.go`)
```go
field.Int("release_year").Positive(),
```

**Step 3 — Regenerate**
```bash
go generate ./ent/...
```

**Step 4 — Repository** (`internal/infrastructure/persistence/ent_book_repository.go`)
```go
// in Create:
SetReleaseYear(b.ReleaseYear).

// in toDomainBook:
ReleaseYear: e.ReleaseYear,
```

**Step 5 — Application layer** — add `ReleaseYear int` to `RegisterBookCommand`, pass it to `NewBook` in the service.

**Step 6 — API layer** — add `ReleaseYear int` to the handler input/output structs with a `minimum: 1450` validation tag.

**Steps 7–9 — Migration**
```bash
go run ./cmd/migrate
docker compose run --rm atlas migrate diff add_release_year --env docker
docker compose run --rm atlas migrate apply --env docker
```

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
