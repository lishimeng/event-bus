/**
 * 与 admin channel 接口一致：传输层为 PEM 整段 UTF-8 的 StdEncoding Base64；
 * 表单内展示 PEM 原文（含 BEGIN/END），保存时再编码。
 */

function utf8ToBase64(pem: string): string {
  const bytes = new TextEncoder().encode(pem);
  let binary = "";
  for (let i = 0; i < bytes.length; i++) {
    binary += String.fromCharCode(bytes[i]);
  }
  return btoa(binary);
}

function base64ToUtf8(b64: string): string {
  const clean = b64.replace(/\s/g, "");
  const binary = atob(clean);
  const bytes = new Uint8Array(binary.length);
  for (let i = 0; i < binary.length; i++) {
    bytes[i] = binary.charCodeAt(i);
  }
  return new TextDecoder().decode(bytes);
}

const pemHeader = /-----BEGIN [A-Z0-9 ]+-----/;

/** 列表/详情里的字段 → 文本框展示（PEM 原文） */
export function channelKeyToDisplay(apiField: string | undefined): string {
  const t = apiField?.trim() ?? "";
  if (!t) return "";
  if (pemHeader.test(t)) return t;
  try {
    const pem = base64ToUtf8(t);
    if (pemHeader.test(pem)) return pem.replace(/\r\n/g, "\n").trimEnd();
  } catch {
    /* 非合法 Base64 时原样展示，便于用户修正 */
  }
  return t;
}

/** 文本框内容 → POST body（API 要求的 Base64 一行） */
export function channelKeyToApiPayload(input: string): string {
  const t = input.trim();
  if (!t) return "";
  if (pemHeader.test(t)) return utf8ToBase64(t);
  try {
    const maybePem = base64ToUtf8(t);
    if (pemHeader.test(maybePem)) return utf8ToBase64(maybePem);
  } catch {
    /* 按已是 API Base64 处理 */
  }
  return t.replace(/\s/g, "");
}

/** 从当前框内容得到 PEM 原文（已含 Base64→PEM）；无法识别时返回 null */
export function channelKeyToPemForFile(input: string): string | null {
  const t = input.trim();
  if (!t) return null;
  if (pemHeader.test(t)) return t.replace(/\r\n/g, "\n").trimEnd();
  try {
    const pem = base64ToUtf8(t);
    if (pemHeader.test(pem)) return pem.replace(/\r\n/g, "\n").trimEnd();
  } catch {
    /* ignore */
  }
  return null;
}

export function downloadPemFile(filename: string, pem: string): void {
  const blob = new Blob([pem], { type: "application/x-pem-file;charset=utf-8" });
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = filename;
  a.rel = "noopener";
  a.click();
  URL.revokeObjectURL(url);
}
