"use client";

import { Folder, FolderOpen, FolderPlus, Pencil, Trash2 } from "lucide-react";
import { type ReactElement, useState } from "react";
import toast from "react-hot-toast";
import { CategoryModal } from "@/components/CategoryModal";
import Modal from "@/components/Modal";
import type { ReservationTreeItem } from "@/lib/api";
import { useClubs, useDeleteReservationCategory, useReservationTree } from "@/lib/hooks";
import { apiErrorMessage, EmptyState, LoadingButton, SkeletonRows } from "./ui";

interface CategoryTreeProps {
  onOpenItem: (itemId: number) => void;
}

/** Categories of the catalogue: create, rename and delete. Items are managed from the catalogue tab. */
export function CategoryTree({ onOpenItem }: CategoryTreeProps) {
  const { data: tree, isLoading } = useReservationTree();
  const { data: clubs } = useClubs();
  const deleteCategory = useDeleteReservationCategory();

  const [expanded, setExpanded] = useState<Set<string>>(new Set());
  const [editing, setEditing] = useState<ReservationTreeItem | null>(null);
  const [modalOpen, setModalOpen] = useState(false);
  const [parentCategoryId, setParentCategoryId] = useState<number | null>(null);
  const [parentClubId, setParentClubId] = useState<number | null>(null);
  const [deleting, setDeleting] = useState<ReservationTreeItem | null>(null);

  const toggle = (id: string) =>
    setExpanded((previous) => {
      const next = new Set(previous);
      if (!next.delete(id)) next.add(id);
      return next;
    });

  const openCreate = (parentId?: number, clubId?: number) => {
    setEditing(null);
    setParentCategoryId(parentId ?? null);
    setParentClubId(clubId ?? null);
    setModalOpen(true);
  };

  const openEdit = (node: ReservationTreeItem) => {
    setEditing(node);
    setParentCategoryId(node.parent_id ?? null);
    setParentClubId(node.club_id ?? null);
    setModalOpen(true);
  };

  const closeModal = () => {
    setModalOpen(false);
    setEditing(null);
    setParentCategoryId(null);
    setParentClubId(null);
  };

  const confirmDelete = async () => {
    if (!deleting?.id) return;
    try {
      await deleteCategory.mutateAsync(deleting.id);
      toast.success("Catégorie supprimée");
      setDeleting(null);
    } catch (error) {
      toast.error(apiErrorMessage(error, "Impossible de supprimer cette catégorie"));
    }
  };

  const renderNodes = (nodes: ReservationTreeItem[], level = 0): ReactElement[] =>
    nodes.map((node) => {
      const nodeId = node.type === "category" ? `cat-${node.id}` : `club-${node.club_id}`;
      const isOpen = expanded.has(nodeId);
      const hasChildren = (node.children?.length ?? 0) > 0;
      const hasItems = (node.items?.length ?? 0) > 0;
      const expandable = hasChildren || hasItems;

      return (
        <div key={nodeId}>
          <div
            className="flex items-center gap-2 px-4 py-2 hover:bg-gray-50"
            style={{ paddingLeft: `${level * 24 + 16}px` }}
          >
            <button
              type="button"
              onClick={() => toggle(nodeId)}
              disabled={!expandable}
              className="flex items-center"
              aria-label={isOpen ? "Replier" : "Déplier"}
            >
              {expandable ? (
                isOpen ? (
                  <FolderOpen className="h-[18px] w-[18px] text-blue-600" />
                ) : (
                  <Folder className="h-[18px] w-[18px] text-gray-500" />
                )
              ) : (
                <span className="inline-block w-[18px]" />
              )}
            </button>
            {node.type === "category" ? (
              <>
                <span className="flex-1 font-medium text-gray-900">{node.name}</span>
                <span className="text-sm text-gray-500">{node.club_name}</span>
                <button
                  type="button"
                  onClick={() => openCreate(node.id, node.club_id)}
                  className="rounded p-1 text-blue-600 hover:bg-blue-50"
                  title="Ajouter une sous-catégorie"
                  aria-label="Ajouter une sous-catégorie"
                >
                  <FolderPlus className="h-4 w-4" />
                </button>
                <button
                  type="button"
                  onClick={() => openEdit(node)}
                  className="rounded p-1 text-indigo-600 hover:bg-indigo-50"
                  title="Renommer"
                  aria-label="Renommer"
                >
                  <Pencil className="h-4 w-4" />
                </button>
                <button
                  type="button"
                  onClick={() => setDeleting(node)}
                  className="rounded p-1 text-red-600 hover:bg-red-50"
                  title="Supprimer"
                  aria-label="Supprimer"
                >
                  <Trash2 className="h-4 w-4" />
                </button>
              </>
            ) : (
              <span className="flex-1 font-medium text-gray-600">
                Éléments sans catégorie · {node.club_name}
              </span>
            )}
          </div>

          {isOpen && hasItems && (
            <ul>
              {node.items?.map((item) => (
                <li
                  key={item.id}
                  className="flex items-center gap-2 px-4 py-1.5 text-sm"
                  style={{ paddingLeft: `${(level + 1) * 24 + 40}px` }}
                >
                  <span className="flex-1 text-gray-700">{item.name}</span>
                  <button
                    type="button"
                    onClick={() => onOpenItem(item.id)}
                    className="text-xs text-blue-600 hover:underline"
                  >
                    Détail
                  </button>
                </li>
              ))}
            </ul>
          )}
          {isOpen && hasChildren && renderNodes(node.children ?? [], level + 1)}
        </div>
      );
    });

  return (
    <div className="space-y-4">
      <div className="flex justify-end">
        <button
          type="button"
          onClick={() => openCreate()}
          className="inline-flex items-center gap-2 rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
        >
          <FolderPlus className="h-4 w-4" />
          Nouvelle catégorie
        </button>
      </div>

      <div className="overflow-hidden rounded-lg bg-white shadow">
        {isLoading ? (
          <SkeletonRows rows={5} columns={3} />
        ) : tree && tree.length > 0 ? (
          <div className="divide-y divide-gray-100">{renderNodes(tree)}</div>
        ) : (
          <EmptyState
            icon={<Folder className="h-10 w-10" />}
            title="Aucune catégorie"
            description="Créez une catégorie pour organiser les salles et le matériel d'un club."
          />
        )}
      </div>

      <CategoryModal
        isOpen={modalOpen}
        onClose={closeModal}
        category={editing}
        parentCategoryId={parentCategoryId}
        parentClubId={parentClubId}
        clubs={clubs ?? []}
      />

      <Modal
        isOpen={deleting !== null}
        onClose={() => setDeleting(null)}
        title="Supprimer la catégorie"
      >
        <div className="space-y-4 text-sm text-gray-700">
          <p>
            Supprimer la catégorie <strong>{deleting?.name}</strong> ? Elle doit être vide (ni
            sous-catégorie, ni élément).
          </p>
          <div className="flex justify-end gap-2">
            <LoadingButton
              onClick={() => setDeleting(null)}
              className="border border-gray-300 text-gray-700 hover:bg-gray-50"
            >
              Annuler
            </LoadingButton>
            <LoadingButton
              onClick={confirmDelete}
              loading={deleteCategory.isPending}
              className="bg-red-600 text-white hover:bg-red-700"
            >
              Supprimer
            </LoadingButton>
          </div>
        </div>
      </Modal>
    </div>
  );
}
