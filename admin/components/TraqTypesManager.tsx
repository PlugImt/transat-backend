"use client";

import { Edit, Plus, Save, Tags, Trash2, X } from "lucide-react";
import { useId, useRef, useState } from "react";
import toast from "react-hot-toast";
import { PageLoading } from "@/components/LoadingSpinner";
import {
  useCreateTraqType,
  useDeleteTraqType,
  useTraqTypes,
  useUpdateTraqType,
} from "@/lib/hooks";
import type { ApiError, TraqType } from "@/lib/types";

export default function TraqTypesManager() {
  const { data: types = [], isLoading, error } = useTraqTypes();
  const createMutation = useCreateTraqType();
  const updateMutation = useUpdateTraqType();
  const deleteMutation = useDeleteTraqType();

  const [newName, setNewName] = useState("");
  const [editingId, setEditingId] = useState<number | null>(null);
  const [editingName, setEditingName] = useState("");
  const editingIdRef = useRef(editingId);
  editingIdRef.current = editingId;
  const newNameId = useId();

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    const name = newName.trim();
    if (!name) return;

    try {
      await createMutation.mutateAsync({ name });
      setNewName("");
      toast.success("Type créé avec succès");
    } catch (err: unknown) {
      toast.error((err as ApiError)?.response?.data?.error || "Échec de la création du type");
    }
  };

  const startEdit = (type: TraqType) => {
    setEditingId(type.id_traq_types);
    setEditingName(type.name);
  };

  const cancelEdit = () => {
    setEditingId(null);
    setEditingName("");
  };

  const handleUpdate = async (id: number) => {
    const name = editingName.trim();
    if (!name) return;

    try {
      await updateMutation.mutateAsync({ id, data: { name } });
      if (editingIdRef.current === id) {
        cancelEdit();
      }
      toast.success("Type mis à jour");
    } catch (err: unknown) {
      toast.error((err as ApiError)?.response?.data?.error || "Échec de la mise à jour");
    }
  };

  const handleDelete = async (type: TraqType) => {
    if (!confirm(`Supprimer le type « ${type.name} » ?`)) return;

    try {
      await deleteMutation.mutateAsync(type.id_traq_types);
      toast.success("Type supprimé");
    } catch (err: unknown) {
      const apiError = err as ApiError;
      toast.error(
        apiError?.response?.data?.error || "Échec de la suppression du type",
      );
    }
  };

  if (isLoading) {
    return <PageLoading text="Chargement des types Traq..." />;
  }

  if (error) {
    return (
      <div className="bg-red-50 border border-red-200 rounded-md p-4">
        <div className="text-sm text-red-700">
          {(error as ApiError)?.message || "Échec de la récupération des types"}
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <form onSubmit={handleCreate} className="flex flex-col sm:flex-row gap-3">
        <div className="flex-1">
          <label htmlFor={newNameId} className="sr-only">
            Nom du type
          </label>
          <input
            id={newNameId}
            type="text"
            value={newName}
            onChange={(e) => setNewName(e.target.value)}
            placeholder="Nouveau type (ex. Bière, Soft...)"
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-amber-500"
          />
        </div>
        <button
          type="submit"
          disabled={!newName.trim() || createMutation.isPending}
          className="inline-flex items-center justify-center px-4 py-2 bg-amber-600 text-white rounded-md hover:bg-amber-700 disabled:opacity-50 text-sm"
        >
          <Plus className="h-4 w-4 mr-2" />
          {createMutation.isPending ? "Création..." : "Ajouter"}
        </button>
      </form>

      <div className="bg-white shadow rounded-lg overflow-hidden">
        {types.length === 0 ? (
          <div className="text-center py-12">
            <Tags className="mx-auto h-12 w-12 text-gray-400" />
            <h3 className="mt-2 text-sm font-medium text-gray-900">Aucun type</h3>
            <p className="mt-1 text-sm text-gray-500">
              Créez un type avant d&apos;ajouter des articles.
            </p>
          </div>
        ) : (
          <ul className="divide-y divide-gray-200">
            {types.map((type) => (
              <li
                key={type.id_traq_types}
                className="px-4 py-3 flex items-center justify-between gap-3 hover:bg-gray-50"
              >
                {editingId === type.id_traq_types ? (
                  <div className="flex-1 flex items-center gap-2">
                    <input
                      type="text"
                      value={editingName}
                      onChange={(e) => setEditingName(e.target.value)}
                      className="flex-1 px-3 py-1.5 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-amber-500 text-sm"
                      autoFocus
                      onKeyDown={(e) => {
                        if (e.key === "Enter") {
                          e.preventDefault();
                          handleUpdate(type.id_traq_types);
                        }
                        if (e.key === "Escape") cancelEdit();
                      }}
                    />
                    <button
                      type="button"
                      onClick={() => handleUpdate(type.id_traq_types)}
                      disabled={updateMutation.isPending}
                      className="p-2 text-green-600 hover:bg-green-50 rounded-full"
                      title="Enregistrer"
                    >
                      <Save className="h-4 w-4" />
                    </button>
                    <button
                      type="button"
                      onClick={cancelEdit}
                      className="p-2 text-gray-400 hover:bg-gray-100 rounded-full"
                      title="Annuler"
                    >
                      <X className="h-4 w-4" />
                    </button>
                  </div>
                ) : (
                  <>
                    <div>
                      <p className="text-sm font-medium text-gray-900">{type.name}</p>
                      <p className="text-xs text-gray-500">ID {type.id_traq_types}</p>
                    </div>
                    <div className="flex items-center space-x-1">
                      <button
                        type="button"
                        onClick={() => startEdit(type)}
                        className="p-2 text-gray-400 hover:text-amber-600 rounded-full hover:bg-gray-100"
                        title="Renommer"
                      >
                        <Edit className="h-4 w-4" />
                      </button>
                      <button
                        type="button"
                        onClick={() => handleDelete(type)}
                        className="p-2 text-gray-400 hover:text-red-600 rounded-full hover:bg-gray-100"
                        title="Supprimer"
                      >
                        <Trash2 className="h-4 w-4" />
                      </button>
                    </div>
                  </>
                )}
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}
