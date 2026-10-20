import type { AdminReservation } from "@/lib/reservationTypes";
import { parisDay, parisHour } from "./ui";

/** The app offers hourly slots from 08:00 to 01:00 (Paris time); the row 24 is the slot starting at midnight. */
export const DEFAULT_HOURS = Array.from({ length: 17 }, (_, index) => 8 + index);

const DAY_MS = 86_400_000;

export const toUtcDate = (day: string) => new Date(`${day}T00:00:00Z`);
export const shiftDay = (day: string, days: number) =>
  new Date(toUtcDate(day).getTime() + days * DAY_MS).toISOString().slice(0, 10);
export const mondayOf = (day: string) => shiftDay(day, -((toUtcDate(day).getUTCDay() + 6) % 7));

export const weekday = new Intl.DateTimeFormat("fr-FR", {
  timeZone: "UTC",
  weekday: "short",
  day: "2-digit",
  month: "2-digit",
});

export const longDay = new Intl.DateTimeFormat("fr-FR", {
  timeZone: "UTC",
  weekday: "long",
  day: "numeric",
  month: "long",
  year: "numeric",
});

/** Paris calendar day and hour row of a slot; the midnight slot belongs to the evening before. */
export const cellOf = (iso: string) => {
  const hour = parisHour(iso);
  const day = parisDay(iso);
  return hour === 0 ? { day: shiftDay(day, -1), hour: 24 } : { day, hour };
};

/** The usual hours, plus any other hour a reservation actually sits on, so that no slot is ever hidden. */
export const hoursWith = (reservations: AdminReservation[], day?: (value: string) => boolean) => {
  const hours = new Set(DEFAULT_HOURS);
  for (const reservation of reservations) {
    if (reservation.status === "cancelled") continue;
    const cell = cellOf(reservation.start_date);
    if (!day || day(cell.day)) hours.add(cell.hour);
  }
  return [...hours].sort((a, b) => a - b);
};

export const hourLabel = (hour: number) => `${String(hour % 24).padStart(2, "0")}:00`;
