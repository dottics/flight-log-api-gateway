# Flight Log API Gateway

## Build
```bash
export VERSION=0.0.0
docker build -f Dockerfile.local -t johannesscr/flight-log-api-gateway:$VERSION-local .
docker build -f Dockerfile.development -t johannesscr/flight-log-api-gateway:$VERSION .
```

## Development
### Flight Log API Gateway
```bash
# Flight Log API Gateway
docker run --name flight-log-api-gateway --net dottics-network --ip 172.18.2.3 -d -p 5030:5030 johannesscr/flight-log-api-gateway:$VERSION
```
