"use client";

import {
  CalendarX2,
  CheckCircle2,
  ChevronLeft,
  ChevronRight,
  ListFilter,
  Search,
  XCircle,
} from "lucide-react";
import { useEffect, useState } from "react";
import { useDebounce } from "@/lib/hooks";
import { useAdminReservations } from "@/lib/reservationHooks";
import type {
  AdminItem,
  AdminReservation,
  ReservationFilters,
  ReservationStatusFilter,
} from "@/lib/reservationTypes";
import { CancelReservationModal, MarkReturnedModal, type ReservationTarget } from "./ActionModals";
import {
  apiErrorMessage,
  CANCEL_KIND_LABELS,
  EmptyState,
  fieldClass,
  formatDateTime,
  formatDuration,
  fullName,
  LoadingButton,
  ReminderBadges,
  ReservationStatusBadge,
  SkeletonRows,
} from "./ui";

const PAGE_SIZE = 25;

const STATUS_TABS: { value: ReservationStatusFilter; label: string }[] = [
  { value: "", label: "Toutes" },
  { value: "active", label: "En cours" },
  { value: "upcoming", label: "À venir" },
  { value: "past", label: "Terminées" },
  { value: "cancelled", label: "Annulées" },
  { value: "overdue", label: "En retard" },
];

const OVERDUE_STEPS = [
  { days: 1, label: "+1 j" },
  { days: 5, label: "+5 j" },
  { days: 7, label: "+7 j" },
];

/** A day typed in the panel, as the instant the day starts for the person using it. */
const dayStart = (day: string, plusDays = 0) => {
  const date = new Date(`${day}T00:00:00`);
  date.setDate(date.getDate() + plusDays);
  return date.toISOString();
};

export const targetOf = (r: AdminReservation): ReservationTarget => ({
  itemId: r.item_id,
  itemName: r.item_name,
  slot: r.slot,
  startDate: r.start_date,
  endDate: r.end_date,
  ownerName: fullName(r.first_name, r.last_name) || r.email,
  ownerEmail: r.email,
});

/** One cancel target for a run of consecutive slots. */
export const targetOfRun = (run: AdminReservation[]): ReservationTarget => ({
  ...targetOf(run[0]),
  endDate: run[run.length - 1].end_date,
  startDates: run.map((r) => r.start_date),
});

interface ReservationsListProps {
  filters: ReservationFilters;
  onFiltersChange: (filters: ReservationFilters) => void;
  items: AdminItem[];
  onOpenItem: (itemId: number) => void;
}

