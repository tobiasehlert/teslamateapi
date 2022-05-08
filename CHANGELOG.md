# Changelog

## [1.14.0] - 2022-05-08

### Added
* adding gzip compression ([#143](https://github.com/tobiasehlert/teslamateapi/pull/143) by [tobiasehlert](https://github.com/tobiasehlert))
* adding 404 for not found endpoints ([#144](https://github.com/tobiasehlert/teslamateapi/pull/144) by [tobiasehlert](https://github.com/tobiasehlert))
* disabling proxy feature of gin ([#145](https://github.com/tobiasehlert/teslamateapi/pull/145) by [tobiasehlert](https://github.com/tobiasehlert))
* adding graceful shutdown to gin ([#146](https://github.com/tobiasehlert/teslamateapi/pull/146) by [tobiasehlert](https://github.com/tobiasehlert))

### Changed
* adding two new fields in status endpoint ([#148](https://github.com/tobiasehlert/teslamateapi/pull/148) by [tobiasehlert](https://github.com/tobiasehlert))
* bump golang from 1.17.6 to 1.18.1 ([#150](https://github.com/tobiasehlert/teslamateapi/pull/150), [#157](https://github.com/tobiasehlert/teslamateapi/pull/157), [#161](https://github.com/tobiasehlert/teslamateapi/pull/161), [#165](https://github.com/tobiasehlert/teslamateapi/pull/165) by [dependabot](https://github.com/dependabot))
* bump github.com/lib/pq from 1.10.4 to 1.10.5 ([#164](https://github.com/tobiasehlert/teslamateapi/pull/164) by [dependabot](https://github.com/dependabot))
* bump github/codeql-action from 1 to 2 ([#166](https://github.com/tobiasehlert/teslamateapi/pull/166) by [dependabot](https://github.com/dependabot))
* bump various workflow versions ([#147](https://github.com/tobiasehlert/teslamateapi/pull/147), [#152](https://github.com/tobiasehlert/teslamateapi/pull/152), [#154](https://github.com/tobiasehlert/teslamateapi/pull/154), [#155](https://github.com/tobiasehlert/teslamateapi/pull/155), [#159](https://github.com/tobiasehlert/teslamateapi/pull/159), [#162](https://github.com/tobiasehlert/teslamateapi/pull/162), [#163](https://github.com/tobiasehlert/teslamateapi/pull/163), [#167](https://github.com/tobiasehlert/teslamateapi/pull/167), [#168](https://github.com/tobiasehlert/teslamateapi/pull/168), [#170](https://github.com/tobiasehlert/teslamateapi/pull/170), [#171](https://github.com/tobiasehlert/teslamateapi/pull/171), [#172](https://github.com/tobiasehlert/teslamateapi/pull/172), [#173](https://github.com/tobiasehlert/teslamateapi/pull/173), [#174](https://github.com/tobiasehlert/teslamateapi/pull/174) by [dependabot](https://github.com/dependabot))

### Fixed
* updating getEnv function log ([#156](https://github.com/tobiasehlert/teslamateapi/pull/156) by [tobiasehlert](https://github.com/tobiasehlert))

## [1.13.3] - 2022-01-27

### Fixed
* fix append of commands to allowList ([#142](https://github.com/tobiasehlert/teslamateapi/pull/142) by [tobiasehlert](https://github.com/tobiasehlert))

## [1.13.2] - 2022-01-21

### Changed
* bump docker/build-push-action from 2.7.0 to 2.8.0 ([#135](https://github.com/tobiasehlert/teslamateapi/pull/135) by [dependabot](https://github.com/dependabot))

### Fixed
* slimming down on the codebase and fixing bug with drivedetails view ([#134](https://github.com/tobiasehlert/teslamateapi/pull/134) by [tobiasehlert](https://github.com/tobiasehlert))

## [1.13.1] - 2022-01-12

### Changed
* bump golang from 1.17.5 to 1.17.6 ([#133](https://github.com/tobiasehlert/teslamateapi/pull/133) by [dependabot](https://github.com/dependabot))
* removing code smell from SonarCloud ([#132](https://github.com/tobiasehlert/teslamateapi/pull/132) by [tobiasehlert](https://github.com/tobiasehlert))

## [1.13.0] - 2022-01-05

### Changed
- simplified response handler for communication ([#110](https://github.com/tobiasehlert/teslamateapi/pull/110) by [alecdoconnor](https://github.com/alecdoconnor), [tobiasehlert](https://github.com/tobiasehlert))

## [1.12.1] - 2021-12-31

### Changed
- using of BasePath function in output and redirect ([#130](https://github.com/tobiasehlert/teslamateapi/pull/130) by [tobiasehlert](https://github.com/tobiasehlert))
- updating build workflow with enhancements ([#131](https://github.com/tobiasehlert/teslamateapi/pull/131) by [tobiasehlert](https://github.com/tobiasehlert))

## [1.12.0] - 2021-12-28

### Fixed
- bug in new installations by changing float64 to NullFloat64 on cars efficiency ([#129](https://github.com/tobiasehlert/teslamateapi/pull/129) by [tobiasehlert](https://github.com/tobiasehlert))

## [1.11.1] - 2021-12-23

### Changed
- bump golang from 1.17.2 to 1.17.5 ([#115](https://github.com/tobiasehlert/teslamateapi/pull/115), [#122](https://github.com/tobiasehlert/teslamateapi/pull/122), [#124](https://github.com/tobiasehlert/teslamateapi/pull/124) by [dependabot](https://github.com/dependabot))
- bump github.com/gin-gonic/gin from 1.7.4 to 1.7.7 ([#119](https://github.com/tobiasehlert/teslamateapi/pull/119), [#120](https://github.com/tobiasehlert/teslamateapi/pull/120) by [dependabot](https://github.com/dependabot))
- bump github.com/lib/pq from 1.10.3 to 1.10.4 ([#116](https://github.com/tobiasehlert/teslamateapi/pull/116) by [dependabot](https://github.com/dependabot))
- bump various workflow versions ([#118](https://github.com/tobiasehlert/teslamateapi/pull/118), [#121](https://github.com/tobiasehlert/teslamateapi/pull/121), [#123](https://github.com/tobiasehlert/teslamateapi/pull/123), [#125](https://github.com/tobiasehlert/teslamateapi/pull/125), [#126](https://github.com/tobiasehlert/teslamateapi/pull/126) by [dependabot](https://github.com/dependabot))
- updating build workflow and Dockerfile (by [tobiasehlert](https://github.com/tobiasehlert))

## [1.11.0] - 2021-11-08

### Added
* add support for new endpoints with 4.2.2 ([#113](https://github.com/tobiasehlert/teslamateapi/pull/113) by @michaeldyrynda, [tobiasehlert](https://github.com/tobiasehlert))

### Changed
* bump docker/metadata-action from 3.5.0 to 3.6.0 ([#111](https://github.com/tobiasehlert/teslamateapi/pull/111) by [dependabot](https://github.com/dependabot))

## [1.10.2] - 2021-10-15

### Changed
- bump golang from 1.17.1 to 1.17.2 ([#108](https://github.com/tobiasehlert/teslamateapi/pull/108) by [dependabot](https://github.com/dependabot))
- updating readme with table for variables

## [1.10.1] - 2021-09-23

### Changed
- bump golang from 1.16.5 to 1.17.1 ([#90](https://github.com/tobiasehlert/teslamateapi/pull/90), [#96](https://github.com/tobiasehlert/teslamateapi/pull/96), [#103](https://github.com/tobiasehlert/teslamateapi/pull/103), [#104](https://github.com/tobiasehlert/teslamateapi/pull/104), [#107](https://github.com/tobiasehlert/teslamateapi/pull/107) by [dependabot](https://github.com/dependabot))
- bump github.com/gin-gonic/gin from 1.7.2 to 1.7.4 ([#94](https://github.com/tobiasehlert/teslamateapi/pull/94), [#100](https://github.com/tobiasehlert/teslamateapi/pull/100) by [dependabot](https://github.com/dependabot))
- bump github.com/lib/pq from 1.10.2 to 1.10.3 ([#105](https://github.com/tobiasehlert/teslamateapi/pull/105) by [dependabot](https://github.com/dependabot))
- bump various workflow versions ([#85](https://github.com/tobiasehlert/teslamateapi/pull/85), [#86](https://github.com/tobiasehlert/teslamateapi/pull/86), [#87](https://github.com/tobiasehlert/teslamateapi/pull/87), [#89](https://github.com/tobiasehlert/teslamateapi/pull/89), [#91](https://github.com/tobiasehlert/teslamateapi/pull/91), [#101](https://github.com/tobiasehlert/teslamateapi/pull/101), [#102](https://github.com/tobiasehlert/teslamateapi/pull/102), [#106](https://github.com/tobiasehlert/teslamateapi/pull/106) by [dependabot](https://github.com/dependabot))

## [1.10.0] - 2021-06-30

### Added
- adding power mqtt value to status endpoint ([#74](https://github.com/tobiasehlert/teslamateapi/pull/74))

### Changed
- improved log message of convert-functions ([#75](https://github.com/tobiasehlert/teslamateapi/pull/75) by [alecdoconnor](https://github.com/alecdoconnor))
- bump various workflow versions ([#78](https://github.com/tobiasehlert/teslamateapi/pull/78), [#82](https://github.com/tobiasehlert/teslamateapi/pull/82), [#83](https://github.com/tobiasehlert/teslamateapi/pull/83) by [dependabot](https://github.com/dependabot))
- minor code adjustments based on go-staticcheck

## [1.9.0] - 2021-06-14

### Added
- option to disable auth token for commands ([#71](https://github.com/tobiasehlert/teslamateapi/pull/71) by [LelandSindt](https://github.com/LelandSindt))

### Changed
- bump golang from 1.16.4 to 1.16.5 ([#72](https://github.com/tobiasehlert/teslamateapi/pull/72) by [dependabot](https://github.com/dependabot))
- bump github.com/eclipse/paho.mqtt.golang from 1.3.4 to 1.3.5 ([#73](https://github.com/tobiasehlert/teslamateapi/pull/73) by [dependabot](https://github.com/dependabot))

## [1.8.0] - 2021-06-01

### Fixed
- fixing sql error when EndDate is null ([#58](https://github.com/tobiasehlert/teslamateapi/pull/58) and [#69](https://github.com/tobiasehlert/teslamateapi/pull/69) by [alecdoconnor](https://github.com/alecdoconnor))

## [1.7.1] - 2021-06-01

### Changed
- bump github.com/gin-gonic/gin from 1.7.1 to 1.7.2 ([#63](https://github.com/tobiasehlert/teslamateapi/pull/63) by [dependabot](https://github.com/dependabot))
- bump various workflow versions ([#64](https://github.com/tobiasehlert/teslamateapi/pull/64), [#65](https://github.com/tobiasehlert/teslamateapi/pull/65), [#66](https://github.com/tobiasehlert/teslamateapi/pull/66), [#67](https://github.com/tobiasehlert/teslamateapi/pull/67) by [dependabot](https://github.com/dependabot))

## [1.7.0] - 2021-05-20

### Added
- feature to resume/suspend logging of TeslaMate through TeslaMateApi ([#34](https://github.com/tobiasehlert/teslamateapi/pull/34), [#45](https://github.com/tobiasehlert/teslamateapi/pull/45) and [#48](https://github.com/tobiasehlert/teslamateapi/pull/48) by [LelandSindt](https://github.com/LelandSindt))

### Changed
- minor code adjustments based on go-staticcheck

## [1.6.2] - 2021-05-19

### Changed
- bump github.com/lib/pq from 1.10.1 to 1.10.2 ([#59](https://github.com/tobiasehlert/teslamateapi/pull/59) by [dependabot](https://github.com/dependabot))
- bump various workflow versions ([#52](https://github.com/tobiasehlert/teslamateapi/pull/52), [#53](https://github.com/tobiasehlert/teslamateapi/pull/53), [#54](https://github.com/tobiasehlert/teslamateapi/pull/54), [#55](https://github.com/tobiasehlert/teslamateapi/pull/55), [#56](https://github.com/tobiasehlert/teslamateapi/pull/56) and [#57](https://github.com/tobiasehlert/teslamateapi/pull/57) by [dependabot](https://github.com/dependabot))

## [1.6.1] - 2021-05-10

### Changed
- bump golang from 1.16.3 to 1.16.4 ([#49](https://github.com/tobiasehlert/teslamateapi/pull/49) by [dependabot](https://github.com/dependabot))
- bump github.com/eclipse/paho.mqtt.golang from 1.3.3 to 1.3.4 ([#46](https://github.com/tobiasehlert/teslamateapi/pull/46) by [dependabot](https://github.com/dependabot))

## [1.6.0] - 2021-05-03

### Added
- doing persistant mqtt connection for status collection ([#16](https://github.com/tobiasehlert/teslamateapi/pull/16) by [LelandSindt](https://github.com/LelandSindt) and [#21](https://github.com/tobiasehlert/teslamateapi/pull/21) by [MattBrittan](https://github.com/MattBrittan))
- adding randomized string to mqtt client ([#15](https://github.com/tobiasehlert/teslamateapi/pull/15) by [LelandSindt](https://github.com/LelandSindt))

### Changed
- fixing sql error when FastChargerBrand is null ([#39](https://github.com/tobiasehlert/teslamateapi/pull/39) by [alecdoconnor](https://github.com/alecdoconnor))
- updating workflow for stale issues/PRs

### Removed
- removing MQTT_SLEEPTIME option for mqtt connection
- removing workflow for no response

## [1.5.0] - 2021-05-03

### Changed
- missing convertion of SpeedMax and SpeedAvg in Drive (km->mi) ([#37](https://github.com/tobiasehlert/teslamateapi/pull/37) by [alecdoconnor](https://github.com/alecdoconnor))

## [1.4.9] - 2021-05-03

### Changed
- bump crazy-max/ghaction-docker-meta from v2.3.0 to v2.4.0 ([#35](https://github.com/tobiasehlert/teslamateapi/pull/35) by [dependabot](https://github.com/dependabot))
- updating build workflow by removing if

## [1.4.8] - 2021-04-22

### Changed
- bump github.com/lib/pq from 1.10.0 to 1.10.1 ([#33](https://github.com/tobiasehlert/teslamateapi/pull/33) by [dependabot](https://github.com/dependabot))

## [1.4.7] - 2021-04-13

### Changed
- bump actions/cache from v2.1.4 to v2.1.5 ([#32](https://github.com/tobiasehlert/teslamateapi/pull/32) by [dependabot](https://github.com/dependabot))
- renaming file to lowercase to match other naming

## [1.4.6] - 2021-04-09

### Changed
- bump github.com/gin-gonic/gin from 1.6.3 to 1.7.1 ([#31](https://github.com/tobiasehlert/teslamateapi/pull/31) by [dependabot](https://github.com/dependabot))
- bump crazy-max/ghaction-docker-meta from v2.2.1 to v2.3.0 ([#30](https://github.com/tobiasehlert/teslamateapi/pull/30) by [dependabot](https://github.com/dependabot))

## [1.4.5] - 2021-04-06

### Changed
- bump crazy-max/ghaction-docker-meta from v2.2.0 to v2.2.1 ([#29](https://github.com/tobiasehlert/teslamateapi/pull/29) by [dependabot](https://github.com/dependabot))

## [1.4.4] - 2021-04-05

### Changed
- bump crazy-max/ghaction-docker-meta from v2.1.1 to v2.2.0 ([#28](https://github.com/tobiasehlert/teslamateapi/pull/28) by [dependabot](https://github.com/dependabot))
- bump golang from 1.16.2 to 1.16.3 ([#27](https://github.com/tobiasehlert/teslamateapi/pull/27) by [dependabot](https://github.com/dependabot))

## [1.4.3] - 2021-04-01

### Changed
- bump github.com/eclipse/paho.mqtt.golang from 1.3.2 to 1.3.3 ([#26](https://github.com/tobiasehlert/teslamateapi/pull/26) by [dependabot](https://github.com/dependabot))

### Removed
- github.com/go-sql-driver/mysql removed, since NullTime isn't used/supported anymore

## [1.4.2] - 2021-03-31

### Changed
- bump crazy-max/ghaction-docker-meta from v2.1.0 to v2.1.1 ([#24](https://github.com/tobiasehlert/teslamateapi/pull/24) by [dependabot](https://github.com/dependabot))

## [1.4.1] - 2021-03-30

### Changed
- bump crazy-max/ghaction-docker-meta from v1 to v2.1.0 ([#23](https://github.com/tobiasehlert/teslamateapi/pull/23) by [dependabot](https://github.com/dependabot))

## [1.4.0] - 2021-03-25

### Added
- added feature commands to proxy POST commands to Tesla owner API ([#22](https://github.com/tobiasehlert/teslamateapi/pull/22))
- support for authentication on command endpoints

## [1.3.1] - 2021-03-23

### Fixed
- fixing sql error when BatteryHeaterNoPower is null ([#19](https://github.com/tobiasehlert/teslamateapi/pull/19) by [LelandSindt](https://github.com/LelandSindt))

## [1.3.0] - 2021-03-17

### Added
- adding mqtt sleep time before doing disconnect ([#17](https://github.com/tobiasehlert/teslamateapi/pull/17))

## [1.2.3] - 2021-03-16

### Added
- adding probot for stale and no-response

### Changed
- bump golang from 1.16.0 to 1.16.2 ([#12](https://github.com/tobiasehlert/teslamateapi/pull/12) by [dependabot](https://github.com/dependabot))
- bump go mod version from 1.15 to 1.16

## [1.2.2] - 2021-03-11

### Changed
- bump github.com/lib/pq from 1.9.0 to 1.10.0 ([#7](https://github.com/tobiasehlert/teslamateapi/pull/7) by [dependabot](https://github.com/dependabot))
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
- bump golang from 1.15.8 to 1.16.0 ([#2](https://github.com/tobiasehlert/teslamateapi/pull/2) by [dependabot](https://github.com/dependabot))
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

Initial commit

[1.14.0]: https://github.com/tobiasehlert/teslamateapi/compare/v1.13.3...v1.14.0
[1.13.3]: https://github.com/tobiasehlert/teslamateapi/compare/v1.13.2...v1.13.3
[1.13.2]: https://github.com/tobiasehlert/teslamateapi/compare/v1.13.1...v1.13.2
[1.13.1]: https://github.com/tobiasehlert/teslamateapi/compare/v1.13.0...v1.13.1
[1.13.0]: https://github.com/tobiasehlert/teslamateapi/compare/v1.12.1...v1.13.0
[1.12.1]: https://github.com/tobiasehlert/teslamateapi/compare/v1.12.0...v1.12.1
[1.12.0]: https://github.com/tobiasehlert/teslamateapi/compare/v1.11.1...v1.12.0
[1.11.1]: https://github.com/tobiasehlert/teslamateapi/compare/v1.11.0...v1.11.1
[1.11.0]: https://github.com/tobiasehlert/teslamateapi/compare/v1.10.2...v1.11.0
[1.10.2]: https://github.com/tobiasehlert/teslamateapi/compare/v1.10.1...v1.10.2
[1.10.1]: https://github.com/tobiasehlert/teslamateapi/compare/v1.10.0...v1.10.1
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
[1.0.0]: https://github.com/tobiasehlert/teslamateapi/compare/31dcb4b...v1.0.0
