"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { Users, Calendar, Building, BarChart3, LogOut } from "lucide-react";
import Image from "next/image";

export default function Sidebar() {
  const pathname = usePathname();

  const menuItems = [
    { href: "/dashboard", label: "Tableau de bord", icon: BarChart3 },
    { href: "/users", label: "Utilisateurs", icon: Users },
    { href: "/events", label: "Événements", icon: Calendar },
    { href: "/clubs", label: "Clubs", icon: Building },
  ];

  const handleLogout = () => {
    localStorage.removeItem("adminToken");
    window.location.href = "/login";
  };

  const getAdminEmail = () => {
    const token = localStorage.getItem("adminToken");
    if (!token) return null;
    const payload = JSON.parse(atob(token.split(".")[1]));
    return payload.email;
  };

  return (
    <div className="h-screen w-64 bg-gray-900 text-white flex flex-col">
      <div className="p-6 border-b border-gray-700 flex flex-col items-center">
        <Image
          src="/transat-admin.png"
          alt="Transat Admin"
          width={100}
          height={100}
        />
        <h1 className="text-xl font-bold">Transat Admin</h1>
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
        </ul>
      </nav>

      <div className="p-4 border-t border-gray-700">
        <p className="text-sm text-gray-400 mb-4">
          Bienvenue {getAdminEmail()}
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
  );
}
