"use client";

import { useMemo, useState } from "react";
import type { AdminItem } from "@/lib/reservationTypes";
import { DayPlanning } from "./DayPlanning";
import { LoanItems } from "./LoanItems";
import { itemHues, PHASE_ALPHA, PHASE_LABELS, type Phase } from "./planningColors";
import { WeekPlanning } from "./WeekPlanning";

interface PlanningViewProps {
  items: AdminItem[];
  onOpenItem: (itemId: number) => void;
}

type Mode = "day" | "week";

const MODES: { id: Mode; label: string }[] = [
  { id: "day", label: "Jour" },
  { id: "week", label: "Semaine" },
];

const PHASES: Phase[] = ["past", "ongoing", "future"];

export function PlanningView({ items, onOpenItem }: PlanningViewProps) {
  const [mode, setMode] = useState<Mode>("day");
  const slotItems = useMemo(
    () =>
      items
        .filter((item) => item.slot)
        .sort(
          (a, b) =>
            a.category_path.localeCompare(b.category_path, "fr") ||
            a.name.localeCompare(b.name, "fr"),
        ),
    [items],
  );
  const hues = useMemo(() => itemHues(items), [items]);

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center gap-x-6 gap-y-2">
        <div
          className="inline-flex rounded-md border border-gray-200 bg-white p-0.5"
          role="tablist"
        >
          {MODES.map(({ id, label }) => (
            <button
              key={id}
              type="button"
              role="tab"
              aria-selected={mode === id}
              onClick={() => setMode(id)}
              className={`rounded px-4 py-1.5 text-sm font-medium ${
                mode === id ? "bg-blue-600 text-white" : "text-gray-600 hover:bg-gray-100"
              }`}
            >
              {label}
            </button>
          ))}
        </div>
        <p className="flex flex-wrap items-center gap-x-4 gap-y-1 text-xs text-gray-500">
          <span>
            Une couleur par élément. Heures de Paris. Cliquez sur une réservation pour l'annuler.
          </span>
          {PHASES.map((phase) => (
            <span key={phase} className="inline-flex items-center gap-1">
              <i
                className="h-3 w-5 rounded border-l-4 border-gray-500"
                style={{ backgroundColor: `rgb(107 114 128 / ${PHASE_ALPHA[phase]})` }}
              />
              {PHASE_LABELS[phase]}
            </span>
          ))}
        </p>
      </div>

      {mode === "day" ? (
        <DayPlanning slotItems={slotItems} hues={hues} onOpenItem={onOpenItem} />
      ) : (
        <WeekPlanning slotItems={slotItems} hues={hues} />
      )}

      <LoanItems items={items} onOpenItem={onOpenItem} />
    </div>
  );
}
