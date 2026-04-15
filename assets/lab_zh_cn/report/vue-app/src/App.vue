<script setup>
import { onMounted, ref } from "vue";

const status = ref("正在加载报告数据...");
const target = ref({});
const summary = ref({});
const tests = ref({});
const challenge = ref({});

onMounted(async () => {
  try {
    const [summaryRes, targetRes, testsRes, challengeRes] = await Promise.all([
      fetch("/data/summary.json"),
      fetch("/data/target.json"),
      fetch("/data/tests.json"),
      fetch("/data/challenge.json"),
    ]);

    summary.value = await summaryRes.json();
    target.value = await targetRes.json();
    tests.value = await testsRes.json();
    challenge.value = await challengeRes.json();
    status.value = "已从 /data/*.json 加载数据";
  } catch (err) {
    status.value = `加载数据失败：${String(err)}`;
  }
});
</script>

<template>
  <main class="wrap">
    <section class="card">
      <p class="eyebrow">d2a 报告应用</p>
      <h1>Vue 报告应用</h1>
      <p>{{ status }}</p>

      <h2>目标</h2>
      <ul>
        <li><code>{{ target.target_repo || "unknown" }}</code></li>
        <li><code>{{ target.repo_root || "unknown" }}</code></li>
        <li><code>{{ target.d2a_path || "unknown" }}</code></li>
      </ul>

      <h2>架构文档</h2>
      <ul>
        <li v-for="item in summary.architecture_docs || []" :key="item">
          <code>{{ item }}</code>
        </li>
      </ul>

      <h2>测试输出</h2>
      <ul>
        <li v-for="item in tests.outputs || []" :key="item">
          <code>{{ item }}</code>
        </li>
      </ul>

      <h2>挑战结论</h2>
      <ul>
        <li><code>{{ challenge.current_stage || "unknown" }}</code></li>
        <li><code>{{ challenge.recommendation || "unknown" }}</code></li>
      </ul>
    </section>
  </main>
</template>
