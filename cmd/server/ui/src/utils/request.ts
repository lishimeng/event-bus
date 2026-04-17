import { API_V1 } from "@/constants/api";

/** 开发环境见 .env.development 的 VITE_API_BASE（直连 Go）；生产构建为空，走同域 /api。 */
export function apiBase(): string {
  return (import.meta.env.VITE_API_BASE ?? "").replace(/\/$/, "");
}

export function apiUrl(path: string): string {
  const p = path.startsWith("/") ? path : `/${path}`;
  return `${apiBase()}${p}`;
}

/** 禁止缓存 API 响应。同源 /api 若被磁盘缓存到旧版 HTML，会误报「收到 HTML 而非 JSON」。 */
export function apiFetch(input: string, init?: RequestInit): Promise<Response> {
  return fetch(input, { ...init, cache: "no-store" });
}

export function unwrapApiData<T>(j: Record<string, unknown>): T | undefined {
  const d = j.data;
  if (d === undefined || d === null) return undefined;
  return d as T;
}

/** Iris 对带尾斜杠的 /path/ 常返回 301 且 body 为 HTML；用无尾斜杠 URL 可避免被当成 JSON 解析失败。 */
export async function readJsonOrThrow(res: Response, context: string): Promise<Record<string, unknown>> {
  const text = await res.text();
  const trimmed = text.trim();
  if (trimmed.startsWith("<") || trimmed.startsWith("<!")) {
    throw new Error(
      `${context}：收到 HTML 而非 JSON。常见原因：① 后端未启动或代理未生效；② 浏览器曾把 /api 缓存成首页 HTML（DevTools 若显示「来自磁盘缓存」请 Ctrl+F5）。前端已对 API 请求使用 cache:no-store，Go 侧也会返回 Cache-Control:no-store。`,
    );
  }
  try {
    return JSON.parse(text) as Record<string, unknown>;
  } catch {
    throw new Error(`${context}：响应不是合法 JSON（前 80 字符）${text.slice(0, 80)}`);
  }
}

export { API_V1 };
