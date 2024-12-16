
FROM golang:latest

WORKDIR /code2test
ADD . .
RUN go mod download
RUN go install github.com/fzipp/gocyclo/cmd/gocyclo@latest

ENTRYPOINT ["./metrics.sh"]

