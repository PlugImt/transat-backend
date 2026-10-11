"use client";

import { Copy, ExternalLink, RefreshCw } from "lucide-react";
import toast from "react-hot-toast";
import Modal from "@/components/Modal";
import { useClubs, usePhotoReport, useRetryPhotoReport, useSetAlbumClub } from "@/lib/hooks";
import {
  DELIVERY_STATUS_LABELS,
  formatDateTime,
  REASON_LABELS,
  safeHttpUrl,
} from "@/lib/photoReports";
import type { PhotoReportDelivery } from "@/lib/types";

const DELIVERY_LABELS: Record<PhotoReportDelivery["status"], string> = {
  pending: "En attente",
  sending: "Envoi en cours",
  sent: "Envoyé",
  failed: "Échec",
};

interface PhotoReportModalProps {
  reportId: number | null;
  onClose: () => void;
}

export default function PhotoReportModal({ reportId, onClose }: PhotoReportModalProps) {
  const { data: report, isLoading, isError } = usePhotoReport(reportId);
  const { data: clubs = [] } = useClubs();
  const retry = useRetryPhotoReport();
  const setAlbumClub = useSetAlbumClub();

  const photoUrl = report ? safeHttpUrl(report.photo_url) : null;
  const albumUrl = report ? safeHttpUrl(report.album_url) : null;
  const canRetry = report?.deliveries?.some((d) => d.status === "failed") ?? false;

  const copyLink = async () => {
    if (!photoUrl) return;
    try {
      await navigator.clipboard.writeText(photoUrl);
      toast.success("Lien copié");
    } catch {
      toast.error("Impossible de copier le lien");
    }
  };

  const handleRetry = async () => {
    if (!report) return;
    try {
      const { retried } = await retry.mutateAsync(report.id);
      toast.success(retried > 0 ? "Nouvel envoi lancé" : "Aucun email en échec");
    } catch {
      toast.error("Échec du nouvel envoi");
    }
  };

  const handleClubChange = async (value: string) => {
    if (!report?.album_id) return;
    try {
      await setAlbumClub.mutateAsync({
        albumId: report.album_id,
        clubId: value === "" ? null : Number(value),
      });
      toast.success("Club de l'album mis à jour (valable pour les prochains signalements)");
    } catch {
      toast.error("Impossible de modifier le club de l'album");
    }
  };

  return (
    <Modal
      isOpen={reportId !== null}
      onClose={onClose}
      title={reportId ? `Signalement #${reportId}` : "Signalement"}
      maxWidth="xl"
    >
      {isLoading && <p className="text-sm text-gray-500">Chargement…</p>}
      {isError && <p className="text-sm text-red-600">Signalement introuvable.</p>}

      {report && (
        <div className="space-y-5 text-sm text-gray-900">
          <section className="space-y-2">
            <h4 className="text-xs font-semibold uppercase tracking-wider text-gray-500">Photo</h4>
            <div className="break-words">{report.photo_name}</div>
            <div className="text-xs text-gray-500 break-all">ID Drive : {report.photo_id}</div>
            {!report.photo_available && (
              <p className="rounded-md bg-gray-100 px-3 py-2 text-xs text-gray-700">
                Cette photo n'est plus dans l'index Drive (déjà supprimée ou déplacée).
              </p>
            )}
            {photoUrl && (
              <div className="space-y-2">
                <div className="flex flex-wrap gap-2">
                  <a
                    href={photoUrl}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="inline-flex items-center rounded-md bg-blue-600 px-3 py-2 text-white hover:bg-blue-700"
                  >
                    <ExternalLink className="h-4 w-4 mr-2" />
                    Ouvrir la photo dans Drive
                  </a>
                  <button
                    type="button"
                    onClick={copyLink}
                    className="inline-flex items-center rounded-md border border-gray-300 px-3 py-2 hover:bg-gray-50"
                  >
                    <Copy className="h-4 w-4 mr-2" />
                    Copier le lien
                  </button>
                </div>
                <p className="text-xs text-gray-500 break-all">{photoUrl}</p>
                <p className="text-xs text-gray-500">
                  La suppression se fait dans Google Drive ; la synchronisation retire ensuite la
                  photo de l'application.
                </p>
              </div>
            )}
          </section>

          <section className="space-y-2">
            <h4 className="text-xs font-semibold uppercase tracking-wider text-gray-500">
              Événement et club
            </h4>
            <div>
              {report.album_title || "Aucun événement"}
              {albumUrl && (
                <a
                  href={albumUrl}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="ml-2 inline-flex items-center text-xs text-blue-600 hover:text-blue-800"
                >
                  <ExternalLink className="h-3 w-3 mr-1" />
                  Dossier Drive
                </a>
              )}
            </div>
            <div className="text-xs text-gray-500">
              Club au moment du signalement :{" "}
              {report.club_name || "aucun (administrateurs uniquement)"}
            </div>
            {report.album_id && (
              <label className="block text-xs text-gray-600">
                Club responsable de cet album (destinataire des prochains signalements)
                <select
                  defaultValue={report.club_id ?? ""}
                  key={`${report.id}-${report.club_id ?? ""}`}
                  disabled={setAlbumClub.isPending}
                  onChange={(e) => handleClubChange(e.target.value)}
                  className="mt-1 block w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900"
                >
                  <option value="">Aucun club (administrateurs uniquement)</option>
                  {clubs.map((club) => (
                    <option key={club.id_clubs} value={club.id_clubs}>
                      {club.name}
                    </option>
                  ))}
                </select>
              </label>
            )}
          </section>

          <section className="space-y-2">
            <h4 className="text-xs font-semibold uppercase tracking-wider text-gray-500">
              Signalement
            </h4>
            <div className="font-medium">{REASON_LABELS[report.reason]}</div>
            <div className="whitespace-pre-wrap break-words rounded-md bg-gray-50 px-3 py-2">
              {report.comment || <span className="text-gray-500">Aucun commentaire</span>}
            </div>
            <div className="text-xs text-gray-500">
              Envoyé le {formatDateTime(report.created_at)}
            </div>
          </section>

          <section className="space-y-1">
            <h4 className="text-xs font-semibold uppercase tracking-wider text-gray-500">
              Signalé par
            </h4>
            <div>{report.reporter.name}</div>
            <div className="text-xs text-gray-600">
              {report.reporter.email ?? "Compte supprimé"}
            </div>
            <div className="text-xs text-gray-500">
              ID {report.reporter.id}
              {report.reporter.formation ? ` · ${report.reporter.formation}` : ""}
              {report.reporter.graduation_year ? ` · promo ${report.reporter.graduation_year}` : ""}
            </div>
            {report.reporter.public_id && (
              <div className="text-xs text-gray-500 break-all">
                ID public : {report.reporter.public_id}
              </div>
            )}
          </section>

          <section className="space-y-2">
            <div className="flex items-center justify-between">
              <h4 className="text-xs font-semibold uppercase tracking-wider text-gray-500">
                Emails — {DELIVERY_STATUS_LABELS[report.delivery_status]}
              </h4>
              {canRetry && (
                <button
                  type="button"
                  onClick={handleRetry}
                  disabled={retry.isPending}
                  className="inline-flex items-center rounded-md border border-gray-300 px-3 py-1.5 text-xs hover:bg-gray-50 disabled:opacity-50"
                >
                  <RefreshCw className="h-3 w-3 mr-1" />
                  Renvoyer les échecs
                </button>
              )}
            </div>
            {report.deliveries?.length === 0 && (
              <p className="text-xs text-gray-500">Aucun destinataire n'a été trouvé.</p>
            )}
            <ul className="divide-y divide-gray-100 rounded-md border border-gray-200">
              {report.deliveries?.map((d) => (
                <li key={d.email} className="px-3 py-2">
                  <div className="flex items-center justify-between gap-3">
                    <span className="break-all">{d.email}</span>
                    <span className="whitespace-nowrap text-xs text-gray-600">
                      {d.role === "club_owner" ? "Responsable du club" : "Admin"} ·{" "}
                      {DELIVERY_LABELS[d.status]}
                      {d.attempts > 0 ? ` (${d.attempts} essai${d.attempts > 1 ? "s" : ""})` : ""}
                    </span>
                  </div>
                  {d.last_error && d.status !== "sent" && (
                    <div className="mt-1 break-words text-xs text-red-600">{d.last_error}</div>
                  )}
                </li>
              ))}
            </ul>
          </section>
        </div>
      )}
    </Modal>
  );
}
