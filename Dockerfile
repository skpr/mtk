FROM alpine:3.21

ARG TARGETPLATFORM
COPY $TARGETPLATFORM/mtk /usr/local/bin/mtk

ENTRYPOINT ["/usr/local/bin/mtk"]
