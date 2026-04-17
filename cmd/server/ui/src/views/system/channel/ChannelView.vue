<script setup lang="ts">
import { ref, onMounted } from "vue";
import { fetchChannelList, postChannel } from "@/api/admin/channel.api";
import type { ChannelRow } from "@/types/models/admin";
import {
  channelKeyToApiPayload,
  channelKeyToDisplay,
  channelKeyToPemForFile,
  downloadPemFile,
} from "@/utils/channelKeyForm";

const chName = ref("示例通道");
const chRoute = ref("demo-topic");
const chCategory = ref<1 | 2>(2);
const chCallback = ref("http://127.0.0.1:9999");
const chPrivateKeyText = ref("");
const chPublicKeyText = ref("");
const chResult = ref("");

const editTarget = ref<ChannelRow | null>(null);
const editName = ref("");
const editRoute = ref("");
const editCategory = ref<1 | 2>(2);
const editCallback = ref("");
const editPrivateKeyText = ref("");
const editPublicKeyText = ref("");
const editResult = ref("");
const editDialogVisible = ref(false);
const channelList = ref<ChannelRow[]>([]);
const channelListErr = ref("");
const channelListLoading = ref(false);

async function loadChannelList() {
  channelListErr.value = "";
  channelListLoading.value = true;
  try {
    const { rows, error } = await fetchChannelList();
    if (error) {
      channelListErr.value = error;
      return;
    }
    channelList.value = rows;
  } finally {
    channelListLoading.value = false;
  }
}

function loadEditForm(row: ChannelRow) {
  editTarget.value = row;
  editName.value = row.name;
  editRoute.value = row.router;
  editCategory.value = row.category === 1 ? 1 : 2;
  editCallback.value = row.callback || "";
  editPrivateKeyText.value = channelKeyToDisplay(row.privateKey);
  editPublicKeyText.value = channelKeyToDisplay(row.publicKey);
  editResult.value = row.privateKey || row.publicKey
    ? `已载入「${row.name}」到表单（私钥/公钥已解码为 PEM 原文）`
    : `已载入「${row.name}」到表单（库中无密钥材料）`;
  editDialogVisible.value = true;
}

function closeEditDialog() {
  editDialogVisible.value = false;
}

function channelKeyDownloadStamp(): string {
  return new Date().toISOString().slice(0, 19).replace(/[:T]/g, "-");
}

function downloadChannelPrivatePemFile(input: string, setResult: (s: string) => void) {
  const pem = channelKeyToPemForFile(input);
  if (!pem) {
    setResult("私钥为空或无法解析为 PEM 原文");
    return;
  }
  downloadPemFile(`channel-private-${channelKeyDownloadStamp()}.pem`, pem);
}

function downloadChannelPublicPemFile(input: string, setResult: (s: string) => void) {
  const pem = channelKeyToPemForFile(input);
  if (!pem) {
    setResult("公钥为空或无法解析为 PEM 原文");
    return;
  }
  downloadPemFile(`channel-public-${channelKeyDownloadStamp()}.pem`, pem);
}

async function saveChannel(
  model: {
    code?: string;
    name: string;
    route: string;
    category: 1 | 2;
    callback: string;
    privateKeyText: string;
    publicKeyText: string;
  },
  setResult: (s: string) => void,
) {
  setResult("");
  const body: Record<string, unknown> = {
    name: model.name,
    route: model.route,
    category: model.category,
    callback: model.callback,
  };
  if (model.code) {
    body.code = model.code;
  }
  if (model.privateKeyText.trim()) {
    body.privateKey = channelKeyToApiPayload(model.privateKeyText);
  }
  if (model.publicKeyText.trim()) {
    body.publicKey = channelKeyToApiPayload(model.publicKeyText);
  }
  try {
    const { ok, text } = await postChannel(body);
    setResult(text);
    if (ok) void loadChannelList();
  } catch (e) {
    const msg = String(e);
    setResult(msg);
  }
}

