<div align="center">

<img src="../.github/assets/banner.svg" alt="Agent-First Feature Flags" width="100%"/>

<br/>

[![Status](https://img.shields.io/badge/🚀%20ステータス-v1.0.0%20リリース可能-10b981?style=flat-square&labelColor=1c2128)](#)
[![ライセンス](https://img.shields.io/badge/ライセンス-MIT-79c0ff?style=flat-square&labelColor=1c2128)](../LICENSE)
[![Go](https://img.shields.io/badge/Go-1.27+-00ADD8?style=flat-square&logo=go&labelColor=1c2128)](#)

<br/>

[English](../README.md) | [Русский](README.ru.md) | [Español](README.es.md) | [中文](README.zh.md) | [日本語](README.ja.md) | [Français](README.fr.md) | [Deutsch](README.de.md)

</div>

---

# MDFLAG

**AIエージェントのための機能フラグ（Feature Flags）。SaaS不要。APIキー不要。Gitリポジトリ内の通常のMarkdownファイルのみ。**

AIコーディングエージェント（Cursor、Claude Code、Copilot、Windsurf）は、コードを生成し直接プロジェクトに反映します。生成されたコードにバグがある場合、アプリケーション全体が障害を起こし、ロールバックの手間とコストが非常に大きくなります。

従来の機能フラグツール（LaunchDarkly、Statsig、Flipt）は、ウェブ登録、APIキー、SDK設定を必要とするため、自立型AIエージェントにはアクセスできません。AIエージェントはブラウザ、メールアドレス、クレジットカードを所有していないためです。

MDFLAGは、機能フラグをプロジェクトディレクトリ内のシンプルなMarkdownファイルにすることで、この問題を解決します。AIエージェントはファイルの作成や編集を得意とするため、フラグをネイティブに扱うことができます。

---

## コアコンセプト

**AIはGitでファイルを作成できますが、SaaSサービスに登録することはできません。**

したがって：
- 従来のフラグサービス（LaunchDarkly） = AIエージェントからは利用不可
- Gitベースのフラグ（MDFLAG） = AIエージェントからネイティブに利用可能

MDFLAGは、複雑な外部連携を教え込もうとするのではなく、AIがすでに得意とする操作を活用する初のソリューションです。

---

## 仕組み

**フラグがない場合:**
AIは古い部屋を解体してから新しい部屋を建てます。新しい部屋が欠陥住宅だった場合、住む場所が完全に失われます。

**フラグがある場合:**
AIは古い部屋の隣に新しい部屋を建てます。どちらに住むかは人間が決定します。新しい部屋に問題があれば、ドアを閉めて古い部屋に住み続ければ安全です。

---

## 主な特徴

- **エージェントネイティブ**: MCP（Model Context Protocol）経由で動作 — ユーザー登録やAPIキーは不要
- **Gitベース**: フラグはリポジトリ内のMarkdownファイルとして管理され、コードと一緒にバージョン管理
- **完全性保護**: SHA-256ハッシュによる改ざん検知で不正な変更を防止
- **段階的ロールアウト**: 5% → 25% → 100% のユーザーへ安全に段階的適用
- **即時ロールバック**: ファイル内の数値を1つ変更するだけで機能を即座に無効化
- **人間によるコントロール**: エージェントがフラグを作成し、人間がロールアウト率を管理

---

## インストール方法

### GitHub Releases からダウンロード

[Releases](https://github.com/legends-deadlines/Agent-First-Markdown-Feature-Flag/releases) ページからお使いのプラットフォーム用の最新バイナリをダウンロードします：

```bash
# Linux/macOS
curl -L https://github.com/legends-deadlines/Agent-First-Markdown-Feature-Flag/releases/latest/download/mdflag-linux-amd64 -o mdflag
chmod +x mdflag
sudo mv mdflag /usr/local/bin/

# インストールの確認
mdflag --help
```

### ソースコードからビルド

```bash
git clone https://github.com/legends-deadlines/Agent-First-Markdown-Feature-Flag.git
cd Agent-First-Markdown-Feature-Flag
go build -o mdflag ./cmd/mdflag/
```

---

## クイックスタート

### 1. フラグの作成

```bash
mdflag create \
  --name new-checkout-flow \
  --percentage 0 \
  --hypothesis "チェックアウトのCVRを5%向上させる" \
  --metrics "conversion_rate,checkout_time_seconds" \
  --author "claude-code-agent" \
  --description "従来の3ステップから1ページ形式のチェックアウトへ刷新"
```

これにより `.mdflag/new-checkout-flow.md` が作成されます：

```markdown
---
name: new-checkout-flow
percentage: 0
targeting: user_id
status: active
created: 2026-09-17T10:00:00Z
author: claude-code-agent
hypothesis: チェックアウトのCVRを5%向上させる
metrics:
    - conversion_rate
    - checkout_time_seconds
agent_section_hash: "sha256:a1b2c3..."
human_section_hash: "sha256:d4e5f6..."
---

## Description
従来の3ステップから1ページ形式のチェックアウトへ刷新
```

### 2. コードでのフラグ利用

```go
import "github.com/legends-deadlines/Agent-First-Markdown-Feature-Flag/pkg/mdflag"

// アプリケーション起動時に一度だけ初期化
client, err := mdflag.New(".mdflag")
if err != nil {
    log.Fatal(err)
}

// リクエストハンドラ内でフラグをチェック
func handleCheckout(w http.ResponseWriter, r *http.Request) {
    userID := getUserID(r)

    if client.Enabled("new-checkout-flow", userID) {
        // 新しいコードパス
        renderOnePageCheckout(w, r)
    } else {
        // 従来の安定版コードパス
        renderThreeStepCheckout(w, r)
    }
}
```

### 3. 段階的ロールアウト

```bash
# 5%のユーザーに有効化
mdflag rollout --name new-checkout-flow --percentage 5

# メトリクスを監視し問題がなければ 25% へ拡大
mdflag rollout --name new-checkout-flow --percentage 25

# 安定が確認できたら全ユーザー（100%）へ全量展開
mdflag rollout --name new-checkout-flow --percentage 100

# 障害検知時は即座に無効化（0%）
mdflag rollout --name new-checkout-flow --percentage 0
```

### 4. 完全性の検証

```bash
mdflag verify
```

出力例：
```
[OK]   .mdflag/new-checkout-flow.md

Results: 1 valid, 0 invalid, 1 total
```

---

## AIエージェントとの連携

MDFLAGはMCP（Model Context Protocol）経由でAIエージェントと接続します。

### Cursorの場合

プロジェクトルートの `.cursor/mcp.json` に追加します：

```json
{
  "mcpServers": {
    "mdflag": {
      "command": "/path/to/mdflag",
      "args": ["serve", "--dir", ".mdflag"]
    }
  }
}
```

### Claude Codeの場合

`.claude/mcp.json` に追加します：

```json
{
  "mcpServers": {
    "mdflag": {
      "command": "/path/to/mdflag",
      "args": ["serve", "--dir", ".mdflag"]
    }
  }
}
```

### Windsurfの場合

`.windsurf/mcp.json` に追加します：

```json
{
  "mcpServers": {
    "mdflag": {
      "command": "/path/to/mdflag",
      "args": ["serve", "--dir", ".mdflag"]
    }
  }
}
```

### エージェントへの指示プロンプト

エージェントのシステムプロンプトや `CLAUDE.md` / `.cursorrules` に記述します：

```
実験的なコード変更を行う場合：
1. コードを書く前に mdflag_create を使用して機能フラグを作成すること
2. 新しいコードを if client.Enabled("flag-name", entityID) { ... } で囲むこと
3. 初期ロールアウト率は 0 に設定すること
4. フラグファイルを直接編集せず、mdflag ツールのみを使用すること
5. フラグファイルとコード変更を一緒にコミットすること
```

---

## 利用可能なMCPツール

| ツール名 | エージェント権限 | 説明 |
|----------|------------------|------|
| `mdflag_create` | 可能 | 新しい機能フラグの作成 |
| `mdflag_list` | 可能 | 全フラグとそのステータスの一覧取得 |
| `mdflag_verify` | Possible | 全フラグの完全性チェック |
| `mdflag rollout` | 不可 | ロールアウト率の変更（CLI限定） |

---

## セキュリティモデル

MDFLAGは3層の保護システムを採用しています：

### 第1層: ハッシュベースの完全性検証

各フラグファイルにはSHA-256ハッシュが含まれます：
- `agent_section_hash`: 作成時にエージェントが設定したフィールドを保護
- `human_section_hash`: ロールアウト率やステータスなどの人間設定フィールドを保護

### 第2層: 権限分離

エージェントには `mdflag_create`, `mdflag_list`, `mdflag_verify` のみ許可されます。`mdflag_rollout` や削除操作の権限はありません。人間のみがCLIを通じて展開率を変更できます。

### 第3層: Git hooks（オプション）

改ざんされたフラグを含むコミットを拒否するPre-commitフックを適用できます。[docs/git-hooks.md](git-hooks.md) を参照。

---

## CLI リファレンス

### `mdflag create`

```bash
mdflag create --name <名前> --hypothesis <仮説> [オプション]
```

### `mdflag rollout`

```bash
mdflag rollout --name <名前> --percentage <0-100>
```

### `mdflag list`

```bash
mdflag list [dir]
```

### `mdflag verify`

```bash
mdflag verify [dir]
```

### `mdflag serve`

```bash
mdflag serve [--dir <path>]
```

---

## ランタイムライブラリ (Go)

```go
package main

import (
    "log"
    "github.com/legends-deadlines/Agent-First-Markdown-Feature-Flag/pkg/mdflag"
)

func main() {
    client, err := mdflag.New(".mdflag")
    if err != nil {
        log.Fatal(err)
    }

    stop, err := client.StartWatcher()
    if err == nil {
        defer stop()
    }

    userID := "user_12345"
    if client.Enabled("new-checkout-flow", userID) {
        // 新コード
    } else {
        // 従来コード
    }
}
```

---

## 決定論的ターゲティング

同じエンティティIDは、指定されたフラグに対して常に同じ結果を受け取ります。FNV-1aハッシュアルゴリズムを使用して均一な分布を実現します。

---

## 開発とテスト

```bash
go build ./...
go test ./... -v
go vet ./...
```

---

## ライセンス

MIT License — 詳細については [LICENSE](../LICENSE) を参照してください。

---

<div align="center">

<img src="../.github/assets/seal.svg" alt="MDFLAG Seal Mascot" width="600"/>

</div>
