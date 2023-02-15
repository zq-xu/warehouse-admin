FROM golang:1.17 as builder

ENV PROJECT=zq-xu/warehouse-admin
ENV WORKDIR=/GoProject/src/$PROJECT

WORKDIR $WORKDIR
ENV GOPATH=/GoProject
ENV GO111MODULE=on

RUN mkdir -p /cicd

COPY cmd $WORKDIR/cmd
COPY pkg $WORKDIR/pkg
COPY internal $WORKDIR/internal
COPY go.mod $WORKDIR/go.mod


RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-w -s -extldflags "-static"' -tags netgo -o /warehouse-admin/webserver ./cmd/webserver/webserver.go

FROM hairyhenderson/upx:3.94 as upx
COPY --from=builder /warehouse-admin /warehouse-admin
RUN upx -9 /warehouse-admin/webserver

FROM alpine:3.15.0 as prod
ARG version
ENV COMMITID=$version

COPY --from=upx /warehouse-admin  /warehouse-admin

ENTRYPOINT ["sh","-c", "/warehouse-admin/webserver"]