export function ReservationsList({
  filters,
  onFiltersChange,
  items,
  onOpenItem,
}: ReservationsListProps) {
  const [page, setPage] = useState(0);
  const [userQuery, setUserQuery] = useState(filters.q ?? "");
  const [fromDay, setFromDay] = useState("");
  const [toDay, setToDay] = useState("");
  const [cancelTarget, setCancelTarget] = useState<ReservationTarget | null>(null);
  const [returnTarget, setReturnTarget] = useState<ReservationTarget | null>(null);

  const debouncedQuery = useDebounce(userQuery, 300);
  const status = filters.status ?? "";

  const query: ReservationFilters = {
    ...filters,
    q: debouncedQuery.trim() || undefined,
    from: fromDay ? dayStart(fromDay) : undefined,
    to: toDay ? dayStart(toDay, 1) : undefined,
    order: filters.order ?? "desc",
    limit: PAGE_SIZE,
    offset: page * PAGE_SIZE,
  };

  const { data, isLoading, isFetching, isError, error, refetch } = useAdminReservations(query);

  // A new filter starts again from the first page.
  // biome-ignore lint/correctness/useExhaustiveDependencies: only the filter values matter.
  useEffect(
    () => setPage(0),
    [
      filters.status,
      filters.kind,
      filters.item_id,
      filters.overdue_days,
      debouncedQuery,
      fromDay,
      toDay,
      filters.order,
    ],
  );

  const set = (patch: Partial<ReservationFilters>) => onFiltersChange({ ...filters, ...patch });
  const total = data?.total ?? 0;
  const lastPage = Math.max(0, Math.ceil(total / PAGE_SIZE) - 1);
  const filtersActive = !!(
    status ||
    filters.kind ||
    filters.item_id ||
    filters.overdue_days ||
    userQuery ||
    fromDay ||
    toDay
  );

  const reset = () => {
    onFiltersChange({});
    setUserQuery("");
    setFromDay("");
    setToDay("");
  };

  return (
    <div className="space-y-4">
      <div className="flex flex-wrap items-center gap-2">
        {STATUS_TABS.map((tab) => (
          <button
            key={tab.value}
            type="button"
            onClick={() =>
              set({
                status: tab.value,
                overdue_days: tab.value === "overdue" ? filters.overdue_days : undefined,
              })
            }
            className={`rounded-full px-3 py-1 text-sm font-medium ${
              status === tab.value
                ? "bg-blue-600 text-white"
                : "bg-white text-gray-700 shadow-sm ring-1 ring-gray-200 hover:bg-gray-50"
            }`}
          >
            {tab.label}
          </button>
        ))}
        {(status === "overdue" || status === "active") && (
          <div className="ml-2 flex items-center gap-1 text-sm text-gray-600">
            Emprunts de
            {OVERDUE_STEPS.map((step) => (
              <button
                key={step.days}
                type="button"
                onClick={() =>
                  set({
                    status: "overdue",
                    overdue_days: filters.overdue_days === step.days ? undefined : step.days,
                  })
                }
                className={`rounded-md px-2 py-0.5 text-xs font-medium ${
                  filters.overdue_days === step.days
                    ? "bg-amber-500 text-white"
                    : "bg-amber-50 text-amber-800 hover:bg-amber-100"
                }`}
              >
                {step.label}
              </button>
            ))}
          </div>
        )}
      </div>

      <div className="flex flex-wrap items-center gap-3">
        <div className="relative min-w-[14rem] flex-1">
          <Search className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
          <input
            type="search"
            value={userQuery}
            onChange={(event) => setUserQuery(event.target.value)}
            placeholder="Utilisateur (nom, prénom, e-mail)…"
            className={`${fieldClass} pl-9`}
            aria-label="Filtrer par utilisateur"
          />
        </div>
        <select
          value={filters.item_id ?? ""}
          onChange={(event) =>
            set({ item_id: event.target.value ? Number(event.target.value) : undefined })
          }
          className={`${fieldClass} w-auto max-w-[16rem]`}
          aria-label="Item ou salle"
        >
          <option value="">Tous les items et salles</option>
          {[...items]
            .sort((a, b) => a.name.localeCompare(b.name, "fr"))
            .map((item) => (
              <option key={item.id} value={item.id}>
                {item.name}
              </option>
            ))}
        </select>
        <select
          value={filters.kind ?? ""}
          onChange={(event) => set({ kind: event.target.value as ReservationFilters["kind"] })}
          className={`${fieldClass} w-auto`}
          aria-label="Type"
        >
          <option value="">Créneaux et emprunts</option>
          <option value="slot">Créneaux</option>
          <option value="loan">Emprunts continus</option>
        </select>
        <label className="flex items-center gap-1 text-sm text-gray-600">
          Du
          <input
            type="date"
            value={fromDay}
            onChange={(event) => setFromDay(event.target.value)}
            className={`${fieldClass} w-auto`}
          />
        </label>
        <label className="flex items-center gap-1 text-sm text-gray-600">
          au
          <input
            type="date"
            value={toDay}
            min={fromDay || undefined}
            onChange={(event) => setToDay(event.target.value)}
            className={`${fieldClass} w-auto`}
          />
        </label>
        {filtersActive && (
          <button type="button" onClick={reset} className="text-sm text-blue-600 hover:underline">
            Réinitialiser
          </button>
        )}
      </div>

      <div
        className={`overflow-hidden rounded-lg bg-white shadow ${isFetching && !isLoading ? "opacity-70" : ""}`}
      >
        {isLoading ? (
          <SkeletonRows rows={8} columns={6} />
        ) : isError ? (
          <div className="p-8 text-center text-sm text-red-600">
            {apiErrorMessage(error, "Impossible de charger les réservations")}{" "}
            <button type="button" onClick={() => refetch()} className="underline">
              Réessayer
            </button>
          </div>
        ) : !data || data.items.length === 0 ? (
          <EmptyState
            icon={
              filtersActive ? (
                <ListFilter className="h-10 w-10" />
              ) : (
                <CalendarX2 className="h-10 w-10" />
              )
            }
            title={filtersActive ? "Aucune réservation ne correspond" : "Aucune réservation"}
            description={
              filtersActive
                ? "Essayez d'élargir les filtres."
                : "Les réservations des utilisateurs apparaîtront ici."
            }
          />
        ) : (
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200 text-sm">
              <thead className="bg-gray-50 text-left text-xs font-medium uppercase text-gray-500">
                <tr>
                  <th className="px-4 py-3">Item ou salle</th>
                  <th className="px-4 py-3">Utilisateur</th>
                  <th className="px-4 py-3">
                    <button
                      type="button"
                      onClick={() =>
                        set({ order: (filters.order ?? "desc") === "desc" ? "asc" : "desc" })
                      }
                      className="uppercase"
                    >
                      Période / durée {(filters.order ?? "desc") === "desc" ? "↓" : "↑"}
                    </button>
                  </th>
                  <th className="px-4 py-3">Statut</th>
                  <th className="px-4 py-3">Rappels envoyés</th>
                  <th className="px-4 py-3" />
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-100">
                {data.items.map((r) => (
                  <tr
                    key={r.key}
                    className={
                      r.overdue_level >= 7
                        ? "bg-red-50/50"
                        : r.overdue_level
                          ? "bg-amber-50/50"
                          : undefined
                    }
                  >
                    <td className="px-4 py-3">
                      <button
                        type="button"
                        onClick={() => onOpenItem(r.item_id)}
                        className="text-left font-medium text-gray-900 hover:text-blue-700 hover:underline"
                      >
                        {r.item_name}
                      </button>
                      <div className="text-xs text-gray-500">
                        {r.slot ? "Créneau" : "Emprunt"} · {r.club_name}
                      </div>
                    </td>
                    <td className="px-4 py-3">
                      <div className="font-medium text-gray-900">
                        {fullName(r.first_name, r.last_name) || "—"}
                      </div>
                      <div className="text-xs text-gray-500">{r.email}</div>
                    </td>
                    <td className="px-4 py-3 text-gray-700">
                      {r.slot ? (
                        <>
                          {formatDateTime(r.start_date)}
                          {r.end_date && (
                            <div className="text-xs text-gray-500">
                              → {formatDateTime(r.end_date)}
                            </div>
                          )}
                        </>
                      ) : (
                        <>
                          {formatDateTime(r.start_date)}
                          {r.elapsed_seconds !== undefined && (
                            <div
                              className={`text-xs ${r.overdue_level ? "font-medium text-amber-700" : "text-gray-500"}`}
                            >
                              {r.status === "active" ? "depuis " : "durée "}
                              {formatDuration(r.elapsed_seconds)}
                            </div>
                          )}
                        </>
                      )}
                    </td>
                    <td className="px-4 py-3">
                      <div className="space-y-1">
                        <ReservationStatusBadge
                          status={r.status}
                          overdueLevel={r.overdue_level}
                          cancelKind={r.cancel_kind}
                        />
                        {r.cancel_kind && r.status === "cancelled" && (
                          <div className="text-xs text-gray-500">
                            {CANCEL_KIND_LABELS[r.cancel_kind]}
                            {r.cancelled_at ? ` · ${formatDateTime(r.cancelled_at)}` : ""}
                          </div>
                        )}
                      </div>
                    </td>
                    <td className="px-4 py-3">
                      {r.cancellable ? (
                        <ReminderBadges reminders={r.reminders} />
                      ) : (
                        <span className="text-xs text-gray-300">—</span>
                      )}
                    </td>
                    <td className="px-4 py-3">
                      {r.cancellable && (
                        <div className="flex justify-end gap-2">
                          {!r.slot && (
                            <LoadingButton
                              size="sm"
                              onClick={() => setReturnTarget(targetOf(r))}
                              className="border border-green-600 text-green-700 hover:bg-green-50"
                            >
                              <CheckCircle2 className="h-4 w-4" />
                              Rendu
                            </LoadingButton>
                          )}
                          <LoadingButton
                            size="sm"
                            onClick={() => setCancelTarget(targetOf(r))}
                            className="border border-red-600 text-red-700 hover:bg-red-50"
                          >
                            <XCircle className="h-4 w-4" />
                            Annuler
                          </LoadingButton>
                        </div>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}

        {data && data.items.length > 0 && (
          <div className="flex items-center justify-between border-t border-gray-100 px-4 py-2 text-sm text-gray-600">
            <span>
              {page * PAGE_SIZE + 1}–{page * PAGE_SIZE + data.items.length} sur {total}
            </span>
            <div className="flex items-center gap-1">
              <button
                type="button"
                disabled={page === 0}
                onClick={() => setPage((p) => Math.max(0, p - 1))}
                className="rounded p-1.5 hover:bg-gray-100 disabled:opacity-40"
                aria-label="Page précédente"
              >
                <ChevronLeft className="h-4 w-4" />
              </button>
              <span>
                Page {page + 1} / {lastPage + 1}
              </span>
              <button
                type="button"
                disabled={page >= lastPage}
                onClick={() => setPage((p) => Math.min(lastPage, p + 1))}
                className="rounded p-1.5 hover:bg-gray-100 disabled:opacity-40"
                aria-label="Page suivante"
              >
                <ChevronRight className="h-4 w-4" />
              </button>
            </div>
          </div>
        )}
      </div>

      <CancelReservationModal target={cancelTarget} onClose={() => setCancelTarget(null)} />
      <MarkReturnedModal target={returnTarget} onClose={() => setReturnTarget(null)} />
    </div>
  );
}
