"use client";

import { useState } from "react";
import { Folder, FolderOpen, Plus, Edit, Trash2, MessageSquare } from "lucide-react";
import {
  useReservationTree,
  useDeleteReservationCategory,
  useDeleteReservationItem,
  useClubs,
} from "@/lib/hooks";
import type { ReservationTreeItem, ReservationItem } from "@/lib/api";
import { ReservationTreeModal } from "@/components/ReservationTreeModal";
import { CategoryModal } from "@/components/CategoryModal";
import { ItemModal } from "@/components/ItemModal";

export default function ReservationsPage() {
  const { data: tree, isLoading } = useReservationTree();
  const { data: clubs } = useClubs();
  const deleteCategory = useDeleteReservationCategory();
  const deleteItem = useDeleteReservationItem();

  const [expandedNodes, setExpandedNodes] = useState<Set<string>>(new Set());
  const [selectedCategory, setSelectedCategory] = useState<ReservationTreeItem | null>(null);
  const [selectedItem, setSelectedItem] = useState<ReservationItem | null>(null);
  const [categoryModalOpen, setCategoryModalOpen] = useState(false);
  const [itemModalOpen, setItemModalOpen] = useState(false);
  const [messageModalOpen, setMessageModalOpen] = useState(false);
  const [editingCategory, setEditingCategory] = useState<ReservationTreeItem | null>(null);
  const [editingItem, setEditingItem] = useState<ReservationItem | null>(null);
  const [parentCategoryId, setParentCategoryId] = useState<number | null>(null);
  const [parentClubId, setParentClubId] = useState<number | null>(null);

  const toggleNode = (nodeId: string) => {
    const newExpanded = new Set(expandedNodes);
    if (newExpanded.has(nodeId)) {
      newExpanded.delete(nodeId);
    } else {
      newExpanded.add(nodeId);
    }
    setExpandedNodes(newExpanded);
  };

  const handleCreateCategory = (parentId?: number, clubId?: number) => {
    setParentCategoryId(parentId || null);
    setParentClubId(clubId || null);
    setEditingCategory(null);
    setCategoryModalOpen(true);
  };

  const handleEditCategory = (category: ReservationTreeItem) => {
    setEditingCategory(category);
    setParentCategoryId(category.parent_id || null);
    setParentClubId(category.club_id || null);
    setCategoryModalOpen(true);
  };

  const handleDeleteCategory = async (id: number) => {
    if (!confirm("Êtes-vous sûr de vouloir supprimer cette catégorie ?")) {
      return;
    }
    try {
      await deleteCategory.mutateAsync(id);
    } catch (error: any) {
      alert(error?.response?.data?.error || "Erreur lors de la suppression");
    }
  };

  const handleCreateItem = (categoryId?: number, clubId?: number) => {
    setParentCategoryId(categoryId || null);
    setParentClubId(clubId || null);
    setEditingItem(null);
    setItemModalOpen(true);
  };

  const handleEditItem = (item: ReservationItem) => {
    setEditingItem(item);
    setItemModalOpen(true);
  };

  const handleDeleteItem = async (id: number) => {
    if (!confirm("Êtes-vous sûr de vouloir supprimer cet élément ?")) {
      return;
    }
    try {
      await deleteItem.mutateAsync(id);
    } catch (error: any) {
      alert(error?.response?.data?.error || "Erreur lors de la suppression");
    }
  };

  const handleEditMessages = (item: ReservationItem) => {
    setSelectedItem(item);
    setMessageModalOpen(true);
  };

  const renderTree = (nodes: ReservationTreeItem[], level = 0): JSX.Element[] => {
    return nodes.map((node) => {
      const nodeId = node.type === "category" ? `cat-${node.id}` : `club-${node.club_id}`;
      const isExpanded = expandedNodes.has(nodeId);
      const hasChildren = node.children && node.children.length > 0;
      const hasItems = node.items && node.items.length > 0;

      return (
        <div key={nodeId} className="select-none">
          <div
            className={`flex items-center gap-2 py-2 px-4 hover:bg-gray-100 ${
              level > 0 ? "ml-8" : ""
            }`}
            style={{ paddingLeft: `${level * 24 + 16}px` }}
          >
            {node.type === "category" && (
              <>
                <button
                  onClick={() => toggleNode(nodeId)}
                  className="flex items-center"
                  disabled={!hasChildren && !hasItems}
                >
                  {hasChildren || hasItems ? (
                    isExpanded ? (
                      <FolderOpen size={18} className="text-blue-600" />
                    ) : (
                      <Folder size={18} className="text-gray-500" />
                    )
                  ) : (
                    <div className="w-[18px]" />
                  )}
                </button>
                <span className="flex-1 font-medium">{node.name}</span>
                <span className="text-sm text-gray-500">({node.club_name})</span>
                <div className="flex gap-1">
                  <button
                    onClick={() => handleCreateCategory(node.id, node.club_id)}
                    className="p-1 text-blue-600 hover:bg-blue-50 rounded"
                    title="Ajouter une sous-catégorie"
                  >
                    <Plus size={16} />
                  </button>
                  <button
                    onClick={() => handleCreateItem(node.id, node.club_id)}
                    className="p-1 text-green-600 hover:bg-green-50 rounded"
                    title="Ajouter un élément"
                  >
                    <Plus size={16} />
                  </button>
                  <button
                    onClick={() => handleEditCategory(node)}
                    className="p-1 text-indigo-600 hover:bg-indigo-50 rounded"
                    title="Modifier"
                  >
                    <Edit size={16} />
                  </button>
                  <button
                    onClick={() => handleDeleteCategory(node.id!)}
                    className="p-1 text-red-600 hover:bg-red-50 rounded"
                    title="Supprimer"
                  >
                    <Trash2 size={16} />
                  </button>
                </div>
              </>
            )}

            {node.type === "club_items" && (
              <>
                <div className="w-[18px]" />
                <span className="flex-1 font-medium text-gray-600">
                  Éléments sans catégorie - {node.club_name}
                </span>
                <button
                  onClick={() => handleCreateItem(undefined, node.club_id)}
                  className="p-1 text-green-600 hover:bg-green-50 rounded"
                  title="Ajouter un élément"
                >
                  <Plus size={16} />
                </button>
              </>
            )}
          </div>

          {/* Render items */}
          {isExpanded && hasItems && (
            <div>
              {node.items!.map((item) => (
                <div
                  key={`item-${item.id}`}
                  className="flex items-center gap-2 py-2 px-4 hover:bg-gray-50"
                  style={{ paddingLeft: `${(level + 1) * 24 + 16}px` }}
                >
                  <div className="w-[18px]" />
                  <span className="flex-1">{item.name}</span>
                  {item.slot && (
                    <span className="text-xs bg-blue-100 text-blue-800 px-2 py-1 rounded">
                      Créneau
                    </span>
                  )}
                  <div className="flex gap-1">
                    <button
                      onClick={() => handleEditMessages(item)}
                      className="p-1 text-purple-600 hover:bg-purple-50 rounded"
                      title="Modifier les messages"
                    >
                      <MessageSquare size={16} />
                    </button>
                    <button
                      onClick={() => handleEditItem(item)}
                      className="p-1 text-indigo-600 hover:bg-indigo-50 rounded"
                      title="Modifier"
                    >
                      <Edit size={16} />
                    </button>
                    <button
                      onClick={() => handleDeleteItem(item.id)}
                      className="p-1 text-red-600 hover:bg-red-50 rounded"
                      title="Supprimer"
                    >
                      <Trash2 size={16} />
                    </button>
                  </div>
                </div>
              ))}
            </div>
          )}

          {/* Render children */}
          {isExpanded && hasChildren && (
            <div>{renderTree(node.children!, level + 1)}</div>
          )}
        </div>
      );
    });
  };

  if (isLoading) {
    return (
      <div className="p-6">
        <div className="animate-pulse space-y-4">
          <div className="h-8 bg-gray-200 rounded w-1/4"></div>
          <div className="space-y-2">
            {[1, 2, 3].map((i) => (
              <div key={i} className="h-12 bg-gray-200 rounded"></div>
            ))}
          </div>
        </div>
      </div>
    );
  }

  // Debug: log tree data
  if (tree) {
    console.log("Reservation tree data:", tree);
    console.log("Tree length:", tree.length);
  }

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold text-gray-900">Gestion des réservations</h1>
        <button
          onClick={() => handleCreateCategory()}
          className="inline-flex items-center px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
        >
          <Plus className="h-4 w-4 mr-2" />
          Nouvelle catégorie
        </button>
      </div>

      <div className="bg-white shadow rounded-lg overflow-hidden">
        {tree && tree.length > 0 ? (
          <div className="divide-y divide-gray-200">{renderTree(tree)}</div>
        ) : (
          <div className="text-center py-12 text-gray-500">
            <p>Aucune catégorie ou élément de réservation trouvé.</p>
            {tree && <p className="text-xs mt-2">Tree data: {JSON.stringify(tree)}</p>}
          </div>
        )}
      </div>

      <CategoryModal
        isOpen={categoryModalOpen}
        onClose={() => {
          setCategoryModalOpen(false);
          setEditingCategory(null);
          setParentCategoryId(null);
          setParentClubId(null);
        }}
        category={editingCategory}
        parentCategoryId={parentCategoryId}
        parentClubId={parentClubId}
        clubs={clubs || []}
      />

      <ItemModal
        isOpen={itemModalOpen}
        onClose={() => {
          setItemModalOpen(false);
          setEditingItem(null);
          setParentCategoryId(null);
          setParentClubId(null);
        }}
        item={editingItem}
        parentCategoryId={parentCategoryId}
        parentClubId={parentClubId}
        clubs={clubs || []}
        categories={tree || []}
      />

      {selectedItem && (
        <ReservationTreeModal
          itemId={selectedItem.id}
          itemName={selectedItem.name}
          isOpen={messageModalOpen}
          onClose={() => {
            setMessageModalOpen(false);
            setSelectedItem(null);
          }}
        />
      )}
    </div>
  );
}
