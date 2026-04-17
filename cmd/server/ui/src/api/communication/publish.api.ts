import { API_V1, apiFetch, apiUrl } from "@/utils/request";

export async function postPublish(body: Record<string, unknown>): Promise<string> {
  const res = await apiFetch(apiUrl(`${API_V1}/communication/publish`), {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  const text = await res.text();
  return `${res.status} ${res.statusText}\n${text}`;
}
