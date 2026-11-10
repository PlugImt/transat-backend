"use client";

import {
  Building,
  CheckCircle2,
  Clock,
  ExternalLink,
  GraduationCap,
  ImagePlus,
  Link2,
  Loader2,
  MapPin,
  MousePointerClick,
  Send,
  Smartphone,
  SmartphoneNfc,
  TriangleAlert,
  UserCheck,
  Users,
  X,
} from "lucide-react";
import { useMemo, useState } from "react";
import toast from "react-hot-toast";
import Combobox from "@/components/Combobox";
import { notificationsApi } from "@/lib/api";
import {
  useAppVersions,
  useClubs,
  useEvents,
  useNotificationRecipients,
  useRecentPushMessages,
  useSendNotification,
  useUsers,
} from "@/lib/hooks";
import type {
  AppVersionStat,
  NotificationAudience,
  NotificationCategory,
  NotificationNavigation,
  PushMessageSummary,
} from "@/lib/types";

type AudienceType = NotificationAudience["type"];
type DestinationType = NotificationNavigation["type"] | "";

const MAX_TITLE = 100;
const MAX_MESSAGE = 500;

const AUDIENCES: { value: AudienceType; label: string; hint: string; icon: typeof Users }[] = [
  { value: "all", label: "Tout le monde", hint: "Tous les utilisateurs", icon: Users },
  { value: "club", label: "Un club", hint: "Ses membres", icon: Building },
  { value: "campus", label: "Un campus", hint: "Nantes, Brest, Rennes", icon: MapPin },
  { value: "cohort", label: "Une promo", hint: "Filière et/ou année", icon: GraduationCap },
  { value: "users", label: "Des personnes", hint: "Choisies une à une", icon: UserCheck },
  { value: "app_version", label: "App à jour", hint: "Versions anciennes", icon: SmartphoneNfc },
];

const VERSION = /^[0-9]{1,4}(\.[0-9]{1,4}){0,3}$/;
const MAX_IMAGE_BYTES = 5 * 1024 * 1024;

const CAMPUSES = ["NANTES", "BREST", "RENNES"];

// Matches the formations the backend accepts.
const FORMATIONS = ["FISE", "FIL", "FIP", "FIT", "FID"];

const WEB_LINK = /^https?:\/\/\S+$/;

const CATEGORIES: { value: NotificationCategory; label: string }[] = [
  { value: "RESTAURANT", label: "Restaurant" },
  { value: "EVENTS", label: "Nouveaux événements" },
  { value: "EVENT_REMINDERS", label: "Rappels d'événements" },
  { value: "RESERVATIONS", label: "Réservations" },
  { value: "TRAQ", label: "Traq" },
];

const DESTINATIONS: { value: DestinationType; label: string }[] = [
  { value: "", label: "Ouvrir l'application" },
  { value: "restaurant", label: "Menu du restaurant" },
  { value: "event", label: "Un événement" },
  { value: "club", label: "Un club" },
  { value: "service", label: "Un service" },
  { value: "url", label: "Un lien web" },
];

// Keys understood by the app's `service` destination.
const SERVICES = [
  { value: "laundry", label: "Laverie" },
  { value: "timetable", label: "Emploi du temps" },
  { value: "homework", label: "Devoirs" },
  { value: "clubs", label: "Clubs" },
  { value: "events", label: "Événements" },
  { value: "traq", label: "Traq" },
  { value: "reservations", label: "Réservations" },
  { value: "my_reservations", label: "Mes réservations" },
  { value: "fourchettas", label: "Fourchettas" },
];

const fieldClass =
  "w-full rounded-xl border border-gray-300 bg-white px-3 py-2.5 text-sm text-gray-900 transition-shadow focus:border-blue-500 focus:outline-none focus:ring-4 focus:ring-blue-500/15";

function errorMessage(error: unknown): string {
  const data = (error as { response?: { data?: { error?: string } } }).response?.data;
  return data?.error ?? "Échec de l'envoi de la notification";
}

const sleep = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

// The server only queues the message: look at what the delivery worker did with it, and say so.
async function reportDelivery(id: number) {
  for (let attempt = 0; attempt < 6; attempt++) {
    await sleep(5000);
    try {
      const summary = await notificationsApi.message(id);
      const failed = summary.deliveries.failed ?? 0;
      const taken = (summary.deliveries.accepted ?? 0) + (summary.deliveries.handed_off ?? 0);
      if (failed > 0) {
        const codes = Object.keys(summary.errors ?? {}).join(", ") || "erreur inconnue";
        toast.error(`Échec de la remise à Expo pour ${failed} appareil(s) : ${codes}`, {
          duration: 10000,
        });
        return;
      }
      if (taken > 0 && (summary.deliveries.pending ?? 0) === 0) {
        toast.success(`Remise à Expo pour ${taken} appareil(s)`);
        return;
      }
    } catch {
      return;
    }
  }
}

