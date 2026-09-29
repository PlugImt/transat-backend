"use client";

import type { ActivityHourPoint } from "@/lib/types";

interface ActivityHeatmapChartProps {
  data: ActivityHourPoint[];
}

const DAY_LABELS: Record<number, string> = {
  1: "Lun",
  2: "Mar",
  3: "Mer",
  4: "Jeu",
  5: "Ven",
  6: "Sam",
  7: "Dim",
};

const HOURS = Array.from({ length: 24 }, (_, i) => i);
const DAYS = [1, 2, 3, 4, 5, 6, 7];

// Interpole entre violet clair et violet soutenu selon l'intensité (0-1)
function cellColor(ratio: number) {
  if (ratio <= 0) return "#f3f4f6";
  const start = { r: 237, g: 233, b: 254 };
  const end = { r: 91, g: 33, b: 182 };
  const r = Math.round(start.r + (end.r - start.r) * ratio);
  const g = Math.round(start.g + (end.g - start.g) * ratio);
  const b = Math.round(start.b + (end.b - start.b) * ratio);
  return `rgb(${r}, ${g}, ${b})`;
}

export default function ActivityHeatmapChart({ data }: ActivityHeatmapChartProps) {
  const counts = new Map<string, number>();
  let max = 0;
  for (const point of data) {
    counts.set(`${point.dayOfWeek}-${point.hour}`, point.count);
    if (point.count > max) max = point.count;
  }

  if (data.length === 0) {
    return <div className="text-sm text-gray-500">Pas assez de données</div>;
  }

  return (
    <div className="overflow-x-auto">
      <div
        className="inline-grid gap-[3px] min-w-[640px]"
        style={{ gridTemplateColumns: "36px repeat(24, minmax(20px, 1fr))" }}
      >
        <div />
        {HOURS.map((h) => (
          <div key={`hour-${h}`} className="text-[10px] text-gray-400 text-center">
            {h % 3 === 0 ? `${h}h` : ""}
          </div>
        ))}
        {DAYS.flatMap((dow) => [
          <div key={`label-${dow}`} className="text-xs text-gray-500 flex items-center">
            {DAY_LABELS[dow]}
          </div>,
          ...HOURS.map((h) => {
            const count = counts.get(`${dow}-${h}`) ?? 0;
            const ratio = max > 0 ? count / max : 0;
            return (
              <div
                key={`cell-${dow}-${h}`}
                className="aspect-square rounded-sm"
                style={{ backgroundColor: cellColor(ratio) }}
                title={`${DAY_LABELS[dow]} ${h}h : ${count} utilisateur${count > 1 ? "s" : ""} actif${count > 1 ? "s" : ""}`}
              />
            );
          }),
        ])}
      </div>
      <div className="flex items-center gap-2 mt-3 text-xs text-gray-500">
        <span>Moins actif</span>
        <div className="flex gap-[3px]">
          {[0, 0.25, 0.5, 0.75, 1].map((ratio) => (
            <div
              key={ratio}
              className="h-3 w-3 rounded-sm"
              style={{ backgroundColor: cellColor(ratio) }}
            />
          ))}
        </div>
        <span>Plus actif</span>
      </div>
    </div>
  );
}
