"use client";

import { useState, useEffect } from "react";
import { X, Save } from "lucide-react";
import { User, ApiError } from "@/lib/types";
import { useCreateUser, useUpdateUser } from "@/lib/hooks";

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
    campus: "",
    formation_name: "",
    graduation_year: "",
    language: "",
    verification_code: "",
    roles: [] as string[],
  });
  const [error, setError] = useState("");

  const createUserMutation = useCreateUser();
  const updateUserMutation = useUpdateUser();
  useEffect(() => {
    if (user) {
      setFormData({
        email: user.email,
        first_name: user.first_name,
        last_name: user.last_name,
        phone_number: user.phone_number || "",
        campus: user.campus || "",
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
        campus: "",
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
        graduation_year: formData.graduation_year
          ? parseInt(formData.graduation_year)
          : undefined,
        // Convertir les chaînes vides en undefined pour les champs optionnels
        campus: formData.campus || undefined,
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
        await updateUserMutation.mutateAsync({
          email: user.email,
          data: filteredData,
        });
      } else {
        await createUserMutation.mutateAsync(filteredData);
      }
      onSave();
      onClose();
    } catch (err: unknown) {
      setError(
        (err as ApiError)?.response?.data?.error || "L'opération a échoué"
      );
    }
  };

  const toggleRole = (role: string) => {
    setFormData((prev) => ({
      ...prev,
      roles: prev.roles.includes(role)
        ? prev.roles.filter((r) => r !== role)
        : [...prev.roles, role],
    }));
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-gray-900/30 backdrop-blur-[2px] flex items-center justify-center z-50 animate-fade-in">
      <div className="bg-white/95 backdrop-blur-sm rounded-lg shadow-lg w-full max-w-md p-6 max-h-[90vh] overflow-y-auto animate-slide-up">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-xl font-bold text-gray-900">
            {user ? "Modifier l'utilisateur" : "Créer un utilisateur"}
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

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Campus
              </label>
              <select
                value={formData.campus}
                onChange={(e) =>
                  setFormData((prev) => ({ ...prev, campus: e.target.value }))
                }
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="">Sélectionner un campus</option>
                <option value="NANTES">Nantes</option>
                <option value="BREST">Brest</option>
                <option value="RENNES">Rennes</option>
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Année de diplôme
              </label>
              <input
                type="number"
                min="2020"
                max="2030"
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
            <input
              type="text"
              value={formData.formation_name}
              onChange={(e) =>
                setFormData((prev) => ({
                  ...prev,
                  formation_name: e.target.value,
                }))
              }
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
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
            <div className="space-y-2">
              {["ADMIN", "NEWF", "VERIFYING"].map((role) => (
                <label key={role} className="flex items-center">
                  <input
                    type="checkbox"
                    checked={formData.roles.includes(role)}
                    onChange={() => toggleRole(role)}
                    className="mr-2 h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                  />
                  <span className="text-sm text-gray-700">{role}</span>
                </label>
              ))}
            </div>
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
    </div>
  );
}
