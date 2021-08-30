# TeslaMateApi

[![GitHub CI](https://github.com/tobiasehlert/teslamateapi/workflows/build/badge.svg?branch=main)](https://github.com/tobiasehlert/teslamateapi/actions?query=workflow%3Abuild)
[![GitHub go.mod version](https://img.shields.io/github/go-mod/go-version/tobiasehlert/teslamateapi)](https://github.com/tobiasehlert/teslamateapi/blob/main/go.mod)
[![Docker version](https://img.shields.io/docker/v/tobiasehlert/teslamateapi/latest)](https://hub.docker.com/r/tobiasehlert/teslamateapi)
[![Docker size](https://img.shields.io/docker/image-size/tobiasehlert/teslamateapi/latest)](https://hub.docker.com/r/tobiasehlert/teslamateapi)
[![GitHub license](https://img.shields.io/github/license/tobiasehlert/teslamateapi)](https://github.com/tobiasehlert/teslamateapi/blob/main/LICENSE)
[![Docker pulls](https://img.shields.io/docker/pulls/tobiasehlert/teslamateapi)](https://hub.docker.com/r/tobiasehlert/teslamateapi)

TeslaMateApi is a RESTful API to get data collected by self-hosted data logger **[TeslaMate](https://github.com/adriankumpf/teslamate)** in JSON.

- Written in **[Golang](https://golang.org/)**
- Data is collected from TeslaMate **Postgres** database and local **MQTT** Broker
- Endpoints return data in JSON format
- Send commands to your Tesla through the TeslaMateApi

### Table of Contents

- [How to use](#how-to-use)
  - [Docker-compose](#docker-compose)
  - [Environment variables](#environment-variables)
- [API documentation](#api-documentation)
  - [Available endpoints](#available-endpoints)
  - [Authentication](#authentication)
  - [Commands](#commands)
- [Security information](#security-information)
- [Credits](#credits)

## How to use

You can either use it in a Docker container or go download the code and deploy it yourself on any server.

### Docker-compose

If you run the simple Docker deployment of TeslaMate, then adding this will do the trick. You'll have TeslaMateApi exposed at port 8080 locally then.

```yaml
services:
  teslamateapi:
    image: tobiasehlert/teslamateapi:latest
    restart: always
    depends_on:
      - database
    environment:
      - DATABASE_USER=teslamate
      - DATABASE_PASS=secret
      - DATABASE_NAME=teslamate
      - DATABASE_HOST=database
      - MQTT_HOST=mosquitto
      - TZ=Europe/Berlin
    ports:
      - 8080:8080
```

If you are using TeslaMate Traefik setup in Docker with environment variables file (.env), then you can simply add this section to the `services:` section of the `docker-compose.yml` file:

```yaml
services:
  teslamateapi:
    image: tobiasehlert/teslamateapi:latest
    restart: always
    depends_on:
      - database
    environment:
      - DATABASE_USER=${TM_DB_USER}
      - DATABASE_PASS=${TM_DB_PASS}
      - DATABASE_NAME=${TM_DB_NAME}
      - DATABASE_HOST=database
      - MQTT_HOST=mosquitto
      - TZ=${TM_TZ}
    labels:
      - "traefik.enable=true"
      - "traefik.port=8080"
      - "traefik.http.middlewares.redirect.redirectscheme.scheme=https"
      - "traefik.http.middlewares.teslamateapi-auth.basicauth.realm=teslamateapi"
      - "traefik.http.middlewares.teslamateapi-auth.basicauth.usersfile=/auth/.htpasswd"
      - "traefik.http.routers.teslamateapi-insecure.rule=Host(`${FQDN_TM}`)"
      - "traefik.http.routers.teslamateapi-insecure.middlewares=redirect"
      - "traefik.http.routers.teslamateapi.rule=Path(`/api`) || PathPrefix(`/api/`)"
      - "traefik.http.routers.teslamateapi.entrypoints=websecure"
      - "traefik.http.routers.teslamateapi.middlewares=teslamateapi-auth"
      - "traefik.http.routers.teslamateapi.tls.certresolver=tmhttpchallenge"
```

In this case, the TeslaMateApi would be accessible at teslamate.example.com/api/

### Environment variables

Basically the same environment variables for the database, mqqt and timezone need to be set for TeslaMateApi as you have for TeslaMate.

**Required** environment variables (even if there are some default values available)

- **DATABASE_USER** string *(default: teslamate)*
- **DATABASE_PASS** string *(default: secret)*
- **DATABASE_NAME** string *(default: teslamate)*
- **DATABASE_HOST** string *(default: database)*
- **MQTT_HOST** string *(default: mosquitto)*
- **TZ** string *(default: Europe/Berlin)*

**Optional** environment variables

- **TESLAMATE_SSL** boolean *(default: false)*
- **TESLAMATE_HOST** string *(default: teslamate)*
- **TESLAMATE_PORT** string *(default: 4000)*
- **API_TOKEN** string *(default: )*
- **API_TOKEN_DISABLE** string *(default: false)*
- **DATABASE_PORT** integer *(default: 5432)*
- **DATABASE_TIMEOUT** integer *(default: 60000)*
- **DATABASE_SSL** boolean *(default: true)*
- **DEBUG_MODE** boolean *(default: false)*
- **DISABLE_MQTT** boolean *(default: false)*
- **MQTT_TLS** boolean *(default: false)*
- **MQTT_PORT** integer *(default: 1883 / default if TLS is true: 8883)*
- **MQTT_USERNAME** string *(default: )*
- **MQTT_PASSWORD** string *(default: )*
- **MQTT_NAMESPACE** string *(default: )*

**Commands** environment variables

- **ENABLE_COMMANDS** boolean *(default: false)*
- **COMMANDS_ALL** boolean *(default: false)*
- **COMMANDS_ALLOWLIST** string *(default: allow_list.json)*
- **COMMANDS_LOGGING** boolean *(deafault: false)*
- **COMMANDS_WAKE** boolean *(default: false)*
- **COMMANDS_ALERT** boolean *(default: false)*
- **COMMANDS_REMOTESTART** boolean *(default: false)*
- **COMMANDS_HOMELINK** boolean *(default: false)*
- **COMMANDS_SPEEDLIMIT** boolean *(default: false)*
- **COMMANDS_VALET** boolean *(default: false)*
- **COMMANDS_SENTRYMODE** boolean *(default: false)*
- **COMMANDS_DOORS** boolean *(default: false)*
- **COMMANDS_TRUNK** boolean *(default: false)*
- **COMMANDS_WINDOWS** boolean *(default: false)*
- **COMMANDS_SUNROOF** boolean *(default: false)*
- **COMMANDS_CHARGING** boolean *(default: false)*
- **COMMANDS_CLIMATE** boolean *(default: false)*
- **COMMANDS_MEDIA** boolean *(default: false)*
- **COMMANDS_SHARING** boolean *(default: false)*
- **COMMANDS_SOFTWAREUPDATE** boolean *(default: false)*

## API documentation

More detailed documentation of every endpoint will come..

### Available endpoints

- GET `/api`
- GET `/api/v1`
- GET `/api/v1/cars`
- GET `/api/v1/cars/:CarID`
- GET `/api/v1/cars/:CarID/charges`
- GET `/api/v1/cars/:CarID/charges/:ChargeID`
- GET `/api/v1/cars/:CarID/command`
- POST `/api/v1/cars/:CarID/command/:Command`
- GET `/api/v1/cars/:CarID/drives`
- GET `/api/v1/cars/:CarID/drives/:DriveID`
- PUT `/api/v1/cars/:CarID/logging/:Command`
- GET `/api/v1/cars/:CarID/logging`
- GET `/api/v1/cars/:CarID/status`
- GET `/api/v1/cars/:CarID/updates`
- POST `/api/v1/cars/:CarID/wake_up`
- GET `/api/v1/globalsettings`
- GET `/api/ping`

### Authentication

If you want to use command or logging endpoints such as `/api/v1/cars/:CarID/command/:Command`, `/api/v1/cars/:CarID/wake_up`, or `/api/v1/cars/:CarID/logging/:Command` you need to add authentication to your request.

You need to specify a token yourself (called **API_TOKEN**) in the environment variables file, to set it. The token has the requirement to be a minimum of 32 characters long.

There are two options available for authentication to be done.

1. Adding extra header `Authorization: Bearer <token>` to your request. (recommended option)

2. Adding URI parameter `?token=<token>` to the endpoint you try to reach. (not a good option)

\* *Note: If you use the second option and your logs get compromised, your token will be leaked.*

### Commands

Commands are not enabled by default.

You need to enable them in your environment variables (with `ENABLE_COMMANDS=true`) and you need to specify which commands you want to use as well.

There are 3 ways of using Commands:

1. Specific groups of commands can be enabled for example `COMMANDS_ALERT=true` will enable the [alert](https://tesla-api.timdorr.com/vehicle/commands/alerts) commands group.

2. If you need a granular set of commands enabled `COMMANDS_ALLOWLIST=/path/to/allow_list.json` can be used to specify a [JSON formatted list of commands](./example/allow_list.json) to enable.

3. The most coarse option `COMMANDS_ALL=true` will enable all commands (specific groups and allow_list will be ignored).

\* *Note: if `COMMANDS_ALL` or any specific group of commands has been enabled `COMMANDS_ALLOWLIST` is ignored.*

A list of possible commands can be found under [environment variables](#environment-variables).

Regarding what fields you need to provide in the commands, we will referr to the [timdorr/tesla-api](https://tesla-api.timdorr.com/vehicle/commands) documentation.

## Security information

There is **no** possibility to get access to your Tesla account tokens by this API and we'll keep it this way!

The data that is accessible is data like the cars, charges, drives, current status, updates and global settings.

Also, apply some authentication on your webserver in front of the container, so your data is not unprotected and too exposed. In the example above, we use the same .htpasswd file as used by TeslaMate.

If you have applied a level of authentication in front of the container `API_TOKEN_DISABLE=true` will allow commands without requiring the header or uri token value.

## Credits

- Authors: Tobias Lindberg – [List of contributors](https://github.com/tobiasehlert/teslamateapi/graphs/contributors)
- Distributed under MIT License
