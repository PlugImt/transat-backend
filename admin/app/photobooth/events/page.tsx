"use client";

import { format, isToday } from "date-fns";
import { fr } from "date-fns/locale";
import {
  AlertTriangle,
  ArrowLeft,
  CalendarDays,
  CheckCircle2,
  Download,
  FileDown,
  ImageOff,
  LayoutGrid,
  Loader2,
  Plus,
  Save,
  Trash2,
  Upload,
} from "lucide-react";
import { useEffect, useMemo, useRef, useState } from "react";
import toast from "react-hot-toast";
import Modal from "@/components/Modal";
import { fetchAuthorizedImage, photoboothEventsApi } from "@/lib/api";
import {
  useDeletePhotoboothEvent,
  usePhotoboothCatalog,
  usePhotoboothEvents,
  usePhotoboothFrame,
  useSavePhotoboothEvent,
} from "@/lib/hooks";
import {
  ASPECT_LABELS,
  checkFrame,
  downloadMask,
  downloadTemplate,
  isStrips,
  shotCount,
  templateSvg,
  templatesPdfUrl,
} from "@/lib/photoboothTemplates";
import type {
  PhotoboothAspect,
  PhotoboothEvent,
  PhotoboothEventSettings,
  PhotoboothLayout,
  PhotoboothLayoutCatalog,
} from "@/lib/types";

const ASPECTS: { id: PhotoboothAspect; hint: string }[] = [
  { id: "RATIO_4_3", hint: "Caméra de la tablette" },
  { id: "RATIO_3_2", hint: "Reflex Canon en mode photo" },
  { id: "RATIO_16_9", hint: "Webcam USB / HDMI" },
];

const dayOf = (date: string) => new Date(`${date}T12:00:00`);
const eventDay = (date: string) => format(dayOf(date), "EEEE d MMMM yyyy", { locale: fr });
const todayIso = () => format(new Date(), "yyyy-MM-dd");

const layoutsOf = (catalog: PhotoboothLayoutCatalog | undefined, aspect: PhotoboothAspect) =>
  catalog?.aspects[aspect] ?? [];

function newSettings(catalog: PhotoboothLayoutCatalog | undefined): PhotoboothEventSettings {
  return {
    name: "",
    date: todayIso(),
    photo_aspect: "RATIO_4_3",
    countdown_seconds: 5,
    max_prints_per_session: 5,
    allow_retakes: true,
    digital_share: true,
    attract_slideshow: true,
    mirror_live_preview: true,
    enabled_layouts: layoutsOf(catalog, "RATIO_4_3")
      .filter((l) => l.enabled_by_default)
      .map((l) => l.type),
  };
}

const settingsOf = (e: PhotoboothEvent): PhotoboothEventSettings => ({
  name: e.name,
  date: e.date,
  photo_aspect: e.photo_aspect,
  countdown_seconds: e.countdown_seconds,
  max_prints_per_session: e.max_prints_per_session,
  allow_retakes: e.allow_retakes,
  digital_share: e.digital_share,
  attract_slideshow: e.attract_slideshow,
  mirror_live_preview: e.mirror_live_preview,
  enabled_layouts: e.enabled_layouts,
});

