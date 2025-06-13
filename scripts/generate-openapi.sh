#!/usr/bin/env bash

SCRIPT_DIR=$(dirname "$0")
cd $SCRIPT_DIR/..

echo "🔄 Génération de la documentation OpenAPI..."

if ! command -v swag &> /dev/null; then
    echo "❌ swag n'est pas installé. Installation en cours..."
    go install github.com/swaggo/swag/v2/cmd/swag@latest
fi

if ! command -v swag &> /dev/null; then
    echo "swag n'est pas installé, vérifiez si le dossier où Go installe ses binaires est dans votre PATH"
    exit 1
fi

swag init -g main.go -o ./docs --v3.1 --parseDependency --parseInternal

if [ $? -eq 0 ]; then
    echo "✅ Documentation générée avec succès dans ./docs/"
else
    echo "Erreur lors de la génération de la documentation"
    exit 1
fi