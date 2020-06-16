# Centrifuge Service

This is the Centrifuge service

Generated with

```
micro new github.com/COVIDEV/viq-chat-services/centrifuge --namespace=com.github.romatroskin.viqchat.centrifuge --type=service
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: com.github.romatroskin.viqchat.centrifuge.service.centrifuge
- Type: service
- Alias: centrifuge

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
./centrifuge-service
```

Build a docker image
```
make docker
```