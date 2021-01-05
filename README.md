# MCPT Central

Simple service that centralizes server/ip address management for wireguard.

Allows automated deployments by interfacing with the wireguard kernel module to dynamically add peers.

See the client at mcpt/central-client

## Sample Config

`config.yml`:

```
interface: wg0
hubip: 10.11.10.1
listen: 127.0.0.1:8080
secret: ecdd4a92cda41933e7852148080591ddf71e5be8f8d5becd574b9cdc79157705
types:
  - name: site
    cidr: 10.11.11.0/24
  - name: judge
    cidr: 10.11.12.0/24
```

`servers.yml` (this is automatically generated)

```
- ip: 10.11.11.1
  type: site
  allowedips:
  - 10.11.10.1/32
  - 10.11.11.0/24
  metadata:
    Hostname: example-1
    IP: 3.1.6.6
  publickey: VuL9pcWE7nbVyXb/MFZQSlmCV0293YP6A1nxyuDwVwI=
  privatekey: 8EwX40omWmmKQWa2zxKD5vS3MNGTeC0t06lOA7HXVGo=
- ip: 10.11.11.2
  type: site
  allowedips:
  - 10.11.10.1/32
  - 10.11.11.0/24
  metadata:
    Hostname: example-2
    IP: 3.1.6.7
  publickey: qLeBJ75THeHuo76AQa/GFSBcRHIofuuV1uYjaOpWGl0=
  privatekey: IIUJV5J8XB7YCzfRp8Pq2zyUMwOXBftA6mnV8x7ZtF4=

```