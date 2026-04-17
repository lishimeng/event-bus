import type { ChannelRow } from "@/types/models/admin";
import { API_V1, apiFetch, apiUrl, readJsonOrThrow, unwrapApiData } from "@/utils/request";

export async function fetchChannelList(): Promise<{ rows: ChannelRow[]; error?: string }> {
  try {
    const res = await apiFetch(apiUrl(`${API_V1}/admin/channel`));
    const j = await readJsonOrThrow(res, "通道列表");
    if (!res.ok) {
      return { rows: [], error: `HTTP ${res.status} ${String(j.message ?? j.Message ?? "")}` };
    }
    const rows = unwrapApiData<ChannelRow[]>(j);
    return { rows: Array.isArray(rows) ? rows : [] };
  } catch (e) {
    return { rows: [], error: String(e) };
  }
}

export async function postChannel(body: Record<string, unknown>): Promise<{ ok: boolean; text: string }> {
  const res = await apiFetch(apiUrl(`${API_V1}/admin/channel`), {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  const text = await res.text();
  return { ok: res.ok, text: `${res.status}\n${text}` };
}
