<div align="center">

<img src="../.github/assets/banner.svg" alt="Agent-First Feature Flags" width="100%"/>

<br/>

[English](../README.md) | [Русский](README.ru.md) | [Español](README.es.md) | [中文](README.zh.md) | [日本語](README.ja.md) | [Français](README.fr.md) | [Deutsch](README.de.md)

</div>

---

# MDFLAG

**专为 AI 代理设计的特性标记（Feature Flags）。无需 SaaS，无需 API 密钥。只需 Git 仓库中的 Markdown 文件。**

AI 编程代理（Cursor、Claude Code、Copilot、Windsurf）直接在您的项目中生成和修改代码。如果代码包含缺陷，整个应用就会崩溃，回滚过程既昂贵又缓慢。

传统的特性标记工具（LaunchDarkly、Statsig、Flipt）对 AI 代理并不友好 — 它们需要网页注册、API 密钥和复杂的 SDK 配置。但 AI 代理没有浏览器、电子邮箱或信用卡。

MDFLAG 将特性标记简化为项目目录中的普通 Markdown 文件，完美解决了这一难题。AI 代理天生擅长创建和编辑文件，因此可以无缝使用特性标记。

---

## 核心理念

**AI 可以在 Git 中创建文件，但无法注册 SaaS 服务。**

因此：
- 传统标记服务（LaunchDarkly） = AI 代理无法访问
- 基于 Git 的标记（MDFLAG） = AI 代理天然原生支持

MDFLAG 是首个立足于 AI 现有能力、而非试图强教其复杂外部集成的解决方案。

---

## 工作原理

**没有特性标记时：**
AI 通过拆除旧房间来建造新房间。如果新房间建歪了，您将无房可住。

**使用特性标记时：**
AI 在旧房间旁建造新房间。由您决定住在哪里。如果新房间出现问题，只需关上房门并继续住在旧房间中。

---

## 主要特性

- **代理原生**：通过 MCP（Model Context Protocol）工作 — 无需注册，无需 API 密钥
- **基于 Git**：标记以 Markdown 文件形式保存在仓库中，与代码一并进行版本控制
- **完整性保护**：基于 SHA-256 哈希的篡改检测，防止未授权修改
- **渐进式发布**：按 5% → 25% → 100% 用户比例逐步启用标记
- **即时回滚**：修改文件中的单个数字即可立即禁用该功能
- **人类掌控**：AI 代理创建标记，人类管理者控制发布进度

---

## 快速开始 (CLI)

```bash
# 创建标记 (初始 0% 发布)
mdflag create --name new-checkout --hypothesis "提高 5% 的转化率"

# 调整发布比例 (25%)
mdflag rollout --name new-checkout --percentage 25

# 验证哈希完整性
mdflag verify

# 列出所有标记
mdflag list
```

---

## 代码中使用 (Go)

```go
import "github.com/legends-deadlines/mdflag/pkg/mdflag"

client, err := mdflag.New(".mdflag")
if err != nil {
    log.Fatal(err)
}

if client.Enabled("new-checkout", userID) {
    // 新功能逻辑
} else {
    // 现有稳定逻辑
}
```

---

## 开源协议

MIT 协议 — 详情请参阅 [LICENSE](../LICENSE)。

---

<div align="center">

<img src="../.github/assets/seal.svg" alt="MDFLAG Seal Mascot" width="600"/>

</div>