export default function PhotoboothEventsPage() {
  const { data: events = [], isLoading, isError } = usePhotoboothEvents();
  const { data: catalog } = usePhotoboothCatalog();
  const [openId, setOpenId] = useState<number | null>(null);
  const [creating, setCreating] = useState(false);
  const open = events.find((e) => e.id === openId) ?? null;
  const noCatalog = catalog !== undefined && Object.keys(catalog.aspects).length === 0;

  useEffect(() => {
    if (openId !== null && !isLoading && !open) setOpenId(null);
  }, [openId, open, isLoading]);

  if (open) return <EventDetail event={open} catalog={catalog} onBack={() => setOpenId(null)} />;

  return (
    <div className="mx-auto max-w-7xl p-6 lg:p-8">
      <header className="mb-6 flex flex-wrap items-end justify-between gap-4 animate-rise">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-gray-900">Événements photobooth</h1>
          <p className="mt-1 max-w-3xl text-sm text-gray-500">
            Préparez ici les soirées du photobooth : réglages, collages proposés aux invités et
            fonds des collages. Les tablettes récupèrent ces événements et choisissent d'elles-mêmes
            celui du jour.
          </p>
        </div>
        <button
          type="button"
          onClick={() => setCreating(true)}
          className="flex items-center gap-2 rounded-xl bg-blue-600 px-4 py-2.5 text-sm font-semibold text-white shadow transition-colors hover:bg-blue-700"
        >
          <Plus className="h-4 w-4" />
          Nouvel événement
        </button>
      </header>

      {noCatalog && <NoCatalog />}

      {isError && (
        <p className="rounded-2xl border border-red-200 bg-red-50 p-6 text-center text-sm text-red-700">
          Impossible de charger les événements.
        </p>
      )}

      {isLoading && (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {Array.from({ length: 3 }, (_, i) => (
            // biome-ignore lint/suspicious/noArrayIndexKey: static placeholders
            <div key={i} className="h-36 animate-pulse rounded-2xl bg-gray-100" />
          ))}
        </div>
      )}

      {!isLoading && !isError && events.length === 0 && (
        <div className="rounded-2xl border border-dashed border-gray-300 py-16 text-center animate-rise">
          <CalendarDays className="mx-auto h-10 w-10 text-gray-400" />
          <h3 className="mt-3 text-sm font-semibold text-gray-900">Aucun événement</h3>
          <p className="mt-1 text-sm text-gray-500">
            Créez la prochaine soirée pour la retrouver sur les tablettes.
          </p>
        </div>
      )}

      <ul className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {events.map((e, i) => {
          const frames = e.layouts.filter((l) => l.enabled && l.frame).length;
          const today = isToday(dayOf(e.date));
          return (
            <li
              key={e.id}
              className="animate-rise"
              style={{ animationDelay: `${Math.min(i, 12) * 35}ms` }}
            >
              <button
                type="button"
                onClick={() => setOpenId(e.id)}
                className="w-full rounded-2xl border border-gray-200 bg-white p-5 text-left shadow-sm transition-all duration-200 hover:-translate-y-0.5 hover:shadow-md"
              >
                <div className="flex items-start justify-between gap-3">
                  <div className="min-w-0">
                    <h2 className="truncate text-base font-semibold text-gray-900">{e.name}</h2>
                    <p className="mt-0.5 text-sm capitalize text-gray-500">{eventDay(e.date)}</p>
                  </div>
                  {today && (
                    <span className="shrink-0 rounded-full bg-green-50 px-2.5 py-1 text-xs font-medium text-green-700">
                      Aujourd'hui
                    </span>
                  )}
                </div>
                <p className="mt-4 flex items-center gap-1.5 text-xs text-gray-500">
                  <LayoutGrid className="h-3.5 w-3.5" />
                  {e.enabled_layouts.length} collage{e.enabled_layouts.length > 1 ? "s" : ""} ·{" "}
                  {frames} fond{frames > 1 ? "s" : ""} · photos {ASPECT_LABELS[e.photo_aspect]}
                </p>
              </button>
            </li>
          );
        })}
      </ul>

      <SettingsModal
        isOpen={creating}
        onClose={() => setCreating(false)}
        initial={newSettings(catalog)}
        catalog={catalog}
        onSaved={(e) => {
          setCreating(false);
          setOpenId(e.id);
        }}
      />
    </div>
  );
}

function NoCatalog() {
  return (
    <p className="mb-6 flex items-start gap-2 rounded-2xl border border-amber-200 bg-amber-50 p-4 text-sm text-amber-800">
      <AlertTriangle className="mt-0.5 h-4 w-4 shrink-0" />
      Aucune tablette n'a encore envoyé la liste de ses collages. Ouvrez l'application du
      photobooth, connectée au serveur Transat (clé de borne renseignée) : les collages, leurs
      gabarits et l'envoi des fonds apparaîtront ici.
    </p>
  );
}

