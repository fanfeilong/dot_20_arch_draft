---
name: d2a-mini-1-scope
description: 内置 d2a 技能，用于 d2a-mini-1-scope 阶段引导与状态更新。
---

# d2a-mini-1-scope

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

如果无法确定当前仓库，立即停止并询问用户要使用哪个仓库。

## Mini 快车道三道 Gate

在 1 小时演讲场景中，mini 阶段必须先通过以下三道 Gate，再进入实现细节：

1. Provider Gate（技术栈）
   - 先识别目标仓库技术栈并匹配内置 provider。
   - 若命中 provider，优先采用 provider 的最小实现方案。
   - 若未命中，退回通用最小方案（单主路径、低依赖、可运行）。
2. Timebox Gate（耗时）
   - 明确 mini 构建时间预算（建议 20 分钟）。
   - 若预计超时，必须降级范围，禁止扩展功能面。
3. Intent Gate（意图对齐）
   - mini 只证明架构意图，不重做业务分析。
   - 必须对齐 arch 阶段的意图锚点：核心对象、状态流转、协作链。

## 阶段 1：分析生成

1. 确认上下文后，调用：

   `d2a skill-state d2a-mini-1-scope --status started --stage mini-derivation-prepared --phase analysis-generation --next-step "选择要保留的单一架构意图。" --next-skill "d2a-mini-2-design" --next-file ".d2a/docs/implementation/00_mini_scope.md" --summary "Started mini-scope derivation."`

2. 将该技能视为 Codex 中 mini 实现阶段的用户入口。
3. 先执行三道 Gate，输出本轮 gate 结论（provider 命中情况、timebox 预算、intent 锚点）。
4. 若实现规划文件尚未准备好，在产出内容前调用 `d2a derive-mini`。
5. 读取：
   - `.d2a/docs/architecture/00_overview.md`
   - `.d2a/docs/architecture/02_driver.md`
   - `.d2a/docs/architecture/03_core_objects.md`
   - `.d2a/docs/architecture/04_state_evolution.md`
   - `.d2a/docs/architecture/05_cooperation.md`
   - `.d2a/docs/architecture/06_constraints.md`
6. 将结果写入 `.d2a/docs/implementation/00_mini_scope.md`.
7. 回答以下原子问题：
   - 必须保留的单一架构意图是什么？
   - Which 可运行的 20% 切片 is enough to demonstrate it?
   - mini 版本将有意省略什么？
   - 哪种技术栈应与原项目保持一致？
8. 初稿完成后，强制执行三轮修订：
   - 压缩修订
   - 去术语修订
   - 口语化简化修订
9. 范围需足够小，以支持首个可运行切片。
10. 当分析草稿稳定后，调用：

   `d2a skill-state d2a-mini-1-scope --status progress --stage mini-derivation-prepared --phase confirmation-questions --question-index 0 --question-total 4 --next-step "开始第 1 题 mini-scope 确认题。" --next-skill "d2a-mini-1-scope" --next-file ".d2a/docs/implementation/00_mini_scope.md" --summary "Mini-scope analysis complete; moving into confirmation questions."`

## 阶段 2：确认题

1. 基于阶段 1 的实际输出生成 `4` 道选择题，不要使用泛化示例。
2. 4 道题应覆盖以下角度：
   - provider 选择与技术栈合理性
   - timebox 下的最小可运行切片
   - intent 锚点对齐（对象/状态/协作链）
   - 有意省略项与范围控制
3. 每轮只问一道题。
4. 每轮提问前，保持与 d2a-step 一致的外壳格式并映射字段：
   - `[d2a]` 行必须包含当前 stage、phase 与 `N/<total>` 进度。
   - `[next]` 行应指向本组问题结束后的 next skill/file。
5. 在提问第 `N` 题前，调用：

   `d2a skill-state d2a-mini-1-scope --status progress --stage mini-derivation-prepared --phase confirmation-questions --question-index <N> --question-total 4 --next-step "继续 mini-scope 确认题。" --next-skill "d2a-mini-2-design" --next-file ".d2a/docs/implementation/01_mini_design.md" --summary "Mini-scope confirmation question <N> is active."`

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

    `d2a skill-state d2a-mini-1-scope --status completed --stage mini-derivation-prepared --phase confirmation-questions --question-index 4 --question-total 4 --next-step "进入 d2a-mini-2-design。" --next-skill "d2a-mini-2-design" --next-file ".d2a/docs/implementation/01_mini_design.md" --summary "Completed mini-scope confirmation questions."`

11. 确认题的题干、用户答案、判定与解释必须写入 `.d2a/qa/<skill>.jsonl`，不得写入 `docs/implementation/*.md` 或 `src/*`。

## 回合结束续接规则

1. 回合结束前必须调用 `d2a skill-state` 持久化当前 phase 与 progress。
2. 回复末尾必须输出：`继续请使用 $d2a-step`。


## 输出（产物）

- 保留的架构意图。
- 可运行的 20% 切片定义。
- 有意省略项。
- 目标技术栈。
- 以上内容写入 `docs/implementation/00_mini_scope.md`。

## 持久化（.d2a）

- Provider/Timebox/Intent gate 决策写入 `.d2a/mini_gate/d2a-mini-1-scope.json`。
- 确认题题干、用户答案、判定、解释写入 `.d2a/qa/<skill>.jsonl`。
- 简短理解回顾与`理解度打分`写入 `.d2a/qa/<skill>.jsonl`。
