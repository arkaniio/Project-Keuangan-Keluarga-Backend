# Database Migrations

This directory is reserved for SQL migration files.

## Strategy

You can run migrations manually or use a migration tool such as
[golang-migrate](https://github.com/golang-migrate/migrate).

### Option 1 — Manual SQL

Connect to your PostgreSQL instance and run your `.sql` files directly:

```bash
psql -h localhost -p 5432 -U appuser2 -d postgres -f migrations/001_create_examples.sql
```

### Option 2 — golang-migrate CLI

1. Install the CLI:
   ```bash
   go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
   ```

2. Create a new migration:
   ```bash
   migrate create -ext sql -dir migrations -seq create_examples_table
   ```

3. Run migrations:
   ```bash
   migrate -path ./migrations \
     -database "postgres://appuser2:app123@localhost:5432/postgres?sslmode=disable" \
     up
   ```

### Naming Convention

Use sequential numbering with descriptive names:
- `000001_create_examples_table.up.sql`
- `000001_create_examples_table.down.sql`
