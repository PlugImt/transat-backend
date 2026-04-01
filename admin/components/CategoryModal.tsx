"use client";

import { useState, useEffect } from "react";
import {
  useCreateReservationCategory,
  useUpdateReservationCategory,
} from "@/lib/hooks";
import type { ReservationTreeItem } from "@/lib/api";
import type { Club } from "@/lib/types";

interface CategoryModalProps {
  isOpen: boolean;
  onClose: () => void;
  category: ReservationTreeItem | null;
  parentCategoryId: number | null;
  parentClubId: number | null;
  clubs: Club[];
}

export const CategoryModal = ({
  isOpen,
  onClose,
  category,
  parentCategoryId,
  parentClubId,
  clubs,
}: CategoryModalProps) => {
  const createCategory = useCreateReservationCategory();
  const updateCategory = useUpdateReservationCategory();

  const [name, setName] = useState("");
  const [selectedClubId, setSelectedClubId] = useState<number | null>(null);

  useEffect(() => {
    if (category) {
      setName(category.name || "");
      setSelectedClubId(category.club_id || null);
    } else {
      setName("");
      setSelectedClubId(parentClubId);
    }
  }, [category, parentClubId, isOpen]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!name.trim()) {
      alert("Le nom de la catégorie est requis");
      return;
    }

    try {
      if (category) {
        await updateCategory.mutateAsync({
          id: category.id!,
          data: { name: name.trim() },
        });
      } else {
        await createCategory.mutateAsync({
          name: name.trim(),
          id_club_parent: selectedClubId || undefined,
          id_category_parent: parentCategoryId || undefined,
        });
      }
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
          <h2 className="text-2xl font-bold">
            {category ? "Modifier la catégorie" : "Nouvelle catégorie"}
          </h2>
        </div>

        <form onSubmit={handleSubmit} className="p-6 space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Nom de la catégorie *
            </label>
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
            />
          </div>

          {!category && (
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Club parent {!parentCategoryId && "*"}
              </label>
              <select
                value={selectedClubId || ""}
                onChange={(e) =>
                  setSelectedClubId(e.target.value ? parseInt(e.target.value) : null)
                }
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                required={!parentCategoryId}
                disabled={!!parentCategoryId}
              >
                <option value="">Sélectionner un club</option>
                {clubs.map((club) => (
                  <option key={club.id_clubs} value={club.id_clubs}>
                    {club.name}
                  </option>
                ))}
              </select>
              {parentCategoryId && (
                <p className="text-xs text-gray-500 mt-1">
                  Le club sera hérité de la catégorie parente
                </p>
              )}
            </div>
          )}

          <div className="flex gap-2 justify-end pt-4">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 border border-gray-300 rounded hover:bg-gray-50"
              disabled={createCategory.isPending || updateCategory.isPending}
            >
              Annuler
            </button>
            <button
              type="submit"
              className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
              disabled={createCategory.isPending || updateCategory.isPending}
            >
              {createCategory.isPending || updateCategory.isPending
                ? "Enregistrement..."
                : category
                  ? "Modifier"
                  : "Créer"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
