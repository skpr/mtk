#!/usr/bin/make -f

export CGO_ENABLED=0

# Builds the project.
build:
	GOOS=linux GOARCH=amd64 go build -o bin/mtk-linux-amd64 -ldflags='-extldflags "-static"' github.com/skpr/mtk/cmd/mtk
	GOOS=linux GOARCH=arm64 go build -o bin/mtk-linux-arm64 -ldflags='-extldflags "-static"' github.com/skpr/mtk/cmd/mtk
	GOOS=darwin GOARCH=amd64 go build -o bin/mtk-darwin-amd64 -ldflags='-extldflags "-static"' github.com/skpr/mtk/cmd/mtk
	GOOS=darwin GOARCH=arm64 go build -o bin/mtk-darwin-arm64 -ldflags='-extldflags "-static"' github.com/skpr/mtk/cmd/mtk

# Run all lint checking with exit codes for CI.
lint:
	golint -set_exit_status `go list ./... | grep -v /vendor/`

# Run tests with coverage reporting.
test:
	go test -cover ./...

.PHONY: *
