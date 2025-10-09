"use client";

import type { ColumnDef } from "@tanstack/react-table";
import {
  Calendar,
  ChefHat,
  Clock,
  MessageSquare,
  Trash2,
  TrendingUp,
  UtensilsCrossed,
} from "lucide-react";
import { useCallback, useMemo, useState } from "react";
import toast from "react-hot-toast";
import DataTable from "@/components/DataTable";
import { PageLoading } from "@/components/LoadingSpinner";
import ReviewsModal from "@/components/ReviewsModal";
import { useDeleteMenuItem, useMenuItems } from "@/lib/hooks";
import type { ApiError, MenuItem } from "@/lib/types";

export default function MenuManager() {
  const [reviewsModalOpen, setReviewsModalOpen] = useState(false);
  const [selectedMenuItem, setSelectedMenuItem] = useState<MenuItem | null>(null);

  const { data: menuItems = [], isLoading, error } = useMenuItems();
  const deleteMenuItemMutation = useDeleteMenuItem();

  const handleDeleteMenuItem = useCallback(
    async (id: number, name: string) => {
      if (!confirm(`Êtes-vous sûr de vouloir supprimer "${name}" ?`)) return;

      try {
        await deleteMenuItemMutation.mutateAsync(id);
        toast.success("Élément du menu supprimé avec succès");
      } catch (err) {
        const error = err as ApiError;
        toast.error(
          error.response?.data?.error || "Erreur lors de la suppression de l'élément du menu",
        );
      }
    },
    [deleteMenuItemMutation],
  );

  const handleViewReviews = useCallback((item: MenuItem) => {
    setSelectedMenuItem(item);
    setReviewsModalOpen(true);
  }, []);

  const handleCloseReviewsModal = useCallback(() => {
    setReviewsModalOpen(false);
    setSelectedMenuItem(null);
  }, []);

  const columns = useMemo<ColumnDef<MenuItem>[]>(
    () => [
      {
        accessorKey: "name",
        header: "Nom du plat",
        cell: ({ row }) => (
          <div className="flex items-center space-x-2">
            <ChefHat className="h-4 w-4 text-orange-500" />
            <span className="font-medium">{row.original.name}</span>
          </div>
        ),
      },
      {
        accessorKey: "average_rating",
        header: "Note moyenne",
        cell: ({ row }) => (
          <div className="flex items-center space-x-1">
            <div className="flex">
              {[1, 2, 3, 4, 5].map((star) => (
                <span
                  key={star}
                  className={`text-sm ${
                    star <= row.original.average_rating ? "text-yellow-400" : "text-gray-300"
                  }`}
                >
                  ★
                </span>
              ))}
            </div>
            <span>
              {row.original.average_rating > 0 ? row.original.average_rating.toFixed(1) : "N/A"}
            </span>
          </div>
        ),
      },
      {
        accessorKey: "total_ratings",
        header: "Nombre d'avis",
        cell: ({ row }) => (
          <div className="flex items-center space-x-1">
            <MessageSquare className="h-4 w-4 text-blue-500" />
            <span>{row.original.total_ratings}</span>
          </div>
        ),
      },
      {
        accessorKey: "times_served",
        header: "Fois servi",
        cell: ({ row }) => (
          <div className="flex items-center space-x-1">
            <TrendingUp className="h-4 w-4 text-blue-500" />
            <span>{row.original.times_served}</span>
          </div>
        ),
      },
      {
        accessorKey: "first_time_served",
        header: "Première fois",
        cell: ({ row }) => (
          <div className="flex items-center space-x-1">
            <Calendar className="h-4 w-4 text-green-500" />
            <span>{new Date(row.original.first_time_served).toLocaleDateString("fr-FR")}</span>
          </div>
        ),
      },
      {
        accessorKey: "last_served",
        header: "Dernière fois",
        cell: ({ row }) => (
          <div className="flex items-center space-x-1">
            <Clock className="h-4 w-4 text-purple-500" />
            <span>
              {row.original.last_served
                ? new Date(row.original.last_served).toLocaleDateString("fr-FR")
                : "N/A"}
            </span>
          </div>
        ),
      },
      {
        id: "actions",
        header: "Actions",
        cell: ({ row }) => (
          <div className="flex space-x-2">
            <button
              type="button"
              onClick={() => handleViewReviews(row.original)}
              className="p-1 text-blue-600 hover:bg-blue-50 rounded"
              title="Voir les avis"
            >
              <MessageSquare size={16} />
            </button>
            <button
              type="button"
              onClick={() =>
                handleDeleteMenuItem(row.original.id_restaurant_articles, row.original.name)
              }
              disabled={deleteMenuItemMutation.isPending}
              className="p-1 text-red-600 hover:bg-red-50 rounded disabled:opacity-50"
              title="Supprimer"
            >
              <Trash2 size={16} />
            </button>
          </div>
        ),
      },
    ],
    [handleDeleteMenuItem, handleViewReviews, deleteMenuItemMutation.isPending],
  );

  if (isLoading) {
    return <PageLoading text="Chargement du menu..." />;
  }

  if (error) {
    return (
      <div className="bg-red-50 border border-red-200 rounded-md p-4">
        <div className="text-sm text-red-700">
          {(error as ApiError)?.message || "Échec de la récupération du menu"}
        </div>
      </div>
    );
  }

  return (
    <>
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center space-x-2 text-gray-600">
          <UtensilsCrossed className="h-5 w-5" />
          <span className="text-sm sm:text-base">
            {menuItems.length} plat{menuItems.length > 1 ? "s" : ""}
          </span>
        </div>
      </div>

      <DataTable
        data={menuItems}
        columns={columns}
        searchPlaceholder="Rechercher par nom de plat..."
        className="mb-6"
      />

      {/* Reviews Modal */}
      {selectedMenuItem && (
        <ReviewsModal
          isOpen={reviewsModalOpen}
          onClose={handleCloseReviewsModal}
          menuItemId={selectedMenuItem.id_restaurant_articles}
          menuItemName={selectedMenuItem.name}
        />
      )}
    </>
  );
}
