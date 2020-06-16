# Otp Service

This is the Otp service

Generated with

```
micro new github.com/COVIDEV/viq-chat-services/otp --namespace=com.github.romatroskin.viqchat.otp --type=service
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: com.github.romatroskin.viqchat.otp.service.otp
- Type: service
- Alias: otp

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
./otp-service
```

Build a docker image
```
make docker
```