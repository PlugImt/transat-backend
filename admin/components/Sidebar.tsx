"use client";

import {
  BarChart3,
  Beer,
  Bell,
  Building,
  Calendar,
  CalendarCheck,
  Flag,
  Gamepad2,
  LogOut,
  Menu,
  Monitor,
  Moon,
  Search,
  Sun,
  Users,
  UtensilsCrossed,
  X,
} from "lucide-react";
import Image from "next/image";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { memo, useMemo } from "react";
import { canOpen } from "@/lib/access";
import { useAppStore } from "@/lib/stores/appStore";
import { useAuthStore } from "@/lib/stores/authStore";
import { type ThemeChoice, useTheme } from "@/lib/theme";

const MENU = [
  { href: "/dashboard", label: "Tableau de bord", icon: BarChart3 },
  { href: "/users", label: "Utilisateurs", icon: Users },
  { href: "/events", label: "Événements", icon: Calendar },
  { href: "/clubs", label: "Clubs", icon: Building },
  { href: "/reservations", label: "Réservations", icon: CalendarCheck },
  { href: "/menu", label: "RU", icon: UtensilsCrossed },
  { href: "/traq", label: "Traq", icon: Beer },
  { href: "/games", label: "Jeux", icon: Gamepad2 },
  { href: "/notifications", label: "Notifications", icon: Bell },
  { href: "/photo-reports", label: "Signalements photos", icon: Flag },
];

const THEMES: { id: ThemeChoice; label: string; icon: typeof Sun }[] = [
  { id: "light", label: "Clair", icon: Sun },
  { id: "system", label: "Système", icon: Monitor },
  { id: "dark", label: "Sombre", icon: Moon },
];

function Sidebar() {
  const pathname = usePathname();
  const { user, logout, access } = useAuthStore();
  const { sidebarOpen, setSidebarOpen, toggleSidebar, setCommandPaletteOpen } = useAppStore();
  const { choice, setChoice } = useTheme();

  // Only the sections this account may open.
  const menuItems = useMemo(() => MENU.filter((item) => canOpen(access, item.href)), [access]);
  const isAdmin = access?.admin ?? false;

  const handleLogout = () => {
    logout();
    window.location.href = "/login";
  };

  const initials = (user?.firstName?.[0] ?? user?.email?.[0] ?? "?").toUpperCase();

  return (
    <>
      <button
        type="button"
        onClick={toggleSidebar}
        aria-label={sidebarOpen ? "Fermer le menu" : "Ouvrir le menu"}
        className="lg:hidden fixed top-4 left-4 z-50 rounded-xl bg-[var(--sidebar)] p-2.5 text-[#ffe6cc] shadow-lg"
      >
        {sidebarOpen ? <X size={20} /> : <Menu size={20} />}
      </button>

      {sidebarOpen && (
        <button
          className="lg:hidden fixed inset-0 z-40 bg-black/60 backdrop-blur-sm animate-fade-in"
          type="button"
          aria-label="Fermer le menu"
          onClick={() => setSidebarOpen(false)}
        />
      )}

      <aside
        className={`fixed lg:static inset-y-0 left-0 z-40 flex w-64 flex-col border-r border-[var(--sidebar-border)] bg-[var(--sidebar)] text-[#ffe6cc] transition-transform duration-300 ease-out ${
          sidebarOpen ? "translate-x-0" : "-translate-x-full lg:translate-x-0"
        }`}
      >
        <div className="flex items-center gap-3 px-5 pt-6 pb-5">
          <Image
            src="/transat-admin.png"
            alt=""
            width={44}
            height={44}
            className="rounded-xl"
            priority
            loading="eager"
          />
          <div className="leading-tight">
            <p className="text-base font-bold tracking-tight">Transat</p>
            <p className="text-xs text-[var(--sidebar-text)]">
              {isAdmin ? "Administration" : "Espace Club Photo"}
            </p>
          </div>
        </div>

        <nav className="flex-1 overflow-y-auto px-3">
          <ul className="space-y-1">
            {menuItems.map((item) => {
              const Icon = item.icon;
              const isActive = pathname === item.href;
              return (
                <li key={item.href}>
                  <Link
                    href={item.href}
                    onClick={() => setSidebarOpen(false)}
                    aria-current={isActive ? "page" : undefined}
                    className={`group relative flex items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-medium transition-all duration-200 ${
                      isActive
                        ? "bg-[#ec7f32]/15 text-[#ffb27a]"
                        : "text-[var(--sidebar-text)] hover:bg-[var(--sidebar-hover)] hover:text-[#ffe6cc]"
                    }`}
                  >
                    <span
                      className={`absolute left-0 top-1/2 h-5 w-1 -translate-y-1/2 rounded-r-full bg-[#ec7f32] transition-opacity ${
                        isActive ? "opacity-100" : "opacity-0"
                      }`}
                    />
                    <Icon
                      size={18}
                      className={`transition-transform duration-200 ${isActive ? "" : "group-hover:scale-110"}`}
                    />
                    <span>{item.label}</span>
                  </Link>
                </li>
              );
            })}
          </ul>

          {isAdmin && (
            <button
              type="button"
              onClick={() => {
                setCommandPaletteOpen(true);
                setSidebarOpen(false);
              }}
              className="mt-4 flex w-full items-center justify-between rounded-xl border border-[var(--sidebar-border)] px-3 py-2 text-sm text-[var(--sidebar-text)] transition-colors hover:bg-[var(--sidebar-hover)] hover:text-[#ffe6cc]"
            >
              <span className="flex items-center gap-3">
                <Search size={16} />
                Recherche
              </span>
              <kbd className="rounded-md bg-[var(--sidebar-hover)] px-1.5 py-0.5 text-[11px] text-[var(--sidebar-text)]">
                ⌘K
              </kbd>
            </button>
          )}
        </nav>

        <div className="space-y-3 border-t border-[var(--sidebar-border)] p-4">
          <fieldset className="grid grid-cols-3 gap-1 rounded-xl bg-[var(--sidebar-hover)] p-1">
            <legend className="sr-only">Thème</legend>
            {THEMES.map(({ id, label, icon: Icon }) => (
              <button
                key={id}
                type="button"
                onClick={() => setChoice(id)}
                aria-pressed={choice === id}
                title={label}
                className={`flex items-center justify-center rounded-lg py-1.5 transition-all ${
                  choice === id
                    ? "bg-[#ec7f32] text-white shadow"
                    : "text-[var(--sidebar-text)] hover:text-[#ffe6cc]"
                }`}
              >
                <Icon size={15} />
                <span className="sr-only">{label}</span>
              </button>
            ))}
          </fieldset>

          <div className="flex items-center gap-3">
            <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-gradient-to-br from-[#ec7f32] to-[#b45417] text-sm font-bold text-white">
              {initials}
            </div>
            <div className="min-w-0 flex-1">
              <p className="truncate text-sm font-medium">
                {user?.firstName ? `${user.firstName} ${user.lastName ?? ""}` : "Connecté"}
              </p>
              <p className="truncate text-xs text-[var(--sidebar-text)]">{user?.email}</p>
            </div>
            <button
              type="button"
              onClick={handleLogout}
              title="Déconnexion"
              aria-label="Déconnexion"
              className="rounded-lg p-2 text-[var(--sidebar-text)] transition-colors hover:bg-[var(--sidebar-hover)] hover:text-[#ffe6cc]"
            >
              <LogOut size={17} />
            </button>
          </div>
        </div>
      </aside>
    </>
  );
}

export default memo(Sidebar);
