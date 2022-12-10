# Flight Log API Gateway

## Exports
```bash
export VERSION=0.0.0
export IP=172.18.2.3
export PORT=5030
export NAME=flight-log-api-gateway
export IMAGE_NAME=flight-log-api-gateway
```

## Build
```bash

docker build -f Dockerfile.local -t flight-log-api-gateway:$VERSION .
docker build -f Dockerfile.development -t johannesscr/flight-log-api-gateway:$VERSION .
```

## Local
### Flight Log API Gateway
```bash
# Flight Log API Gateway
docker run --name flight-log-api-gateway --net dottics-network --ip $IP -d -p $PORT:$PORT flight-log-api-gateway:$VERSION
```

## Development
### Flight Log API Gateway
```bash
# Flight Log API Gateway
docker run --name flight-log-api-gateway --net dottics-network --ip $IP -d -p $PORT:$PORT johannesscr/flight-log-api-gateway:$VERSION
```
