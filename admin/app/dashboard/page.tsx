"use client";

import {
  Activity,
  Building,
  Calendar,
  CircleAlert,
  Clock,
  Gauge,
  ShieldCheck,
  TrendingUp,
  UserPlus,
  Users,
} from "lucide-react";
import { useMemo, useState } from "react";
import ActiveUsersChart from "@/components/ActiveUsersChart";
import ActivityHeatmapChart from "@/components/ActivityHeatmapChart";
import UserGrowthChart from "@/components/UserGrowthChart";
import {
  useActiveUsersOverTime,
  useActivityHeatmap,
  useDashboardStats,
  useEndpointApiStats,
  useGlobalApiStats,
} from "@/lib/hooks";
import type { ActiveUsersPeriod, ActivityHeatmapRange, ApiError } from "@/lib/types";

const RANGES = [
  { label: "7 j", days: 7 },
  { label: "30 j", days: 30 },
  { label: "90 j", days: 90 },
  { label: "Tout", days: 0 },
] as const;

const ACTIVE_USERS_PERIODS: { label: string; value: ActiveUsersPeriod }[] = [
  { label: "Jour", value: "day" },
  { label: "Semaine", value: "week" },
  { label: "Mois", value: "month" },
  { label: "3 mois", value: "quarter" },
  { label: "Année", value: "year" },
];

const HEATMAP_RANGES: { label: string; value: ActivityHeatmapRange }[] = [
  { label: "Cette semaine", value: "week" },
  { label: "Ce mois", value: "month" },
  { label: "Cette année", value: "year" },
  { label: "Tout", value: "all" },
];

const nf = new Intl.NumberFormat("fr-FR");
const pf = new Intl.NumberFormat("fr-FR", { maximumFractionDigits: 1 });

function rateColor(rate: number) {
  if (rate >= 99) return "text-green-600";
  if (rate >= 95) return "text-yellow-600";
  return "text-red-600";
}

function StatCard({
  title,
  value,
  hint,
  icon: Icon,
  color,
  bgColor,
}: {
  title: string;
  value: string;
  hint?: string;
  icon: typeof Users;
  color: string;
  bgColor: string;
}) {
  return (
    <div className="bg-white overflow-hidden shadow rounded-lg p-5 flex items-center">
      <div className={`p-3 rounded-md shrink-0 ${bgColor}`}>
        <Icon className={`h-6 w-6 ${color}`} />
      </div>
      <div className="ml-5 w-0 flex-1">
        <div className="text-sm font-medium text-gray-500 truncate">{title}</div>
        <div className="text-lg font-semibold text-gray-900">{value}</div>
        {hint && <div className="text-xs text-gray-500 truncate">{hint}</div>}
      </div>
    </div>
  );
}

