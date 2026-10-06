"use client";

import { Bell, Send, X } from "lucide-react";
import { useMemo, useState } from "react";
import toast from "react-hot-toast";
import Combobox from "@/components/Combobox";
import {
  useClubs,
  useEvents,
  useNotificationRecipients,
  useSendNotification,
  useUsers,
} from "@/lib/hooks";
import type {
  NotificationAudience,
  NotificationCategory,
  NotificationNavigation,
} from "@/lib/types";

type AudienceType = NotificationAudience["type"];
type DestinationType = NotificationNavigation["type"] | "";

const MAX_TITLE = 100;
const MAX_MESSAGE = 500;

const AUDIENCES: { value: AudienceType; label: string }[] = [
  { value: "all", label: "Tous les utilisateurs" },
  { value: "club", label: "Membres d'un club" },
  { value: "campus", label: "Un campus" },
  { value: "users", label: "Utilisateurs précis" },
];

const CAMPUSES = ["NANTES", "BREST", "RENNES"];

const CATEGORIES: { value: NotificationCategory; label: string }[] = [
  { value: "RESTAURANT", label: "Restaurant" },
  { value: "EVENTS", label: "Nouveaux événements" },
  { value: "EVENT_REMINDERS", label: "Rappels d'événements" },
  { value: "RESERVATIONS", label: "Réservations" },
  { value: "TRAQ", label: "Traq" },
];