/** Creation form, and the settings of an existing event. */
function SettingsModal({
  isOpen,
  onClose,
  initial,
  catalog,
  eventId = null,
  onSaved,
}: {
  isOpen: boolean;
  onClose: () => void;
  initial: PhotoboothEventSettings;
  catalog: PhotoboothLayoutCatalog | undefined;
  eventId?: number | null;
  onSaved: (e: PhotoboothEvent) => void;
}) {
  const save = useSavePhotoboothEvent();
  const [draft, setDraft] = useState(initial);
  const [wasOpen, setWasOpen] = useState(false);
  // Fresh values every time the dialog opens.
  if (isOpen !== wasOpen) {
    setWasOpen(isOpen);
    if (isOpen) setDraft(initial);
  }
  const set = <K extends keyof PhotoboothEventSettings>(k: K, v: PhotoboothEventSettings[K]) =>
    setDraft((d) => ({ ...d, [k]: v }));

  const setAspect = (aspect: PhotoboothAspect) =>
    setDraft((d) => {
      const known = new Set(layoutsOf(catalog, aspect).map((l) => l.type));
      // Every ratio offers the same layouts; keep the choice when they are known for it.
      const kept = d.enabled_layouts.filter((t) => known.size === 0 || known.has(t));
      return { ...d, photo_aspect: aspect, enabled_layouts: kept };
    });

  const submit = (ev: React.FormEvent) => {
    ev.preventDefault();
    save.mutate(
      { id: eventId, settings: draft },
      {
        onSuccess: (e) => {
          toast.success(eventId === null ? "Événement créé" : "Réglages enregistrés");
          onSaved(e);
        },
        onError: (e) => toast.error(e.message || "L'enregistrement a échoué"),
      },
    );
  };

  const layouts = layoutsOf(catalog, draft.photo_aspect);
  const input =
    "w-full rounded-xl border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500/20";

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={eventId === null ? "Nouvel événement photobooth" : "Réglages de l'événement"}
      maxWidth="xl"
    >
      <form onSubmit={submit} className="space-y-5 text-left">
        <div className="grid gap-4 sm:grid-cols-[1fr_auto]">
          <label className="block text-sm font-medium text-gray-700">
            Nom
            <input
              className={`${input} mt-1`}
              value={draft.name}
              onChange={(e) => set("name", e.target.value)}
              placeholder="Soirée d'intégration"
              maxLength={80}
              required
            />
          </label>
          <label className="block text-sm font-medium text-gray-700">
            Date
            <input
              type="date"
              className={`${input} mt-1`}
              value={draft.date}
              onChange={(e) => set("date", e.target.value)}
              required
            />
          </label>
        </div>

        <fieldset>
          <legend className="text-sm font-medium text-gray-700">Format des photos</legend>
          <div className="mt-2 grid grid-cols-3 gap-2">
            {ASPECTS.map((a) => (
              <button
                key={a.id}
                type="button"
                onClick={() => setAspect(a.id)}
                aria-pressed={draft.photo_aspect === a.id}
                className={`rounded-xl border px-3 py-2 text-left transition-colors ${
                  draft.photo_aspect === a.id
                    ? "border-blue-500 bg-blue-50 text-blue-800"
                    : "border-gray-200 text-gray-700 hover:bg-gray-50"
                }`}
              >
                <span className="block text-sm font-semibold">{ASPECT_LABELS[a.id]}</span>
                <span className="block text-[11px] text-gray-500">{a.hint}</span>
              </button>
            ))}
          </div>
          <p className="mt-1.5 text-xs text-gray-500">
            Celui de l'appareil utilisé, pour que les photos ne soient jamais recadrées.
          </p>
        </fieldset>

        <div className="grid gap-4 sm:grid-cols-2">
          <label className="block text-sm font-medium text-gray-700">
            Compte à rebours (s)
            <input
              type="number"
              min={1}
              max={15}
              className={`${input} mt-1`}
              value={draft.countdown_seconds}
              onChange={(e) => set("countdown_seconds", Number(e.target.value))}
            />
          </label>
          <label className="block text-sm font-medium text-gray-700">
            Impressions max. par passage
            <input
              type="number"
              min={1}
              max={20}
              className={`${input} mt-1`}
              value={draft.max_prints_per_session}
              onChange={(e) => set("max_prints_per_session", Number(e.target.value))}
            />
          </label>
        </div>

        <div className="space-y-2">
          <Toggle
            label="Autoriser à reprendre une photo"
            checked={draft.allow_retakes}
            onChange={(v) => set("allow_retakes", v)}
          />
          <Toggle
            label="Galerie en ligne et QR code de téléchargement"
            checked={draft.digital_share}
            onChange={(v) => set("digital_share", v)}
          />
          <Toggle
            label="Diaporama des photos sur l'écran d'accueil"
            checked={draft.attract_slideshow}
            onChange={(v) => set("attract_slideshow", v)}
          />
          <Toggle
            label="Aperçu en miroir"
            checked={draft.mirror_live_preview}
            onChange={(v) => set("mirror_live_preview", v)}
          />
        </div>

        {eventId === null && layouts.length > 0 && (
          <fieldset>
            <legend className="text-sm font-medium text-gray-700">Collages proposés</legend>
            <div className="mt-2 flex flex-wrap gap-2">
              {layouts.map((l) => {
                const on = draft.enabled_layouts.includes(l.type);
                return (
                  <button
                    key={l.type}
                    type="button"
                    aria-pressed={on}
                    onClick={() =>
                      set(
                        "enabled_layouts",
                        on
                          ? draft.enabled_layouts.filter((t) => t !== l.type)
                          : [...draft.enabled_layouts, l.type],
                      )
                    }
                    className={`rounded-full border px-3 py-1 text-xs font-medium transition-colors ${
                      on
                        ? "border-blue-500 bg-blue-50 text-blue-800"
                        : "border-gray-200 text-gray-600 hover:bg-gray-50"
                    }`}
                  >
                    {l.title}
                  </button>
                );
              })}
            </div>
            <p className="mt-1.5 text-xs text-gray-500">
              Modifiables ensuite, avec les fonds et les gabarits de chaque collage.
            </p>
          </fieldset>
        )}

        <div className="flex justify-end gap-2 pt-2">
          <button
            type="button"
            onClick={onClose}
            className="rounded-xl border border-gray-300 bg-white px-4 py-2 text-sm text-gray-700 hover:bg-gray-50"
          >
            Annuler
          </button>
          <button
            type="submit"
            disabled={save.isPending || draft.enabled_layouts.length === 0}
            className="flex items-center gap-2 rounded-xl bg-blue-600 px-4 py-2 text-sm font-semibold text-white hover:bg-blue-700 disabled:opacity-50"
          >
            {save.isPending && <Loader2 className="h-4 w-4 animate-spin" />}
            {eventId === null ? "Créer" : "Enregistrer"}
          </button>
        </div>
      </form>
    </Modal>
  );
}

