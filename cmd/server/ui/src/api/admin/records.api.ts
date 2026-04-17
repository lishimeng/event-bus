import type { DataRecordRow } from "@/types/models/admin";
import { API_V1, apiFetch, apiUrl, readJsonOrThrow, unwrapApiData } from "@/utils/request";

export type RecordsQuery = { limit: number; offset: number; source?: string };

export async function fetchRecords(
  q: RecordsQuery,
): Promise<{ total: number; items: DataRecordRow[]; error?: string }> {
  const params = new URLSearchParams();
  params.set("limit", String(q.limit));
  params.set("offset", String(q.offset));
  if (q.source?.trim()) params.set("source", q.source.trim());
  try {
    const res = await apiFetch(apiUrl(`${API_V1}/admin/records?${params.toString()}`));
    const j = await readJsonOrThrow(res, "消息记录");
    if (!res.ok) {
      return { total: 0, items: [], error: `HTTP ${res.status} ${String(j.message ?? j.Message ?? "")}` };
    }
    const payload = unwrapApiData<{ total?: number; items?: DataRecordRow[] }>(j);
    return {
      total: Number(payload?.total ?? 0),
      items: Array.isArray(payload?.items) ? payload.items : [],
    };
  } catch (e) {
    return { total: 0, items: [], error: String(e) };
  }
}
