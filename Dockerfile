FROM golang:1.23 as build

WORKDIR /go/src/github.com/skpr/mtk
COPY . /go/src/github.com/skpr/mtk

ENV CGO_ENABLED=0

RUN go build -o bin/mtk -ldflags='-extldflags "-static"' github.com/skpr/mtk/cmd/mtk

FROM alpine:3.18

COPY --from=build /go/src/github.com/skpr/mtk/bin/mtk /usr/local/bin/mtk

ENTRYPOINT ["/usr/local/bin/mtk"]
