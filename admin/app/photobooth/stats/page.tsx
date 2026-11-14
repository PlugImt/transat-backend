"use client";

import { format } from "date-fns";
import { fr } from "date-fns/locale";
import {
  CalendarHeart,
  ExternalLink,
  Images,
  Info,
  Loader2,
  QrCode,
  ScanLine,
  Users,
} from "lucide-react";
import type { ReactNode } from "react";
import { usePhotoboothStats } from "@/lib/hooks";
import type { PhotoboothEventStats } from "@/lib/types";

const number = new Intl.NumberFormat("fr-FR", { maximumFractionDigits: 1 });
const percent = new Intl.NumberFormat("fr-FR", { style: "percent", maximumFractionDigits: 0 });
const eventDay = (date: string) =>
  format(new Date(`${date}T12:00:00`), "EEE d MMM yyyy", { locale: fr });

function Tile({
  icon,
  label,
  value,
  hint,
}: {
  icon: ReactNode;
  label: string;
  value: string;
  hint: string;
}) {
  return (
    <div className="rounded-2xl border border-gray-200 bg-white p-5 shadow-sm animate-rise">
      <div className="flex items-center gap-2 text-sm font-medium text-gray-500">
        {icon}
        {label}
      </div>
      <p className="mt-2 text-3xl font-bold tracking-tight text-gray-900 tabular-nums">{value}</p>
      <p className="mt-1 text-xs text-gray-500">{hint}</p>
    </div>
  );
}

/** A count with a bar scaled on the evening with the most of it. */
function Bar({ value, max, sub }: { value: number; max: number; sub?: string }) {
  return (
    <div className="min-w-28">
      <div className="flex items-baseline gap-1.5">
        <span className="font-semibold text-gray-900 tabular-nums">{value}</span>
        {sub && <span className="text-xs text-gray-500">{sub}</span>}
      </div>
      <div className="mt-1 h-1.5 rounded-full bg-gray-100">
        <div
          className="h-1.5 rounded-full bg-blue-500"
          style={{ width: `${max > 0 ? Math.max(value > 0 ? 3 : 0, (value / max) * 100) : 0}%` }}
        />
      </div>
    </div>
  );
}

function Row({
  e,
  maxScanners,
  maxAlbum,
}: {
  e: PhotoboothEventStats;
  maxScanners: number;
  maxAlbum: number;
}) {
  const rate = e.sessions > 0 ? e.scanned_sessions / e.sessions : 0;
  return (
    <tr className="border-t border-gray-100 align-middle">
      <td className="px-4 py-3">
        <p className="font-medium text-gray-900">{e.name}</p>
        <p className="text-xs text-gray-500 capitalize">{eventDay(e.date)}</p>
      </td>
      <td className="px-4 py-3 tabular-nums text-gray-900">{e.sessions}</td>
      <td className="px-4 py-3 tabular-nums">
        {e.sessions > 0 ? (
          <>
            <span className="text-gray-900">
              {e.scanned_sessions} / {e.sessions}
            </span>
            <span className="ml-1.5 text-xs text-gray-500">{percent.format(rate)}</span>
          </>
        ) : (
          <span className="text-gray-400">—</span>
        )}
      </td>
      <td className="px-4 py-3">
        <Bar
          value={e.scanners}
          max={maxScanners}
          sub={e.scan_views > e.scanners ? `${e.scan_views} ouvertures` : undefined}
        />
      </td>
      <td className="px-4 py-3">
        {e.album_id ? (
          <div className="flex items-center gap-3">
            <Bar
              value={e.album_visitors}
              max={maxAlbum}
              sub={e.album_views > e.album_visitors ? `${e.album_views} vues` : undefined}
            />
            <a
              href={`https://transat.dev/albums/${e.album_id}`}
              target="_blank"
              rel="noreferrer"
              title={e.album_title}
              className="text-gray-400 transition-colors hover:text-blue-600"
            >
              <ExternalLink className="h-4 w-4" />
            </a>
          </div>
        ) : (
          <span className="text-xs text-gray-400">Pas encore d'album</span>
        )}
      </td>
    </tr>
  );
}

