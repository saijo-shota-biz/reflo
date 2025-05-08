# Reflo Design Document

---

## 1. プロジェクト概要

**reflo**は、日々の作業を振り返り、行動を記録し、自分の行動を資産化するための軽量なCLIツールである。

---

## 2. システム構成

```plaintext
reflo/
├── cmd/reflo/          # CLI エントリポイント
│   └── main.go         # コマンド分岐と対話フロー
├── internal/
│   ├── timer/          # タイマー機能（ポモドーロ相当）
│   │   └── timer.go
│   └── logger/         # セッションログの永続化
│       └── json_logger.go
├── CHANGELOG.md
└── README.md
```

---

## 3. 主要コンポーネント

### 3.1 `cmd/reflo/main.go`

* CLI コマンド取り扱い (`start`, `end‑day`, `help`)
* ユーザー入力の読み取り（`bufio.Scanner`）
* `signal.NotifyContext` で Ctrl‑C を graceful exit

### 3.2 `internal/timer`

* `New(duration)` でワンショットタイマー生成
* `Wait(ctx)` で満了またはキャンセル待機
* 25min フォーカス → 5min 休憩をデフォルト値で提供

### 3.3 `internal/logger`

* `Write(session)` : JSON へ追記保存
   `~/.reflo/YYYY-MM-DD.json`（ファイル名は **ローカル日付**、エントリは **ISO‑8601 / UTC**）
* `ReadDay()`: 当日ファイルを読み出し、`[]Session` へアンマーシャル

### 3.4 データ構造

```go
// internal/logger/logger.go

// Session represents a single focus cycle and reflection.
type Session struct {
    StartTime time.Time // ISO‑8601 / UTC
    EndTime   time.Time // ISO‑8601 / UTC
    Goal      string    // 作業宣言
    Retro     string    // 振り返りコメント
}
```


---

## 4. データモデル

### 4.1 行動ログJSON構成

```json
[
  {
    "StartTime": "2025-05-08T09:00:00Z",
    "EndTime":   "2025-05-08T09:25:00Z",
    "Goal":      "API 設計ドキュメント作成",
    "Retro":     "Slack を見過ぎて進捗 80%"
  },
  {
    "StartTime": "2025-05-08T09:30:00Z",
    "EndTime":   "2025-05-08T09:55:00Z",
    "Goal":      "コードレビュー",
    "Retro":     "想定よりコメント少なめ"
  }
]
```

---

## 5. 実装方針

- コマンド分岐はGo標準`flag`パッケージを使用
- ログ操作は実行日に基づく日単分割
- コマンド実装はTDDで進める
- 基本はローカル動作に絞り、クラウド連携やネット通信を使用しない

---

## 6. 拡張性

- TUI化 (Bubble Tea)
- AI APIとの連携による振り返り分析
- マルチデバイス対応 (Windows/Linux)

---

# End of Document

