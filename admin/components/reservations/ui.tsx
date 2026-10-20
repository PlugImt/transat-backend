"use client";

import clsx from "clsx";
import { Loader2 } from "lucide-react";
import type { ButtonHTMLAttributes, ReactNode } from "react";
import type {
  CancelKind,
  ItemAvailability,
  ItemStatus,
  ReminderMilestone,
  ReservationStatus,
} from "@/lib/reservationTypes";

const PARIS = "Europe/Paris";

const dateTime = new Intl.DateTimeFormat("fr-FR", {
  timeZone: PARIS,
  day: "2-digit",
  month: "2-digit",
  year: "numeric",
  hour: "2-digit",
  minute: "2-digit",
});
const dateOnly = new Intl.DateTimeFormat("fr-FR", {
  timeZone: PARIS,
  weekday: "short",
  day: "2-digit",
  month: "2-digit",
});
const timeOnly = new Intl.DateTimeFormat("fr-FR", {
  timeZone: PARIS,
  hour: "2-digit",
  minute: "2-digit",
});
const dayKey = new Intl.DateTimeFormat("en-CA", { timeZone: PARIS });
const hourOnly = new Intl.DateTimeFormat("en-GB", {
  timeZone: PARIS,
  hour: "2-digit",
  hour12: false,
});

export const formatDateTime = (iso: string) => dateTime.format(new Date(iso));
export const formatDay = (iso: string) => dateOnly.format(new Date(iso));
export const formatTime = (iso: string) => timeOnly.format(new Date(iso));
/** Paris calendar day (YYYY-MM-DD) of an instant. */
export const parisDay = (iso: string | Date) =>
  dayKey.format(typeof iso === "string" ? new Date(iso) : iso);
export const parisHour = (iso: string) => Number(hourOnly.format(new Date(iso))) % 24;

export const formatDuration = (seconds: number) => {
  const minutes = Math.max(0, Math.floor(seconds / 60));
  const days = Math.floor(minutes / 1440);
  const hours = Math.floor((minutes % 1440) / 60);
  if (days > 0) return hours > 0 ? `${days} j ${hours} h` : `${days} j`;
  if (hours > 0) return `${hours} h ${String(minutes % 60).padStart(2, "0")}`;
  return `${minutes} min`;
};

export const fullName = (first?: string, last?: string) => `${first ?? ""} ${last ?? ""}`.trim();

type Tone = "green" | "amber" | "blue" | "red" | "gray" | "purple";

const TONES: Record<Tone, string> = {
  green: "bg-green-100 text-green-800",
  amber: "bg-amber-100 text-amber-800",
  blue: "bg-blue-100 text-blue-800",
  red: "bg-red-100 text-red-800",
  gray: "bg-gray-100 text-gray-700",
  purple: "bg-purple-100 text-purple-800",
};

export function Badge({
  tone,
  children,
  title,
}: {
  tone: Tone;
  children: ReactNode;
  title?: string;
}) {
  return (
    <span
      title={title}
      className={clsx(
        "inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-xs font-medium whitespace-nowrap",
        TONES[tone],
      )}
    >
      {children}
    </span>
  );
}

const AVAILABILITY: Record<ItemAvailability, { label: string; tone: Tone }> = {
  available: { label: "Disponible", tone: "green" },
  borrowed: { label: "Emprunté", tone: "amber" },
  occupied: { label: "Occupé", tone: "amber" },
  maintenance: { label: "Maintenance", tone: "purple" },
  inactive: { label: "Inactif", tone: "gray" },
};

export const AVAILABILITY_LABELS = Object.fromEntries(
  Object.entries(AVAILABILITY).map(([key, value]) => [key, value.label]),
) as Record<ItemAvailability, string>;

export function AvailabilityBadge({ value }: { value: ItemAvailability }) {
  const { label, tone } = AVAILABILITY[value];
  return <Badge tone={tone}>{label}</Badge>;
}

export const ITEM_STATUS_LABELS: Record<ItemStatus, string> = {
  active: "Actif",
  maintenance: "En maintenance",
  inactive: "Inactif",
};

