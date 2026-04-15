<script setup>
import { onMounted, ref } from "vue";

const status = ref("Loading report data ...");
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
    status.value = "Loaded data from /data/*.json";
  } catch (err) {
    status.value = `Failed to load data: ${String(err)}`;
  }
});
</script>

<template>
  <main class="wrap">
    <section class="card">
      <p class="eyebrow">d2a report app</p>
      <h1>Vue Report App</h1>
      <p>{{ status }}</p>

      <h2>Target</h2>
      <ul>
        <li><code>{{ target.target_repo || "unknown" }}</code></li>
        <li><code>{{ target.repo_root || "unknown" }}</code></li>
        <li><code>{{ target.d2a_path || "unknown" }}</code></li>
      </ul>

      <h2>Architecture Docs</h2>
      <ul>
        <li v-for="item in summary.architecture_docs || []" :key="item">
          <code>{{ item }}</code>
        </li>
      </ul>

      <h2>Test Outputs</h2>
      <ul>
        <li v-for="item in tests.outputs || []" :key="item">
          <code>{{ item }}</code>
        </li>
      </ul>

      <h2>Challenge</h2>
      <ul>
        <li><code>{{ challenge.current_stage || "unknown" }}</code></li>
        <li><code>{{ challenge.recommendation || "unknown" }}</code></li>
      </ul>
    </section>
  </main>
</template>
