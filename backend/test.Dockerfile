FROM golang:bookworm

RUN apt-get update && apt-get upgrade && apt-get install -y ca-certificates

ADD go.mod /go/
ADD go.sum /go/

RUN go mod download

COPY crts/keeper.crt crts/* /usr/local/share/ca-certificates/
RUN rm /usr/local/share/ca-certificates/keeper.crt && update-ca-certificates

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

RUN echo '#!/bin/bash\n exec make ${TASK}' > ./entrypoint.sh
RUN chmod +x ./entrypoint.sh

ENTRYPOINT ["./entrypoint.sh"]

