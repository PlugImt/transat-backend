"use client";

import { ChevronLeft, ChevronRight, ExternalLink, Flag } from "lucide-react";
import { useEffect, useState } from "react";
import PhotoReportModal from "@/components/PhotoReportModal";
import { usePhotoReports } from "@/lib/hooks";
import {
  DELIVERY_STATUS_LABELS,
  formatDateTime,
  REASON_LABELS,
  safeHttpUrl,
} from "@/lib/photoReports";
import type { PhotoReportDeliveryStatus, PhotoReportReason } from "@/lib/types";

const PAGE_SIZE = 25;

const DELIVERY_BADGE: Record<PhotoReportDeliveryStatus, string> = {
  sent: "bg-green-100 text-green-800",
  pending: "bg-yellow-100 text-yellow-800",
  failed: "bg-red-100 text-red-800",
  no_recipients: "bg-gray-100 text-gray-800",
};

export default function PhotoReportsPage() {
  const [reason, setReason] = useState<PhotoReportReason | "">("");
  const [delivery, setDelivery] = useState<PhotoReportDeliveryStatus | "">("");
  const [offset, setOffset] = useState(0);
  const [openId, setOpenId] = useState<number | null>(null);

  const { data, isLoading, isError, isFetching } = usePhotoReports({
    reason,
    delivery,
    limit: PAGE_SIZE,
    offset,
  });
  const reports = data?.reports ?? [];
  const total = data?.total ?? 0;

  // Emails link here with ?report=<id>.
  useEffect(() => {
    const id = Number(new URLSearchParams(window.location.search).get("report"));
    if (Number.isInteger(id) && id > 0) setOpenId(id);
  }, []);

  const closeModal = () => {
    setOpenId(null);
    if (window.location.search) window.history.replaceState(null, "", window.location.pathname);
  };

  const changeFilter =
    <T,>(setter: (value: T) => void) =>
    (value: T) => {
      setter(value);
      setOffset(0);
    };

  const last = Math.min(offset + PAGE_SIZE, total);

  return (
    <div className="p-6">
      <div className="flex flex-wrap items-center justify-between gap-4 mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Signalements de photos</h1>
          <p className="text-sm text-gray-500 mt-1">
            Un signalement est une demande de vérification : la photo n'est ni masquée ni supprimée
            automatiquement.
          </p>
        </div>
        <div className="flex flex-wrap gap-3">
          <select
            value={reason}
            onChange={(e) => changeFilter(setReason)(e.target.value as PhotoReportReason | "")}
            className="rounded-md border border-gray-300 bg-white px-3 py-2 text-sm"
            aria-label="Filtrer par motif"
          >
            <option value="">Tous les motifs</option>
            {Object.entries(REASON_LABELS).map(([value, label]) => (
              <option key={value} value={value}>
                {label}
              </option>
            ))}
          </select>
          <select
            value={delivery}
            onChange={(e) =>
              changeFilter(setDelivery)(e.target.value as PhotoReportDeliveryStatus | "")
            }
            className="rounded-md border border-gray-300 bg-white px-3 py-2 text-sm"
            aria-label="Filtrer par envoi des emails"
          >
            <option value="">Tous les envois</option>
            {Object.entries(DELIVERY_STATUS_LABELS).map(([value, label]) => (
              <option key={value} value={value}>
                {label}
              </option>
            ))}
          </select>
        </div>
      </div>

      <div className="bg-white shadow rounded-lg overflow-hidden">
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                {["Date", "Photo", "Événement / club", "Motif", "Signalé par", "Emails", ""].map(
                  (title) => (
                    <th
                      key={title}
                      className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                    >
                      {title}
                    </th>
                  ),
                )}
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {reports.map((report) => {
                const photoUrl = safeHttpUrl(report.photo_url);
                return (
                  <tr key={report.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {formatDateTime(report.created_at)}
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-900 max-w-xs">
                      <div className="truncate" title={report.photo_name}>
                        {report.photo_name}
                      </div>
                      {!report.photo_available && (
                        <div className="text-xs text-gray-500">Plus dans l'index Drive</div>
                      )}
                      {photoUrl && (
                        <a
                          href={photoUrl}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="inline-flex items-center text-xs text-blue-600 hover:text-blue-800"
                        >
                          <ExternalLink className="h-3 w-3 mr-1" />
                          Ouvrir dans Drive
                        </a>
                      )}
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-900 max-w-xs">
                      <div className="truncate">{report.album_title || "—"}</div>
                      <div className="text-xs text-gray-500">
                        {report.club_name || "Aucun club associé"}
                      </div>
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-900">
                      {REASON_LABELS[report.reason]}
                      {report.comment && (
                        <div className="text-xs text-gray-500 max-w-xs truncate">
                          {report.comment}
                        </div>
                      )}
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-900">
                      <div>{report.reporter.name}</div>
                      <div className="text-xs text-gray-500">
                        {report.reporter.email ?? "Compte supprimé"}
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm">
                      <span
                        className={`inline-flex rounded-full px-2 py-1 text-xs font-medium ${DELIVERY_BADGE[report.delivery_status]}`}
                      >
                        {DELIVERY_STATUS_LABELS[report.delivery_status]}
                      </span>
                      <div className="text-xs text-gray-500 mt-1">
                        {report.recipients_sent}/{report.recipients} destinataires
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-right text-sm">
                      <button
                        type="button"
                        onClick={() => setOpenId(report.id)}
                        className="text-blue-600 hover:text-blue-800 font-medium"
                      >
                        Détails
                      </button>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>

        {isLoading && <div className="py-12 text-center text-sm text-gray-500">Chargement…</div>}
        {isError && (
          <div className="py-12 text-center text-sm text-red-600">
            Impossible de charger les signalements.
          </div>
        )}
        {!isLoading && !isError && reports.length === 0 && (
          <div className="text-center py-12">
            <Flag className="mx-auto h-12 w-12 text-gray-400" />
            <h3 className="mt-2 text-sm font-medium text-gray-900">Aucun signalement</h3>
            <p className="mt-1 text-sm text-gray-500">
              Les photos signalées par les utilisateurs apparaîtront ici.
            </p>
          </div>
        )}

        {total > 0 && (
          <div className="flex items-center justify-between border-t border-gray-200 px-6 py-3 text-sm text-gray-600">
            <span>
              {offset + 1}–{last} sur {total}
              {isFetching ? " …" : ""}
            </span>
            <div className="flex gap-2">
              <button
                type="button"
                disabled={offset === 0}
                onClick={() => setOffset(Math.max(0, offset - PAGE_SIZE))}
                className="rounded-md border border-gray-300 p-1.5 disabled:opacity-40"
                aria-label="Page précédente"
              >
                <ChevronLeft className="h-4 w-4" />
              </button>
              <button
                type="button"
                disabled={last >= total}
                onClick={() => setOffset(offset + PAGE_SIZE)}
                className="rounded-md border border-gray-300 p-1.5 disabled:opacity-40"
                aria-label="Page suivante"
              >
                <ChevronRight className="h-4 w-4" />
              </button>
            </div>
          </div>
        )}
      </div>

      <PhotoReportModal reportId={openId} onClose={closeModal} />
    </div>
  );
}
