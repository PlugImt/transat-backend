"use client";

import { Trash2, UserPlus, X } from "lucide-react";
import { useEffect, useId, useState } from "react";
import { useAddClubOwner, useClubOwners, useRemoveClubOwner } from "@/lib/hooks";

interface ClubOwnersModalProps {
  isOpen: boolean;
  onClose: () => void;
  clubId: number | null;
  clubName: string;
}

export default function ClubOwnersModal({
  isOpen,
  onClose,
  clubId,
  clubName,
}: ClubOwnersModalProps) {
  const [email, setEmail] = useState("");
  const emailId = useId();

  const enabled = isOpen && !!clubId;
  const {
    data: owners = [],
    isLoading,
    refetch,
  } = useClubOwners(clubId ?? 0, { enabled });

  const addOwnerMutation = useAddClubOwner();
  const removeOwnerMutation = useRemoveClubOwner();

  useEffect(() => {
    if (!isOpen) {
      setEmail("");
    }
  }, [isOpen]);

  if (!isOpen || !clubId) return null;

  const handleAddOwner = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!email.trim()) return;

    await addOwnerMutation.mutateAsync({ clubId, email: email.trim() });
    setEmail("");
    await refetch();
  };

  const handleRemoveOwner = async (ownerEmail: string) => {
    await removeOwnerMutation.mutateAsync({ clubId, email: ownerEmail });
    await refetch();
  };

  return (
    <div className="fixed inset-0 bg-gray-900/30 backdrop-blur-[2px] flex items-center justify-center z-50 animate-fade-in">
      <div className="bg-white/95 backdrop-blur-sm rounded-lg shadow-lg w-full max-w-lg p-6 max-h-[90vh] overflow-y-auto animate-slide-up">
        <div className="flex items-center justify-between mb-4">
          <div>
            <h2 className="text-xl font-bold text-gray-900">Responsables du club</h2>
            <p className="text-sm text-gray-500 mt-1">{clubName}</p>
          </div>
          <button
            type="button"
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600"
          >
            <X className="h-6 w-6" />
          </button>
        </div>

        <div className="space-y-4">
          <form onSubmit={handleAddOwner} className="space-y-2">
            <label htmlFor={emailId} className="block text-sm font-medium text-gray-700">
              Ajouter un responsable (email)
            </label>
            <div className="flex space-x-2">
              <input
                id={emailId}
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="prenom.nom@imt-atlantique.net"
                className="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
              <button
                type="submit"
                disabled={addOwnerMutation.isPending}
                className="inline-flex items-center px-3 py-2 rounded-md bg-blue-600 text-white text-sm font-medium hover:bg-blue-700 disabled:opacity-50"
              >
                <UserPlus className="h-4 w-4 mr-1" />
                Ajouter
              </button>
            </div>
            <p className="text-xs text-gray-500">
              L&apos;utilisateur doit exister dans Transat. Il sera automatiquement ajouté
              comme membre du club si nécessaire.
            </p>
          </form>

          <div className="border-t border-gray-200 pt-4">
            <h3 className="text-sm font-medium text-gray-700 mb-2">
              Responsables actuels ({owners.length})
            </h3>
            {isLoading ? (
              <p className="text-sm text-gray-500">Chargement des responsables...</p>
            ) : owners.length === 0 ? (
              <p className="text-sm text-gray-500">Aucun responsable défini pour ce club.</p>
            ) : (
              <ul className="divide-y divide-gray-200 rounded-md border border-gray-200">
                {owners.map((owner) => (
                  <li
                    key={owner.email}
                    className="flex items-center justify-between px-3 py-2 bg-white"
                  >
                    <div>
                      <p className="text-sm font-medium text-gray-900">
                        {owner.first_name} {owner.last_name}
                      </p>
                      <p className="text-xs text-gray-500">{owner.email}</p>
                    </div>
                    <button
                      type="button"
                      onClick={() => handleRemoveOwner(owner.email)}
                      disabled={removeOwnerMutation.isPending}
                      className="inline-flex items-center px-2 py-1 rounded-md text-xs font-medium text-red-600 hover:text-red-800 hover:bg-red-50 disabled:opacity-50"
                    >
                      <Trash2 className="h-3 w-3 mr-1" />
                      Retirer
                    </button>
                  </li>
                ))}
              </ul>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

