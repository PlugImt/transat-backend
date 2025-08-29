"use client";

import { useState, useMemo, useEffect, useCallback } from "react";
import {
  Users,
  Mail,
  Check,
  X,
  Edit,
  Trash2,
  Plus,
  Search,
  ChevronUp,
  ChevronDown,
} from "lucide-react";
import { User, ApiError } from "@/lib/types";
import { useUsers, useDeleteUser } from "@/lib/hooks";
import UserModal from "@/components/UserModal";
import LanguageFlag from "@/components/LanguageFlag";

// Hook personnalisé pour le debounce
function useDebounce<T>(value: T, delay: number): T {
  const [debouncedValue, setDebouncedValue] = useState<T>(value);

  useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedValue(value);
    }, delay);

    return () => {
      clearTimeout(handler);
    };
  }, [value, delay]);

  return debouncedValue;
}

export default function UsersPage() {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingUser, setEditingUser] = useState<User | null>(null);
  const [searchTerm, setSearchTerm] = useState("");
  const [sortField, setSortField] = useState<string>("");
  const [sortDirection, setSortDirection] = useState<"asc" | "desc">("asc");

  // Debounce du terme de recherche avec un délai de 300ms
  const debouncedSearchTerm = useDebounce(searchTerm, 300);

  const { data: users = [], isLoading, error } = useUsers();
  const deleteUserMutation = useDeleteUser();

  // Fonction de tri mémorisée avec useCallback
  const sortUsers = useCallback(
    (usersToSort: User[]) => {
      if (!sortField) return usersToSort;

      return [...usersToSort].sort((a, b) => {
        let aValue: string | number;
        let bValue: string | number;

        switch (sortField) {
          case "name":
            aValue = `${a.first_name || ""} ${a.last_name || ""}`.toLowerCase();
            bValue = `${b.first_name || ""} ${b.last_name || ""}`.toLowerCase();
            break;
          case "email":
            aValue = a.email.toLowerCase();
            bValue = b.email.toLowerCase();
            break;
          case "campus":
            aValue = (a.campus || "").toLowerCase();
            bValue = (b.campus || "").toLowerCase();
            break;
          case "formation":
            aValue = (a.formation_name || "").toLowerCase();
            bValue = (b.formation_name || "").toLowerCase();
            break;
          case "roles":
            aValue = (a.roles || []).join(", ").toLowerCase();
            bValue = (b.roles || []).join(", ").toLowerCase();
            break;
          case "verified":
            aValue = a.verification_code ? 0 : 1;
            bValue = b.verification_code ? 0 : 1;
            break;
          default:
            return 0;
        }

        if (aValue < bValue) return sortDirection === "asc" ? -1 : 1;
        if (aValue > bValue) return sortDirection === "asc" ? 1 : -1;
        return 0;
      });
    },
    [sortField, sortDirection]
  );

  // Filtrer et trier les utilisateurs
  const filteredAndSortedUsers = useMemo(() => {
    if (!debouncedSearchTerm.trim()) return sortUsers(users);

    const searchLower = debouncedSearchTerm.toLowerCase();
    const filtered = users.filter((user: User) => {
      return (
        user.first_name?.toLowerCase().includes(searchLower) ||
        user.last_name?.toLowerCase().includes(searchLower) ||
        user.email.toLowerCase().includes(searchLower) ||
        user.campus?.toLowerCase().includes(searchLower) ||
        user.formation_name?.toLowerCase().includes(searchLower) ||
        user.roles?.some((role: string) =>
          role.toLowerCase().includes(searchLower)
        )
      );
    });

    return sortUsers(filtered);
  }, [users, debouncedSearchTerm, sortUsers]);

  // Gérer le changement de tri
  const handleSort = (field: string) => {
    if (sortField === field) {
      setSortDirection(sortDirection === "asc" ? "desc" : "asc");
    } else {
      setSortField(field);
      setSortDirection("asc");
    }
  };

  // Rendu de l'icône de tri
  const renderSortIcon = (field: string) => {
    if (sortField !== field) {
      return <ChevronUp className="h-4 w-4 text-gray-400" />;
    }
    return sortDirection === "asc" ? (
      <ChevronUp className="h-4 w-4 text-blue-600" />
    ) : (
      <ChevronDown className="h-4 w-4 text-blue-600" />
    );
  };

  const handleCreateUser = () => {
    setEditingUser(null);
    setIsModalOpen(true);
  };

  const handleEditUser = (user: User) => {
    setEditingUser(user);
    setIsModalOpen(true);
  };

  const handleDeleteUser = async (email: string) => {
    if (!confirm("Êtes-vous sûr de vouloir supprimer cet utilisateur ?"))
      return;

    try {
      await deleteUserMutation.mutateAsync(email);
    } catch (err: unknown) {
      alert(
        (err as ApiError)?.response?.data?.error ||
          "Échec de la suppression de l'utilisateur"
      );
    }
  };

  if (isLoading) {
    return (
      <div className="p-6">
        <div className="animate-pulse space-y-4">
          <div className="h-8 bg-gray-200 rounded w-1/4"></div>
          <div className="space-y-3">
            {[1, 2, 3].map((i) => (
              <div key={i} className="h-16 bg-gray-200 rounded"></div>
            ))}
          </div>
        </div>
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
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold text-gray-900">
          Gestion des utilisateurs
        </h1>
        <div className="flex items-center space-x-4">
          <div className="flex items-center space-x-2 text-gray-600">
            <Users className="h-5 w-5" />
            <span>{filteredAndSortedUsers.length} utilisateurs</span>
          </div>
          <button
            onClick={handleCreateUser}
            className="inline-flex items-center px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
          >
            <Plus className="h-4 w-4 mr-2" />
            Créer un utilisateur
          </button>
        </div>
      </div>

      {/* Barre de recherche */}
      <div className="mb-6">
        <div className="relative max-w-md">
          <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <Search className="h-5 w-5 text-gray-400" />
          </div>
          <input
            type="text"
            placeholder="Rechercher par nom, email, campus..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md leading-5 bg-white placeholder-gray-500 focus:outline-none focus:placeholder-gray-400 focus:ring-1 focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
          />
          {/* Indicateur de chargement de la recherche */}
          {searchTerm !== debouncedSearchTerm && (
            <div className="absolute inset-y-0 right-0 pr-3 flex items-center">
              <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-blue-600"></div>
            </div>
          )}
        </div>
      </div>

      {/* En-têtes de colonnes avec tri */}
      <div className="bg-white shadow overflow-hidden sm:rounded-md">
        <div className="px-6 py-3 border-b border-gray-200 bg-gray-50">
          <div className="grid grid-cols-12 gap-4 text-xs font-medium text-gray-500 uppercase tracking-wider">
            <div className="col-span-3">
              <button
                onClick={() => handleSort("name")}
                className="flex items-center space-x-1 hover:text-gray-700 transition-colors"
              >
                <span>Nom</span>
                {renderSortIcon("name")}
              </button>
            </div>
            <div className="col-span-3">
              <button
                onClick={() => handleSort("email")}
                className="flex items-center space-x-1 hover:text-gray-700 transition-colors"
              >
                <span>Email</span>
                {renderSortIcon("email")}
              </button>
            </div>
            <div className="col-span-2">
              <button
                onClick={() => handleSort("campus")}
                className="flex items-center space-x-1 hover:text-gray-700 transition-colors"
              >
                <span>Campus</span>
                {renderSortIcon("campus")}
              </button>
            </div>
            <div className="col-span-2">
              <button
                onClick={() => handleSort("roles")}
                className="flex items-center space-x-1 hover:text-gray-700 transition-colors"
              >
                <span>Rôles</span>
                {renderSortIcon("roles")}
              </button>
            </div>
            <div className="col-span-2 text-right">Actions</div>
          </div>
        </div>

        <ul className="divide-y divide-gray-200">
          {filteredAndSortedUsers.map((user) => (
            <li key={user.email} className="px-6 py-4">
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-4">
                  <div className="flex-shrink-0">
                    <div className="h-10 w-10 rounded-full bg-gray-300 flex items-center justify-center">
                      <span className="text-sm font-medium text-gray-700">
                        {user.first_name?.[0]}
                        {user.last_name?.[0]}
                      </span>
                    </div>
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center space-x-2">
                      <p className="text-sm font-medium text-gray-900 truncate">
                        {user.first_name} {user.last_name}
                      </p>
                      {user.verification_code ? (
                        <div title="Non vérifié">
                          <X className="h-4 w-4 text-red-500" />
                        </div>
                      ) : (
                        <div title="Vérifié">
                          <Check className="h-4 w-4 text-green-500" />
                        </div>
                      )}
                    </div>
                    <div className="flex items-center space-x-4 text-sm text-gray-500">
                      <div className="flex items-center space-x-1">
                        <Mail className="h-4 w-4" />
                        <span>{user.email}</span>
                      </div>
                      {user.language && (
                        <div className="flex items-center space-x-1">
                          <LanguageFlag
                            languageCode={user.language}
                            className="text-base"
                          />
                          <span>{user.language?.toUpperCase()}</span>
                        </div>
                      )}
                      {user.campus && <span>Campus : {user.campus}</span>}
                      {user.formation_name && (
                        <span>Formation : {user.formation_name}</span>
                      )}
                    </div>
                  </div>
                </div>

                <div className="flex items-center space-x-2">
                  {user.roles && user.roles.length > 0 && (
                    <div
                      className={`grid gap-1 ${
                        user.roles.length > 3 ? "grid-cols-2" : "flex space-x-1"
                      }`}
                    >
                      {user.roles.map((role: string, index: number) => (
                        <span
                          key={index}
                          className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800 whitespace-nowrap"
                        >
                          {role}
                        </span>
                      ))}
                    </div>
                  )}

                  <div className="flex space-x-1">
                    <button
                      onClick={() => handleEditUser(user)}
                      className="p-1 text-gray-400 hover:text-blue-600"
                      title="Modifier l'utilisateur"
                    >
                      <Edit className="h-4 w-4" />
                    </button>
                    <button
                      onClick={() => handleDeleteUser(user.email)}
                      className="p-1 text-gray-400 hover:text-red-600"
                      title="Supprimer l'utilisateur"
                    >
                      <Trash2 className="h-4 w-4" />
                    </button>
                  </div>
                </div>
              </div>
            </li>
          ))}
        </ul>
      </div>

      {filteredAndSortedUsers.length === 0 && (
        <div className="text-center py-12">
          <Users className="mx-auto h-12 w-12 text-gray-400" />
          <h3 className="mt-2 text-sm font-medium text-gray-900">
            {searchTerm.trim()
              ? "Aucun utilisateur trouvé"
              : "Aucun utilisateur"}
          </h3>
          <p className="mt-1 text-sm text-gray-500">
            {searchTerm.trim()
              ? `Aucun utilisateur ne correspond à la recherche "${searchTerm}".`
              : "Aucun utilisateur trouvé dans le système."}
          </p>
        </div>
      )}

      <UserModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        user={editingUser}
        onSave={() => {}}
      />
    </div>
  );
}
