"use client";

import { useId } from "react";
import {
  Area,
  Bar,
  CartesianGrid,
  ComposedChart,
  Legend,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";

interface UserGrowthChartProps {
  data: Array<{
    date: string;
    count: number;
    cumulativeCount: number;
  }>;
}

const SERIES_LABELS: Record<string, string> = {
  cumulativeCount: "Total des comptes",
  count: "Nouveaux comptes",
};

export default function UserGrowthChart({ data }: UserGrowthChartProps) {
  const gradientId = useId();
  const formatDate = (dateString: string) =>
    new Date(dateString).toLocaleDateString("fr-FR", { month: "short", day: "numeric" });

  const formatTooltipDate = (dateString: string) =>
    new Date(dateString).toLocaleDateString("fr-FR", {
      weekday: "long",
      year: "numeric",
      month: "long",
      day: "numeric",
    });

  return (
    <ResponsiveContainer width="100%" height="100%">
      <ComposedChart data={data} margin={{ top: 5, right: 10, left: 0, bottom: 5 }}>
        <defs>
          <linearGradient id={gradientId} x1="0" y1="0" x2="0" y2="1">
            <stop offset="0%" stopColor="#2563eb" stopOpacity={0.25} />
            <stop offset="100%" stopColor="#2563eb" stopOpacity={0} />
          </linearGradient>
        </defs>
        <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" />
        <XAxis
          dataKey="date"
          tickFormatter={formatDate}
          stroke="#6b7280"
          fontSize={12}
          minTickGap={24}
        />
        <YAxis yAxisId="total" stroke="#6b7280" fontSize={12} width={40} />
        <YAxis yAxisId="daily" orientation="right" stroke="#6b7280" fontSize={12} width={30} />
        <Tooltip
          labelFormatter={(label) => formatTooltipDate(String(label ?? ""))}
          formatter={(value, name) => [value, SERIES_LABELS[String(name)] ?? name]}
          contentStyle={{
            backgroundColor: "#fff",
            border: "1px solid #e5e7eb",
            borderRadius: "6px",
            fontSize: "12px",
          }}
        />
        <Legend formatter={(value) => SERIES_LABELS[String(value)] ?? value} />
        <Bar yAxisId="daily" dataKey="count" fill="#059669" radius={[2, 2, 0, 0]} maxBarSize={14} />
        <Area
          yAxisId="total"
          type="monotone"
          dataKey="cumulativeCount"
          stroke="#2563eb"
          strokeWidth={2}
          fill={`url(#${gradientId})`}
          dot={false}
          activeDot={{ r: 4 }}
        />
      </ComposedChart>
    </ResponsiveContainer>
  );
}
