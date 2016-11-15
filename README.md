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
# ladybug info

Running containers: 1

# ladybug containers

Handle        IP Address  Process Name  Created At                      Port Mappings
my-container  10.254.0.6  ruby          2016-11-14 13:48:57  60000->80, 60001->81
```

## Testing and CI

ladybug uses [concourse](https://concourse.ci) for both testing and CI.
Assuming you are running concourse-lite, and that you have [garden-runc-release](https://github.com/cloudfoundry/garden-runc-release) checked out to `"$HOME/workspace/garden-runc-release"`,
the test suite can be run via `scripts/test`.