export default function NotificationsPage() {
  const [title, setTitle] = useState("");
  const [message, setMessage] = useState("");
  const [audienceType, setAudienceType] = useState<AudienceType>("all");
  const [clubId, setClubId] = useState("");
  const [campus, setCampus] = useState(CAMPUSES[0]);
  const [formation, setFormation] = useState("");
  const [graduationYear, setGraduationYear] = useState("");
  const [emails, setEmails] = useState<string[]>([]);
  const [service, setService] = useState<NotificationCategory | "">("");
  const [destinationType, setDestinationType] = useState<DestinationType>("");
  const [destinationId, setDestinationId] = useState("");
  const [confirming, setConfirming] = useState(false);
  const [belowVersion, setBelowVersion] = useState("");
  const [platform, setPlatform] = useState<"" | "ios" | "android">("");
  const [imageUrl, setImageUrl] = useState("");
  const [uploading, setUploading] = useState(false);
  const { data: versions = [] } = useAppVersions(audienceType === "app_version");

  const pickImage = async (file: File | undefined) => {
    if (!file) return;
    if (!file.type.startsWith("image/")) return toast.error("Choisissez une image (JPEG, PNG…)");
    if (file.size > MAX_IMAGE_BYTES) return toast.error("Image trop lourde : 5 Mo au maximum");
    setUploading(true);
    try {
      setImageUrl(await notificationsApi.uploadImage(file));
    } catch (e) {
      toast.error((e as Error).message);
    } finally {
      setUploading(false);
    }
  };

  const { data: clubs = [] } = useClubs();
  const { data: events = [] } = useEvents();
  const { data: users = [] } = useUsers();
  const sendMutation = useSendNotification();

  const audience = useMemo<NotificationAudience | null>(() => {
    switch (audienceType) {
      case "all":
        return { type: "all" };
      case "club":
        return clubId ? { type: "club", clubId: Number(clubId) } : null;
      case "campus":
        return { type: "campus", campus };
      case "users":
        return emails.length > 0 ? { type: "users", emails } : null;
      case "app_version":
        return VERSION.test(belowVersion)
          ? { type: "app_version", belowVersion, platform: platform || undefined }
          : null;
      case "cohort": {
        const year = Number(graduationYear);
        const validYear = Number.isInteger(year) && year >= 2000 && year <= 2100;
        if (graduationYear !== "" && !validYear) return null;
        if (!formation && graduationYear === "") return null;
        return {
          type: "cohort",
          formation: formation || undefined,
          graduationYear: graduationYear === "" ? undefined : year,
        };
      }
    }
  }, [audienceType, clubId, campus, emails, formation, graduationYear, belowVersion, platform]);

  const navigation = useMemo<NotificationNavigation | undefined | null>(() => {
    if (destinationType === "") return undefined;
    if (destinationType === "restaurant") return { type: "restaurant" };
    if (destinationType === "url") {
      const link = destinationId.trim();
      return WEB_LINK.test(link) ? { type: "url", id: link } : null;
    }
    return destinationId ? { type: destinationType, id: destinationId } : null;
  }, [destinationType, destinationId]);

  const {
    data: recipients,
    isFetching: isCounting,
    isError: countFailed,
  } = useNotificationRecipients(audience, service || undefined);

  const clubOptions = clubs.map((club) => ({ value: String(club.id_clubs), label: club.name }));
  const eventOptions = events.map((event) => ({
    value: String(event.id_events),
    label: `${event.name} (${new Date(event.start_date).toLocaleDateString("fr-FR")})`,
  }));
  const userOptions = users
    .filter((user) => !emails.includes(user.email))
    .map((user) => ({
      value: user.email,
      label: `${user.first_name} ${user.last_name} (${user.email})`,
    }));

  const reachable = recipients ? recipients.users - recipients.withoutDevice.length : 0;
  const canSend =
    title.trim() !== "" &&
    audience !== null &&
    navigation !== null &&
    (recipients?.devices ?? 0) > 0 &&
    !sendMutation.isPending;

  const send = async () => {
    if (!audience || navigation === null || !recipients) return;
    try {
      const result = await sendMutation.mutateAsync({
        title: title.trim(),
        message: message.trim(),
        audience,
        service: service || undefined,
        navigation,
        imageUrl: imageUrl || undefined,
      });
      toast.success(
        `Notification mise en file d'attente pour ${result.users - result.withoutDevice.length} utilisateur(s)`,
      );
      if (result.id) void reportDelivery(result.id);
      setTitle("");
      setMessage("");
      setImageUrl("");
      setConfirming(false);
    } catch (error) {
      toast.error(errorMessage(error));
    }
  };

  return (
    <div className="mx-auto max-w-7xl p-6 lg:p-8">
      <header className="mb-8 animate-rise">
        <h1 className="text-2xl font-bold tracking-tight text-gray-900">Notifications</h1>
        <p className="mt-1 text-sm text-gray-500">
          Écrivez le message, choisissez qui le reçoit, vérifiez l'aperçu, envoyez.
        </p>
      </header>

      <div className="grid gap-8 lg:grid-cols-[minmax(0,1fr)_360px]">
        <form
          onSubmit={(e) => {
            e.preventDefault();
            if (canSend) setConfirming(true);
          }}
          className="space-y-6"
        >
          <Step n={1} title="Le message" delay={60}>
            <div className="space-y-4">
              <Field label="Titre" counter={`${title.length}/${MAX_TITLE}`}>
                <input
                  type="text"
                  value={title}
                  maxLength={MAX_TITLE}
                  onChange={(e) => setTitle(e.target.value)}
                  placeholder="Ex. : Le menu de ce midi est en ligne"
                  className={fieldClass}
                  required
                />
              </Field>
              <Field label="Message" counter={`${message.length}/${MAX_MESSAGE}`}>
                <textarea
                  value={message}
                  maxLength={MAX_MESSAGE}
                  rows={4}
                  onChange={(e) => setMessage(e.target.value)}
                  placeholder="Quelques mots de plus (facultatif)"
                  className={`${fieldClass} resize-none`}
                />
              </Field>
              <div>
                <span className="mb-1.5 block text-sm font-medium text-gray-700">
                  Image (facultative)
                </span>
                {imageUrl ? (
                  <div className="flex items-center gap-3 rounded-xl border border-gray-200 p-2 animate-fade-in">
                    {/* biome-ignore lint/performance/noImgElement: remote preview of the uploaded image */}
                    <img src={imageUrl} alt="" className="h-16 w-24 rounded-lg object-cover" />
                    <span className="flex-1 truncate text-xs text-gray-500">{imageUrl}</span>
                    <button
                      type="button"
                      onClick={() => setImageUrl("")}
                      className="rounded-lg p-2 text-gray-500 hover:bg-gray-100"
                      aria-label="Retirer l'image"
                    >
                      <X className="h-4 w-4" />
                    </button>
                  </div>
                ) : (
                  <label
                    className={`flex cursor-pointer items-center justify-center gap-2 rounded-xl border-2 border-dashed border-gray-300 px-4 py-5 text-sm text-gray-500 transition-colors hover:border-blue-400 hover:text-blue-600 ${
                      uploading ? "pointer-events-none opacity-60" : ""
                    }`}
                  >
                    {uploading ? (
                      <Loader2 className="h-4 w-4 animate-spin" />
                    ) : (
                      <ImagePlus className="h-4 w-4" />
                    )}
                    {uploading ? "Envoi de l'image…" : "Ajouter une image (5 Mo max.)"}
                    <input
                      type="file"
                      accept="image/jpeg,image/png,image/webp,image/gif"
                      className="sr-only"
                      onChange={(e) => {
                        void pickImage(e.target.files?.[0]);
                        e.target.value = "";
                      }}
                    />
                  </label>
                )}
                <p className="mt-1 text-xs text-gray-500">
                  Affichée en grand sur Android et en vignette sur iPhone ; le texte part quand même
                  si l'image ne charge pas.
                </p>
              </div>
              <p className="text-xs text-gray-500">
                Le même texte est envoyé quelle que soit la langue de l'application.
              </p>
            </div>
          </Step>

          <Step n={2} title="Les destinataires" delay={120}>
            <div className="grid grid-cols-[repeat(auto-fill,minmax(9rem,1fr))] gap-2">
              {AUDIENCES.map(({ value, label, hint, icon: Icon }) => {
                const active = audienceType === value;
                return (
                  <button
                    key={value}
                    type="button"
                    onClick={() => setAudienceType(value)}
                    aria-pressed={active}
                    className={`flex flex-col items-start gap-2 rounded-2xl border p-3 text-left transition-all duration-200 ${
                      active
                        ? "border-blue-500 bg-blue-50 shadow-sm ring-4 ring-blue-500/10"
                        : "border-gray-200 bg-white hover:-translate-y-0.5 hover:border-gray-300"
                    }`}
                  >
                    <Icon className={`h-5 w-5 ${active ? "text-blue-600" : "text-gray-400"}`} />
                    <span className="text-sm font-semibold text-gray-900">{label}</span>
                    <span className="text-xs leading-tight text-gray-500">{hint}</span>
                  </button>
                );
              })}
            </div>

            <div className="mt-4 space-y-4">
              {audienceType === "club" && (
                <Combobox
                  label="Club"
                  options={clubOptions}
                  value={clubId}
                  onChange={setClubId}
                  placeholder="Choisir un club"
                />
              )}

              {audienceType === "campus" && (
                <div className="flex gap-2">
                  {CAMPUSES.map((name) => (
                    <button
                      key={name}
                      type="button"
                      onClick={() => setCampus(name)}
                      className={`rounded-full px-4 py-2 text-sm font-medium transition-colors ${
                        campus === name
                          ? "bg-blue-600 text-white"
                          : "bg-gray-100 text-gray-700 hover:bg-gray-200"
                      }`}
                    >
                      {name.charAt(0) + name.slice(1).toLowerCase()}
                    </button>
                  ))}
                </div>
              )}

              {audienceType === "cohort" && (
                <div className="grid grid-cols-2 gap-4">
                  <select
                    aria-label="Filière"
                    value={formation}
                    onChange={(e) => setFormation(e.target.value)}
                    className={fieldClass}
                  >
                    <option value="">Toutes les filières</option>
                    {FORMATIONS.map((name) => (
                      <option key={name} value={name}>
                        {name}
                      </option>
                    ))}
                  </select>
                  <input
                    aria-label="Année de diplôme"
                    type="number"
                    min={2000}
                    max={2100}
                    placeholder="Année de diplôme (ex. 2027)"
                    value={graduationYear}
                    onChange={(e) => setGraduationYear(e.target.value)}
                    className={fieldClass}
                  />
                </div>
              )}

              {audienceType === "users" && (
                <div className="space-y-2">
                  <Combobox
                    label="Ajouter une personne"
                    options={userOptions}
                    value=""
                    onChange={(email) => email && setEmails((current) => [...current, email])}
                    placeholder="Rechercher un utilisateur"
                    clearable={false}
                  />
                  <div className="flex flex-wrap gap-2">
                    {emails.map((email) => (
                      <span
                        key={email}
                        className="inline-flex items-center gap-1 rounded-full bg-blue-100 py-1 pl-3 pr-1.5 text-xs text-blue-800 animate-fade-in"
                      >
                        {email}
                        <button
                          type="button"
                          aria-label={`Retirer ${email}`}
                          onClick={() => setEmails((current) => current.filter((e) => e !== email))}
                          className="rounded-full p-0.5 hover:bg-blue-200"
                        >
                          <X className="h-3 w-3" />
                        </button>
                      </span>
                    ))}
                  </div>
                </div>
              )}

              {audienceType === "app_version" && (
                <AppVersionPicker
                  versions={versions}
                  belowVersion={belowVersion}
                  onBelowVersion={setBelowVersion}
                  platform={platform}
                  onPlatform={setPlatform}
                  onTemplate={() => {
                    setTitle("Une mise à jour de Transat est disponible");
                    setMessage(
                      "Mettez l'application à jour pour profiter des dernières nouveautés et corrections.",
                    );
                  }}
                />
              )}

              <Field label="Préférences de notification">
                <select
                  value={service}
                  onChange={(e) => setService(e.target.value as NotificationCategory | "")}
                  className={fieldClass}
                >
                  <option value="">Ignorer : envoyer à tous les appareils du groupe</option>
                  {CATEGORIES.map((option) => (
                    <option key={option.value} value={option.value}>
                      Seulement ceux qui ont activé « {option.label} »
                    </option>
                  ))}
                </select>
              </Field>
            </div>
          </Step>

          <Step n={3} title="Au toucher de la notification" delay={180} optional>
            <div className="flex flex-wrap gap-2">
              {DESTINATIONS.map((option) => (
                <button
                  key={option.value}
                  type="button"
                  onClick={() => {
                    setDestinationType(option.value);
                    setDestinationId("");
                  }}
                  className={`rounded-full px-3.5 py-1.5 text-sm font-medium transition-colors ${
                    destinationType === option.value
                      ? "bg-blue-600 text-white"
                      : "bg-gray-100 text-gray-700 hover:bg-gray-200"
                  }`}
                >
                  {option.label}
                </button>
              ))}
            </div>
            <div className="mt-4">
              {destinationType === "event" && (
                <Combobox
                  label="Événement"
                  options={eventOptions}
                  value={destinationId}
                  onChange={setDestinationId}
                  placeholder="Choisir un événement"
                />
              )}
              {destinationType === "club" && (
                <Combobox
                  label="Club"
                  options={clubOptions}
                  value={destinationId}
                  onChange={setDestinationId}
                  placeholder="Choisir un club"
                />
              )}
              {destinationType === "service" && (
                <select
                  aria-label="Service"
                  value={destinationId}
                  onChange={(e) => setDestinationId(e.target.value)}
                  className={fieldClass}
                >
                  <option value="">Choisir un service</option>
                  {SERVICES.map((option) => (
                    <option key={option.value} value={option.value}>
                      {option.label}
                    </option>
                  ))}
                </select>
              )}
              {destinationType === "url" && (
                <div>
                  <input
                    aria-label="Lien"
                    type="url"
                    placeholder="https://exemple.fr/page"
                    value={destinationId}
                    onChange={(e) => setDestinationId(e.target.value)}
                    className={fieldClass}
                  />
                  {destinationId !== "" && navigation === null && (
                    <p className="mt-1 text-xs text-red-600">
                      Le lien doit commencer par http:// ou https://
                    </p>
                  )}
                </div>
              )}
            </div>
          </Step>
        </form>

        <aside className="space-y-4 lg:sticky lg:top-8 lg:self-start">
          <PhonePreview
            title={title}
            message={message}
            imageUrl={imageUrl}
            destination={DESTINATIONS.find((d) => d.value === destinationType)?.label}
          />

          <div
            className="rounded-2xl border border-gray-200 bg-white p-5 animate-rise"
            style={{ animationDelay: "120ms" }}
          >
            <p className="text-xs font-semibold uppercase tracking-wider text-gray-500">Portée</p>
            <div aria-live="polite" className="mt-3">
              {audience === null ? (
                <p className="text-sm text-gray-500">Choisissez les destinataires.</p>
              ) : countFailed ? (
                <p className="text-sm text-red-600">Impossible de calculer les destinataires.</p>
              ) : (
                <>
                  <div
                    className={`grid grid-cols-2 gap-3 transition-opacity ${isCounting ? "opacity-50" : ""}`}
                  >
                    <Stat icon={Users} value={recipients ? reachable : "…"} label="personnes" />
                    <Stat icon={Smartphone} value={recipients?.devices ?? "…"} label="appareils" />
                  </div>
                  {recipients && recipients.withoutDevice.length > 0 && (
                    <UnreachableList
                      names={recipients.withoutDeviceNames ?? []}
                      total={recipients.withoutDevice.length}
                    />
                  )}
                </>
              )}
            </div>

            {confirming ? (
              <div className="mt-5 space-y-3 rounded-xl bg-blue-50 p-3 animate-fade-in">
                <p className="text-sm text-gray-900">
                  Envoyer à <strong>{reachable}</strong> personne(s),{" "}
                  <strong>{recipients?.devices}</strong> appareil(s) ?
                </p>
                <div className="flex gap-2">
                  <button
                    type="button"
                    onClick={send}
                    disabled={sendMutation.isPending}
                    className="inline-flex flex-1 items-center justify-center gap-2 rounded-xl bg-blue-600 px-4 py-2.5 text-sm font-semibold text-white hover:bg-blue-700 disabled:opacity-50"
                  >
                    <Send className="h-4 w-4" />
                    {sendMutation.isPending ? "Envoi…" : "Confirmer"}
                  </button>
                  <button
                    type="button"
                    onClick={() => setConfirming(false)}
                    className="rounded-xl border border-gray-300 bg-white px-3 py-2.5 text-sm"
                  >
                    Annuler
                  </button>
                </div>
              </div>
            ) : (
              <button
                type="button"
                onClick={() => setConfirming(true)}
                disabled={!canSend}
                className="mt-5 inline-flex w-full items-center justify-center gap-2 rounded-xl bg-gradient-to-r from-blue-500 to-blue-700 px-4 py-3 text-sm font-semibold text-white shadow-lg shadow-blue-600/20 transition-all hover:-translate-y-0.5 hover:shadow-xl active:translate-y-0 disabled:translate-y-0 disabled:opacity-40 disabled:shadow-none"
              >
                <Send className="h-4 w-4" />
                Envoyer la notification
              </button>
            )}
          </div>

          <RecentSends />
        </aside>
      </div>
    </div>
  );
}