const DESTINATIONS: { value: DestinationType; label: string }[] = [
  { value: "", label: "Aucune (ouvre simplement l'application)" },
  { value: "restaurant", label: "Menu du restaurant" },
  { value: "event", label: "Un événement" },
  { value: "club", label: "Un club" },
  { value: "service", label: "Un service" },
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
  "w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500";

function errorMessage(error: unknown): string {
  const data = (error as { response?: { data?: { error?: string } } }).response?.data;
  return data?.error ?? "Échec de l'envoi de la notification";
}

export default function NotificationsPage() {
  const [title, setTitle] = useState("");
  const [message, setMessage] = useState("");
  const [audienceType, setAudienceType] = useState<AudienceType>("all");
  const [clubId, setClubId] = useState("");
  const [campus, setCampus] = useState(CAMPUSES[0]);
  const [emails, setEmails] = useState<string[]>([]);
  const [service, setService] = useState<NotificationCategory | "">("");
  const [destinationType, setDestinationType] = useState<DestinationType>("");
  const [destinationId, setDestinationId] = useState("");

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
    }
  }, [audienceType, clubId, campus, emails]);

  const navigation = useMemo<NotificationNavigation | undefined | null>(() => {
    if (destinationType === "") return undefined;
    if (destinationType === "restaurant") return { type: "restaurant" };
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

  const canSend =
    title.trim() !== "" &&
    audience !== null &&
    navigation !== null &&
    (recipients?.devices ?? 0) > 0 &&
    !sendMutation.isPending;

  const handleDestinationTypeChange = (value: DestinationType) => {
    setDestinationType(value);
    setDestinationId("");
  };

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();
    if (!audience || navigation === null || !recipients) return;

    if (
      !confirm(
        `Envoyer cette notification à ${recipients.users - recipients.withoutDevice.length} utilisateur(s) (${recipients.devices} appareil(s)) ?`,
      )
    ) {
      return;
    }

    try {
      const result = await sendMutation.mutateAsync({
        title: title.trim(),
        message: message.trim(),
        audience,
        service: service || undefined,
        navigation,
      });
      toast.success(
        `Notification envoyée à ${result.users - result.withoutDevice.length} utilisateur(s)`,
      );
      setTitle("");
      setMessage("");
    } catch (error) {
      toast.error(errorMessage(error));
    }
  };

  return (
    <div className="p-6 max-w-3xl">
      <div className="flex items-center mb-6">
        <Bell className="h-6 w-6 text-gray-700 mr-3" />
        <h1 className="text-2xl font-bold text-gray-900">Notifications</h1>
      </div>

      <form onSubmit={handleSubmit} className="bg-white shadow rounded-lg p-6 space-y-6">
        <section className="space-y-4">
          <div>
            <label htmlFor="notif-title" className="block text-sm font-medium text-gray-700 mb-1">
              Titre
            </label>
            <input
              id="notif-title"
              type="text"
              value={title}
              maxLength={MAX_TITLE}
              onChange={(e) => setTitle(e.target.value)}
              className={fieldClass}
              required
            />
          </div>
          <div>
            <label htmlFor="notif-message" className="block text-sm font-medium text-gray-700 mb-1">
              Message
            </label>
            <textarea
              id="notif-message"
              value={message}
              maxLength={MAX_MESSAGE}
              rows={4}
              onChange={(e) => setMessage(e.target.value)}
              className={fieldClass}
            />
            <p className="text-xs text-gray-500 mt-1">
              {message.length}/{MAX_MESSAGE} — le même texte est envoyé dans toutes les langues.
            </p>
          </div>
        </section>

        <section className="space-y-4">
          <h2 className="text-lg font-semibold text-gray-900">Destinataires</h2>
          <div>
            <label
              htmlFor="notif-audience"
              className="block text-sm font-medium text-gray-700 mb-1"
            >
              Groupe
            </label>
            <select
              id="notif-audience"
              value={audienceType}
              onChange={(e) => setAudienceType(e.target.value as AudienceType)}
              className={fieldClass}
            >
              {AUDIENCES.map((option) => (
                <option key={option.value} value={option.value}>
                  {option.label}
                </option>
              ))}
            </select>
          </div>

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
            <select
              aria-label="Campus"
              value={campus}
              onChange={(e) => setCampus(e.target.value)}
              className={fieldClass}
            >
              {CAMPUSES.map((name) => (
                <option key={name} value={name}>
                  {name}
                </option>
              ))}
            </select>
          )}

          {audienceType === "users" && (
            <div className="space-y-2">
              <Combobox
                label="Ajouter un utilisateur"
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
                    className="inline-flex items-center px-2 py-1 rounded-full text-xs bg-blue-100 text-blue-800"
                  >
                    {email}
                    <button
                      type="button"
                      aria-label={`Retirer ${email}`}
                      onClick={() => setEmails((current) => current.filter((e) => e !== email))}
                      className="ml-1"
                    >
                      <X className="h-3 w-3" />
                    </button>
                  </span>
                ))}
              </div>
            </div>
          )}

          <div>
            <label htmlFor="notif-service" className="block text-sm font-medium text-gray-700 mb-1">
              Respecter les préférences
            </label>
            <select
              id="notif-service"
              value={service}
              onChange={(e) => setService(e.target.value as NotificationCategory | "")}
              className={fieldClass}
            >
              <option value="">Non : envoyer à tous les appareils du groupe</option>
              {CATEGORIES.map((option) => (
                <option key={option.value} value={option.value}>
                  Seulement ceux qui ont activé : {option.label}
                </option>
              ))}
            </select>
          </div>

          <div className="text-sm text-gray-600 space-y-1" aria-live="polite">
            {audience === null ? (
              <p>Sélectionnez les destinataires.</p>
            ) : isCounting ? (
              <p>Calcul des destinataires…</p>
            ) : countFailed ? (
              <p className="text-red-600">Impossible de calculer les destinataires.</p>
            ) : (
              recipients && (
                <>
                  <p>
                    {recipients.users} utilisateur(s), {recipients.devices} appareil(s)
                    enregistré(s)
                  </p>
                  {recipients.withoutDevice.length > 0 && (
                    <p className="text-amber-700">
                      Sans appareil enregistré (n'a jamais autorisé les notifications ou ouvert
                      l'application sur un téléphone), donc injoignable :{" "}
                      {recipients.withoutDevice.slice(0, 10).join(", ")}
                      {recipients.withoutDevice.length > 10 &&
                        ` et ${recipients.withoutDevice.length - 10} autre(s)`}
                    </p>
                  )}
                </>
              )
            )}
          </div>
        </section>

        <section className="space-y-4">
          <h2 className="text-lg font-semibold text-gray-900">Page ouverte au clic (optionnel)</h2>
          <select
            aria-label="Destination"
            value={destinationType}
            onChange={(e) => handleDestinationTypeChange(e.target.value as DestinationType)}
            className={fieldClass}
          >
            {DESTINATIONS.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>

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
        </section>

        <button
          type="submit"
          disabled={!canSend}
          className="inline-flex items-center px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <Send className="h-4 w-4 mr-2" />
          {sendMutation.isPending ? "Envoi…" : "Envoyer la notification"}
        </button>
      </form>
    </div>
  );
}
