artifact_name := "practical-golang"

run: && exec
  @go build -o {{artifact_name}} main.go

search:
  @go build -o {{artifact_name}} main.go && ./{{artifact_name}} --default-language=lan search test

check:
  errcheck ch5/errcheck.go

run-mac-m1: && exec
  @CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o {{artifact_name}} -trimpath -ldflags='-s -w -X main.version=1.0.0' main.go

exec:
  @./{{artifact_name}}

create-db:
  docker run -d --name my-postgres -e POSTGRES_USER=testuser -e POSTGRES_PASSWORD=pass -e POSTGRES_DB=testdb -p 5432:5432 postgres

restart-db:
  docker start my-postgres

stop-db:
  docker stop my-postgres

db-in:
  docker exec -it my-postgres bash -c "psql testdb -U testuser"

sqlboiler:
  sqlboiler psql -c ch9/sqlboiler.toml -o ch9/models -p models --no-tests --wipe

sqlc:
  sqlc generate -f ch9/sqlc.json

start-otel:
  docker run -d --name jaeger --rm -p 14268:14268 -p 16686:16686 jaegertracing/all-in-one:1.36
  echo "view Jaeger UI at http://localhost:16686"

stop-otel:
  docker stop jaeger

cover: && html
  go test ./ch13 -coverprofile=./ch13/coverage.out

html:
  go tool cover -html=./ch13/coverage.out