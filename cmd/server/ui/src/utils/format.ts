/** 将后端 RFC3339 格式化为本地墙钟时间（YYYY-MM-DD HH:mm:ss），不显示 T 与 +08:00 */
export function formatRecordTime(iso?: string): string {
  if (!iso?.trim()) return "—";
  const d = new Date(iso.trim());
  if (Number.isNaN(d.getTime())) return iso;
  const y = d.getFullYear();
  const m = String(d.getMonth() + 1).padStart(2, "0");
  const day = String(d.getDate()).padStart(2, "0");
  const h = String(d.getHours()).padStart(2, "0");
  const min = String(d.getMinutes()).padStart(2, "0");
  const sec = String(d.getSeconds()).padStart(2, "0");
  return `${y}-${m}-${day} ${h}:${min}:${sec}`;
}

export function prettyJson(s: string) {
  if (!s?.trim()) return "（空）";
  try {
    return JSON.stringify(JSON.parse(s) as unknown, null, 2);
  } catch {
    return s;
  }
}
