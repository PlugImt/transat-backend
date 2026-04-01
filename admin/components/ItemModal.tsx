"use client";

import { useState, useEffect } from "react";
import {
  useCreateReservationItem,
  useUpdateReservationItem,
} from "@/lib/hooks";
import type { ReservationItem, ReservationTreeItem } from "@/lib/api";
import type { Club } from "@/lib/types";

interface ItemModalProps {
  isOpen: boolean;
  onClose: () => void;
  item: ReservationItem | null;
  parentCategoryId: number | null;
  parentClubId: number | null;
  clubs: Club[];
  categories: ReservationTreeItem[];
}

export const ItemModal = ({
  isOpen,
  onClose,
  item,
  parentCategoryId,
  parentClubId,
  clubs,
  categories,
}: ItemModalProps) => {
  const createItem = useCreateReservationItem();
  const updateItem = useUpdateReservationItem();

  const [name, setName] = useState("");
  const [slot, setSlot] = useState(false);
  const [description, setDescription] = useState("");
  const [location, setLocation] = useState("");
  const [selectedClubId, setSelectedClubId] = useState<number | null>(null);
  const [selectedCategoryId, setSelectedCategoryId] = useState<number | null>(null);

  useEffect(() => {
    if (item) {
      setName(item.name || "");
      setSlot(item.slot || false);
      setDescription(item.description || "");
      setLocation(item.location || "");
    } else {
      setName("");
      setSlot(false);
      setDescription("");
      setLocation("");
      setSelectedClubId(parentClubId);
      setSelectedCategoryId(parentCategoryId);
    }
  }, [item, parentCategoryId, parentClubId, isOpen]);

  const flattenCategories = (
    cats: ReservationTreeItem[],
    level = 0,
  ): Array<{ id: number; name: string; level: number }> => {
    const result: Array<{ id: number; name: string; level: number }> = [];
    cats.forEach((cat) => {
      if (cat.type === "category" && cat.id) {
        result.push({ id: cat.id, name: "  ".repeat(level) + cat.name, level });
        if (cat.children) {
          result.push(...flattenCategories(cat.children, level + 1));
        }
      }
    });
    return result;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!name.trim()) {
      alert("Le nom de l'élément est requis");
      return;
    }

    try {
      if (item) {
        await updateItem.mutateAsync({
          id: item.id,
          data: {
            name: name.trim(),
            slot,
            description: description.trim() || null,
            location: location.trim() || null,
          },
        });
      } else {
        await createItem.mutateAsync({
          name: name.trim(),
          slot,
          description: description.trim() || undefined,
          location: location.trim() || undefined,
          id_club_parent: selectedClubId || undefined,
          id_category_parent: selectedCategoryId || undefined,
        });
      }
      onClose();
    } catch (error: any) {
      alert(error?.response?.data?.error || "Erreur lors de l'enregistrement");
    }
  };

  if (!isOpen) return null;

  const flatCategories = flattenCategories(categories);

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
      <div className="bg-white rounded-lg shadow-xl w-full max-w-md max-h-[90vh] overflow-y-auto">
        <div className="p-6 border-b">
          <h2 className="text-2xl font-bold">
            {item ? "Modifier l'élément" : "Nouvel élément"}
          </h2>
        </div>

        <form onSubmit={handleSubmit} className="p-6 space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Nom de l'élément *
            </label>
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
            />
          </div>

          <div>
            <label className="flex items-center gap-2">
              <input
                type="checkbox"
                checked={slot}
                onChange={(e) => setSlot(e.target.checked)}
                className="rounded"
              />
              <span className="text-sm font-medium text-gray-700">
                Réservation par créneau
              </span>
            </label>
          </div>

          {!item && (
            <>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Catégorie parente
                </label>
                <select
                  value={selectedCategoryId || ""}
                  onChange={(e) =>
                    setSelectedCategoryId(
                      e.target.value ? parseInt(e.target.value) : null,
                    )
                  }
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value="">Aucune (élément racine)</option>
                  {flatCategories.map((cat) => (
                    <option key={cat.id} value={cat.id}>
                      {cat.name}
                    </option>
                  ))}
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Club parent {!selectedCategoryId && "*"}
                </label>
                <select
                  value={selectedClubId || ""}
                  onChange={(e) =>
                    setSelectedClubId(e.target.value ? parseInt(e.target.value) : null)
                  }
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  required={!selectedCategoryId}
                  disabled={!!selectedCategoryId}
                >
                  <option value="">Sélectionner un club</option>
                  {clubs.map((club) => (
                    <option key={club.id_clubs} value={club.id_clubs}>
                      {club.name}
                    </option>
                  ))}
                </select>
                {selectedCategoryId && (
                  <p className="text-xs text-gray-500 mt-1">
                    Le club sera hérité de la catégorie parente
                  </p>
                )}
              </div>
            </>
          )}

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Description
            </label>
            <textarea
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              rows={3}
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Localisation
            </label>
            <input
              type="text"
              value={location}
              onChange={(e) => setLocation(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <div className="flex gap-2 justify-end pt-4">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 border border-gray-300 rounded hover:bg-gray-50"
              disabled={createItem.isPending || updateItem.isPending}
            >
              Annuler
            </button>
            <button
              type="submit"
              className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
              disabled={createItem.isPending || updateItem.isPending}
            >
              {createItem.isPending || updateItem.isPending
                ? "Enregistrement..."
                : item
                  ? "Modifier"
                  : "Créer"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
