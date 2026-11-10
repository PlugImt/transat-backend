"use client";

import { CircleCheck, CircleX, Copy, ExternalLink, RefreshCw, Trash2 } from "lucide-react";
import { useRef, useState } from "react";
import toast from "react-hot-toast";
import Modal from "@/components/Modal";
import { PhotoPreview } from "@/components/photoReports/PhotoPreview";
import {
  useClubs,
  useDecidePhotoReport,
  usePhotoReport,
  useRetryPhotoReport,
  useSetAlbumClub,
} from "@/lib/hooks";
import {
  DELIVERY_STATUS_LABELS,
  formatDateTime,
  REASON_LABELS,
  STATUS_BADGE,
  safeHttpUrl,
} from "@/lib/photoReports";
import { useAuthStore } from "@/lib/stores/authStore";
import type { PhotoReport, PhotoReportDelivery } from "@/lib/types";

const DELIVERY_LABELS: Record<PhotoReportDelivery["status"], string> = {
  pending: "En attente",
  sending: "Envoi en cours",
  sent: "Envoyé",
  failed: "Échec",
};

const ROLE_LABELS: Record<PhotoReportDelivery["role"], string> = {
  club_owner: "Responsable du club",
  photo_club: "Club Photo",
  admin: "Admin",
};

const DECISION_TITLES = {
  pending: "En attente de décision",
  approved: "Approuvé",
  rejected: "Refusé",
};

interface PhotoReportModalProps {
  reportId: number | null;
  onClose: () => void;
  /** False when the server cannot delete from Drive. */
  canRemove: boolean;
}

export default function PhotoReportModal({ reportId, onClose, canRemove }: PhotoReportModalProps) {
  // The dialog keeps showing the last report while it fades out after closing.
  const lastId = useRef(reportId);
  if (reportId !== null) lastId.current = reportId;
  const shownId = reportId ?? lastId.current;
  const { data: report, isLoading, isError } = usePhotoReport(shownId);
  const isAdmin = useAuthStore((s) => s.access?.admin ?? false);

  return (
    <Modal
      isOpen={reportId !== null}
      onClose={onClose}
      title={shownId ? `Signalement #${shownId}` : "Signalement"}
      maxWidth="4xl"
    >
      {isLoading && <div className="h-80 animate-pulse rounded-2xl bg-gray-100" />}
      {isError && <p className="text-sm text-red-600">Signalement introuvable.</p>}
      {report && <ReportDetail report={report} canRemove={canRemove} isAdmin={isAdmin} />}
    </Modal>
  );
}

function ReportDetail({
  report,
  canRemove,
  isAdmin,
}: {
  report: PhotoReport;
  canRemove: boolean;
  isAdmin: boolean;
}) {
  const photoUrl = safeHttpUrl(report.photo_url);
  const albumUrl = safeHttpUrl(report.album_url);
  const removed = report.status === "approved" || !report.photo_available;

  return (
    <div className="grid gap-6 text-sm text-gray-900 md:grid-cols-[minmax(0,1.1fr)_minmax(0,1fr)]">
      <div className="space-y-3">
        <PhotoPreview
          photoId={report.photo_id}
          alt={report.photo_name}
          width={1200}
          removed={removed}
          className="aspect-[4/3] w-full rounded-2xl"
        />
        <div className="flex flex-wrap items-center gap-2">
          {photoUrl && !removed && (
            <a
              href={photoUrl}
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-1.5 rounded-xl border border-gray-300 px-3 py-1.5 text-xs font-medium transition-colors hover:bg-gray-50"
            >
              <ExternalLink className="h-3.5 w-3.5" />
              Ouvrir dans Drive
            </a>
          )}
          {photoUrl && (
            <button
              type="button"
              onClick={async () => {
                try {
                  await navigator.clipboard.writeText(photoUrl);
                  toast.success("Lien copié");
                } catch {
                  toast.error("Impossible de copier le lien");
                }
              }}
              className="inline-flex items-center gap-1.5 rounded-xl border border-gray-300 px-3 py-1.5 text-xs font-medium transition-colors hover:bg-gray-50"
            >
              <Copy className="h-3.5 w-3.5" />
              Copier le lien
            </button>
          )}
        </div>
        <p className="break-all text-xs text-gray-500">{report.photo_name}</p>
      </div>

      <div className="space-y-5">
        <Decision report={report} canRemove={canRemove} />

        <Section title="Signalement">
          <p className="font-medium">{REASON_LABELS[report.reason]}</p>
          <p className="whitespace-pre-wrap break-words rounded-xl bg-gray-50 px-3 py-2">
            {report.comment || <span className="text-gray-500">Aucun commentaire</span>}
          </p>
          <p className="text-xs text-gray-500">
            Par {report.reporter.name}
            {report.reporter.email ? ` (${report.reporter.email})` : " (compte supprimé)"}
            {report.reporter.formation ? ` · ${report.reporter.formation}` : ""}
            {report.reporter.graduation_year ? ` · promo ${report.reporter.graduation_year}` : ""} —
            le {formatDateTime(report.created_at)}
          </p>
        </Section>

        <Section title="Événement et club">
          <p>
            {report.album_title || "Aucun événement"}
            {albumUrl && (
              <a
                href={albumUrl}
                target="_blank"
                rel="noopener noreferrer"
                className="ml-2 inline-flex items-center text-xs text-blue-600 hover:text-blue-800"
              >
                <ExternalLink className="mr-1 h-3 w-3" />
                Dossier Drive
              </a>
            )}
          </p>
          <p className="text-xs text-gray-500">Club : {report.club_name || "aucun"}</p>
          {isAdmin && report.album_id && <AlbumClubSelect report={report} />}
        </Section>

        <Deliveries report={report} />
      </div>
    </div>
  );
}

