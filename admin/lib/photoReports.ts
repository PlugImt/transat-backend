import { format } from "date-fns";
import type { PhotoReportDeliveryStatus, PhotoReportReason, PhotoReportStatus } from "./types";

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

export const STATUS_LABELS: Record<PhotoReportStatus, string> = {
  pending: "À traiter",
  approved: "Approuvés",
  rejected: "Refusés",
};

/** Drive serves a resized copy of the photo; the gallery of the app uses the same address. */
export const photoPreviewUrl = (driveFileId: string, width: number) =>
  `https://lh3.googleusercontent.com/d/${encodeURIComponent(driveFileId)}=w${width}`;

export const photoPreviewFallbackUrl = (driveFileId: string, width: number) =>
  `https://drive.google.com/thumbnail?id=${encodeURIComponent(driveFileId)}&sz=w${width}`;

export const STATUS_BADGE: Record<PhotoReportStatus, string> = {
  pending: "bg-amber-100 text-amber-800",
  approved: "bg-green-100 text-green-800",
  rejected: "bg-gray-100 text-gray-700",
};
