FROM golang:bookworm

RUN mkdir /allure; cd /allure; \
    wget "https://github.com/allure-framework/allure2/releases/download/2.32.0/allure_2.32.0-1_all.deb"; \
    apt-get update; apt-get upgrade; \
    dpkg -i ./allure_2.32.0-1_all.deb; \
    apt-get install -f -y; \
    cd /go;

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
ADD allure-report/Makefile /go/allure-report/Makefile
ADD Makefile /go/Makefile

ENV ALLURE_OUTPUT_PATH=/go/allure-report

ENTRYPOINT make $TASK

