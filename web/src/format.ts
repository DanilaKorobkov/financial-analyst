// Утилиты форматирования значений из ответа Connect-сервера для отображения.
//
// Источник — типизированный protobuf-JSON: int64 приходит строкой,
// timestamp — строкой ISO-8601, числовые поля — числами.

const ISO_DATE_RE = /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}/;
const NUMERIC_STRING_RE = /^-?\d+(\.\d+)?$/;

export function humanizeKey(key: string): string {
  // snake_case или camelCase → "Title Case".
  const withSpaces = key
    .replace(/_/g, ' ')
    .replace(/([a-z0-9])([A-Z])/g, '$1 $2')
    .trim();
  return withSpaces.charAt(0).toUpperCase() + withSpaces.slice(1);
}

export function formatValue(v: unknown): string {
  if (v === null || v === undefined || v === '') return '—';
  if (typeof v === 'boolean') return v ? '✓' : '✗';

  if (typeof v === 'string') {
    if (ISO_DATE_RE.test(v)) return formatDate(v);
    if (NUMERIC_STRING_RE.test(v)) return formatNumber(Number(v));
    return v;
  }

  if (typeof v === 'number') return formatNumber(v);

  if (Array.isArray(v)) return `${v.length} элем.`;
  if (typeof v === 'object') return JSON.stringify(v);
  return String(v);
}

function formatDate(iso: string): string {
  const d = new Date(iso);
  if (Number.isNaN(d.getTime())) return iso;
  return d.toLocaleDateString('ru-RU', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  });
}

function formatNumber(n: number): string {
  if (!Number.isFinite(n)) return String(n);
  if (Number.isInteger(n)) return n.toLocaleString('ru-RU');
  const abs = Math.abs(n);
  const digits = abs >= 100 ? 2 : abs >= 1 ? 3 : 4;
  return n.toLocaleString('ru-RU', { maximumFractionDigits: digits });
}

export function isEmptyValue(v: unknown): boolean {
  if (v === null || v === undefined) return true;
  if (typeof v === 'string' && v === '') return true;
  if (typeof v === 'number' && v === 0) return false; // 0 — валидное значение
  if (Array.isArray(v) && v.length === 0) return true;
  if (typeof v === 'object' && v !== null && Object.keys(v).length === 0) return true;
  return false;
}
