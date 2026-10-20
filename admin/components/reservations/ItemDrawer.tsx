"use client";

import { Dialog, DialogPanel, Transition, TransitionChild } from "@headlessui/react";
import { CalendarClock, CheckCircle2, MapPin, Package, Pencil, X, XCircle } from "lucide-react";
import { Fragment, useState } from "react";
import { useAdminReservations } from "@/lib/reservationHooks";
import type { AdminItem, AdminReservation } from "@/lib/reservationTypes";
import { CancelReservationModal, MarkReturnedModal, type ReservationTarget } from "./ActionModals";
import {
  AvailabilityBadge,
  Badge,
  CANCEL_KIND_LABELS,
  EmptyState,
  formatDateTime,
  formatDuration,
  fullName,
  ITEM_STATUS_LABELS,
  LoadingButton,
  ReminderBadges,
  ReservationStatusBadge,
  SkeletonRows,
} from "./ui";
import { WaitlistPanel } from "./WaitlistPanel";

interface ItemDrawerProps {
  item: AdminItem | null;
  onClose: () => void;
  onEdit: (item: AdminItem) => void;
}

const targetOf = (item: AdminItem, r: AdminReservation): ReservationTarget => ({
  itemId: item.id,
  itemName: item.name,
  slot: item.slot,
  startDate: r.start_date,
  endDate: r.end_date,
  ownerName: fullName(r.first_name, r.last_name) || r.email,
  ownerEmail: r.email,
});

function Section({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <section className="space-y-3">
      <h3 className="text-sm font-semibold uppercase tracking-wide text-gray-500">{title}</h3>
      {children}
    </section>
  );
}

function ReservationLine({
  item,
  reservation,
  onCancel,
  onReturn,
}: {
  item: AdminItem;
  reservation: AdminReservation;
  onCancel: (target: ReservationTarget) => void;
  onReturn: (target: ReservationTarget) => void;
}) {
  const name = fullName(reservation.first_name, reservation.last_name) || reservation.email;
  return (
    <li className="flex flex-wrap items-center justify-between gap-3 px-4 py-3">
      <div className="min-w-0 space-y-1">
        <div className="flex flex-wrap items-center gap-2">
          <span className="font-medium text-gray-900">{name}</span>
          <ReservationStatusBadge
            status={reservation.status}
            overdueLevel={reservation.overdue_level}
            cancelKind={reservation.cancel_kind}
          />
        </div>
        <div className="text-xs text-gray-500">{reservation.email}</div>
        <div className="text-sm text-gray-700">
          {reservation.slot
            ? `${formatDateTime(reservation.start_date)}${reservation.end_date ? ` → ${formatDateTime(reservation.end_date)}` : ""}`
            : `Depuis le ${formatDateTime(reservation.start_date)}${
                reservation.elapsed_seconds !== undefined
                  ? ` · ${formatDuration(reservation.elapsed_seconds)}`
                  : ""
              }`}
        </div>
        {reservation.cancel_kind && (
          <div className="text-xs text-gray-500">
            {CANCEL_KIND_LABELS[reservation.cancel_kind]}
            {reservation.cancelled_by ? ` (${reservation.cancelled_by})` : ""}
            {reservation.cancelled_at ? ` · ${formatDateTime(reservation.cancelled_at)}` : ""}
          </div>
        )}
        {reservation.cancellable && (
          <div className="flex items-center gap-2 text-xs text-gray-500">
            Rappels : <ReminderBadges reminders={reservation.reminders} />
          </div>
        )}
      </div>
      {reservation.cancellable && (
        <div className="flex gap-2">
          {!item.slot && (
            <LoadingButton
              size="sm"
              onClick={() => onReturn(targetOf(item, reservation))}
              className="border border-green-600 text-green-700 hover:bg-green-50"
            >
              <CheckCircle2 className="h-4 w-4" />
              Rendu
            </LoadingButton>
          )}
          <LoadingButton
            size="sm"
            onClick={() => onCancel(targetOf(item, reservation))}
            className="border border-red-600 text-red-700 hover:bg-red-50"
          >
            <XCircle className="h-4 w-4" />
            Annuler
          </LoadingButton>
        </div>
      )}
    </li>
  );
}

