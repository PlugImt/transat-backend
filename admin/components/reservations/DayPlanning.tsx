"use client";

import clsx from "clsx";
import { CalendarDays, ChevronLeft, ChevronRight } from "lucide-react";
import { useMemo, useState } from "react";
import { useSlotReservations } from "@/lib/reservationHooks";
import type { AdminItem, AdminReservation } from "@/lib/reservationTypes";
import { CancelReservationModal, type ReservationTarget } from "./ActionModals";
import { solidColor } from "./planningColors";
import { cellOf, hourLabel, hoursWith, longDay, shiftDay, toUtcDate } from "./planningDates";
import { ReservationChip } from "./ReservationChip";
import { targetOf } from "./ReservationsList";
import { Badge, EmptyState, ITEM_STATUS_LABELS, parisDay, parisHour, SkeletonRows } from "./ui";

interface DayPlanningProps {
  /** The items with slots, in display order. */
  slotItems: AdminItem[];
  hues: Map<number, number>;
  onOpenItem: (itemId: number) => void;
}

/** One day: a column per item, a row per hour. */
export function DayPlanning({ slotItems, hues, onOpenItem }: DayPlanningProps) {
  const today = parisDay(new Date());
  const [day, setDay] = useState(today);
  const [target, setTarget] = useState<ReservationTarget | null>(null);

  // Wider than the day: the days are cut in Paris time on this side.
  const { data, isLoading, isError } = useSlotReservations(shiftDay(day, -1), shiftDay(day, 2));

  const { cells, hours, itemsBooked } = useMemo(() => {
    const map = new Map<string, AdminReservation>();
    const booked = new Set<number>();
    for (const reservation of data ?? []) {
      if (reservation.status === "cancelled") continue;
      const cell = cellOf(reservation.start_date);
      if (cell.day !== day) continue;
      map.set(`${reservation.item_id}|${cell.hour}`, reservation);
      booked.add(reservation.item_id);
    }
    return {
      cells: map,
      hours: hoursWith(data ?? [], (value) => value === day),
      itemsBooked: booked.size,
    };
  }, [data, day]);

  const nowHour = parisHour(new Date().toISOString());
  const isPast = (hour: number) => day < today || (day === today && hour !== 24 && hour < nowHour);

  return (
    <div className="space-y-4">
      <div className="flex flex-wrap items-center gap-3">
        <div className="flex items-center gap-1">
          <button
            type="button"
            onClick={() => setDay(shiftDay(day, -1))}
            className="rounded p-2 hover:bg-gray-100"
            aria-label="Jour précédent"
          >
            <ChevronLeft className="h-4 w-4" />
          </button>
          <button
            type="button"
            onClick={() => setDay(today)}
            className="rounded-md px-3 py-1.5 text-sm font-medium text-blue-700 hover:bg-blue-50"
          >
            Aujourd'hui
          </button>
          <button
            type="button"
            onClick={() => setDay(shiftDay(day, 1))}
            className="rounded p-2 hover:bg-gray-100"
            aria-label="Jour suivant"
          >
            <ChevronRight className="h-4 w-4" />
          </button>
        </div>
        <input
          type="date"
          value={day}
          onChange={(event) => event.target.value && setDay(event.target.value)}
          className="rounded-md border border-gray-300 px-3 py-1.5 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500"
          aria-label="Choisir un jour"
        />
        <span className="text-sm font-medium capitalize text-gray-800">
          {longDay.format(toUtcDate(day))}
        </span>
        <span className="ml-auto text-sm text-gray-600">
          {isLoading
            ? "Chargement…"
            : `${cells.size} créneau${cells.size > 1 ? "x" : ""} réservé${cells.size > 1 ? "s" : ""} sur ${itemsBooked} élément${itemsBooked > 1 ? "s" : ""}`}
        </span>
      </div>

      {slotItems.length === 0 ? (
        <div className="rounded-lg bg-white shadow">
          <EmptyState
            icon={<CalendarDays className="h-10 w-10" />}
            title="Aucun élément à créneaux"
            description="Les salles et le matériel réservés par créneaux apparaissent ici."
          />
        </div>
      ) : (
        <div className="max-h-[70vh] overflow-auto rounded-lg bg-white shadow">
          {isLoading ? (
            <SkeletonRows rows={10} columns={6} />
          ) : isError ? (
            <p className="p-8 text-center text-sm text-red-600">
              Impossible de charger le planning.
            </p>
          ) : (
            <table className="min-w-full border-separate border-spacing-0 text-xs">
              <thead>
                <tr>
                  <th className="sticky top-0 left-0 z-20 w-16 border-b border-gray-200 bg-gray-50 px-2 py-2 text-right font-medium text-gray-500">
                    Heure
                  </th>
                  {slotItems.map((item) => (
                    <th
                      key={item.id}
                      scope="col"
                      style={{ borderTop: `3px solid ${solidColor(hues.get(item.id) ?? 0)}` }}
                      className="sticky top-0 z-10 min-w-[10rem] border-b border-l border-gray-200 bg-gray-50 px-2 py-2 text-left font-medium"
                    >
                      <button
                        type="button"
                        onClick={() => onOpenItem(item.id)}
                        className="block max-w-full truncate text-left text-gray-800 hover:text-blue-700 hover:underline"
                        title={`${item.name}${item.category_path ? ` — ${item.category_path}` : ""}`}
                      >
                        {item.name}
                      </button>
                      <div className="mt-0.5 flex items-center gap-1 font-normal text-gray-500">
                        <span className="truncate">{item.location || item.category_path}</span>
                        {item.status !== "active" && (
                          <Badge tone={item.status === "maintenance" ? "purple" : "gray"}>
                            {ITEM_STATUS_LABELS[item.status]}
                          </Badge>
                        )}
                      </div>
                    </th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {hours.map((hour) => (
                  <tr key={hour}>
                    <th
                      scope="row"
                      className={clsx(
                        "sticky left-0 z-10 border-b border-gray-100 bg-white px-2 py-1 text-right font-normal",
                        day === today && hour === nowHour
                          ? "font-semibold text-blue-700"
                          : "text-gray-400",
                      )}
                    >
                      {hourLabel(hour)}
                    </th>
                    {slotItems.map((item) => {
                      const reservation = cells.get(`${item.id}|${hour}`);
                      const unavailable = item.status !== "active";
                      return (
                        <td
                          key={item.id}
                          className={clsx(
                            "h-11 border-b border-l border-gray-100 p-0.5 align-middle",
                            unavailable && !reservation && "bg-gray-50",
                          )}
                          title={
                            unavailable && !reservation
                              ? ITEM_STATUS_LABELS[item.status]
                              : undefined
                          }
                        >
                          {reservation ? (
                            <ReservationChip
                              reservation={reservation}
                              hue={hues.get(item.id) ?? 0}
                              onSelect={(value) => setTarget(targetOf(value))}
                            />
                          ) : (
                            !unavailable &&
                            !isPast(hour) && (
                              <span className="block px-1.5 text-gray-300">Libre</span>
                            )
                          )}
                        </td>
                      );
                    })}
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      )}

      <CancelReservationModal target={target} onClose={() => setTarget(null)} />
    </div>
  );
}
