# TeslaMateApi

[![GitHub CI](https://github.com/tobiasehlert/teslamateapi/workflows/build/badge.svg?branch=main)](https://github.com/tobiasehlert/teslamateapi/actions?query=workflow%3Abuild)
![GitHub go.mod version](https://img.shields.io/github/go-mod/go-version/tobiasehlert/teslamateapi)
[![Docker version](https://img.shields.io/docker/v/tobiasehlert/teslamateapi/latest)](https://microbadger.com/images/tobiasehlert/teslamateapi)
[![Docker size](https://img.shields.io/docker/image-size/tobiasehlert/teslamateapi/latest)](https://microbadger.com/images/tobiasehlert/teslamateapi)
[![Docker layers](https://img.shields.io/microbadger/layers/tobiasehlert/teslamateapi/latest)](https://microbadger.com/images/tobiasehlert/teslamateapi)
[![GitHub license](https://img.shields.io/github/license/tobiasehlert/teslamateapi)](https://microbadger.com/images/tobiasehlert/teslamateapi)
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

- **API_TOKEN** string *(default: )*
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
- **MQTT_SLEEPTIME** integer *(default: 100)*

**Commands** environment variables

- **ENABLE_COMMANDS** boolean *(default: false)*
- **COMMANDS_ALL** boolean *(default: false)*
- **COMMANDS_ALLOWLIST** string *(default: allow_list.json)*
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
- GET `/api/v1/cars/:CarID/status`
- GET `/api/v1/cars/:CarID/updates`
- POST `/api/v1/cars/:CarID/wake_up`
- GET `/api/v1/globalsettings`
- GET `/api/ping`

## Security information

There is **no** possibility to get access to your Tesla account tokens by this API and we'll keep it this way!

The data that is accessible is data like the cars, charges, drives, current status, updates and global settings.

Also, apply some authentication on your webserver in front of the container, so your data is not unprotected and too exposed. In the example above, we use the same .htpasswd file as used by TeslaMate.

## Credits

- Authors: Tobias Lindberg – [List of contributors](https://github.com/tobiasehlert/teslamateapi/graphs/contributors)
- Distributed under MIT License
