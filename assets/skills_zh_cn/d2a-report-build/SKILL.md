---
name: d2a-report-build
description: 内置 d2a 技能，用于 d2a-report-build 阶段引导与状态更新。
---

# d2a-report-build

## 目标

将当前仓库的 d2a 工作区整理为可本地展示的报告包。

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
3. 确认上下文后，调用 `d2a skill-state d2a-report-build --status started --stage report-prepared --phase analysis-generation --next-step "Refine the report summary and report artifacts." --summary "Started report-build work."`.
4. 将该技能视为 Codex 中报告阶段的用户入口。
5. 若报告数据缺失或过期，在完善报告前调用 `d2a report`。
6. 读取 `.d2a/report/index.md` 与 `.d2a/report/data/*.json`。
7. 使用 `docs/`、`.d2a/src/`、`.d2a/tests/` 作为报告内容来源。
8. 报告聚焦架构、mini 实现、测试与教学叙事。
9. 将 `.d2a/report/data/*.json` 视为未来 Vue 应用的稳定输入契约。
10. 必须生成双页简报产物：`report/brief.md` 与 `report/brief.html`（A4 打印样式）。
11. 收尾前必须显式执行一次 `d2a report`，确保最新报告产物已落盘（不要只停留在对话层结论）。
12. 若任一 DoD 未满足，或 `d2a report` 未成功执行，不得标记 completed。
13. 本轮完成后，调用 `d2a skill-state d2a-report-build --status completed --stage report-ready --phase analysis-generation --next-step "Run d2a serve to open the report." --summary "Completed report-build work and refreshed artifacts via d2a report."`.

## DoD（必须全部满足）

1. 严格双页结构（2 张 A4）：
   - 第 1 页：1 张状态机/架构图 + 6 要素极简表（边界/驱动/核心对象/状态机/模块协作/约束）。
   - 第 2 页：mini 实现简报（技术栈、20%%切片、构建摘要、测试证据、刻意未实现项）。
2. 内容超长时必须压缩，不允许扩展到第 3 页。
3. 报告必须是“可讲解提纲”，禁止长段落灌水。
4. 产物文件必须存在：
   - `report/brief.md`
   - `report/brief.html`
   - `report/index.html`

## 回合结束续接规则

1. 回合结束前必须调用 `d2a skill-state` 持久化当前 phase 与 progress。
2. 回复末尾必须输出：`继续请使用 $d2a-step`。


## 输出（产物）

- `report/brief.md`（双页 A4 结构化简报）
- `report/brief.html`（可打印 A4 双页简报）
- `report/index.md`（总览索引）

## 持久化（.d2a）

- 本技能的推进状态、下一步路由通过 `d2a skill-state` 持久化到 `.d2a/state.json` 与 `.d2a/history.jsonl`。
