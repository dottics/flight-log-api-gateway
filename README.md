# Flight Log API Gateway

## Exports
```bash
export VERSION=0.5.0
export IP=172.18.2.3
export PORT=5030
export NAME=flight-log-api-gateway
export IMAGE_NAME=flight-log-api-gateway
```

## Build
```bash

docker build -f Dockerfile.local -t $NAME:$VERSION .
docker build -f Dockerfile.development -t johannesscr/$NAME:$VERSION .
```

## Local
### Flight Log API Gateway
```bash
# Flight Log API Gateway
docker run --name $NAME --net dottics-network --ip $IP -d -p $PORT:$PORT $IMAGE_NAME:$VERSION
```

## Development
### Flight Log API Gateway
```bash
# Flight Log API Gateway
docker run --name $NAME --net dottics-network --ip $IP -d -p $PORT:$PORT johannesscr/$IMAGE_NAME:$VERSION
```

## Basic HTTP
###
```bash
export TOKEN=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzLXRva2VuIjoiNWVhM2Q0MzctMzQzMC00YTgwLTk2MGEtNjQwNTVlNzcwMTJjIiwicy1pZCI6Ijc4OWU4NDAxMTIyYTBmYmQ2M2NkM2JjNjhkMTQ5NzlmODc3NjZiMTk1MzdiZThkYmRjNDFmNTE4ZDFjZWViY2QiLCJ1LWlkIjoiMWNhMGFlNjgtMWJmMi00YTE4LWE4MTktYmU1YWE4MGVkOThlIiwiY3JlYXRlZCI6IjEyLzIwLzIwMjIsIDA4OjE0OjEyIn0.lxqBTlDlSFLGwiaMjrXm1Fdoh4-Zhp6cNTUCbBkYZhw

curl -i -X POST http://localhost:5030/login -H "content-type=application/json" -d '{"email":"t@test.dottics.com","password":"test"}'

curl http://localhost:5030/aircraft-type -H "x-token=$TOKEN"
```
