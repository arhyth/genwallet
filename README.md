Genwallet
===

Genwallet is a generic wallet service with a REST API. It is implemented in Go, makes use of [go-kit](https://github.com/go-kit/kit) components when possible[^*], and follows DDD.

It is designed to be a stateless service so that it can be scaled to multiple instances without requiring some "control plane" overhead.

All account/wallet transactions in Genwallet are processed and kept in a postgreSQL backed repository. Concurrent account processes are guaranteed equivalent to some serial order with use of `Serializable` isolation level. There is some performance penalty incurred for this as concurrent transactions targeting similar row/s will fail except for the succeeding one. For simplicity, it is left to the API user to retry the request.

Roadmap
---
- [x] Design and documentation
- [x] Skeleton and fake service implementation
- [ ] Request validation
- [ ] Error handling
- [ ] Unit tests
- [ ] Service implementation (postgres backed)
- [ ] Integration tests
- [ ] Dev setup conveniences (docker-compose etc)

How it works
---
Genwallet works by recording to a single "global" ledger of `transfer`s between wallet `account`s. `Payment`s are simply representations of `transfer` with respect to an `account`.
![how Genwallet works](./DB-entities.svg)

API
---
Details of request and response structure for each endpoint are listed in [separate document](API.md).

### Summary

| Method | Path | Description |
| :--- | :---: | :---: |
| `GET` | `/wallets` | list all wallets |
| `POST` | `/wallets` | create wallet |
| `GET` | `/wallets/{id}` | show wallet |
| `GET` | `/wallets/{id}/payments` | list all transfers from/to wallet |
| `POST` | `/wallets/{id}/payments` | make transfer from one wallet to another |
| `GET` | `/transfers` | list all transfers |

Getting Started
---
### Config
- **ADDR_PORT** : `address:port` where service listens (defaults to `:8000`)
- **DB_URL** (required) : postgres database connection string

### Development
**Setup**
- Build `go build -mod=vendor -o gw-bin cmd/main.go`
- Set environment variables as listed in section [`Config`](#config)
- Run migrations if applicable
```sh
$ goose -dir migrations postgres <DB_URL> up
```
- Run `./gw-bin` with your set env vars

### Testing
We make use of the standard library `testing` package as well as some small 3rd party helper packages such as [testify](https://github.com/stretchr/testify).
Tests that require setup of dependencies must be kept separate to default `go test` call by using Go build tags. This allows anyone new to the project to contribute and add corresponding unit tests without burdening them to setup dependencies for unrelated tests. This establishes bias against complexity.

[^*]: Personally, although I agree with the design ideals of go-kit, I am not a fan of its much use of the empty `interface{}`. I think this is caused by requiring a certain structure to a very wide use case (microservices) while wanting to keep user codebase/s DRY.
