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
$ vedran-daemon -h

vedran-daemon is a command line interface for ....

Usage:
  vedran-daemon [command]

Available Commands:
  version     Show the current version of Vedran deamon app

Use "vedran-daemon [command] --help" for more information about a command.
```

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
