"use client";

import { AlertTriangle, BellRing, Clock, PackageOpen, XCircle } from "lucide-react";
import type { ReactNode } from "react";
import { useReservationStats } from "@/lib/reservationHooks";

interface StatsCardsProps {
  /** Opens the reservations list on the matching filter. */
  onShowOverdue: (days: 1 | 5 | 7) => void;
  onShowActive: () => void;
  onShowCancelled: () => void;
}

function Card({
  icon,
  label,
  value,
  hint,
  onClick,
  tone = "text-gray-900",
}: {
  icon: ReactNode;
  label: string;
  value: ReactNode;
  hint?: ReactNode;
  onClick?: () => void;
  tone?: string;
}) {
  const content = (
    <>
      <div className="flex items-center gap-2 text-xs font-medium uppercase tracking-wide text-gray-500">
        {icon}
        {label}
      </div>
      <div className={`mt-1 text-2xl font-bold ${tone}`}>{value}</div>
      {hint && <div className="mt-1 text-xs text-gray-500">{hint}</div>}
    </>
  );
  const className = "rounded-lg bg-white p-4 text-left shadow";
  return onClick ? (
    <button
      type="button"
      onClick={onClick}
      className={`${className} hover:ring-2 hover:ring-blue-200`}
    >
      {content}
    </button>
  ) : (
    <div className={className}>{content}</div>
  );
}

export function StatsCards({ onShowOverdue, onShowActive, onShowCancelled }: StatsCardsProps) {
  const { data } = useReservationStats();

  if (!data) {
    return (
      <div className="grid animate-pulse grid-cols-2 gap-4 lg:grid-cols-5" aria-hidden="true">
        {Array.from({ length: 5 }, (_, index) => (
          // biome-ignore lint/suspicious/noArrayIndexKey: static placeholders
          <div key={index} className="h-24 rounded-lg bg-gray-200" />
        ))}
      </div>
    );
  }

  return (
    <div className="grid grid-cols-2 gap-4 lg:grid-cols-5">
      <Card
        icon={<PackageOpen className="h-4 w-4" />}
        label="Emprunts en cours"
        value={data.active_loans}
        hint={`${data.ongoing_slots} créneau${data.ongoing_slots > 1 ? "x" : ""} en cours · ${data.upcoming_slots} à venir`}
        onClick={onShowActive}
      />
      <Card
        icon={<Clock className="h-4 w-4" />}
        label="Emprunts > 1 j"
        value={data.overdue_1d}
        tone={data.overdue_1d > 0 ? "text-amber-600" : undefined}
        onClick={() => onShowOverdue(1)}
      />
      <Card
        icon={<AlertTriangle className="h-4 w-4" />}
        label="> 5 j / > 7 j"
        value={
          <>
            {data.overdue_5d}
            <span className="text-gray-300"> / </span>
            <span className={data.overdue_7d > 0 ? "text-red-600" : undefined}>
              {data.overdue_7d}
            </span>
          </>
        }
        onClick={() => onShowOverdue(data.overdue_7d > 0 ? 7 : 5)}
      />
      <Card
        icon={<BellRing className="h-4 w-4" />}
        label="Files d'attente"
        value={data.waitlist_entries}
        hint={
          data.waitlist_awaiting_response > 0
            ? `${data.waitlist_awaiting_response} en attente de réponse`
            : "Personne n'a la main"
        }
      />
      <Card
        icon={<XCircle className="h-4 w-4" />}
        label="Annulations (30 j)"
        value={data.cancelled_30d}
        hint={
          data.unavailable_items > 0
            ? `${data.unavailable_items} élément(s) indisponible(s)`
            : undefined
        }
        onClick={onShowCancelled}
      />
    </div>
  );
}
