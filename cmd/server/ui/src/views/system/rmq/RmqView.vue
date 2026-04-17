<script setup lang="ts">
import { ref, onMounted } from "vue";
import { fetchRmqConfig, postRmqConfig } from "@/api/admin/rmq.api";

const rmqEndpoint = ref("169.254.24.18:10911");
const rmqAppId = ref("");
const rmqSecret = ref("");
const rmqMessageGroup = ref("event-bus-pub");
const rmqConsumerGroup = ref("event-bus-cg");
const rmqSubscriberTopic = ref("demo-topic");
const rmqPublisherTopicsCsv = ref("demo-topic");
const rmqResult = ref("");
const rmqFromDbHint = ref("");
const rmqLoadErr = ref("");

async function loadRmqFromDb() {
  rmqLoadErr.value = "";
  rmqFromDbHint.value = "";
  const r = await fetchRmqConfig();
  if (r.hint) rmqFromDbHint.value = r.hint;
  if (r.error) rmqLoadErr.value = r.error;
  if (!r.configured) {
    return;
  }
  const cfg = r.config;
  if (!cfg) {
    rmqFromDbHint.value = r.hint || "已标记 configured，但 config 解析失败，请检查库中 JSON。";
    return;
  }
  rmqEndpoint.value = cfg.endpoint ?? "";
  rmqAppId.value = cfg.appId ?? "";
  rmqSecret.value = cfg.secret ?? "";
  rmqMessageGroup.value = cfg.publisher?.messageGroup ?? "event-bus-publisher";
  const sub0 = cfg.subscribers?.[0];
  rmqConsumerGroup.value = sub0?.consumerGroup ?? "";
  rmqSubscriberTopic.value = sub0?.topic ?? "";
  const topics = cfg.publisher?.topics ?? [];
  rmqPublisherTopicsCsv.value = topics.length ? topics.join(", ") : "";
  rmqFromDbHint.value = r.httpMessage
    ? `已从库加载（${r.httpMessage}）`
    : "已从数据库加载当前 rmq_config 到表单。";
}

async function saveRmqToDb() {
  rmqResult.value = "";
  const topics = rmqPublisherTopicsCsv.value
    .split(",")
    .map((s) => s.trim())
    .filter(Boolean);
  const body = {
    endpoint: rmqEndpoint.value.trim(),
    appId: rmqAppId.value.trim(),
    secret: rmqSecret.value.trim(),
    publisher: {
      messageGroup: rmqMessageGroup.value.trim() || "event-bus-publisher",
      topics,
    },
    subscribers: [
      {
        consumerGroup: rmqConsumerGroup.value.trim(),
        topic: rmqSubscriberTopic.value.trim(),
      },
    ],
  };
  if (!body.endpoint) {
    rmqResult.value = "请填写 endpoint";
    return;
  }
  // if (!body.subscribers[0].consumerGroup || !body.subscribers[0].topic) {
  //   rmqResult.value = "请填写 consumerGroup";
  //   return;
  // }
  try {
    const { ok, text } = await postRmqConfig(body);
    rmqResult.value = text;
    if (ok) void loadRmqFromDb();
  } catch (e) {
    rmqResult.value = String(e);
  }
}

onMounted(() => {
  void loadRmqFromDb();
});
</script>

<template>
  <section class="card">
    <h2>RMQ 配置入库</h2>
    <p class="hint">
      <code>GET /api/v1/admin/rmq_config</code> 读取库中当前配置并填入下方表单；
      <code>POST /api/v1/admin/rmq_config</code> 写入
      <code>sys_config</code>（name=rmq_config）。重启服务后 RocketMQ 连接才生效。
    </p>
    <div class="toolbar">
      <button type="button" class="secondary" @click="loadRmqFromDb">从数据库重新加载</button>
    </div>
    <p v-if="rmqFromDbHint" class="hint okhint">{{ rmqFromDbHint }}</p>
    <p v-if="rmqLoadErr" class="warn">{{ rmqLoadErr }}</p>
    <div class="grid">
      <label class="wide">
        <span>endpoint（v5 客户端多为 Proxy gRPC 端口，非 10911）</span>
        <input v-model="rmqEndpoint" type="text" autocomplete="off" />
      </label>
      <label>
        <span>appId（未开启鉴权可留空；已开启鉴权时必填）</span>
        <input v-model="rmqAppId" type="text" />
      </label>
      <label>
        <span>secret（未开启鉴权可留空；已开启鉴权时必填）</span>
        <input v-model="rmqSecret" type="text" />
      </label>
      <label>
        <span>publisher.messageGroup</span>
        <input v-model="rmqMessageGroup" type="text" />
      </label>
      <label>
        <span>consumerGroup</span>
        <input v-model="rmqConsumerGroup" type="text" />
      </label>
    </div>
    <button type="button" class="primary" @click="saveRmqToDb">保存到 sys_config</button>
    <pre v-if="rmqResult" class="out mono">{{ rmqResult }}</pre>
  </section>
</template>
