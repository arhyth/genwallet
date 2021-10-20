Genwallet
===

Genwallet is a generic wallet service with a REST API. It is implemented in Go, makes use of [go-kit](https://github.com/go-kit/kit) components when possible[^*], and follows DDD.

Genwallet codebase has minimal number of concepts/mechanisms. This is not a requirement or objective but we try to introduce as little indirection as possible in order to achieve the least amount of coupling.

All of its endpoints offer a synchronous API including payment transactions. This design choice provides predictability to the user. This is both a pro and a con. In a sync system, the user immediately knows if the system is slow or when it encounters an error. But in an otherwise async system, initial interactions such as submitting a payment request will almost always succeed but as the system hits a bottleneck somewhere, the lack of backpressure can "bury" the system into a failure loop.

Genwallet is also designed to be a stateless service so that it can be scaled to multiple instances without the overhead of some "control plane". All account/wallet transactions in Genwallet are handled by a postgreSQL database. Concurrent account processes are guaranteed equivalent to some serial order with use of `Serializable` isolation level. There is some performance penalty incurred for this as concurrent transactions targeting similar row/s will fail except for the succeeding one. For simplicity, it is left to the API user to retry the request. This also serves as a feedback mechanism.

Roadmap
---
- [x] Design and documentation
- [x] Skeleton and fake service implementation
- [x] Request validation
- [x] Error handling
- [x] Unit tests (needs a lot more work tbh)
- [x] Service implementation (postgres backed)
- [ ] Integration tests
- [ ] Dev setup conveniences (docker-compose etc)

How it works
---
Genwallet works by recording to a single "global" ledger of `transfer`s between wallet `account`s. `Payment`s are simply representations of `transfer` with respect to an `account`. In a production service, wallet `account`s could benefit from partitioning by currency but for simplicity, Genwallet uses only a single table.
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
- Run migrations if applicable. Migrations are written for [`goose`](https://github.com/pressly/goose)
```sh
$ goose -dir db/migrations postgres <DB_URL> up
```
- Run `./gw-bin` with your set env vars

### Testing

We make use of the standard library `testing` package as well as some small 3rd party helper packages such as [testify](https://github.com/stretchr/testify).
Tests that require setup of dependencies must be kept separate to default `go test` call by using Go build tags. This establishes some bias against complexity and allows anyone new to the project to contribute and add corresponding unit tests without burdening them to setup dependencies for unrelated tests.

[^*]: Personally, although I agree with the design ideals of go-kit, I am not a fan of its much use of the empty `interface{}`. I think this is caused by requiring a certain structure to a very wide use case (microservices) while wanting to keep user codebase/s DRY.
