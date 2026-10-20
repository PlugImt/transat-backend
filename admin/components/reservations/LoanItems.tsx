"use client";

import { useMemo } from "react";
import type { AdminItem } from "@/lib/reservationTypes";
import {
  AvailabilityBadge,
  formatDateTime,
  formatDuration,
  fullName,
  ReservationStatusBadge,
} from "./ui";

const overdueLevel = (seconds: number) =>
  seconds >= 7 * 86_400 ? 7 : seconds >= 5 * 86_400 ? 5 : seconds >= 86_400 ? 1 : 0;

interface LoanItemsProps {
  items: AdminItem[];
  onOpenItem: (itemId: number) => void;
}

/** The items lent without a slot, latest borrowers first: they have no place in a calendar. */
export function LoanItems({ items, onOpenItem }: LoanItemsProps) {
  const loanItems = useMemo(
    () =>
      items
        .filter((item) => !item.slot)
        .sort(
          (a, b) =>
            (b.current?.elapsed_seconds ?? -1) - (a.current?.elapsed_seconds ?? -1) ||
            a.name.localeCompare(b.name, "fr"),
        ),
    [items],
  );

  return (
    <section className="space-y-3">
      <h3 className="text-sm font-semibold text-gray-800">
        Éléments prêtés sans créneau
        <span className="ml-2 font-normal text-gray-500">{loanItems.length}</span>
      </h3>
      <div className="overflow-x-auto rounded-lg bg-white shadow">
        {loanItems.length === 0 ? (
          <p className="p-6 text-center text-sm text-gray-500">Aucun élément prêté sans créneau.</p>
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
  );
}
