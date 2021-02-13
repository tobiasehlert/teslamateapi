# TeslaMateApi

TeslaMateApi is a RESTful API to get data collected by self-hosted data logger **[TeslaMate](https://github.com/adriankumpf/teslamate)** in JSON.

- Written in **[Golang](https://golang.org/)**
- Data is collected from TeslaMate **Postgres** database and local **MQTT** Broker
- Endpoints return data in JSON format

## How to use

You can either use it in a Docker container or go download the code and deploy it yourself on any server.

If you are using TeslaMate as Docker with environment variables file (.env), then you can simply add this section to the `services:` section of the `docker-compose.yml` file:

```
services:
  teslamateapi:
    image: tobiasehlert/teslamateapi:latest
    restart: always
    environment:
      - DATABASE_USER=${TM_DB_USER}
      - DATABASE_PASS=${TM_DB_PASS}
      - DATABASE_NAME=${TM_DB_NAME}
      - DATABASE_HOST=database
      - MQTT_HOST=mosquitto
      - TZ=${TM_TZ}
```

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

- **DATABASE_PORT** integer *(default: 5432)*
- **DATABASE_TIMEOUT** integer *(default: 60000)*
- **DATABASE_SSL** boolean *(default: true)*
- **DISABLE_MQTT** boolean *(default: false)*
- **MQTT_TLS** boolean *(default: false)*
- **MQTT_PORT** integer *(default: 1883 / default if TLS is true: 8883)*
- **MQTT_USERNAME** string *(default: )*
- **MQTT_PASSWORD** string *(default: )*
- **MQTT_NAMESPACE** string *(default: )*

## API documentation

More detailed documentation of every endpoint will come..

- GET `/cars`
- GET `/cars/:CarID`
- GET `/cars/:CarID/charges`
- GET `/cars/:CarID/charges/:ChargeID`
- GET `/cars/:CarID/drives`
- GET `/cars/:CarID/drives/:DriveID`
- GET `/cars/:CarID/status`
- GET `/cars/:CarID/updates`
- GET `/globalsettings`
- GET `/ping`

## Security information

There is **no** possibility to get access to your Tesla account tokens by this API and we'll keep it this way!

The data that is accessible is data like the cars, charges, drives, current status, updates and global settings.

## Credits

- Authors: Tobias Lindberg – [List of contributors](https://github.com/tobiasehlert/teslamateapi/graphs/contributors)
- Distributed under MIT License
