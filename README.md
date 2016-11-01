## Installation

Prerequisites:

* A working [Go](https://golang.org/) development environment

```
go get -u github.com/teddyking/ladybug
cd $GOPATH/src/github.com/teddyking/ladybug
go install
```

## Usage

```
ladybug info
```

## Testing and CI

ladybug uses [concourse](https://concourse.ci) for both testing and CI.
Assuming you are running concourse-lite, and that you have [garden-runc-release](https://github.com/cloudfoundry/garden-runc-release) checked out to `"$HOME/workspace/garden-runc-release"`,
the test suite can be run via `scripts/test`.
