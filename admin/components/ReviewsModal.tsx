"use client";

import { Calendar, MessageSquare, Star, Trash2, User, X } from "lucide-react";
import Image from "next/image";
import { useCallback } from "react";
import toast from "react-hot-toast";
import { PageLoading } from "@/components/LoadingSpinner";
import { useDeleteMenuItemReview, useMenuItemReviews } from "@/lib/hooks";
import type { ApiError } from "@/lib/types";

interface ReviewsModalProps {
  isOpen: boolean;
  onClose: () => void;
  menuItemId: number;
  menuItemName: string;
}

export default function ReviewsModal({
  isOpen,
  onClose,
  menuItemId,
  menuItemName,
}: ReviewsModalProps) {
  const { data: reviews = [], isLoading, error } = useMenuItemReviews(menuItemId);
  const deleteReviewMutation = useDeleteMenuItemReview();

  const handleDeleteReview = useCallback(
    async (email: string, userName: string) => {
      if (!confirm(`Êtes-vous sûr de vouloir supprimer l'avis de ${userName} ?`)) return;

      try {
        await deleteReviewMutation.mutateAsync({ menuItemId, email });
        toast.success("Avis supprimé avec succès");
      } catch (err) {
        const error = err as ApiError;
        toast.error(error.response?.data?.error || "Erreur lors de la suppression de l'avis");
      }
    },
    [deleteReviewMutation, menuItemId],
  );

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

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-lg shadow-xl max-w-4xl w-full max-h-[90vh] overflow-hidden">
        {/* Header */}
        <div className="p-6 border-b border-gray-200">
          <div className="flex items-center justify-between">
            <div>
              <h2 className="text-xl font-semibold text-gray-900">
                Avis pour &quot;{menuItemName}&quot;
              </h2>
              <p className="text-sm text-gray-500 mt-1">{reviews.length} avis au total</p>
            </div>
            <button
              type="button"
              onClick={onClose}
              className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
            >
              <X className="h-5 w-5" />
            </button>
          </div>
        </div>

        {/* Content */}
        <div className="p-6 overflow-y-auto max-h-[calc(90vh-120px)]">
          {isLoading ? (
            <div className="flex justify-center py-8">
              <PageLoading />
            </div>
          ) : error ? (
            <div className="bg-red-50 border border-red-200 rounded-lg p-4">
              <h3 className="text-red-800 font-semibold">Erreur</h3>
              <p className="text-red-600">
                {error.message || "Erreur lors du chargement des avis"}
              </p>
            </div>
          ) : reviews.length === 0 ? (
            <div className="text-center py-8">
              <MessageSquare className="h-12 w-12 text-gray-400 mx-auto mb-4" />
              <h3 className="text-lg font-medium text-gray-900 mb-2">Aucun avis</h3>
              <p className="text-gray-500">Ce plat n&apos;a pas encore reçu d&apos;avis.</p>
            </div>
          ) : (
            <div className="space-y-4">
              {reviews.map((review) => (
                <div
                  key={`${review.email}-${review.date}`}
                  className="bg-gray-50 rounded-lg p-4 hover:bg-gray-100 transition-colors"
                >
                  <div className="flex items-start justify-between">
                    <div className="flex-1">
                      {/* User Info */}
                      <div className="flex items-center space-x-3 mb-3">
                        <div className="flex-shrink-0">
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
                          <div className="flex items-center space-x-1 text-sm text-gray-500 mt-1">
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

                      {/* Comment */}
                      {review.comment && (
                        <div className="bg-white rounded-lg p-3 border">
                          <p className="text-gray-700">{review.comment}</p>
                        </div>
                      )}
                    </div>

                    {/* Delete Button */}
                    <button
                      type="button"
                      onClick={() =>
                        handleDeleteReview(review.email, `${review.first_name} ${review.last_name}`)
                      }
                      disabled={deleteReviewMutation.isPending}
                      className="ml-4 p-2 text-red-600 hover:bg-red-50 rounded-lg transition-colors disabled:opacity-50"
                      title="Supprimer l'avis"
                    >
                      <Trash2 className="h-4 w-4" />
                    </button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Footer */}
        <div className="p-6 border-t border-gray-200">
          <div className="flex justify-end">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors"
            >
              Fermer
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
