"use client";

import { Minus, Plus } from "lucide-react";
import { useEffect, useId, useState } from "react";
import toast from "react-hot-toast";
import Modal from "@/components/Modal";
import { useUpdateBassineScore } from "@/lib/hooks";
import type { ApiError, BassineScore } from "@/lib/types";

interface BassineScoreModalProps {
  isOpen: boolean;
  onClose: () => void;
  user: BassineScore;
}

export default function BassineScoreModal({ isOpen, onClose, user }: BassineScoreModalProps) {
  const [scoreChange, setScoreChange] = useState(0);
  const updateScoreMutation = useUpdateBassineScore();
  const scoreInputId = useId();

  useEffect(() => {
    if (isOpen) {
      setScoreChange(0);
    }
  }, [isOpen]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (scoreChange === 0) {
      toast.error("Veuillez saisir un changement de score");
      return;
    }

    try {
      await updateScoreMutation.mutateAsync({
        userEmail: user.user_email,
        scoreChange,
      });

      toast.success("Score mis à jour avec succès");
      onClose();
    } catch (err) {
      const error = err as ApiError;
      toast.error(error.response?.data?.error || "Erreur lors de la mise à jour du score");
    }
  };

  const newScore = user.current_score + scoreChange;

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="Modifier le score">
      <form onSubmit={handleSubmit} className="space-y-6">
        {/* User Info */}
        <div className="bg-gray-50 rounded-lg p-4">
          <h3 className="text-sm font-medium text-gray-900 mb-2">Joueur sélectionné</h3>
          <div className="text-sm text-gray-600">
            <p className="font-medium">
              {user.user_first_name} {user.user_last_name}
            </p>
            <p>{user.user_email}</p>
            <p className="mt-1">
              Score actuel: <span className="font-bold text-yellow-600">{user.current_score}</span>
            </p>
          </div>
        </div>

        {/* Score Change */}
        <div>
          <label htmlFor={scoreInputId} className="block text-sm font-medium text-gray-700 mb-2">
            Changement de score
          </label>
          <div className="flex items-center space-x-3">
            <button
              type="button"
              onClick={() => setScoreChange((prev) => prev - 1)}
              className="p-2 bg-red-100 text-red-600 rounded hover:bg-red-200"
            >
              <Minus size={16} />
            </button>
            <input
              id={scoreInputId}
              type="number"
              value={scoreChange}
              onChange={(e) => setScoreChange(parseInt(e.target.value, 10) || 0)}
              className="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 text-center"
              placeholder="0"
            />
            <button
              type="button"
              onClick={() => setScoreChange((prev) => prev + 1)}
              className="p-2 bg-green-100 text-green-600 rounded hover:bg-green-200"
            >
              <Plus size={16} />
            </button>
          </div>
          {scoreChange !== 0 && (
            <p className="mt-2 text-sm text-gray-600">
              Nouveau score: <span className="font-bold text-blue-600">{newScore}</span>
              {scoreChange > 0 ? (
                <span className="text-green-600 ml-2">(+{scoreChange})</span>
              ) : (
                <span className="text-red-600 ml-2">({scoreChange})</span>
              )}
            </p>
          )}
        </div>

        {/* Actions */}
        <div className="flex justify-end space-x-3 pt-4 border-t">
          <button
            type="button"
            onClick={onClose}
            className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200"
          >
            Annuler
          </button>
          <button
            type="submit"
            disabled={updateScoreMutation.isPending || scoreChange === 0}
            className="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {updateScoreMutation.isPending ? "Mise à jour..." : "Mettre à jour"}
          </button>
        </div>
      </form>
    </Modal>
  );
}
