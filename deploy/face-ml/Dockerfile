# Railway service "face-ml": Immich's machine-learning container, used only for face detection + embeddings.
# It keeps no photo: the backend sends one image per request and stores the result itself.
# Tested with v3.4.0 (services/faces/realmodel_test.go): bump only after running that test again.
ARG IMMICH_VERSION=v3.4.0
FROM ghcr.io/immich-app/immich-machine-learning:${IMMICH_VERSION}

# Listen on every interface, IPv6 included (Railway's private network), on the port the backend expects.
ENV IMMICH_HOST="::"
ENV IMMICH_PORT=3003
# Models (~300 MB for buffalo_l) are downloaded once into the /cache volume.
ENV MACHINE_LEARNING_CACHE_FOLDER=/cache
# One worker: each worker holds its own copy of the models in RAM.
ENV MACHINE_LEARNING_WORKERS=1
# Unload the models after 10 idle minutes to give the RAM back between scans.
ENV MACHINE_LEARNING_MODEL_TTL=600