export default function DashboardPage() {
  const { data: stats, isLoading, error } = useDashboardStats();
  const { data: apiGlobal } = useGlobalApiStats();
  const { data: endpoints } = useEndpointApiStats();
  const [rangeDays, setRangeDays] = useState<number>(30);
  const [activeUsersPeriod, setActiveUsersPeriod] = useState<ActiveUsersPeriod>("day");
  const [heatmapRange, setHeatmapRange] = useState<ActivityHeatmapRange>("month");
  const { data: activeUsersOverTime } = useActiveUsersOverTime(activeUsersPeriod);
  const { data: activityHeatmap } = useActivityHeatmap(heatmapRange);

  const growth = stats?.userGrowth;

  const chartData = useMemo(() => {
    if (!growth) return [];
    return rangeDays > 0 ? growth.slice(-rangeDays) : growth;
  }, [growth, rangeDays]);

  const growthSummary = useMemo(() => {
    if (!growth || growth.length === 0) return null;
    const sum = (arr: typeof growth) => arr.reduce((acc, d) => acc + d.count, 0);
    const last7 = sum(growth.slice(-7));
    const prev7 = sum(growth.slice(-14, -7));
    const last30 = sum(growth.slice(-30));
    const best = growth.reduce((a, b) => (b.count > a.count ? b : a), growth[0]);
    return {
      last7,
      last30,
      // null quand la semaine précédente est vide (variation non définie)
      weekDelta: prev7 > 0 ? ((last7 - prev7) / prev7) * 100 : null,
      avgPerDay: sum(growth) / growth.length,
      best,
    };
  }, [growth]);

  const topEndpoints = useMemo(
    () => [...(endpoints ?? [])].sort((a, b) => b.request_count - a.request_count).slice(0, 8),
    [endpoints],
  );
  const slowEndpoints = useMemo(
    () =>
      [...(endpoints ?? [])]
        .filter((e) => e.request_count >= 5)
        .sort((a, b) => b.avg_duration_ms - a.avg_duration_ms)
        .slice(0, 5),
    [endpoints],
  );
  const failingEndpoints = useMemo(
    () =>
      [...(endpoints ?? [])]
        .filter((e) => e.error_count > 0)
        .sort((a, b) => b.error_count - a.error_count)
        .slice(0, 5),
    [endpoints],
  );

  if (isLoading) {
    return (
      <div className="p-6">
        <div className="animate-pulse space-y-6">
          <div className="h-8 bg-gray-200 rounded w-1/4"></div>
          <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
            {[1, 2, 3, 4].map((i) => (
              <div key={i} className="h-32 bg-gray-200 rounded"></div>
            ))}
          </div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="p-6">
        <div className="bg-red-50 border border-red-200 rounded-md p-4">
          <div className="text-sm text-red-700">
            {(error as ApiError)?.message || "Échec de la récupération des statistiques"}
          </div>
        </div>
      </div>
    );
  }

  const totalUsers = stats?.totalUsers || 0;
  const unverified = stats?.unverifiedUsers || 0;
  const verifiedRate = totalUsers > 0 ? ((totalUsers - unverified) / totalUsers) * 100 : 0;

  const statCards = [
    {
      title: "Utilisateurs",
      value: nf.format(totalUsers),
      hint: growthSummary ? `+${nf.format(growthSummary.last7)} cette semaine` : undefined,
      icon: Users,
      color: "text-blue-600",
      bgColor: "bg-blue-50",
    },
    {
      title: "Non vérifiés",
      value: nf.format(unverified),
      hint: totalUsers > 0 ? `${pf.format(100 - verifiedRate)} % des comptes` : undefined,
      icon: CircleAlert,
      color: "text-red-600",
      bgColor: "bg-red-50",
    },
    {
      title: "Événements",
      value: nf.format(stats?.totalEvents || 0),
      icon: Calendar,
      color: "text-green-600",
      bgColor: "bg-green-50",
    },
    {
      title: "Clubs",
      value: nf.format(stats?.totalClubs || 0),
      icon: Building,
      color: "text-purple-600",
      bgColor: "bg-purple-50",
    },
  ];

  return (
    <div className="p-6 space-y-8">
      <h1 className="text-2xl font-bold text-gray-900">Tableau de bord</h1>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {statCards.map((stat) => (
          <StatCard key={stat.title} {...stat} />
        ))}
      </div>

      {growthSummary && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <StatCard
            title="Comptes vérifiés"
            value={`${pf.format(verifiedRate)} %`}
            icon={ShieldCheck}
            color="text-emerald-600"
            bgColor="bg-emerald-50"
          />
          <StatCard
            title="Inscrits sur 7 jours"
            value={nf.format(growthSummary.last7)}
            hint={
              growthSummary.weekDelta === null
                ? "Pas de comparaison possible"
                : `${growthSummary.weekDelta >= 0 ? "+" : ""}${pf.format(growthSummary.weekDelta)} % vs semaine précédente`
            }
            icon={UserPlus}
            color="text-sky-600"
            bgColor="bg-sky-50"
          />
          <StatCard
            title="Inscrits sur 30 jours"
            value={nf.format(growthSummary.last30)}
            hint={`${pf.format(growthSummary.avgPerDay)} / jour en moyenne`}
            icon={TrendingUp}
            color="text-indigo-600"
            bgColor="bg-indigo-50"
          />
          <StatCard
            title="Meilleur jour"
            value={nf.format(growthSummary.best.count)}
            hint={new Date(growthSummary.best.date).toLocaleDateString("fr-FR", {
              day: "numeric",
              month: "long",
              year: "numeric",
            })}
            icon={Calendar}
            color="text-amber-600"
            bgColor="bg-amber-50"
          />
        </div>
      )}

      {growth && growth.length > 0 && (
        <div className="bg-white shadow rounded-lg p-6">
          <div className="flex flex-wrap items-center justify-between gap-3 mb-4">
            <h2 className="text-lg font-medium text-gray-900">Comptes créés</h2>
            <div className="inline-flex rounded-md border border-gray-200 overflow-hidden">
              {RANGES.map((r) => (
                <button
                  key={r.label}
                  type="button"
                  onClick={() => setRangeDays(r.days)}
                  className={`px-3 py-1 text-sm ${
                    rangeDays === r.days
                      ? "bg-blue-600 text-white"
                      : "bg-white text-gray-600 hover:bg-gray-50"
                  }`}
                >
                  {r.label}
                </button>
              ))}
            </div>
          </div>
          <div className="h-72">
            <UserGrowthChart data={chartData} />
          </div>
        </div>
      )}

      {stats?.activeUsers && (
        <section>
          <h2 className="text-lg font-medium text-gray-900 mb-4">Utilisateurs actifs</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-6 mb-6">
            <StatCard
              title="DAU (24 h)"
              value={nf.format(stats.activeUsers.dau)}
              icon={Users}
              color="text-violet-600"
              bgColor="bg-violet-50"
            />
            <StatCard
              title="WAU (7 j)"
              value={nf.format(stats.activeUsers.wau)}
              icon={Users}
              color="text-violet-600"
              bgColor="bg-violet-50"
            />
            <StatCard
              title="MAU (30 j)"
              value={nf.format(stats.activeUsers.mau)}
              hint={
                totalUsers > 0
                  ? `${pf.format((stats.activeUsers.mau / totalUsers) * 100)} % des comptes`
                  : undefined
              }
              icon={Users}
              color="text-violet-600"
              bgColor="bg-violet-50"
            />
            <StatCard
              title="YAU (365 j)"
              value={nf.format(stats.activeUsers.yau)}
              hint={
                totalUsers > 0
                  ? `${pf.format((stats.activeUsers.yau / totalUsers) * 100)} % des comptes`
                  : undefined
              }
              icon={Users}
              color="text-violet-600"
              bgColor="bg-violet-50"
            />
            <StatCard
              title="Stickiness (DAU/MAU)"
              value={
                stats.activeUsers.mau > 0
                  ? `${pf.format((stats.activeUsers.dau / stats.activeUsers.mau) * 100)} %`
                  : "–"
              }
              hint="Part des actifs mensuels qui reviennent chaque jour"
              icon={Gauge}
              color="text-fuchsia-600"
              bgColor="bg-fuchsia-50"
            />
          </div>
          <div className="bg-white shadow rounded-lg p-6 mb-6">
            <div className="flex flex-wrap items-center justify-between gap-3 mb-4">
              <h3 className="text-base font-medium text-gray-900">Utilisateurs actifs</h3>
              <div className="inline-flex rounded-md border border-gray-200 overflow-hidden">
                {ACTIVE_USERS_PERIODS.map((p) => (
                  <button
                    key={p.value}
                    type="button"
                    onClick={() => setActiveUsersPeriod(p.value)}
                    className={`px-3 py-1 text-sm ${
                      activeUsersPeriod === p.value
                        ? "bg-blue-600 text-white"
                        : "bg-white text-gray-600 hover:bg-gray-50"
                    }`}
                  >
                    {p.label}
                  </button>
                ))}
              </div>
            </div>
            <div className="h-64">
              {activeUsersOverTime && activeUsersOverTime.length > 0 ? (
                <ActiveUsersChart data={activeUsersOverTime} period={activeUsersPeriod} />
              ) : (
                <div className="h-full flex items-center justify-center text-sm text-gray-400">
                  Chargement…
                </div>
              )}
            </div>
          </div>
          <div className="bg-white shadow rounded-lg p-6">
            <div className="flex flex-wrap items-center justify-between gap-3 mb-4">
              <h3 className="text-base font-medium text-gray-900">
                Heures d&apos;activité par jour de la semaine
              </h3>
              <div className="inline-flex rounded-md border border-gray-200 overflow-hidden">
                {HEATMAP_RANGES.map((r) => (
                  <button
                    key={r.value}
                    type="button"
                    onClick={() => setHeatmapRange(r.value)}
                    className={`px-3 py-1 text-sm ${
                      heatmapRange === r.value
                        ? "bg-blue-600 text-white"
                        : "bg-white text-gray-600 hover:bg-gray-50"
                    }`}
                  >
                    {r.label}
                  </button>
                ))}
              </div>
            </div>
            {activityHeatmap && activityHeatmap.length > 0 ? (
              <ActivityHeatmapChart data={activityHeatmap} />
            ) : (
              <div className="text-sm text-gray-400">Chargement…</div>
            )}
          </div>
        </section>
      )}

      {apiGlobal && (
        <section>
          <h2 className="text-lg font-medium text-gray-900 mb-4">Activité de l&apos;API</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-6">
            <StatCard
              title="Requêtes"
              value={nf.format(apiGlobal.total_request_count)}
              hint={`depuis le ${new Date(apiGlobal.first_request).toLocaleDateString("fr-FR")}`}
              icon={Activity}
              color="text-blue-600"
              bgColor="bg-blue-50"
            />
            <StatCard
              title="Taux de succès"
              value={`${pf.format(apiGlobal.global_success_rate_percent)} %`}
              hint={`${nf.format(apiGlobal.error_count)} erreurs`}
              icon={ShieldCheck}
              color={rateColor(apiGlobal.global_success_rate_percent)}
              bgColor="bg-gray-50"
            />
            <StatCard
              title="Durée moyenne"
              value={`${pf.format(apiGlobal.global_avg_duration_ms)} ms`}
              hint={`min ${apiGlobal.global_min_duration_ms} ms · max ${nf.format(apiGlobal.global_max_duration_ms)} ms`}
              icon={Clock}
              color="text-orange-600"
              bgColor="bg-orange-50"
            />
            <StatCard
              title="Dernière requête"
              value={new Date(apiGlobal.last_request).toLocaleTimeString("fr-FR", {
                hour: "2-digit",
                minute: "2-digit",
              })}
              hint={new Date(apiGlobal.last_request).toLocaleDateString("fr-FR")}
              icon={Clock}
              color="text-gray-600"
              bgColor="bg-gray-100"
            />
          </div>

          {topEndpoints.length > 0 && (
            <div className="grid grid-cols-1 xl:grid-cols-3 gap-6">
              <div className="bg-white shadow rounded-lg p-6 xl:col-span-2">
                <h3 className="text-base font-medium text-gray-900 mb-3">
                  Endpoints les plus sollicités
                </h3>
                <div className="overflow-x-auto">
                  <table className="min-w-full text-sm">
                    <thead>
                      <tr className="text-left text-xs uppercase text-gray-500">
                        <th className="py-2 pr-4">Endpoint</th>
                        <th className="py-2 pr-4 text-right">Requêtes</th>
                        <th className="py-2 pr-4 text-right">Moy.</th>
                        <th className="py-2 text-right">Succès</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-100">
                      {topEndpoints.map((e) => (
                        <tr key={`${e.method}-${e.endpoint}`}>
                          <td className="py-2 pr-4 font-mono text-xs text-gray-800">
                            <span className="font-semibold text-gray-500 mr-2">{e.method}</span>
                            {e.endpoint}
                          </td>
                          <td className="py-2 pr-4 text-right">{nf.format(e.request_count)}</td>
                          <td className="py-2 pr-4 text-right">
                            {pf.format(e.avg_duration_ms)} ms
                          </td>
                          <td
                            className={`py-2 text-right font-medium ${rateColor(e.success_rate_percent)}`}
                          >
                            {pf.format(e.success_rate_percent)} %
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>

              <div className="space-y-6">
                <div className="bg-white shadow rounded-lg p-6">
                  <h3 className="text-base font-medium text-gray-900 mb-3">Les plus lents</h3>
                  <ul className="space-y-2 text-sm">
                    {slowEndpoints.length === 0 && (
                      <li className="text-gray-500">Pas assez de données</li>
                    )}
                    {slowEndpoints.map((e) => (
                      <li key={`${e.method}-${e.endpoint}`} className="flex justify-between gap-3">
                        <span className="font-mono text-xs truncate">
                          {e.method} {e.endpoint}
                        </span>
                        <span className="font-medium text-orange-600 shrink-0">
                          {pf.format(e.avg_duration_ms)} ms
                        </span>
                      </li>
                    ))}
                  </ul>
                </div>
                <div className="bg-white shadow rounded-lg p-6">
                  <h3 className="text-base font-medium text-gray-900 mb-3">
                    Le plus d&apos;erreurs
                  </h3>
                  <ul className="space-y-2 text-sm">
                    {failingEndpoints.length === 0 && (
                      <li className="text-gray-500">Aucune erreur 🎉</li>
                    )}
                    {failingEndpoints.map((e) => (
                      <li key={`${e.method}-${e.endpoint}`} className="flex justify-between gap-3">
                        <span className="font-mono text-xs truncate">
                          {e.method} {e.endpoint}
                        </span>
                        <span className="font-medium text-red-600 shrink-0">
                          {nf.format(e.error_count)}
                        </span>
                      </li>
                    ))}
                  </ul>
                </div>
              </div>
            </div>
          )}
        </section>
      )}
    </div>
  );
}
