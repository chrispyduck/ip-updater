VERSION = 0.1.2
REGISTRY = your.registry
IMAGE = chrispyduck/ip-updater

.PHONY: docker local push

local: ip-updater
ip-updater: 
	go get -d -v 
	go build -v .

docker:
	docker buildx build -t $(REGISTRY)/$(IMAGE):$(VERSION) . --progress plain

push: docker 
	docker push $(REGISTRY)/$(IMAGE):$(VERSION)