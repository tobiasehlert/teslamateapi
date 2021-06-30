# Changelog

## [Unreleased]

## [1.10.0] - 2021-06-30

### Added
- adding power mqtt value to status endpoint (#74)

### Changed
- improved log message of convert-functions (#75 by @alecdoconnor)
- bump various workflow versions (#78, #82, #83 by dependabot)
- minor code adjustments based on go-staticcheck

## [1.9.0] - 2021-06-14

### Added
- option to disable auth token for commands (#71 by @LelandSindt)

### Changed
- bump golang from 1.16.4 to 1.16.5 (#72 by dependabot)
- bump github.com/eclipse/paho.mqtt.golang from 1.3.4 to 1.3.5 (#73 by dependabot)

## [1.8.0] - 2021-06-01

### Fixed
- fixing sql error when EndDate is null (#58 and #69 by @alecdoconnor)

## [1.7.1] - 2021-06-01

### Changed
- bump github.com/gin-gonic/gin from 1.7.1 to 1.7.2 (#63 by dependabot)
- bump various workflow versions (#64, #65, #66, #67 by dependabot)

## [1.7.0] - 2021-05-20

### Added
- feature to resume/suspend logging of TeslaMate through TeslaMateApi (#34, #45 and #48 by @LelandSindt)

### Changed
- minor code adjustments based on go-staticcheck

## [1.6.2] - 2021-05-19

### Changed
- bump github.com/lib/pq from 1.10.1 to 1.10.2 (#59 by dependabot)
- bump various workflow versions (#52, #53, #54, #55, #56 and #57 by dependabot)

## [1.6.1] - 2021-05-10

### Changed
- bump golang from 1.16.3 to 1.16.4 (#49 by dependabot)
- bump github.com/eclipse/paho.mqtt.golang from 1.3.3 to 1.3.4 (#46 by dependabot)

## [1.6.0] - 2021-05-03

### Added
- doing persistant mqtt connection for status collection (#16 by @LelandSindt and #21 by @MattBrittan)
- adding randomized string to mqtt client (#15 by @LelandSindt)

### Changed
- fixing sql error when FastChargerBrand is null (#39 by @alecdoconnor)
- updating workflow for stale issues/PRs

### Removed
- removing MQTT_SLEEPTIME option for mqtt connection
- removing workflow for no response

## [1.5.0] - 2021-05-03

### Changed
- missing convertion of SpeedMax and SpeedAvg in Drive (km->mi) (#37 by @alecdoconnor)

## [1.4.9] - 2021-05-03

### Changed
- bump crazy-max/ghaction-docker-meta from v2.3.0 to v2.4.0 (#35 by dependabot)
- updating build workflow by removing if

## [1.4.8] - 2021-04-22

### Changed
- bump github.com/lib/pq from 1.10.0 to 1.10.1 (#33 by dependabot)

## [1.4.7] - 2021-04-13

### Changed
- bump actions/cache from v2.1.4 to v2.1.5 (#32 by dependabot)
- renaming file to lowercase to match other naming

## [1.4.6] - 2021-04-09

### Changed
- bump github.com/gin-gonic/gin from 1.6.3 to 1.7.1 (#31 by dependabot)
- bump crazy-max/ghaction-docker-meta from v2.2.1 to v2.3.0 (#30 by dependabot)

## [1.4.5] - 2021-04-06

### Changed
- bump crazy-max/ghaction-docker-meta from v2.2.0 to v2.2.1 (#29 by dependabot)

## [1.4.4] - 2021-04-05

### Changed
- bump crazy-max/ghaction-docker-meta from v2.1.1 to v2.2.0 (#28 by dependabot)
- bump golang from 1.16.2 to 1.16.3 (#27 by dependabot)

## [1.4.3] - 2021-04-01

### Changed
- bump github.com/eclipse/paho.mqtt.golang from 1.3.2 to 1.3.3 (#26 by dependabot)

### Removed
- github.com/go-sql-driver/mysql removed, since NullTime isn't used/supported anymore

## [1.4.2] - 2021-03-31

### Changed
- bump crazy-max/ghaction-docker-meta from v2.1.0 to v2.1.1 (#24 by dependabot)

## [1.4.1] - 2021-03-30

### Changed
- bump crazy-max/ghaction-docker-meta from v1 to v2.1.0 (#23 by dependabot)

## [1.4.0] - 2021-03-25

### Added
- added feature commands to proxy POST commands to Tesla owner API (#22)
- support for authentication on command endpoints

## [1.3.1] - 2021-03-23

### Fixed
- fixing sql error when BatteryHeaterNoPower is null (#19 by @LelandSindt)

## [1.3.0] - 2021-03-17

### Added
- adding mqtt sleep time before doing disconnect (#17)

## [1.2.3] - 2021-03-16

### Added
- adding probot for stale and no-response

### Changed
- bump golang from 1.16.0 to 1.16.2 (#12 by dependabot)
- bump go mod version from 1.15 to 1.16

## [1.2.2] - 2021-03-11

### Changed
- bump github.com/lib/pq from 1.9.0 to 1.10.0 (#7 by dependabot)
- adjustment in logging

## [1.2.1] - 2021-03-02

### Fixed
- fixing endpoint redirect to /api/v1 destinations
- resolving path issue with traefik

## [1.2.0] - 2021-03-02

### Added
- adding version into URL for better versioning of api

### Changed
- previous endpoints (without versioning) return 301 towards new uri
- renaming of all go files to see version number

## [1.1.1] - 2021-02-18

### Added
- setting mqtt cleansession flag for unsubscribe on disconnect

### Changed
- changing to one multi-subscribe instead of 46 separate subscribes on mqtt

## [1.1.0] - 2021-02-18

### Added
- adding codeql-analysis workflow
- adding dependabot for gomod and docker
- using go mod now

### Changed
- calling on functions without params and using gin.Context in functions instead
- logging for better readability (some rows based on DEBUG_MODE)
- merged TeslaMateAPICars and TeslaMateAPICarsID into one file
- updating Dockerfile a little
- renaming of functions
- bump golang from 1.15.8 to 1.16.0 (#2 by dependabot)
- some code cleanup

### Fixed
- sql query issue with TeslaMateAPICars

## [1.0.2] - 2021-02-15

### Fixed
- sql query error

## [1.0.1] - 2021-02-15

### Added
- / endpoint saying API is running
- DEBUG_MODE variable (printing out debug of TeslaMateApi if set to true)

### Changed
- specifying port 8080 in Run()
- updated Traefik example in README
- code cleanup

### Fixed
- added missing tzdata package in Dockerfile

## [1.0.0] - 2021-02-15

[Unreleased]: https://github.com/tobiasehlert/teslamateapi/compare/v1.10.0...HEAD
[1.10.0]: https://github.com/tobiasehlert/teslamateapi/compare/v1.9.0...v1.10.0
[1.9.0]: https://github.com/tobiasehlert/teslamateapi/compare/v1.8.0...v1.9.0
[1.8.0]: https://github.com/tobiasehlert/teslamateapi/compare/v1.7.1...v1.8.0
[1.7.1]: https://github.com/tobiasehlert/teslamateapi/compare/v1.7.0...v1.7.1
[1.7.0]: https://github.com/tobiasehlert/teslamateapi/compare/v1.6.2...v1.7.0
[1.6.2]: https://github.com/tobiasehlert/teslamateapi/compare/v1.6.1...v1.6.2
[1.6.1]: https://github.com/tobiasehlert/teslamateapi/compare/v1.6.0...v1.6.1
[1.6.0]: https://github.com/tobiasehlert/teslamateapi/compare/v1.5.0...v1.6.0
[1.5.0]: https://github.com/tobiasehlert/teslamateapi/compare/v1.4.9...v1.5.0
[1.4.9]: https://github.com/tobiasehlert/teslamateapi/compare/v1.4.8...v1.4.9
[1.4.8]: https://github.com/tobiasehlert/teslamateapi/compare/v1.4.7...v1.4.8
[1.4.7]: https://github.com/tobiasehlert/teslamateapi/compare/v1.4.6...v1.4.7
[1.4.6]: https://github.com/tobiasehlert/teslamateapi/compare/v1.4.5...v1.4.6
[1.4.5]: https://github.com/tobiasehlert/teslamateapi/compare/v1.4.4...v1.4.5
[1.4.4]: https://github.com/tobiasehlert/teslamateapi/compare/v1.4.3...v1.4.4
[1.4.3]: https://github.com/tobiasehlert/teslamateapi/compare/v1.4.2...v1.4.3
[1.4.2]: https://github.com/tobiasehlert/teslamateapi/compare/v1.4.1...v1.4.2
[1.4.1]: https://github.com/tobiasehlert/teslamateapi/compare/v1.4.0...v1.4.1
[1.4.0]: https://github.com/tobiasehlert/teslamateapi/compare/v1.3.1...v1.4.0
[1.3.1]: https://github.com/tobiasehlert/teslamateapi/compare/v1.3.0...v1.3.1
[1.3.0]: https://github.com/tobiasehlert/teslamateapi/compare/v1.2.3...v1.3.0
[1.2.3]: https://github.com/tobiasehlert/teslamateapi/compare/v1.2.2...v1.2.3
[1.2.2]: https://github.com/tobiasehlert/teslamateapi/compare/v1.2.1...v1.2.2
[1.2.1]: https://github.com/tobiasehlert/teslamateapi/compare/v1.2.0...v1.2.1
[1.2.0]: https://github.com/tobiasehlert/teslamateapi/compare/v1.1.1...v1.2.0
[1.1.1]: https://github.com/tobiasehlert/teslamateapi/compare/v1.1.0...v1.1.1
[1.1.0]: https://github.com/tobiasehlert/teslamateapi/compare/v1.0.2...v1.1.0
[1.0.2]: https://github.com/tobiasehlert/teslamateapi/compare/v1.0.1...v1.0.2
[1.0.1]: https://github.com/tobiasehlert/teslamateapi/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/tobiasehlert/teslamateapi/releases/tag/v1.0.0
