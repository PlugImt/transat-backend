"use client";

import React, { useState, useMemo, useCallback } from "react";
import { type ColumnDef } from "@tanstack/react-table";
import {
  UtensilsCrossed,
  Trash2,
  Star,
  Calendar,
  ChefHat,
  TrendingUp,
  Clock,
  MessageSquare,
} from "lucide-react";
import { MenuItem, ApiError } from "@/lib/types";
import { useMenuItems, useDeleteMenuItem } from "@/lib/hooks";
import { PageLoading } from "@/components/LoadingSpinner";
import ErrorBoundary from "@/components/ErrorBoundary";
import DataTable from "@/components/DataTable";
import { useAppStore } from "@/lib/stores/appStore";
import ReviewsModal from "@/components/ReviewsModal";
import toast from "react-hot-toast";

function MenuPageContent() {
  const { setCurrentPage } = useAppStore();
  const [reviewsModalOpen, setReviewsModalOpen] = useState(false);
  const [selectedMenuItem, setSelectedMenuItem] = useState<MenuItem | null>(
    null
  );

  const { data: menuItems = [], isLoading, error } = useMenuItems();
  const deleteMenuItemMutation = useDeleteMenuItem();

  React.useEffect(() => {
    setCurrentPage("Menu du RU");
  }, [setCurrentPage]);

  const handleDeleteMenuItem = useCallback(
    async (id: number, name: string) => {
      if (!confirm(`Êtes-vous sûr de vouloir supprimer "${name}" ?`)) return;

      try {
        await deleteMenuItemMutation.mutateAsync(id);
        toast.success("Élément du menu supprimé avec succès");
      } catch (err) {
        const error = err as ApiError;
        toast.error(
          error.response?.data?.error ||
            "Erreur lors de la suppression de l'élément du menu"
        );
      }
    },
    [deleteMenuItemMutation]
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
                    star <= row.original.average_rating
                      ? "text-yellow-400"
                      : "text-gray-300"
                  }`}
                >
                  ★
                </span>
              ))}
            </div>
            <span>
              {row.original.average_rating > 0
                ? row.original.average_rating.toFixed(1)
                : "N/A"}
            </span>
            <span className="text-gray-500 text-sm">
              ({row.original.total_ratings} avis)
            </span>
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
            <span>
              {new Date(row.original.first_time_served).toLocaleDateString(
                "fr-FR"
              )}
            </span>
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
              onClick={() => handleViewReviews(row.original)}
              className="p-1 text-blue-600 hover:bg-blue-50 rounded"
              title="Voir les avis"
            >
              <MessageSquare size={16} />
            </button>
            <button
              onClick={() =>
                handleDeleteMenuItem(
                  row.original.id_restaurant_articles,
                  row.original.name
                )
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
    [handleDeleteMenuItem, handleViewReviews, deleteMenuItemMutation.isPending]
  );

  if (isLoading) return <PageLoading />;

  if (error) {
    return (
      <div className="p-6">
        <div className="bg-red-50 border border-red-200 rounded-lg p-4">
          <h2 className="text-red-800 font-semibold">Erreur</h2>
          <p className="text-red-600">
            {error.message || "Erreur lors du chargement des données"}
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center space-x-3">
          <UtensilsCrossed className="h-8 w-8 text-orange-500" />
          <div>
            <h1 className="text-2xl font-bold text-gray-900">Menu du RU</h1>
            <p className="text-gray-600">
              Gestion des plats du restaurant universitaire
            </p>
          </div>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
        <div className="bg-white p-6 rounded-lg shadow-sm border">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Total plats</p>
              <p className="text-3xl font-bold text-gray-900">
                {menuItems.length}
              </p>
            </div>
            <ChefHat className="h-12 w-12 text-orange-500" />
          </div>
        </div>

        <div className="bg-white p-6 rounded-lg shadow-sm border">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Note moyenne</p>
              <p className="text-3xl font-bold text-gray-900">
                {menuItems.length > 0
                  ? (
                      menuItems
                        .filter((item) => item.total_ratings > 0)
                        .reduce((acc, item) => acc + item.average_rating, 0) /
                      menuItems.filter((item) => item.total_ratings > 0).length
                    ).toFixed(1)
                  : "N/A"}
              </p>
            </div>
            <Star className="h-12 w-12 text-yellow-500" />
          </div>
        </div>

        <div className="bg-white p-6 rounded-lg shadow-sm border">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Total avis</p>
              <p className="text-3xl font-bold text-gray-900">
                {menuItems.reduce((acc, item) => acc + item.total_ratings, 0)}
              </p>
            </div>
            <TrendingUp className="h-12 w-12 text-blue-500" />
          </div>
        </div>

        <div className="bg-white p-6 rounded-lg shadow-sm border">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">
                Plats bien notés
              </p>
              <p className="text-3xl font-bold text-gray-900">
                {menuItems.filter((item) => item.average_rating >= 4).length}
              </p>
            </div>
            <Star className="h-12 w-12 text-green-500" />
          </div>
        </div>
      </div>

      {/* Table */}
      <div className="bg-white rounded-lg shadow-sm border p-2">
        <DataTable
          columns={columns}
          data={menuItems}
          globalFilterColumn="name"
          searchPlaceholder="Rechercher un plat..."
        />
      </div>

      {/* Reviews Modal */}
      {selectedMenuItem && (
        <ReviewsModal
          isOpen={reviewsModalOpen}
          onClose={handleCloseReviewsModal}
          menuItemId={selectedMenuItem.id_restaurant_articles}
          menuItemName={selectedMenuItem.name}
        />
      )}
    </div>
  );
}

export default function MenuPage() {
  return (
    <ErrorBoundary>
      <MenuPageContent />
    </ErrorBoundary>
  );
}
