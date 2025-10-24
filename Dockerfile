FROM alpine:3.21

# dockers_v2 version temporarily disabled.
# ARG TARGETPLATFORM
# COPY $TARGETPLATFORM/mtk /usr/local/bin/mtk

COPY mtk /usr/local/bin/mtk

ENTRYPOINT ["/usr/local/bin/mtk"]
