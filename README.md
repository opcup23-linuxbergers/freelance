# freelance

## Deployment

### Manual

API:
```
make build
./out
```

Client:
```
bun run install
bun build
bun .output/server/index.mjs
```

### Docker Compose
При необходимости изменить config.json, а затем:
```
docker-compose up --build
```

По умолчанию апи и клиент будут запущены на портах 23001 и 23000 соответственно.

### Docker Container

API:
```
cd backend
docker build -t fl-api .
docker volume create {fl-api,fl-db}
docker network create fl
docker run  --name fl-api --rm --detach --mount source=fl-data,target=/app/data --mount type=bind,source=/path/to/config.json,target=/app/config.json --network fl -p 23001:18000 fl-api
```

Client:
```
cd frontend
docker build -t fl-ui .
docker run --name fl-ui --rm --detach -p 23000:3000 fl-ui
```

DB:
```
docker run -d --name postgres --mount source=fl-db,target=/var/lib/postgresql/data --network fl --rm -e POSTGRES_PASSWORD=admin -e POSTGRES_DB=fl -e POSTGRES_USER=admin postgres:latest
```
