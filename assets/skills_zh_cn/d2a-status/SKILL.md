---
name: d2a-status
description: 内置 d2a 技能，用于 d2a-status 阶段引导与状态更新。
---

# d2a-status

## 目标

从 `.d2a/state.json` 与 `.d2a/history.jsonl` 重新展示当前 d2a 工作流状态。

## 必需输出格式（与 d2a-step 一致外壳）

所有非 `d2a-step` 技能回复都必须复用与 `d2a-step` 一致的外壳格式：

```text
==================================================
【<一级阶段> <N/总数>｜<二级步骤>】<repo>
next: <next_skill> → <next_file>
==================================================

<本技能正文>

--------------------------------------------------
done: <本轮完成动作>
state: <当前骨架位置> → <下一骨架位置> · 继续请使用 $d2a-step
--------------------------------------------------
```

规则：

1. 开场区（分割线之间）必须严格两行。
2. 结尾区（分割线之间）必须严格两行。
3. 不再输出旧的多行头信息清单（repo/path/stage/flow 逐行罗列）。
4. 若 `next_file` 未知，输出 `unknown`，且保持单行格式不变。


## 正文排版硬规则

1. 正文必须使用 `- ` 列表输出，最少 2 条、最多 4 条。
2. 每条要点独占一行，不得写成长段落。
3. 单行不超过 80 个中文字符；超长必须拆行。
4. 正文禁止使用 Markdown 强调符号（如 `` `...` ``、`**...**`）。

5. 对条目、编号列表、选择题选项（A/B/C/D）等结构化内容，不要为满足字数限制而合并或破坏结构；保持一项一行。

如果无法确定当前仓库，立即停止并询问用户要使用哪个仓库。

## Human In Loop 标记规则

当本回合包含“向用户提问并等待回答”的动作时，回复正文最后一行必须追加：

`[human_in_loop]`

## 执行说明

1. 先确认当前仓库上下文。将 repo/path 信息写入统一外壳，不要额外输出独立头信息清单。
2. 如果无法确定当前仓库，立即停止并询问用户要使用哪个仓库。
3. 确认上下文后，调用 `d2a skill-state d2a-status --status started --phase analysis-generation --next-step "Read the latest state and recent history." --summary "Started status 复审."`.
4. 读取 `.d2a/state.json`。
5. 读取 `.d2a/history.jsonl` 的最近记录。
6. 重述当前阶段、最近命令、最近技能与下一步建议。
7. 摘要保持简短且可执行。
8. 若状态文件缺失或过期，告知用户应先执行的命令。
9. 状态摘要完成后，调用 `d2a skill-state d2a-status --status completed --phase analysis-generation --summary "Completed status 复审."`.

## 回合结束续接规则

1. 回合结束前必须调用 `d2a skill-state` 持久化当前 phase 与 progress。
2. 回复末尾必须输出：`继续请使用 $d2a-step`。


## 输出

- 当前阶段
- 最近命令
- 最近技能
- 下一步
- 最近历史摘要
