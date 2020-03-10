#!/usr/bin/make -f

IMAGE=skpr/mtk
VERSION=$(shell git describe --tags --always)

define build_image
	docker build -t $(IMAGE)-${1}:$(VERSION) -t $(IMAGE)-${1}:latest ${1}
endef

define test_image
	container-structure-test test --image $(IMAGE)-${1}:${VERSION} --config ${1}/tests.yml
endef

define push_image
	docker push $(IMAGE)-${1}:${VERSION}
	docker push $(IMAGE)-${1}:latest
endef

build:
	$(call build_image,build)
	$(call build_image,dump)
	$(call build_image,mysql)

test:
	$(call test_image,build)
	$(call test_image,dump)
	$(call test_image,mysql)

push:
	$(call push_image,build)
	$(call push_image,dump)
	$(call push_image,mysql)

.PHONY: *
