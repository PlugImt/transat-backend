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

/** Consecutive hourly slots of one person on one item, shown as a single long event. */
export interface SlotRun {
  day: string;
  hour: number;
  /** The slots in chronological order. */
  reservations: AdminReservation[];
}

export const mergeRuns = (reservations: AdminReservation[]): SlotRun[] => {
  const sorted = reservations
    .map((reservation) => ({ reservation, cell: cellOf(reservation.start_date) }))
    .sort(
      (a, b) =>
        a.reservation.item_id - b.reservation.item_id ||
        a.reservation.email.localeCompare(b.reservation.email) ||
        a.cell.day.localeCompare(b.cell.day) ||
        a.cell.hour - b.cell.hour,
    );
  const runs: SlotRun[] = [];
  let last: { run: SlotRun; hour: number } | null = null;
  for (const { reservation, cell } of sorted) {
    const previous = last?.run.reservations[0];
    if (
      last &&
      previous &&
      previous.item_id === reservation.item_id &&
      previous.email === reservation.email &&
      last.run.day === cell.day &&
      last.hour + 1 === cell.hour
    ) {
      last.run.reservations.push(reservation);
      last.hour = cell.hour;
    } else {
      const run = { day: cell.day, hour: cell.hour, reservations: [reservation] };
      runs.push(run);
      last = { run, hour: cell.hour };
    }
  }
  return runs;
};
