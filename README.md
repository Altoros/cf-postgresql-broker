# cf-postgresql-broker

PostgreSQL service broker for the Cloud Foundry (Diego compatible).

## Installation

Push the code base to the Cloud Foundry.

```
$ go get github.com/altoros/cf-postgresql-broker
$ cd $GOPATH/src/github.com/altoros/cf-postgresql-broker
$ cf push postgresql --no-start -m 128M -k 256M
```

Set environment.

* `PG_SERVICES` can be customized according to [this go library](https://github.com/pivotal-cf/brokerapi/blob/master/catalog.go#L3)
* `{GUID}` for `id` attributes will be replaced with its runtime value

```
$ AUTH_USER=admin
$ AUTH_PASSWORD=admin
$ PG_SOURCE=postgresql://username:password@host:port/dbname
$ PG_SERVICES='[{
  "id": "service-1-{GUID}",
  "name": "postgresql",
  "description": "DBaaS",
  "bindable": true,
  "plan_updateable": false,
  "plans": [{
    "id": "plan-1-{GUID}",
    "name": "basic",
    "description": "Shared plan"
  }]
}]'

$ cf set-env postgresql PG_SOURCE "$PG_SOURCE"
$ cf set-env postgresql PG_SERVICES "$PG_SERVICES"
$ cf set-env postgresql AUTH_USER "$AUTH_USER"
$ cf set-env postgresql AUTH_PASSWORD "$AUTH_PASSWORD"
```

Start the broker.

```
$ cf start postgresql
```

Register a service broker.

```
$ BROKER_URL=$(cf app postgresql | grep urls: | awk '{print $2}')
$ cf create-service-broker postgresql $AUTH_USER $AUTH_PASSWORD http://$BROKER_URL
$ cf enable-service-access postgresql
```

Then you should see `postgresql` in `$ cf marketplace`

## Binding applications

```
$ cf create-service postgresql basic pgsql
$ cf bind-service my-app pgsql
$ cf restage my-app
```

## Development

1. Copy `.envrc.example` to `.envrc`, then load it by `$ source .envrc` if you don't have the [direnv](http://direnv.net) package installed.

2. Install [godep](https://github.com/tools/godep) and fire `$ godep restore`
