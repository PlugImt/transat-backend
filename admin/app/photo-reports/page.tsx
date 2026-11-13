"use client";

import { useQueryClient } from "@tanstack/react-query";
import {
  AlertTriangle,
  CheckCheck,
  ChevronLeft,
  ChevronRight,
  CircleCheck,
  CircleX,
  Clock,
  Flag,
  Inbox,
  Loader2,
} from "lucide-react";
import { useEffect, useRef, useState } from "react";
import toast from "react-hot-toast";
import Modal from "@/components/Modal";
import PhotoReportModal from "@/components/PhotoReportModal";
import { PhotoPreview } from "@/components/photoReports/PhotoPreview";
import { photoReportsApi } from "@/lib/api";
import { usePhotoReports } from "@/lib/hooks";
import { formatDateTime, REASON_LABELS, STATUS_LABELS } from "@/lib/photoReports";
import type { PhotoReport, PhotoReportReason, PhotoReportStatus } from "@/lib/types";

const PAGE_SIZE = 24;

const TABS: { id: PhotoReportStatus; icon: typeof Clock }[] = [
  { id: "pending", icon: Clock },
  { id: "approved", icon: CircleCheck },
  { id: "rejected", icon: CircleX },
];

const REASON_BADGE: Record<PhotoReportReason, string> = {
  appears_in_photo: "bg-violet-100 text-violet-800",
  inappropriate: "bg-red-100 text-red-800",
  poor_quality: "bg-sky-100 text-sky-800",
  other: "bg-gray-100 text-gray-700",
};