function Step({
  n,
  title,
  delay,
  optional,
  children,
}: {
  n: number;
  title: string;
  delay: number;
  optional?: boolean;
  children: React.ReactNode;
}) {
  return (
    <section
      className="rounded-2xl border border-gray-200 bg-white p-6 animate-rise"
      style={{ animationDelay: `${delay}ms` }}
    >
      <div className="mb-5 flex items-center gap-3">
        <span className="flex h-7 w-7 items-center justify-center rounded-full bg-blue-600 text-xs font-bold text-white">
          {n}
        </span>
        <h2 className="text-base font-semibold text-gray-900">{title}</h2>
        {optional && <span className="text-xs text-gray-400">facultatif</span>}
      </div>
      {children}
    </section>
  );
}

function Field({
  label,
  counter,
  children,
}: {
  label: string;
  counter?: string;
  children: React.ReactNode;
}) {
  return (
    // biome-ignore lint/a11y/noLabelWithoutControl: the control is the child
    <label className="block">
      <span className="mb-1.5 flex items-center justify-between text-sm font-medium text-gray-700">
        {label}
        {counter && <span className="text-xs font-normal text-gray-400">{counter}</span>}
      </span>
      {children}
    </label>
  );
}

function Stat({
  icon: Icon,
  value,
  label,
}: {
  icon: typeof Users;
  value: number | string;
  label: string;
}) {
  return (
    <div className="rounded-xl bg-gray-50 p-3">
      <Icon className="h-4 w-4 text-gray-400" />
      <p className="mt-1 text-2xl font-bold tabular-nums text-gray-900">{value}</p>
      <p className="text-xs text-gray-500">{label}</p>
    </div>
  );
}

