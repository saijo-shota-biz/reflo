# Reflo

> Reflect your actions, flow your life, and log your growth.

---

## 📚 Overview

**reflo** is a lightweight CLI tool designed to help you build focused work rhythms, reflect on your daily activities, and log your actions for future personal growth.

- Manage your focus sessions with Pomodoro-based cycles
- Reflect on your completed tasks naturally
- Save action logs as personal assets for future AI analysis
- Generate daily summaries for easy review and reporting

---

## 🚀 Features

- Minimalist CLI design for instant startup
- Session declaration and reflection input
- 25-minute focus timer with an audible bell on completion
- 5-minute break management
- **Automatic JSON session logging** to `~/.reflo/YYYY-MM-DD.json`
  (filename uses local date; timestamps are ISO-8601 / UTC)
- Daily summary report (`reflo end-day`)
- Graceful interruption with Ctrl-C

---

## 🛠 Installation

### Homebrew (recommended)

```bash
brew tap saijo-shota-biz/homebrew-reflo
brew install reflo
```

_(Manual binary download from GitHub Releases is also available.)_

---

## 📝 Usage

### Start a focus session

```bash
$ reflo start
```
- Declare the task
- Start a 25-minute focus timer

### End of day summary

```bash
$ reflo end-day
```
- Output a summary of the day's actions
- Automatically copy the report to the clipboard

---

## 🔧 Troubleshooting

### macOS: Notifications not appearing

On first use, you need to enable notifications for reflo:

1. Open Script Editor
   ```bash
   open -a "Script Editor"
   ```

2. Paste and run this script (click the ▶️ play button):
   ```applescript
   display notification "Test notification" with title "Reflo"
   ```

3. When prompted, click "Allow" to grant notification permissions

4. Notifications from reflo should now work properly

**Note**: This is a one-time setup required due to macOS security policies for CLI applications.

---

## 🛤 Future Plans

- TUI version using Bubble Tea
- Weekly and monthly report generation
- Auto-logging based on active window tracking
- AI integration for personal action analysis and advice

---

## 📜 License

This project is licensed under the MIT License.
