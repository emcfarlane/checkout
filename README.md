Checkout
===

Checkout implements a simple authorization flow service.

- gRPC service with HTTP bindings by [graphpb](https://github.com/emcfarlane/graphpb)
- Postgres service with library from [go cloud cdk](https://github.com/google/go-cloud)

Build
---

Build with go:
```
go install ./checkout
```

### Protobuf (optional)

Protobuffers can be rebuilt with protoc:
```sh
protoc -I checkoutpb --go_out=paths=source_relative:checkoutpb --go-grpc_out=paths=source_relative:checkoutpb checkoutpb/checkout.proto
```
NB: Requires both go & grpc protobuf generates installed locally as well as protoc with google apis included.

Run
---

Run must specify a DATABASE address.
```
POSTGRES=postgres://<user>:<password>@localhost/<database> checkout
```
The database URL pattern is described [here](https://gocloud.dev/howto/sql/).

Test
---

Unit tests can be run:
```
go test -v .
```

Integration tests require a postgresql database connection:
```
POSTGRES=postgres://<user>:<password>@localhost/<database> go test -v -tags=integration .
```

Examples
---

Create a new authorization:
```
curl -X POST -H "Content-Type: application/json" -d '{"pan":"4532111111111112","exp_month":"1","exp_year":"2022","cvv":"123","
amount":500,"currency":"gbp"}' localhost:8080/authorize
```

Capture an authorization:
```
curl -X PATCH -H "Content-Type: application/json" -d '{"id":"<id>","amount":250}' localhost:8080/capture
```

Refund an authorization:
```
curl -X PATCH -H "Content-Type: application/json" -d '{"id":"<id>","amount":250}' localhost:8080/refund
```

Void an authorization:
```
curl -X POST -H "Content-Type: application/json" -d '{"id":"<id>"}' localhost:8080/void
```
