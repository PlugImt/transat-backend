"use client";

import { useState, useEffect } from "react";
import { Save } from "lucide-react";
import { User, ApiError } from "@/lib/types";
import { useCreateUser, useUpdateUser, useRoles } from "@/lib/hooks";
import TagAutocomplete from "./TagAutocomplete";
import Modal from "./Modal";
import toast from "react-hot-toast";

interface UserModalProps {
  isOpen: boolean;
  onClose: () => void;
  user?: User | null;
  onSave: () => void;
}

const languageNames: Record<string, string> = {
  en: "Anglais",
  es: "Espagnol",
  fr: "Français",
  de: "Allemand",
  pt: "Português",
  zh: "Chinois",
};

export default function UserModal({
  isOpen,
  onClose,
  user,
  onSave,
}: UserModalProps) {
  const [formData, setFormData] = useState({
    email: "",
    first_name: "",
    last_name: "",
    phone_number: "",
    formation_name: "",
    graduation_year: "",
    language: "",
    verification_code: "",
    roles: [] as string[],
  });
  const [error, setError] = useState("");

  const createUserMutation = useCreateUser();
  const updateUserMutation = useUpdateUser();
  const { data: roles = [], isLoading: rolesLoading } = useRoles();
  useEffect(() => {
    if (user) {
      setFormData({
        email: user.email,
        first_name: user.first_name,
        last_name: user.last_name,
        phone_number: user.phone_number || "",
        formation_name: user.formation_name || "",
        graduation_year: user.graduation_year?.toString() || "",
        language: user.language || "",
        verification_code: user.verification_code || "",
        roles: user.roles || [],
      });
    } else {
      setFormData({
        email: "",
        first_name: "",
        last_name: "",
        phone_number: "",
        formation_name: "",
        graduation_year: "",
        language: "",
        verification_code: "",
        roles: [],
      });
    }
    setError("");
  }, [user, isOpen]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    try {
      const userData = {
        ...formData,
        graduation_year:
          formData.graduation_year && parseInt(formData.graduation_year) !== 0
            ? parseInt(formData.graduation_year)
            : 0,
        // Convertir les chaînes vides en undefined pour les champs optionnels
        phone_number: formData.phone_number || undefined,
        formation_name: formData.formation_name || undefined,
        language: formData.language || undefined,
      };

      const filteredData = Object.fromEntries(
        Object.entries(userData).filter(
          ([, value]) => value !== "" && value !== null && value !== undefined
        )
      );

      if (user) {
        await toast.promise(
          updateUserMutation.mutateAsync({
            email: user.email,
            data: filteredData,
          }),
          {
            loading: "Mise à jour en cours...",
            success: "Utilisateur modifié avec succès !",
            error: (err: ApiError) =>
              err?.response?.data?.error || "Erreur lors de la modification",
          }
        );
      } else {
        await toast.promise(createUserMutation.mutateAsync(filteredData), {
          loading: "Création en cours...",
          success: "Utilisateur créé avec succès !",
          error: (err: ApiError) =>
            err?.response?.data?.error || "Erreur lors de la création",
        });
      }
      onSave();
      onClose();
    } catch (err: unknown) {
      // Toast promise handles the error display, we just catch to avoid uncaught errors
      setError(
        (err as ApiError)?.response?.data?.error || "L'opération a échoué"
      );
    }
  };

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={user ? "Modifier l'utilisateur" : "Créer un utilisateur"}
      maxWidth="md"
    >
      <div className="max-h-[80vh] overflow-y-auto">
        {/* Scrollable content wrapper */}

        {error && (
          <div className="mb-4 bg-red-50 border border-red-200 rounded-md p-4">
            <div className="text-sm text-red-700">{error}</div>
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-3 sm:space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Email *
            </label>
            <input
              type="email"
              required
              disabled={!!user}
              value={formData.email}
              onChange={(e) =>
                setFormData((prev) => ({ ...prev, email: e.target.value }))
              }
              pattern=".*@imt-atlantique\.net$"
              placeholder="prenom.nom@imt-atlantique.net"
              title="L'email doit se terminer par @imt-atlantique.net"
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:bg-gray-100"
            />
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Prénom *
              </label>
              <input
                type="text"
                required
                value={formData.first_name}
                onChange={(e) =>
                  setFormData((prev) => ({
                    ...prev,
                    first_name: e.target.value,
                  }))
                }
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Nom *
              </label>
              <input
                type="text"
                required
                value={formData.last_name}
                onChange={(e) =>
                  setFormData((prev) => ({
                    ...prev,
                    last_name: e.target.value,
                  }))
                }
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Numéro de téléphone
            </label>
            <input
              type="tel"
              value={formData.phone_number}
              onChange={(e) =>
                setFormData((prev) => ({
                  ...prev,
                  phone_number: e.target.value,
                }))
              }
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <div className="grid grid-cols-1 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Année de diplôme
              </label>
              <input
                type="number"
                min="2025"
                value={formData.graduation_year}
                onChange={(e) =>
                  setFormData((prev) => ({
                    ...prev,
                    graduation_year: e.target.value,
                  }))
                }
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Formation
            </label>
            <select
              value={formData.formation_name}
              onChange={(e) =>
                setFormData((prev) => ({
                  ...prev,
                  formation_name: e.target.value,
                }))
              }
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="">Sélectionner une formation</option>
              <option value="FIL">FIL</option>
              <option value="FISE">FISE</option>
              <option value="FIT">FIT</option>
              <option value="FIP">FIP</option>
              <option value="FID">FID</option>
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Langue
            </label>
            <select
              value={formData.language}
              onChange={(e) =>
                setFormData((prev) => ({ ...prev, language: e.target.value }))
              }
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="">Sélectionner une langue</option>
              {Object.entries(languageNames).map(([code, name]) => (
                <option key={code} value={code}>
                  {name}
                </option>
              ))}
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Code de vérification (mettre 0 pour effacer le code)
            </label>
            <input
              type="text"
              value={formData.verification_code}
              onChange={(e) =>
                setFormData((prev) => ({
                  ...prev,
                  verification_code: e.target.value,
                }))
              }
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Rôles
            </label>
            {rolesLoading ? (
              <div className="h-10 bg-gray-200 rounded animate-pulse"></div>
            ) : (
              <TagAutocomplete
                options={roles.map((role) => ({
                  value: role.name,
                  label: role.name,
                }))}
                selectedTags={formData.roles}
                onChange={(newRoles) =>
                  setFormData((prev) => ({ ...prev, roles: newRoles }))
                }
                placeholder="Rechercher et sélectionner des rôles..."
              />
            )}
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
                createUserMutation.isPending || updateUserMutation.isPending
              }
              className="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 disabled:opacity-50"
            >
              <Save className="h-4 w-4 mr-2" />
              {createUserMutation.isPending || updateUserMutation.isPending
                ? "Enregistrement..."
                : user
                ? "Mettre à jour"
                : "Créer"}
            </button>
          </div>
        </form>
      </div>
    </Modal>
  );
}
