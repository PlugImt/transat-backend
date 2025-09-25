"use client";

import { Calendar, History, TrendingDown, TrendingUp, User2 } from "lucide-react";
import { useMemo } from "react";
import { PageLoading } from "@/components/LoadingSpinner";
import Modal from "@/components/Modal";
import { useBassineHistory } from "@/lib/hooks";
import type { ApiError } from "@/lib/types";

interface BassineHistoryModalProps {
  isOpen: boolean;
  onClose: () => void;
  userEmail: string;
  userName: string;
}

export default function BassineHistoryModal({
  isOpen,
  onClose,
  userEmail,
  userName,
}: BassineHistoryModalProps) {
  const { data: history = [], isLoading, error } = useBassineHistory(userEmail);

  // Statistiques de l'historique
  const stats = useMemo(() => {
    if (!history.length) return null;

    const totalChanges = history.length;
    const positiveChanges = history.filter((h) => h.score_change > 0).length;
    const negativeChanges = history.filter((h) => h.score_change < 0).length;
    const totalScoreChange = history.reduce((acc, h) => acc + h.score_change, 0);

    return {
      totalChanges,
      positiveChanges,
      negativeChanges,
      totalScoreChange,
    };
  }, [history]);

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString("fr-FR", {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  const getChangeColor = (change: number) => {
    if (change > 0) return "text-green-600";
    if (change < 0) return "text-red-600";
    return "text-gray-600";
  };

  const getChangeIcon = (change: number) => {
    if (change > 0) return <TrendingUp className="h-4 w-4 text-green-500" />;
    if (change < 0) return <TrendingDown className="h-4 w-4 text-red-500" />;
    return null;
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} title={`Historique de ${userName}`} maxWidth="xl">
      <div className="max-h-[80vh] overflow-y-auto space-y-6">
        {isLoading ? (
          <PageLoading text="Chargement de l'historique..." />
        ) : error ? (
          <div className="bg-red-50 border border-red-200 rounded-md p-4">
            <div className="text-sm text-red-700">
              {(error as ApiError)?.message || "Erreur lors du chargement de l'historique"}
            </div>
          </div>
        ) : (
          <>
            {/* Stats */}
            {stats && (
              <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                <div className="bg-blue-50 rounded-lg p-4">
                  <p className="text-sm font-medium text-blue-600">Total modifications</p>
                  <p className="text-xl font-bold text-blue-900">{stats.totalChanges}</p>
                </div>
                <div className="bg-green-50 rounded-lg p-4">
                  <p className="text-sm font-medium text-green-600">Gains</p>
                  <p className="text-xl font-bold text-green-900">{stats.positiveChanges}</p>
                </div>
                <div className="bg-red-50 rounded-lg p-4">
                  <p className="text-sm font-medium text-red-600">Pertes</p>
                  <p className="text-xl font-bold text-red-900">{stats.negativeChanges}</p>
                </div>
                <div className="bg-purple-50 rounded-lg p-4">
                  <p className="text-sm font-medium text-purple-600">Total changement</p>
                  <p className={`text-xl font-bold ${getChangeColor(stats.totalScoreChange)}`}>
                    {stats.totalScoreChange > 0 ? "+" : ""}
                    {stats.totalScoreChange}
                  </p>
                </div>
              </div>
            )}

            {/* Timeline */}
            {history.length > 0 ? (
              <div className="space-y-4">
                <h3 className="text-lg font-medium text-gray-900">Timeline des modifications</h3>
                <div className="space-y-3">
                  {history.map((entry) => (
                    <div
                      key={entry.id}
                      className="flex items-start space-x-3 p-4 bg-gray-50 rounded-lg"
                    >
                      <div className="flex-shrink-0 mt-0.5">
                        {getChangeIcon(entry.score_change)}
                      </div>
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center justify-between">
                          <div className="flex items-center space-x-2">
                            <span className={`font-semibold ${getChangeColor(entry.score_change)}`}>
                              {entry.score_change > 0 ? "+" : ""}
                              {entry.score_change}
                            </span>
                            <span className="text-gray-600">→</span>
                            <span className="font-medium text-gray-900">{entry.new_total}</span>
                          </div>
                          <div className="flex items-center space-x-1 text-sm text-gray-500">
                            <Calendar className="h-4 w-4" />
                            <span>{formatDate(entry.game_date)}</span>
                          </div>
                        </div>
                        {entry.notes && <p className="mt-1 text-sm text-gray-600">{entry.notes}</p>}
                        {entry.admin_email && (
                          <div className="mt-2 flex items-center space-x-1 text-xs text-gray-500">
                            <User2 className="h-3 w-3" />
                            <span>Modifié par: {entry.admin_email}</span>
                          </div>
                        )}
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            ) : (
              <div className="text-center py-12">
                <History className="mx-auto h-12 w-12 text-gray-400" />
                <h3 className="mt-2 text-sm font-medium text-gray-900">Aucun historique</h3>
                <p className="mt-1 text-sm text-gray-500">
                  Aucune modification de score n&apos;a encore été enregistrée pour ce joueur.
                </p>
              </div>
            )}
          </>
        )}
      </div>
    </Modal>
  );
}
