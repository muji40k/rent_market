
FROM golang:latest as build

ADD internal/ /go/internal/
ADD server/ /go/server/
ADD cmd/ /go/cmd/
ADD misc/ /go/misc/

ADD go.mod /go/
ADD go.sum /go/

RUN go build ./cmd/main.go

FROM golang:latest

COPY --from=build /go/main /main
ADD config.json /

ENTRYPOINT /main

