# Changelog

All notable changes to **reflo** will be documented in this file.

This project follows the guidelines of [Keep a Changelog](https://keepachangelog.com)
and adheres to [Semantic Versioning](https://semver.org).

---

## \[Unreleased]

### Added

* *Add new features here.*

### Changed

* *Add changes that modify existing behaviour here.*

### Fixed

* *Add bug fixes here.*

### Removed

* *Add features that have been removed here.*

---

## \[1.0.0] - 2025-05-08

### Added

* **Focus session timer** with default 25‑minute focus / 5‑minute break cycle (`reflo start`).
* **Graceful cancellation**: `Ctrl‑C` (SIGINT) cleanly stops timers via the new `Wait(ctx)` implementation.
* **Automatic break timer** followed by an audible bell when focus or break completes.
* **JSON session logging** to `~/.reflo/YYYY-MM-DD.json` (file name uses *local* date, entries are ISO‑8601 timestamps in UTC) with goal & retrospective notes.
