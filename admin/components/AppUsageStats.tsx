"use client";

import {
  Camera,
  CheckCircle2,
  Hourglass,
  ScanFace,
  Smartphone,
  UserCheck,
  type Users,
} from "lucide-react";
import { useAppUsageStats } from "@/lib/hooks";
import type { AppUsageStats as Stats } from "@/lib/types";

const nf = new Intl.NumberFormat("fr-FR");
const pf = new Intl.NumberFormat("fr-FR", { maximumFractionDigits: 1 });

const plural = (n: number, one: string, many: string) => `${nf.format(n)} ${n > 1 ? many : one}`;

const PLATFORM_LABEL: Record<string, string> = { ios: "iOS", android: "Android" };
const PLATFORM_COLOR: Record<string, string> = { ios: "bg-gray-800", android: "bg-green-600" };

function percent(part: number, total: number) {
  return total > 0 ? (100 * part) / total : 0;
}

function Card({
  title,
  value,
  hint,
  icon: Icon,
  color,
}: {
  title: string;
  value: string;
  hint?: string;
  icon: typeof Users;
  color: string;
}) {
  return (
    <div className="bg-white overflow-hidden shadow rounded-lg p-5 flex items-center">
      <div className={`p-3 rounded-md shrink-0 bg-gray-50`}>
        <Icon className={`h-6 w-6 ${color}`} />
      </div>
      <div className="ml-5 w-0 flex-1">
        <dt className="text-sm font-medium text-gray-500 truncate">{title}</dt>
        <dd className="text-2xl font-semibold text-gray-900">{value}</dd>
        {hint && <dd className="text-xs text-gray-500 mt-0.5">{hint}</dd>}
      </div>
    </div>
  );
}

