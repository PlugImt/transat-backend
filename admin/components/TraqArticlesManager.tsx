"use client";

import { Beer, Edit, Plus, Trash2 } from "lucide-react";
import Image from "next/image";
import { useCallback, useMemo, useState } from "react";
import toast from "react-hot-toast";
import DataTable, { type DataTableColumnDef } from "@/components/DataTable";
import { PageLoading } from "@/components/LoadingSpinner";
import { useDeleteTraqArticle, useTraqArticles } from "@/lib/hooks";
import { useAppStore } from "@/lib/stores/appStore";
import type { ApiError, TraqArticle } from "@/lib/types";

function StatusBadge({
  active,
  label,
  activeClass,
}: {
  active: boolean;
  label: string;
  activeClass: string;
}) {
  if (!active) return null;
  return (
    <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${activeClass}`}>
      {label}
    </span>
  );
}

export default function TraqArticlesManager() {
  const { data: articles = [], isLoading, error } = useTraqArticles();
  const deleteMutation = useDeleteTraqArticle();
  const { openTraqArticleModal } = useAppStore();
  const [showDisabledOnly, setShowDisabledOnly] = useState(false);

  const handleCreate = useCallback(() => {
    openTraqArticleModal();
  }, [openTraqArticleModal]);

  const handleEdit = useCallback(
    (article: TraqArticle) => {
      openTraqArticleModal(article);
    },
    [openTraqArticleModal],
  );

  const handleDelete = useCallback(
    async (article: TraqArticle) => {
      if (!confirm(`Êtes-vous sûr de vouloir supprimer « ${article.name} » ?`)) return;

      toast.promise(deleteMutation.mutateAsync(article.id_traq), {
        loading: "Suppression en cours...",
        success: "Article supprimé avec succès",
        error: (err: ApiError) =>
          err?.response?.data?.error || "Échec de la suppression de l'article",
      });
    },
    [deleteMutation],
  );

  const filteredArticles = useMemo(() => {
    if (!showDisabledOnly) return articles;
    return articles.filter((a) => a.disabled);
  }, [articles, showDisabledOnly]);

  const columns = useMemo<DataTableColumnDef<TraqArticle>[]>(
    () => [
      {
        id: "article",
        header: "Article",
        accessorFn: (row) => `${row.name} ${row.description} ${row.traq_type}`.toLowerCase(),
        cell: ({ row }) => {
          const article = row.original;
          return (
            <div className="flex items-center space-x-3">
              {article.picture ? (
                <Image
                  src={article.picture}
                  alt={article.name}
                  width={40}
                  height={40}
                  className="h-10 w-10 rounded object-cover shrink-0"
                />
              ) : (
                <div className="h-10 w-10 rounded bg-amber-100 flex items-center justify-center shrink-0">
                  <Beer className="h-5 w-5 text-amber-600" />
                </div>
              )}
              <div className="min-w-0">
                <p className="text-sm font-medium text-gray-900 truncate">{article.name}</p>
                <p className="text-xs text-gray-500 truncate max-w-xs">{article.description}</p>
              </div>
            </div>
          );
        },
      },
      {
        accessorKey: "traq_type",
        header: "Type",
        cell: ({ row }) => (
          <span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-amber-100 text-amber-800">
            {row.original.traq_type || "—"}
          </span>
        ),
      },
      {
        id: "prices",
        header: "Prix",
        accessorFn: (row) => row.price,
        cell: ({ row }) => {
          const { price, price_half } = row.original;
          return (
            <div className="text-sm text-gray-900">
              <div>{Number(price).toFixed(2)} €</div>
              {price_half > 0 && (
                <div className="text-xs text-gray-500">demi {Number(price_half).toFixed(2)} €</div>
              )}
            </div>
          );
        },
      },
      {
        accessorKey: "alcohol",
        header: "Alcool",
        cell: ({ row }) => (
          <span className="text-sm text-gray-900">{Number(row.original.alcohol).toFixed(1)} %</span>
        ),
      },
      {
        id: "status",
        header: "Statut",
        accessorFn: (row) =>
          [row.disabled && "disabled", row.limited && "limited", row.out_of_stock && "oos"]
            .filter(Boolean)
            .join(" "),
        cell: ({ row }) => {
          const a = row.original;
          return (
            <div className="flex flex-wrap gap-1">
              <StatusBadge
                active={a.disabled}
                label="Désactivé"
                activeClass="bg-gray-100 text-gray-700"
              />
              <StatusBadge
                active={a.limited}
                label="Limité"
                activeClass="bg-purple-100 text-purple-800"
              />
              <StatusBadge
                active={a.out_of_stock}
                label="Rupture"
                activeClass="bg-red-100 text-red-800"
              />
              {!a.disabled && !a.limited && !a.out_of_stock && (
                <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                  Dispo
                </span>
              )}
            </div>
          );
        },
      },
      {
        id: "actions",
        header: "Actions",
        enableSorting: false,
        cell: ({ row }) => {
          const article = row.original;
          return (
            <div className="flex items-center space-x-1">
              <button
                type="button"
                onClick={() => handleEdit(article)}
                className="p-2 text-gray-400 hover:text-amber-600 rounded-full hover:bg-gray-100"
                title="Modifier"
              >
                <Edit className="h-4 w-4" />
              </button>
              <button
                type="button"
                onClick={() => handleDelete(article)}
                className="p-2 text-gray-400 hover:text-red-600 rounded-full hover:bg-gray-100"
                title="Supprimer"
              >
                <Trash2 className="h-4 w-4" />
              </button>
            </div>
          );
        },
      },
    ],
    [handleEdit, handleDelete],
  );

  if (isLoading) {
    return <PageLoading text="Chargement des articles Traq..." />;
  }

  if (error) {
    return (
      <div className="bg-red-50 border border-red-200 rounded-md p-4">
        <div className="text-sm text-red-700">
          {(error as ApiError)?.message || "Échec de la récupération des articles"}
        </div>
      </div>
    );
  }

  return (
    <div>
      <div className="flex flex-col sm:flex-row sm:justify-between sm:items-center mb-4 gap-3">
        <div className="flex items-center space-x-2 text-gray-600">
          <Beer className="h-5 w-5" />
          <span className="text-sm">
            {showDisabledOnly
              ? `${filteredArticles.length} article${filteredArticles.length > 1 ? "s" : ""} désactivé${filteredArticles.length > 1 ? "s" : ""}`
              : `${articles.length} article${articles.length > 1 ? "s" : ""}`}
          </span>
        </div>
        <button
          type="button"
          onClick={handleCreate}
          className="inline-flex items-center justify-center px-4 py-2 bg-amber-600 text-white rounded-md hover:bg-amber-700 text-sm"
        >
          <Plus className="h-4 w-4 mr-2" />
          Créer un article
        </button>
      </div>

      <DataTable
        data={filteredArticles}
        columns={columns}
        searchPlaceholder="Rechercher par nom, description, type..."
        showFilter
        filterActive={showDisabledOnly}
        onFilterToggle={() => setShowDisabledOnly(!showDisabledOnly)}
        filterLabel="Désactivés uniquement"
        filterActiveLabel="Articles désactivés"
      />
    </div>
  );
}
