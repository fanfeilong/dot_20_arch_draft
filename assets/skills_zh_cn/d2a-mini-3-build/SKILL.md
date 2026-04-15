---
name: d2a-mini-3-build
description: 内置 d2a 技能，用于 d2a-mini-3-build 阶段引导与状态更新。
---

# d2a-mini-3-build

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

## Mini 快车道三道 Gate

在 1 小时演讲场景中，mini 阶段必须先通过以下三道 Gate，再进入实现细节：

1. Provider Gate（技术栈）
   - 先识别目标仓库技术栈并匹配内置 provider。
   - 若命中 provider，直接采用 provider 的最小实现模板与文件布局。
2. Timebox Gate（耗时）
   - 明确 mini 构建时间预算（建议 20 分钟）。
   - 若预计超时，必须立即降级为单主路径实现。
3. Intent Gate（意图对齐）
   - 实现目标仅限证明意图，不做业务扩展。
   - 仅围绕 arch 锚点（对象/状态/协作链）落代码。

## 阶段 1：分析生成

1. 确认上下文后，调用：

   `d2a skill-state d2a-mini-3-build --status started --stage mini-design-complete --phase analysis-generation --next-step "实现首个可运行 mini 切片。" --next-skill "d2a-mini-4-test" --next-file ".d2a/src/ARCHITECTURE.md" --summary "Started mini-build work."`

2. 先执行三道 Gate，输出本轮 gate 结论（provider 命中情况、timebox 预算、intent 锚点）。
3. 若实现规划文件尚未准备好，调用 `d2a derive-mini`。
4. 读取：
   - `.d2a/docs/implementation/00_mini_scope.md`
   - `.d2a/docs/implementation/01_mini_design.md`
   - `.d2a/docs/implementation/02_build_plan.md`
   - `.d2a/src/ARCHITECTURE.md`
5. Implement only the first runnable slice described in the build plan.
6. Prefer a small but executable result over broad coverage.
7. Keep the 选定技术栈 aligned with the original project when practical.
8. Update `.d2a/src/ARCHITECTURE.md` if implementation reality forces a design correction.
9. Do not expand scope until the first runnable slice is working.
10. After the implementation is stable, write a brief note on what remains intentionally unimplemented.
11. When the implementation pass is stable, call:

   `d2a skill-state d2a-mini-3-build --status progress --stage mini-design-complete --phase confirmation-questions --question-index 0 --question-total 4 --next-step "开始第 1 题 mini-build 确认题。" --next-skill "d2a-mini-3-build" --next-file ".d2a/src/ARCHITECTURE.md" --summary "Mini-build implementation complete; moving into confirmation questions."`

## 阶段 2：确认题

1. 基于阶段 1 的实际输出生成 `4` 道选择题，不要使用泛化示例。
2. 4 道题应覆盖以下角度：
   - provider 模板与技术栈执行情况
   - timebox 下的单主路径实现取舍
   - intent 锚点（对象/状态/协作链）是否被证明
   - 有意未实现项与范围控制
3. 每轮只问一道题。
4. 每轮提问前，保持与 d2a-step 一致的外壳格式并映射字段：
   - `[d2a]` 行必须包含当前 stage、phase 与 `N/<total>` 进度。
   - `[next]` 行应指向本组问题结束后的 next skill/file。
5. 在提问第 `N` 题前，调用：

   `d2a skill-state d2a-mini-3-build --status progress --stage mini-design-complete --phase confirmation-questions --question-index <N> --question-total 4 --next-step "继续 mini-build 确认题。" --next-skill "d2a-mini-4-test" --next-file ".d2a/tests/README.md" --summary "Mini-build confirmation question <N> is active."`

6. 给出一道选择题。
7. 等待学习者作答。
8. 学习者作答后：
   - 判断答案是正确、部分正确还是错误
   - 给出一句简短解释
   - 即使答错也继续下一题
9. 第 4 题评估后：
   - 输出简短回顾
   - 输出 `理解度打分`
   - `理解度打分` 控制在 100 字以内
10. 确认题阶段结束时，调用：

    `d2a skill-state d2a-mini-3-build --status completed --stage mini-design-complete --phase confirmation-questions --question-index 4 --question-total 4 --next-step "进入 d2a-mini-4-test。" --next-skill "d2a-mini-4-test" --next-file ".d2a/tests/README.md" --summary "Completed mini-build confirmation questions."`

11. 确认题的题干、用户答案、判定与解释必须写入 `.d2a/qa/<skill>.jsonl`，不得写入 `docs/implementation/*.md` 或 `src/*`。

## 回合结束续接规则

1. 回合结束前必须调用 `d2a skill-state` 持久化当前 phase 与 progress。
2. 回复末尾必须输出：`继续请使用 $d2a-step`。


## 输出（产物）

- 可运行的 mini 实现（写入 `src/*`）。
- 更新后的 `src/ARCHITECTURE.md`。
- 关于仍刻意未实现内容的简要说明（写入 `docs/implementation/02_build_plan.md` 或 `src/ARCHITECTURE.md` 的实现备注段）。

## 持久化（.d2a）

- Provider/Timebox/Intent gate 决策写入 `.d2a/mini_gate/d2a-mini-3-build.json`。
- 确认题题干、用户答案、判定、解释写入 `.d2a/qa/<skill>.jsonl`。
- 简短理解回顾与`理解度打分`写入 `.d2a/qa/<skill>.jsonl`。
