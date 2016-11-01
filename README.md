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
Assuming you have concourse installed with a [fly target](https://concourse.ci/fly-targets.html)
named `lite`, the test suite can be run with a simple `scripts/test`.
