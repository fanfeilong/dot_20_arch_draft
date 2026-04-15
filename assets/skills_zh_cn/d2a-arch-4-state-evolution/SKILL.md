---
name: d2a-arch-4-state-evolution
description: 内置 d2a 技能，用于 d2a-arch-4-state-evolution 阶段引导与状态更新。
---

# d2a-arch-4-state-evolution

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

## 阶段 1：原子问题对齐（一次补充机会）

1. 确认上下文后，先调用：

   `d2a skill-state d2a-arch-4-state-evolution --status started --stage architecture-in-progress --phase atomic-question-alignment --question-index 0 --question-total 1 --next-step "展示原子分析问题，并进行一次性补充询问（是/否）。" --next-skill "d2a-arch-4-state-evolution" --next-file "docs/1.架构拆解/04_状态演化.md" --summary "已启动 state-evolution 原子问题对齐。"`

2. 在开始分析前，必须先向用户输出：

   `接下来我会扫描并回答这些问题：<列出本技能的原子问题>；请问这些问题有需要补充么？（是/否）`

3. 本阶段只允许一次补充交互：
   - 用户答 `是`：收集补充问题并与原子问题合并，然后回显“已合并问题清单”。
   - 用户答 `否`：直接使用原子问题进入分析。
4. 不允许在用户确认前写入 `docs/1.架构拆解/04_状态演化.md`。
5. 完成对齐后，调用：

   `d2a skill-state d2a-arch-4-state-evolution --status progress --stage architecture-in-progress --phase analysis-generation --question-index 1 --question-total 1 --next-step "使用合并后的原子问题进入分析生成。" --next-skill "d2a-arch-4-state-evolution" --next-file "docs/1.架构拆解/04_状态演化.md" --summary "state-evolution 原子问题对齐完成。"`

## 阶段 2：分析生成

1. 确认上下文后，调用：

   `d2a skill-state d2a-arch-4-state-evolution --status progress --stage architecture-in-progress --phase analysis-generation --next-step "Track the main 状态阶段 and transitions." --next-skill "d2a-arch-5-module-view" --next-file "docs/1.架构拆解/04_状态演化.md" --summary "Started state-evolution analysis."`

2. 基于真实仓库进行分析，并将结果写入 `docs/1.架构拆解/04_状态演化.md`.
3. 回答以下合并后的原子问题（基础问题 + 用户可选补充）：
   - 应该跟踪的单一对象或流程是什么？
   - What are its three to six 状态阶段?
   - 主要状态迁移由什么触发？
   - 状态在哪里存储、观测或重建？
4. 最终写入 `docs/1.架构拆解/04_状态演化.md` 的内容必须仅包含一个 Mermaid `stateDiagram-v2` 代码块，不要附加段落、列表或解释文字。
5. 初稿完成后，强制执行三轮修订：
   - 压缩修订
   - 去术语修订
   - 口语化简化修订
6. 优先关注状态推进，而不是实现细节。
7. 当分析草稿稳定后，调用：

   `d2a skill-state d2a-arch-4-state-evolution --status progress --stage architecture-in-progress --phase confirmation-questions --question-index 0 --question-total 4 --next-step "开始第 1 题 state-evolution 确认题。" --next-skill "d2a-arch-4-state-evolution" --next-file "docs/1.架构拆解/04_状态演化.md" --summary "State-evolution analysis complete; moving into confirmation questions."`

## 阶段 3：确认题

1. 基于阶段 1 的实际输出生成 `4` 道选择题，不要使用泛化示例。
   - 干扰项必须与真实实现或常见误解高度相似，不能用明显荒谬选项凑数。
   - 正确项与干扰项在字面上应有迷惑性，判定应依赖对本项目结论的理解，而非关键词匹配。
   - 每题至少 `2` 个干扰项必须使用本项目里真实出现过的概念、模块或流程名，但语义需错位。
2. 4 道题应覆盖以下角度：
   - 被跟踪对象或流程
   - 状态阶段
   - 迁移触发器
   - 状态存储/观测点
3. 每轮只问一道题。
4. 每轮提问前，保持与 d2a-step 一致的外壳格式并映射字段：
   - `[d2a]` 行必须包含当前 stage、phase 与 `N/<total>` 进度。
   - `[next]` 行应指向本组问题结束后的 next skill/file。
5. 在提问第 `N` 题前，调用：

   `d2a skill-state d2a-arch-4-state-evolution --status progress --stage architecture-in-progress --phase confirmation-questions --question-index <N> --question-total 4 --next-step "继续 state-evolution 确认题。" --next-skill "d2a-arch-5-module-view" --next-file "docs/1.架构拆解/05_协作.md" --summary "State-evolution confirmation question <N> is active."`

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

    `d2a skill-state d2a-arch-4-state-evolution --status completed --stage architecture-in-progress --phase confirmation-questions --question-index 4 --question-total 4 --next-step "进入 d2a-arch-5-module-view。" --next-skill "d2a-arch-5-module-view" --next-file "docs/1.架构拆解/05_协作.md" --summary "Completed state-evolution confirmation questions."`

13. 确认题的题干、用户答案、判定与解释必须写入 `.d2a/qa/<skill>.jsonl`，不得写入 `docs/architecture/*.md`。

## 回合结束续接规则

1. 回合结束前必须调用 `d2a skill-state` 持久化当前 phase 与 progress。
2. 回复末尾必须输出：`继续请使用 $d2a-step`。


## 输出（产物）

- 一个 Mermaid `stateDiagram-v2` 状态机图（且仅此一个代码块），写入 `docs/architecture/04_state_evolution.md`。

## 持久化（.d2a）

- 对齐后的原子问题集（基础 + 用户补充）写入 `.d2a/qa/<skill>.json`。
- 确认题题干、用户答案、判定、解释写入 `.d2a/qa/<skill>.jsonl`。
- 简短理解回顾与`理解度打分`写入 `.d2a/qa/<skill>.jsonl`。