function Toggle({
  label,
  checked,
  onChange,
}: {
  label: string;
  checked: boolean;
  onChange: (v: boolean) => void;
}) {
  return (
    <label className="flex cursor-pointer items-center justify-between gap-4 rounded-xl border border-gray-200 px-3 py-2">
      <span className="text-sm text-gray-700">{label}</span>
      <input
        type="checkbox"
        checked={checked}
        onChange={(e) => onChange(e.target.checked)}
        className="h-4 w-4 accent-blue-600"
      />
    </label>
  );
}

function EventDetail({
  event,
  catalog,
  onBack,
}: {
  event: PhotoboothEvent;
  catalog: PhotoboothLayoutCatalog | undefined;
  onBack: () => void;
}) {
  const save = useSavePhotoboothEvent();
  const remove = useDeletePhotoboothEvent();
  const [editing, setEditing] = useState(false);
  const [confirmDelete, setConfirmDelete] = useState(false);
  const layouts = layoutsOf(catalog, event.photo_aspect);
  const state = useMemo(() => new Map(event.layouts.map((l) => [l.type, l])), [event.layouts]);
  const enabled = layouts.filter((l) => state.get(l.type)?.enabled);
  const missingFrames = enabled.filter((l) => !state.get(l.type)?.frame).length;

  const toggle = (type: string, on: boolean) => {
    const next = on
      ? [...event.enabled_layouts, type]
      : event.enabled_layouts.filter((t) => t !== type);
    if (next.length === 0) {
      toast.error("Gardez au moins un collage");
      return;
    }
    save.mutate(
      { id: event.id, settings: { ...settingsOf(event), enabled_layouts: next } },
      { onError: (e) => toast.error(e.message || "L'enregistrement a échoué") },
    );
  };

  const doDelete = () =>
    remove.mutate(event.id, {
      onSuccess: () => {
        toast.success("Événement supprimé");
        onBack();
      },
      onError: (e) => toast.error(e.message || "La suppression a échoué"),
    });

  const settings = [
    `Photos ${ASPECT_LABELS[event.photo_aspect]}`,
    `compte à rebours ${event.countdown_seconds} s`,
    `${event.max_prints_per_session} impression${event.max_prints_per_session > 1 ? "s" : ""} max.`,
    event.allow_retakes ? "reprises autorisées" : "sans reprise",
    event.digital_share ? "galerie et QR code" : "sans galerie en ligne",
  ];

  return (
    <div className="mx-auto max-w-7xl p-6 lg:p-8">
      <button
        type="button"
        onClick={onBack}
        className="mb-4 flex items-center gap-1.5 text-sm text-gray-500 transition-colors hover:text-gray-900"
      >
        <ArrowLeft className="h-4 w-4" />
        Tous les événements
      </button>

      <header className="mb-6 flex flex-wrap items-end justify-between gap-4 animate-rise">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-gray-900">{event.name}</h1>
          <p className="mt-1 text-sm capitalize text-gray-500">{eventDay(event.date)}</p>
          <p className="mt-1 text-xs text-gray-500">{settings.join(" · ")}</p>
        </div>
        <div className="flex gap-2">
          <button
            type="button"
            onClick={() => setConfirmDelete(true)}
            className="flex items-center gap-2 rounded-xl border border-gray-300 bg-white px-4 py-2.5 text-sm text-gray-700 transition-colors hover:border-red-300 hover:text-red-700"
          >
            <Trash2 className="h-4 w-4" />
            Supprimer
          </button>
          <button
            type="button"
            onClick={() => setEditing(true)}
            className="flex items-center gap-2 rounded-xl bg-blue-600 px-4 py-2.5 text-sm font-semibold text-white shadow transition-colors hover:bg-blue-700"
          >
            <Save className="h-4 w-4" />
            Réglages
          </button>
        </div>
      </header>

      {layouts.length === 0 ? (
        <NoCatalog />
      ) : (
        <section className="mb-6 flex flex-wrap items-center justify-between gap-3 rounded-2xl border border-gray-200 bg-white p-4 text-sm text-gray-600 animate-rise">
          <p>
            Activez les collages proposés aux invités, téléchargez leur gabarit, puis envoyez un
            fond PNG aux dimensions indiquées : les zones roses sont les photos, la zone bleue est
            libre pour le titre ou le logo. Les tablettes l'appliquent à leur prochaine
            synchronisation.
            {missingFrames > 0 && (
              <span className="ml-1 font-medium text-amber-700">
                {missingFrames} collage{missingFrames > 1 ? "s" : ""} activé
                {missingFrames > 1 ? "s" : ""} sans fond.
              </span>
            )}
          </p>
          <a
            href={templatesPdfUrl(event.photo_aspect)}
            target="_blank"
            rel="noreferrer"
            className="flex items-center gap-1.5 font-medium text-blue-600 hover:underline"
          >
            <FileDown className="h-4 w-4" />
            Tous les gabarits en PDF
          </a>
        </section>
      )}

      <ul className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {layouts.map((l, i) => (
          <li
            key={l.type}
            className="animate-rise"
            style={{ animationDelay: `${Math.min(i, 12) * 35}ms` }}
          >
            <LayoutCard
              event={event}
              layout={l}
              enabled={state.get(l.type)?.enabled ?? false}
              frameSha={state.get(l.type)?.frame?.sha256 ?? null}
              busy={save.isPending}
              onToggle={(on) => toggle(l.type, on)}
            />
          </li>
        ))}
      </ul>

      <SettingsModal
        isOpen={editing}
        onClose={() => setEditing(false)}
        initial={settingsOf(event)}
        catalog={catalog}
        eventId={event.id}
        onSaved={() => setEditing(false)}
      />

      <Modal
        isOpen={confirmDelete}
        onClose={() => setConfirmDelete(false)}
        title="Supprimer cet événement ?"
      >
        <p className="text-sm text-gray-600">
          « {event.name} » et ses fonds sont retirés du serveur. Les tablettes qui l'ont déjà
          récupéré le gardent, avec les photos prises.
        </p>
        <div className="mt-6 flex justify-end gap-2">
          <button
            type="button"
            onClick={() => setConfirmDelete(false)}
            className="rounded-xl border border-gray-300 bg-white px-4 py-2 text-sm text-gray-700 hover:bg-gray-50"
          >
            Annuler
          </button>
          <button
            type="button"
            disabled={remove.isPending}
            onClick={doDelete}
            className="flex items-center gap-2 rounded-xl bg-red-600 px-4 py-2 text-sm font-semibold text-white hover:bg-red-700 disabled:opacity-50"
          >
            {remove.isPending && <Loader2 className="h-4 w-4 animate-spin" />}
            Supprimer
          </button>
        </div>
      </Modal>
    </div>
  );
}

