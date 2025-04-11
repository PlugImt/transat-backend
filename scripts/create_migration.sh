#!/bin/bash

# Script pour créer une nouvelle migration SQL avec goose

# Vérifier si un nom de migration est fourni
if [ -z "$1" ]; then
  echo "Usage: $0 <nom_de_la_migration>"
  exit 1
fi

MIGRATION_NAME=$1
MIGRATIONS_DIR="db/migrations"

# Vérifier si le répertoire de migrations existe, sinon le créer
if [ ! -d "$MIGRATIONS_DIR" ]; then
  echo "Le répertoire de migrations $MIGRATIONS_DIR n'existe pas. Création..."
  mkdir -p "$MIGRATIONS_DIR"
  if [ $? -ne 0 ]; then
    echo "Erreur lors de la création du répertoire de migrations."
    exit 1
  fi
fi

# Vérifier si goose est installé
if ! command -v goose &> /dev/null
then
    echo "goose n'est pas installé. Installation en cours..."
    # Essayer d'installer goose via go install
    go install github.com/pressly/goose/v3/cmd/goose@latest
    # Vérifier à nouveau si l'installation a réussi
    if ! command -v goose &> /dev/null
    then
        echo "Impossible d'installer goose automatiquement."
        echo "Veuillez l'installer manuellement: go install github.com/pressly/goose/v3/cmd/goose@latest"
        exit 1
    fi
    echo "goose installé avec succès."
fi

# Créer la migration
echo "Création de la migration : $MIGRATION_NAME"
goose -dir "$MIGRATIONS_DIR" create "$MIGRATION_NAME" sql

if [ $? -eq 0 ]; then
  echo "Migration '$MIGRATION_NAME' créée avec succès dans $MIGRATIONS_DIR/"
else
  echo "Erreur lors de la création de la migration '$MIGRATION_NAME'."
  exit 1
fi

exit 0 