export CGO_ENABLED=0
export VERSION=$(shell git describe --tags --abbrev=0)

all: sacloud-otel-collector

ocb:
	curl --proto '=https' --tlsv1.2 -fL -o ocb \
		https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/cmd%2Fbuilder%2Fv0.125.0/ocb_0.125.0_linux_amd64
	chmod +x ocb

.PHONY: build-src
build-src: ocb
	./ocb \
		--config=builder-config.yaml \
		--skip-compilation

sacloud-otel-collector: cmd/sacloud-otel-collector/*.go cmd/sacloud-otel-collector/go.*
	cd cmd/sacloud-otel-collector && go build -o ../../sacloud-otel-collector .

.PHONY: dist
dist:
	goreleaser build --snapshot --clean

.PHONY: release
release:
	goreleaser release --clean

.PHONY: clean
clean:
	rm -rf dist sacloud-otel-collector

.PHONY: docker-push
docker-push:
	docker buildx build \
		--build-arg VERSION=$(VERSION) \
		--platform linux/amd64,linux/arm64 \
		-t ghcr.io/sacloud/sacloud-otel-collector:$(VERSION) \
		-t ghcr.io/sacloud/sacloud-otel-collector:latest \
		--push \
		.
