## Installation

Prerequisites:

* A working [Go](https://golang.org/) development environment

```
go get -u github.com/teddyking/ladybug
cd $GOPATH/src/github.com/teddyking/ladybug
go install
```

## Testing and CI

ladybug uses [concourse](https://concourse.ci) for both testing and CI.
Assuming you have concourse installed with a [fly target](https://concourse.ci/fly-targets.html)
named `lite`, the test suite can be run as follows:

```
fly -t lite e -c ci/test.yml -x -i ladybug-src=.
```

And the ladybug pipeline can be configured by running:

```
fly -t lite sp -p ladybug -c ci/pipeline.yml -l ci/secrets.yml
```

Alternatively, simply use the `scripts/test` and `scripts/set-pipeline` scripts.
