# cf-postgresql-broker

PostgreSQL service broker for the cloud foundry.

## Installation

```bash
BROKER_NAME=cf-postgresql-broker
PLAN_NAME=basic
SERVICE_NAME=pgsql

AUTH_USER=admin
AUTH_PASSWORD=admin
PG_SOURCE=postgresql://username:password@host:port/dbname
PG_SERVICES="[{
  \"id\": \"service-1-{GUID}\",
  \"name\": \"$BROKER_NAME\",
  \"description\": \"DBaaS\",
  \"bindable\": true,
  \"plan_updateable\": false,
  \"plans\": [{
    \"id\": \"plan-1-{GUID}\",
    \"name\": \"$PLAN_NAME\",
    \"description\": \"Shared plan\"
  }]
}]"

# Deploy to Cloud Foundry
go get github.com/altoros/cf-postgresql-broker
cd $GOPATH/src/github.com/altoros/cf-postgresql-broker
cf push $BROKER_NAME --no-start -m 128M -k 256M
cf set-env $BROKER_NAME PG_SOURCE "$PG_SOURCE"
cf set-env $BROKER_NAME PG_SERVICES "$PG_SERVICES"
cf set-env $BROKER_NAME AUTH_USER "$AUTH_USER"
cf set-env $BROKER_NAME AUTH_PASSWORD "$AUTH_PASSWORD"
cf start $BROKER_NAME

# Register service broker
BROKER_URL=$(cf app $BROKER_NAME | grep urls: | awk '{print $2}')
cf create-service-broker $BROKER_NAME $AUTH_USER $AUTH_PASSWORD http://$BROKER_URL
cf enable-service-access $BROKER_NAME

# Bind an application
cf create-service $BROKER_NAME $PLAN_NAME $SERVICE_NAME
cf bind-service my-app $SERVICE_NAME
cf restage
```

## Development

1. Copy `.envrc.example` to `.envrc`, then load it by `$ source .envrc` if you don't have the **direnv** package installed.

2. `$ godep restore`
