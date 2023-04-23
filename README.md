# Database Proxy

`db-proxy` is database server to share db connections with the same connection parameters,
and limit the parallellism of access to the real database server. `db-proxy` can be accessed with protocols
named HTTP, net/rpc, or net/rpc over websocket.

There is a corresponding golang db driver named [db-proxy-driver](https://github.com/rosbit/db-proxy-driver)
which can be used to access `db-proxy`. If you want to make use of the power of `db-proxy`,
what you can do is just replace your driver with `db-proxy-driver`, and paste the real driver name aheader of
your origin DSN string with colon as the delimeter. e.g., `mysql:user:password@tcp(ip:port)/db?charset=utf8mb4`.

Currently, only mysql driver is included in `db-proxy`. Other driver can be included easily
if necessary.

## Install

run `go get github.com/rosbit/db-proxy`, an executable `db-proxy` should be found.

There is a sample configration file named [db-proxy.sample.yaml](db-proxy.sample.yaml), which is
needed to run `db-proxy`.
```yaml
---
listen-host: ""
http-listen-port: 7080
jsonl-rpc-listen-prot: 7081
q-len: 5 # the limit parallellism to access the real db server.
```

## Run

type `CONF_FILE=./db-proxy.sample.yaml ./db-proxy` to run `db-proxy`

## Status

`db-proxy` is fully tested with `db-proxy-driver`.
