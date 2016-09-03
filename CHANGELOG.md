# Change Log
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## [0.4.0] - 2016-09-03

### Added
- debug flag `-d`
- github wildcard support
- a short summary at the bottom of the message
- pagination for github repos

### Changed
- log output is formatted differently
- asynchronous fetching of pull requests
- sort gitlab pull requests on updated date

## [0.3.0] - 2016-09-02
### Added
- JSON file for configuration
- User Whitelist

## [0.2.0] - 2016-08-31
### Added
- Find pull requests from Gitlab
- Skip pull requests that contains `[WIP]` or `WIP` in the title
- Show "time since updated" format

### Changed
- Errors will be printed to the STDOUT instead of full panic

## [0.1.0] - 2016-08-28

Initial release

### Added
- Find pull request from Github
- Send pull requests to a slack channel
