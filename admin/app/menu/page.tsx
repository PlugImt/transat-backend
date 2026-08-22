"use client";

import { MessageSquare, UtensilsCrossed } from "lucide-react";
import { useState } from "react";
import { MessageSquare, UtensilsCrossed } from "lucide-react";
import { useState } from "react";
import ErrorBoundary from "@/components/ErrorBoundary";
import MenuManager from "@/components/MenuManager";
import ReviewsManager from "@/components/ReviewsManager";

type MenuTab = "menu" | "reviews";
import MenuManager from "@/components/MenuManager";
import ReviewsManager from "@/components/ReviewsManager";

type MenuTab = "menu" | "reviews";

function MenuPageContent() {
  const [activeTab, setActiveTab] = useState<MenuTab>("menu");
  const [activeTab, setActiveTab] = useState<MenuTab>("menu");

  const tabs = [
    {
      id: "menu" as MenuTab,
      label: "Menu",
      icon: UtensilsCrossed,
      description: "Gestion des plats du restaurant universitaire",
    },
    {
      id: "reviews" as MenuTab,
      label: "Avis",
      icon: MessageSquare,
      description: "Gestion des avis laissés par les utilisateurs",
    },
  ];
  const tabs = [
    {
      id: "menu" as MenuTab,
      label: "Menu",
      icon: UtensilsCrossed,
      description: "Gestion des plats du restaurant universitaire",
    },
    {
      id: "reviews" as MenuTab,
      label: "Avis",
      icon: MessageSquare,
      description: "Gestion des avis laissés par les utilisateurs",
    },
  ];

  return (
    <div className="p-4 sm:p-6 pt-16 lg:pt-6">
      <div className="flex flex-col sm:flex-row sm:justify-between sm:items-center mb-6 space-y-4 sm:space-y-0">
        <div className="flex items-center space-x-3">
          <UtensilsCrossed className="h-6 w-6 text-orange-600" />
          <h1 className="text-xl sm:text-2xl font-bold text-gray-900">Restaurant</h1>
  if (isLoading) {
    return (
      <div className="p-4 sm:p-6 pt-16 lg:pt-6">
        <div className="flex flex-col sm:flex-row sm:justify-between sm:items-center mb-6 space-y-4 sm:space-y-0">
          <h1 className="text-xl sm:text-2xl font-bold text-gray-900">Menu du RU</h1>
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
                  className={`flex items-center space-x-2 py-4 px-1 border-b-2 font-medium text-sm transition-colors cursor-pointer ${
                    isActive
                      ? "border-orange-500 text-orange-600"
                      : "border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300"
                  }`}
                >
                  <Icon className="h-4 w-4" />
                  <span>{tab.label}</span>
                </button>
              );
            })}
          </nav>
  if (error) {
    return (
      <div className="p-6">
        <div className="bg-red-50 border border-red-200 rounded-md p-4">
          <div className="text-sm text-red-700">
            {(error as ApiError)?.message || "Échec de la récupération du menu"}
          </div>
        </div>

        {/* Tab Description */}
        <div className="mt-4">
          <p className="text-sm text-gray-600">
            {tabs.find((tab) => tab.id === activeTab)?.description}
          </p>
        </div>
      </div>

      {/* Tab Content */}
      <div className="min-h-[400px]">
        {activeTab === "menu" && <MenuManager />}
        {activeTab === "reviews" && <ReviewsManager />}
      </div>
    );
  }

  return (
    <div className="p-4 sm:p-6 pt-16 lg:pt-6">
      <div className="flex flex-col sm:flex-row sm:justify-between sm:items-center mb-6 space-y-4 sm:space-y-0">
        <h1 className="text-xl sm:text-2xl font-bold text-gray-900">Menu du RU</h1>
        <div className="flex items-center space-x-2 text-gray-600">
          <UtensilsCrossed className="h-5 w-5" />
          <span className="text-sm sm:text-base">
            {menuItems.length} plat{menuItems.length > 1 ? "s" : ""}
          </span>
        </div>
      </div>
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
