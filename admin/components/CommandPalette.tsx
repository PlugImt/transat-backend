"use client";

import { Command } from "cmdk";
import { BarChart3, Building, Calendar, LogOut, Plus, Search, Users } from "lucide-react";
import { useRouter } from "next/navigation";
import { memo, useEffect, useMemo, useState } from "react";
import { useClubs, useEvents, useKeyboardShortcuts, useUsers } from "@/lib/hooks";
import { useAppStore } from "@/lib/stores/appStore";
import { useAuthStore } from "@/lib/stores/authStore";
import type { Club, Event, User } from "@/lib/types";

interface CommandItem {
  id: string;
  label: string;
  description?: string;
  icon: React.ComponentType<{ className?: string; size?: number }>;
  action: () => void;
  keywords?: string[];
}

function CommandPalette() {
  const router = useRouter();
  const {
    commandPaletteOpen,
    setCommandPaletteOpen,
    openUserModal,
    openClubModal,
    openEventModal,
  } = useAppStore();
  const { logout } = useAuthStore();
  const [search, setSearch] = useState("");

  // Fetch data for search
  const { data: users = [] } = useUsers();
  const { data: clubs = [] } = useClubs();
  const { data: events = [] } = useEvents();

  // Generate dynamic commands for users, clubs, and events
  const dynamicCommands = useMemo(() => {
    const userCommands: CommandItem[] = users.map((user: User) => ({
      id: `user-${user.email}`,
      label:
        `${user.first_name || ""} ${user.last_name || ""}`.trim() || user.email || "Utilisateur",
      description: `Utilisateur • ${user.email}`,
      icon: Users,
      action: () => {
        openUserModal(user);
        router.push("/users");
      },
      keywords: [
        user.first_name || "",
        user.last_name || "",
        user.email || "",
        user.campus || "",
        user.formation_name || "",
        ...(user.roles || []),
      ].filter(Boolean),
    }));

    const clubCommands: CommandItem[] = clubs.map((club: Club, index: number) => ({
      id: `club-${club.id_clubs || `unknown-${index}`}`,
      label: club.name || "Club sans nom",
      description: `Club • ${club.description || "Pas de description"}`,
      icon: Building,
      action: () => {
        openClubModal(club);
        router.push("/clubs");
      },
      keywords: [
        club.name || "",
        club.description || "",
        club.location || "",
        "club",
        "association",
      ].filter(Boolean),
    }));

    const eventCommands: CommandItem[] = events.map((event: Event, index: number) => ({
      id: `event-${event.id_events || `unknown-${index}`}`,
      label: event.name || "Événement sans titre",
      description: `Événement • ${
        event.start_date
          ? new Date(event.start_date).toLocaleDateString("fr-FR")
          : "Date non définie"
      }`,
      icon: Calendar,
      action: () => {
        openEventModal(event);
        router.push("/events");
      },
      keywords: [
        event.name || "",
        event.description || "",
        event.location || "",
        event.creator || "",
        "event",
        "evenement",
      ].filter(Boolean),
    }));

    return [...userCommands, ...clubCommands, ...eventCommands];
  }, [users, clubs, events, router, openUserModal, openClubModal, openEventModal]);

  // Static commands
  const staticCommands: CommandItem[] = useMemo(
    () => [
      // Navigation
      {
        id: "nav-dashboard",
        label: "Tableau de bord",
        description: "Voir les statistiques et métriques",
        icon: BarChart3,
        action: () => router.push("/dashboard"),
        keywords: ["dashboard", "stats", "statistiques", "home"],
      },
      {
        id: "nav-users",
        label: "Utilisateurs",
        description: "Gérer les utilisateurs",
        icon: Users,
        action: () => router.push("/users"),
        keywords: ["users", "utilisateurs", "gestion"],
      },
      {
        id: "nav-events",
        label: "Événements",
        description: "Gérer les événements",
        icon: Calendar,
        action: () => router.push("/events"),
        keywords: ["events", "evenements", "calendar", "calendrier"],
      },
      {
        id: "nav-clubs",
        label: "Clubs",
        description: "Gérer les clubs",
        icon: Building,
        action: () => router.push("/clubs"),
        keywords: ["clubs", "organisations"],
      },

      // Actions
      {
        id: "action-new-user",
        label: "Créer un utilisateur",
        description: "Ajouter un nouvel utilisateur",
        icon: Plus,
        action: () => {
          openUserModal();
          router.push("/users");
        },
        keywords: ["create", "creer", "nouveau", "new", "user", "utilisateur"],
      },
      {
        id: "action-new-club",
        label: "Créer un club",
        description: "Ajouter un nouveau club",
        icon: Plus,
        action: () => {
          openClubModal();
          router.push("/clubs");
        },
        keywords: ["create", "creer", "nouveau", "new", "club", "association"],
      },
      {
        id: "action-new-event",
        label: "Créer un événement",
        description: "Ajouter un nouvel événement",
        icon: Plus,
        action: () => {
          openEventModal();
          router.push("/events");
        },
        keywords: ["create", "creer", "nouveau", "new", "event", "evenement"],
      },
      {
        id: "system-logout",
        label: "Déconnexion",
        description: "Se déconnecter de l'administration",
        icon: LogOut,
        action: () => {
          logout();
          router.push("/login");
        },
        keywords: ["logout", "deconnexion", "disconnect", "exit", "sortir"],
      },
    ],
    [router, logout, openUserModal, openClubModal, openEventModal],
  );

  // Combine all commands
  const allCommands = useMemo(() => {
    return [...staticCommands, ...dynamicCommands];
  }, [staticCommands, dynamicCommands]);

  // Memoize filtered commands to avoid recalculation on every render
  const filteredCommands = useMemo(() => {
    if (!search.trim()) return staticCommands; // Show only static commands by default

    const searchLower = search.toLowerCase();
    return allCommands
      .filter(
        (command) =>
          command.label?.toLowerCase().includes(searchLower) ||
          command.description?.toLowerCase().includes(searchLower) ||
          command.keywords?.some((keyword) => keyword?.toLowerCase().includes(searchLower)),
      )
      .slice(0, 50); // Show more results since we include all data
  }, [allCommands, staticCommands, search]);

  // Handle keyboard shortcuts using the custom hook
  useKeyboardShortcuts(
    [
      {
        key: "k",
        ctrlKey: true,
        callback: () => setCommandPaletteOpen(!commandPaletteOpen),
      },
      {
        key: "k",
        metaKey: true,
        callback: () => setCommandPaletteOpen(!commandPaletteOpen),
      },
      {
        key: "Escape",
        callback: () => commandPaletteOpen && setCommandPaletteOpen(false),
      },
    ],
    [commandPaletteOpen],
  );

  // Reset search when opening
  useEffect(() => {
    if (commandPaletteOpen) {
      setSearch("");
    }
  }, [commandPaletteOpen]);

  if (!commandPaletteOpen) return null;

  const handleSelect = (commandId: string) => {
    const command = allCommands.find((cmd) => cmd.id === commandId);
    if (command) {
      command.action();
      setCommandPaletteOpen(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black/50 z-50 flex items-start justify-center pt-[20vh]">
      <Command className="bg-white rounded-lg shadow-2xl w-full max-w-2xl mx-4 overflow-hidden">
        <div className="flex items-center border-b px-3">
          <Search className="h-5 w-5 text-gray-400 mr-2" />
          <Command.Input
            value={search}
            onValueChange={setSearch}
            placeholder="Tapez une commande ou recherchez..."
            className="flex-1 py-3 text-lg bg-transparent border-none outline-none placeholder-gray-400"
            autoFocus
          />
          <kbd className="ml-2 text-xs text-gray-500 bg-gray-100 px-2 py-1 rounded">ESC</kbd>
        </div>

        <Command.List className="max-h-96 overflow-y-auto">
          <Command.Empty className="py-6 text-center text-sm text-gray-500">
            Aucune commande trouvée pour &quot;{search}&quot;.
          </Command.Empty>

          {/* Render command groups dynamically */}
          {filteredCommands.filter((cmd) => cmd.id.startsWith("nav-")).length > 0 && (
            <Command.Group heading="Navigation" className="p-2">
              {filteredCommands
                .filter((cmd) => cmd.id.startsWith("nav-"))
                .map((command) => (
                  <CommandItemComponent
                    key={command.id}
                    command={command}
                    onSelect={handleSelect}
                  />
                ))}
            </Command.Group>
          )}

          {filteredCommands.filter((cmd) => cmd.id.startsWith("action-")).length > 0 && (
            <Command.Group heading="Actions" className="p-2">
              {filteredCommands
                .filter((cmd) => cmd.id.startsWith("action-"))
                .map((command) => (
                  <CommandItemComponent
                    key={command.id}
                    command={command}
                    onSelect={handleSelect}
                  />
                ))}
            </Command.Group>
          )}

          {filteredCommands.filter((cmd) => cmd.id.startsWith("user-")).length > 0 && (
            <Command.Group heading="Utilisateurs" className="p-2">
              {filteredCommands
                .filter((cmd) => cmd.id.startsWith("user-"))
                .map((command) => (
                  <CommandItemComponent
                    key={command.id}
                    command={command}
                    onSelect={handleSelect}
                  />
                ))}
            </Command.Group>
          )}

          {filteredCommands.filter((cmd) => cmd.id.startsWith("club-")).length > 0 && (
            <Command.Group heading="Clubs" className="p-2">
              {filteredCommands
                .filter((cmd) => cmd.id.startsWith("club-"))
                .map((command) => (
                  <CommandItemComponent
                    key={command.id}
                    command={command}
                    onSelect={handleSelect}
                  />
                ))}
            </Command.Group>
          )}

          {filteredCommands.filter((cmd) => cmd.id.startsWith("event-")).length > 0 && (
            <Command.Group heading="Événements" className="p-2">
              {filteredCommands
                .filter((cmd) => cmd.id.startsWith("event-"))
                .map((command) => (
                  <CommandItemComponent
                    key={command.id}
                    command={command}
                    onSelect={handleSelect}
                  />
                ))}
            </Command.Group>
          )}

          {filteredCommands.filter((cmd) => cmd.id.startsWith("system-")).length > 0 && (
            <Command.Group heading="Système" className="p-2">
              {filteredCommands
                .filter((cmd) => cmd.id.startsWith("system-"))
                .map((command) => (
                  <CommandItemComponent
                    key={command.id}
                    command={command}
                    onSelect={handleSelect}
                  />
                ))}
            </Command.Group>
          )}
        </Command.List>

        <div className="border-t px-4 py-2 text-xs text-gray-500 flex justify-between">
          <div>
            <kbd className="bg-gray-100 px-1 py-0.5 rounded text-xs mr-1">↑↓</kbd>
            pour naviguer
          </div>
          <div>
            <kbd className="bg-gray-100 px-1 py-0.5 rounded text-xs mr-1">↵</kbd>
            pour sélectionner
          </div>
          <div>
            <kbd className="bg-gray-100 px-1 py-0.5 rounded text-xs">⌘K</kbd>
            pour basculer
          </div>
        </div>
      </Command>

      {/* Background overlay to close on click */}
      <button
        className="absolute inset-0 -z-10"
        type="button"
        tabIndex={0}
        onClick={() => setCommandPaletteOpen(false)}
        onKeyDown={(e) => {
          if (e.key === "Enter" || e.key === " ") {
            e.preventDefault();
            setCommandPaletteOpen(false);
          }
        }}
      />
    </div>
  );
}

// Memoized CommandItemComponent for better performance
const CommandItemComponent = memo<{
  command: CommandItem;
  onSelect: (commandId: string) => void;
}>(({ command, onSelect }) => {
  const Icon = command.icon;

  return (
    <Command.Item
      key={command.id}
      value={`${command.label} ${command.description ?? ""} ${(command.keywords ?? []).join(" ")}`}
      onSelect={() => onSelect(command.id)}
      className="flex items-center space-x-3 px-3 py-3 cursor-pointer hover:bg-gray-100 rounded-md"
    >
      <Icon className="h-4 w-4 text-gray-500 shrink-0" />
      <div className="flex-1 min-w-0">
        <div className="text-sm font-medium text-gray-900 truncate">{command.label}</div>
        {command.description && (
          <div className="text-xs text-gray-500 truncate">{command.description}</div>
        )}
      </div>
    </Command.Item>
  );
});

CommandItemComponent.displayName = "CommandItemComponent";

export default memo(CommandPalette);
