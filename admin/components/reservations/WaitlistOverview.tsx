"use client";

import { BellRing, Hourglass, Search, Trash2 } from "lucide-react";
import { useMemo, useState } from "react";
import toast from "react-hot-toast";
import Modal from "@/components/Modal";
import { useAllWaitlists, useRemoveWaitlistEntry } from "@/lib/reservationHooks";
import type { AdminQueuedEntry } from "@/lib/reservationTypes";
import {
  apiErrorMessage,
  Badge,
  EmptyState,
  formatDateTime,
  fullName,
  LoadingButton,
  SkeletonRows,
} from "./ui";

/** Everyone in a waiting list, all items and slots together, with the people currently being waited for first. */
export function WaitlistOverview({ onOpenItem }: { onOpenItem: (id: number) => void }) {
  const { data, isLoading, isError, error } = useAllWaitlists();
  const remove = useRemoveWaitlistEntry();
  const [query, setQuery] = useState("");
  const [removing, setRemoving] = useState<AdminQueuedEntry | null>(null);

  const entries = useMemo(() => {
    const q = query.trim().toLowerCase();
    const all = data ?? [];
    if (!q) return all;
    return all.filter((e) =>
      [e.first_name, e.last_name, e.email, e.item_name, e.club_name].some((v) =>
        v?.toLowerCase().includes(q),
      ),
    );
  }, [data, query]);

  const awaiting = (data ?? []).filter((e) => e.notified_at).length;
  const people = new Set((data ?? []).map((e) => e.email)).size;

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

  return (
    <div className="space-y-4">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div className="flex flex-wrap gap-2 text-sm">
          <span className="rounded-full bg-white px-3 py-1 text-gray-700 ring-1 ring-gray-200">
            <strong>{data?.length ?? 0}</strong> place{(data?.length ?? 0) > 1 ? "s" : ""} en
            attente
          </span>
          <span className="rounded-full bg-white px-3 py-1 text-gray-700 ring-1 ring-gray-200">
            <strong>{people}</strong> personne{people > 1 ? "s" : ""}
          </span>
          {awaiting > 0 && (
            <span className="rounded-full bg-amber-100 px-3 py-1 text-amber-800">
              <strong>{awaiting}</strong> réponse{awaiting > 1 ? "s" : ""} attendue
              {awaiting > 1 ? "s" : ""}
            </span>
          )}
        </div>
        <label className="relative">
          <Search className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
          <input
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Personne, objet ou club"
            className="w-64 rounded-xl border border-gray-300 bg-white py-2 pl-9 pr-3 text-sm"
          />
        </label>
      </div>

      <div className="overflow-hidden rounded-2xl border border-gray-200 bg-white">
        {isLoading && <SkeletonRows rows={5} columns={5} />}
        {isError && (
          <p className="px-4 py-6 text-sm text-red-600">
            {apiErrorMessage(error, "Files d'attente indisponibles")}
          </p>
        )}
        {!isLoading && !isError && entries.length === 0 && (
          <EmptyState
            icon={<Hourglass className="h-10 w-10" />}
            title={query ? "Aucun résultat" : "Personne n'attend"}
            description={
              query ? undefined : "Les personnes inscrites sur une file d'attente apparaîtront ici."
            }
          />
        )}
        {entries.length > 0 && (
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200 text-sm">
              <thead className="bg-gray-50 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
                <tr>
                  <th className="px-4 py-3">Personne</th>
                  <th className="px-4 py-3">Attend</th>
                  <th className="px-4 py-3">Place</th>
                  <th className="px-4 py-3">Inscrit le</th>
                  <th className="px-4 py-3">État</th>
                  <th className="px-4 py-3" />
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-100">
                {entries.map((e) => (
                  <tr key={e.id} className="transition-colors hover:bg-gray-50">
                    <td className="px-4 py-3">
                      <p className="font-medium text-gray-900">
                        {fullName(e.first_name, e.last_name) || e.email}
                      </p>
                      <p className="text-xs text-gray-500">{e.email}</p>
                    </td>
                    <td className="min-w-[14rem] px-4 py-3">
                      <button
                        type="button"
                        onClick={() => onOpenItem(e.item_id)}
                        className="text-left font-medium text-blue-600 hover:text-blue-800"
                      >
                        {e.item_name}
                      </button>
                      <p className="text-xs text-gray-500">
                        {e.club_name}
                        {e.slot_start
                          ? ` · créneau du ${formatDateTime(e.slot_start)}`
                          : e.slot
                            ? ""
                            : " · emprunt"}
                      </p>
                    </td>
                    <td className="whitespace-nowrap px-4 py-3 text-gray-900">n° {e.position}</td>
                    <td className="px-4 py-3 whitespace-nowrap text-gray-600">
                      {formatDateTime(e.joined_at)}
                    </td>
                    <td className="whitespace-nowrap px-4 py-3">
                      {e.notified_at ? (
                        <Badge
                          tone="amber"
                          title={
                            e.turn_expires_at
                              ? `Tour perdu le ${formatDateTime(e.turn_expires_at)}`
                              : undefined
                          }
                        >
                          <BellRing className="h-3 w-3" />
                          Prévenu, réponse attendue
                        </Badge>
                      ) : (
                        <Badge tone="gray">
                          <Hourglass className="h-3 w-3" />
                          En attente
                        </Badge>
                      )}
                    </td>
                    <td className="px-4 py-3 text-right">
                      <button
                        type="button"
                        onClick={() => setRemoving(e)}
                        className="rounded-lg p-2 text-gray-400 transition-colors hover:bg-red-50 hover:text-red-600"
                        aria-label={`Retirer ${fullName(e.first_name, e.last_name)} de la file`}
                      >
                        <Trash2 className="h-4 w-4" />
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      <Modal
        isOpen={removing !== null}
        onClose={() => setRemoving(null)}
        title="Retirer de la file d'attente"
      >
        {removing && (
          <div className="space-y-4 text-sm text-gray-700">
            <p>
              Retirer{" "}
              <strong>{fullName(removing.first_name, removing.last_name) || removing.email}</strong>{" "}
              de la file de <strong>{removing.item_name}</strong> ?
              {removing.notified_at &&
                " Cette personne était prévenue : la suivante le sera à sa place."}
            </p>
            <div className="flex justify-end gap-2">
              <button
                type="button"
                onClick={() => setRemoving(null)}
                className="rounded-xl border border-gray-300 px-4 py-2 hover:bg-gray-50"
              >
                Annuler
              </button>
              <LoadingButton
                loading={remove.isPending}
                onClick={confirmRemove}
                className="rounded-xl bg-red-600 text-white hover:bg-red-700"
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
