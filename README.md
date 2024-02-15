# Fruits API

A simple Fruits REST API built in `golang` using Labstack's [Echo](https://https://echo.labstack.com/]).

- For RDBMS(PostgreSQL,MySQL/MariaDB) demo use the [main](../../tree/main) branch
- NoSQL(__MongoDB__) please switch to [mongodb](../../tree/mongodb) branch.

| ℹ️ Note |
|---------|
| In addition to the application itself, this project can be used to demonstrate [AIDA](https://www.harness.io/products/aida) and [Remote Debug](https://developer.harness.io/docs/continuous-integration/troubleshoot-ci/debug-mode/) with [Harness CI](https://www.harness.io/products/continuous-integration). Follow [these instructions](.harness/README.md) to experiment with these features in your Harness account. |

## Pre-requisites

- [Docker Desktop](https://docs.docker.com/desktop/)
- [kubectl](https://kubernetes.io/docs/tasks/tools)
- [Drone CI CLI](https://docs.drone.io/cli/install/)

## Environment Setup

Copy the `.env.example` to `.env` and update the following variables to suit your settings.

- `PLUGIN_REGISTRY` - the docker registry to use
- `PLUGIN_TAG`      - the tag to push the image to docker registry
- `PLUGIN_REPO`     - the docker registry repository
- `PLUGIN_USERNAME` - the docker Registry username
- `PLUGIN_PASSWORD` - the docker registry password

### Backend Database to use

- `FRUITS_DB_TYPE` - the database to use with fruits api, defaults: `sqlite`

### Postgresql DB Settings

- `POSTGRES_HOST` - the postgresql host usually the docker or kubernetes service name e.g. `postgresql`
- `POSTGRES_PORT` - the postgresql port e.g. `5432`
- `POSTGRES_USER` - the postgresql user e.g. `demo`
- `POSTGRES_PASSWORD` - the postgresql password e.g `pa55Word!`
- `POSTGRES_DB` - the postgresql database to use e.g `demodb`

### MariaDB/MySQL Settings

- `MYSQL_HOST` - the MySQL host usually the docker or kubernetes service name e.g.`mysql`
- `MYSQL_PORT` - the MySQL port e.g. `3306`
- `MYSQL_ROOT_PASSWORD` - the MySQL root password `superS3cret!`
- `MYSQL_PASSWORD` - the MySQL password `pa55Word!`
- `MYSQL_USER` - the MySQL user e.g `demo`
- `MYSQL_DATABASE` - the MySQL database to use e.g `demodb`

### SQLite

- `FRUITS_DB_FILE` - the default database file to use.
  
>__NOTE:__
>
> - Most of the above database settings comes from how you setup the datbase. Please update accordingly
> - If you use FRUITS_DB_FILE  to use for testing.

## Build the Application

Set the `FRUIT_DB_TYPE` to `pgsql` or `mysql` to run tests against those databases. As by default all the tests are performed against `SQLite`.

```shell
drone exec --trusted --env-file=.env
```

The command will test, build and push the container image to the `$PLUGIN_REPO:$PLUGIN_TAG`.

## Run Application

### Locally

```shell
docker-compose up
```

## Testing

The application provides a [Swagger UI](http://localhost:8080/swagger/index.html) that can be used to used to play with the API.