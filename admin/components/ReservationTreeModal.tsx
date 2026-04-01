"use client";

import { useState, useEffect } from "react";
import { useUpdateReservationItem, useReservationTree } from "@/lib/hooks";

interface ReservationTreeModalProps {
  itemId: number;
  itemName: string;
  isOpen: boolean;
  onClose: () => void;
}

const findItemInTree = (
  tree: any[],
  itemId: number,
): { warning_message?: string; confirmation_message?: string } | null => {
  for (const node of tree) {
    if (node.items) {
      const item = node.items.find((i: any) => i.id === itemId);
      if (item) {
        return {
          warning_message: item.warning_message,
          confirmation_message: item.confirmation_message,
        };
      }
    }
    if (node.children) {
      const found = findItemInTree(node.children, itemId);
      if (found) return found;
    }
  }
  return null;
};

export const ReservationTreeModal = ({
  itemId,
  itemName,
  isOpen,
  onClose,
}: ReservationTreeModalProps) => {
  const updateItem = useUpdateReservationItem();
  const { data: tree } = useReservationTree();

  const [warningMessage, setWarningMessage] = useState("");
  const [confirmationMessage, setConfirmationMessage] = useState("");

  useEffect(() => {
    if (isOpen && tree) {
      const itemData = findItemInTree(tree, itemId);
      if (itemData) {
        setWarningMessage(itemData.warning_message || "");
        setConfirmationMessage(itemData.confirmation_message || "");
      } else {
        setWarningMessage("");
        setConfirmationMessage("");
      }
    } else if (!isOpen) {
      setWarningMessage("");
      setConfirmationMessage("");
    }
  }, [isOpen, tree, itemId]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      await updateItem.mutateAsync({
        id: itemId,
        data: {
          warning_message: warningMessage.trim() || null,
          confirmation_message: confirmationMessage.trim() || null,
        },
      });
      onClose();
    } catch (error: any) {
      alert(error?.response?.data?.error || "Erreur lors de l'enregistrement");
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
      <div className="bg-white rounded-lg shadow-xl w-full max-w-md">
        <div className="p-6 border-b">
          <h2 className="text-2xl font-bold">Messages de réservation</h2>
          <p className="text-sm text-gray-600 mt-1">{itemName}</p>
        </div>

        <form onSubmit={handleSubmit} className="p-6 space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Message d&apos;avertissement
            </label>
            <textarea
              value={warningMessage}
              onChange={(e) => setWarningMessage(e.target.value)}
              placeholder="Message affiché dans le dialogue de réservation (laisser vide pour utiliser le message par défaut)"
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              rows={3}
            />
            <p className="text-xs text-gray-500 mt-1">
              Ce message sera affiché avant la confirmation de réservation.
            </p>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Message de confirmation
            </label>
            <textarea
              value={confirmationMessage}
              onChange={(e) => setConfirmationMessage(e.target.value)}
              placeholder="Message affiché après la confirmation (laisser vide pour utiliser le message par défaut)"
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              rows={3}
            />
            <p className="text-xs text-gray-500 mt-1">
              Ce message sera affiché après la confirmation de réservation.
            </p>
          </div>

          <div className="flex gap-2 justify-end pt-4">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 border border-gray-300 rounded hover:bg-gray-50"
              disabled={updateItem.isPending}
            >
              Annuler
            </button>
            <button
              type="submit"
              className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
              disabled={updateItem.isPending}
            >
              {updateItem.isPending ? "Enregistrement..." : "Enregistrer"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
