# Changelog

## [Unreleased]

### Added

### Changed

### Fixed

## [0.8.0] - 2026-07-25

### Added
- Added `--version`/`-v` flag to display the current release version
- Added `--update`/`-u` flag to update helm to the latest release (via Homebrew or the install script)

### Changed
- Config storage now uses XDG-aware paths on Linux (`$XDG_CONFIG_HOME/helm` or `~/.config/helm`) instead of always `~/.helm`; existing `~/.helm` configs are migrated automatically on first load. macOS/Windows behavior is unchanged.

### Fixed

## [0.7.0] - 2026-07-24

### Added
- Added `--quick` flag to start a single-stage timer without the menu (e.g. `helm --quick 10`)

### Changed

### Fixed

## [0.6.0] - 2025-12-11

### Added
- Added support for sound alerts on timers (Thanks @muyiwaolurin!)

### Changed

### Fixed
- Fixed customize menu to display actual workflow names

## [0.5.0] - 2025-12-03

### Added

### Changed
- Updated color scheme to helm logo brand colors

### Fixed
- Prevent panic when advancing past completed session
- Prevent panic when creating session with empty workflow

## [0.4.0] - 2025-12-02

### Added
- Add auto-transition feature for all workflows

### Changed
- Update styling and colours

### Fixed

## [0.3.2] - 2025-12-01

### Added

### Changed

### Fixed

## [0.3.1] - 2025-12-01

### Changed
- Improved timer package efficiency by reusing timer instances
- Simplified state management logic
- Added comprehensive unit test coverage for timer and session

### Added
- Initial release of Helm TUI Pomodoro timer
- Pomodoro workflow with classic 25/5 minute cycles
- Design Interview workflow for structured practice
- Custom workflow support with persistence
- Large ASCII timer display
- Progress bar visualization
- Terminal title updates with countdown
- Sound notifications on step completion
- Workflow customization via TUI