export default function PhotoReportsPage() {
  const [status, setStatus] = useState<PhotoReportStatus>("pending");
  const [reason, setReason] = useState<PhotoReportReason | "">("");
  const [offset, setOffset] = useState(0);
  const [openId, setOpenId] = useState<number | null>(null);
  const [approveAllOpen, setApproveAllOpen] = useState(false);

  const { data, isLoading, isError, isFetching } = usePhotoReports({
    status,
    reason,
    limit: PAGE_SIZE,
    offset,
  });
  const reports = data?.reports ?? [];
  const total = data?.total ?? 0;
  const counts = data?.counts;

  // Emails link here with ?report=<id>.
  useEffect(() => {
    const id = Number(new URLSearchParams(window.location.search).get("report"));
    if (Number.isInteger(id) && id > 0) setOpenId(id);
  }, []);

  const closeModal = () => {
    setOpenId(null);
    if (window.location.search) window.history.replaceState(null, "", window.location.pathname);
  };

  const last = Math.min(offset + PAGE_SIZE, total);
  const canRemove = data?.can_remove ?? true;
  // On the pending tab, total counts the reports the current reason filter shows: what "approve all" would take.
  const pendingShown = status === "pending" ? total : 0;

  return (
    <div className="mx-auto max-w-7xl p-6 lg:p-8">
      <header className="mb-6 flex flex-wrap items-end justify-between gap-4 animate-rise">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-gray-900">
            Signalements de photos
          </h1>
          <p className="mt-1 text-sm text-gray-500">
            Approuver un signalement supprime la photo du Drive et de l'application ; le refuser la
            laisse en ligne.
          </p>
        </div>
        <div className="flex flex-wrap items-center gap-2">
          {status === "pending" && pendingShown > 0 && (
            <button
              type="button"
              onClick={() => setApproveAllOpen(true)}
              disabled={!canRemove}
              title={
                canRemove
                  ? "Approuver tous les signalements en attente affichés"
                  : "La suppression de photos n'est pas configurée sur le serveur"
              }
              className="flex items-center gap-2 rounded-xl bg-red-600 px-4 py-2 text-sm font-semibold text-white shadow-sm transition-colors hover:bg-red-700 disabled:cursor-not-allowed disabled:opacity-50"
            >
              <CheckCheck className="h-4 w-4" />
              Tout approuver ({pendingShown})
            </button>
          )}
          <select
            value={reason}
            onChange={(e) => {
              setReason(e.target.value as PhotoReportReason | "");
              setOffset(0);
            }}
            className="rounded-xl border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900"
            aria-label="Filtrer par motif"
          >
            <option value="">Tous les motifs</option>
            {Object.entries(REASON_LABELS).map(([value, label]) => (
              <option key={value} value={value}>
                {label}
              </option>
            ))}
          </select>
        </div>
      </header>

      <div
        role="tablist"
        aria-label="Statut"
        className="mb-6 inline-flex gap-1 rounded-2xl border border-gray-200 bg-white p-1"
      >
        {TABS.map(({ id, icon: Icon }) => {
          const active = status === id;
          return (
            <button
              key={id}
              type="button"
              role="tab"
              aria-selected={active}
              onClick={() => {
                setStatus(id);
                setOffset(0);
              }}
              className={`flex items-center gap-2 rounded-xl px-4 py-2 text-sm font-medium transition-all duration-200 ${
                active ? "bg-blue-600 text-white shadow" : "text-gray-600 hover:bg-gray-100"
              }`}
            >
              <Icon className="h-4 w-4" />
              {STATUS_LABELS[id]}
              {counts && (
                <span
                  className={`rounded-full px-2 py-0.5 text-xs ${active ? "bg-white/20" : "bg-gray-100 text-gray-600"}`}
                >
                  {counts[id]}
                </span>
              )}
            </button>
          );
        })}
      </div>

      {isError && (
        <p className="rounded-2xl border border-red-200 bg-red-50 p-6 text-center text-sm text-red-700">
          Impossible de charger les signalements.
        </p>
      )}

      {isLoading && (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
          {Array.from({ length: 8 }, (_, i) => (
            // biome-ignore lint/suspicious/noArrayIndexKey: static placeholders
            <div key={i} className="h-72 animate-pulse rounded-2xl bg-gray-100" />
          ))}
        </div>
      )}

      {!isLoading && !isError && reports.length === 0 && (
        <div className="rounded-2xl border border-dashed border-gray-300 py-16 text-center animate-rise">
          {status === "pending" ? (
            <Inbox className="mx-auto h-10 w-10 text-gray-400" />
          ) : (
            <Flag className="mx-auto h-10 w-10 text-gray-400" />
          )}
          <h3 className="mt-3 text-sm font-semibold text-gray-900">
            {status === "pending" ? "Rien à traiter" : "Aucun signalement ici"}
          </h3>
          <p className="mt-1 text-sm text-gray-500">
            {status === "pending"
              ? "Les nouveaux signalements de l'application apparaîtront ici."
              : "Les signalements traités sont conservés avec leur décision."}
          </p>
        </div>
      )}

      <ul
        className={`grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 transition-opacity ${isFetching ? "opacity-70" : ""}`}
      >
        {reports.map((report, i) => (
          <li
            key={report.id}
            className="animate-rise"
            style={{ animationDelay: `${Math.min(i, 12) * 35}ms` }}
          >
            <ReportCard report={report} onOpen={() => setOpenId(report.id)} />
          </li>
        ))}
      </ul>

      {total > PAGE_SIZE && (
        <div className="mt-6 flex items-center justify-between text-sm text-gray-600">
          <span>
            {offset + 1}–{last} sur {total}
          </span>
          <div className="flex gap-2">
            <button
              type="button"
              disabled={offset === 0}
              onClick={() => setOffset(Math.max(0, offset - PAGE_SIZE))}
              className="rounded-xl border border-gray-300 bg-white p-2 transition-colors hover:bg-gray-50 disabled:opacity-40"
              aria-label="Page précédente"
            >
              <ChevronLeft className="h-4 w-4" />
            </button>
            <button
              type="button"
              disabled={last >= total}
              onClick={() => setOffset(offset + PAGE_SIZE)}
              className="rounded-xl border border-gray-300 bg-white p-2 transition-colors hover:bg-gray-50 disabled:opacity-40"
              aria-label="Page suivante"
            >
              <ChevronRight className="h-4 w-4" />
            </button>
          </div>
        </div>
      )}

      <PhotoReportModal reportId={openId} onClose={closeModal} canRemove={canRemove} />

      <ApproveAllModal
        isOpen={approveAllOpen}
        onClose={() => setApproveAllOpen(false)}
        reason={reason}
        count={pendingShown}
      />
    </div>
  );
}

