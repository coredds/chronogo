SHELL := /usr/bin/env bash

.PHONY: test cover race lint bench ci fmt

test:
	go test ./...

cover:
	go test -covermode=count -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

race:
	CGO_ENABLED=1 go test -race ./...

lint:
	@echo "golangci-lint is optional; install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
	@command -v golangci-lint >/dev/null 2>&1 && golangci-lint run || echo "golangci-lint not installed"

bench:
	go test -bench=. ./...

fmt:
	go fmt ./...

ci: fmt test cover


