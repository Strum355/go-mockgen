# Changelog

## [Unreleased]

Nothing yet!

## [v1.3.6] - 2022-10-26

- Config file loading now disambiguates interfaces with the same names in multiple search packages. [#36](https://github.com/derision-test/go-mockgen/pull/36)

## [v1.3.5] - 2022-10-26

- No updates.

## [v1.3.4] - 2022-08-11

- Added `NeedDeps` mode to package loader (prevents fatal log from within `golang.org/x/tools/go/packages` in some circumstances). [3ae60a2](https://github.com/derision-test/go-mockgen/commit/3ae60a20c75f7eb1ae85fc6af66f237f5ee1a04d)

## [v1.3.3] - 2022-06-09

- Added support for `include-config-paths` in config file. [#35](https://github.com/derision-test/go-mockgen/pull/35)

## [v1.3.2] - 2022-06-09

### Added

- Added support for `sources` in config file. [#33](https://github.com/derision-test/go-mockgen/pull/33)

### Fixed

- Fixed broken `import-path` flag. [#34](https://github.com/derision-test/go-mockgen/pull/34)

## [v1.3.1] - 2022-06-06

- Added `--file-prefix` flag. [#32](https://github.com/derision-test/go-mockgen/pull/32)

## [v1.3.0] - 2022-06-06

### Added

- Added support for configuration files. [#31](https://github.com/derision-test/go-mockgen/pull/31)
- Added `--constructor-prefix` flag. [#28](https://github.com/derision-test/go-mockgen/pull/28)

## [v1.2.0] - 2022-03-28

### Changed

- Fixed generation of code with inline interface definitions. [#23](https://github.com/derision-test/go-mockgen/pull/23)
- Added basic support for generic interfaces - now requires Go 1.18 or above. [#20](https://github.com/derision-test/go-mockgen/pull/20)

## [v1.1.5] - 2022-04-08

### Changed

- Updated x/tools for Go 1.18 support. [#22](https://github.com/derision-test/go-mockgen/pull/22)

## [v1.1.4] - 2022-02-01

### Changed

- Fixed generation for nested package on Windows. [#19](https://github.com/derision-test/go-mockgen/pull/19)
- Fixed support for array types in method signatures. [#21](https://github.com/derision-test/go-mockgen/pull/21)

## [v1.1.3] - 2022-02-21

### Added

- Added `--exclude`/`-e` flag to support exclusion of target interfaces. [#13](https://github.com/derision-test/go-mockgen/pull/13)
- Added `--for-test` flag. [#14](https://github.com/derision-test/go-mockgen/pull/14)
- Added `NewStrictMockX` constructor. [#16](https://github.com/derision-test/go-mockgen/pull/16)

## [v1.1.2] - 2021-06-14

No significant changes (only corrected version output).

## [v1.1.1] - 2021-06-14

### Added

- Added `--goimports` flag. [0f4ed82](https://github.com/derision-test/go-mockgen/commit/0f4ed82247eff5446b885c3ea48f48b870a9ee4a)

## [v1.0.0] - 2021-06-14

### Added

- Added support for testify assertions. [#3](https://github.com/derision-test/go-mockgen/pull/3), [#8](https://github.com/derision-test/go-mockgen/pull/8)

### Changed

- Migrated from [efritz/go-mockgen](https://github.com/efritz/go-mockgen). [#1](https://github.com/derision-test/go-mockgen/pull/1)
- We now run `goimports` over rendered files. [096f848](https://github.com/derision-test/go-mockgen/commit/096f848333579e185c8018ff2d17688e4b5f6f27)
- Fixed output paths when directories are generated. [#10](https://github.com/derision-test/go-mockgen/pull/10)

[Unreleased]: https://github.com/derision-test/go-mockgen/compare/v1.3.6...HEAD
[v1.0.0]: https://github.com/derision-test/go-mockgen/releases/tag/v1.0.0
[v1.1.1]: https://github.com/derision-test/go-mockgen/compare/v1.0.0...v1.1.1
[v1.1.2]: https://github.com/derision-test/go-mockgen/compare/v1.1.1...v1.1.2
[v1.1.3]: https://github.com/derision-test/go-mockgen/compare/v1.1.2...v1.1.3
[v1.1.4]: https://github.com/derision-test/go-mockgen/compare/v1.1.3...v1.1.4
[v1.2.0]: https://github.com/derision-test/go-mockgen/compare/v1.1.4...v1.2.0
[v1.3.0]: https://github.com/derision-test/go-mockgen/compare/v1.2.0...v1.3.0
[v1.3.1]: https://github.com/derision-test/go-mockgen/compare/v1.3.0...v1.3.1
[v1.3.2]: https://github.com/derision-test/go-mockgen/compare/v1.3.1...v1.3.2
[v1.3.3]: https://github.com/derision-test/go-mockgen/compare/v1.3.2...v1.3.3
[v1.3.4]: https://github.com/derision-test/go-mockgen/compare/v1.3.3...v1.3.4
[v1.3.5]: https://github.com/derision-test/go-mockgen/compare/v1.3.4...v1.3.5
[v1.3.6]: https://github.com/derision-test/go-mockgen/compare/v1.3.5...v1.3.6
