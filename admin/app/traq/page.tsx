"use client";

import { Beer, Tags } from "lucide-react";
import { useEffect, useState } from "react";
import ErrorBoundary from "@/components/ErrorBoundary";
import TraqArticleModal from "@/components/TraqArticleModal";
import TraqArticlesManager from "@/components/TraqArticlesManager";
import TraqTypesManager from "@/components/TraqTypesManager";
import { useAppStore } from "@/lib/stores/appStore";

type TraqTab = "articles" | "types";

function TraqPageContent() {
  const [activeTab, setActiveTab] = useState<TraqTab>("articles");
  const { traqArticleModalOpen, editingTraqArticle, closeTraqArticleModal } = useAppStore();

  // Command palette can open the modal while Types is active — switch to Articles.
  useEffect(() => {
    if (traqArticleModalOpen) {
      setActiveTab("articles");
    }
  }, [traqArticleModalOpen]);

  const tabs = [
    {
      id: "articles" as TraqTab,
      label: "Articles",
      icon: Beer,
      description: "Gestion des boissons et articles Traq",
    },
    {
      id: "types" as TraqTab,
      label: "Types",
      icon: Tags,
      description: "Catégories d'articles (bière, soft, etc.)",
    },
  ];

  return (
    <div className="p-4 sm:p-6 pt-16 lg:pt-6">
      <div className="flex flex-col sm:flex-row sm:justify-between sm:items-center mb-6 space-y-4 sm:space-y-0">
        <div className="flex items-center space-x-3">
          <Beer className="h-6 w-6 text-amber-600" />
          <h1 className="text-xl sm:text-2xl font-bold text-gray-900">Traq</h1>
        </div>
      </div>

      <div className="mb-6">
        <div className="border-b border-gray-200">
          <nav className="-mb-px flex space-x-8">
            {tabs.map((tab) => {
              const Icon = tab.icon;
              const isActive = activeTab === tab.id;

              return (
                <button
                  key={tab.id}
                  type="button"
                  onClick={() => setActiveTab(tab.id)}
                  className={`flex items-center space-x-2 py-4 px-1 border-b-2 font-medium text-sm transition-colors cursor-pointer ${
                    isActive
                      ? "border-amber-500 text-amber-600"
                      : "border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300"
                  }`}
                >
                  <Icon className="h-4 w-4" />
                  <span>{tab.label}</span>
                </button>
              );
            })}
          </nav>
        </div>

        <div className="mt-4">
          <p className="text-sm text-gray-600">
            {tabs.find((tab) => tab.id === activeTab)?.description}
          </p>
        </div>
      </div>

      <div className="min-h-[400px]">
        {activeTab === "articles" && <TraqArticlesManager />}
        {activeTab === "types" && <TraqTypesManager />}
      </div>

      <TraqArticleModal
        isOpen={traqArticleModalOpen}
        onClose={closeTraqArticleModal}
        article={editingTraqArticle}
      />
    </div>
  );
}

export default function TraqPage() {
  return (
    <ErrorBoundary>
      <TraqPageContent />
    </ErrorBoundary>
  );
}
