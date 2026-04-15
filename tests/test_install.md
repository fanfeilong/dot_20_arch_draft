# d2a Install Test

在 `build` 目录下创建一个独立的测试子目录，验证 `d2a` 的基本功能。

## 测试目录

示例：

```bash
mkdir -p build/install-smoke
```

后续测试都在 `build/install-smoke` 下进行，避免污染仓库根目录。

## 测试步骤

1. 使用 `curl` 安装 `d2a`，建议将安装目录指向测试子目录中的 `bin`：

```bash
mkdir -p build/install-smoke/bin
curl -fsSL https://raw.githubusercontent.com/fanfeilong/dot_20_arch_draft/main/install.sh | D2A_INSTALL_DIR="$PWD/build/install-smoke/bin" sh
```

2. 验证帮助命令可用：

```bash
build/install-smoke/bin/d2a help
```

3. 在测试子目录下创建一个仓库目录并执行 `init`：

```bash
mkdir -p build/install-smoke/project
build/install-smoke/bin/d2a init build/install-smoke/project
```

如需英文语言包：

```bash
build/install-smoke/bin/d2a init build/install-smoke/project-en --lang en
```

4. 检查初始化结果，确认 `.d2a` 工作目录已经生成：

```bash
find build/install-smoke/project -maxdepth 3 -type f | sort
```

重点确认存在：

- `LAB.md`
- `docs/architecture/00_overview.md`
- `docs/implementation/00_mini_scope.md`
- `docs/report/00_report_outline.md`
- `src/README.md`
- `tests/README.md`
- `report/README.md`

5. 检查初始化结果，确认至少生成以下内置 skill：

```bash
find build/install-smoke/project -type f | sort
```

重点确认存在：

- `.codex/skills/d2a-arch-1-project-scope/SKILL.md`
- `.codex/skills/d2a-arch-2-runtime-view/SKILL.md`
- `.codex/skills/d2a-arch-3-core-objects/SKILL.md`
- `.codex/skills/d2a-arch-4-state-evolution/SKILL.md`
- `.codex/skills/d2a-arch-5-module-view/SKILL.md`
- `.codex/skills/d2a-arch-6-tradeoff-view/SKILL.md`
- `.codex/skills/d2a-mini-1-scope/SKILL.md`
- `.codex/skills/d2a-mini-2-design/SKILL.md`
- `.codex/skills/d2a-mini-3-build/SKILL.md`
- `.codex/skills/d2a-mini-4-test/SKILL.md`
- `.codex/skills/d2a-step/SKILL.md`
- `.codex/skills/d2a-report-build/SKILL.md`
- `.codex/skills/d2a-status/SKILL.md`
- `.codex/skills/d2a-challenge-architecture/SKILL.md`

## 通过标准

- `curl` 安装成功，生成可执行文件 `build/install-smoke/bin/d2a`
- `d2a help` 正常输出帮助信息
- `d2a init` 执行成功
- 目标目录下成功生成 `.d2a` 工作目录和模板文件
- 目标目录下成功生成多家 agent 的 `skills` 目录和内置 `d2a-*` skill 文件
