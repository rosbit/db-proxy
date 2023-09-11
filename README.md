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
q-len: 5         # the limit parallellism to access the real db server.
base32-chars: "" # 32 charactors used as base32 base.
dsn-params:
  js-file: name2dsn.js
  js-func: getDSNByName
common-endpoints:
  health: /health
  websocket: /websocket
```

There is a also a sample JavaScript file named [name2dsn.js](name2dsn.js), which can provide dsn-name to real dsn mapping dynamically. The sample content of `name2dsn.js`:
```javascript
var name2dsn = {
	'test': 'mysql:user:password@tcp(172.16.10.240:3306)/?charset=utf8mb4',
	'163': 'mysql:user:password@tcp(192.168.0.163:3306)/?charset=utf8mb4'
}

function getDSNByName(name) {
	var dsn = name2dsn[name];
	if (dsn) {
		return dsn
	}
	return ""
}
```

With the helper of `name2dsn.js`, you can specify dsn like `dbproxy:name-of-dsn` if you want to access a real database instance.

## Run

type `CONF_FILE=./db-proxy.sample.yaml ./db-proxy` to run `db-proxy`

## Status

`db-proxy` is fully tested with `db-proxy-driver`.