/** How the notification shows up on a phone's lock screen, as it is typed. */
function PhonePreview({
  title,
  message,
  imageUrl,
  destination,
}: {
  title: string;
  message: string;
  imageUrl?: string;
  destination?: string;
}) {
  return (
    <div className="animate-rise overflow-hidden rounded-[2rem] border border-gray-200 bg-gradient-to-br from-[#2a1712] via-[#120a08] to-[#0d0505] p-4 shadow-xl">
      <div className="mb-8 mt-2 text-center text-[#ffe6cc]">
        <p className="text-xs opacity-70">
          {new Date().toLocaleDateString("fr-FR", {
            weekday: "long",
            day: "numeric",
            month: "long",
          })}
        </p>
        <p className="text-5xl font-semibold tracking-tight tabular-nums">
          {new Date().toLocaleTimeString("fr-FR", { hour: "2-digit", minute: "2-digit" })}
        </p>
      </div>
      <div className="rounded-2xl bg-white/12 p-3 text-white backdrop-blur-md ring-1 ring-white/10 transition-all">
        <div className="mb-1 flex items-center gap-2 text-[11px] uppercase tracking-wide text-white/60">
          {/* biome-ignore lint/performance/noImgElement: tiny static icon */}
          <img src="/transat-admin.png" alt="" className="h-4 w-4 rounded" />
          Transat
          <span className="ml-auto normal-case">maintenant</span>
        </div>
        <p className="break-words text-sm font-semibold">{title || "Titre de la notification"}</p>
        <p className="line-clamp-4 break-words text-sm text-white/80">
          {message || "Le message apparaît ici."}
        </p>
        {imageUrl && (
          // biome-ignore lint/performance/noImgElement: remote preview of the uploaded image
          <img
            src={imageUrl}
            alt=""
            className="mt-2 max-h-40 w-full rounded-xl object-cover animate-fade-in"
          />
        )}
      </div>
      <p className="mt-3 flex items-center justify-center gap-1.5 text-[11px] text-white/50">
        <MousePointerClick className="h-3 w-3" />
        {destination ?? "Ouvrir l'application"}
      </p>
    </div>
  );
}

