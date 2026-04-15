# dot_20_arch_draft

`dot_20_arch_draft` 的目标不是只讨论“如何拆解开源项目架构”，而是最终交付一套可直接安装的 `d2a` skills 套装。

这个仓库将使用 Go 开发并发布为一个 CLI 程序：`d2a`。

## What Is d2a

`d2a` 是一个面向开源项目架构拆解的 workspace-root 工作流工具。

它的核心目标是降低获取门槛：

- 用户不需要手动拷贝 prompt 或模板
- 用户不需要自己维护多套 agent 配置
- 用户只需要执行一个初始化命令，就可以创建独立的 `<target>_d2a` 工作区，并安装内置的 `d2a-*` skills

## Product Goal

通过 `d2a`，用户可以快速得到一个用于“架构拆解”的独立工作区目录，其中包含内置 skills、分析文档模板、实现目录、测试目录、报告目录和 `repos/` 输入仓库目录，并在 Codex、Claude、Cursor、OpenCode、Trae、NeoCode 等目录约定下复用。

这些 skills 最终服务于一个统一目标：

- 拿一个开源项目
- 用一组原子化、极简的提问步骤逐步拆解其架构
- 将输出控制得足够小
- 在大约一小时内，向学生和 IT 工程师展示如何理解一个项目的架构内核

## CLI

当前计划先保持极小接口：

```text
d2a help
d2a -U
d2a init <target-repo-git-url> [--lang <zh|en>]
d2a analyze [<target-repo>] [--repo <repo-dir>]
d2a derive-mini [--repo <repo-dir>] [--skip-challenge-reason <text>]
d2a test-mini [--repo <repo-dir>]
d2a report [--repo <repo-dir>]
d2a serve [--repo <repo-dir>]
d2a status [--repo <repo-dir>]
d2a skill-state <skill-name> [--repo <repo-dir>] [--status <started|progress|completed>] [--stage <stage>] [--phase <phase>] [--question-index <n>] [--question-total <n>] [--next-step <text>] [--next-skill <name>] [--next-file <path>] [--decision <label>] [--strength <strong|partial|weak>] [--recommendation <proceed|review|revisit architecture>] [--objection <text>] [--summary <text>]
d2a version
```

快速自更新到 GitHub 最新版本：

```bash
d2a -U
```

## Install

发布 GitHub Release 后，用户可以直接安装：

```bash
curl -fsSL https://raw.githubusercontent.com/fanfeilong/dot_20_arch_draft/main/install.sh | sh
```

Windows PowerShell:

```powershell
irm https://raw.githubusercontent.com/fanfeilong/dot_20_arch_draft/main/install.ps1 | iex
```

也可以安装指定版本：

```bash
curl -fsSL https://raw.githubusercontent.com/fanfeilong/dot_20_arch_draft/main/install.sh | D2A_VERSION=v0.0.1 sh
```

```powershell
$env:D2A_VERSION="v0.0.1"
irm https://raw.githubusercontent.com/fanfeilong/dot_20_arch_draft/main/install.ps1 | iex
```

默认安装到 `/usr/local/bin/d2a`。

Windows 默认安装到 `$env:LOCALAPPDATA\d2a\bin\d2a.exe`。

可选环境变量：

- `D2A_VERSION`: 指定 Release 版本，默认 `latest`
- `D2A_INSTALL_DIR`: 指定安装目录，默认 `/usr/local/bin`
- `D2A_REPO`: 指定 GitHub 仓库，默认 `fanfeilong/dot_20_arch_draft`
- `D2A_BASE_URL`: 指定 Release 资产下载基址，主要用于 CI 或自托管镜像测试

## init Command

`d2a init <target-repo-git-url>` 会在当前目录创建 `<target-repo-name>_d2a` 工作区，并将目标仓库浅克隆到 `repos/<target-repo-name>/`。

语言包：

- `--lang zh`：安装中文 skill 包（默认）
- `--lang en`：安装英文 skill 包

当前会创建：

- 多家 agent 的 `skills` 目录
- `docs/architecture`
- `docs/implementation`
- `docs/report`
- `src`
- `tests`
- `report`
- `repos/<target-repo-name>`（浅克隆）
- `LAB.md`
- `.d2a/state.json`、`.d2a/history.jsonl`（状态元数据）

在 d2a 初始化完成后，可以继续运行：

```bash
d2a analyze [<target-repo>] [--repo <repo-dir>]
```

当前这一步会：

- 将目标仓库记录到 `<repo-dir>/.d2a/target.json`（`init` 已写入，`analyze` 可覆盖）
- 将 `docs/architecture/*.md` 改写成带目标仓库、主 skill、原子问题和输出约束的分析入口文件
- 为后续在 AI Coding 工具中逐步填写 architecture docs 提供稳定落点

