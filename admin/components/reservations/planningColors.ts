import type { CSSProperties } from "react";
import type { AdminItem, AdminReservation } from "@/lib/reservationTypes";

export type Phase = "past" | "ongoing" | "future";

export const PHASE_LABELS: Record<Phase, string> = {
  past: "Passée",
  ongoing: "En cours",
  future: "À venir",
};

/** Opacity of a reservation's fill: faded once over, full while it takes place, in between before it starts. */
export const PHASE_ALPHA: Record<Phase, number> = { past: 0.22, future: 0.55, ongoing: 1 };

export const phaseOf = (reservation: AdminReservation): Phase =>
  reservation.status === "ongoing" || reservation.status === "active"
    ? "ongoing"
    : reservation.status === "upcoming"
      ? "future"
      : "past";

/** One stable hue per item, spread by the golden angle so that neighbours in the list never look alike. */
export const itemHues = (items: AdminItem[]) => {
  const hues = new Map<number, number>();
  [...items]
    .filter((item) => item.slot)
    .sort((a, b) => a.id - b.id)
    .forEach((item, index) => {
      hues.set(item.id, Math.round((index * 137.508) % 360));
    });
  return hues;
};

export const solidColor = (hue: number) => `hsl(${hue} 70% 42%)`;

export const chipStyle = (hue: number, phase: Phase): CSSProperties => ({
  backgroundColor: `hsl(${hue} 80% ${phase === "ongoing" ? 62 : 58}% / ${PHASE_ALPHA[phase]})`,
  borderLeft: `4px solid ${solidColor(hue)}`,
  color: `hsl(${hue} 70% ${phase === "ongoing" ? 10 : 18}%)`,
  boxShadow: phase === "ongoing" ? `0 0 0 1px ${solidColor(hue)}` : undefined,
});
