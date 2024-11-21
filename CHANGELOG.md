# Changelog

## [1.19.0] - 2024-11-21

### Added

- ability to use custom API endpoint ([#307](https://github.com/tobiasehlert/teslamateapi/pull/307) by [jlestel](https://github.com/jlestel))

### Changed

- bump golang from 1.23.2 to 1.23.3 ([#311](https://github.com/tobiasehlert/teslamateapi/pull/311) by [dependabot](https://github.com/dependabot))

## [1.18.3] - 2024-10-15

### Changed

- bump alpine from 3.20.2 to 3.20.3 ([#305](https://github.com/tobiasehlert/teslamateapi/pull/305) by [dependabot](https://github.com/dependabot))
- bump golang from 1.22.5 to 1.23.2 ([#302](https://github.com/tobiasehlert/teslamateapi/pull/302), [#306](https://github.com/tobiasehlert/teslamateapi/pull/306), [#308](https://github.com/tobiasehlert/teslamateapi/pull/308) by [dependabot](https://github.com/dependabot))
- decrease Dockerfile layers ([#303](https://github.com/tobiasehlert/teslamateapi/pull/303) by [tobiasehlert](https://github.com/tobiasehlert))

## [1.18.2] - 2024-08-12

### Changed

- bump alpine from 3.20.1 to 3.20.2 ([#297](https://github.com/tobiasehlert/teslamateapi/pull/297) by [dependabot](https://github.com/dependabot))
- bump github.com/eclipse/paho.mqtt.golang from 1.4.3 to 1.5.0 ([#298](https://github.com/tobiasehlert/teslamateapi/pull/298) by [dependabot](https://github.com/dependabot))

### Fixed

- fix converts the current speed if unit is "mi" on status endpoint ([#299](https://github.com/tobiasehlert/teslamateapi/pull/299) by [ckanoab](https://github.com/ckanoab))
- fix status drivingdetails.speed to be int and not float64 ([#301](https://github.com/tobiasehlert/teslamateapi/pull/301) by [tobiasehlert](https://github.com/tobiasehlert))

## [1.18.1] - 2024-07-24

### Fixed

- invalid memory address or nil pointer dereference ([#296](https://github.com/tobiasehlert/teslamateapi/pull/296) by [tobiasehlert](https://github.com/tobiasehlert))

## [1.18.0] - 2024-07-24

### Added

- add support for k8s health endpoints ([#191](https://github.com/tobiasehlert/teslamateapi/pull/191) by [tobiasehlert](https://github.com/tobiasehlert))
- add cosign of images in build workflow ([#280](https://github.com/tobiasehlert/teslamateapi/pull/280) by [tobiasehlert](https://github.com/tobiasehlert))

### Changed

- add and update of MQTT topics ([#289](https://github.com/tobiasehlert/teslamateapi/pull/289) by [tobiasehlert](https://github.com/tobiasehlert))
- add center_display_state mqtt topic ([#295](https://github.com/tobiasehlert/teslamateapi/pull/295) by [tobiasehlert](https://github.com/tobiasehlert))

### Fixed

- resolving fatal when disabling of mqtt ([#288](https://github.com/tobiasehlert/teslamateapi/pull/288) by [tobiasehlert](https://github.com/tobiasehlert))
- update GitHub action workflow and go mod ([#279](https://github.com/tobiasehlert/teslamateapi/pull/279) by [tobiasehlert](https://github.com/tobiasehlert))
- bump alpine from 3.19.1 to 3.20.1 ([#285](https://github.com/tobiasehlert/teslamateapi/pull/285), [#290](https://github.com/tobiasehlert/teslamateapi/pull/290) by [dependabot](https://github.com/dependabot))
- bump docker/build-push-action from 5 to 6 ([#287](https://github.com/tobiasehlert/teslamateapi/pull/287) by [dependabot](https://github.com/dependabot))
- bump github.com/gin-contrib/gzip from 1.0.0 to 1.0.1 ([#282](https://github.com/tobiasehlert/teslamateapi/pull/282) by [dependabot](https://github.com/dependabot))
- bump github.com/gin-gonic/gin from 1.9.1 to 1.10.0 ([#283](https://github.com/tobiasehlert/teslamateapi/pull/283) by [dependabot](https://github.com/dependabot))
- bump golang from 1.22.2 to 1.22.5 ([#284](https://github.com/tobiasehlert/teslamateapi/pull/284), [#286](https://github.com/tobiasehlert/teslamateapi/pull/286), [#293](https://github.com/tobiasehlert/teslamateapi/pull/293) by [dependabot](https://github.com/dependabot))

## [1.17.2] - 2024-03-30

### Changed

- bump github.com/gin-contrib/gzip from 0.0.6 to 1.0.0 ([#276](https://github.com/tobiasehlert/teslamateapi/pull/276) by [dependabot](https://github.com/dependabot))

### Fixed

- fix: issue after location mqtt implementation ([#278](https://github.com/tobiasehlert/teslamateapi/pull/278) by [tobiasehlert](https://github.com/tobiasehlert))

## [1.17.1] - 2024-03-15

### Fixed

- fix permission issues with Dockerfile nonroot implementation ([#274](https://github.com/tobiasehlert/teslamateapi/pull/274) by [tobiasehlert](https://github.com/tobiasehlert))

## [1.17.0] - 2024-03-15

### Added

- add tire pressure warning and active route from mqtt ([#270](https://github.com/tobiasehlert/teslamateapi/pull/270) by [tobiasehlert](https://github.com/tobiasehlert))

### Changed

- alignment of commands to version 4.23.6 ([#265](https://github.com/tobiasehlert/teslamateapi/pull/265) by [tobiasehlert](https://github.com/tobiasehlert))
- bump actions/cache from 3 to 4 ([#261](https://github.com/tobiasehlert/teslamateapi/pull/261) by [dependabot](https://github.com/dependabot))
- bump golang from 1.21.5 to 1.22.1 ([#260](https://github.com/tobiasehlert/teslamateapi/pull/260), [#268](https://github.com/tobiasehlert/teslamateapi/pull/268), [#271](https://github.com/tobiasehlert/teslamateapi/pull/271) by [dependabot](https://github.com/dependabot))
- bump google.golang.org/protobuf from 1.32.0 to 1.33.0 ([#272](https://github.com/tobiasehlert/teslamateapi/pull/272) by [dependabot](https://github.com/dependabot))
- bump peter-evans/dockerhub-description from 3 to 4 ([#262](https://github.com/tobiasehlert/teslamateapi/pull/262) by [dependabot](https://github.com/dependabot))
- update Dockerfile to specific version and use of nonroot user ([#266](https://github.com/tobiasehlert/teslamateapi/pull/266) by [tobiasehlert](https://github.com/tobiasehlert))
- updating go mods and linting markdown files ([#264](https://github.com/tobiasehlert/teslamateapi/pull/264) by [tobiasehlert](https://github.com/tobiasehlert))

### Fixed

- fix Dockerfile alpine container typo ([#267](https://github.com/tobiasehlert/teslamateapi/pull/267) by [tobiasehlert](https://github.com/tobiasehlert))

## [1.16.6] - 2023-12-19

### Changed

- bump golang from 1.21.2 to 1.21.5 ([#251](https://github.com/tobiasehlert/teslamateapi/pull/251), [#253](https://github.com/tobiasehlert/teslamateapi/pull/253), [#257](https://github.com/tobiasehlert/teslamateapi/pull/257) by [dependabot](https://github.com/dependabot))
- bump actions/setup-go from 4 to 5 ([#256](https://github.com/tobiasehlert/teslamateapi/pull/256) by [dependabot](https://github.com/dependabot))
- bump github/codeql-action from 2 to 3 ([#258](https://github.com/tobiasehlert/teslamateapi/pull/258) by [dependabot](https://github.com/dependabot))
- bump golang.org/x/crypto from 0.14.0 to 0.17.0 ([#259](https://github.com/tobiasehlert/teslamateapi/pull/259) by [dependabot](https://github.com/dependabot))

## [1.16.5] - 2023-10-12

### Changed

- bump golang from 1.21.1 to 1.21.2 ([#247](https://github.com/tobiasehlert/teslamateapi/pull/247) by [dependabot](https://github.com/dependabot))
- bump golang.org/x/net from 0.15.0 to 0.17.0 ([#248](https://github.com/tobiasehlert/teslamateapi/pull/248) by [dependabot](https://github.com/dependabot))
- workflow setup go version by go.mod ([#249](https://github.com/tobiasehlert/teslamateapi/pull/249) by [tobiasehlert](https://github.com/tobiasehlert))

## [1.16.4] - 2023-09-23

### Fixed

- fix NullString issue in several endpoints ([#245](https://github.com/tobiasehlert/teslamateapi/pull/245) by [tobiasehlert](https://github.com/tobiasehlert))

## [1.16.3] - 2023-09-22

### Changed

- bump actions/checkout from 3 to 4 ([#233](https://github.com/tobiasehlert/teslamateapi/pull/233) by [dependabot](https://github.com/dependabot))
- bump docker/build-push-action from 4 to 5 ([#239](https://github.com/tobiasehlert/teslamateapi/pull/239) by [dependabot](https://github.com/dependabot))
- bump docker/login-action from 2 to 3 ([#238](https://github.com/tobiasehlert/teslamateapi/pull/238) by [dependabot](https://github.com/dependabot))
- bump docker/metadata-action from 4 to 5 ([#237](https://github.com/tobiasehlert/teslamateapi/pull/237) by [dependabot](https://github.com/dependabot))
- bump docker/setup-buildx-action from 2 to 3 ([#236](https://github.com/tobiasehlert/teslamateapi/pull/236) by [dependabot](https://github.com/dependabot))
- bump docker/setup-qemu-action from 2 to 3 ([#235](https://github.com/tobiasehlert/teslamateapi/pull/235) by [dependabot](https://github.com/dependabot))
- bump golang from 1.20.6 to 1.21.1 ([#231](https://github.com/tobiasehlert/teslamateapi/pull/231), [#234](https://github.com/tobiasehlert/teslamateapi/pull/234) by [dependabot](https://github.com/dependabot))
- cleaning, bumping and enhancing ([#243](https://github.com/tobiasehlert/teslamateapi/pull/243) by [tobiasehlert](https://github.com/tobiasehlert))

### Fixed

- change CarName to be NullString instead of string ([#242](https://github.com/tobiasehlert/teslamateapi/pull/242) by [tobiasehlert](https://github.com/tobiasehlert))

## [1.16.2] - 2023-07-28

### Changed

- bump github.com/thanhpk/randstr from 1.0.5 to 1.0.6 ([#222](https://github.com/tobiasehlert/teslamateapi/pull/222) by [dependabot](https://github.com/dependabot))
- bump github.com/gin-gonic/gin from 1.9.0 to 1.9.1 ([#223](https://github.com/tobiasehlert/teslamateapi/pull/223) by [dependabot](https://github.com/dependabot))
- bump golang from 1.20.4 to 1.20.6 ([#225](https://github.com/tobiasehlert/teslamateapi/pull/225), [#230](https://github.com/tobiasehlert/teslamateapi/pull/230) by [dependabot](https://github.com/dependabot))
- bump github.com/eclipse/paho.mqtt.golang from 1.4.2 to 1.4.3 ([#229](https://github.com/tobiasehlert/teslamateapi/pull/229) by [dependabot](https://github.com/dependabot))

## [1.16.1] - 2023-05-22

### Changed

- bump docker/build-push-action from 3 to 4 ([#208](https://github.com/tobiasehlert/teslamateapi/pull/208) by [dependabot](https://github.com/dependabot))
- bump github.com/eclipse/paho.mqtt.golang from 1.4.1 to 1.4.2 ([#202](https://github.com/tobiasehlert/teslamateapi/pull/202) by [dependabot](https://github.com/dependabot))
- bump github.com/gin-gonic/gin from 1.8.1 to 1.9.0 ([#206](https://github.com/tobiasehlert/teslamateapi/pull/206), [#211](https://github.com/tobiasehlert/teslamateapi/pull/211) by [dependabot](https://github.com/dependabot))
- bump github.com/lib/pq from 1.10.7 to 1.10.9 ([#220](https://github.com/tobiasehlert/teslamateapi/pull/220) by [dependabot](https://github.com/dependabot))
- bump github.com/thanhpk/randstr from 1.0.4 to 1.0.5 ([#215](https://github.com/tobiasehlert/teslamateapi/pull/215) by [dependabot](https://github.com/dependabot))
- bump golang from 1.19.2 to 1.20.4 ([#204](https://github.com/tobiasehlert/teslamateapi/pull/204), [#205](https://github.com/tobiasehlert/teslamateapi/pull/205), [#207](https://github.com/tobiasehlert/teslamateapi/pull/207), [#210](https://github.com/tobiasehlert/teslamateapi/pull/210), [#213](https://github.com/tobiasehlert/teslamateapi/pull/213), [#216](https://github.com/tobiasehlert/teslamateapi/pull/216), [#221](https://github.com/tobiasehlert/teslamateapi/pull/221) by [dependabot](https://github.com/dependabot))

## [1.16.0] - 2022-10-12

### Changed

- updating `go build` step in dockerfile ([#192](https://github.com/tobiasehlert/teslamateapi/pull/192) by [tobiasehlert](https://github.com/tobiasehlert))
- bump golang from 1.18.3 to 1.19.2 ([#193](https://github.com/tobiasehlert/teslamateapi/pull/193), [#194](https://github.com/tobiasehlert/teslamateapi/pull/194), [#196](https://github.com/tobiasehlert/teslamateapi/pull/196), [#200](https://github.com/tobiasehlert/teslamateapi/pull/200) by [dependabot](https://github.com/dependabot))
- bump github.com/lib/pq from 1.10.6 to 1.10.7 ([#195](https://github.com/tobiasehlert/teslamateapi/pull/195) by [dependabot](https://github.com/dependabot))
- fix mqtt reconnection issue #176 ([#199](https://github.com/tobiasehlert/teslamateapi/pull/199) by [virusbrain](https://github.com/virusbrain) and [LelandSindt](https://github.com/LelandSindt))

## [1.15.0] - 2022-07-15

### Information

🔓 **Encryption of API tokens** was added in [1.27.0](https://github.com/teslamate-org/teslamate/releases/tag/v1.27.0) of TeslaMate.
You therefore need to adjust your TeslaMateApi deployment with the new added environment variables `ENCRYPTION_KEY`.
The `ENCRYPTION_KEY` needs to have the same value as the key in the environment variables of your TeslaMate.

### Added

- TeslaMate encryption of API tokens ([#141](https://github.com/tobiasehlert/teslamateapi/pull/141) by [LelandSindt](https://github.com/LelandSindt), [tobiasehlert](https://github.com/tobiasehlert))
- support execute commands for China region cars ([#184](https://github.com/tobiasehlert/teslamateapi/pull/184) by [richard1122](https://github.com/richard1122))
- support for tire pressure metrics from MQTT ([#186](https://github.com/tobiasehlert/teslamateapi/pull/186) by [tobiasehlert](https://github.com/tobiasehlert))
- support for new commands ([#187](https://github.com/tobiasehlert/teslamateapi/pull/187) by [tobiasehlert](https://github.com/tobiasehlert))

### Changed

- removing **v** in container image tag ([#188](https://github.com/tobiasehlert/teslamateapi/pull/188) by [tobiasehlert](https://github.com/tobiasehlert))
- bump golang from 1.18.1 to 1.18.3 ([#177](https://github.com/tobiasehlert/teslamateapi/pull/177), [#181](https://github.com/tobiasehlert/teslamateapi/pull/181) by [dependabot](https://github.com/dependabot))
- bump github.com/eclipse/paho.mqtt.golang from 1.3.5 to 1.4.1 ([#182](https://github.com/tobiasehlert/teslamateapi/pull/182) by [dependabot](https://github.com/dependabot))
- bump github.com/gin-contrib/gzip from 0.0.5 to 0.0.6 ([#189](https://github.com/tobiasehlert/teslamateapi/pull/189) by [dependabot](https://github.com/dependabot))
- bump github.com/gin-gonic/gin from 1.7.7 to 1.8.1 ([#179](https://github.com/tobiasehlert/teslamateapi/pull/179), [#183](https://github.com/tobiasehlert/teslamateapi/pull/183) by [dependabot](https://github.com/dependabot))
- bump github.com/lib/pq from 1.10.5 to 1.10.6 ([#178](https://github.com/tobiasehlert/teslamateapi/pull/178) by [dependabot](https://github.com/dependabot))
- some go mod and workflow build updates ([#180](https://github.com/tobiasehlert/teslamateapi/pull/180) by [dependabot](https://github.com/dependabot))

## [1.14.0] - 2022-05-08

### Added

- adding gzip compression ([#143](https://github.com/tobiasehlert/teslamateapi/pull/143) by [tobiasehlert](https://github.com/tobiasehlert))
- adding 404 for not found endpoints ([#144](https://github.com/tobiasehlert/teslamateapi/pull/144) by [tobiasehlert](https://github.com/tobiasehlert))
- disabling proxy feature of gin ([#145](https://github.com/tobiasehlert/teslamateapi/pull/145) by [tobiasehlert](https://github.com/tobiasehlert))
- adding graceful shutdown to gin ([#146](https://github.com/tobiasehlert/teslamateapi/pull/146) by [tobiasehlert](https://github.com/tobiasehlert))

### Changed

- adding two new fields in status endpoint ([#148](https://github.com/tobiasehlert/teslamateapi/pull/148) by [tobiasehlert](https://github.com/tobiasehlert))
- bump golang from 1.17.6 to 1.18.1 ([#150](https://github.com/tobiasehlert/teslamateapi/pull/150), [#157](https://github.com/tobiasehlert/teslamateapi/pull/157), [#161](https://github.com/tobiasehlert/teslamateapi/pull/161), [#165](https://github.com/tobiasehlert/teslamateapi/pull/165) by [dependabot](https://github.com/dependabot))
- bump github.com/lib/pq from 1.10.4 to 1.10.5 ([#164](https://github.com/tobiasehlert/teslamateapi/pull/164) by [dependabot](https://github.com/dependabot))
- bump github/codeql-action from 1 to 2 ([#166](https://github.com/tobiasehlert/teslamateapi/pull/166) by [dependabot](https://github.com/dependabot))
- bump various workflow versions ([#147](https://github.com/tobiasehlert/teslamateapi/pull/147), [#152](https://github.com/tobiasehlert/teslamateapi/pull/152), [#154](https://github.com/tobiasehlert/teslamateapi/pull/154), [#155](https://github.com/tobiasehlert/teslamateapi/pull/155), [#159](https://github.com/tobiasehlert/teslamateapi/pull/159), [#162](https://github.com/tobiasehlert/teslamateapi/pull/162), [#163](https://github.com/tobiasehlert/teslamateapi/pull/163), [#167](https://github.com/tobiasehlert/teslamateapi/pull/167), [#168](https://github.com/tobiasehlert/teslamateapi/pull/168), [#170](https://github.com/tobiasehlert/teslamateapi/pull/170), [#171](https://github.com/tobiasehlert/teslamateapi/pull/171), [#172](https://github.com/tobiasehlert/teslamateapi/pull/172), [#173](https://github.com/tobiasehlert/teslamateapi/pull/173), [#174](https://github.com/tobiasehlert/teslamateapi/pull/174) by [dependabot](https://github.com/dependabot))

### Fixed

- updating getEnv function log ([#156](https://github.com/tobiasehlert/teslamateapi/pull/156) by [tobiasehlert](https://github.com/tobiasehlert))

## [1.13.3] - 2022-01-27

### Fixed

- fix append of commands to allowList ([#142](https://github.com/tobiasehlert/teslamateapi/pull/142) by [tobiasehlert](https://github.com/tobiasehlert))

## [1.13.2] - 2022-01-21

### Changed

- bump docker/build-push-action from 2.7.0 to 2.8.0 ([#135](https://github.com/tobiasehlert/teslamateapi/pull/135) by [dependabot](https://github.com/dependabot))

### Fixed

- slimming down on the codebase and fixing bug with drivedetails view ([#134](https://github.com/tobiasehlert/teslamateapi/pull/134) by [tobiasehlert](https://github.com/tobiasehlert))

## [1.13.1] - 2022-01-12

### Changed

- bump golang from 1.17.5 to 1.17.6 ([#133](https://github.com/tobiasehlert/teslamateapi/pull/133) by [dependabot](https://github.com/dependabot))
- removing code smell from SonarCloud ([#132](https://github.com/tobiasehlert/teslamateapi/pull/132) by [tobiasehlert](https://github.com/tobiasehlert))

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

- add support for new endpoints with 4.2.2 ([#113](https://github.com/tobiasehlert/teslamateapi/pull/113) by @michaeldyrynda, [tobiasehlert](https://github.com/tobiasehlert))

### Changed

- bump docker/metadata-action from 3.5.0 to 3.6.0 ([#111](https://github.com/tobiasehlert/teslamateapi/pull/111) by [dependabot](https://github.com/dependabot))

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

[1.19.0]: https://github.com/tobiasehlert/teslamateapi/compare/v1.18.3...v1.19.0
[1.18.3]: https://github.com/tobiasehlert/teslamateapi/compare/v1.18.2...v1.18.3
[1.18.2]: https://github.com/tobiasehlert/teslamateapi/compare/v1.18.1...v1.18.2
[1.18.1]: https://github.com/tobiasehlert/teslamateapi/compare/v1.18.0...v1.18.1
[1.18.0]: https://github.com/tobiasehlert/teslamateapi/compare/v1.17.2...v1.18.0
[1.17.2]: https://github.com/tobiasehlert/teslamateapi/compare/v1.17.1...v1.17.2
[1.17.1]: https://github.com/tobiasehlert/teslamateapi/compare/v1.17.0...v1.17.1
[1.17.0]: https://github.com/tobiasehlert/teslamateapi/compare/v1.16.6...v1.17.0
[1.16.6]: https://github.com/tobiasehlert/teslamateapi/compare/v1.16.5...v1.16.6
[1.16.5]: https://github.com/tobiasehlert/teslamateapi/compare/v1.16.4...v1.16.5
[1.16.4]: https://github.com/tobiasehlert/teslamateapi/compare/v1.16.3...v1.16.4
[1.16.3]: https://github.com/tobiasehlert/teslamateapi/compare/v1.16.2...v1.16.3
[1.16.2]: https://github.com/tobiasehlert/teslamateapi/compare/v1.16.1...v1.16.2
[1.16.1]: https://github.com/tobiasehlert/teslamateapi/compare/v1.16.0...v1.16.1
[1.16.0]: https://github.com/tobiasehlert/teslamateapi/compare/v1.15.0...v1.16.0
[1.15.0]: https://github.com/tobiasehlert/teslamateapi/compare/v1.14.0...v1.15.0
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
