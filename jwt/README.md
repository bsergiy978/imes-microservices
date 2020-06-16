# Jwt Service

This is the Jwt service

Generated with

```
micro new github.com/COVIDEV/viq-chat-services/jwt --namespace=com.github.romatroskin.viqchat.jwt --type=service
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: com.github.romatroskin.viqchat.jwt.service.jwt
- Type: service
- Alias: jwt

## Dependencies

Micro services depend on service discovery. The default is multicast DNS, a zeroconf system.

In the event you need a resilient multi-host setup we recommend etcd.

```
# install etcd
brew install etcd

# run etcd
etcd
```

## Usage

A Makefile is included for convenience

Build the binary

```
make build
```

Run the service
```
./jwt-service
```

Build a docker image
```
make docker
```