function ReportCard({ report, onOpen }: { report: PhotoReport; onOpen: () => void }) {
  return (
    <button
      type="button"
      onClick={onOpen}
      className="group flex h-full w-full flex-col overflow-hidden rounded-2xl border border-gray-200 bg-white text-left shadow-sm transition-all duration-200 hover:-translate-y-0.5 hover:border-blue-300 hover:shadow-lg"
    >
      <div className="relative">
        <PhotoPreview
          photoId={report.photo_id}
          alt={report.photo_name}
          width={480}
          removed={report.status === "approved" || !report.photo_available}
          className="aspect-[4/3] w-full transition-transform duration-500 group-hover:scale-[1.03]"
        />
        <span
          className={`absolute left-3 top-3 rounded-full px-2.5 py-1 text-xs font-medium shadow-sm ${REASON_BADGE[report.reason]}`}
        >
          {REASON_LABELS[report.reason]}
        </span>
      </div>
      <div className="flex flex-1 flex-col gap-1 p-4">
        <p className="truncate text-sm font-semibold text-gray-900">
          {report.album_title || "Sans événement"}
        </p>
        <p className="truncate text-xs text-gray-500">{report.photo_name}</p>
        {report.comment && (
          <p className="mt-1 line-clamp-2 text-xs text-gray-600">« {report.comment} »</p>
        )}
        <div className="mt-auto flex items-center justify-between pt-3 text-xs text-gray-500">
          <span className="truncate">{report.reporter.name}</span>
          <span className="shrink-0">{formatDateTime(report.created_at)}</span>
        </div>
      </div>
    </button>
  );
}

type ApprovePhase = "confirm" | "running" | "done";

/**
 * Approves every pending report the list shows (all, or those of the selected reason). Each approval deletes a
 * photo from Drive, so the server works in batches of a few seconds; this keeps calling it and shows the progress.
 */
