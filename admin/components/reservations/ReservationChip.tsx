"use client";

import { useState } from "react";
import type { AdminReservation } from "@/lib/reservationTypes";
import { chipStyle, PHASE_LABELS, phaseOf, solidColor } from "./planningColors";
import { formatDateTime, formatTime, fullName, ReminderBadges, ReservationStatusBadge } from "./ui";

interface ReservationChipProps {
  reservation: AdminReservation;
  hue: number;
  onSelect: (reservation: AdminReservation) => void;
}

/** A reservation in the calendar: item and person always visible, everything else on hover. */
export function ReservationChip({ reservation, hue, onSelect }: ReservationChipProps) {
  const [anchor, setAnchor] = useState<DOMRect | null>(null);
  const phase = phaseOf(reservation);
  const person = fullName(reservation.first_name, reservation.last_name) || reservation.email;

  return (
    <>
      <button
        type="button"
        aria-disabled={!reservation.cancellable}
        aria-label={`${reservation.item_name}, ${person}, ${PHASE_LABELS[phase]}`}
        onClick={() => reservation.cancellable && onSelect(reservation)}
        onMouseEnter={(event) => setAnchor(event.currentTarget.getBoundingClientRect())}
        onMouseLeave={() => setAnchor(null)}
        onFocus={(event) => setAnchor(event.currentTarget.getBoundingClientRect())}
        onBlur={() => setAnchor(null)}
        style={chipStyle(hue, phase)}
        className={`flex h-full min-h-9 w-full flex-col justify-center overflow-hidden rounded px-1.5 py-0.5 text-left leading-tight ${
          reservation.cancellable ? "cursor-pointer" : "cursor-default"
        }`}
      >
        <span className="truncate text-[11px] font-semibold">{reservation.item_name}</span>
        <span className="truncate text-[10px]">{person}</span>
      </button>
      {anchor && <HoverCard anchor={anchor} reservation={reservation} hue={hue} />}
    </>
  );
}

const CARD_WIDTH = 288;

function HoverCard({
  anchor,
  reservation,
  hue,
}: {
  anchor: DOMRect;
  reservation: AdminReservation;
  hue: number;
}) {
  const left = Math.max(8, Math.min(anchor.left, window.innerWidth - CARD_WIDTH - 8));
  const below = anchor.bottom + 220 < window.innerHeight;
  const position = below
    ? { top: anchor.bottom + 6 }
    : { bottom: window.innerHeight - anchor.top + 6 };

  return (
    <div
      role="tooltip"
      style={{ left, width: CARD_WIDTH, ...position }}
      className="pointer-events-none fixed z-50 space-y-2 rounded-lg border border-gray-200 bg-white p-3 text-xs text-gray-700 shadow-xl"
    >
      <div className="flex items-start gap-2">
        <i
          className="mt-0.5 h-3 w-3 shrink-0 rounded-sm"
          style={{ backgroundColor: solidColor(hue) }}
        />
        <div className="min-w-0">
          <p className="truncate text-sm font-semibold text-gray-900">{reservation.item_name}</p>
          {reservation.club_name && (
            <p className="truncate text-gray-500">{reservation.club_name}</p>
          )}
        </div>
        <span className="ml-auto shrink-0">
          <ReservationStatusBadge status={reservation.status} />
        </span>
      </div>
      <dl className="grid grid-cols-[4.5rem_1fr] gap-x-2 gap-y-1">
        <dt className="text-gray-500">Réservé par</dt>
        <dd className="min-w-0 font-medium text-gray-900">
          {fullName(reservation.first_name, reservation.last_name) || "—"}
          <span className="block truncate font-normal text-gray-500">{reservation.email}</span>
        </dd>
        <dt className="text-gray-500">Créneau</dt>
        <dd>
          {formatDateTime(reservation.start_date)}
          {reservation.end_date ? ` – ${formatTime(reservation.end_date)}` : ""}
        </dd>
        <dt className="text-gray-500">Rappels</dt>
        <dd>
          <ReminderBadges reminders={reservation.reminders} />
        </dd>
      </dl>
      {reservation.cancellable && <p className="text-gray-500">Cliquer pour annuler.</p>}
    </div>
  );
}