function Section({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <section className="space-y-2">
      <h4 className="text-xs font-semibold uppercase tracking-wider text-gray-500">{title}</h4>
      {children}
    </section>
  );
}

function Decision({ report, canRemove }: { report: PhotoReport; canRemove: boolean }) {
  const decide = useDecidePhotoReport();
  const [note, setNote] = useState("");
  const [confirming, setConfirming] = useState(false);

  const run = async (decision: "approve" | "reject") => {
    try {
      await decide.mutateAsync({ id: report.id, decision, note });
      toast.success(
        decision === "approve" ? "Photo supprimée, signalement approuvé" : "Signalement refusé",
      );
      setConfirming(false);
    } catch (e) {
      toast.error((e as Error).message || "La décision n'a pas pu être enregistrée");
    }
  };

  if (report.status !== "pending") {
    return (
      <div className="rounded-2xl border border-gray-200 bg-gray-50 p-4">
        <span
          className={`inline-flex rounded-full px-2.5 py-1 text-xs font-semibold ${STATUS_BADGE[report.status]}`}
        >
          {DECISION_TITLES[report.status]}
        </span>
        <p className="mt-2 text-xs text-gray-600">
          {report.decided_by ? `Par ${report.decided_by}` : ""}
          {report.decided_at ? ` le ${formatDateTime(report.decided_at)}` : ""}
        </p>
        {report.decision_note && (
          <p className="mt-2 whitespace-pre-wrap text-sm">{report.decision_note}</p>
        )}
      </div>
    );
  }

  return (
    <div className="space-y-3 rounded-2xl border border-blue-200 bg-blue-50 p-4">
      <p className="text-sm font-semibold">Décision</p>
      <textarea
        value={note}
        onChange={(e) => setNote(e.target.value)}
        maxLength={1000}
        rows={2}
        placeholder="Note interne (facultative)"
        className="w-full resize-none rounded-xl border border-gray-300 bg-white px-3 py-2 text-sm"
      />
      {confirming ? (
        <div className="space-y-2 rounded-xl border border-red-200 bg-red-50 p-3 animate-fade-in">
          <p className="text-xs text-red-800">
            La photo sera mise à la corbeille du Drive (récupérable 30 jours) et retirée de
            l'application tout de suite.
          </p>
          <div className="flex gap-2">
            <button
              type="button"
              onClick={() => run("approve")}
              disabled={decide.isPending}
              className="inline-flex flex-1 items-center justify-center gap-1.5 rounded-xl bg-red-600 px-3 py-2 text-sm font-semibold text-white transition-colors hover:bg-red-700 disabled:opacity-50"
            >
              <Trash2 className="h-4 w-4" />
              {decide.isPending ? "Suppression…" : "Supprimer la photo"}
            </button>
            <button
              type="button"
              onClick={() => setConfirming(false)}
              className="rounded-xl border border-gray-300 bg-white px-3 py-2 text-sm"
            >
              Annuler
            </button>
          </div>
        </div>
      ) : (
        <div className="flex gap-2">
          <button
            type="button"
            onClick={() => setConfirming(true)}
            disabled={!canRemove || decide.isPending}
            title={
              canRemove ? undefined : "La suppression Drive n'est pas configurée sur le serveur"
            }
            className="inline-flex flex-1 items-center justify-center gap-1.5 rounded-xl bg-blue-600 px-3 py-2 text-sm font-semibold text-white transition-all hover:bg-blue-700 active:scale-[0.98] disabled:opacity-50"
          >
            <CircleCheck className="h-4 w-4" />
            Approuver
          </button>
          <button
            type="button"
            onClick={() => run("reject")}
            disabled={decide.isPending}
            className="inline-flex flex-1 items-center justify-center gap-1.5 rounded-xl border border-gray-300 bg-white px-3 py-2 text-sm font-semibold transition-all hover:bg-gray-50 active:scale-[0.98] disabled:opacity-50"
          >
            <CircleX className="h-4 w-4" />
            Refuser
          </button>
        </div>
      )}
      {!canRemove && (
        <p className="text-xs text-gray-600">
          La suppression de photos n'est pas configurée sur le serveur : les signalements peuvent
          seulement être refusés.
        </p>
      )}
    </div>
  );
}

