#!/bin/bash

# Charge le jeu de données de démonstration des réservations (db/seeds/reservations_demo.sql).
#
#   ./scripts/seed_reservations.sh                 # base décrite par DB_* dans .env
#   DATABASE_URL="postgres://user:mdp@hote:5432/base" ./scripts/seed_reservations.sh
#   ./scripts/seed_reservations.sh --yes           # sans demande de confirmation
#
# Le serveur doit avoir tourné au moins une fois sur cette base (migrations à jour).

set -euo pipefail
cd "$(dirname "$0")/.."

if ! command -v psql >/dev/null 2>&1; then
  echo "psql est introuvable (brew install libpq, ou installer PostgreSQL)." >&2
  exit 1
fi

if [ -n "${DATABASE_URL:-}" ]; then
  TARGET="${DATABASE_URL#*@}"
  CONNECT=("$DATABASE_URL")
else
  [ -f .env ] || { echo "Aucun .env ni DATABASE_URL." >&2; exit 1; }
  # Seules les variables DB_* sont lues, sans exécuter le reste du fichier.
  while IFS='=' read -r key value; do
    value="${value%%#*}"
    value="${value%"${value##*[![:space:]]}"}"
    value="${value%\"}"; value="${value#\"}"
    case "$key" in
      DB_HOST) DB_HOST="$value" ;;
      DB_PORT) DB_PORT="$value" ;;
      DB_USER) DB_USER="$value" ;;
      DB_PASS) export PGPASSWORD="$value" ;;
      DB_NAME) DB_NAME="$value" ;;
    esac
  done < <(grep -E '^DB_(HOST|PORT|USER|PASS|NAME)=' .env)
  : "${DB_HOST:?DB_HOST manquant dans .env}" "${DB_USER:?}" "${DB_NAME:?}"
  TARGET="${DB_HOST}:${DB_PORT:-5432}/${DB_NAME}"
  CONNECT=(-h "$DB_HOST" -p "${DB_PORT:-5432}" -U "$DB_USER" -d "$DB_NAME")
fi

echo "Base cible : $TARGET"
echo "Les éléments « [Démo] ... » existants sont remplacés, les autres données ne sont pas modifiées."
if [ "${1:-}" != "--yes" ]; then
  read -r -p "Continuer ? (oui/non) " answer
  [ "$answer" = "oui" ] || { echo "Annulé."; exit 1; }
fi

psql "${CONNECT[@]}" -v ON_ERROR_STOP=1 -f db/seeds/reservations_demo.sql