在 architecture docs 有了稳定落点后，可以继续运行：

```bash
d2a derive-mini [--repo <repo-dir>] [--skip-challenge-reason <text>]
```

当前这一步会：

- 默认要求 challenge phase 已经完成
- 将 `docs/implementation/*.md` 改写成面向 mini implementation 的任务入口文件
- 为 mini 版本的 scope、design、build plan、test plan 建立固定文档落点
- 生成 `src/ARCHITECTURE.md`，作为后续 `src/` 实现的总纲入口
- 生成首个可运行栈样板 `src/go-mini/`（Go）

如果确实需要跳过 challenge phase，必须显式传：

```bash
d2a derive-mini --skip-challenge-reason "..."
```

这个跳过理由会被记录到 challenge 状态文件中，而不是静默绕过。

在 mini implementation 规划完成后，可以继续运行：

```bash
d2a test-mini [--repo <repo-dir>]
```

当前这一步会：

- 将 `tests/README.md` 改写成面向 mini implementation 的测试总入口
- 生成 `tests/01_integration_tasks.md`，用于定义第一个端到端测试和后续集成场景
- 生成 `.d2a/test-mini.json`，明确测试阶段的输入和输出文件

在测试规划完成后，可以继续运行：

```bash
d2a report [--repo <repo-dir>]
```

当前这一步会：

- 生成 `report/data/summary.json`
- 生成 `report/data/target.json`
- 生成 `report/data/tests.json`
- 生成 `report/data/challenge.json`
- 生成 `report/index.md`
- 生成 `report/index.html`
- 确保 `report/vue-app/` Vue 骨架存在（缺失时自动补齐）

当前这一步的目标是建立稳定数据接口，并提供两个展示层：

- 默认展示层：`index.html` 运行时直接加载 `./data/*.json`，不依赖 Node
- 可选开发层：`vue-app/` 用于后续前端迭代

在报告产物生成后，可以继续运行：

```bash
d2a serve [--repo <repo-dir>]
```

当前这一步会：

- 校验 `report/index.html` 已存在
- 启动本地静态服务
- 暴露 `report/index.html` 与 `report/data/*.json`

当前这一步提供的是一个最小可浏览的本地报告页，为后续 Vue 报告应用留出稳定路径和数据接口。

还可以随时运行：

```bash
d2a status [--repo <repo-dir>]
```

当前这一步会：

- 读取 `.d2a/state.json`
- 读取 `.d2a/history.jsonl`
- 输出当前阶段、上一个命令、下一个建议动作以及最近历史摘要

另外还提供一个偏内部用途的命令：

```bash
d2a skill-state <skill-name> ...
```

它主要供 `d2a-*` skills 在 Codex 中调用，用来记录：

- 当前 skill
- 当前 phase
- question progress
- next step / next skill / next file
- skill 级历史事件
- challenge decision / strength / recommendation / objection

当前覆盖的目录：

- `.codex/skills`
- `.claude/skills`
- `.cursor/skills`
- `.opencode/skills`
- `.trae/skills`
- `.neocode/skills`

安装后的结构示例：

```text
<repo-dir>/
  .codex/skills/d2a-arch-1-project-scope/SKILL.md
  .codex/skills/d2a-arch-2-runtime-view/SKILL.md
  .codex/skills/d2a-arch-3-core-objects/SKILL.md
  .codex/skills/d2a-arch-4-state-evolution/SKILL.md
  .codex/skills/d2a-arch-5-module-view/SKILL.md
  .codex/skills/d2a-arch-6-tradeoff-view/SKILL.md
  .codex/skills/d2a-mini-1-scope/SKILL.md
  .codex/skills/d2a-mini-2-design/SKILL.md
  .codex/skills/d2a-mini-3-build/SKILL.md
  .codex/skills/d2a-mini-4-test/SKILL.md
  .codex/skills/d2a-step/SKILL.md
  .codex/skills/d2a-report-build/SKILL.md
  .codex/skills/d2a-status/SKILL.md
  .codex/skills/d2a-challenge-architecture/SKILL.md
  LAB.md
  docs/architecture/00_overview.md
  docs/implementation/00_mini_scope.md
  docs/report/00_report_outline.md
  src/README.md
  tests/README.md
  report/README.md
  .claude/skills/d2a-arch-1-project-scope/SKILL.md
```

## Built-in Skills

当前脚手架先内置一组最小 skills：

