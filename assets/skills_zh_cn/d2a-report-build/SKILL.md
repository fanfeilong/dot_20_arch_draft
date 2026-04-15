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

如果无法确定当前仓库，立即停止并询问用户要使用哪个仓库。

## 执行说明

1. 先确认当前仓库上下文。将 repo/path 信息写入统一外壳，不要额外输出独立头信息清单。
2. 如果无法确定当前仓库，立即停止并询问用户要使用哪个仓库。
3. 确认上下文后，调用 `d2a skill-state d2a-report-build --status started --stage report-prepared --phase analysis-generation --next-step "Refine the report summary and report artifacts." --summary "Started report-build work."`.
4. 将该技能视为 Codex 中报告阶段的用户入口。
5. 若报告数据缺失或过期，在完善报告前调用 `d2a report`。
6. 读取 `.d2a/report/index.md` 与 `.d2a/report/data/*.json`。
7. 使用 `.d2a/docs/`、`.d2a/src/`、`.d2a/tests/` 作为报告内容来源。
8. 报告聚焦架构、mini 实现、测试与教学叙事。
9. 将 `.d2a/report/data/*.json` 视为未来 Vue 应用的稳定输入契约。
10. 可完善 `.d2a/report/index.md` 或后续报告资产，但不要改变阶段契约。
11. 本轮完成后，调用 `d2a skill-state d2a-report-build --status completed --stage report-ready --phase analysis-generation --next-step "Review the local report or run d2a serve." --summary "Completed report-build work."`.

## 回合结束续接规则

1. 回合结束前必须调用 `d2a skill-state` 持久化当前 phase 与 progress。
2. 回复末尾必须输出：`继续请使用 $d2a-step`。


## 输出

- 报告摘要
- 报告数据一致性
- 下一步展示改进点
