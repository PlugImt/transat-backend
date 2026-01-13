"use client";

import { useState } from "react";
import { useReservationItemsForClub, useUpdateReservationItemMessages } from "../lib/hooks";
import type { ReservationItem } from "../lib/api";

interface ReservationItemsModalProps {
  clubId: number;
  clubName: string;
  isOpen: boolean;
  onClose: () => void;
}

export const ReservationItemsModal = ({
  clubId,
  clubName,
  isOpen,
  onClose,
}: ReservationItemsModalProps) => {
  const { data: items, isLoading } = useReservationItemsForClub(clubId, {
    enabled: isOpen,
  });
  const updateMessages = useUpdateReservationItemMessages();

  const [editingItem, setEditingItem] = useState<number | null>(null);
  const [warningMessage, setWarningMessage] = useState<string>("");
  const [confirmationMessage, setConfirmationMessage] = useState<string>("");

  const handleEdit = (item: ReservationItem) => {
    setEditingItem(item.id);
    setWarningMessage(item.warning_message || "");
    setConfirmationMessage(item.confirmation_message || "");
  };

  const handleCancel = () => {
    setEditingItem(null);
    setWarningMessage("");
    setConfirmationMessage("");
  };

  const handleSave = async (itemId: number) => {
    await updateMessages.mutateAsync({
      itemId,
      messages: {
        warning_message: warningMessage.trim() || null,
        confirmation_message: confirmationMessage.trim() || null,
      },
    });
    setEditingItem(null);
    setWarningMessage("");
    setConfirmationMessage("");
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
      <div className="bg-white rounded-lg shadow-xl w-full max-w-4xl max-h-[90vh] overflow-hidden flex flex-col">
        <div className="p-6 border-b flex justify-between items-center">
          <h2 className="text-2xl font-bold">
            Messages de réservation - {clubName}
          </h2>
          <button
            onClick={onClose}
            className="text-gray-500 hover:text-gray-700 text-2xl"
          >
            ×
          </button>
        </div>

        <div className="p-6 overflow-y-auto flex-1">
          {isLoading ? (
            <div className="text-center py-8">Chargement...</div>
          ) : !items || items.length === 0 ? (
            <div className="text-center py-8 text-gray-500">
              Aucun élément de réservation trouvé pour ce club.
            </div>
          ) : (
            <div className="space-y-6">
              {items.map((item) => (
                <div
                  key={item.id}
                  className="border rounded-lg p-4 bg-gray-50"
                >
                  <div className="flex justify-between items-start mb-4">
                    <div>
                      <h3 className="font-semibold text-lg">{item.name}</h3>
                      {item.description && (
                        <p className="text-sm text-gray-600 mt-1">
                          {item.description}
                        </p>
                      )}
                      {item.slot && (
                        <span className="inline-block mt-2 px-2 py-1 text-xs bg-blue-100 text-blue-800 rounded">
                          Réservation par créneau
                        </span>
                      )}
                    </div>
                    {editingItem !== item.id && (
                      <button
                        onClick={() => handleEdit(item)}
                        className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
                      >
                        Modifier
                      </button>
                    )}
                  </div>

                  {editingItem === item.id ? (
                    <div className="space-y-4">
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
                          Ce message sera affiché avant la confirmation de
                          réservation.
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
                          Ce message sera affiché après la confirmation de
                          réservation.
                        </p>
                      </div>

                      <div className="flex gap-2 justify-end">
                        <button
                          onClick={handleCancel}
                          className="px-4 py-2 border border-gray-300 rounded hover:bg-gray-50"
                          disabled={updateMessages.isPending}
                        >
                          Annuler
                        </button>
                        <button
                          onClick={() => handleSave(item.id)}
                          className="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700"
                          disabled={updateMessages.isPending}
                        >
                          {updateMessages.isPending ? "Enregistrement..." : "Enregistrer"}
                        </button>
                      </div>
                    </div>
                  ) : (
                    <div className="space-y-2 text-sm">
                      <div>
                        <span className="font-medium">Message d&apos;avertissement:</span>
                        <p className="text-gray-600 mt-1">
                          {item.warning_message || (
                            <span className="italic text-gray-400">
                              (Message par défaut)
                            </span>
                          )}
                        </p>
                      </div>
                      <div>
                        <span className="font-medium">Message de confirmation:</span>
                        <p className="text-gray-600 mt-1">
                          {item.confirmation_message || (
                            <span className="italic text-gray-400">
                              (Message par défaut)
                            </span>
                          )}
                        </p>
                      </div>
                    </div>
                  )}
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
