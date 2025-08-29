"use client";

import { useState, useEffect } from "react";
import { X, Save } from "lucide-react";
import { Club, ApiError } from "@/lib/types";
import { useCreateClub, useUpdateClub } from "@/lib/hooks";

interface ClubModalProps {
  isOpen: boolean;
  onClose: () => void;
  club?: Club | null;
  onSave: () => void;
}

export default function ClubModal({
  isOpen,
  onClose,
  club,
  onSave,
}: ClubModalProps) {
  const [formData, setFormData] = useState({
    name: "",
    picture: "",
    description: "",
    location: "",
    link: "",
  });
  const [error, setError] = useState("");

  const createClubMutation = useCreateClub();
  const updateClubMutation = useUpdateClub();

  useEffect(() => {
    if (club) {
      setFormData({
        name: club.name,
        picture: club.picture,
        description: club.description || "",
        location: club.location || "",
        link: club.link || "",
      });
    } else {
      setFormData({
        name: "",
        picture: "",
        description: "",
        location: "",
        link: "",
      });
    }
    setError("");
  }, [club, isOpen]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    try {
      if (club) {
        await updateClubMutation.mutateAsync({
          id: club.id_clubs,
          data: formData,
        });
      } else {
        await createClubMutation.mutateAsync(formData);
      }
      onSave();
      onClose();
    } catch (err: unknown) {
      setError(
        (err as ApiError)?.response?.data?.error || "L'opération a échoué"
      );
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-gray-900/30 backdrop-blur-[2px] flex items-center justify-center z-50 animate-fade-in">
      <div className="bg-white/95 backdrop-blur-sm rounded-lg shadow-lg w-full max-w-md p-6 max-h-[90vh] overflow-y-auto animate-slide-up">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-xl font-bold text-gray-900">
            {club ? "Modifier le club" : "Créer un club"}
          </h2>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600"
          >
            <X className="h-6 w-6" />
          </button>
        </div>

        {error && (
          <div className="mb-4 bg-red-50 border border-red-200 rounded-md p-4">
            <div className="text-sm text-red-700">{error}</div>
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Nom *
            </label>
            <input
              type="text"
              required
              value={formData.name}
              onChange={(e) =>
                setFormData((prev) => ({ ...prev, name: e.target.value }))
              }
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              URL de l&apos;image *
            </label>
            <input
              type="url"
              required
              value={formData.picture}
              onChange={(e) =>
                setFormData((prev) => ({ ...prev, picture: e.target.value }))
              }
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Description
            </label>
            <textarea
              value={formData.description}
              onChange={(e) =>
                setFormData((prev) => ({
                  ...prev,
                  description: e.target.value,
                }))
              }
              rows={3}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Localisation
            </label>
            <input
              type="text"
              value={formData.location}
              onChange={(e) =>
                setFormData((prev) => ({ ...prev, location: e.target.value }))
              }
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Lien
            </label>
            <input
              type="url"
              value={formData.link}
              onChange={(e) =>
                setFormData((prev) => ({ ...prev, link: e.target.value }))
              }
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <div className="flex justify-end space-x-3 pt-6">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-50"
            >
              Annuler
            </button>
            <button
              type="submit"
              disabled={
                createClubMutation.isPending || updateClubMutation.isPending
              }
              className="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 disabled:opacity-50"
            >
              <Save className="h-4 w-4 mr-2" />
              {createClubMutation.isPending || updateClubMutation.isPending
                ? "Enregistrement..."
                : club
                ? "Mettre à jour"
                : "Créer"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
