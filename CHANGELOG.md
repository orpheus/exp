# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [2022-06-18]
### Added
- create admin role/user on startup
- login by email __or username

## [2022-06-13]
### Changed
- refactored app to use 'Clean Arch' based on Robert Martin's Clean design philosophy
  - figured I'd try it out. so far I love the separation of the domain

## [2022-03-22]
### Added
- db postgres setup (pgx)
- router setup (gin)
- skillConfig controller, repository, crud methods
- skill controller, repository, crud methods
- add txp to skill with increasing level system
- users, login, sign up
- permission, roles, authorized routes
- allow one skill type per user
- date_last_txp_add db column under skill
- restrict adding txp more than time diff of last add date
- custom `migrate.go` script to load sql migrations into db

### Fixed
### Removed

[unreleased]: https://github.com/olivierlacan/keep-a-changelog/compare/v1.1.0...HEAD
[1.1.0]: https://github.com/olivierlacan/keep-a-changelog/compare/v1.0.0...v1.1.0