const STATUS_STYLE: Record<string, { label: string; icon: typeof Clock; tone: string }> = {
  completed: { label: "Envoyée", icon: CheckCircle2, tone: "text-green-600" },
  done: { label: "Envoyée", icon: CheckCircle2, tone: "text-green-600" },
  failed: { label: "Échec", icon: TriangleAlert, tone: "text-red-600" },
};

function RecentSends() {
  const { data, isError } = useRecentPushMessages();
  if (isError || !data || data.length === 0) return null;
  return (
    <div
      className="rounded-2xl border border-gray-200 bg-white p-5 animate-rise"
      style={{ animationDelay: "200ms" }}
    >
      <p className="mb-3 text-xs font-semibold uppercase tracking-wider text-gray-500">
        Derniers envois
      </p>
      <ul className="space-y-3">
        {data.map((m) => (
          <RecentRow key={m.id} message={m} />
        ))}
      </ul>
    </div>
  );
}

function RecentRow({ message: m }: { message: PushMessageSummary }) {
  const delivered = (m.deliveries.accepted ?? 0) + (m.deliveries.handed_off ?? 0);
  const failed = m.deliveries.failed ?? 0;
  const total = Object.values(m.deliveries).reduce((a, b) => a + b, 0);
  const style = STATUS_STYLE[m.status] ?? {
    label: "En cours",
    icon: Clock,
    tone: "text-amber-600",
  };
  const Icon = style.icon;
  return (
    <li className="text-sm">
      <div className="flex items-center justify-between gap-2">
        <span className={`inline-flex items-center gap-1.5 font-medium ${style.tone}`}>
          <Icon className="h-3.5 w-3.5" />
          {style.label}
        </span>
        <span className="text-xs text-gray-500">
          {m.createdAt
            ? new Date(m.createdAt).toLocaleString("fr-FR", {
                dateStyle: "short",
                timeStyle: "short",
              })
            : `#${m.id}`}
        </span>
      </div>
      <div className="mt-1.5 flex h-1.5 overflow-hidden rounded-full bg-gray-100">
        {total > 0 && (
          <>
            <div
              className="bg-green-500 transition-all"
              style={{ width: `${(delivered / total) * 100}%` }}
            />
            <div
              className="bg-red-500 transition-all"
              style={{ width: `${(failed / total) * 100}%` }}
            />
          </>
        )}
      </div>
      <p className="mt-1 flex items-center gap-1 text-xs text-gray-500">
        {m.kind === "admin" ? <Link2 className="h-3 w-3" /> : <ExternalLink className="h-3 w-3" />}
        {delivered}/{total} appareil(s){m.category ? ` · ${m.category}` : ""}
        {m.createdBy ? ` · ${m.createdBy}` : ""}
      </p>
    </li>
  );
}

