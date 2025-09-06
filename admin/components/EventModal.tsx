"use client";

import { useState, useEffect } from "react";
import { X, Save } from "lucide-react";
import { Event, ApiError } from "@/lib/types";
import {
  useCreateEvent,
  useUpdateEvent,
  useClubs,
  useUsers,
} from "@/lib/hooks";

interface EventModalProps {
  isOpen: boolean;
  onClose: () => void;
  event?: Event | null;
  onSave: () => void;
}

export default function EventModal({
  isOpen,
  onClose,
  event,
  onSave,
}: EventModalProps) {
  const [formData, setFormData] = useState({
    name: "",
    description: "",
    link: "",
    start_date: "",
    end_date: "",
    location: "",
    picture: "",
    creator: "",
    id_club: "",
  });
  const [error, setError] = useState("");
  const [showEmailSuggestions, setShowEmailSuggestions] = useState(false);
  const [emailFilter, setEmailFilter] = useState("");

  const { data: clubs = [] } = useClubs();
  const { data: users = [] } = useUsers();
  const createEventMutation = useCreateEvent();
  const updateEventMutation = useUpdateEvent();

  // Filtrer les emails des utilisateurs selon la saisie
  const filteredEmails = users
    .filter((user) =>
      user.email.toLowerCase().includes(emailFilter.toLowerCase())
    )
    .slice(0, 5); // Limiter à 5 suggestions

  const handleEmailSelect = (email: string) => {
    setFormData((prev) => ({ ...prev, creator: email }));
    setEmailFilter(email);
    setShowEmailSuggestions(false);
  };

  const handleEmailInputChange = (value: string) => {
    setFormData((prev) => ({ ...prev, creator: value }));
    setEmailFilter(value);
    setShowEmailSuggestions(value.length > 0);
  };

  useEffect(() => {
    if (event) {
      setFormData({
        name: event.name,
        description: event.description || "",
        link: event.link || "",
        start_date: event.start_date
          ? new Date(event.start_date).toISOString().slice(0, 16)
          : "",
        end_date: event.end_date
          ? new Date(event.end_date).toISOString().slice(0, 16)
          : "",
        location: event.location,
        picture: event.picture || "",
        creator: event.creator,
        id_club: event.id_club.toString(),
      });
      setEmailFilter(event.creator);
    } else {
      setFormData({
        name: "",
        description: "",
        link: "",
        start_date: "",
        end_date: "",
        location: "",
        picture: "",
        creator: "",
        id_club: "",
      });
      setEmailFilter("");
    }
    setError("");
    setShowEmailSuggestions(false);
  }, [event, isOpen]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    try {
      const eventData = {
        ...formData,
        id_club: parseInt(formData.id_club),
      };

      if (event) {
        await updateEventMutation.mutateAsync({
          id: event.id_events,
          data: eventData,
        });
      } else {
        await createEventMutation.mutateAsync(eventData);
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
      <div className="bg-white/95 backdrop-blur-sm rounded-lg shadow-lg w-full max-w-2xl p-6 max-h-[90vh] overflow-y-auto animate-slide-up">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-xl font-bold text-gray-900">
            {event ? "Modifier l'événement" : "Créer un événement"}
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

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Date de début *
              </label>
              <input
                type="datetime-local"
                required
                value={formData.start_date}
                onChange={(e) =>
                  setFormData((prev) => ({
                    ...prev,
                    start_date: e.target.value,
                  }))
                }
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Date de fin
              </label>
              <input
                type="datetime-local"
                value={formData.end_date}
                onChange={(e) =>
                  setFormData((prev) => ({ ...prev, end_date: e.target.value }))
                }
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Localisation *
            </label>
            <input
              type="text"
              required
              value={formData.location}
              onChange={(e) =>
                setFormData((prev) => ({ ...prev, location: e.target.value }))
              }
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Club *
            </label>
            <select
              required
              value={formData.id_club}
              onChange={(e) =>
                setFormData((prev) => ({ ...prev, id_club: e.target.value }))
              }
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="">Sélectionner un club</option>
              {clubs.map((club) => (
                <option key={club.id_clubs} value={club.id_clubs}>
                  {club.name}
                </option>
              ))}
            </select>
          </div>

          <div className="relative">
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Email du créateur *
            </label>
            <input
              type="email"
              required
              value={formData.creator}
              onChange={(e) => handleEmailInputChange(e.target.value)}
              onFocus={() =>
                setShowEmailSuggestions(formData.creator.length > 0)
              }
              onBlur={() =>
                setTimeout(() => setShowEmailSuggestions(false), 200)
              }
              placeholder="Commencez à taper pour voir les suggestions..."
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />

            {/* Suggestions d'emails */}
            {showEmailSuggestions && filteredEmails.length > 0 && (
              <div className="absolute z-10 w-full mt-1 bg-white border border-gray-300 rounded-md shadow-lg max-h-48 overflow-y-auto">
                {filteredEmails.map((user) => (
                  <div
                    key={user.id_newf}
                    className="px-3 py-2 hover:bg-gray-100 cursor-pointer text-sm"
                    onClick={() => handleEmailSelect(user.email)}
                  >
                    <div className="font-medium text-gray-900">
                      {user.email}
                    </div>
                    <div className="text-gray-500 text-xs">
                      {user.first_name} {user.last_name}
                    </div>
                  </div>
                ))}
              </div>
            )}
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

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              URL de l&apos;image
            </label>
            <input
              type="url"
              value={formData.picture}
              onChange={(e) =>
                setFormData((prev) => ({ ...prev, picture: e.target.value }))
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
                createEventMutation.isPending || updateEventMutation.isPending
              }
              className="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 disabled:opacity-50"
            >
              <Save className="h-4 w-4 mr-2" />
              {createEventMutation.isPending || updateEventMutation.isPending
                ? "Enregistrement..."
                : event
                ? "Mettre à jour"
                : "Créer"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
