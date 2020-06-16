# Messages Service

This is the Messages service

Generated with

```
micro new github.com/COVIDEV/viq-chat-services/messages --namespace=com.github.romatroskin.viqchat.messages --type=service
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: com.github.romatroskin.viqchat.messages.service.messages
- Type: service
- Alias: messages

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
./messages-service
```

Build a docker image
```
make docker
```