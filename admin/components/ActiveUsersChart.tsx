"use client";

import { Bar, BarChart, CartesianGrid, ResponsiveContainer, Tooltip, XAxis, YAxis } from "recharts";

interface ActiveUsersChartProps {
  data: Array<{ date: string; count: number }>;
}

export default function ActiveUsersChart({ data }: ActiveUsersChartProps) {
  return (
    <ResponsiveContainer width="100%" height="100%">
      <BarChart data={data} margin={{ top: 5, right: 10, left: 0, bottom: 5 }}>
        <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" />
        <XAxis
          dataKey="date"
          tickFormatter={(d) =>
            new Date(d).toLocaleDateString("fr-FR", { month: "short", day: "numeric" })
          }
          stroke="#6b7280"
          fontSize={12}
          minTickGap={24}
        />
        <YAxis stroke="#6b7280" fontSize={12} width={40} allowDecimals={false} />
        <Tooltip
          labelFormatter={(label) =>
            new Date(String(label)).toLocaleDateString("fr-FR", {
              weekday: "long",
              day: "numeric",
              month: "long",
              year: "numeric",
            })
          }
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
