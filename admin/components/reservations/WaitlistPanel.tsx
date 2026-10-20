"use client";

import { BellRing, Hourglass, SkipForward, Trash2 } from "lucide-react";
import { useState } from "react";
import toast from "react-hot-toast";
import Modal from "@/components/Modal";
import {
  useAdminWaitlist,
  useAdvanceWaitlist,
  useRemoveWaitlistEntry,
} from "@/lib/reservationHooks";
import type { AdminItem, AdminWaitlistEntry } from "@/lib/reservationTypes";
import {
  apiErrorMessage,
  Badge,
  EmptyState,
  formatDateTime,
  fullName,
  LoadingButton,
  SkeletonRows,
} from "./ui";

interface QueueGroupProps {
  item: AdminItem;
  /** ISO start of the slot; none for an item without slot. */
  start?: string;
}

function QueueGroup({ item, start }: QueueGroupProps) {
  const { data, isLoading, isError, error } = useAdminWaitlist(item.id, start);
  const remove = useRemoveWaitlistEntry();
  const advance = useAdvanceWaitlist();
  const [removing, setRemoving] = useState<AdminWaitlistEntry | null>(null);

  if (isLoading) return <SkeletonRows rows={2} columns={4} />;
  if (isError || !data) {
    return (
      <p className="px-4 py-3 text-sm text-red-600">
        {apiErrorMessage(error, "File d'attente indisponible")}
      </p>
    );
  }
  if (data.entries.length === 0) return null;

  const holder = data.entries.find((entry) => entry.notified_at);

  const confirmRemove = async () => {
    if (!removing) return;
    try {
      const result = await remove.mutateAsync(removing.id);
      toast.success(
        result.next_notified
          ? `${fullName(removing.first_name, removing.last_name)} a été retiré(e) ; la personne suivante est prévenue.`
          : `${fullName(removing.first_name, removing.last_name)} a été retiré(e) de la file.`,
      );
      setRemoving(null);
    } catch (e) {
      toast.error(apiErrorMessage(e, "Impossible de retirer cette personne"));
    }
  };

  const passTurn = async () => {
    try {
      const result = await advance.mutateAsync({ itemId: item.id, startDate: start });
      toast.success(
        result.next_notified
          ? "La personne suivante a été prévenue."
          : "Plus personne à prévenir dans cette file.",
      );
    } catch (e) {
      toast.error(apiErrorMessage(e, "Impossible de passer au suivant"));
    }
  };

  return (
    <div className="rounded-lg border border-gray-200">
      <div className="flex flex-wrap items-center justify-between gap-2 border-b border-gray-200 bg-gray-50 px-4 py-2">
        <div className="text-sm font-medium text-gray-800">
          {start ? formatDateTime(start) : "File de l'élément"}
          <span className="ml-2 font-normal text-gray-500">
            {data.entries.length} personne{data.entries.length > 1 ? "s" : ""}
          </span>
        </div>
        <LoadingButton
          size="sm"
          onClick={passTurn}
          loading={advance.isPending}
          disabled={!data.free}
          title={
            data.free
              ? "Retire la personne qui n'a pas répondu et prévient la suivante"
              : "Disponible uniquement quand l'élément ou le créneau est libre"
          }
          className="border border-gray-300 bg-white text-gray-700 hover:bg-gray-50"
        >
          <SkipForward className="h-4 w-4" />
          Passer au suivant
        </LoadingButton>
      </div>

      {holder ? (
        <div className="flex items-start gap-2 border-b border-amber-200 bg-amber-50 px-4 py-2 text-sm text-amber-900">
          <Hourglass className="mt-0.5 h-4 w-4 shrink-0" />
          <span>
            En attente de la réponse de{" "}
            <strong>{fullName(holder.first_name, holder.last_name)}</strong> : les boutons{" "}
            <em>Réserver</em> / <em>Annuler</em> lui ont été envoyés le{" "}
            {formatDateTime(holder.notified_at as string)}.
            {holder.turn_expires_at && (
              <>
                {" "}
                Sans réponse avant le {formatDateTime(holder.turn_expires_at)}, la personne suivante
                sera prévenue automatiquement.
              </>
            )}
          </span>
        </div>
      ) : data.free ? (
        <div className="border-b border-blue-200 bg-blue-50 px-4 py-2 text-sm text-blue-900">
          Libre, mais personne n'a pu être prévenu (aucun appareil joignable ou notifications de
          réservation désactivées). Les personnes de la file peuvent encore réserver directement.
        </div>
      ) : null}

      <table className="min-w-full divide-y divide-gray-100 text-sm">
        <thead className="text-left text-xs uppercase text-gray-500">
          <tr>
            <th className="w-12 px-4 py-2">#</th>
            <th className="px-4 py-2">Utilisateur</th>
            <th className="px-4 py-2">Ajouté le</th>
            <th className="px-4 py-2">Notification</th>
            <th className="w-12 px-4 py-2" />
          </tr>
        </thead>
        <tbody className="divide-y divide-gray-100">
          {data.entries.map((entry) => (
            <tr key={entry.id}>
              <td className="px-4 py-2 font-semibold text-gray-700">{entry.position}</td>
              <td className="px-4 py-2">
                <div className="font-medium text-gray-900">
                  {fullName(entry.first_name, entry.last_name)}
                </div>
                <div className="text-xs text-gray-500">{entry.email}</div>
              </td>
              <td className="px-4 py-2 text-gray-600">{formatDateTime(entry.joined_at)}</td>
              <td className="px-4 py-2">
                {entry.notified_at ? (
                  <Badge tone="amber">
                    <BellRing className="h-3 w-3" />
                    Prévenu(e) · en attente
                  </Badge>
                ) : (
                  <Badge tone="gray">Pas encore prévenu(e)</Badge>
                )}
              </td>
              <td className="px-4 py-2 text-right">
                <button
                  type="button"
                  onClick={() => setRemoving(entry)}
                  className="rounded p-1 text-red-600 hover:bg-red-50"
                  title="Retirer de la file"
                  aria-label={`Retirer ${fullName(entry.first_name, entry.last_name)} de la file`}
                >
                  <Trash2 className="h-4 w-4" />
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      <Modal
        isOpen={removing !== null}
        onClose={() => setRemoving(null)}
        title="Retirer de la file d'attente"
      >
        {removing && (
          <div className="space-y-4 text-sm text-gray-700">
            <p>
              Retirer <strong>{fullName(removing.first_name, removing.last_name)}</strong> de la
              file d'attente de <strong>{item.name}</strong> ? Cette personne ne sera pas prévenue.
            </p>
            {removing.notified_at && (
              <p className="rounded-md bg-amber-50 p-2 text-amber-900">
                Elle attend actuellement le résultat d'une libération : la personne suivante sera
                prévenue.
              </p>
            )}
            <div className="flex justify-end gap-2">
              <LoadingButton
                onClick={() => setRemoving(null)}
                className="border border-gray-300 text-gray-700 hover:bg-gray-50"
              >
                Annuler
              </LoadingButton>
              <LoadingButton
                onClick={confirmRemove}
                loading={remove.isPending}
                className="bg-red-600 text-white hover:bg-red-700"
              >
                Retirer
              </LoadingButton>
            </div>
          </div>
        )}
      </Modal>
    </div>
  );
}

/** The queues of interest of an item: one for an item without slot, one per slot for an item with slots. */
export function WaitlistPanel({ item }: { item: AdminItem }) {
  const overview = useAdminWaitlist(item.id, undefined, { enabled: item.slot });

  if (!item.slot) {
    return (
      <div className="space-y-3">
        <QueueGroup item={item} />
        {item.waitlist_count === 0 && (
          <EmptyState
            icon={<BellRing className="h-8 w-8" />}
            title="Aucune file d'attente"
            description="Personne n'attend que cet élément se libère."
          />
        )}
      </div>
    );
  }

  if (overview.isLoading) return <SkeletonRows rows={3} columns={3} />;
  const starts = [
    ...new Set((overview.data?.entries ?? []).map((entry) => entry.slot_start).filter(Boolean)),
  ] as string[];
  if (starts.length === 0) {
    return (
      <EmptyState
        icon={<BellRing className="h-8 w-8" />}
        title="Aucune file d'attente"
        description="Personne n'attend qu'un créneau de cette salle se libère."
      />
    );
  }
  return (
    <div className="space-y-3">
      {starts.map((start) => (
        <QueueGroup key={start} item={item} start={start} />
      ))}
    </div>
  );
}
