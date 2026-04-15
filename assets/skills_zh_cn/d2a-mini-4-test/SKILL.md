---
name: d2a-mini-4-test
description: 内置 d2a 技能，用于 d2a-mini-4-test 阶段引导与状态更新。
---

# d2a-mini-4-test

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
   - 测试设计优先采用 provider 的最小验证契约。
2. Timebox Gate（耗时）
   - 明确测试阶段预算（建议 10 分钟内完成首个端到端场景）。
   - 若预计超时，仅保留一条成功信号 + 一条失败信号。
3. Intent Gate（意图对齐）
   - 测试只验证 mini 意图实现，不扩展业务覆盖面。
   - 检查对象/状态/协作链三个锚点是否可观察。

## 阶段 1：分析生成

1. 确认上下文后，调用：

   `d2a skill-state d2a-mini-4-test --status started --stage testing-in-progress --phase analysis-generation --next-step "创建首个集成测试。" --next-skill "d2a-report-build" --next-file ".d2a/tests/README.md" --summary "Started mini-test work."`

2. 先执行三道 Gate，输出本轮 gate 结论（provider 命中情况、timebox 预算、intent 锚点）。
3. 将该技能视为 Codex 中测试阶段的用户入口。
4. 若测试任务文件尚未准备好，在写测试前调用 `d2a test-mini`。
5. 读取：
   - `.d2a/docs/implementation/03_test_plan.md`
   - `.d2a/tests/README.md`
   - `.d2a/tests/01_integration_tasks.md`
   - `.d2a/src/ARCHITECTURE.md`
6. 从首个可运行切片的一条端到端测试开始。
7. 优先验证可观察行为，而不是内部单元细节。
8. 在首个场景明确后，仅补充下一个最有价值的失败场景。
9. 测试需与要展示的架构意图保持一致。
10. 明确说明：
   - 可观察的成功信号
   - 可观察的失败信号
11. 当测试轮次稳定后，调用：

   `d2a skill-state d2a-mini-4-test --status progress --stage testing-in-progress --phase confirmation-questions --question-index 0 --question-total 4 --next-step "开始第 1 题 mini-test 确认题。" --next-skill "d2a-mini-4-test" --next-file ".d2a/tests/01_integration_tasks.md" --summary "Mini-test work complete; moving into confirmation questions."`

## 阶段 2：确认题

1. 基于阶段 1 的实际输出生成 `4` 道选择题，不要使用泛化示例。
2. 4 道题应覆盖以下角度：
   - provider 测试契约是否命中
   - timebox 下的最小测试集是否成立
   - 成功/失败信号是否可观察
   - intent 锚点是否被有效验证
3. 每轮只问一道题。
4. 每轮提问前，保持与 d2a-step 一致的外壳格式并映射字段：
   - `[d2a]` 行必须包含当前 stage、phase 与 `N/<total>` 进度。
   - `[next]` 行应指向本组问题结束后的 next skill/file。
5. 在提问第 `N` 题前，调用：

   `d2a skill-state d2a-mini-4-test --status progress --stage testing-in-progress --phase confirmation-questions --question-index <N> --question-total 4 --next-step "继续 mini-test 确认题。" --next-skill "d2a-report-build" --next-file ".d2a/report/index.md" --summary "Mini-test confirmation question <N> is active."`

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

    `d2a skill-state d2a-mini-4-test --status completed --stage testing-complete --phase confirmation-questions --question-index 4 --question-total 4 --next-step "进入 d2a-report-build。" --next-skill "d2a-report-build" --next-file ".d2a/report/index.md" --summary "Completed mini-test confirmation questions."`

11. 确认题的题干、用户答案、判定与解释必须写入 `.d2a/qa/<skill>.jsonl`，不得写入 `tests/*` 或 `report/*`。

## 回合结束续接规则

1. 回合结束前必须调用 `d2a skill-state` 持久化当前 phase 与 progress。
2. 回复末尾必须输出：`继续请使用 $d2a-step`。


## 输出（产物）

- 首个集成测试（写入 `tests/*`）。
- 下一集成场景（写入 `tests/*`）。
- 可观察成功信号与可观察失败信号（写入 `tests/01_integration_tasks.md`）。

## 持久化（.d2a）

- Provider/Timebox/Intent gate 决策写入 `.d2a/mini_gate/d2a-mini-4-test.json`。
- 确认题题干、用户答案、判定、解释写入 `.d2a/qa/<skill>.jsonl`。
- 简短理解回顾与`理解度打分`写入 `.d2a/qa/<skill>.jsonl`。
