"use client";

import { CalendarHeart, ChartColumnIncreasing, Images } from "lucide-react";
import { useRouter, useSearchParams } from "next/navigation";
import { Suspense, useCallback } from "react";
import { EventsTab } from "@/components/photobooth/Events";
import { ReviewTab } from "@/components/photobooth/Review";
import { StatsTab } from "@/components/photobooth/Stats";
import { usePhotoboothPending } from "@/lib/hooks";

type Tab = "events" | "review" | "stats";

interface View {
  tab: Tab;
  event: number | null;
  folder: string | null;
}

const TABS: { id: Tab; label: string; icon: typeof Images }[] = [
  { id: "events", label: "Soirées", icon: CalendarHeart },
  { id: "review", label: "À approuver", icon: Images },
  { id: "stats", label: "Statistiques", icon: ChartColumnIncreasing },
];

/** ?tab=events|review|stats&event=<id>&folder=<id>: links and the back button land on the same view. */
function viewOf(q: URLSearchParams): View {
  const tab = q.get("tab");
  const event = Number(q.get("event"));
  return {
    tab: tab === "review" || tab === "stats" ? tab : "events",
    event: Number.isSafeInteger(event) && event > 0 ? event : null,
    folder: q.get("folder") || null,
  };
}

function urlOf(v: View) {
  const q = new URLSearchParams();
  if (v.tab !== "events") q.set("tab", v.tab);
  if (v.tab === "events" && v.event !== null) q.set("event", String(v.event));
  if (v.tab === "review" && v.folder) q.set("folder", v.folder);
  return `/photobooth${q.size ? `?${q}` : ""}`;
}

export default function PhotoboothPage() {
  return (
    <Suspense>
      <Photobooth />
    </Suspense>
  );
}

/**
 * The photobooth in one place: the evenings prepared for the booths (settings, collages, frames, then the photos to
 * approve and the audience of each), every folder awaiting approval, and the statistics.
 */
function Photobooth() {
  const router = useRouter();
  const params = useSearchParams();
  const view = viewOf(params);
  const { data: pending = [] } = usePhotoboothPending();
  const toApprove = pending.reduce((n, f) => n + f.pictures, 0);

  const go = useCallback(
    (next: View, push = true) => {
      const url = urlOf(next);
      if (push) router.push(url, { scroll: false });
      else router.replace(url, { scroll: false });
    },
    [router],
  );
  const openEvent = useCallback(
    (event: number | null) => go({ tab: "events", event, folder: null }, event !== null),
    [go],
  );
  const openFolder = useCallback(
    (folder: string | null) => go({ tab: "review", event: null, folder }, folder !== null),
    [go],
  );

  return (
    <div className="mx-auto max-w-7xl p-6 lg:p-8">
      <header className="mb-5 animate-rise">
        <h1 className="text-2xl font-bold tracking-tight text-gray-900">Photobooth</h1>
      </header>

      <div
        role="tablist"
        aria-label="Photobooth"
        className="mb-6 inline-flex max-w-full gap-1 overflow-x-auto rounded-2xl border border-gray-200 bg-white p-1"
      >
        {TABS.map(({ id, label, icon: Icon }) => {
          const active = view.tab === id;
          return (
            <button
              key={id}
              type="button"
              role="tab"
              aria-selected={active}
              onClick={() => go({ tab: id, event: null, folder: null })}
              className={`flex shrink-0 items-center gap-2 rounded-xl px-4 py-2 text-sm font-medium transition-all duration-200 ${
                active ? "bg-blue-600 text-white shadow" : "text-gray-600 hover:bg-gray-100"
              }`}
            >
              <Icon className="h-4 w-4" />
              {label}
              {id === "review" && toApprove > 0 && (
                <span
                  className={`rounded-full px-2 py-0.5 text-xs ${active ? "bg-white/20" : "bg-amber-100 text-amber-800"}`}
                >
                  {toApprove}
                </span>
              )}
            </button>
          );
        })}
      </div>

      {view.tab === "events" && <EventsTab openId={view.event} onOpen={openEvent} />}
      {view.tab === "review" && (
        <ReviewTab openId={view.folder} onOpen={openFolder} onOpenEvent={openEvent} />
      )}
      {view.tab === "stats" && <StatsTab onOpenEvent={openEvent} />}
    </div>
  );
}
