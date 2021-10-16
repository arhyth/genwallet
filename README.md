# Genwallet

Genwallet is a generic wallet service with a REST API. It is implemented in Go, makes use of [go-kit](https://github.com/go-kit/kit) components when possible[^*], and follows DDD.

It is designed to be a stateless service so that it can be scaled to multiple instances without requiring some "control plane" overhead.

All account/wallet transactions in Genwallet are processed and kept in a postgreSQL backed repository. Concurrent account processes are guaranteed equivalent to some serial order with use of `Serializable` isolation level.

## Getting Started

#### Testing
We make use of the standard library `testing` package as well as some small 3rd party helper packages such as [testify](https://github.com/stretchr/testify).
Tests that require setup of dependencies must be kept separate to default `go test` call by using Go build tags. This allows anyone new to the project to contribute and add corresponding unit tests without burdening them to setup dependencies for unrelated tests. This establishes bias against complexity.

[^*]: Personally, although I agree with the design ideals of go-kit, I am not a fan of its much use of the empty `interface{}`. I think this is caused by requiring a certain structure to a very wide use case (microservices) while wanting to keep user codebase/s DRY.
