#!/usr/bin/make -f

export CGO_ENABLED=0

define build_step
	GOOS=$(1) GOARCH=$(2) go build -o bin/mtk-$(1)-$(2) -ldflags='-extldflags "-static"' github.com/skpr/mtk/cmd/mtk
endef

# Builds the project.
build:
	$(call build_step,linux,amd64)
	$(call build_step,linux,arm64)
	$(call build_step,darwin,amd64)
	$(call build_step,darwin,arm64)

.PHONY: *
