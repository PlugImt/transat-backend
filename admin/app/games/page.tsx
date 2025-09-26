"use client";

import { Gamepad2 } from "lucide-react";
import { useState } from "react";
import BassineScores from "@/components/BassineScores";
import ErrorBoundary from "@/components/ErrorBoundary";

type GameTab = "bassine";

function GamesPageContent() {
  const [activeTab, setActiveTab] = useState<GameTab>("bassine");

  const tabs = [
    {
      id: "bassine" as GameTab,
      label: "Bassine",
      icon: () => <span className="text-lg">🥣</span>,
      description: "Gestion des scores de bassine",
    },
  ];

  return (
    <div className="p-4 sm:p-6 pt-16 lg:pt-6">
      <div className="flex flex-col sm:flex-row sm:justify-between sm:items-center mb-6 space-y-4 sm:space-y-0">
        <div className="flex items-center space-x-3">
          <Gamepad2 className="h-6 w-6 text-purple-600" />
          <h1 className="text-xl sm:text-2xl font-bold text-gray-900">Jeux</h1>
        </div>
      </div>

      {/* Tabs Navigation */}
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
                  className={`flex items-center space-x-2 py-4 px-1 border-b-2 font-medium text-sm transition-colors ${
                    isActive
                      ? "border-purple-500 text-purple-600"
                      : "border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300"
                  }`}
                >
                  <Icon />
                  <span>{tab.label}</span>
                </button>
              );
            })}
          </nav>
        </div>

        {/* Tab Description */}
        <div className="mt-4">
          <p className="text-sm text-gray-600">
            {tabs.find((tab) => tab.id === activeTab)?.description}
          </p>
        </div>
      </div>

      {/* Tab Content */}
      <div className="min-h-[400px]">{activeTab === "bassine" && <BassineScores />}</div>
    </div>
  );
}

export default function GamesPage() {
  return (
    <ErrorBoundary>
      <GamesPageContent />
    </ErrorBoundary>
  );
}