const RESERVATION_STATUS: Record<ReservationStatus, { label: string; tone: Tone }> = {
  upcoming: { label: "À venir", tone: "blue" },
  ongoing: { label: "En cours", tone: "green" },
  active: { label: "En cours", tone: "green" },
  past: { label: "Terminée", tone: "gray" },
  cancelled: { label: "Annulée", tone: "red" },
};

export function ReservationStatusBadge({
  status,
  overdueLevel,
  cancelKind,
}: {
  status: ReservationStatus;
  overdueLevel?: number;
  cancelKind?: CancelKind;
}) {
  if (status === "active" && overdueLevel) {
    return <Badge tone={overdueLevel >= 7 ? "red" : "amber"}>En retard +{overdueLevel} j</Badge>;
  }
  if (status === "past" && cancelKind === "returned_by_admin") {
    return (
      <Badge tone="gray" title="Marqué comme rendu par un administrateur">
        Rendu (admin)
      </Badge>
    );
  }
  const { label, tone } = RESERVATION_STATUS[status];
  return <Badge tone={tone}>{label}</Badge>;
}

export const CANCEL_KIND_LABELS: Record<CancelKind, string> = {
  cancelled: "Annulée par l'utilisateur",
  cancelled_by_admin: "Annulée par un administrateur",
  returned_by_admin: "Marquée comme rendue par un administrateur",
};

const REMINDER_LABELS: Record<ReminderMilestone, string> = {
  slot_1h: "1 h avant",
  loan_d1: "J+1",
  loan_d5: "J+5",
  loan_d7: "J+7",
};

export function ReminderBadges({ reminders }: { reminders: ReminderMilestone[] }) {
  if (reminders.length === 0) return <span className="text-xs text-gray-400">Aucun envoyé</span>;
  return (
    <span className="inline-flex flex-wrap gap-1">
      {reminders.map((milestone) => (
        <Badge key={milestone} tone="blue" title="Rappel envoyé">
          {REMINDER_LABELS[milestone] ?? milestone}
        </Badge>
      ))}
    </span>
  );
}

export function LoadingButton({
  loading,
  children,
  className,
  disabled,
  size = "md",
  ...props
}: ButtonHTMLAttributes<HTMLButtonElement> & { loading?: boolean; size?: "sm" | "md" }) {
  return (
    <button
      type="button"
      {...props}
      disabled={disabled || loading}
      className={clsx(
        "inline-flex items-center justify-center gap-2 rounded-md text-sm font-medium disabled:cursor-not-allowed disabled:opacity-60",
        size === "sm" ? "px-2.5 py-1" : "px-3 py-2",
        className,
      )}
    >
      {loading && <Loader2 className="h-4 w-4 animate-spin" />}
      {children}
    </button>
  );
}

export function SkeletonRows({ rows = 6, columns = 5 }: { rows?: number; columns?: number }) {
  return (
    <div className="animate-pulse divide-y divide-gray-100" aria-hidden="true">
      {Array.from({ length: rows }, (_, row) => (
        // biome-ignore lint/suspicious/noArrayIndexKey: static placeholder rows
        <div key={row} className="flex gap-4 px-4 py-3">
          {Array.from({ length: columns }, (_, column) => (
            // biome-ignore lint/suspicious/noArrayIndexKey: static placeholder cells
            <div key={column} className="h-4 flex-1 rounded bg-gray-200" />
          ))}
        </div>
      ))}
    </div>
  );
}

export function EmptyState({
  icon,
  title,
  description,
}: {
  icon: ReactNode;
  title: string;
  description?: string;
}) {
  return (
    <div className="flex flex-col items-center justify-center gap-2 px-6 py-12 text-center text-gray-500">
      <div className="text-gray-300">{icon}</div>
      <p className="font-medium text-gray-700">{title}</p>
      {description && <p className="max-w-md text-sm">{description}</p>}
    </div>
  );
}

export const apiErrorMessage = (error: unknown, fallback: string) => {
  const data = (error as { response?: { data?: { error?: string } } })?.response?.data;
  return data?.error ?? fallback;
};

export const fieldClass =
  "w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500";
