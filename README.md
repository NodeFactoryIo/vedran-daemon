# Vedran daemon

> Daemon application for interacting with vedran load balancer.

### Get `vedran-daemon` package
1. Install [Golang](https://golang.org/doc/install) **1.13 or greater**
2. Run the command below
```
go get github.com/NodeFactoryIo/vedran-daemon
```
3. Run vedran-daemon from your Go bin directory. For linux systems it will likely be:
```
~/go/bin/vedran-daemon
```
Note that if you need to do this, you probably want to add your Go bin directory to your $PATH to make things easier!

## Usage

```
$ ./vedran-daemon -h
Register vedran-daemon with load balancer and start sending telemetry

Usage:
  vedran-daemon [flags]

Flags:
  -h, --help                    help for vedran-daemon
      --id string               Vedran-daemon id string (required)
      --lb string               Target load balancer url (required)
      --log-file string         Path to logfile. If not set defaults to stdout
      --log-level string        Level of logging (eg. debug, info, warn, error) (default "info")
      --node-metrics string     Polkadot node metrics url (default "localhost:9615")
      --node-rpc string         Polkadot node rpc url (default "localhost:9933")
      --payout-address string   Payout address to which reward tokens will be sent (required)
```
## Starting daemon

First download latest prebuilt binary of vedran daemon [from releases](https://github.com/NodeFactoryIo/vedran-daemon/releases) and binary for [node](https://github.com/paritytech/polkadot/releases).
Daemon is expected to be started in conjuction with node and will wait for node if it is unavailable.


### Node
For starting node see [instructions](https://github.com/paritytech/polkadot/blob/master/README.md)

**NOTE node should be started with rpc cors disabled

Example:
`
  ./polkadot --rpc-cors=all
`

### Daemon
Daemon is started by invoking binary.

For example:
```
  ./vedran-daemon-linux-amd64 --id UuPrCMnkni --lb https://load-balancer.com --payout-address 15MCkjt3B59dNo5reMCWWpxY8QB8VpEbYLo2xHEjuuWsSmTU
```

It will register to load balancer and start sending pings and metrics to load balancer.
Port forwarding to local node is not needed and node and daemon can be in a private network because [http tunnel](https://en.wikipedia.org/wiki/HTTP_tunnel) is created between
node and load balancer on registration which communicate via daemon used as a proxy server.

### Required flags

`--id` - id string by which load balancer will distinguish between nodes - **CAUTION** this should be a unique string and should not be shared

`--lb` - public url of vedran load balancer

`--payout-address` - address of wallet to which reward tokens should be set - **CAUTION** - use valid address depending on network

### Other flags

`--log-level` - log level (debug, info, warn, error) - **DEFAULT** [error]

`--log-file` - path to file in which logs will be saved - **DEFAULT** [stdout]

`--node-metrics` - local url to node metrics - **DEFAULT** [http://localhost:9615]

`--node-rpc` - local url to node rpc endpoint - **DEFAULT** [http://localhost:9933]

## Development


### Clone

```bash
git clone git@github.com:NodeFactoryIo/vedran-daemon.git
```

### Lint
[Golangci-lint](https://golangci-lint.run/usage/install/#local-installation) is expected to be installed.

```bash
make lint
```

### Build

```bash
make build
```

### Test

```bash
make test
```

Run daemon app with `go run main.go [command]`.

More about different _commands_ can be found in [Usage](#Usage).

## License

This project is licensed under Apache 2.0:
- Apache License, Version 2.0, [LICENSE](LICENSE) or http://www.apache.org/licenses/LICENSE-2.0
