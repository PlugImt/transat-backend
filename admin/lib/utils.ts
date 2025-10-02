export function toLocalDateTimeInputValue(date: Date) {
  const pad = (n: number) => `${n}`.padStart(2, "0");
  const y = date.getFullYear();
  const m = pad(date.getMonth() + 1);
  const d = pad(date.getDate());
  const h = pad(date.getHours());
  const mi = pad(date.getMinutes());
  return `${y}-${m}-${d}T${h}:${mi}`;
}
