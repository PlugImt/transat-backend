"use client";

import { Building, Edit, ExternalLink, MapPin, Plus, Trash2, Users } from "lucide-react";
import Image from "next/image";
import ClubModal from "@/components/ClubModal";
import { useClubs, useDeleteClub } from "@/lib/hooks";
import { useAppStore } from "@/lib/stores/appStore";
import type { Club, ClubWithResponsible } from "@/lib/types";

export default function ClubsPage() {
  const { data: clubs = [], isLoading } = useClubs();
  const deleteClubMutation = useDeleteClub();

  const { clubModalOpen, editingClub, openClubModal, closeClubModal } = useAppStore();

  const handleCreateClub = () => {
    openClubModal();
  };

  const handleEditClub = (club: Club) => {
    openClubModal(club);
  };

  const handleDeleteClub = async (club: Club) => {
    if (!confirm(`Êtes-vous sûr de vouloir supprimer "${club.name}" ?`)) {
      return;
    }

    try {
      await deleteClubMutation.mutateAsync(club.id_clubs);
    } catch (error) {
      console.error("Failed to delete club:", error);
      alert("Échec de la suppression du club");
    }
  };

  if (isLoading) {
    return (
      <div className="p-6">
        <div className="animate-pulse space-y-4">
          <div className="h-8 bg-gray-200 rounded w-1/4"></div>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {[1, 2, 3].map((i) => (
              <div key={i} className="h-64 bg-gray-200 rounded"></div>
            ))}
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold text-gray-900">Gestion des clubs</h1>
        <button
          type="button"
          onClick={handleCreateClub}
          className="inline-flex items-center px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
        >
          <Plus className="h-4 w-4 mr-2" />
          Créer un club
        </button>
      </div>

      <div className="bg-white shadow rounded-lg overflow-hidden">
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Club
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Description
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Localisation
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Membres
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Responsable
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {clubs?.map((club) => (
                <tr key={club.id_clubs} className="hover:bg-gray-50">
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="flex items-center">
                      {club.picture && (
                        <Image
                          src={club.picture}
                          alt={club.name}
                          width={40}
                          height={40}
                          className="h-10 w-10 rounded object-cover mr-3"
                        />
                      )}
                      <div>
                        <div className="text-sm font-medium text-gray-900">{club.name}</div>
                        {club.link && (
                          <a
                            href={club.link}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="text-xs text-blue-600 hover:text-blue-800 flex items-center"
                          >
                            <ExternalLink className="h-3 w-3 mr-1" />
                            Site web
                          </a>
                        )}
                      </div>
                    </div>
                  </td>
                  <td className="px-6 py-4">
                    <div className="text-sm text-gray-900 max-w-xs truncate">
                      {club.description || "-"}
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="text-sm text-gray-900 flex items-center">
                      {club.location ? (
                        <>
                          <MapPin className="h-4 w-4 mr-1" />
                          {club.location}
                        </>
                      ) : (
                        "-"
                      )}
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="flex items-center text-sm text-gray-900">
                      <Users className="h-4 w-4 mr-1" />
                      {club.member_count || 0}
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="text-sm text-gray-900">
                      {(club as ClubWithResponsible).responsible
                        ? `${(club as ClubWithResponsible).responsible?.first_name} ${
                            (club as ClubWithResponsible).responsible?.last_name
                          }`
                        : "-"}
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                    <div className="flex space-x-2">
                      <button
                        type="button"
                        onClick={() => handleEditClub(club)}
                        className="text-indigo-600 hover:text-indigo-900"
                      >
                        <Edit className="h-4 w-4" />
                      </button>
                      <button
                        type="button"
                        onClick={() => handleDeleteClub(club)}
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

        {clubs && clubs.length === 0 && (
          <div className="text-center py-12">
            <Building className="mx-auto h-12 w-12 text-gray-400" />
            <h3 className="mt-2 text-sm font-medium text-gray-900">Aucun club</h3>
            <p className="mt-1 text-sm text-gray-500">Commencez par créer un nouveau club.</p>
          </div>
        )}
      </div>

      <ClubModal
        isOpen={clubModalOpen}
        onClose={closeClubModal}
        club={editingClub}
        onSave={() => {}}
      />
    </div>
  );
}
