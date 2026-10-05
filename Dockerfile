
FROM golang:1.23-alpine

ADD . /line-login-go
WORKDIR /line-login-go

RUN export GOFLAGS=-mod=vendor
RUN cd /line-login-go && go build -o line-login-go

ENTRYPOINT ["/line-login-go/line-login-go"]
