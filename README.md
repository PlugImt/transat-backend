# Transat-backend

Le backend de Transat est un serveur Fiber, développé en Go, connecté à une base de données PostgreSQL.

# Comment lancer le projet
```bash
# Copier le fichier .env.example en .env
cp .env.example .env
# Remplacez les variables d'environnement par les vôtres si nécessaire, mais ça fonctionnera de base (sans mail)

# Lancer la base de données PostgreSQL (optionnel, sinon il faut configurer la base de données manuellement dans le .env)
docker compose up -d db

# Lancer le serveur
go run main.go
```
 
## Docker
Pour simplifier le développement, il est possible d'utiliser Docker Compose pour lancer une base de données PostgreSQL (configurée via le fichier `.env`).

```bash
# Lancer le serveur
docker compose up -d
```

## Schéma MCD
Le schéma MCD est disponible dans le fichier `MCD.lo1`, conçu avec le logiciel [Looping](https://www.looping-mcd.fr/).

## Sentry
Sentry est utilisé pour le suivi des erreurs et le monitoring des crons. Nous les remercions pour leur offre spéciale.
Demandez à être ajouté à la team Transat !

## Informations utiles
- Version de Go : `1.24.3`
