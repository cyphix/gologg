# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.0] - 2026-06-24

### Changed
- Module path moved from `gitea.cyphix.dev/kade/go.logg` to `github.com/cyphix/gologg`.

## [0.1.0] - 2026-02-14

### Added
- Initial release of `logg`.
- Thread-safe `zerolog` wrapper with global state management.
- `proxyWriter` for dynamic log redirection (supports loggers created before `Init`).
- `Init`, `InitConsole`, `InitWithWriter`, and `InitWithWriters` for flexible setup.
- `SetSilent` for global log suppression.
- `Ctx` helper for standardized structured logging with package and component context.
- Custom key mapping support via `SetKeys`.
