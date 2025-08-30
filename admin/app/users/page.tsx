"use client";

import { useState, useMemo, useCallback } from "react";
import { type ColumnDef } from "@tanstack/react-table";
import {
  Users,
  Mail,
  Check,
  X,
  Edit,
  Trash2,
  Plus,
  CheckCircle,
  Phone,
} from "lucide-react";
import { User, ApiError } from "@/lib/types";
import { useUsers, useDeleteUser, useValidateUser } from "@/lib/hooks";
import LanguageFlag from "@/components/LanguageFlag";
import { PageLoading } from "@/components/LoadingSpinner";
import ErrorBoundary from "@/components/ErrorBoundary";
import DataTable from "@/components/DataTable";
import { UserModal } from "@/components/LazyComponents";
import { useAppStore } from "@/lib/stores/appStore";
import toast from "react-hot-toast";

function UsersPageContent() {
  const { userModalOpen, editingUser, openUserModal, closeUserModal } =
    useAppStore();
  const [showUnverifiedOnly, setShowUnverifiedOnly] = useState(false);

  const { data: users = [], isLoading, error } = useUsers();
  const deleteUserMutation = useDeleteUser();
  const validateUserMutation = useValidateUser();

  const handleCreateUser = useCallback(() => {
    openUserModal();
  }, [openUserModal]);

  const handleEditUser = useCallback(
    (user: User) => {
      openUserModal(user);
    },
    [openUserModal]
  );

  const handleDeleteUser = useCallback(
    async (email: string) => {
      if (!confirm("Êtes-vous sûr de vouloir supprimer cet utilisateur ?"))
        return;

      toast.promise(deleteUserMutation.mutateAsync(email), {
        loading: "Suppression en cours...",
        success: "Utilisateur supprimé avec succès !",
        error: (err: ApiError) =>
          err?.response?.data?.error ||
          "Échec de la suppression de l'utilisateur",
      });
    },
    [deleteUserMutation]
  );

  const handleValidateUser = useCallback(
    async (email: string) => {
      if (
        !confirm(
          "Êtes-vous sûr de vouloir valider cet utilisateur ? Cela changera son rôle VERIFYING en NEWF."
        )
      )
        return;

      toast.promise(validateUserMutation.mutateAsync(email), {
        loading: "Validation en cours...",
        success: "Utilisateur validé avec succès !",
        error: (err: ApiError) =>
          err?.response?.data?.error ||
          "Échec de la validation de l'utilisateur",
      });
    },
    [validateUserMutation]
  );

  // Filter users based on verification status
  const filteredUsers = useMemo(() => {
    if (!showUnverifiedOnly) {
      return users;
    }
    return users.filter((user: User) => user.roles?.includes("VERIFYING"));
  }, [users, showUnverifiedOnly]);

  // Define columns for the DataTable
  const columns = useMemo<ColumnDef<User>[]>(
    () => [
      {
        id: "user_info",
        header: "Utilisateur",
        accessorFn: (row) => {
          // Cette fonction définit quelle donnée sera utilisée pour le tri et le filtrage
          return `${row.first_name || ""} ${row.last_name || ""} ${row.email} ${
            row.phone_number || ""
          }`.toLowerCase();
        },
        cell: ({ row }) => {
          const user = row.original;
          return (
            <div className="flex items-center space-x-3">
              <div className="flex-shrink-0">
                <div className="h-10 w-10 rounded-full bg-gray-300 flex items-center justify-center">
                  <span className="text-sm font-medium text-gray-700">
                    {user.first_name?.[0]}
                    {user.last_name?.[0]}
                  </span>
                </div>
              </div>
              <div>
                <div className="flex items-center space-x-2">
                  <p className="text-sm font-medium text-gray-900">
                    {user.first_name} {user.last_name}
                  </p>
                  {user.roles?.includes("VERIFYING") ? (
                    <div title="Non vérifié">
                      <X className="h-4 w-4 text-red-500" />
                    </div>
                  ) : (
                    <div title="Vérifié">
                      <Check className="h-4 w-4 text-green-500" />
                    </div>
                  )}
                </div>
                <div className="flex flex-col space-y-1">
                  <div className="flex items-center space-x-1 text-sm text-gray-500">
                    <Mail className="h-4 w-4" />
                    <span>{user.email}</span>
                  </div>
                  {user.phone_number && (
                    <div className="flex items-center space-x-1 text-sm text-gray-500">
                      <Phone className="h-4 w-4" />
                      <span>{user.phone_number}</span>
                    </div>
                  )}
                </div>
              </div>
            </div>
          );
        },
        sortingFn: (a, b) => {
          const nameA = `${a.original.first_name || ""} ${
            a.original.last_name || ""
          }`.toLowerCase();
          const nameB = `${b.original.first_name || ""} ${
            b.original.last_name || ""
          }`.toLowerCase();
          return nameA.localeCompare(nameB);
        },
      },
      {
        accessorKey: "formation_name",
        header: "Formation",
        cell: ({ row }) => (
          <span className="text-sm text-gray-900">
            {row.original.formation_name || "-"}
          </span>
        ),
      },
      {
        accessorKey: "language",
        header: "Langue",
        cell: ({ row }) => {
          const user = row.original;
          return user.language ? (
            <div className="flex items-center space-x-1">
              <LanguageFlag
                languageCode={user.language}
                className="text-base"
              />
              <span className="text-sm text-gray-900">
                {user.language.toUpperCase()}
              </span>
            </div>
          ) : (
            <span className="text-sm text-gray-500">-</span>
          );
        },
      },
      {
        id: "roles",
        header: "Rôles",
        accessorFn: (row) => {
          // Cette fonction définit quelle donnée sera utilisée pour le tri et le filtrage
          return (row.roles || []).join(" ").toLowerCase();
        },
        cell: ({ row }) => {
          const roles = row.original.roles || [];
          return roles.length > 0 ? (
            <div className="grid grid-cols-2 gap-1">
              {roles.map((role: string, index: number) => (
                <span
                  key={index}
                  className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800"
                >
                  {role}
                </span>
              ))}
            </div>
          ) : (
            <span className="text-sm text-gray-500">-</span>
          );
        },
        sortingFn: (a, b) => {
          const rolesA = (a.original.roles || []).join(", ").toLowerCase();
          const rolesB = (b.original.roles || []).join(", ").toLowerCase();
          return rolesA.localeCompare(rolesB);
        },
      },
      {
        id: "actions",
        header: "Actions",
        cell: ({ row }) => {
          const user = row.original;
          return (
            <div className="flex items-center space-x-1">
              <button
                onClick={() => handleEditUser(user)}
                className="p-2 text-gray-400 hover:text-blue-600 rounded-full hover:bg-gray-100"
                title="Modifier l'utilisateur"
              >
                <Edit className="h-4 w-4" />
              </button>
              {user.roles?.includes("VERIFYING") && (
                <button
                  onClick={() => handleValidateUser(user.email)}
                  className="p-2 text-gray-400 hover:text-green-600 rounded-full hover:bg-gray-100"
                  title="Valider l'utilisateur (VERIFYING → NEWF)"
                >
                  <CheckCircle className="h-4 w-4" />
                </button>
              )}
              <button
                onClick={() => handleDeleteUser(user.email)}
                className="p-2 text-gray-400 hover:text-red-600 rounded-full hover:bg-gray-100"
                title="Supprimer l'utilisateur"
              >
                <Trash2 className="h-4 w-4" />
              </button>
            </div>
          );
        },
        enableSorting: false,
      },
    ],
    [handleEditUser, handleValidateUser, handleDeleteUser]
  );

  if (isLoading) {
    return (
      <div className="p-4 sm:p-6 pt-16 lg:pt-6">
        <div className="flex flex-col sm:flex-row sm:justify-between sm:items-center mb-6 space-y-4 sm:space-y-0">
          <h1 className="text-xl sm:text-2xl font-bold text-gray-900">
            Gestion des utilisateurs
          </h1>
        </div>
        <PageLoading text="Chargement des utilisateurs..." />
      </div>
    );
  }

  if (error) {
    return (
      <div className="p-6">
        <div className="bg-red-50 border border-red-200 rounded-md p-4">
          <div className="text-sm text-red-700">
            {(error as ApiError)?.message ||
              "Échec de la récupération des utilisateurs"}
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="p-4 sm:p-6 pt-16 lg:pt-6">
      <div className="flex flex-col sm:flex-row sm:justify-between sm:items-center mb-6 space-y-4 sm:space-y-0">
        <h1 className="text-xl sm:text-2xl font-bold text-gray-900">
          Gestion des utilisateurs
        </h1>
        <div className="flex flex-col sm:flex-row sm:items-center space-y-2 sm:space-y-0 sm:space-x-4">
          <div className="flex items-center space-x-2 text-gray-600">
            <Users className="h-5 w-5" />
            <span className="text-sm sm:text-base">
              {showUnverifiedOnly
                ? `${filteredUsers.length} utilisateur${
                    filteredUsers.length > 1 ? "s" : ""
                  } non validé${filteredUsers.length > 1 ? "s" : ""}`
                : `${users.length} utilisateurs`}
            </span>
          </div>
          <button
            onClick={handleCreateUser}
            className="inline-flex items-center justify-center px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 text-sm sm:text-base"
          >
            <Plus className="h-4 w-4 mr-2" />
            Créer un utilisateur
          </button>
        </div>
      </div>

      <DataTable
        data={filteredUsers}
        columns={columns}
        searchPlaceholder="Rechercher par nom, email, formation..."
        className="mb-6"
        showFilter={true}
        filterActive={showUnverifiedOnly}
        onFilterToggle={() => setShowUnverifiedOnly(!showUnverifiedOnly)}
        filterLabel="Non validés uniquement"
        filterActiveLabel="Utilisateurs non validés (VERIFYING)"
      />

      <UserModal
        isOpen={userModalOpen}
        onClose={closeUserModal}
        user={editingUser}
        onSave={() => {}}
      />
    </div>
  );
}

export default function UsersPage() {
  return (
    <ErrorBoundary>
      <UsersPageContent />
    </ErrorBoundary>
  );
}
