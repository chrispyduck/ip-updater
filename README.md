# Cloudflare Dynamic DNS Updater

## Why?
1. Get better at Golang
2. Existing tools don't do the job well or the way I like

## Installation

### 1. Build Docker Image
```bash
$ make docker REGISTRY=your.registry
```

### 2. Create Kustomization to install into Kubernetes

#### `kustomize.yaml`
```yaml
kind: Kustomize
resources: git::github.com/chrispyduck/ip-updater.git?ref=v0.1
images:
- name: ip-updater
imageName: your.registry/ip-updater
secretGenerator:
- name: ip-updater
files:
- ip-updater.yaml
``` 

#### `ip-updater.yaml`
```yaml
debug: false
queryUrls:
  v4: https://api4.my-ip.io/ip.txt
  v6: https://api6.my-ip.io/ip.txt
hostname: your.dynamic.hostname.hosted.on.your.cloudflare.domain
cloudflare:
  apiToken: ...
  zone: your.cloudflare.domain
```