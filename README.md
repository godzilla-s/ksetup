# Ksetup

ksetup is a tool for help us to setup k3s cluster quickly and easily.

# Usage

```
ksetup install -f config.yaml
```

## config.yaml 

`clusrer` define cluster config:

```yaml
cluster1: &k3s-1 
  master:
    config:
      disable:
        - traefik
        - servicelb
        - local-storage
      flannel-backend: "vxlan"
    members:
      - *node-01 
  agent:
    config:
      kubelet-arg:
        - --max-pods=100
    members: 
      - *node-02
      - *node-03
  charts:
    - metallb 
    - longhorn
```

there are 3 parts: 

+ master/agent: define master/agent nodes

+ charts: define charts to install

+ manifest: define manifest to install

nodes can be define as :

```yaml 
node1: &node-01
  host: 192.168.100.10
  hostname: k3s-master-1
  root_password: <PASSWORD>
  pre_install:
    - k3s 
    - installsh
  post_install:
    - docker 
```