/** Platforms and app versions of the installations used in the last 30 days. */
function DevicesSection({ devices }: { devices: Stats["devices"] }) {
  const total = devices.devices;
  return (
    <section>
      <h2 className="text-lg font-medium text-gray-900 mb-1">Application mobile</h2>
      <p className="text-sm text-gray-500 mb-4">
        Appareils connectés au moins une fois sur les {devices.active_window_days} derniers jours.
      </p>
      <div className="grid grid-cols-1 xl:grid-cols-3 gap-6">
        <div className="bg-white shadow rounded-lg p-6 space-y-4">
          <div className="flex items-baseline justify-between">
            <h3 className="text-base font-medium text-gray-900">iOS / Android</h3>
            <span className="text-sm text-gray-500">
              {plural(total, "appareil", "appareils")} ·{" "}
              {plural(devices.users, "personne", "personnes")}
            </span>
          </div>
          {total === 0 ? (
            <p className="text-sm text-gray-500">Aucun appareil récent</p>
          ) : (
            <>
              <div className="flex h-4 w-full overflow-hidden rounded-full bg-gray-100">
                {devices.platforms.map((p) => (
                  <div
                    key={p.platform}
                    className={PLATFORM_COLOR[p.platform] ?? "bg-gray-400"}
                    style={{ width: `${percent(p.devices, total)}%` }}
                    title={`${PLATFORM_LABEL[p.platform] ?? p.platform} : ${pf.format(percent(p.devices, total))} %`}
                  />
                ))}
              </div>
              <ul className="space-y-2 text-sm">
                {devices.platforms.map((p) => (
                  <li key={p.platform} className="flex items-center justify-between gap-3">
                    <span className="flex items-center gap-2">
                      <span
                        className={`h-3 w-3 rounded-full ${PLATFORM_COLOR[p.platform] ?? "bg-gray-400"}`}
                      />
                      {PLATFORM_LABEL[p.platform] ?? p.platform}
                    </span>
                    <span className="text-gray-700">
                      <span className="font-semibold">
                        {pf.format(percent(p.devices, total))} %
                      </span>
                      <span className="text-gray-500">
                        {" "}
                        · {plural(p.devices, "appareil", "appareils")} ·{" "}
                        {plural(p.users, "personne", "personnes")}
                      </span>
                    </span>
                  </li>
                ))}
              </ul>
            </>
          )}
        </div>

        <div className="bg-white shadow rounded-lg p-6 xl:col-span-2">
          <h3 className="text-base font-medium text-gray-900 mb-3">
            Versions de l&apos;application
          </h3>
          {devices.versions.length === 0 ? (
            <p className="text-sm text-gray-500">Aucun appareil récent</p>
          ) : (
            <div className="overflow-x-auto max-h-80">
              <table className="min-w-full text-sm">
                <thead>
                  <tr className="text-left text-xs uppercase text-gray-500">
                    <th className="py-2 pr-4">Version</th>
                    <th className="py-2 pr-4">Plateforme</th>
                    <th className="py-2 pr-4 text-right">Appareils</th>
                    <th className="py-2 pr-4 text-right">Personnes</th>
                    <th className="py-2 text-right">Part</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-100">
                  {devices.versions.map((v) => (
                    <tr key={`${v.platform}-${v.version}`}>
                      <td className="py-2 pr-4 font-mono text-xs text-gray-800">
                        {v.version || <span className="text-gray-400">inconnue</span>}
                      </td>
                      <td className="py-2 pr-4">{PLATFORM_LABEL[v.platform] ?? v.platform}</td>
                      <td className="py-2 pr-4 text-right">{nf.format(v.devices)}</td>
                      <td className="py-2 pr-4 text-right">{nf.format(v.users)}</td>
                      <td className="py-2 text-right text-gray-600">
                        {pf.format(percent(v.devices, total))} %
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>
      </div>
    </section>
  );
}

/** Face recognition usage and library progress. Aggregates only: no user is named. */
function FacesSection({ faces }: { faces: Stats["faces"] }) {
  const done = percent(faces.photos_analysed, faces.photos_total);
  return (
    <section>
      <h2 className="text-lg font-medium text-gray-900 mb-1">Reconnaissance faciale</h2>
      <p className="text-sm text-gray-500 mb-4">
        {!faces.available
          ? "Désactivée sur ce serveur (FACES_ML_URL non défini)."
          : faces.scoped
            ? "Analyse limitée aux dossiers de test (FACES_SCOPE_FOLDERS)."
            : "Analyse de toute la galerie."}
      </p>
      <div className="grid grid-cols-1 md:grid-cols-2 2xl:grid-cols-4 gap-6 mb-6">
        <Card
          title="Personnes inscrites"
          value={nf.format(faces.users_enabled)}
          hint={plural(faces.users_with_selfie, "avec un selfie", "avec un selfie")}
          icon={UserCheck}
          color="text-blue-600"
        />
        <Card
          title="Avec au moins une photo"
          value={nf.format(faces.users_with_photos)}
          hint={`${pf.format(faces.average_photos_by_user)} photos en moyenne par personne`}
          icon={Camera}
          color="text-orange-600"
        />
        <Card
          title="Visages détectés"
          value={nf.format(faces.faces_found)}
          hint={`${nf.format(faces.rejected_matches)} « Ce n'est pas moi »`}
          icon={ScanFace}
          color="text-purple-600"
        />
        <Card
          title="Photos analysées"
          value={`${nf.format(faces.photos_analysed)} / ${nf.format(faces.photos_total)}`}
          hint={`${pf.format(done)} %`}
          icon={CheckCircle2}
          color="text-green-600"
        />
      </div>
      <div className="bg-white shadow rounded-lg p-6 space-y-3">
        <div className="flex flex-wrap items-baseline justify-between gap-2">
          <h3 className="text-base font-medium text-gray-900 flex items-center gap-2">
            <Hourglass className="h-4 w-4 text-gray-500" /> Analyse de la galerie
          </h3>
          <span className="text-sm text-gray-500">
            {nf.format(faces.photos_pending)} en attente · {nf.format(faces.photos_failed)} en échec
            {faces.last_scan_at &&
              ` · dernière analyse ${new Date(faces.last_scan_at).toLocaleString("fr-FR", {
                dateStyle: "short",
                timeStyle: "short",
              })}`}
          </span>
        </div>
        <div className="flex h-3 w-full overflow-hidden rounded-full bg-gray-100">
          <div className="bg-green-500" style={{ width: `${done}%` }} />
          <div
            className="bg-red-400"
            style={{ width: `${percent(faces.photos_failed, faces.photos_total)}%` }}
          />
        </div>
      </div>
    </section>
  );
}

export default function AppUsageStats() {
  const { data, isLoading, isError } = useAppUsageStats();
  if (isLoading) {
    return (
      <div className="text-sm text-gray-400">
        Chargement des statistiques de l&apos;application…
      </div>
    );
  }
  if (isError || !data) {
    return (
      <div className="flex items-center gap-2 text-sm text-gray-500">
        <Smartphone className="h-4 w-4" /> Statistiques de l&apos;application indisponibles
      </div>
    );
  }
  return (
    <>
      <DevicesSection devices={data.devices} />
      <FacesSection faces={data.faces} />
    </>
  );
}
