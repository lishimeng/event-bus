import type { RmqCfg } from "@/types/models/admin";
import { API_V1, apiFetch, apiUrl, readJsonOrThrow, unwrapApiData } from "@/utils/request";

export type RmqLoadResult = {
  configured: boolean;
  config?: RmqCfg;
  hint?: string;
  error?: string;
  httpMessage?: string;
  httpCode: number;
  ok: boolean;
};

export async function fetchRmqConfig(): Promise<RmqLoadResult> {
  try {
    const res = await apiFetch(apiUrl(`${API_V1}/admin/rmq_config`));
    const j = await readJsonOrThrow(res, "RMQ 配置");
    const code = Number(j.code ?? j.Code ?? res.status);
    const msg = String(j.message ?? j.Message ?? "");
    const payload = unwrapApiData<{ configured?: boolean; config?: RmqCfg }>(j);
    if (!payload) {
      return { configured: false, httpCode: code, ok: res.ok, error: "响应无 data" };
    }
    if (!payload.configured) {
      return {
        configured: false,
        httpCode: code,
        ok: res.ok,
        hint: "数据库中尚无 rmq_config，以下为当前表单默认值或上次内容。",
        error: !res.ok && code !== 200 ? msg || `HTTP ${code}` : undefined,
      };
    }
    const cfg = payload.config;
    if (!cfg) {
      return {
        configured: true,
        httpCode: code,
        ok: res.ok,
        hint: msg || "已标记 configured，但 config 解析失败，请检查库中 JSON。",
      };
    }
    return {
      configured: true,
      config: cfg,
      httpCode: code,
      ok: res.ok,
      httpMessage: msg,
    };
  } catch (e) {
    return { configured: false, httpCode: 0, ok: false, error: String(e) };
  }
}

export async function postRmqConfig(body: Record<string, unknown>): Promise<{ ok: boolean; text: string }> {
  const res = await apiFetch(apiUrl(`${API_V1}/admin/rmq_config`), {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  const text = await res.text();
  return { ok: res.ok, text: `${res.status}\n${text}` };
}
