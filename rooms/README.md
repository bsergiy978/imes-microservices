# Rooms Service

This is the Rooms service

Generated with

```
micro new github.com/COVIDEV/viq-chat-services/rooms --namespace=com.github.romatroskin.viqchat.rooms --type=service
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: com.github.romatroskin.viqchat.rooms.service.rooms
- Type: service
- Alias: rooms

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
./rooms-service
```

Build a docker image
```
make docker
```