export function ItemDrawer({ item, onClose, onEdit }: ItemDrawerProps) {
  const [cancelTarget, setCancelTarget] = useState<ReservationTarget | null>(null);
  const [returnTarget, setReturnTarget] = useState<ReservationTarget | null>(null);
  const reservations = useAdminReservations(
    { item_id: item?.id, limit: 100, order: "desc" },
    { enabled: !!item },
  );

  const lines = reservations.data?.items ?? [];
  const live = lines
    .filter((r) => r.cancellable)
    .sort((a, b) => a.start_date.localeCompare(b.start_date));
  const history = lines.filter((r) => !r.cancellable);

  return (
    <Transition show={item !== null} as={Fragment}>
      <Dialog as="div" className="relative z-40" onClose={onClose}>
        <TransitionChild
          as={Fragment}
          enter="ease-out duration-200"
          enterFrom="opacity-0"
          enterTo="opacity-100"
          leave="ease-in duration-150"
          leaveFrom="opacity-100"
          leaveTo="opacity-0"
        >
          <div className="fixed inset-0 bg-black/30" />
        </TransitionChild>

        <div className="fixed inset-0 overflow-hidden">
          <div className="absolute inset-y-0 right-0 flex max-w-full">
            <TransitionChild
              as={Fragment}
              enter="transform transition ease-out duration-200"
              enterFrom="translate-x-full"
              enterTo="translate-x-0"
              leave="transform transition ease-in duration-150"
              leaveFrom="translate-x-0"
              leaveTo="translate-x-full"
            >
              <DialogPanel className="flex h-full w-screen max-w-2xl flex-col bg-white shadow-xl">
                {item && (
                  <>
                    <header className="flex items-start justify-between gap-4 border-b border-gray-200 p-5">
                      <div className="min-w-0 space-y-2">
                        <div className="flex items-center gap-2 text-gray-500">
                          {item.slot ? (
                            <CalendarClock className="h-5 w-5" />
                          ) : (
                            <Package className="h-5 w-5" />
                          )}
                          <span className="text-sm">
                            {item.slot ? "Réservation par créneaux" : "Emprunt continu"}
                          </span>
                        </div>
                        <h2 className="truncate text-xl font-bold text-gray-900">{item.name}</h2>
                        <div className="flex flex-wrap items-center gap-2">
                          <AvailabilityBadge value={item.availability} />
                          {item.status !== "active" && (
                            <Badge tone="purple">{ITEM_STATUS_LABELS[item.status]}</Badge>
                          )}
                          {item.waitlist_count > 0 && (
                            <Badge tone="blue">{item.waitlist_count} en file d'attente</Badge>
                          )}
                        </div>
                        <p className="text-sm text-gray-600">
                          {item.club_name}
                          {item.category_path ? ` · ${item.category_path}` : ""}
                        </p>
                        {item.location && (
                          <p className="flex items-center gap-1 text-sm text-gray-600">
                            <MapPin className="h-4 w-4" />
                            {item.location}
                          </p>
                        )}
                        {item.description && (
                          <p className="text-sm text-gray-600">{item.description}</p>
                        )}
                        {item.status !== "active" && item.status_note && (
                          <p className="text-sm text-purple-700">Motif : {item.status_note}</p>
                        )}
                      </div>
                      <div className="flex shrink-0 items-center gap-1">
                        <LoadingButton
                          size="sm"
                          onClick={() => onEdit(item)}
                          className="border border-gray-300 text-gray-700 hover:bg-gray-50"
                        >
                          <Pencil className="h-4 w-4" />
                          Modifier
                        </LoadingButton>
                        <button
                          type="button"
                          onClick={onClose}
                          className="rounded p-1.5 text-gray-400 hover:bg-gray-100 hover:text-gray-600"
                          aria-label="Fermer"
                        >
                          <X className="h-5 w-5" />
                        </button>
                      </div>
                    </header>

                    <div className="flex-1 space-y-8 overflow-y-auto p-5">
                      <Section
                        title={item.slot ? "Réservations en cours et à venir" : "Emprunt en cours"}
                      >
                        {reservations.isLoading ? (
                          <SkeletonRows rows={3} columns={3} />
                        ) : live.length === 0 ? (
                          <EmptyState
                            icon={
                              item.slot ? (
                                <CalendarClock className="h-8 w-8" />
                              ) : (
                                <Package className="h-8 w-8" />
                              )
                            }
                            title={item.slot ? "Aucune réservation à venir" : "Disponible"}
                            description={
                              item.slot
                                ? "Aucun créneau n'est réservé pour le moment."
                                : "Personne n'a cet objet actuellement."
                            }
                          />
                        ) : (
                          <ul className="divide-y divide-gray-100 rounded-lg border border-gray-200">
                            {live.map((r) => (
                              <ReservationLine
                                key={r.key}
                                item={item}
                                reservation={r}
                                onCancel={setCancelTarget}
                                onReturn={setReturnTarget}
                              />
                            ))}
                          </ul>
                        )}
                      </Section>

                      <Section title="File d'attente">
                        <WaitlistPanel item={item} />
                      </Section>

                      <Section title="Historique récent">
                        {reservations.isLoading ? (
                          <SkeletonRows rows={2} columns={3} />
                        ) : history.length === 0 ? (
                          <p className="text-sm text-gray-500">Aucun historique.</p>
                        ) : (
                          <ul className="divide-y divide-gray-100 rounded-lg border border-gray-200">
                            {history.slice(0, 15).map((r) => (
                              <ReservationLine
                                key={r.key}
                                item={item}
                                reservation={r}
                                onCancel={setCancelTarget}
                                onReturn={setReturnTarget}
                              />
                            ))}
                          </ul>
                        )}
                      </Section>
                    </div>
                    <CancelReservationModal
                      target={cancelTarget}
                      onClose={() => setCancelTarget(null)}
                    />
                    <MarkReturnedModal
                      target={returnTarget}
                      onClose={() => setReturnTarget(null)}
                    />
                  </>
                )}
              </DialogPanel>
            </TransitionChild>
          </div>
        </div>
      </Dialog>
    </Transition>
  );
}
