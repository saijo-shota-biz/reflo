# Reflo Design Document

---

## 1. プロジェクト概要

**reflo**は、日々の作業を振り返り、行動を記録し、自分の行動を資産化するための軽量なCLIツールである。

---

## 2. システム構成

```plaintext
reflo/
├── main.go           # エントリーポイント (CLI管理)
├── pomodoro.go       # ポモドーロタイマー処理
├── logger.go         # ログ保存/読込処理
├── clipboard.go      # クリップボード操作
├── notification.go   # macOS通知処理
├── utils.go          # 澄用関数（ビープ音など）
├── logs/             # ログ格納用ディレクトリ
└── go.mod            # Go module file
```

---

## 3. 主要コンポーネント

### 3.1 `main.go`
- CLIの入力パース
- `start`、`end-day`コマンド分岐処理

### 3.2 `pomodoro.go`
- タイマーセッション管理
    - 作業定義入力
    - タイマーカウント処理
    - 振り返り入力

### 3.3 `logger.go`
- 日単のログ作成
- JSONログの書き込み、読み込み

### 3.4 `clipboard.go`
- 振り返りまとめのクリップボードコピー

### 3.5 `notification.go`
- macOSでポップアップ通知を送信

### 3.6 `utils.go`
- 時間管理
- ビープ音鳴らし
- フォーマット関数

---

## 4. データモデル

### 4.1 行動ログJSON構成

```json
{
  "date": "2024-04-30",
  "sessions": [
    {
      "start_time": "09:00",
      "task": "API設計ドキュメント作成",
      "reflection": "8割経過。Slack見過ぎた"
    },
    {
      "start_time": "10:00",
      "task": "コードレビュー完了",
      "reflection": "問題なし"
    }
  ]
}
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

