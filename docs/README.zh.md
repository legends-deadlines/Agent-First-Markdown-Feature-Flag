<div align="center">

<img src="../.github/assets/banner.svg?v=1.0.2" alt="Agent-First Feature Flags" width="100%"/>

<br/>

[![Status](https://img.shields.io/badge/🚀%20状态-v1.0.2%20就绪-10b981?style=flat-square&labelColor=1c2128)](#)
[![开源协议](https://img.shields.io/badge/协议-MIT-79c0ff?style=flat-square&labelColor=1c2128)](../LICENSE)
[![Go](https://img.shields.io/badge/Go-1.27+-00ADD8?style=flat-square&logo=go&labelColor=1c2128)](#)

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

## 安装说明

### 从 GitHub Releases 下载

从 [Releases](https://github.com/legends-deadlines/Agent-First-Markdown-Feature-Flag/releases) 下载适用于您平台的最新可执行文件：

```bash
# Linux/macOS
curl -L https://github.com/legends-deadlines/Agent-First-Markdown-Feature-Flag/releases/latest/download/mdflag-linux-amd64 -o mdflag
chmod +x mdflag
sudo mv mdflag /usr/local/bin/

# 验证安装
mdflag --help
```

### 从源码编译

```bash
git clone https://github.com/legends-deadlines/Agent-First-Markdown-Feature-Flag.git
cd Agent-First-Markdown-Feature-Flag
go build -o mdflag ./cmd/mdflag/
```

---

## 快速开始

### 1. 创建标记

```bash
mdflag create \
  --name new-checkout-flow \
  --percentage 0 \
  --hypothesis "提高 5% 的结账转化率" \
  --metrics "conversion_rate,checkout_time_seconds" \
  --author "claude-code-agent" \
  --description "用单页结账流程替换现有的 3 步流程"
```

这将在 `.mdflag/new-checkout-flow.md` 创建文件：

```markdown
---
name: new-checkout-flow
percentage: 0
targeting: user_id
status: active
created: 2026-09-17T10:00:00Z
author: claude-code-agent
hypothesis: 提高 5% 的结账转化率
metrics:
    - conversion_rate
    - checkout_time_seconds
agent_section_hash: "sha256:a1b2c3..."
human_section_hash: "sha256:d4e5f6..."
---

## Description
用单页结账流程替换现有的 3 步流程
```

### 2. 在代码中使用标记

```go
import "github.com/legends-deadlines/Agent-First-Markdown-Feature-Flag/pkg/mdflag"

// 在应用启动时初始化客户端
client, err := mdflag.New(".mdflag")
if err != nil {
    log.Fatal(err)
}

// 在请求处理函数中检查标记状态
func handleCheckout(w http.ResponseWriter, r *http.Request) {
    userID := getUserID(r)

    if client.Enabled("new-checkout-flow", userID) {
        // 新代码分支
        renderOnePageCheckout(w, r)
    } else {
        // 现有稳定分支
        renderThreeStepCheckout(w, r)
    }
}
```

### 3. 渐进式发布

```bash
# 为 5% 的用户启用
mdflag rollout --name new-checkout-flow --percentage 5

# 观察指标无误后调高至 25%
mdflag rollout --name new-checkout-flow --percentage 25

# 运行正常后为所有用户全量启用 (100%)
mdflag rollout --name new-checkout-flow --percentage 100

# 如遇异常立即禁用 (0%)
mdflag rollout --name new-checkout-flow --percentage 0
```

### 4. 完整性校验

```bash
mdflag verify
```

预期输出：
```
[OK]   .mdflag/new-checkout-flow.md

Results: 1 valid, 0 invalid, 1 total
```

---

## AI 代理集成

MDFLAG 通过 MCP (Model Context Protocol) 协议与 AI 代理无缝集成。

### 针对 Cursor

在项目根目录下的 `.cursor/mcp.json` 添加：

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

### 针对 Claude Code

在 `.claude/mcp.json` 添加：

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

### 针对 Windsurf

在 `.windsurf/mcp.json` 添加：

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

### 代理指令说明

在您的 AI 代理系统提示词或 `CLAUDE.md` / `.cursorrules` 中加入：

```
在编写实验性代码更改时：
1. 在编写代码之前，必须先使用 mdflag_create 创建特性标记。
2. 使用 if client.Enabled("flag-name", entityID) { ... } 包裹新代码。
3. 将初始发布比例 (percentage) 设置为 0。
4. 切勿直接修改标记 Markdown 文件 — 只能使用 mdflag 提供的工具。
5. 将标记文件和代码修改一同提交至 Git。
```

---

## 可用的 MCP 工具

| 工具名称 | 代理权限 | 功能描述 |
|----------|----------|----------|
| `mdflag_create` | 可用 | 创建新的特性标记 |
| `mdflag_list` | 可用 | 列出所有标记及其当前状态 |
| `mdflag_verify` | 可用 | 校验所有标记的哈希完整性 |
| `mdflag rollout` | 禁用 | 调整标记发布比例（仅限 CLI 命令） |

此权限隔离机制确保 AI 代理可以自由创建标记，但无法自行开启标记或篡改已有规则。

---

## 安全架构模型

MDFLAG 采用三层安全防护机制：

### 第一层：基于哈希的完整性保护

每个标记文件均包含由 SHA-256 计算的内容校验和：
- `agent_section_hash`：保护代理创建时填写的字段。
- `human_section_hash`：保护控制开关与比例的关键字段。

未经 CLI 工具擅自修改文件内容将导致哈希比对失败，`mdflag verify` 会立即报错。

### 第二层：工具权限隔离

AI 代理仅获得 `mdflag_create`、`mdflag_list`、`mdflag_verify` 的授权，无权调用 `mdflag_rollout` 或执行删除。发布比例的变更全权由人类掌控。

### 第三层：Git 钩子保护 (可选)

可启用 Pre-commit 钩子拦截包含篡改标记的提交，详情参见 [docs/git-hooks.md](git-hooks.md)。

---

## CLI 命令参考

### `mdflag create`

创建新的特性标记。

```bash
mdflag create --name <名称> --hypothesis <假设> [参数]

选项参数:
  --name          标记名称 (必填, kebab-case)
  --hypothesis    实验假设 (必填)
  --percentage    初始比例 0-100 (默认: 0)
  --targeting     目标维度: user_id, session_id, random (默认: user_id)
  --author        作者名称
  --metrics       评估指标 (多个用逗号隔开)
  --expires       过期时间 (RFC3339 格式)
  --description   人类可读描述
  --dir           标记存储目录 (默认: .mdflag)
```

### `mdflag rollout`

调整特性标记的启用比例。

```bash
mdflag rollout --name <名称> --percentage <0-100>

选项参数:
  --name          标记名称 (必填)
  --percentage    新比例 0-100 (必填)
  --dir           标记存储目录 (默认: .mdflag)
```

### `mdflag list`

列出指定目录下的所有标记。

```bash
mdflag list [目录]
```

### `mdflag verify`

检查所有标记文件的完整性。

```bash
mdflag verify [目录]
```

### `mdflag serve`

启动面向 AI 代理的 MCP 服务端。

```bash
mdflag serve [--dir <路径>]
```

---

## 运行时 SDK

Go 语言 SDK 用于在应用内部实时评估标记。

```go
package main

import (
    "log"
    "github.com/legends-deadlines/Agent-First-Markdown-Feature-Flag/pkg/mdflag"
)

func main() {
    // 启动时初始化
    client, err := mdflag.New(".mdflag")
    if err != nil {
        log.Fatal(err)
    }

    // 可选：开启文件热加载
    stop, err := client.StartWatcher()
    if err == nil {
        defer stop()
    }

    // 在请求处理中使用
    userID := "user_12345"
    if client.Enabled("new-checkout-flow", userID) {
        // 新功能分支
    } else {
        // 旧功能分支
    }
}
```

### 确定性分发 (Deterministic Targeting)

同一实体 ID 对同一标记的评估结果始终一致：
- 用户在不同请求间不会遭遇界面闪烁
- A/B 测试流量分组稳定
- 问题排查可精准复现

### 比例算法

对于 1 至 99 之间的比例，MDFLAG 采用 FNV-1a 哈希算法：

```
bucket = FNV-1a(flagName + ":" + entityID) % 100
enabled = bucket < percentage
```

确保流量分布均匀可靠。

---

## 标记生命周期

1. **创建标记** (percentage: 0) — AI 代理创建标记并包裹新代码
2. **代码评审** — 开发者审查 PR 及标记文件
3. **渐进发布** — 开发者逐步提高比例：5% → 25% → 100%
4. **指标监控** — 观察标记设定的衡量指标
5. **清理归档** — 二选一：
   - 移除标记并保留新代码 (成功)
   - 比例重置为 0 并回滚代码 (失败)

---

## 方案对比

| 工具名称 | 工作原理 | AI 代理面临的痛点 |
|----------|----------|-------------------|
| LaunchDarkly / Statsig | SaaS 云服务，需 API 密钥 | AI 代理无法自动注册账号 |
| Flipt v2 | GitOps 模式，但需部署本地服务器 | AI 代理无法自行启动服务器 |
| dif.sh | Markdown 标记 | 缺乏针对 AI 篡改的保护机制 |
| **MDFLAG** | Markdown 标记 + 哈希校验保护 | **全自主无缝运行** |

---

## 设计哲学

**不要试图让 AI 变得万无一失。** AI 偶尔写出有缺陷的代码是不可避免的。

**正确的做法是：** 构建一套安全机制，确保 AI 的缺陷代码绝不会摧毁生产系统。隔离 AI 代码，渐进式开启，秒级回滚。

---

## 市场验证

- **dif.sh** 在 Product Hunt 登顶榜首 (#1) — 证明了该方向的强劲需求
- **GitClear 研究**：包含 AI 介入的项目代码流失率 (churn) 升至 28-40%
- **46.4%** 由 AI 代理生成的 PR 在评审环节被驳回

---

## 编译与开发

### 编译

```bash
go build ./...
```

### 测试

```bash
go test ./... -v
```

### 静态检查

```bash
go vet ./...
gofmt -l .
```

---

## 开源协议

MIT 协议 — 详情参阅 [LICENSE](../LICENSE)。

---

## 路线图 Roadmap

- [x] 核心库：文件格式、解析器、哈希计算、校验器
- [x] CLI 工具集：create, rollout, verify, list
- [x] 确定性分发 SDK 运行时
- [x] 面向 AI 代理的 MCP 服务端
- [ ] Git Hooks 提交拦截保护
- [ ] 多语言 SDK 支持 (TypeScript, Python)
- [ ] 标记数据分析与指标集成
- [ ] 团队协作：CODEOWNERS、CI 检查
- [ ] 网页版 Flag 管理控制台

---

<div align="center">

<img src="../.github/assets/seal.svg" alt="MDFLAG Seal Mascot" width="600"/>

</div>
