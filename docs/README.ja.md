<div align="center">

<img src="../.github/assets/banner.svg" alt="Agent-First Feature Flags" width="100%"/>

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

## クイックスタート (CLI)

```bash
# フラグの作成 (初期値 0% ロールアウト)
mdflag create --name new-checkout --hypothesis "コンバージョン率を5%向上"

# ロールアウト率の変更 (25%)
mdflag rollout --name new-checkout --percentage 25

# ハッシュの完全性を検証
mdflag verify

# 全フラグの一覧表示
mdflag list
```

---

## コードでの使用例 (Go)

```go
import "github.com/legends-deadlines/mdflag/pkg/mdflag"

client, err := mdflag.New(".mdflag")
if err != nil {
    log.Fatal(err)
}

if client.Enabled("new-checkout", userID) {
    // 新機能のコードパス
} else {
    // 既存の安定版コードパス
}
```

---

## ライセンス

MIT License — 詳細については [LICENSE](../LICENSE) を参照してください。

---

<div align="center">

<img src="../.github/assets/seal.svg" alt="MDFLAG Seal Mascot" width="600"/>

</div>
