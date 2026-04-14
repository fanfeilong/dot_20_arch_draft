# dot_20_arch_draft

`dot_20_arch_draft` 的目标不是只讨论“如何拆解开源项目架构”，而是最终交付一套可直接安装的 `d2a` skills 套装。

这个仓库将使用 Go 开发并发布为一个 CLI 程序：`d2a`。

## What Is d2a

`d2a` 是一个面向开源项目架构拆解的 skills 分发工具。

它的核心目标是降低获取门槛：

- 用户不需要手动拷贝 prompt 或模板
- 用户不需要自己维护多套 agent 配置
- 用户只需要执行一个初始化命令，就可以把内置的 `d2a-*` skills 安装到目标目录

## Product Goal

通过 `d2a`，用户可以在一个项目目录中快速得到一套用于“架构拆解”的内置 skills，并在 Codex、Claude、Cursor、OpenCode、Trae、NeoCode 等目录约定下复用。

这些 skills 最终服务于一个统一目标：

- 拿一个开源项目
- 用一组原子化、极简的提问步骤逐步拆解其架构
- 将输出控制得足够小
- 在大约一小时内，向学生和 IT 工程师展示如何理解一个项目的架构内核

## CLI

当前计划先保持极小接口：

```text
d2a help
d2a init <target-dir>
d2a version
```

## Install

发布 GitHub Release 后，用户可以直接安装：

```bash
curl -fsSL https://raw.githubusercontent.com/fanfeilong/dot_20_arch_draft/main/install.sh | sh
```

也可以安装指定版本：

```bash
curl -fsSL https://raw.githubusercontent.com/fanfeilong/dot_20_arch_draft/main/install.sh | D2A_VERSION=v0.0.1 sh
```

默认安装到 `/usr/local/bin/d2a`。

可选环境变量：

- `D2A_VERSION`: 指定 Release 版本，默认 `latest`
- `D2A_INSTALL_DIR`: 指定安装目录，默认 `/usr/local/bin`
- `D2A_REPO`: 指定 GitHub 仓库，默认 `fanfeilong/dot_20_arch_draft`

## init Command

`d2a init <target-dir>` 会在目标目录下初始化多家 agent 的 `skills` 目录，并安装内置的一组 `d2a-*` skills。

当前覆盖的目录：

- `.codex/skills`
- `.claude/skills`
- `.cursor/skills`
- `.opencode/skills`
- `.trae/skills`
- `.neocode/skills`

安装后的结构示例：

```text
<target-dir>/
  .codex/skills/d2a-project-scope/SKILL.md
  .codex/skills/d2a-runtime-view/SKILL.md
  .codex/skills/d2a-module-view/SKILL.md
  .claude/skills/d2a-project-scope/SKILL.md
  ...
```

## Built-in Skills

当前脚手架先内置一组最小 skills：

- `d2a-project-scope`
- `d2a-runtime-view`
- `d2a-module-view`

这些 skills 只是第一版骨架，后续会继续扩展为完整套装。

## Development Direction

这个仓库会沿着两个方向同时推进：

1. 打磨 `d2a` CLI 本身，让安装和分发足够简单。
2. 持续收敛 `d2a-*` skills 的内容，让“架构拆解”流程真正可复用。

## Current Status

当前阶段完成的是第一版项目脚手架：

- Go CLI 入口
- `help` 与 `init` 命令
- 内置 skills 资源
- 多目录初始化逻辑

下一步会继续补：

- skills 的版本管理
- 更完整的 `d2a-*` skills 套装
- 面向真实开源项目的演示样例
- 更稳定的发布与安装方式

## Release

仓库内置了 GitHub Actions release workflow。

当推送形如 `v0.0.1` 的 tag 时，workflow 会：

- 运行测试
- 构建 macOS 和 Linux 二进制
- 打包为 Release 资产
- 创建对应的 GitHub Release