function saveCreate() {
  void saveChannel(
    {
      name: chName.value,
      route: chRoute.value,
      category: chCategory.value,
      callback: chCallback.value,
      privateKeyText: chPrivateKeyText.value,
      publicKeyText: chPublicKeyText.value,
    },
    (s) => (chResult.value = s),
  );
}

function saveEdit() {
  void saveChannel(
    {
      code: editTarget.value?.code,
      name: editName.value,
      route: editRoute.value,
      category: editCategory.value,
      callback: editCallback.value,
      privateKeyText: editPrivateKeyText.value,
      publicKeyText: editPublicKeyText.value,
    },
    (s) => (editResult.value = s),
  );
}

onMounted(() => {
  void loadChannelList();
});
</script>

<template>
  <section class="card">
    <h2>通道配置</h2>
    <p class="hint">
      <code>GET /api/v1/admin/channel</code> 列出库中全部通道；
      <code>POST /api/v1/admin/channel</code> 按 <strong>route + category</strong> 新建或更新（同一 topic 可同时存在 PublishTo 与 Subscribe 各一条）。开发时请求走
      <code>VITE_API_BASE</code>（见 <code>ui/.env.development</code>）直连后端，避免代理返回 HTML。
      category：<strong>1</strong> = 订阅，<strong>2</strong> = 发布。点「编辑」会打开弹窗加载该条记录；请勿把管理接口暴露公网。
    </p>
    <div class="toolbar">
      <button type="button" class="secondary" :disabled="channelListLoading" @click="loadChannelList">
        {{ channelListLoading ? "加载中…" : "刷新列表" }}
      </button>
    </div>
    <p v-if="channelListErr" class="warn">{{ channelListErr }}</p>
    <div v-if="channelList.length" class="table-wrap">
      <table class="data-table">
        <thead>
          <tr>
            <th>名称</th>
            <th>Route</th>
            <th>类别</th>
            <th>TLS</th>
            <th>密钥</th>
            <th />
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in channelList" :key="row.id">
            <td>{{ row.name }}</td>
            <td class="mono">{{ row.router }}</td>
            <td>{{ row.categoryLabel }}</td>
            <td>{{ row.useSecurity ? "开" : "关" }}</td>
            <td class="small">
              <span v-if="row.hasPrivateKey">私</span>
              <span v-if="row.hasPublicKey">公</span>
              <span v-if="!row.hasPrivateKey && !row.hasPublicKey">—</span>
            </td>
            <td>
              <button type="button" class="linkish" @click="loadEditForm(row)">编辑</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <p v-else-if="!channelListLoading" class="hint">暂无通道记录。</p>

    <h3 class="subh">新建</h3>
    <div class="grid">
      <label>
        <span>名称</span>
        <input v-model="chName" type="text" />
      </label>
      <label>
        <span>Route / Topic</span>
        <input v-model="chRoute" type="text" />
      </label>
      <label>
        <span>category</span>
        <select v-model.number="chCategory">
          <option :value="1">1 — Subscribe</option>
          <option :value="2">2 — PublishTo</option>
        </select>
      </label>
      <label class="wide">
        <span>Callback 基地址（订阅类常用）</span>
        <input v-model="chCallback" type="text" />
      </label>
      <div class="wide key-field-row">
        <span class="key-field-label">私钥 PEM（可选）</span>
        <textarea
          v-model="chPrivateKeyText"
          class="mono key-field-textarea"
          rows="6"
          spellcheck="false"
          placeholder="粘贴 -----BEGIN RSA PRIVATE KEY----- … 整段 PEM 原文。"
        />
        <div class="key-field-actions">
          <button
            type="button"
            class="secondary key-field-dl"
            :disabled="!chPrivateKeyText.trim()"
            @click="downloadChannelPrivatePemFile(chPrivateKeyText, (s) => (chResult = s))"
          >
            下载私钥 .pem
          </button>
        </div>
      </div>
      <div class="wide key-field-row">
        <span class="key-field-label">公钥 PEM（可选）</span>
        <textarea
          v-model="chPublicKeyText"
          class="mono key-field-textarea"
          rows="5"
          spellcheck="false"
          placeholder="粘贴 -----BEGIN PUBLIC KEY----- … 整段 PEM 原文。"
        />
        <div class="key-field-actions">
          <button
            type="button"
            class="secondary key-field-dl"
            :disabled="!chPublicKeyText.trim()"
            @click="downloadChannelPublicPemFile(chPublicKeyText, (s) => (chResult = s))"
          >
            下载公钥 .pem
          </button>
        </div>
      </div>
    </div>
    <button type="button" class="primary" @click="saveCreate">新建并保存到数据库</button>
    <pre v-if="chResult" class="out mono">{{ chResult }}</pre>

    <Teleport to="body">
      <div v-if="editDialogVisible && editTarget" class="edit-modal-mask" @click.self="closeEditDialog">
        <div class="edit-modal">
          <div class="edit-modal-header">
            <h3>编辑通道</h3>
            <button type="button" class="secondary" @click="closeEditDialog">关闭</button>
          </div>
          <p class="hint">
            当前编辑：<strong>{{ editTarget.name }}</strong>
            （{{ editTarget.router }} / {{ editTarget.categoryLabel }}）
          </p>
          <div class="grid">
        <label>
          <span>名称</span>
          <input v-model="editName" type="text" />
        </label>
        <label>
          <span>Route / Topic</span>
          <input v-model="editRoute" type="text" />
        </label>
        <label>
          <span>category</span>
          <select v-model.number="editCategory">
            <option :value="1">1 — Subscribe</option>
            <option :value="2">2 — PublishTo</option>
          </select>
        </label>
        <label class="wide">
          <span>Callback 基地址（订阅类常用）</span>
          <input v-model="editCallback" type="text" />
        </label>
        <div class="wide key-field-row">
          <span class="key-field-label">私钥 PEM（可选）</span>
          <textarea
            v-model="editPrivateKeyText"
            class="mono key-field-textarea"
            rows="6"
            spellcheck="false"
            placeholder="编辑时请粘贴 RSA PRIVATE KEY 的 PEM 原文。"
          />
          <div class="key-field-actions">
            <button
              type="button"
              class="secondary key-field-dl"
              :disabled="!editPrivateKeyText.trim()"
              @click="downloadChannelPrivatePemFile(editPrivateKeyText, (s) => (editResult = s))"
            >
              下载私钥 .pem
            </button>
          </div>
        </div>
        <div class="wide key-field-row">
          <span class="key-field-label">公钥 PEM（可选）</span>
          <textarea
            v-model="editPublicKeyText"
            class="mono key-field-textarea"
            rows="5"
            spellcheck="false"
            placeholder="编辑时请粘贴 PUBLIC KEY 的 PEM 原文。"
          />
          <div class="key-field-actions">
            <button
              type="button"
              class="secondary key-field-dl"
              :disabled="!editPublicKeyText.trim()"
              @click="downloadChannelPublicPemFile(editPublicKeyText, (s) => (editResult = s))"
            >
              下载公钥 .pem
            </button>
          </div>
        </div>
          </div>
          <div class="edit-modal-footer">
            <button type="button" class="primary" @click="saveEdit">保存编辑到数据库</button>
          </div>
          <pre v-if="editResult" class="out mono">{{ editResult }}</pre>
        </div>
      </div>
    </Teleport>
  </section>
</template>

<style scoped>
.edit-modal-mask {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 1rem;
}

.edit-modal {
  width: min(920px, 100%);
  max-height: 90vh;
  overflow: auto;
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: 12px;
  padding: 1rem 1rem 1.25rem;
}

.edit-modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 0.5rem;
}

.edit-modal-header h3 {
  margin: 0;
}

.edit-modal-footer {
  margin-top: 0.75rem;
}
</style>