export default function PhotoboothStatsPage() {
  const { data, isLoading, isError } = usePhotoboothStats();
  const t = data?.totals;
  const events = data?.events ?? [];
  const maxScanners = Math.max(0, ...events.map((e) => e.scanners));
  const maxAlbum = Math.max(0, ...events.map((e) => e.album_visitors));

  return (
    <div className="mx-auto max-w-7xl p-6 lg:p-8">
      <header className="mb-6 animate-rise">
        <h1 className="text-2xl font-bold tracking-tight text-gray-900">Statistiques photobooth</h1>
        <p className="mt-1 max-w-3xl text-sm text-gray-500">
          Combien de personnes ont scanné les QR codes imprimés sur les photos, et combien ont
          ouvert l'album public de la soirée. Une personne est comptée une fois par soirée, même si
          elle ouvre plusieurs liens ou revient plus tard.
        </p>
      </header>

      {isLoading && (
        <div className="flex justify-center py-16 text-gray-400">
          <Loader2 className="h-6 w-6 animate-spin" />
        </div>
      )}
      {isError && (
        <p className="rounded-2xl border border-red-200 bg-red-50 p-6 text-center text-sm text-red-700">
          Impossible de charger les statistiques.
        </p>
      )}

      {t && (
        <>
          <section className="mb-6 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
            <Tile
              icon={<Users className="h-4 w-4" />}
              label="Personnes ayant scanné"
              value={number.format(t.scanners)}
              hint={`${number.format(t.scan_views)} ouvertures de liens au total`}
            />
            <Tile
              icon={<ScanLine className="h-4 w-4" />}
              label="Moyenne par soirée"
              value={number.format(t.avg_scanners_per_event)}
              hint={`personnes qui scannent, sur ${t.events} soirée${t.events > 1 ? "s" : ""}`}
            />
            <Tile
              icon={<QrCode className="h-4 w-4" />}
              label="QR codes scannés"
              value={percent.format(t.scan_rate)}
              hint={`${t.scanned_sessions} sessions sur ${t.sessions} · ${number.format(t.avg_scanners_per_session)} pers. par session`}
            />
            <Tile
              icon={<Images className="h-4 w-4" />}
              label="Visiteurs des albums"
              value={number.format(t.album_visitors)}
              hint={`${number.format(t.avg_album_visitors_per_event)} par soirée · ${number.format(t.album_views)} vues`}
            />
          </section>

          <p className="mb-4 flex items-start gap-2 text-xs text-gray-500">
            <Info className="mt-0.5 h-3.5 w-3.5 shrink-0" />
            {data.tracking_since
              ? `Visites comptées depuis le ${format(new Date(data.tracking_since), "d MMMM yyyy", { locale: fr })} : les soirées plus anciennes n'ont que leurs sessions.`
              : "Aucune visite comptée pour l'instant : les compteurs démarrent avec la mise en ligne du suivi sur le site."}{" "}
            Les visiteurs sans stockage local (navigation privée stricte) ne sont pas comptés.
          </p>

          {events.length === 0 ? (
            <div className="rounded-2xl border border-dashed border-gray-300 p-10 text-center text-sm text-gray-500">
              <CalendarHeart className="mx-auto mb-2 h-6 w-6 text-gray-400" />
              Aucune soirée pour l'instant.
            </div>
          ) : (
            <div className="overflow-x-auto rounded-2xl border border-gray-200 bg-white shadow-sm">
              <table className="w-full text-left text-sm">
                <thead className="bg-gray-50 text-xs font-semibold uppercase tracking-wide text-gray-500">
                  <tr>
                    <th className="px-4 py-3">Soirée</th>
                    <th className="px-4 py-3">Sessions</th>
                    <th className="px-4 py-3">QR scannés</th>
                    <th className="px-4 py-3">Personnes ayant scanné</th>
                    <th className="px-4 py-3">Visiteurs de l'album</th>
                  </tr>
                </thead>
                <tbody>
                  {events.map((e) => (
                    <Row
                      key={`${e.date}-${e.name}`}
                      e={e}
                      maxScanners={maxScanners}
                      maxAlbum={maxAlbum}
                    />
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </>
      )}
    </div>
  );
}
