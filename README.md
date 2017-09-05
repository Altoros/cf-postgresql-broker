# cf-postgresql-broker

PostgreSQL single-database-instance service broker for the Cloud Foundry.

## Installation

Clone the git repository and compile the package:
```
$ git clone https://github.com/Altoros/cf-postgresql-broker.git
$ cd cf-postgresql-broker
```


Push the compiled binary to the Cloud Foundry:
```
$ cf push --no-start
```

Set up the environment:
```
$ BASIC_AUTH_USERNAME=admin
$ BASIC_AUTH_PASSWORD=admin
$ PG_SOURCE=postgresql://username:password@host:port/dbname
$ PG_SERVICES='[{
  "id": "service-postgresql-{GUID}",
  "name": "postgres",
  "description": "DBaaS",
  "bindable": true,
  "plan_updateable": false,
  "plans": [{
    "id": "plan-basic-{GUID}",
    "name": "basic",
    "description": "Single database plan"
  }]
}]'

$ cf set-env cf-postgresql-broker BASIC_AUTH_USERNAME "$BASIC_AUTH_USERNAME"
$ cf set-env cf-postgresql-broker BASIC_AUTH_PASSWORD "$BASIC_AUTH_PASSWORD"
$ cf set-env cf-postgresql-broker PG_SOURCE "$PG_SOURCE"
$ cf set-env cf-postgresql-broker PG_SERVICES "$PG_SERVICES"
```

* `PG_SERVICES` can be customized according to [this go library](https://github.com/pivotal-cf/brokerapi/blob/master/catalog.go#L3)
* `{GUID}` will be replaced with its runtime value

Optional:  If your PostgreSQL server does not use SSL, disable it:
```
$ cf set-env cf-postgresql-broker PGSSLMODE "disable"
```

Start the application and register a service broker:
```$AUTH_PASSWORD
$ cf start cf-postgresql-broker
$ BROKER_URL=$(cf app cf-postgresql-broker | grep urls: | awk '{print $2}')
$ cf create-service-broker cf-postgresql-broker $BASIC_AUTH_USERNAME $BASIC_AUTH_PASSWORD http://$BROKER_URL
$ cf enable-service-access postgres
```

Then you should see `cf-postgresql-broker` in `$ cf marketplace`

## Binding applications

```
$ cf create-service postgres basic my-psql-db
$ cf bind-service my-app my-psql-db
$ cf restage my-app
```
