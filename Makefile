#!/usr/bin/make -f

IMAGE=skpr/mtk
VERSION=$(shell git describe --tags --always)

define build_and_push
	docker build -t $(IMAGE)-${1}:$(VERSION) -t $(IMAGE)-${1}:latest ${1}
	docker push $(IMAGE)-${1}:${VERSION}
	docker push $(IMAGE)-${1}:latest
endef

build:
	$(call build_and_push,dump)
	$(call build_and_push,package)
	$(call build_and_push,mysql)

.PHONY: *
