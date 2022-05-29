run:
  @go build && ./practical-golang

search:
  @go build && ./practical-golang --default-language=lan search test

check:
  errcheck ch5/errcheck.go