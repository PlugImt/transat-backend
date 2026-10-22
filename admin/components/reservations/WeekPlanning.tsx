"use client";

import clsx from "clsx";
import { CalendarDays, ChevronLeft, ChevronRight } from "lucide-react";
import { useMemo, useState } from "react";
import { useSlotReservations } from "@/lib/reservationHooks";
import type { AdminItem, AdminReservation } from "@/lib/reservationTypes";
import { CancelReservationModal, type ReservationTarget } from "./ActionModals";
import { solidColor } from "./planningColors";
import {
  cellOf,
  hourLabel,
  hoursWith,
  mergeRuns,
  mondayOf,
  shiftDay,
  type SlotRun,
  toUtcDate,
  weekday,
} from "./planningDates";
import { ReservationChip } from "./ReservationChip";
import { targetOfRun } from "./ReservationsList";
import { EmptyState, parisDay, parisHour, SkeletonRows } from "./ui";

interface WeekPlanningProps {
  slotItems: AdminItem[];
  hues: Map<number, number>;
}

/** One week: a column per day, a row per hour, and in each cell every item booked at that hour. */
export function WeekPlanning({ slotItems, hues }: WeekPlanningProps) {
  const today = parisDay(new Date());
  const [monday, setMonday] = useState(() => mondayOf(today));
  const [hidden, setHidden] = useState<Set<number>>(new Set());
  const [target, setTarget] = useState<ReservationTarget | null>(null);

  const days = useMemo(
    () => Array.from({ length: 7 }, (_, index) => shiftDay(monday, index)),
    [monday],
  );

  // Wider than the week: the days are cut in Paris time on this side.
  const { data, isLoading, isError } = useSlotReservations(
    shiftDay(monday, -1),
    shiftDay(monday, 8),
  );

  const { cells, hours, perDay } = useMemo(() => {
    const map = new Map<string, { run: SlotRun; first: boolean }[]>();
    const counts = new Map<string, number>();
    const names = new Map(slotItems.map((item) => [item.id, item.name]));
    const visible: AdminReservation[] = [];
    for (const reservation of data ?? []) {
      if (reservation.status === "cancelled" || hidden.has(reservation.item_id)) continue;
      if (!days.includes(cellOf(reservation.start_date).day)) continue;
      visible.push(reservation);
    }
    for (const run of mergeRuns(visible)) {
      counts.set(run.day, (counts.get(run.day) ?? 0) + run.reservations.length);
      // The event sits on its first hour; the following hours only carry a continuation strip.
      for (let offset = 0; offset < run.reservations.length; offset++) {
        const key = `${run.day}|${run.hour + offset}`;
        map.set(key, [...(map.get(key) ?? []), { run, first: offset === 0 }]);
      }
    }
    const nameOf = (run: SlotRun) =>
      names.get(run.reservations[0].item_id) ?? run.reservations[0].item_name;
    for (const list of map.values()) {
      list.sort((a, b) => nameOf(a.run).localeCompare(nameOf(b.run), "fr"));
    }
    return {
      cells: map,
      hours: hoursWith(
        (data ?? []).filter((reservation) => !hidden.has(reservation.item_id)),
        (value) => days.includes(value),
      ),
      perDay: counts,
    };
  }, [data, days, hidden, slotItems]);

  const toggle = (itemId: number) =>
    setHidden((previous) => {
      const next = new Set(previous);
      if (!next.delete(itemId)) next.add(itemId);
      return next;
    });

  const nowHour = parisHour(new Date().toISOString());

  if (slotItems.length === 0) {
    return (
      <div className="rounded-lg bg-white shadow">
        <EmptyState
          icon={<CalendarDays className="h-10 w-10" />}
          title="Aucun élément à créneaux"
          description="Les salles et le matériel réservés par créneaux apparaissent ici."
        />
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <div className="flex flex-wrap items-center gap-3">
        <div className="flex items-center gap-1">
          <button
            type="button"
            onClick={() => setMonday(shiftDay(monday, -7))}
            className="rounded p-2 hover:bg-gray-100"
            aria-label="Semaine précédente"
          >
            <ChevronLeft className="h-4 w-4" />
          </button>
          <button
            type="button"
            onClick={() => setMonday(mondayOf(today))}
            className="rounded-md px-3 py-1.5 text-sm font-medium text-blue-700 hover:bg-blue-50"
          >
            Cette semaine
          </button>
          <button
            type="button"
            onClick={() => setMonday(shiftDay(monday, 7))}
            className="rounded p-2 hover:bg-gray-100"
            aria-label="Semaine suivante"
          >
            <ChevronRight className="h-4 w-4" />
          </button>
        </div>
        <span className="text-sm font-medium text-gray-800">
          Semaine du {weekday.format(toUtcDate(days[0]))} au {weekday.format(toUtcDate(days[6]))}
        </span>
      </div>

      <fieldset className="flex flex-wrap items-center gap-1.5">
        <legend className="sr-only">Éléments affichés</legend>
        {slotItems.map((item) => {
          const shown = !hidden.has(item.id);
          return (
            <button
              key={item.id}
              type="button"
              aria-pressed={shown}
              onClick={() => toggle(item.id)}
              title={shown ? "Masquer cet élément" : "Afficher cet élément"}
              className={clsx(
                "inline-flex items-center gap-1.5 rounded-full border px-2.5 py-1 text-xs font-medium",
                shown
                  ? "border-gray-300 bg-white text-gray-800"
                  : "border-gray-200 bg-gray-50 text-gray-400 line-through",
              )}
            >
              <i
                className="h-2.5 w-2.5 rounded-full"
                style={{ backgroundColor: shown ? solidColor(hues.get(item.id) ?? 0) : "#d1d5db" }}
              />
              {item.name}
            </button>
          );
        })}
        {hidden.size > 0 && (
          <button
            type="button"
            onClick={() => setHidden(new Set())}
            className="px-2 text-xs text-blue-600 hover:underline"
          >
            Tout afficher
          </button>
        )}
      </fieldset>

      <div className="max-h-[70vh] overflow-auto rounded-lg bg-white shadow">
        {isLoading ? (
          <SkeletonRows rows={10} columns={7} />
        ) : isError ? (
          <p className="p-8 text-center text-sm text-red-600">Impossible de charger le planning.</p>
        ) : (
          <table className="w-full min-w-[60rem] table-fixed border-separate border-spacing-0 text-xs">
            <thead>
              <tr>
                <th className="sticky top-0 left-0 z-20 w-16 border-b border-gray-200 bg-gray-50 px-2 py-2 text-right font-medium text-gray-500">
                  Heure
                </th>
                {days.map((day) => (
                  <th
                    key={day}
                    scope="col"
                    className={clsx(
                      "sticky top-0 z-10 border-b border-l border-gray-200 px-2 py-2 text-center font-medium",
                      day === today ? "bg-blue-50 text-blue-700" : "bg-gray-50 text-gray-600",
                    )}
                  >
                    <span className="capitalize">{weekday.format(toUtcDate(day))}</span>
                    <span className="block font-normal text-gray-500">
                      {perDay.get(day) ?? 0} réservation{(perDay.get(day) ?? 0) > 1 ? "s" : ""}
                    </span>
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
                      "sticky left-0 z-10 border-b border-gray-100 bg-white px-2 py-1 text-right align-top font-normal",
                      hour === nowHour ? "font-semibold text-blue-700" : "text-gray-400",
                    )}
                  >
                    {hourLabel(hour)}
                  </th>
                  {days.map((day) => (
                    <td
                      key={day}
                      className={clsx(
                        "h-11 border-b border-l border-gray-100 p-0.5 align-top",
                        day === today && "bg-blue-50/40",
                      )}
                    >
                      <div className="flex flex-col gap-0.5">
                        {(cells.get(`${day}|${hour}`) ?? []).map(({ run, first }) => {
                          const head = run.reservations[0];
                          const hue = hues.get(head.item_id) ?? 0;
                          return first ? (
                            <ReservationChip
                              key={head.key}
                              run={run.reservations}
                              hue={hue}
                              onSelect={(value) => setTarget(targetOfRun(value))}
                            />
                          ) : (
                            <i
                              key={head.key}
                              aria-hidden
                              title={`${head.item_name} (suite)`}
                              className="block h-2 rounded"
                              style={{ backgroundColor: solidColor(hue), opacity: 0.45 }}
                            />
                          );
                        })}
                      </div>
                    </td>
                  ))}
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>

      <CancelReservationModal target={target} onClose={() => setTarget(null)} />
    </div>
  );
}
