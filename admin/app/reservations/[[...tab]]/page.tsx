"use client";

import { CalendarDays, FolderTree, ListChecks, Package } from "lucide-react";
import { useParams, useRouter } from "next/navigation";
import { useState } from "react";
import { CatalogueTab } from "@/components/reservations/CatalogueTab";
import { CategoryTree } from "@/components/reservations/CategoryTree";
import { ItemDrawer } from "@/components/reservations/ItemDrawer";
import { ItemFormModal } from "@/components/reservations/ItemFormModal";
import { PlanningView } from "@/components/reservations/PlanningView";
import { ReservationsList } from "@/components/reservations/ReservationsList";
import { StatsCards } from "@/components/reservations/StatsCards";
import { useClubs, useReservationTree } from "@/lib/hooks";
import { useReservationCatalogue } from "@/lib/reservationHooks";
import type { AdminItem, ReservationFilters } from "@/lib/reservationTypes";

type Tab = "catalogue" | "reservations" | "planning" | "categories";

const TABS: { id: Tab; label: string; icon: typeof Package }[] = [
  { id: "catalogue", label: "Catalogue", icon: Package },
  { id: "reservations", label: "Réservations", icon: ListChecks },
  { id: "planning", label: "Planning", icon: CalendarDays },
  { id: "categories", label: "Catégories", icon: FolderTree },
];

const tabPath = (id: Tab) => (id === "catalogue" ? "/reservations" : `/reservations/${id}`);

export default function ReservationsPage() {
  const router = useRouter();
  const segment = useParams<{ tab?: string[] }>().tab?.[0];
  const tab = TABS.find(({ id }) => id === segment)?.id ?? "catalogue";
  const setTab = (id: Tab) => router.push(tabPath(id));
  const [filters, setFilters] = useState<ReservationFilters>({});
  const [drawerItemId, setDrawerItemId] = useState<number | null>(null);
  const [form, setForm] = useState<{ open: boolean; item: AdminItem | null }>({
    open: false,
    item: null,
  });

  const { data: catalogue } = useReservationCatalogue();
  const { data: tree } = useReservationTree();
  const { data: clubs } = useClubs();
  const items = catalogue ?? [];

  // The drawer reads the live item, so a refresh after an action is reflected in it.
  const drawerItem = items.find((item) => item.id === drawerItemId) ?? null;

  const showReservations = (next: ReservationFilters) => {
    setFilters(next);
    setTab("reservations");
  };

  return (
    <div className="space-y-6 p-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Gestion des réservations</h1>
        <p className="mt-1 text-sm text-gray-600">
          Catalogue des salles et du matériel, réservations en cours, annulations et files
          d'attente.
        </p>
      </div>

      <StatsCards
        onShowActive={() => showReservations({ status: "active", kind: "loan" })}
        onShowOverdue={(days) => showReservations({ status: "overdue", overdue_days: days })}
        onShowCancelled={() => showReservations({ status: "cancelled" })}
      />

      <div className="border-b border-gray-200">
        <nav className="-mb-px flex gap-6" aria-label="Sections">
          {TABS.map(({ id, label, icon: Icon }) => (
            <button
              key={id}
              type="button"
              onClick={() => setTab(id)}
              aria-current={tab === id ? "page" : undefined}
              className={`inline-flex items-center gap-2 border-b-2 px-1 py-3 text-sm font-medium ${
                tab === id
                  ? "border-blue-600 text-blue-700"
                  : "border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700"
              }`}
            >
              <Icon className="h-4 w-4" />
              {label}
            </button>
          ))}
        </nav>
      </div>

      {tab === "catalogue" && (
        <CatalogueTab
          onOpen={(item) => setDrawerItemId(item.id)}
          onEdit={(item) => setForm({ open: true, item })}
          onCreate={() => setForm({ open: true, item: null })}
        />
      )}
      {tab === "reservations" && (
        <ReservationsList
          filters={filters}
          onFiltersChange={setFilters}
          items={items}
          onOpenItem={setDrawerItemId}
        />
      )}
      {tab === "planning" && <PlanningView items={items} onOpenItem={setDrawerItemId} />}
      {tab === "categories" && <CategoryTree onOpenItem={setDrawerItemId} />}

      <ItemDrawer
        item={drawerItem}
        onClose={() => setDrawerItemId(null)}
        onEdit={(item) => setForm({ open: true, item })}
      />

      <ItemFormModal
        isOpen={form.open}
        onClose={() => setForm({ open: false, item: null })}
        item={form.item}
        clubs={clubs ?? []}
        tree={tree ?? []}
      />
    </div>
  );
}
