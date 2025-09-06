"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useMemo, memo } from "react";
import {
  Users,
  Calendar,
  Building,
  BarChart3,
  LogOut,
  Menu,
  X,
  Search,
  UtensilsCrossed,
  Gamepad2,
} from "lucide-react";
import Image from "next/image";
import { useAuthStore } from "@/lib/stores/authStore";
import { useAppStore } from "@/lib/stores/appStore";

function Sidebar() {
  const pathname = usePathname();
  const { user, logout } = useAuthStore();
  const { sidebarOpen, setSidebarOpen, toggleSidebar, setCommandPaletteOpen } =
    useAppStore();

  // Memoize menu items to prevent recreation on every render
  const menuItems = useMemo(
    () => [
      { href: "/dashboard", label: "Tableau de bord", icon: BarChart3 },
      { href: "/users", label: "Utilisateurs", icon: Users },
      { href: "/events", label: "Événements", icon: Calendar },
      { href: "/clubs", label: "Clubs", icon: Building },
      { href: "/menu", label: "Menu du RU", icon: UtensilsCrossed },
      { href: "/games", label: "Jeux", icon: Gamepad2 },
    ],
    []
  );

  const handleLogout = () => {
    logout();
    window.location.href = "/login";
  };

  return (
    <>
      {/* Mobile menu button */}
      <button
        onClick={toggleSidebar}
        className="lg:hidden fixed top-4 left-4 z-50 p-2 bg-gray-900 text-white rounded-md"
      >
        {sidebarOpen ? <X size={20} /> : <Menu size={20} />}
      </button>

      {/* Mobile overlay */}
      {sidebarOpen && (
        <div
          className="lg:hidden fixed inset-0 bg-black/50 z-40"
          onClick={() => setSidebarOpen(false)}
        />
      )}

      {/* Sidebar */}
      <div
        className={`
        fixed lg:static inset-y-0 left-0 z-40 w-64 bg-gray-900 text-white flex flex-col transition-transform duration-300 ease-in-out
        ${sidebarOpen ? "translate-x-0" : "-translate-x-full lg:translate-x-0"}
      `}
      >
        <div className="p-4 lg:p-6 border-b border-gray-700 flex flex-col items-center">
          <Image
            src="/transat-admin.png"
            alt="Transat Admin"
            width={80}
            height={80}
            className="lg:w-[100px] lg:h-[100px]"
          />
          <h1 className="text-lg lg:text-xl font-bold">Transat Admin</h1>
        </div>

        <nav className="flex-1 p-4">
          <ul className="space-y-2">
            {menuItems.map((item) => {
              const Icon = item.icon;
              const isActive = pathname === item.href;

              return (
                <li key={item.href}>
                  <Link
                    href={item.href}
                    onClick={() => setSidebarOpen(false)}
                    className={`flex items-center space-x-3 px-4 py-2 rounded-lg transition-colors ${
                      isActive
                        ? "bg-blue-600 text-white"
                        : "text-gray-300 hover:bg-gray-800 hover:text-white"
                    }`}
                  >
                    <Icon size={20} />
                    <span>{item.label}</span>
                  </Link>
                </li>
              );
            })}

            {/* Command Palette Button */}
            <li className="pt-2 border-t border-gray-700 mt-4">
              <button
                onClick={() => {
                  setCommandPaletteOpen(true);
                  setSidebarOpen(false);
                }}
                className="flex items-center justify-between w-full px-4 py-2 text-gray-300 hover:bg-gray-800 hover:text-white rounded-lg transition-colors"
              >
                <div className="flex items-center space-x-3">
                  <Search size={20} />
                  <span>Recherche</span>
                </div>
                <kbd className="text-xs bg-gray-700 px-2 py-1 rounded">⌘K</kbd>
              </button>
            </li>
          </ul>
        </nav>

        <div className="p-4 border-t border-gray-700">
          <p className="text-xs lg:text-sm text-gray-400 mb-4 truncate">
            Bienvenue {user?.email || "Utilisateur"}
          </p>
          <button
            onClick={handleLogout}
            className="flex items-center space-x-3 px-4 py-2 rounded-lg text-gray-300 hover:bg-gray-800 hover:text-white transition-colors w-full"
          >
            <LogOut size={20} />
            <span>Déconnexion</span>
          </button>
        </div>
      </div>
    </>
  );
}

export default memo(Sidebar);
