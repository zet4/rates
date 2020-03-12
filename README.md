# Rates

Rates is a microservice implemented per spec as a test project for Tet, this service is split in two parts, first is the API server that simply fetches data from database and second is the collect command that fetches currency data from RSS feed and updates the database.

As this project was implemented per spec, it does not use redis for caching queries to database, real usage would see that changed.

This is also a very simple service and thus the project layout was kept simple, more complex services would see a better laid out package structure.

## Setup

TL;DR;
```
docker-compose up
cat ./schema.sql | docker exec -i rates_db_1 sh -c "exec mysql -urates -prates rates"
docker exec -i rates_rates_service_1 sh -c "exec /rates collect"
```

Start the project by running `docker-compose up`, and populate the database by running `cat ./schema.sql | docker exec -i rates_db_1 sh -c "exec mysql -urates -prates rates"`, after which either run `docker exec -i rates_rates_service_1 sh -c "exec /rates collect"` to collect data from RSS or alternatively wait for cron task to execute (in theory)

The service should now be running and available on `localhost:8080`.

## Endpoints

Per spec, the service currently exposes 2 endpoints,
* `/api/v1/rates/latest` for up-to-date currency rates
* `/api/v1/rates/history/:currency` for historical data points for given currency.

Examples with service running:
* `http://localhost:8080/api/v1/rates/latest`
* `http://localhost:8080/api/v1/rates/history/gbp`

## Commands

Per spec, the service currently exposes 2 commands,
* `rates serve` starts the REST API server, by default on `:3333`
* `rates collect` fetches data from defined source and updates database
For more details and options check `rates -h` and `rates collect -h`