function ApproveAllModal({
  isOpen,
  onClose,
  reason,
  count,
}: {
  isOpen: boolean;
  onClose: () => void;
  reason: PhotoReportReason | "";
  count: number;
}) {
  const queryClient = useQueryClient();
  const [phase, setPhase] = useState<ApprovePhase>("confirm");
  const [progress, setProgress] = useState({ approved: 0, photos: 0, failed: 0, remaining: 0 });
  const [error, setError] = useState<string | null>(null);
  const stop = useRef(false);
  // The count when the run started: the list refreshes underneath while it goes.
  const [target, setTarget] = useState(count);

  useEffect(() => {
    if (isOpen) {
      setPhase("confirm");
      setProgress({ approved: 0, photos: 0, failed: 0, remaining: 0 });
      setError(null);
      setTarget(count);
    }
  }, [isOpen, count]);

  const run = async () => {
    stop.current = false;
    setPhase("running");
    setError(null);
    const skip: number[] = [];
    let approved = 0;
    let photos = 0;
    try {
      for (;;) {
        const r = await photoReportsApi.approveAll(reason, skip);
        approved += r.approved;
        photos += r.photos;
        skip.push(...r.failed);
        setProgress({ approved, photos, failed: skip.length, remaining: r.remaining });
        queryClient.invalidateQueries({ queryKey: ["photo-reports"] });
        // Done, stopped, or a batch that moved nothing (it would loop forever).
        if (r.remaining === 0 || stop.current || (r.approved === 0 && r.failed.length === 0)) break;
      }
    } catch {
      setError("Le serveur n'a pas pu terminer. Les signalements déjà approuvés le restent.");
    }
    queryClient.invalidateQueries({ queryKey: ["photo-reports"] });
    queryClient.invalidateQueries({ queryKey: ["photo-report"] });
    setPhase("done");
    if (approved > 0)
      toast.success(
        `${approved} signalement${approved > 1 ? "s" : ""} approuvé${approved > 1 ? "s" : ""}`,
      );
  };

  const done = progress.approved + progress.failed;
  const percent = target > 0 ? Math.min(100, Math.round((done / target) * 100)) : 0;
  const scope = reason ? `« ${REASON_LABELS[reason]} »` : "en attente";

  return (
    <Modal
      isOpen={isOpen}
      // Closing mid-run would hide a destructive job that keeps going: stop it first.
      onClose={phase === "running" ? () => {} : onClose}
      title="Tout approuver"
    >
      {phase === "confirm" && (
        <div className="space-y-4 text-sm text-gray-600">
          <div className="flex gap-3 rounded-xl border border-red-200 bg-red-50 p-3 text-red-800">
            <AlertTriangle className="mt-0.5 h-5 w-5 shrink-0" />
            <p>
              Les <strong>{count}</strong> signalements {scope} vont être approuvés : chaque photo
              signalée est supprimée du Drive et de l'application. Le Drive la garde 30 jours dans
              sa corbeille.
            </p>
          </div>
          <div className="flex justify-end gap-2">
            <button
              type="button"
              onClick={onClose}
              className="rounded-xl border border-gray-300 px-4 py-2 font-medium text-gray-700 hover:bg-gray-50"
            >
              Annuler
            </button>
            <button
              type="button"
              onClick={run}
              className="flex items-center gap-2 rounded-xl bg-red-600 px-4 py-2 font-semibold text-white hover:bg-red-700"
            >
              <CheckCheck className="h-4 w-4" />
              Approuver les {count} signalements
            </button>
          </div>
        </div>
      )}

      {phase !== "confirm" && (
        <div className="space-y-4 text-sm text-gray-600">
          <div>
            <div className="mb-2 flex items-center justify-between">
              <span className="flex items-center gap-2 font-medium text-gray-900">
                {phase === "running" && <Loader2 className="h-4 w-4 animate-spin" />}
                {phase === "running" ? "Approbation en cours…" : "Terminé"}
              </span>
              <span>
                {done} / {target}
              </span>
            </div>
            <div className="h-2 overflow-hidden rounded-full bg-gray-100">
              <div
                className="h-full rounded-full bg-red-600 transition-[width] duration-500"
                style={{ width: `${percent}%` }}
              />
            </div>
          </div>
          <ul className="space-y-1">
            <li>
              {progress.approved} signalement{progress.approved > 1 ? "s" : ""} approuvé
              {progress.approved > 1 ? "s" : ""}, {progress.photos} photo
              {progress.photos > 1 ? "s" : ""} supprimée{progress.photos > 1 ? "s" : ""}
            </li>
            {progress.failed > 0 && (
              <li className="text-amber-700">
                {progress.failed} refusé{progress.failed > 1 ? "s" : ""} par le Drive : ils restent
                en attente, à reprendre un par un.
              </li>
            )}
            {phase === "done" && progress.remaining > 0 && !error && (
              <li className="text-amber-700">{progress.remaining} restent en attente.</li>
            )}
          </ul>
          {error && <p className="rounded-xl bg-red-50 p-3 text-red-700">{error}</p>}
          <div className="flex justify-end">
            {phase === "running" ? (
              <button
                type="button"
                onClick={() => {
                  stop.current = true;
                }}
                className="rounded-xl border border-gray-300 px-4 py-2 font-medium text-gray-700 hover:bg-gray-50"
              >
                Arrêter après ce lot
              </button>
            ) : (
              <button
                type="button"
                onClick={onClose}
                className="rounded-xl bg-blue-600 px-4 py-2 font-semibold text-white hover:bg-blue-700"
              >
                Fermer
              </button>
            )}
          </div>
        </div>
      )}
    </Modal>
  );
}
