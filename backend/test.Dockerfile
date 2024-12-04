FROM golang:latest

ADD go.mod /go/
ADD go.sum /go/

RUN go mod download

ADD application/ /go/application/
ADD builders/ /go/builders
ADD constructors/ /go/constructors
ADD logger/ /go/logger
ADD internal/ /go/internal/
ADD server/ /go/server/
ADD cmd/ /go/cmd/
ADD misc/ /go/misc/
ADD tests/ /go/tests/
ADD allure-report/Makefile /go/allure-report/Makefile
ADD Makefile /go/Makefile

ENV ALLURE_OUTPUT_PATH=/go/allure-report

ENTRYPOINT ["make", "$TASK"]

