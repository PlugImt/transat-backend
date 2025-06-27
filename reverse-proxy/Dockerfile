FROM nginx:alpine

# Variables d'environnement par défaut
ENV PORT=8080
ENV BACKEND_URL=http://backend:3000

# Copier la configuration nginx personnalisée
COPY nginx.conf /etc/nginx/nginx.conf

# Exposer le port (Railway utilise la variable d'environnement PORT)
EXPOSE 8080

# Démarrer nginx
CMD ["sh", "-c", "envsubst '$$PORT,$$BACKEND_URL' < /etc/nginx/nginx.conf > /tmp/nginx.conf && cp /tmp/nginx.conf /etc/nginx/nginx.conf && nginx -g 'daemon off;'"]
