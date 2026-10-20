"use client";

import { CalendarDays, ChevronLeft, ChevronRight } from "lucide-react";
import { useMemo, useState } from "react";
import { useAdminReservations } from "@/lib/reservationHooks";
import type { AdminItem, AdminReservation } from "@/lib/reservationTypes";
import { CancelReservationModal, type ReservationTarget } from "./ActionModals";
import { DayPlanning } from "./DayPlanning";
import {
  cellOf,
  hourLabel,
  hoursWith,
  mondayOf,
  shiftDay,
  toUtcDate,
  weekday,
} from "./planningDates";
import { targetOf } from "./ReservationsList";
import {
  EmptyState,
  fieldClass,
  formatTime,
  fullName,
  parisDay,
  parisHour,
  SkeletonRows,
} from "./ui";

interface PlanningViewProps {
  items: AdminItem[];
  onOpenItem: (itemId: number) => void;
}

type Mode = "day" | "week";

const MODES: { id: Mode; label: string }[] = [
  { id: "day", label: "Jour · tous les éléments" },
  { id: "week", label: "Semaine · un élément" },
];

export function PlanningView({ items, onOpenItem }: PlanningViewProps) {
  const [mode, setMode] = useState<Mode>("day");
  return (
    <div className="space-y-4">
      <div className="inline-flex rounded-md border border-gray-200 bg-white p-0.5" role="tablist">
        {MODES.map(({ id, label }) => (
          <button
            key={id}
            type="button"
            role="tab"
            aria-selected={mode === id}
            onClick={() => setMode(id)}
            className={`rounded px-3 py-1.5 text-sm font-medium ${
              mode === id ? "bg-blue-600 text-white" : "text-gray-600 hover:bg-gray-100"
            }`}
          >
            {label}
          </button>
        ))}
      </div>
      {mode === "day" ? (
        <DayPlanning items={items} onOpenItem={onOpenItem} />
      ) : (
        <WeekPlanning items={items} onOpenItem={onOpenItem} />
      )}
    </div>
  );
}

function WeekPlanning({ items, onOpenItem }: PlanningViewProps) {
  const slotItems = useMemo(
    () => items.filter((item) => item.slot).sort((a, b) => a.name.localeCompare(b.name, "fr")),
    [items],
  );
  const today = parisDay(new Date());
  const [itemId, setItemId] = useState<number | undefined>(undefined);
  const [monday, setMonday] = useState(() => mondayOf(today));
  const [target, setTarget] = useState<ReservationTarget | null>(null);

  const selected = slotItems.find((item) => item.id === itemId) ?? slotItems[0];
  const days = Array.from({ length: 7 }, (_, index) => shiftDay(monday, index));

  // Wider than the week: the days are cut in Paris time on this side.
  const { data, isLoading, isError } = useAdminReservations(
    {
      item_id: selected?.id,
      kind: "slot",
      from: shiftDay(monday, -1),
      to: shiftDay(monday, 9),
      order: "asc",
      limit: 200,
    },
    { enabled: !!selected },
  );

  const cells = useMemo(() => {
    const map = new Map<string, AdminReservation>();
    for (const reservation of data?.items ?? []) {
      if (reservation.status === "cancelled") continue;
      const cell = cellOf(reservation.start_date);
      map.set(`${cell.day}|${cell.hour}`, reservation);
    }
    return map;
  }, [data]);

  const hours = useMemo(
    () => hoursWith(data?.items ?? [], (value) => days.includes(value)),
    [data, days],
  );

  if (slotItems.length === 0) {
    return (
      <div className="rounded-lg bg-white shadow">
        <EmptyState
          icon={<CalendarDays className="h-10 w-10" />}
          title="Aucun item à créneaux"
          description="Le planning est disponible pour les salles et le matériel réservés par créneaux."
        />
      </div>
    );
  }

  const nowHour = parisHour(new Date().toISOString());

  return (
    <div className="space-y-4">
      <div className="flex flex-wrap items-center gap-3">
        <select
          value={selected?.id ?? ""}
          onChange={(event) => setItemId(Number(event.target.value))}
          className={`${fieldClass} w-auto min-w-[14rem]`}
          aria-label="Salle ou item"
        >
          {slotItems.map((item) => (
            <option key={item.id} value={item.id}>
              {item.name}
              {item.category_path ? ` — ${item.category_path}` : ""}
            </option>
          ))}
        </select>
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
        <span className="text-sm text-gray-600">
          Semaine du {weekday.format(toUtcDate(days[0]))} au {weekday.format(toUtcDate(days[6]))}
        </span>
        {selected && (
          <button
            type="button"
            onClick={() => onOpenItem(selected.id)}
            className="ml-auto text-sm text-blue-600 hover:underline"
          >
            Détail de {selected.name}
          </button>
        )}
      </div>

      <div className="overflow-x-auto rounded-lg bg-white shadow">
        {isLoading ? (
          <SkeletonRows rows={10} columns={7} />
        ) : isError ? (
          <p className="p-8 text-center text-sm text-red-600">Impossible de charger le planning.</p>
        ) : (
          <table className="min-w-[56rem] table-fixed border-collapse text-xs">
            <thead>
              <tr className="bg-gray-50">
                <th className="w-16 border-b border-gray-200 px-2 py-2 text-right font-medium text-gray-500">
                  Heure
                </th>
                {days.map((day) => (
                  <th
                    key={day}
                    className={`border-b border-l border-gray-200 px-2 py-2 text-center font-medium capitalize ${
                      day === today ? "bg-blue-50 text-blue-700" : "text-gray-600"
                    }`}
                  >
                    {weekday.format(toUtcDate(day))}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {hours.map((hour) => (
                <tr key={hour}>
                  <td className="border-b border-gray-100 px-2 py-1 text-right text-gray-400">
                    {hourLabel(hour)}
                  </td>
                  {days.map((day) => {
                    const reservation = cells.get(`${day}|${hour}`);
                    const past = day < today || (day === today && hour <= nowHour && hour !== 24);
                    return (
                      <td
                        key={day}
                        className={`h-9 border-b border-l border-gray-100 p-0.5 align-middle ${day === today ? "bg-blue-50/40" : ""}`}
                      >
                        {reservation ? (
                          <button
                            type="button"
                            disabled={!reservation.cancellable}
                            onClick={() => setTarget(targetOf(reservation))}
                            title={`${fullName(reservation.first_name, reservation.last_name)} · ${reservation.email}\n${formatTime(reservation.start_date)}${
                              reservation.cancellable ? "\nCliquer pour annuler" : ""
                            }`}
                            className={`h-full w-full truncate rounded px-1.5 text-left font-medium ${
                              reservation.cancellable
                                ? "bg-blue-100 text-blue-900 hover:bg-red-100 hover:text-red-900"
                                : "bg-gray-100 text-gray-500"
                            } ${past ? "opacity-70" : ""}`}
                          >
                            {fullName(reservation.first_name, reservation.last_name) ||
                              reservation.email}
                          </button>
                        ) : null}
                      </td>
                    );
                  })}
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
      <p className="text-xs text-gray-500">
        Heures de Paris. Cliquez sur un créneau réservé pour l'annuler.
      </p>

      <CancelReservationModal target={target} onClose={() => setTarget(null)} />
    </div>
  );
}
