# Change Log
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## [Unreleased]

### Changed

 - Replace glide with go mod
 - Filter out Github [Draft Pull Requests](https://github.blog/2019-02-14-introducing-draft-pull-requests/)

## [0.8.0] - 2018-09-18

- Use gitlab API v4 instead of v3

## [0.7.0] - 2017-12-07

- Configuration format for whitelisting users has been changed to make it easier to
  add further global filters
- enable and disable filters via the config
- output the number of filtered pull requests
- deduplicate github and gitlab repositories

## [0.6.0] - 2017-10-13

### Changed

- Configuration to get private and public repositories for Organisations
- Configuration to get users public repositories

### Removed

- Configuration for getting repositories with wildcard syntax

## [0.5.1] - 2017-10-11

### Added
- use of glide as a dependency manager to lock down the gitlab client to v3 of the gitlab API

### Changed
- escape &, > and < characters for slack
- better Work-In-Progress (WIP) detection from the pull request titles

## [0.5.0] - 2017-05-09

### Changed
- Exclude Github PR that are currently under review

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
