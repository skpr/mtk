FROM golang:1.12 as build
ADD . /go/src/github.com/skpr/mtk/cli
WORKDIR /go/src/github.com/skpr/mtk
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o bin/mtk github.com/skpr/mtk

FROM skpr/mtk-mysql:latest
USER root
RUN apk --no-cache add bash ca-certificates
COPY --from=build /go/src/github.com/skpr/mtk/bin/mtk /usr/local/bin/mtk
ADD scripts/database-sanitize /usr/local/bin/database-sanitize
RUN chmod +x /usr/local/bin/database-sanitize
USER mysql
WORKDIR /workspace
CMD ["/bin/bash", "-c", "mtk-dump > db.sql"]
