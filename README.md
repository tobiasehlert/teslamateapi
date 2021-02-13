# TeslaMateApi

TeslaMateApi is a RESTful API to get data collected by self-hosted data logger **[TeslaMate](https://github.com/adriankumpf/teslamate)** in JSON.

- Written in **[Golang](https://golang.org/)**
- Data is collected from TeslaMate **Postgres** database and local **MQTT** Broker
- Endpoints return data in JSON format

## How to use

*Information will be updated..*

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