- `d2a-arch-1-project-scope`: 对应六要素中的“边界”
- `d2a-arch-2-runtime-view`: 对应六要素中的“驱动”
- `d2a-arch-3-core-objects`: 对应六要素中的“核心对象”
- `d2a-arch-4-state-evolution`: 对应六要素中的“状态演化 / 状态机”
- `d2a-arch-5-module-view`: 对应六要素中的“协作”
- `d2a-arch-6-tradeoff-view`: 对应六要素中的“约束 / 取舍”
- `d2a-mini-1-scope`: mini implementation 范围收敛入口
- `d2a-mini-2-design`: mini implementation 设计入口
- `d2a-mini-3-build`: mini implementation 实现入口
- `d2a-mini-4-test`: mini integration testing 入口
- `d2a-step`: 状态驱动推进入口，自动选择下一子 skill，并支持会话中断后的断点续接
- `d2a-report-build`: report 汇总与展示入口
- `d2a-status`: 当前 d2a 工作流状态查看入口
- `d2a-challenge-architecture`: 分析完成后的架构质疑阶段入口，负责记录挑战而不默认改写架构结论

这些 skills 只是第一版骨架，后续会继续扩展为完整套装。

推荐默认使用方式：

1. 在 Codex 中先调用 `$d2a-step`
2. 每轮回答后继续调用 `$d2a-step`

当前其中这两个 skill 已作为第一批双阶段 skill 样板，明确区分：

- [d2a-arch-1-project-scope](/Users/feilong/Desktop/dev/zigslice/dot_20_arch_draft/assets/skills_zh_cn/d2a-arch-1-project-scope/SKILL.md)
- [d2a-arch-2-runtime-view](/Users/feilong/Desktop/dev/zigslice/dot_20_arch_draft/assets/skills_zh_cn/d2a-arch-2-runtime-view/SKILL.md)
- [d2a-arch-3-core-objects](/Users/feilong/Desktop/dev/zigslice/dot_20_arch_draft/assets/skills_zh_cn/d2a-arch-3-core-objects/SKILL.md)
- [d2a-arch-4-state-evolution](/Users/feilong/Desktop/dev/zigslice/dot_20_arch_draft/assets/skills_zh_cn/d2a-arch-4-state-evolution/SKILL.md)
- [d2a-arch-5-module-view](/Users/feilong/Desktop/dev/zigslice/dot_20_arch_draft/assets/skills_zh_cn/d2a-arch-5-module-view/SKILL.md)
- [d2a-arch-6-tradeoff-view](/Users/feilong/Desktop/dev/zigslice/dot_20_arch_draft/assets/skills_zh_cn/d2a-arch-6-tradeoff-view/SKILL.md)

- analysis-generation
- confirmation-questions

并要求通过 `d2a skill-state` 写入 question progress 和下一步建议。

mini 主链当前也开始按同样方式推进，已升级的样板包括：

- [d2a-mini-1-scope](/Users/feilong/Desktop/dev/zigslice/dot_20_arch_draft/assets/skills_zh_cn/d2a-mini-1-scope/SKILL.md)
- [d2a-mini-2-design](/Users/feilong/Desktop/dev/zigslice/dot_20_arch_draft/assets/skills_zh_cn/d2a-mini-2-design/SKILL.md)
- [d2a-mini-3-build](/Users/feilong/Desktop/dev/zigslice/dot_20_arch_draft/assets/skills_zh_cn/d2a-mini-3-build/SKILL.md)
- [d2a-mini-4-test](/Users/feilong/Desktop/dev/zigslice/dot_20_arch_draft/assets/skills_zh_cn/d2a-mini-4-test/SKILL.md)

## Development Direction

这个仓库会沿着两个方向同时推进：

1. 打磨 `d2a` CLI 本身，让安装和分发足够简单。
2. 持续收敛 `d2a-*` skills 的内容，让“架构拆解”流程真正可复用。

## Current Status

当前阶段完成的是第一版项目脚手架：

- Go CLI 入口
- `help`、`init`、`analyze`、`derive-mini`、`test-mini`、`report` 与 `serve` 命令
- `status` 命令
- 内置 skills 资源
- workspace-root 独立工作区初始化逻辑（含 `repos/` 浅克隆）
- analysis task 文件生成逻辑
- mini derivation task 文件生成逻辑
- test planning task 文件生成逻辑
- report data 与 report index 生成逻辑
- 本地静态报告服务逻辑

下一步会继续补：

- skills 的版本管理
- 更完整的 `d2a-*` skills 套装
- 面向真实开源项目的演示样例
- 更稳定的发布与安装方式

## Release

仓库内置了 GitHub Actions release workflow。

当推送形如 `v0.0.1` 的 tag 时，workflow 会：

- 运行测试
- 构建 macOS、Linux 和 Windows 二进制
- 打包为 Release 资产
- 创建对应的 GitHub Release
