import { format } from "date-fns";
import type { PhotoReportDeliveryStatus, PhotoReportReason } from "./types";

export const REASON_LABELS: Record<PhotoReportReason, string> = {
  appears_in_photo: "Retrait demandé (apparaît sur la photo)",
  inappropriate: "Contenu inapproprié ou NSFW",
  poor_quality: "Photo de mauvaise qualité",
  other: "Autre",
};

export const DELIVERY_STATUS_LABELS: Record<PhotoReportDeliveryStatus, string> = {
  sent: "Envoyé",
  pending: "En cours / nouvelle tentative",
  failed: "Échec",
  no_recipients: "Aucun destinataire",
};

export const formatDateTime = (value: string) => format(new Date(value), "dd/MM/yyyy HH:mm");

// Only web links become clickable: the value comes from the database, not from the page.
export const safeHttpUrl = (value: string): string | null =>
  /^https?:\/\//i.test(value) ? value : null;
