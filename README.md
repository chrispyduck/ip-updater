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

#### `kustomization.yaml`
```yaml
kind: Kustomization
resources: 
- github.com/chrispyduck/ip-updater.git?ref=0.1.2
images:
- name: ip-updater
  newName: your.registry/ip-updater
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