/** People of the audience no push can reach, by name. */
function UnreachableList({ names, total }: { names: string[]; total: number }) {
  const [all, setAll] = useState(false);
  const shown = all ? names : names.slice(0, 12);
  return (
    <details className="mt-3 text-xs">
      <summary className="cursor-pointer text-amber-700">
        {total} personne{total > 1 ? "s" : ""} sans notifications (injoignable
        {total > 1 ? "s" : ""})
      </summary>
      <p className="mt-1 text-gray-500">
        Notifications refusées, ou application jamais ouverte sur un téléphone.
      </p>
      <ul className="mt-2 flex flex-wrap gap-1.5">
        {shown.map((name, i) => (
          // biome-ignore lint/suspicious/noArrayIndexKey: two people may share a name
          <li key={`${name}-${i}`} className="rounded-full bg-gray-100 px-2 py-0.5 text-gray-700">
            {name}
          </li>
        ))}
      </ul>
      {names.length > shown.length && (
        <button
          type="button"
          onClick={() => setAll(true)}
          className="mt-2 text-blue-600 hover:text-blue-800"
        >
          Voir les {names.length - shown.length} autres
        </button>
      )}
    </details>
  );
}

const PLATFORMS = [
  ["", "Tous"],
  ["ios", "iPhone"],
  ["android", "Android"],
] as const;

