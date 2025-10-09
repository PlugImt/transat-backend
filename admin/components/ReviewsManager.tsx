"use client";

import { Calendar, MessageSquare, Star, User, X } from "lucide-react";
import Image from "next/image";
import { useMemo, useState } from "react";
import Combobox from "@/components/Combobox";
import { PageLoading } from "@/components/LoadingSpinner";
import { useAllReviews, useUsers } from "@/lib/hooks";

export default function ReviewsManager() {
  const [selectedUserEmail, setSelectedUserEmail] = useState<string>("");
  const [selectedDishName, setSelectedDishName] = useState<string>("");

  const { data: allReviews = [], isLoading: reviewsLoading, error } = useAllReviews();
  const { data: users = [] } = useUsers();

  // Filtrer les avis selon les critères sélectionnés
  const filteredReviews = useMemo(() => {
    let filtered = allReviews;

    if (selectedUserEmail) {
      filtered = filtered.filter((review) => review.email === selectedUserEmail);
    }

    if (selectedDishName) {
      filtered = filtered.filter((review) => review.dish_name === selectedDishName);
    }

    return filtered;
  }, [allReviews, selectedUserEmail, selectedDishName]);

  const renderStars = (rating: number) => {
    return (
      <div className="flex">
        {[1, 2, 3, 4, 5].map((star) => (
          <Star
            key={star}
            className={`h-4 w-4 ${
              star <= rating ? "text-yellow-400 fill-current" : "text-gray-300"
            }`}
          />
        ))}
      </div>
    );
  };

  // Obtenir les utilisateurs qui ont au moins un avis
  const usersWithReviews = useMemo(() => {
    const userEmails = new Set(allReviews.map((review) => review.email));
    return users.filter((user) => userEmails.has(user.email));
  }, [users, allReviews]);

  // Obtenir la liste unique des plats avec des avis
  const dishesWithReviews = useMemo(() => {
    const dishMap = new Map<string, string>();
    allReviews.forEach((review) => {
      if (review.dish_name) {
        dishMap.set(review.dish_name, review.dish_name);
      }
    });
    return Array.from(dishMap.values()).sort();
  }, [allReviews]);

  // Options pour les combobox
  const userOptions = useMemo(
    () =>
      usersWithReviews.map((user) => ({
        value: user.email,
        label: `${user.first_name} ${user.last_name} (${user.email})`,
      })),
    [usersWithReviews],
  );

  const dishOptions = useMemo(
    () =>
      dishesWithReviews.map((dish) => ({
        value: dish,
        label: dish,
      })),
    [dishesWithReviews],
  );

  return (
    <>
      {/* Filtres */}
      <div className="bg-white rounded-lg shadow p-6 mb-6">
        <h2 className="text-lg font-semibold mb-4">Filtres</h2>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {/* Filtre par utilisateur */}
          <Combobox
            options={userOptions}
            value={selectedUserEmail}
            onChange={setSelectedUserEmail}
            label="Filtrer par utilisateur"
            placeholder="Tous les utilisateurs"
            searchPlaceholder="Rechercher un utilisateur..."
            emptyText="Aucun utilisateur trouvé"
          />

          {/* Filtre par plat */}
          <Combobox
            options={dishOptions}
            value={selectedDishName}
            onChange={setSelectedDishName}
            label="Filtrer par plat"
            placeholder="Tous les plats"
            searchPlaceholder="Rechercher un plat..."
            emptyText="Aucun plat trouvé"
          />
        </div>

        {/* Indicateurs de filtres actifs */}
        {(selectedUserEmail || selectedDishName) && (
          <div className="mt-4 flex flex-wrap gap-2">
            {selectedUserEmail && (
              <span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                Utilisateur:{" "}
                {usersWithReviews.find((u) => u.email === selectedUserEmail)?.first_name}{" "}
                {usersWithReviews.find((u) => u.email === selectedUserEmail)?.last_name}
                <button
                  type="button"
                  onClick={() => setSelectedUserEmail("")}
                  className="ml-1 text-blue-600 hover:text-blue-800"
                >
                  <X className="h-3 w-3" />
                </button>
              </span>
            )}
            {selectedDishName && (
              <span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-green-100 text-green-800">
                Plat: "{selectedDishName}"
                <button
                  type="button"
                  onClick={() => setSelectedDishName("")}
                  className="ml-1 text-green-600 hover:text-green-800"
                >
                  <X className="h-3 w-3" />
                </button>
              </span>
            )}
          </div>
        )}
      </div>

      {/* Liste des avis */}
      <div className="bg-white rounded-lg shadow">
        <div className="px-6 py-4 border-b border-gray-200">
          <h2 className="text-lg font-semibold">Avis ({filteredReviews.length})</h2>
        </div>

        <div className="p-6">
          {reviewsLoading ? (
            <div className="flex justify-center py-8">
              <PageLoading />
            </div>
          ) : error ? (
            <div className="bg-red-50 border border-red-200 rounded-lg p-4">
              <h3 className="text-red-800 font-semibold">Erreur</h3>
              <p className="text-red-600">Erreur lors du chargement des avis</p>
            </div>
          ) : filteredReviews.length === 0 ? (
            <div className="text-center py-8">
              <MessageSquare className="h-12 w-12 text-gray-400 mx-auto mb-4" />
              <h3 className="text-lg font-medium text-gray-900 mb-2">
                {selectedUserEmail || selectedDishName ? "Aucun avis trouvé" : "Aucun avis"}
              </h3>
              <p className="text-gray-500">
                {selectedUserEmail || selectedDishName
                  ? "Aucun avis ne correspond aux critères de recherche."
                  : "Aucun avis n'a encore été laissé."}
              </p>
            </div>
          ) : (
            <div className="space-y-4">
              {filteredReviews.map((review, index: number) => (
                <div
                  key={`${review.email}-${review.dish_id || "unknown"}-${index}`}
                  className="bg-gray-50 rounded-lg p-4 hover:bg-gray-100 transition-colors"
                >
                  <div className="flex items-start justify-between">
                    <div className="flex-1">
                      {/* En-tête de l'avis */}
                      <div className="flex items-center justify-between mb-3">
                        <div className="flex items-center space-x-3">
                          <div className="shrink-0">
                            {review.profile_picture ? (
                              <Image
                                src={review.profile_picture}
                                alt={`${review.first_name} ${review.last_name}`}
                                width={40}
                                height={40}
                                className="h-10 w-10 rounded-full object-cover"
                              />
                            ) : (
                              <div className="h-10 w-10 bg-gray-300 rounded-full flex items-center justify-center">
                                <User className="h-6 w-6 text-gray-600" />
                              </div>
                            )}
                          </div>
                          <div className="flex-1">
                            <div className="flex items-center space-x-2">
                              <h4 className="font-medium text-gray-900">
                                {review.first_name} {review.last_name}
                              </h4>
                              {renderStars(review.note)}
                            </div>
                            <p className="text-sm text-gray-500">{review.email}</p>
                          </div>
                        </div>
                        <div className="text-right">
                          <div className="flex items-center space-x-1 text-sm text-gray-500">
                            <Calendar className="h-3 w-3" />
                            <span>
                              {new Date(review.date).toLocaleDateString("fr-FR", {
                                year: "numeric",
                                month: "long",
                                day: "numeric",
                                hour: "2-digit",
                                minute: "2-digit",
                              })}
                            </span>
                          </div>
                        </div>
                      </div>

                      {/* Plat */}
                      {review.dish_name && (
                        <div className="mb-3">
                          <span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                            {review.dish_name}
                          </span>
                        </div>
                      )}

                      {/* Commentaire */}
                      {review.comment && (
                        <div className="bg-white rounded-lg p-3 border">
                          <p className="text-gray-700">{review.comment}</p>
                        </div>
                      )}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </>
  );
}
