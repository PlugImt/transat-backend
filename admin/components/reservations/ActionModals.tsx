"use client";

import { AlertTriangle, Mail, Smartphone, Users } from "lucide-react";
import toast from "react-hot-toast";
import Modal from "@/components/Modal";
import { useAdminCancelReservation, useMarkItemReturned } from "@/lib/reservationHooks";
import { useAuthStore } from "@/lib/stores/authStore";
import { apiErrorMessage, formatDateTime, LoadingButton } from "./ui";

export interface ReservationTarget {
  itemId: number;
  itemName: string;
  slot: boolean;
  /** ISO start of the slot; the loan of an item without slot needs none. */
  startDate?: string;
  endDate?: string;
  /** Every slot start of a merged event; cancelled one after the other. */
  startDates?: string[];
  ownerName: string;
  ownerEmail: string;
}

interface ActionModalProps {
  target: ReservationTarget | null;
  onClose: () => void;
}

export function CancelReservationModal({ target, onClose }: ActionModalProps) {
  const cancel = useAdminCancelReservation();
  const admin = useAuthStore((state) => state.user);
  const adminName = `${admin?.firstName ?? ""} ${admin?.lastName ?? ""}`.trim() || "vous";

  const confirm = async () => {
    if (!target) return;
    try {
      const starts = target.slot ? (target.startDates ?? [target.startDate]) : [undefined];
      let waitlistNotified = false;
      for (const start of starts) {
        const result = await cancel.mutateAsync({ item_id: target.itemId, start_date: start });
        waitlistNotified ||= result.waitlist_notified;
      }
      toast.success(
        waitlistNotified
          ? `Réservation annulée. ${target.ownerName} a été prévenu(e), la première personne de la file d'attente aussi.`
          : `Réservation annulée. ${target.ownerName} a été prévenu(e).`,
      );
      onClose();
    } catch (error) {
      toast.error(apiErrorMessage(error, "Impossible d'annuler la réservation"));
    }
  };

  return (
    <Modal
      isOpen={target !== null}
      onClose={cancel.isPending ? () => {} : onClose}
      title="Annuler la réservation"
      maxWidth="lg"
    >
      {target && (
        <div className="space-y-4">
          <div className="rounded-md border border-gray-200 bg-gray-50 p-3 text-sm">
            <p className="font-medium text-gray-900">{target.itemName}</p>
            <p className="text-gray-600">
              {target.ownerName} · {target.ownerEmail}
            </p>
            {target.slot && target.startDate && (
              <p className="text-gray-600">
                {formatDateTime(target.startDate)}
                {target.endDate ? ` → ${formatDateTime(target.endDate)}` : ""}
              </p>
            )}
          </div>

          <ul className="space-y-2 text-sm text-gray-700">
            <li className="flex gap-2">
              <Mail className="mt-0.5 h-4 w-4 shrink-0 text-blue-600" />
              <span>
                {target.ownerName} recevra automatiquement un <strong>e-mail</strong> (au thème
                Transat) indiquant que <strong>{adminName}</strong> a annulé sa réservation.
              </span>
            </li>
            <li className="flex gap-2">
              <Smartphone className="mt-0.5 h-4 w-4 shrink-0 text-blue-600" />
              <span>
                Une <strong>notification push</strong> mentionnant votre nom lui sera aussi envoyée.
              </span>
            </li>
            <li className="flex gap-2">
              <Users className="mt-0.5 h-4 w-4 shrink-0 text-blue-600" />
              <span>
                {target.slot ? "Le créneau" : "L'élément"} sera libéré et, s'il y a une file
                d'attente, sa première personne sera prévenue.
              </span>
            </li>
          </ul>

          <div className="flex justify-end gap-2 pt-2">
            <LoadingButton
              onClick={onClose}
              disabled={cancel.isPending}
              className="border border-gray-300 text-gray-700 hover:bg-gray-50"
            >
              Conserver
            </LoadingButton>
            <LoadingButton
              onClick={confirm}
              loading={cancel.isPending}
              className="bg-red-600 text-white hover:bg-red-700"
            >
              Annuler la réservation
            </LoadingButton>
          </div>
        </div>
      )}
    </Modal>
  );
}

export function MarkReturnedModal({ target, onClose }: ActionModalProps) {
  const markReturned = useMarkItemReturned();

  const confirm = async () => {
    if (!target) return;
    try {
      const result = await markReturned.mutateAsync(target.itemId);
      toast.success(
        result.waitlist_notified
          ? `${target.itemName} est de nouveau disponible. La première personne de la file a été prévenue.`
          : `${target.itemName} est de nouveau disponible.`,
      );
      onClose();
    } catch (error) {
      toast.error(apiErrorMessage(error, "Impossible de marquer cet élément comme rendu"));
    }
  };

  return (
    <Modal
      isOpen={target !== null}
      onClose={markReturned.isPending ? () => {} : onClose}
      title="Marquer comme rendu"
      maxWidth="lg"
    >
      {target && (
        <div className="space-y-4">
          <div className="rounded-md border border-gray-200 bg-gray-50 p-3 text-sm">
            <p className="font-medium text-gray-900">{target.itemName}</p>
            <p className="text-gray-600">
              Emprunté par {target.ownerName} · {target.ownerEmail}
            </p>
          </div>
          <div className="flex gap-2 rounded-md bg-amber-50 p-3 text-sm text-amber-900">
            <AlertTriangle className="mt-0.5 h-4 w-4 shrink-0" />
            <p>
              À utiliser quand l'objet a été restitué sans que l'utilisateur le valide dans
              l'application. L'élément est libéré et les rappels automatiques (J+1, J+5, J+7)
              s'arrêtent. L'utilisateur n'est <strong>pas notifié</strong> ; la première personne de
              la file d'attente, s'il y en a une, est prévenue.
            </p>
          </div>
          <div className="flex justify-end gap-2 pt-2">
            <LoadingButton
              onClick={onClose}
              disabled={markReturned.isPending}
              className="border border-gray-300 text-gray-700 hover:bg-gray-50"
            >
              Annuler
            </LoadingButton>
            <LoadingButton
              onClick={confirm}
              loading={markReturned.isPending}
              className="bg-green-600 text-white hover:bg-green-700"
            >
              Marquer comme rendu
            </LoadingButton>
          </div>
        </div>
      )}
    </Modal>
  );
}
