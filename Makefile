#!/usr/bin/make -f

IMAGE_REPO_BASE=skpr/mtk
ARCH=amd64
VERSION_TAG=v3-latest

define build_image
	docker build --build-arg ARCH=${ARCH} -t ${IMAGE_REPO_BASE}-${1}:${VERSION_TAG}-${ARCH} ${1}
endef

define test_image
	container-structure-test test --image ${IMAGE_REPO_BASE}-${1}:${VERSION_TAG}-${ARCH} --config ${1}/tests.yml
endef

define push_image
	docker push ${IMAGE_REPO_BASE}-${1}:${VERSION_TAG}-${ARCH}
endef

define manifest
	$(eval IMAGE=${IMAGE_REPO_BASE}-${1}:${VERSION_TAG})
	docker manifest create ${IMAGE} --amend ${IMAGE}-arm64 --amend ${IMAGE}-amd64
    docker manifest push ${IMAGE}
endef

build:
	$(call build_image,build)
	$(call build_image,mysql)
	$(call build_image,dump)
	$(call build_image,empty)

test:
	$(call test_image,build)
	$(call test_image,mysql)
	$(call test_image,dump)
	$(call test_image,empty)

push:
#	$(call push_image,build)
#	$(call push_image,mysql)
#	$(call push_image,dump)
#	$(call push_image,empty)

manifest:
	$(call manifest,build)
	$(call manifest,mysql)
	$(call manifest,dump)
	$(call manifest,empty)

.PHONY: *
