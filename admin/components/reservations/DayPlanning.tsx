"use client";

import clsx from "clsx";
import { CalendarDays, ChevronLeft, ChevronRight } from "lucide-react";
import { useMemo, useState } from "react";
import { useSlotReservations } from "@/lib/reservationHooks";
import type { AdminItem, AdminReservation } from "@/lib/reservationTypes";
import { CancelReservationModal, type ReservationTarget } from "./ActionModals";
import { cellOf, hourLabel, hoursWith, longDay, shiftDay, toUtcDate } from "./planningDates";
import { targetOf } from "./ReservationsList";
import {
  AvailabilityBadge,
  Badge,
  EmptyState,
  fieldClass,
  formatDateTime,
  formatDuration,
  formatTime,
  fullName,
  ITEM_STATUS_LABELS,
  parisDay,
  parisHour,
  ReservationStatusBadge,
  SkeletonRows,
} from "./ui";

const overdueLevel = (seconds: number) =>
  seconds >= 7 * 86_400 ? 7 : seconds >= 5 * 86_400 ? 5 : seconds >= 86_400 ? 1 : 0;

const byCategoryThenName = (a: AdminItem, b: AdminItem) =>
  a.category_path.localeCompare(b.category_path, "fr") || a.name.localeCompare(b.name, "fr");

interface DayPlanningProps {
  items: AdminItem[];
  onOpenItem: (itemId: number) => void;
}

/** One day, every item with slots as a column and every hour as a row, then the items lent without slot. */
export function DayPlanning({ items, onOpenItem }: DayPlanningProps) {
  const slotItems = useMemo(
    () => items.filter((item) => item.slot).sort(byCategoryThenName),
    [items],
  );
  const loanItems = useMemo(
    () =>
      items
        .filter((item) => !item.slot)
        .sort(
          (a, b) =>
            (b.current?.elapsed_seconds ?? -1) - (a.current?.elapsed_seconds ?? -1) ||
            byCategoryThenName(a, b),
        ),
    [items],
  );

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

  const now = new Date();
  const nowHour = parisHour(now.toISOString());
  const isPast = (hour: number) => day < today || (day === today && hour !== 24 && hour < nowHour);

  return (
    <div className="space-y-6">
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
            className={`${fieldClass} w-auto`}
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
                        className="sticky top-0 z-10 min-w-[9.5rem] border-b border-l border-gray-200 bg-gray-50 px-2 py-2 text-left font-medium"
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
                              "h-9 border-b border-l border-gray-100 p-0.5 align-middle",
                              unavailable && !reservation && "bg-gray-50",
                            )}
                            title={
                              unavailable && !reservation
                                ? ITEM_STATUS_LABELS[item.status]
                                : undefined
                            }
                          >
                            {reservation ? (
                              <button
                                type="button"
                                disabled={!reservation.cancellable}
                                onClick={() => setTarget(targetOf(reservation))}
                                title={`${fullName(reservation.first_name, reservation.last_name)} · ${reservation.email}\n${item.name} · ${formatTime(reservation.start_date)}${
                                  reservation.cancellable ? "\nCliquer pour annuler" : ""
                                }`}
                                className={clsx(
                                  "h-full w-full truncate rounded px-1.5 text-left font-medium",
                                  reservation.status === "ongoing" &&
                                    "bg-green-100 text-green-900 hover:bg-red-100 hover:text-red-900",
                                  reservation.status === "upcoming" &&
                                    "bg-blue-100 text-blue-900 hover:bg-red-100 hover:text-red-900",
                                  reservation.status !== "ongoing" &&
                                    reservation.status !== "upcoming" &&
                                    "bg-gray-100 text-gray-500",
                                )}
                              >
                                {fullName(reservation.first_name, reservation.last_name) ||
                                  reservation.email}
                              </button>
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
        <p className="flex flex-wrap items-center gap-x-4 gap-y-1 text-xs text-gray-500">
          <span>Heures de Paris. Cliquez sur un créneau réservé pour l'annuler.</span>
          <span className="inline-flex items-center gap-1">
            <i className="h-3 w-3 rounded bg-blue-100" /> À venir
          </span>
          <span className="inline-flex items-center gap-1">
            <i className="h-3 w-3 rounded bg-green-100" /> En cours
          </span>
          <span className="inline-flex items-center gap-1">
            <i className="h-3 w-3 rounded bg-gray-100 ring-1 ring-gray-200" /> Passé ou indisponible
          </span>
        </p>
      </div>

      <section className="space-y-3">
        <h3 className="text-sm font-semibold text-gray-800">
          Éléments prêtés sans créneau
          <span className="ml-2 font-normal text-gray-500">{loanItems.length}</span>
        </h3>
        <div className="overflow-x-auto rounded-lg bg-white shadow">
          {loanItems.length === 0 ? (
            <p className="p-6 text-center text-sm text-gray-500">
              Aucun élément prêté sans créneau.
            </p>
          ) : (
            <table className="min-w-full divide-y divide-gray-100 text-sm">
              <thead className="bg-gray-50 text-left text-xs uppercase text-gray-500">
                <tr>
                  <th className="px-4 py-2">Élément</th>
                  <th className="px-4 py-2">État</th>
                  <th className="px-4 py-2">Emprunteur</th>
                  <th className="px-4 py-2">Depuis</th>
                  <th className="px-4 py-2">Retard</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-100">
                {loanItems.map((item) => (
                  <tr key={item.id}>
                    <td className="px-4 py-2">
                      <button
                        type="button"
                        onClick={() => onOpenItem(item.id)}
                        className="font-medium text-gray-900 hover:text-blue-700 hover:underline"
                      >
                        {item.name}
                      </button>
                      {item.category_path && (
                        <div className="text-xs text-gray-500">{item.category_path}</div>
                      )}
                    </td>
                    <td className="px-4 py-2">
                      <AvailabilityBadge value={item.availability} />
                    </td>
                    <td className="px-4 py-2 text-gray-700">
                      {item.current
                        ? fullName(item.current.first_name, item.current.last_name) ||
                          item.current.email
                        : "—"}
                    </td>
                    <td className="px-4 py-2 text-gray-600">
                      {item.current
                        ? `${formatDateTime(item.current.start_date)} · ${formatDuration(item.current.elapsed_seconds)}`
                        : "—"}
                    </td>
                    <td className="px-4 py-2">
                      {item.current ? (
                        <ReservationStatusBadge
                          status="active"
                          overdueLevel={overdueLevel(item.current.elapsed_seconds)}
                        />
                      ) : null}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      </section>

      <CancelReservationModal target={target} onClose={() => setTarget(null)} />
    </div>
  );
}
