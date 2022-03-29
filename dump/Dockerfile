ARG IMAGE_TAG=v2-latest
FROM golang:1.12 as build
ARG ARCH=amd64
ARG IMAGE_TAG

ADD . /go/src/github.com/skpr/mtk/dump
WORKDIR /go/src/github.com/skpr/mtk/dump
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${ARCH} go build -a -o bin/mtk-dump github.com/skpr/mtk/dump

FROM skpr/mtk-mysql:${IMAGE_TAG}
USER root
RUN apk --no-cache add bash ca-certificates
COPY --from=build /go/src/github.com/skpr/mtk/dump/bin/mtk-dump /usr/local/bin/mtk-dump
ADD scripts/database-sanitize /usr/local/bin/database-sanitize
RUN chmod +x /usr/local/bin/database-sanitize
USER mysql
WORKDIR /workspace
CMD ["/bin/bash", "-c", "mtk-dump > db.sql"]
