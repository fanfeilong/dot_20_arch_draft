---
name: d2a-challenge-architecture
description: 内置 d2a 技能，用于 d2a-challenge-architecture 阶段引导与状态更新。
---

# d2a-challenge-architecture

## 目标

该技能用于当前阶段的中文引导、状态推进与教学问答。

## 语言规则

所有面向用户的文本都必须使用简体中文。

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
3. 单行不超过 100 个中文字符；超长必须拆行。
4. 正文禁止使用 Markdown 强调符号（如 `` `...` ``、`**...**`）。

如果无法确定当前仓库，立即停止并询问用户要使用哪个仓库。

## Human In Loop 标记规则

当本回合包含“向用户提问并等待回答”的动作时，回复正文最后一行必须追加：

`[human_in_loop]`

## 输入

开始前读取以下文件：

- `.d2a/docs/architecture/00_overview.md`
- `.d2a/docs/architecture/01_boundary.md`
- `.d2a/docs/architecture/02_driver.md`
- `.d2a/docs/architecture/03_core_objects.md`
- `.d2a/docs/architecture/04_state_evolution.md`
- `.d2a/docs/architecture/05_cooperation.md`
- `.d2a/docs/architecture/06_constraints.md`
- `.d2a/docs/architecture/99_code_map.md`

## 阶段 1：挑战准备

1. 确认上下文后，调用：

   `d2a skill-state d2a-challenge-architecture --status started --stage architecture-challenge-prepared --phase challenge-preparation --next-step "Prepare the six architecture decisions for challenge dialogue." --next-skill "d2a-challenge-architecture" --next-file ".d2a/docs/architecture/00_overview.md" --summary "Started architecture challenge preparation."`

2. Extract these `6` architecture decisions to challenge:
   - system boundary
   - primary driver
   - core objects
   - state evolution
   - cooperation pattern
   - dominant constraint / tradeoff
3. Do not rewrite any architecture file in this phase.
4. If needed, create a temporary challenge checklist in the conversation, but keep the architecture docs unchanged.
5. When the challenge set is ready, call:

   `d2a skill-state d2a-challenge-architecture --status progress --stage architecture-challenge-in-progress --phase challenge-dialogue --question-index 0 --question-total 6 --next-step "Start the first architecture challenge round." --next-skill "d2a-challenge-architecture" --next-file ".d2a/docs/architecture/00_overview.md" --summary "Architecture challenge set prepared; moving into challenge dialogue."`

## 阶段 2：挑战对话

1. 运行 `6` 轮挑战，每个架构决策一轮。
2. 顺序应为：
   - boundary
   - driver
   - core objects
   - state evolution
   - cooperation
   - dominant constraint
3. 每轮挑战前，保持与 d2a-step 一致的外壳格式并映射字段：
   - `[d2a]` 行必须包含 `architecture-challenge-in-progress`、`challenge-dialogue` 与 `N/6` 进度。
   - `[next]` 行应指向挑战结束后的续接目标。
4. 在第 `N` 轮挑战前，调用：

   `d2a skill-state d2a-challenge-architecture --status progress --stage architecture-challenge-in-progress --phase challenge-dialogue --question-index <N> --question-total 6 --next-step "Continue architecture challenge dialogue." --next-skill "d2a-mini-1-scope" --next-file ".d2a/docs/implementation/00_mini_scope.md" --summary "Architecture challenge round <N> is active."`

5. 每轮中：
   - 给出一个架构决策
   - 请学习者对其提出质疑
   - 接收学习者的质疑
   - 回应该质疑
   - 将质疑归类为 `strong`、`partial` 或 `weak`
   - 简要解释归类原因
6. The AI must not directly revise the architecture docs during challenge dialogue.
7. If a challenge is strong, mark it for later 复审 rather than silently editing architecture output.
8. Continue even when the learner objection is weak.

## 阶段 3：挑战收尾

1. 第 6 轮结束后：
   - 总结哪些质疑属于 strong、partial、weak
   - 列出未解决问题（如有）
   - 给出一条建议：
     - `继续推进`
     - `复审`
     - `回到架构重审`
2. 在该收尾阶段保持架构文档不变。
3. 挑战阶段结束时，调用：

   `d2a skill-state d2a-challenge-architecture --status completed --stage architecture-challenge-complete --phase challenge-dialogue --question-index 6 --question-total 6 --next-step "Proceed to d2a-mini-1-scope unless a 复审 is required." --next-skill "d2a-mini-1-scope" --next-file ".d2a/docs/implementation/00_mini_scope.md" --summary "Completed architecture challenge phase."`

## 回合结束续接规则

1. 回合结束前必须调用 `d2a skill-state` 持久化当前 phase 与 progress。
2. 回复末尾必须输出：`继续请使用 $d2a-step`。


## 输出（产物）

- 本阶段不新增 `docs/architecture/*` 产物文件。

## 持久化（.d2a）

- 挑战轮次、质疑强度评估、未解决问题、推进建议写入 `.d2a/challenge.json` 与 `.d2a/challenge_log.jsonl`。
