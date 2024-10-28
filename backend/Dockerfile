
FROM golang:latest AS build

ADD go.mod /go/
ADD go.sum /go/

RUN go mod download

ADD application/ /go/application/
ADD builders/ /go/builders
ADD constructors/ /go/constructors
ADD internal/ /go/internal/
ADD server/ /go/server/
ADD cmd/ /go/cmd/
ADD misc/ /go/misc/

RUN go build ./cmd/main.go

FROM golang:latest

COPY --from=build /go/main /main

ENTRYPOINT /main

