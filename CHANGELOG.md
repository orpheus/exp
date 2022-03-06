# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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

### Changed
### Fixed
### Removed

[unreleased]: https://github.com/olivierlacan/keep-a-changelog/compare/v1.1.0...HEAD
[1.1.0]: https://github.com/olivierlacan/keep-a-changelog/compare/v1.0.0...v1.1.0