/** Which old versions to target, with the versions in use for reference. */
function AppVersionPicker({
  versions,
  belowVersion,
  onBelowVersion,
  platform,
  onPlatform,
  onTemplate,
}: {
  versions: AppVersionStat[];
  belowVersion: string;
  onBelowVersion: (v: string) => void;
  platform: "" | "ios" | "android";
  onPlatform: (p: "" | "ios" | "android") => void;
  onTemplate: () => void;
}) {
  const latest = versions.find((v) => VERSION.test(v.version))?.version;
  const total = versions.reduce((n, v) => n + v.devices, 0);
  const valid = VERSION.test(belowVersion);
  return (
    <div className="space-y-4 rounded-2xl bg-gray-50 p-4 animate-fade-in">
      <div className="flex flex-wrap items-end gap-3">
        <label className="block">
          <span className="mb-1.5 block text-sm font-medium text-gray-700">
            Versions antérieures à
          </span>
          <input
            value={belowVersion}
            onChange={(e) => onBelowVersion(e.target.value.trim())}
            placeholder={latest ?? "2.4.0"}
            className="w-36 rounded-xl border border-gray-300 bg-white px-3 py-2.5 text-sm"
          />
        </label>
        {latest && belowVersion !== latest && (
          <button
            type="button"
            onClick={() => onBelowVersion(latest)}
            className="rounded-xl border border-gray-300 bg-white px-3 py-2.5 text-sm hover:bg-gray-100"
          >
            Tout sauf la dernière ({latest})
          </button>
        )}
        <div className="flex gap-1 rounded-xl bg-white p-1 ring-1 ring-gray-200">
          {PLATFORMS.map(([value, label]) => (
            <button
              key={value}
              type="button"
              onClick={() => onPlatform(value)}
              className={`rounded-lg px-3 py-1.5 text-sm ${
                platform === value ? "bg-blue-600 text-white" : "text-gray-600 hover:bg-gray-100"
              }`}
            >
              {label}
            </button>
          ))}
        </div>
        <button
          type="button"
          onClick={onTemplate}
          className="text-sm text-blue-600 hover:text-blue-800"
        >
          Message type « mise à jour »
        </button>
      </div>
      {belowVersion !== "" && !valid && (
        <p className="text-xs text-red-600">Une version s'écrit en chiffres, par exemple 2.4.0.</p>
      )}
      {versions.length > 0 && (
        <div>
          <p className="mb-2 text-xs font-semibold uppercase tracking-wider text-gray-500">
            Versions utilisées
          </p>
          <ul className="space-y-1.5">
            {versions.map((v) => {
              const targeted =
                valid &&
                (v.version === "inconnue" || compareVersions(v.version, belowVersion) < 0) &&
                (!platform || v.platform === platform);
              return (
                // On a narrow screen the bar takes its own line under the version.
                <li
                  key={`${v.version}-${v.platform}`}
                  className="flex flex-wrap items-center gap-x-3 gap-y-1 text-sm"
                >
                  <span
                    className={`w-20 shrink-0 font-mono ${targeted ? "font-semibold text-blue-700" : "text-gray-700"}`}
                  >
                    {v.version}
                  </span>
                  <span className="w-16 shrink-0 text-xs text-gray-500">
                    {v.platform === "ios"
                      ? "iPhone"
                      : v.platform === "android"
                        ? "Android"
                        : v.platform}
                  </span>
                  <div className="order-last h-2 basis-full overflow-hidden rounded-full bg-gray-200 sm:order-none sm:basis-0 sm:flex-1">
                    <div
                      className={`h-full rounded-full transition-all ${targeted ? "bg-blue-500" : "bg-gray-400"}`}
                      style={{ width: `${total ? (v.devices / total) * 100 : 0}%` }}
                    />
                  </div>
                  <span className="ml-auto text-right text-xs tabular-nums text-gray-600 sm:w-24">
                    {v.devices} appareil{v.devices > 1 ? "s" : ""}
                  </span>
                </li>
              );
            })}
          </ul>
          <p className="mt-2 text-xs text-gray-500">
            Seuls les appareils concernés reçoivent la notification. « inconnue » : appareils
            enregistrés avant le suivi des versions.
          </p>
        </div>
      )}
    </div>
  );
}

/** Numeric comparison of "2.10.0" and "2.9.1"; extra text such as "-beta" is ignored, like on the server. */
function compareVersions(a: string, b: string) {
  const parse = (v: string) => (v.match(/^[0-9]+(\.[0-9]+)*/)?.[0] ?? "0").split(".").map(Number);
  const x = parse(a);
  const y = parse(b);
  for (let i = 0; i < Math.max(x.length, y.length); i++) {
    const d = (x[i] ?? 0) - (y[i] ?? 0);
    if (d !== 0) return d;
  }
  return 0;
}
