---
name: d2a-mini-2-design
description: 内置 d2a 技能，用于 d2a-mini-2-design 阶段引导与状态更新。
---

# d2a-mini-2-design

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
3. 单行不超过 80 个中文字符；超长必须拆行。
4. 正文禁止使用 Markdown 强调符号（如 `` `...` ``、`**...**`）。

5. 对条目、编号列表、选择题选项（A/B/C/D）等结构化内容，不要为满足字数限制而合并或破坏结构；保持一项一行。

如果无法确定当前仓库，立即停止并询问用户要使用哪个仓库。

## Human In Loop 标记规则

当本回合包含“向用户提问并等待回答”的动作时，回复正文最后一行必须追加：

`[human_in_loop]`

## Mini 快车道三道 Gate

在 1 小时演讲场景中，mini 阶段必须先通过以下三道 Gate，再进入实现细节：

1. Provider Gate（技术栈）
   - 先识别目标仓库技术栈并匹配内置 provider。
   - 若命中 provider，优先采用 provider 的最小实现方案。
2. Timebox Gate（耗时）
   - 明确 mini 构建时间预算（建议 20 分钟）。
   - 若预计超时，必须缩减设计复杂度与模块数量。
3. Intent Gate（意图对齐）
   - 设计必须直接服务 arch 阶段意图锚点，不做新业务发散。

## 阶段 1：分析生成

1. 确认上下文后，调用：

   `d2a skill-state d2a-mini-2-design --status started --stage mini-design-in-progress --phase analysis-generation --next-step "设计最小可用的可运行 mini 架构。" --next-skill "d2a-mini-3-build" --next-file "docs/2.mini实现/01_最小设计.md" --summary "Started mini-design work."`

2. 先执行三道 Gate，输出本轮 gate 结论（provider 命中情况、timebox 预算、intent 锚点）。
3. 若实现规划文件尚未准备好，调用 `d2a derive-mini`。
4. 读取：
   - `docs/2.mini实现/00_最小范围.md`
   - 其引用的架构文档
5. 将结果写入 `docs/2.mini实现/01_最小设计.md`.
6. 若 `.d2a/src/ARCHITECTURE.md` 摘要已过期，请与选定设计保持一致并更新。
7. 回答以下原子问题：
   - What are the 主要模块 of the mini version?
   - What 接口或入口点 are required?
   - What is the 运行流程 of the mini version?
   - 必须保留的状态模型是什么？
8. 初稿完成后，强制执行三轮修订：
   - 压缩修订
   - 去术语修订
   - 口语化简化修订
9. 设计规模需足够小，以支持首个可运行切片。
10. 设计需绑定核心架构意图，而不是追求宽泛功能覆盖。
11. 当分析草稿稳定后，调用：

   `d2a skill-state d2a-mini-2-design --status progress --stage mini-design-in-progress --phase confirmation-questions --question-index 0 --question-total 4 --next-step "开始第 1 题 mini-design 确认题。" --next-skill "d2a-mini-2-design" --next-file "docs/2.mini实现/01_最小设计.md" --summary "Mini-design analysis complete; moving into confirmation questions."`

## 阶段 2：确认题

1. 基于阶段 1 的实际输出生成 `4` 道选择题，不要使用泛化示例。
   - 干扰项必须与真实实现或常见误解高度相似，不能用明显荒谬选项凑数。
   - 正确项与干扰项在字面上应有迷惑性，判定应依赖对本项目结论的理解，而非关键词匹配。
   - 每题至少 `2` 个干扰项必须使用本项目里真实出现过的概念、模块或流程名，但语义需错位。
2. 4 道题应覆盖以下角度：
   - provider 约束下的设计合理性
   - timebox 下的模块/接口最小化
   - 运行流程是否服务核心意图
   - 状态模型与 arch 锚点一致性
3. 每轮只问一道题。
4. 每轮提问前，保持与 d2a-step 一致的外壳格式并映射字段：
   - `[d2a]` 行必须包含当前 stage、phase 与 `N/<total>` 进度。
   - `[next]` 行应指向本组问题结束后的 next skill/file。
5. 在提问第 `N` 题前，调用：

   `d2a skill-state d2a-mini-2-design --status progress --stage mini-design-in-progress --phase confirmation-questions --question-index <N> --question-total 4 --next-step "继续 mini-design 确认题。" --next-skill "d2a-mini-3-build" --next-file ".d2a/src/ARCHITECTURE.md" --summary "Mini-design confirmation question <N> is active."`

6. 给出一道选择题。
7. 题目展示必须保持结构：
   - 题干独占一行
   - 选项 `A.`、`B.`、`C.`、`D.` 各自独占一行
   - `[human_in_loop]` 必须独占一行
8. 禁止将题干与选项合并到同一行。
9. 等待学习者作答。
10. 学习者作答后：
   - 判断答案是正确、部分正确还是错误
   - 给出一句简短解释
   - 即使答错也继续下一题
11. 第 4 题评估后：
   - 输出简短回顾
   - 输出 `理解度打分`
   - `理解度打分` 控制在 80 字以内
12. 确认题阶段结束时，调用：

    `d2a skill-state d2a-mini-2-design --status completed --stage mini-design-complete --phase confirmation-questions --question-index 4 --question-total 4 --next-step "进入 d2a-mini-3-build。" --next-skill "d2a-mini-3-build" --next-file ".d2a/src/ARCHITECTURE.md" --summary "Completed mini-design confirmation questions."`

13. 确认题的题干、用户答案、判定与解释必须写入 `.d2a/qa/<skill>.jsonl`，不得写入 `docs/implementation/*.md` 或 `src/*`。

## 回合结束续接规则

1. 回合结束前必须调用 `d2a skill-state` 持久化当前 phase 与 progress。
2. 回复末尾必须输出：`继续请使用 $d2a-step`。


## 输出（产物）

- 主要模块。
- 主要接口或入口点。
- 运行流程。
- 状态模型。
- 以上内容写入 `docs/implementation/01_mini_design.md`。

## 持久化（.d2a）

- Provider/Timebox/Intent gate 决策写入 `.d2a/mini_gate/d2a-mini-2-design.json`。
- 确认题题干、用户答案、判定、解释写入 `.d2a/qa/<skill>.jsonl`。
- 简短理解回顾与`理解度打分`写入 `.d2a/qa/<skill>.jsonl`。
