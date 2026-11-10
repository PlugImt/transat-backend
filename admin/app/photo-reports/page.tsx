"use client";

import { ChevronLeft, ChevronRight, CircleCheck, CircleX, Clock, Flag, Inbox } from "lucide-react";
import { useEffect, useState } from "react";
import PhotoReportModal from "@/components/PhotoReportModal";
import { PhotoPreview } from "@/components/photoReports/PhotoPreview";
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

      <PhotoReportModal
        reportId={openId}
        onClose={closeModal}
        canRemove={data?.can_remove ?? true}
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
