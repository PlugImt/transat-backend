"use client";

import { Bar, BarChart, CartesianGrid, ResponsiveContainer, Tooltip, XAxis, YAxis } from "recharts";
import type { ActiveUsersPeriod } from "@/lib/types";

interface ActiveUsersChartProps {
  data: Array<{ date: string; count: number }>;
  period: ActiveUsersPeriod;
}

function quarterOf(date: Date) {
  return Math.floor(date.getMonth() / 3) + 1;
}

function formatTick(dateString: string, period: ActiveUsersPeriod) {
  const date = new Date(dateString);
  switch (period) {
    case "week":
      return date.toLocaleDateString("fr-FR", { day: "numeric", month: "short" });
    case "month":
      return date.toLocaleDateString("fr-FR", { month: "short", year: "2-digit" });
    case "quarter":
      return `T${quarterOf(date)} ${date.getFullYear()}`;
    case "year":
      return `${date.getFullYear()}`;
    default:
      return date.toLocaleDateString("fr-FR", { month: "short", day: "numeric" });
  }
}

function formatLabel(dateString: string, period: ActiveUsersPeriod) {
  const date = new Date(dateString);
  switch (period) {
    case "week":
      return `Semaine du ${date.toLocaleDateString("fr-FR", { day: "numeric", month: "long", year: "numeric" })}`;
    case "month":
      return date.toLocaleDateString("fr-FR", { month: "long", year: "numeric" });
    case "quarter":
      return `Trimestre T${quarterOf(date)} ${date.getFullYear()}`;
    case "year":
      return `${date.getFullYear()}`;
    default:
      return date.toLocaleDateString("fr-FR", {
        weekday: "long",
        day: "numeric",
        month: "long",
        year: "numeric",
      });
  }
}

export default function ActiveUsersChart({ data, period }: ActiveUsersChartProps) {
  return (
    <ResponsiveContainer width="100%" height="100%">
      <BarChart data={data} margin={{ top: 5, right: 10, left: 0, bottom: 5 }}>
        <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" />
        <XAxis
          dataKey="date"
          tickFormatter={(d) => formatTick(d, period)}
          stroke="#6b7280"
          fontSize={12}
          minTickGap={24}
        />
        <YAxis stroke="#6b7280" fontSize={12} width={40} allowDecimals={false} />
        <Tooltip
          labelFormatter={(label) => formatLabel(String(label), period)}
          formatter={(value) => [value, "Utilisateurs actifs"]}
          contentStyle={{
            backgroundColor: "#fff",
            border: "1px solid #e5e7eb",
            borderRadius: "6px",
            fontSize: "12px",
          }}
        />
        <Bar dataKey="count" fill="#7c3aed" radius={[2, 2, 0, 0]} maxBarSize={14} />
      </BarChart>
    </ResponsiveContainer>
  );
}
