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