function LayoutCard({
  event,
  layout: l,
  enabled,
  frameSha,
  busy,
  onToggle,
}: {
  event: PhotoboothEvent;
  layout: PhotoboothLayout;
  enabled: boolean;
  frameSha: string | null;
  busy: boolean;
  onToggle: (on: boolean) => void;
}) {
  const frame = usePhotoboothFrame();
  const input = useRef<HTMLInputElement>(null);
  const [error, setError] = useState<string | null>(null);
  const shots = shotCount(l);

  const pick = async (file: File | undefined) => {
    if (!file) return;
    setError(null);
    const problem = await checkFrame(file, l);
    if (problem) {
      setError(problem);
      return;
    }
    frame.mutate(
      { id: event.id, type: l.type, file },
      {
        onSuccess: () => toast.success(`Fond de « ${l.title} » envoyé`),
        onError: (e) => setError(e.message || "L'envoi a échoué"),
      },
    );
  };

  const removeFrame = () =>
    frame.mutate(
      { id: event.id, type: l.type, file: null },
      {
        onSuccess: () => toast.success("Fond retiré"),
        onError: (e) => toast.error(e.message || "La suppression a échoué"),
      },
    );

  return (
    <div
      className={`flex h-full flex-col rounded-2xl border bg-white shadow-sm transition-opacity ${
        enabled ? "border-gray-200" : "border-dashed border-gray-300 opacity-70"
      }`}
    >
      <div className="flex items-start justify-between gap-3 px-4 pt-4">
        <div className="min-w-0">
          <h2 className="truncate text-sm font-semibold text-gray-900">{l.title}</h2>
          <p className="text-xs text-gray-500">
            {l.width} × {l.height} px · {shots} photo{shots > 1 ? "s" : ""}
            {isStrips(l) ? ", 2 bandes" : ""}
          </p>
        </div>
        <label className="flex shrink-0 cursor-pointer items-center gap-2 text-xs font-medium text-gray-600">
          {enabled ? "Activé" : "Désactivé"}
          <input
            type="checkbox"
            checked={enabled}
            disabled={busy}
            onChange={(e) => onToggle(e.target.checked)}
            className="h-4 w-4 accent-blue-600"
          />
        </label>
      </div>

      <div className="flex flex-1 items-center justify-center p-4">
        <div className={l.height > l.width ? "w-2/3" : "w-full"}>
          <Preview eventId={event.id} layout={l} frameSha={frameSha} />
        </div>
      </div>

      <div className="flex gap-2 px-4 text-xs">
        <button
          type="button"
          onClick={() => void downloadTemplate(event.name, l)}
          className="flex flex-1 items-center justify-center gap-1.5 rounded-lg border border-gray-200 px-2 py-2 font-medium text-gray-700 hover:bg-gray-50"
        >
          <Download className="h-3.5 w-3.5" /> Gabarit PNG
        </button>
        <button
          type="button"
          onClick={() => void downloadMask(event.name, l)}
          title="Calque blanc avec des trous transparents à la place des photos"
          className="flex flex-1 items-center justify-center gap-1.5 rounded-lg border border-gray-200 px-2 py-2 font-medium text-gray-700 hover:bg-gray-50"
        >
          <Download className="h-3.5 w-3.5" /> Masque PNG
        </button>
      </div>

      <div className="p-4">
        <input
          ref={input}
          type="file"
          accept="image/png"
          className="hidden"
          onChange={(e) => {
            void pick(e.target.files?.[0]);
            e.target.value = "";
          }}
        />
        <div className="flex gap-2">
          <button
            type="button"
            disabled={frame.isPending}
            onClick={() => input.current?.click()}
            className={`flex flex-1 items-center justify-center gap-2 rounded-xl px-4 py-2.5 text-sm font-semibold transition-colors disabled:opacity-60 ${
              frameSha
                ? "border border-gray-300 bg-white text-gray-700 hover:bg-gray-50"
                : "bg-blue-600 text-white hover:bg-blue-700"
            }`}
          >
            {frame.isPending ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : (
              <Upload className="h-4 w-4" />
            )}
            {frameSha ? "Remplacer le fond" : "Envoyer le fond (PNG)"}
          </button>
          {frameSha && (
            <button
              type="button"
              disabled={frame.isPending}
              onClick={removeFrame}
              aria-label={`Retirer le fond de ${l.title}`}
              className="rounded-xl border border-gray-300 bg-white px-3 text-gray-500 hover:border-red-300 hover:text-red-700"
            >
              <Trash2 className="h-4 w-4" />
            </button>
          )}
        </div>
        {frameSha && !error && (
          <p className="mt-2 flex items-center gap-1.5 text-xs text-green-700">
            <CheckCircle2 className="h-3.5 w-3.5" /> Fond en place
          </p>
        )}
        {error && (
          <p role="alert" className="mt-2 flex gap-1.5 text-xs text-red-700">
            <AlertTriangle className="mt-0.5 h-3.5 w-3.5 shrink-0" /> {error}
          </p>
        )}
      </div>
    </div>
  );
}

