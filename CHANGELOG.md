# Changelog

All notable changes to **reflo** are documented in this file.

This project follows the guidelines of [Keep a Changelog](https://keepachangelog.com)
and adheres to [Semantic Versioning](https://semver.org).

---

## [Unreleased]

### Added
* _Add new features here._

### Changed
* _Add changes that modify existing behaviour here._

### Fixed
* _Add bug fixes here._

### Removed
* _Add features that have been removed here._

---

## [1.2.0] – 2025-05-28

### Added
* **Human-readable durations** – session summaries now show `25m`, `1h07m`, etc. instead of raw seconds.

### Changed
* **End-time capture moved** – the session now finishes *after* the retrospective, so the reported duration includes your reflection time.
* **Unified Japanese CLI messages** with clear icons (⏳ / ✅ / ⚠️ …) for better at-a-glance feedback.
* **Refined `Ctrl+C` behaviour** – cancelling during focus skips only that focus; cancelling during break skips only that break.

### Fixed
* Break timers no longer get canceled when a focus timer is skipped with `Ctrl+C`.
* A session is still logged even if you cancel the retrospective prompt.

---

## [1.1.0] – 2025-05-17

### Added
* **Interactive line editing & history** – switched to `github.com/chzyer/readline`, enabling arrow-key navigation, multiline editing, and **Ctrl + D** submission.
* **Contextual help hint** shown above every prompt: `Enter to add newline · Ctrl+D to submit · Ctrl+C to quit`.

### Changed
* Replaced the previous `bufio.Scanner` input loop with **readline**, improving Unicode handling, long-line support, and overall UX.
* Streamlined the session flow: after a break finishes, the next focus cycle starts automatically; the old `Start another session? [y/n]` confirmation has been removed.
* Prompt messages have been consolidated and localised for greater clarity.

### Fixed
* Goals and retrospectives containing spaces or multibyte characters are now logged without truncation.
* Removed stray `\r` characters on Windows that previously caused parsing issues.

### Removed
* Legacy `bufio`-based input code.
* The `y/n` prompt that asked whether to start another session at the end of each cycle.

---

## [1.0.0] – 2025-05-08

### Added
* **Focus session timer** with default 25-minute focus / 5-minute break cycle (`reflo start`).
* **Graceful cancellation**: `Ctrl-C` (SIGINT) cleanly stops timers via the new `Wait(ctx)` implementation.
* **Automatic break timer** followed by an audible bell when focus or break completes.
* **JSON session logging** to `~/.reflo/YYYY-MM-DD.json` (file name uses local date, entries are ISO-8601 timestamps in UTC) with goal & retrospective notes.
