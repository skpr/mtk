FROM golang:1.12 as build
ADD . /go/src/github.com/skpr/mtk/dump
WORKDIR /go/src/github.com/skpr/mtk/dump
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o bin/mtk-dump github.com/skpr/mtk/dump

FROM alpine:3.10
RUN apk --no-cache add bash ca-certificates
COPY --from=build /go/src/github.com/skpr/mtk/dump/bin/mtk-dump /usr/local/bin/mtk-dump
WORKDIR /workspace
CMD ["/bin/bash", "-c", "mtk-dump > db.sql"]
