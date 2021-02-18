# Changelog

## [Unreleased]

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
- bumping dependabot versions
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

[Unreleased]: https://github.com/tobiasehlert/teslamateapi/compare/v1.1.0...HEAD
[1.1.0]: https://github.com/tobiasehlert/teslamateapi/compare/v1.0.2...v1.1.0
[1.0.2]: https://github.com/tobiasehlert/teslamateapi/compare/v1.0.1...v1.0.2
[1.0.1]: https://github.com/tobiasehlert/teslamateapi/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/tobiasehlert/teslamateapi/releases/tag/v1.0.0