function AlbumClubSelect({ report }: { report: PhotoReport }) {
  const { data: clubs = [] } = useClubs();
  const setAlbumClub = useSetAlbumClub();
  return (
    <label className="block text-xs text-gray-600">
      Club responsable de cet album (prochains signalements)
      <select
        defaultValue={report.club_id ?? ""}
        key={`${report.id}-${report.club_id ?? ""}`}
        disabled={setAlbumClub.isPending}
        onChange={async (e) => {
          if (!report.album_id) return;
          try {
            await setAlbumClub.mutateAsync({
              albumId: report.album_id,
              clubId: e.target.value === "" ? null : Number(e.target.value),
            });
            toast.success("Club de l'album mis à jour");
          } catch {
            toast.error("Impossible de modifier le club de l'album");
          }
        }}
        className="mt-1 block w-full rounded-xl border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900"
      >
        <option value="">Club par défaut (Club Photo)</option>
        {clubs.map((club) => (
          <option key={club.id_clubs} value={club.id_clubs}>
            {club.name}
          </option>
        ))}
      </select>
    </label>
  );
}

function Deliveries({ report }: { report: PhotoReport }) {
  const retry = useRetryPhotoReport();
  const canRetry = report.deliveries?.some((d) => d.status === "failed") ?? false;
  return (
    <details className="group rounded-2xl border border-gray-200">
      <summary className="flex cursor-pointer list-none items-center justify-between px-4 py-3 text-xs font-semibold uppercase tracking-wider text-gray-500">
        Emails — {DELIVERY_STATUS_LABELS[report.delivery_status]} ({report.recipients_sent}/
        {report.recipients})<span className="transition-transform group-open:rotate-180">▾</span>
      </summary>
      <div className="space-y-2 border-t border-gray-200 px-4 py-3">
        {canRetry && (
          <button
            type="button"
            onClick={async () => {
              try {
                const { retried } = await retry.mutateAsync(report.id);
                toast.success(retried > 0 ? "Nouvel envoi lancé" : "Aucun email en échec");
              } catch {
                toast.error("Échec du nouvel envoi");
              }
            }}
            disabled={retry.isPending}
            className="inline-flex items-center gap-1 rounded-xl border border-gray-300 px-3 py-1.5 text-xs hover:bg-gray-50 disabled:opacity-50"
          >
            <RefreshCw className="h-3 w-3" />
            Renvoyer les échecs
          </button>
        )}
        {report.deliveries?.length === 0 && (
          <p className="text-xs text-gray-500">Aucun destinataire trouvé.</p>
        )}
        <ul className="divide-y divide-gray-100 text-xs">
          {report.deliveries?.map((d) => (
            <li key={d.email} className="py-1.5">
              <div className="flex items-center justify-between gap-3">
                <span className="break-all">{d.email}</span>
                <span className="whitespace-nowrap text-gray-600">
                  {ROLE_LABELS[d.role]} · {DELIVERY_LABELS[d.status]}
                </span>
              </div>
              {d.last_error && d.status !== "sent" && (
                <p className="mt-1 break-words text-red-600">{d.last_error}</p>
              )}
            </li>
          ))}
        </ul>
      </div>
    </details>
  );
}
