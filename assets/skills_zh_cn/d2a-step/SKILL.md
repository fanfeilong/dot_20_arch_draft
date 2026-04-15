---
name: d2a-step
description: State-driven orchestrator skill that resumes from .d2a state and routes to the correct next d2a sub-skill.
---

# d2a-step

## 目标

该技能用于当前阶段的中文引导、状态推进与教学问答。

## 语言规则

所有面向用户的文本都必须使用简体中文。

## 必需输出格式

每次 `d2a-step` 回复都必须使用以下紧凑格式：

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
3. 正文保持简洁可执行；除非用户明确要求调试，不要输出冗长原始状态。
4. 正文至少两行，不得只输出一行：
   - 第 1 行：当前动作（正在做什么）
   - 第 2 行：本轮目的/依据（为什么这样做或产出落点）
5. 正文必须自动分行，单行不超过 100 字；超长内容要拆成多行短句。
6. 正文禁止使用 Markdown 强调符号（如 `` `...` ``、`**...**`），避免在 Cursor 中呈现为代码感整块文本。
7. 若 `next_file` 未知，输出 `unknown`，且保持单行格式不变。


## 正文排版硬规则

1. 正文必须使用 `- ` 列表输出，最少 2 条、最多 4 条。
2. 每条要点独占一行，不得写成长段落。
3. 单行不超过 100 个中文字符；超长必须拆行。
4. 正文禁止使用 Markdown 强调符号（如 `` `...` ``、`**...**`）。

如果无法确定当前仓库，立即停止并询问用户要使用哪个仓库。

## Human In Loop 标记规则

当本回合包含“向用户提问并等待回答”的动作时，回复正文最后一行必须追加：

`[human_in_loop]`

### 真实输出示例（用于验收骨架可读性）

```text
==================================================
【分析 3/6｜核心对象】n8n_d2a
next: d2a-arch-4-state-evolution -> docs/architecture/04_state_evolution.md
==================================================

- 已恢复到确认题阶段，将继续第 3 题。
- 本轮仅推进题目与判定，不改写架构文档。

--------------------------------------------------
done: 完成核心对象确认题第 3 题判定
state: 分析 3/6｜核心对象 -> 分析 4/6｜状态机 · 继续请使用 $d2a-step
--------------------------------------------------
```

```text
==================================================
【分析 1/6｜边界】n8n_d2a
next: d2a-arch-1-project-scope -> .d2a/docs/architecture/01_boundary.md
==================================================

- 即将执行：d2a-arch-1-project-scope（文件：.d2a/docs/architecture/01_boundary.md）。
- 请确认是否继续本动作。（是/否）
[human_in_loop]

--------------------------------------------------
done: 已完成下一个动作的执行前确认
state: 分析 1/6｜边界 -> 分析 1/6｜边界 · 继续请使用 $d2a-step
--------------------------------------------------
```

```text
==================================================
【实现 2/4｜最小设计】n8n_d2a
next: d2a-mini-3-build -> docs/implementation/02_build_plan.md
==================================================

- 已读取最小范围结论，进入最小设计草案收敛。
- 本轮目标是确定模块边界与主流程。

--------------------------------------------------
done: 完成最小设计草案并进入确认题
state: 实现 2/4｜最小设计 -> 实现 3/4｜最小构建 · 继续请使用 $d2a-step
--------------------------------------------------
```

```text
==================================================
【报告 1/1｜报告构建】n8n_d2a
next: d2a-report-build -> report/index.md
==================================================

- 已汇总分析、实现、测试产物并刷新报告数据文件。
- 可继续优化报告叙事与展示结构。

--------------------------------------------------
done: 完成报告构建主流程
state: 报告 1/1｜报告构建 -> 报告 1/1｜报告构建 · 继续请使用 $d2a-step
--------------------------------------------------
```

## 状态恢复

1. 调用 `d2a status` 或读取 `.d2a/state.json` 获取：
   - `current_stage`
   - `current_phase`
   - `question_index`
   - `question_total`
   - `current_skill`
   - `next_skill`
   - `next_file`
2. 当续接细节不清晰时，读取最近 `.d2a/history.jsonl`。
3. 若 `.d2a/state.json` 缺失，提示先执行 `d2a init` 并停止。

## 路由规则

### 二层骨架映射（中文显示）

1. 第一层固定为：
   - `分析`
   - `实现`
   - `报告`
2. 第二层固定为：
   - 分析：`边界` / `驱动` / `核心对象` / `状态机` / `模块协作` / `约束`
   - 实现：`最小范围` / `最小设计` / `最小构建` / `最小测试`
   - 报告：`报告构建`
3. 头部与尾部中的骨架位置必须使用这组中文名，不得回退为英文步骤名。

4. 若 `current_phase` 为 `confirmation-questions` 或 `challenge-dialogue` 且 `question_index < question_total`，则续接 `current_skill`。
5. 否则按阶段路由：
   - `initialized` or `analysis-prepared` or `architecture-in-progress` -> `d2a-arch-1-project-scope`
   - `architecture-complete` or `architecture-challenge-prepared` or `architecture-challenge-in-progress` -> `d2a-challenge-architecture`
   - `architecture-challenge-complete` or `mini-derivation-prepared` -> `d2a-mini-1-scope`
   - `mini-design-in-progress` -> `d2a-mini-2-design`
   - `mini-design-complete` -> `d2a-mini-3-build`
   - `test-plan-prepared` or `testing-in-progress` -> `d2a-mini-4-test`
   - `testing-complete` or `report-prepared` or `report-ready` -> `d2a-report-build`
   - `serving` -> `d2a-status`
6. 若 `next_skill` 存在且不与阶段路由冲突，优先使用 `next_skill`。
7. 若路由存在歧义，提出一个简短澄清问题，不要猜测。

## 推进执行

1. 在给出路由结果前，先调用：

   `d2a skill-state d2a-step --status started --phase analysis-generation --next-step "Resume from persisted d2a state and route to the next skill." --summary "Started d2a-step orchestration."`

2. 进入下一个动作前必须先输出“执行意图确认”并等待用户确认，不得直接推进执行：
   - 固定提示：`即将执行：<ROUTED_SKILL>（文件：<ROUTED_FILE>）。是否继续？（是/否）`
   - 该回合正文最后一行必须追加：`[human_in_loop]`
   - 用户答 `否`：停止推进，询问用户希望改为哪个技能或阶段。
   - 用户答 `是`：继续第 3 步。
3. 明确告诉用户当前应执行哪个技能，并说明理由（阶段 + 相位证据）。
4. 如有必要，给出应继续的精确文件（`next_file`）。
5. 在用户确认继续后，再调用：

   `d2a skill-state d2a-step --status completed --phase analysis-generation --next-step "Continue with the routed skill." --next-skill "<ROUTED_SKILL>" --next-file "<ROUTED_FILE>" --summary "d2a-step routed to <ROUTED_SKILL> based on persisted state."`

6. 结尾必须输出：`继续请使用 $d2a-step`。

## 输出

- 当前骨架位置（第一层 + 第二层 + 阶段内进度）
- 已选择下一技能
- 路由原因（简短）
- 下一续接文件
- `d2a-step` 续接提示
