<script setup lang="ts">
import { ref } from "vue";
import { postPublish } from "@/api/communication/publish.api";

const publishRoute = ref("demo-topic");
const publishMethod = ref("POST");
const publishAction = ref("/callback/handler");
const publishParamsJson = ref('{\n  "hello": "event-bus"\n}');
const publishReferId = ref("");
const publishResult = ref("");

async function doPublish() {
  publishResult.value = "";
  let params: Record<string, unknown>;
  try {
    params = JSON.parse(publishParamsJson.value || "{}") as Record<string, unknown>;
  } catch {
    publishResult.value = "参数 JSON 无法解析";
    return;
  }
  const body: Record<string, unknown> = {
    route: publishRoute.value,
    biz: {
      apiPath: publishAction.value,
      method: publishMethod.value,
      params,
    },
  };
  if (publishReferId.value.trim()) {
    body.referId = publishReferId.value.trim();
  }
  try {
    publishResult.value = await postPublish(body);
  } catch (e) {
    publishResult.value = String(e);
  }
}
</script>

<template>
  <section class="card">
    <h2>发布消息</h2>
    <p class="hint">
      对应 <code>POST /api/v1/communication/publish</code>，组装
      <code>sdk.Request</code> 中的 <code>route</code> 与 <code>biz</code>。
    </p>
    <div class="grid">
      <label>
        <span>Route（Topic）</span>
        <input v-model="publishRoute" type="text" autocomplete="off" />
      </label>
      <label>
        <span>HTTP Method</span>
        <select v-model="publishMethod">
          <option>POST</option>
          <option>GET</option>
        </select>
      </label>
      <label class="wide">
        <span>biz.apiPath（回调路径后缀）</span>
        <input v-model="publishAction" type="text" />
      </label>
      <label class="wide">
        <span>biz.params（JSON）</span>
        <textarea v-model="publishParamsJson" class="mono" rows="6" spellcheck="false" />
      </label>
      <label class="wide">
        <span>referId（可选）</span>
        <input v-model="publishReferId" type="text" placeholder="关联上游 messageId" />
      </label>
    </div>
    <button type="button" class="primary" @click="doPublish">发送</button>
    <pre v-if="publishResult" class="out mono">{{ publishResult }}</pre>
  </section>
</template>
