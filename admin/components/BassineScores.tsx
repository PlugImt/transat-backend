"use client";

import type { ColumnDef } from "@tanstack/react-table";
import { Award, Calendar, Edit, History, TrendingUp, User } from "lucide-react";
import { useCallback, useMemo, useState } from "react";
import BassineHistoryModal from "@/components/BassineHistoryModal";
import BassineScoreModal from "@/components/BassineScoreModal";
import DataTable from "@/components/DataTable";
import { PageLoading } from "@/components/LoadingSpinner";
import { useBassineScores } from "@/lib/hooks";
import type { ApiError, BassineScore } from "@/lib/types";

export default function BassineScores() {
  const [scoreModalOpen, setScoreModalOpen] = useState(false);
  const [historyModalOpen, setHistoryModalOpen] = useState(false);
  const [selectedUser, setSelectedUser] = useState<BassineScore | null>(null);

  const { data: scores = [], isLoading, error } = useBassineScores();

  const handleEditScore = useCallback((score: BassineScore) => {
    setSelectedUser(score);
    setScoreModalOpen(true);
  }, []);

  const handleViewHistory = useCallback((score: BassineScore) => {
    setSelectedUser(score);
    setHistoryModalOpen(true);
  }, []);

  const handleCloseScoreModal = useCallback(() => {
    setScoreModalOpen(false);
    setSelectedUser(null);
  }, []);

  const handleCloseHistoryModal = useCallback(() => {
    setHistoryModalOpen(false);
    setSelectedUser(null);
  }, []);

  // Stats calculées
  const stats = useMemo(() => {
    if (!scores.length) return null;

    const totalPlayers = scores.length;
    const totalGames = scores.reduce((acc, score) => acc + score.total_games_played, 0);
    const averageScore = scores.reduce((acc, score) => acc + score.current_score, 0) / totalPlayers;
    const topScore = Math.max(...scores.map((score) => score.current_score));

    return { totalPlayers, totalGames, averageScore, topScore };
  }, [scores]);

  const columns = useMemo<ColumnDef<BassineScore>[]>(
    () => [
      {
        accessorKey: "user_info",
        header: "Joueur",
        accessorFn: (row) => {
          return `${row.user_first_name || ""} ${row.user_last_name || ""} ${row.user_email}`.toLowerCase();
        },
        cell: ({ row }) => {
          const score = row.original;
          return (
            <div className="flex items-center space-x-3">
              <div className="flex-shrink-0">
                <div className="h-10 w-10 rounded-full bg-purple-100 flex items-center justify-center">
                  <User className="h-5 w-5 text-purple-600" />
                </div>
              </div>
              <div>
                <div className="text-sm font-medium text-gray-900">
                  {score.user_first_name} {score.user_last_name}
                </div>
                <div className="text-sm text-gray-500">{score.user_email}</div>
              </div>
            </div>
          );
        },
      },
      {
        accessorKey: "current_score",
        header: "Score actuel",
        cell: ({ row }) => (
          <div className="flex items-center space-x-2">
            <span className="text-lg font-bold text-gray-900">{row.original.current_score}</span>
          </div>
        ),
      },
      {
        accessorKey: "last_updated",
        header: "Dernière activité",
        cell: ({ row }) => (
          <div className="flex items-center space-x-1">
            <Calendar className="h-4 w-4 text-green-500" />
            <span className="text-sm text-gray-600">
              {new Date(row.original.last_updated).toLocaleDateString("fr-FR")}
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
              onClick={() => handleEditScore(row.original)}
              className="p-1 text-blue-600 hover:bg-blue-50 rounded"
              title="Modifier le score"
            >
              <Edit size={16} />
            </button>
            <button
              type="button"
              onClick={() => handleViewHistory(row.original)}
              className="p-1 text-purple-600 hover:bg-purple-50 rounded"
              title="Voir l'historique"
            >
              <History size={16} />
            </button>
          </div>
        ),
      },
    ],
    [handleEditScore, handleViewHistory],
  );

  if (isLoading) {
    return <PageLoading text="Chargement des scores..." />;
  }

  if (error) {
    return (
      <div className="bg-red-50 border border-red-200 rounded-md p-4">
        <div className="text-sm text-red-700">
          {(error as ApiError)?.message || "Échec de la récupération des scores"}
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Stats Cards */}
      {stats && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          {[
            {
              title: "Alcooliques",
              value: stats.totalPlayers,
              icon: User,
              color: "text-blue-600",
              bgColor: "bg-blue-50",
            },
            {
              title: "Bassines bues en tout",
              value: stats.totalGames,
              icon: TrendingUp,
              color: "text-green-600",
              bgColor: "bg-green-50",
            },
            {
              title: "Moyenne de bassine par joueur",
              value: stats.averageScore.toFixed(1),
              icon: () => <span className="text-2xl">🥣</span>,
              color: "text-yellow-600",
              bgColor: "bg-yellow-50",
            },
            {
              title: "Plus grand nombre de bassine",
              value: stats.topScore,
              icon: Award,
              color: "text-purple-600",
              bgColor: "bg-purple-50",
            },
          ].map((stat) => {
            const Icon = stat.icon;
            return (
              <div key={stat.title} className="bg-white overflow-hidden shadow rounded-lg">
                <div className="p-5">
                  <div className="flex items-center">
                    <div className="flex-shrink-0">
                      <div className={`p-3 rounded-md ${stat.bgColor}`}>
                        <Icon className={`h-6 w-6 ${stat.color}`} />
                      </div>
                    </div>
                    <div className="ml-5 w-0 flex-1">
                      <dl>
                        <dt className="text-sm font-medium text-gray-500 truncate">{stat.title}</dt>
                        <dd className="text-lg font-medium text-gray-900">
                          {typeof stat.value === "number" && stat.title !== "Score moyen"
                            ? stat.value.toLocaleString()
                            : stat.value}
                        </dd>
                      </dl>
                    </div>
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      )}

      {/* Table */}
      <DataTable
        data={scores}
        columns={columns}
        searchPlaceholder="Rechercher par nom ou email..."
        className="mb-6"
      />

      {/* Modals */}
      {selectedUser && (
        <>
          <BassineScoreModal
            isOpen={scoreModalOpen}
            onClose={handleCloseScoreModal}
            user={selectedUser}
          />
          <BassineHistoryModal
            isOpen={historyModalOpen}
            onClose={handleCloseHistoryModal}
            userEmail={selectedUser.user_email}
            userName={`${selectedUser.user_first_name} ${selectedUser.user_last_name}`}
          />
        </>
      )}
    </div>
  );
}
