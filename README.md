# cf-postgresql-broker

PostgreSQL service broker for the Cloud Foundry.

## Requirements

* [govendor](https://github.com/kardianos/govendor)
* [go-buildpack](https://github.com/cloudfoundry/go-buildpack) >= 1.7.7

## Installation

Download the source code and install dependencies.

```
$ go get github.com/altoros/cf-postgresql-broker
$ cd $GOPATH/src/github.com/altoros/cf-postgresql-broker
$ govendor sync
```

Push the code base to the Cloud Foundry.

```
$ cf push postgresql --no-start -m 128M -k 256M [-b https://github.com/cloudfoundry/go-buildpack#v1.7.7]
```

Set the environment.

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

$ cf set-env postgresql AUTH_USER "$AUTH_USER"
$ cf set-env postgresql AUTH_PASSWORD "$AUTH_PASSWORD"
$ cf set-env postgresql PG_SOURCE "$PG_SOURCE"
$ cf set-env postgresql PG_SERVICES "$PG_SERVICES"
```

Start the application and register a service broker

```
$ cf start postgresql
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
