"use client";

import { Save, X } from "lucide-react";
import { useEffect, useId, useState } from "react";
import { useCreateTraqArticle, useTraqTypes, useUpdateTraqArticle } from "@/lib/hooks";
import type { ApiError, TraqArticle } from "@/lib/types";

interface TraqArticleModalProps {
  isOpen: boolean;
  onClose: () => void;
  article?: TraqArticle | null;
  onSave?: () => void;
}

const emptyForm = {
  name: "",
  description: "",
  picture: "",
  traq_type: "",
  price: 0,
  price_half: 0,
  alcohol: 0,
  limited: false,
  out_of_stock: false,
  disabled: false,
};

export default function TraqArticleModal({
  isOpen,
  onClose,
  article,
  onSave,
}: TraqArticleModalProps) {
  const [formData, setFormData] = useState(emptyForm);
  const [error, setError] = useState("");

  const { data: types = [] } = useTraqTypes();
  const createMutation = useCreateTraqArticle();
  const updateMutation = useUpdateTraqArticle();

  const nameId = useId();
  const descriptionId = useId();
  const pictureId = useId();
  const typeId = useId();
  const priceId = useId();
  const priceHalfId = useId();
  const alcoholId = useId();

  useEffect(() => {
    if (article) {
      setFormData({
        name: article.name,
        description: article.description || "",
        picture: article.picture || "",
        traq_type: article.traq_type || "",
        price: article.price ?? 0,
        price_half: article.price_half ?? 0,
        alcohol: article.alcohol ?? 0,
        limited: article.limited ?? false,
        out_of_stock: article.out_of_stock ?? false,
        disabled: article.disabled ?? false,
      });
    } else {
      setFormData(emptyForm);
    }
    setError("");
  }, [article, isOpen]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    if (!formData.name || !formData.description || !formData.picture || !formData.traq_type) {
      setError("Nom, description, image et type sont requis");
      return;
    }

    try {
      if (article) {
        await updateMutation.mutateAsync({
          id: article.id_traq,
          data: formData,
        });
      } else {
        await createMutation.mutateAsync(formData);
      }
      onSave?.();
      onClose();
    } catch (err: unknown) {
      setError((err as ApiError)?.response?.data?.error || "L'opération a échoué");
    }
  };

  if (!isOpen) return null;

  const isPending = createMutation.isPending || updateMutation.isPending;

  return (
    <div className="fixed inset-0 bg-gray-900/30 backdrop-blur-[2px] flex items-center justify-center z-50 animate-fade-in">
      <div className="bg-white/95 backdrop-blur-sm rounded-lg shadow-lg w-full max-w-lg p-6 max-h-[90vh] overflow-y-auto animate-slide-up">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-xl font-bold text-gray-900">
            {article ? "Modifier l'article Traq" : "Créer un article Traq"}
          </h2>
          <button type="button" onClick={onClose} className="text-gray-400 hover:text-gray-600">
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
            <label htmlFor={nameId} className="block text-sm font-medium text-gray-700 mb-1">
              Nom *
            </label>
            <input
              id={nameId}
              type="text"
              required
              value={formData.name}
              onChange={(e) => setFormData((prev) => ({ ...prev, name: e.target.value }))}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-amber-500"
            />
          </div>

          <div>
            <label htmlFor={descriptionId} className="block text-sm font-medium text-gray-700 mb-1">
              Description *
            </label>
            <textarea
              id={descriptionId}
              required
              rows={3}
              value={formData.description}
              onChange={(e) => setFormData((prev) => ({ ...prev, description: e.target.value }))}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-amber-500"
            />
          </div>

          <div>
            <label htmlFor={pictureId} className="block text-sm font-medium text-gray-700 mb-1">
              URL de l&apos;image *
            </label>
            <input
              id={pictureId}
              type="url"
              required
              value={formData.picture}
              onChange={(e) => setFormData((prev) => ({ ...prev, picture: e.target.value }))}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-amber-500"
            />
          </div>

          <div>
            <label htmlFor={typeId} className="block text-sm font-medium text-gray-700 mb-1">
              Type *
            </label>
            <select
              id={typeId}
              required
              value={formData.traq_type}
              onChange={(e) => setFormData((prev) => ({ ...prev, traq_type: e.target.value }))}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-amber-500"
            >
              <option value="">Sélectionner un type</option>
              {types.map((type) => (
                <option key={type.id_traq_types} value={type.name}>
                  {type.name}
                </option>
              ))}
            </select>
            {types.length === 0 && (
              <p className="mt-1 text-xs text-amber-600">
                Aucun type disponible. Créez-en un dans l&apos;onglet Types.
              </p>
            )}
          </div>

          <div className="grid grid-cols-3 gap-3">
            <div>
              <label htmlFor={priceId} className="block text-sm font-medium text-gray-700 mb-1">
                Prix (€)
              </label>
              <input
                id={priceId}
                type="number"
                min={0}
                step={0.01}
                value={formData.price}
                onChange={(e) =>
                  setFormData((prev) => ({ ...prev, price: parseFloat(e.target.value) || 0 }))
                }
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-amber-500"
              />
            </div>
            <div>
              <label htmlFor={priceHalfId} className="block text-sm font-medium text-gray-700 mb-1">
                Demi (€)
              </label>
              <input
                id={priceHalfId}
                type="number"
                min={0}
                step={0.01}
                value={formData.price_half}
                onChange={(e) =>
                  setFormData((prev) => ({
                    ...prev,
                    price_half: parseFloat(e.target.value) || 0,
                  }))
                }
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-amber-500"
              />
            </div>
            <div>
              <label htmlFor={alcoholId} className="block text-sm font-medium text-gray-700 mb-1">
                Alcool (%)
              </label>
              <input
                id={alcoholId}
                type="number"
                min={0}
                step={0.1}
                value={formData.alcohol}
                onChange={(e) =>
                  setFormData((prev) => ({ ...prev, alcohol: parseFloat(e.target.value) || 0 }))
                }
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-amber-500"
              />
            </div>
          </div>

          <div className="flex flex-wrap gap-4">
            <label className="inline-flex items-center space-x-2 text-sm text-gray-700">
              <input
                type="checkbox"
                checked={formData.limited}
                onChange={(e) => setFormData((prev) => ({ ...prev, limited: e.target.checked }))}
                className="rounded border-gray-300 text-amber-600 focus:ring-amber-500"
              />
              <span>Édition limitée</span>
            </label>
            <label className="inline-flex items-center space-x-2 text-sm text-gray-700">
              <input
                type="checkbox"
                checked={formData.out_of_stock}
                onChange={(e) =>
                  setFormData((prev) => ({ ...prev, out_of_stock: e.target.checked }))
                }
                className="rounded border-gray-300 text-amber-600 focus:ring-amber-500"
              />
              <span>Rupture de stock</span>
            </label>
            <label className="inline-flex items-center space-x-2 text-sm text-gray-700">
              <input
                type="checkbox"
                checked={formData.disabled}
                onChange={(e) => setFormData((prev) => ({ ...prev, disabled: e.target.checked }))}
                className="rounded border-gray-300 text-amber-600 focus:ring-amber-500"
              />
              <span>Désactivé</span>
            </label>
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
              disabled={isPending}
              className="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-amber-600 hover:bg-amber-700 disabled:opacity-50"
            >
              <Save className="h-4 w-4 mr-2" />
              {isPending ? "Enregistrement..." : article ? "Mettre à jour" : "Créer"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