/** The template, or the frame with the photo zones laid over it as they will be printed. */
function Preview({
  eventId,
  layout: l,
  frameSha,
}: {
  eventId: number;
  layout: PhotoboothLayout;
  frameSha: string | null;
}) {
  const [src, setSrc] = useState<string | null>(null);
  const [failed, setFailed] = useState(false);

  useEffect(() => {
    setSrc(null);
    setFailed(false);
    if (!frameSha) return;
    const controller = new AbortController();
    let url: string | null = null;
    fetchAuthorizedImage(
      photoboothEventsApi.framePath(eventId, l.type, frameSha),
      controller.signal,
    )
      .then((u) => {
        url = u;
        setSrc(u);
      })
      .catch(() => {
        if (!controller.signal.aborted) setFailed(true);
      });
    return () => {
      controller.abort();
      if (url) URL.revokeObjectURL(url);
    };
  }, [eventId, l.type, frameSha]);

  if (!frameSha) {
    return (
      <div
        className="overflow-hidden rounded-md border border-gray-200 shadow-sm [&>svg]:block [&>svg]:h-auto [&>svg]:w-full"
        // biome-ignore lint/security/noDangerouslySetInnerHtml: SVG built locally from numeric geometry
        dangerouslySetInnerHTML={{
          __html: templateSvg(l, false).replace(/ width="\d+" height="\d+"/, ""),
        }}
      />
    );
  }
  if (failed) {
    return (
      <div
        className="flex flex-col items-center justify-center gap-2 rounded-md bg-gray-100 text-gray-400"
        style={{ aspectRatio: `${l.width} / ${l.height}` }}
      >
        <ImageOff className="h-7 w-7" />
        <span className="text-xs">Aperçu indisponible</span>
      </div>
    );
  }
  if (!src) {
    return (
      <div
        className="flex items-center justify-center rounded-md bg-gray-100"
        style={{ aspectRatio: `${l.width} / ${l.height}` }}
      >
        <Loader2 className="h-6 w-6 animate-spin text-gray-400" />
      </div>
    );
  }
  return (
    <svg
      viewBox={`0 0 ${l.width} ${l.height}`}
      className="block h-auto w-full overflow-hidden rounded-md border border-gray-200 shadow-sm"
      role="img"
      aria-label={`Fond de ${l.title}`}
    >
      <rect width={l.width} height={l.height} fill="#ffffff" />
      <image href={src} x={0} y={0} width={l.width} height={l.height} preserveAspectRatio="none" />
      {l.photos.map((z) => (
        <g key={`${z.x}-${z.y}`}>
          <rect x={z.x} y={z.y} width={z.w} height={z.h} fill="rgba(30,21,21,0.55)" />
          <text
            x={z.x + z.w / 2}
            y={z.y + z.h / 2}
            fontSize={Math.min(l.width, l.height) / 14}
            fill="#ffe6cc"
            textAnchor="middle"
            dominantBaseline="middle"
            fontWeight={600}
          >
            {z.shot ?? 1}
          </text>
        </g>
      ))}
    </svg>
  );
}
