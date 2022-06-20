VERSION = 0.2.0
REGISTRY = your.registry
IMAGE = chrispyduck/ip-updater

.PHONY: docker local push

local: ip-updater
ip-updater: 
	go get -d -v 
	go build -v .

docker:
	docker buildx build \
		--platform arm64,amd64 \
		--push \
		-t $(REGISTRY)/$(IMAGE):$(VERSION) \
		.
