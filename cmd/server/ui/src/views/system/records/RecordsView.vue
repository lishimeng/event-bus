<script setup lang="ts">
import { ref, onMounted } from "vue";
import { fetchRecords } from "@/api/admin/records.api";
import type { DataRecordRow } from "@/types/models/admin";
import { formatRecordTime, prettyJson } from "@/utils/format";

const recordsItems = ref<DataRecordRow[]>([]);
const recordsTotal = ref(0);
const recordsLimit = ref(20);
const recordsOffset = ref(0);
const recordsSourceFilter = ref("");
const recordsLoading = ref(false);
const recordsErr = ref("");
const expandedRecordId = ref<number | null>(null);

async function loadRecords() {
  recordsErr.value = "";
  recordsLoading.value = true;
  try {
    const { total, items, error } = await fetchRecords({
      limit: recordsLimit.value,
      offset: recordsOffset.value,
      source: recordsSourceFilter.value,
    });
    if (error) {
      recordsErr.value = error;
      return;
    }
    recordsTotal.value = total;
    recordsItems.value = items;
  } finally {
    recordsLoading.value = false;
  }
}

function recordsNextPage() {
  if (recordsOffset.value + recordsLimit.value < recordsTotal.value) {
    recordsOffset.value += recordsLimit.value;
    void loadRecords();
  }
}

function recordsPrevPage() {
  recordsOffset.value = Math.max(0, recordsOffset.value - recordsLimit.value);
  void loadRecords();
}

function recordsFirstPage() {
  recordsOffset.value = 0;
  void loadRecords();
}

function toggleRecordExpand(id: number) {
  expandedRecordId.value = expandedRecordId.value === id ? null : id;
}

onMounted(() => {
  void loadRecords();
});
</script>

<template>
  <section class="card">
    <h2>消息记录</h2>
    <p class="hint">
      数据表 <code>data_record</code>，由收发链路异步写入。
      <code>GET /api/v1/admin/records</code> 支持 <code>limit</code>、<code>offset</code>、<code>source</code>（按来源/Topic 精确过滤）。请求请<strong>不要</strong>多写尾斜杠，避免 301 返回 HTML 导致解析失败。
    </p>
    <div class="toolbar toolbar-wrap">
      <label class="inline-field">
        <span>source 过滤</span>
        <input v-model="recordsSourceFilter" type="text" placeholder="留空表示全部" autocomplete="off" />
      </label>
      <label class="inline-field">
        <span>每页</span>
        <select v-model.number="recordsLimit" class="narrow-select" @change="recordsFirstPage">
          <option :value="10">10</option>
          <option :value="20">20</option>
          <option :value="50">50</option>
        </select>
      </label>
      <button type="button" class="secondary" :disabled="recordsLoading" @click="recordsFirstPage">
        {{ recordsLoading ? "加载中…" : "查询 / 刷新" }}
      </button>
    </div>
    <p v-if="recordsErr" class="warn">{{ recordsErr }}</p>
    <p v-if="!recordsErr && recordsTotal > 0" class="hint pager-hint">
      共 <strong>{{ recordsTotal }}</strong> 条，当前第
      {{ recordsOffset + 1 }}–{{ Math.min(recordsOffset + recordsLimit, recordsTotal) }} 条
    </p>
    <div v-if="recordsItems.length" class="table-wrap">
      <table class="data-table records-table">
        <thead>
          <tr>
            <th class="col-id">ID</th>
            <th>Code</th>
            <th>Refer</th>
            <th>Source</th>
            <th>方向</th>
            <th class="col-time">时间</th>
            <th class="col-act">详情</th>
          </tr>
        </thead>
        <tbody>
          <template v-for="row in recordsItems" :key="row.id">
            <tr class="rec-row" @click="toggleRecordExpand(row.id)">
              <td class="mono">{{ row.id }}</td>
              <td class="mono clip">{{ row.code }}</td>
              <td class="mono clip">{{ row.referCode || "—" }}</td>
              <td class="mono clip">{{ row.source }}</td>
              <td>{{ row.routeLabel }}</td>
              <td class="small">{{ formatRecordTime(row.createTime || row.updateTime) }}</td>
              <td>
                <button type="button" class="linkish" @click.stop="toggleRecordExpand(row.id)">
                  {{ expandedRecordId === row.id ? "收起" : "展开" }}
                </button>
              </td>
            </tr>
            <tr v-if="expandedRecordId === row.id" class="detail-row">
              <td colspan="7">
                <div class="detail-grid">
                  <div>
                    <h4 class="detail-title">Payload</h4>
                    <pre class="detail-pre mono">{{ prettyJson(row.payload) }}</pre>
                  </div>
                  <div>
                    <h4 class="detail-title">BizPayload</h4>
                    <pre class="detail-pre mono">{{ prettyJson(row.bizPayload) }}</pre>
                  </div>
                </div>
              </td>
            </tr>
          </template>
        </tbody>
      </table>
    </div>
    <p v-else-if="!recordsLoading && !recordsErr" class="hint">暂无记录，可先发布或消费消息后再刷新。</p>
    <div v-if="recordsTotal > recordsLimit" class="pager">
      <button
        type="button"
        class="secondary"
        :disabled="recordsOffset === 0 || recordsLoading"
        @click="recordsPrevPage"
      >
        上一页
      </button>
      <button
        type="button"
        class="secondary"
        :disabled="recordsOffset + recordsLimit >= recordsTotal || recordsLoading"
        @click="recordsNextPage"
      >
        下一页
      </button>
    </div>
  </section>
</template>
