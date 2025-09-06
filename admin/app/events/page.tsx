"use client";

import { useState, useMemo } from "react";
import {
  Plus,
  Edit,
  Trash2,
  Calendar,
  MapPin,
  Users,
  Filter,
} from "lucide-react";
import { Event, EventWithClubName } from "@/lib/types";
import { useEvents, useDeleteEvent } from "@/lib/hooks";
import { useAppStore } from "@/lib/stores/appStore";
import EventModal from "@/components/EventModal";

type DateFilter = "all" | "past" | "future";

export default function EventsPage() {
  const [dateFilter, setDateFilter] = useState<DateFilter>("future");

  const { data: events = [], isLoading } = useEvents();
  const deleteEventMutation = useDeleteEvent();

  const { eventModalOpen, editingEvent, openEventModal, closeEventModal } =
    useAppStore();

  // Fonction pour vérifier si un événement est en cours
  const isEventOngoing = (event: Event) => {
    const now = new Date();
    const startDate = new Date(event.start_date);
    const endDate = new Date(event.end_date);
    return now >= startDate && now <= endDate;
  };

  // Filtrer les événements selon le filtre de date sélectionné
  const filteredEvents = useMemo(() => {
    if (dateFilter === "all") return events;

    const now = new Date();
    return events.filter((event) => {
      const eventStartDate = new Date(event.start_date);

      if (dateFilter === "past") {
        return eventStartDate < now;
      } else if (dateFilter === "future") {
        // Inclure les événements futurs ET en cours
        return eventStartDate >= now || isEventOngoing(event);
      }

      return true;
    });
  }, [events, dateFilter]);

  const handleCreateEvent = () => {
    openEventModal();
  };

  const handleEditEvent = (event: Event) => {
    openEventModal(event);
  };

  const handleDeleteEvent = async (event: Event) => {
    if (!confirm(`Êtes-vous sûr de vouloir supprimer "${event.name}" ?`)) {
      return;
    }

    try {
      await deleteEventMutation.mutateAsync(event.id_events);
    } catch (error) {
      console.error("Failed to delete event:", error);
      alert("Échec de la suppression de l'événement");
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString("fr-FR", {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  if (isLoading) {
    return (
      <div className="p-6">
        <div className="animate-pulse space-y-4">
          <div className="h-8 bg-gray-200 rounded w-1/4"></div>
          <div className="space-y-3">
            {[1, 2, 3].map((i) => (
              <div key={i} className="h-24 bg-gray-200 rounded"></div>
            ))}
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold text-gray-900">
          Gestion des événements
        </h1>
        <div className="flex items-center space-x-4">
          {/* Sélecteur de filtre par date */}
          <div className="flex items-center space-x-2">
            <Filter className="h-4 w-4 text-gray-500" />
            <select
              value={dateFilter}
              onChange={(e) => setDateFilter(e.target.value as DateFilter)}
              className="px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="all">Tous les événements</option>
              <option value="future">Événements futurs</option>
              <option value="past">Événements passés</option>
            </select>
          </div>

          <button
            onClick={handleCreateEvent}
            className="inline-flex items-center px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
          >
            <Plus className="h-4 w-4 mr-2" />
            Créer un événement
          </button>
        </div>
      </div>

      <div className="bg-white shadow rounded-lg overflow-hidden">
        {/* Indicateur du nombre d'événements */}
        <div className="px-6 py-3 bg-gray-50 border-b border-gray-200">
          <div className="flex items-center justify-between">
            <span className="text-sm text-gray-600">
              {filteredEvents.length} événement
              {filteredEvents.length !== 1 ? "s" : ""}
              {dateFilter !== "all" && (
                <span className="ml-1">
                  ({dateFilter === "future" ? "futurs et en cours" : "passés"})
                </span>
              )}
            </span>
            {dateFilter !== "all" && (
              <span className="text-xs text-gray-500">
                sur {events.length} au total
              </span>
            )}
          </div>
        </div>

        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Événement
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Club
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Créateur
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Date début
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Date fin
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Participants
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {filteredEvents &&
                filteredEvents.map((event) => (
                  <tr
                    key={event.id_events}
                    className={`hover:bg-gray-50 ${
                      isEventOngoing(event) ? "bg-orange-50" : ""
                    }`}
                  >
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="flex items-center">
                        {event.picture && (
                          <img
                            src={event.picture}
                            alt={event.name}
                            className="h-10 w-10 rounded object-cover mr-3"
                          />
                        )}
                        <div>
                          <div className="text-sm font-medium text-gray-900 flex items-center">
                            {event.name}
                            {isEventOngoing(event) && (
                              <span className="ml-2 inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-orange-100 text-orange-800">
                                En cours
                              </span>
                            )}
                          </div>
                          <div className="text-sm text-gray-500 flex items-center">
                            <MapPin className="h-3 w-3 mr-1" />
                            {event.location}
                          </div>
                        </div>
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="text-sm text-gray-900">
                        {(event as EventWithClubName).club_name ||
                          `Club ${event.id_club}`}
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="text-sm text-gray-900">
                        {event.creator}
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="flex items-center text-sm text-gray-900">
                        <Calendar className="h-4 w-4 mr-1" />
                        {formatDate(event.start_date)}
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="flex items-center text-sm text-gray-900">
                        <Calendar className="h-4 w-4 mr-1" />
                        {formatDate(event.end_date)}
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="flex items-center text-sm text-gray-900">
                        <Users className="h-4 w-4 mr-1" />
                        {event.attendee_count || 0}
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                      <div className="flex space-x-2">
                        <button
                          onClick={() => handleEditEvent(event)}
                          className="text-indigo-600 hover:text-indigo-900"
                        >
                          <Edit className="h-4 w-4" />
                        </button>
                        <button
                          onClick={() => handleDeleteEvent(event)}
                          className="text-red-600 hover:text-red-900"
                        >
                          <Trash2 className="h-4 w-4" />
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
            </tbody>
          </table>
        </div>

        {filteredEvents && filteredEvents.length === 0 && (
          <div className="text-center py-12">
            <Calendar className="mx-auto h-12 w-12 text-gray-400" />
            <h3 className="mt-2 text-sm font-medium text-gray-900">
              {dateFilter === "all"
                ? "Aucun événement"
                : dateFilter === "future"
                ? "Aucun événement futur ou en cours"
                : "Aucun événement passé"}
            </h3>
            <p className="mt-1 text-sm text-gray-500">
              {dateFilter === "all"
                ? "Commencez par créer un nouvel événement."
                : `Aucun événement ${
                    dateFilter === "future" ? "futur ou en cours" : "passé"
                  } trouvé.`}
            </p>
          </div>
        )}
      </div>

      <EventModal
        isOpen={eventModalOpen}
        onClose={closeEventModal}
        event={editingEvent}
        onSave={() => {}}
      />
    </div>
  );
}
