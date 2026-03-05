#!/bin/bash
# run_migrations.sh — apply all TeslaPay SQL migrations in order.
# Usage:
#   DB_HOST=localhost DB_USER=teslapay DB_PASS=teslapay DB_NAME=teslapay ./run_migrations.sh
set -euo pipefail

DB_HOST="${DB_HOST:-localhost}"
DB_USER="${DB_USER:-teslapay}"
DB_PASS="${DB_PASS:-teslapay}"
DB_NAME="${DB_NAME:-teslapay}"

export PGPASSWORD="$DB_PASS"
PSQL="psql -h $DB_HOST -U $DB_USER -d $DB_NAME"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "Running TeslaPay migrations against $DB_NAME@$DB_HOST..."
for f in "$SCRIPT_DIR"/0*.sql; do
  echo "  Applying $(basename "$f")..."
  $PSQL -f "$f"
done
echo "All migrations applied!"
