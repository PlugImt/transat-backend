"use client";

import {
  CalendarClock,
  ChevronDown,
  ChevronUp,
  Package,
  PackageSearch,
  Pencil,
  Plus,
  Search,
  Trash2,
} from "lucide-react";
import { useMemo, useState } from "react";
import toast from "react-hot-toast";
import Modal from "@/components/Modal";
import { useDeleteReservationItem } from "@/lib/hooks";
import { useReservationCatalogue } from "@/lib/reservationHooks";
import type { AdminItem, ItemAvailability } from "@/lib/reservationTypes";
import {
  AVAILABILITY_LABELS,
  AvailabilityBadge,
  apiErrorMessage,
  Badge,
  EmptyState,
  fieldClass,
  formatDay,
  formatDuration,
  formatTime,
  fullName,
  LoadingButton,
  SkeletonRows,
} from "./ui";

interface CatalogueTabProps {
  onOpen: (item: AdminItem) => void;
  onEdit: (item: AdminItem) => void;
  onCreate: () => void;
}

type SortKey = "name" | "availability" | "waitlist";

const norm = (value: string | undefined) =>
  (value ?? "")
    .normalize("NFD")
    .replace(/\p{Diacritic}/gu, "")
    .toLowerCase();

const DAY = 86_400;

function StatusCell({ item }: { item: AdminItem }) {
  if (item.availability === "maintenance" || item.availability === "inactive") {
    return <span className="text-xs text-gray-500">{item.status_note ?? "—"}</span>;
  }
  const holder = item.current;
  if (item.availability === "borrowed" && holder) {
    const late = holder.elapsed_seconds >= DAY;
    return (
      <div className="text-sm">
        <div className="font-medium text-gray-900">
          {fullName(holder.first_name, holder.last_name)}
        </div>
        <div className={late ? "text-xs font-medium text-amber-700" : "text-xs text-gray-500"}>
          depuis {formatDuration(holder.elapsed_seconds)}
          {holder.elapsed_seconds >= 7 * DAY
            ? " · +7 j"
            : holder.elapsed_seconds >= 5 * DAY
              ? " · +5 j"
              : late
                ? " · +1 j"
                : ""}
        </div>
      </div>
    );
  }
  if (item.availability === "occupied" && holder?.end_date) {
    return (
      <div className="text-sm">
        <div className="font-medium text-gray-900">
          {fullName(holder.first_name, holder.last_name)}
        </div>
        <div className="text-xs text-gray-500">jusqu'à {formatTime(holder.end_date)}</div>
      </div>
    );
  }
  if (item.slot && item.next) {
    return (
      <div className="text-xs text-gray-600">
        Prochain : {formatDay(item.next.start_date)} {formatTime(item.next.start_date)}
        <div className="text-gray-500">{fullName(item.next.first_name, item.next.last_name)}</div>
      </div>
    );
  }
  return <span className="text-xs text-gray-400">Rien de prévu</span>;
}

export function CatalogueTab({ onOpen, onEdit, onCreate }: CatalogueTabProps) {
  const { data, isLoading, isError, error, refetch } = useReservationCatalogue();
  const deleteItem = useDeleteReservationItem();

  const [search, setSearch] = useState("");
  const [category, setCategory] = useState("");
  const [kind, setKind] = useState<"" | "slot" | "loan">("");
  const [availability, setAvailability] = useState<"" | ItemAvailability>("");
  const [sort, setSort] = useState<{ key: SortKey; desc: boolean }>({ key: "name", desc: false });
  const [deleting, setDeleting] = useState<AdminItem | null>(null);

  const items = data ?? [];
  const categories = useMemo(
    () =>
      [...new Set(items.map((item) => item.category_path || `${item.club_name} (racine)`))].sort(),
    [items],
  );

  const visible = useMemo(() => {
    const needle = norm(search.trim());
    const filtered = items.filter((item) => {
      if (kind === "slot" && !item.slot) return false;
      if (kind === "loan" && item.slot) return false;
      if (availability && item.availability !== availability) return false;
      if (category && (item.category_path || `${item.club_name} (racine)`) !== category)
        return false;
      if (!needle) return true;
      const haystack = norm(
        [
          item.name,
          item.description,
          item.location,
          item.category_path,
          item.club_name,
          item.current && fullName(item.current.first_name, item.current.last_name),
          item.current?.email,
        ]
          .filter(Boolean)
          .join(" "),
      );
      return haystack.includes(needle);
    });
    const order = ["borrowed", "occupied", "available", "maintenance", "inactive"];
    const compare = (a: AdminItem, b: AdminItem) => {
      switch (sort.key) {
        case "availability":
          return (
            order.indexOf(a.availability) - order.indexOf(b.availability) ||
            a.name.localeCompare(b.name)
          );
        case "waitlist":
          return a.waitlist_count - b.waitlist_count || a.name.localeCompare(b.name);
        default:
          return a.name.localeCompare(b.name, "fr");
      }
    };
    return filtered.sort((a, b) => (sort.desc ? -compare(a, b) : compare(a, b)));
  }, [items, search, category, kind, availability, sort]);

  const filtersActive = !!(search || category || kind || availability);
  const toggleSort = (key: SortKey) =>
    setSort((previous) => ({ key, desc: previous.key === key ? !previous.desc : false }));

  const SortHeader = ({ id, children }: { id: SortKey; children: React.ReactNode }) => (
    <button
      type="button"
      onClick={() => toggleSort(id)}
      className="inline-flex items-center gap-1 uppercase"
    >
      {children}
      {sort.key === id &&
        (sort.desc ? <ChevronDown className="h-3 w-3" /> : <ChevronUp className="h-3 w-3" />)}
    </button>
  );

  const confirmDelete = async () => {
    if (!deleting) return;
    try {
      await deleteItem.mutateAsync(deleting.id);
      toast.success(`${deleting.name} supprimé`);
      setDeleting(null);
    } catch (e) {
      toast.error(apiErrorMessage(e, "Impossible de supprimer cet élément"));
    }
  };

  return (
    <div className="space-y-4">
      <div className="flex flex-wrap items-center gap-3">
        <div className="relative min-w-[16rem] flex-1">
          <Search className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
          <input
            type="search"
            value={search}
            onChange={(event) => setSearch(event.target.value)}
            placeholder="Rechercher un élément, une salle, un emprunteur…"
            className={`${fieldClass} pl-9`}
            aria-label="Rechercher dans le catalogue"
          />
        </div>
        <select
          value={category}
          onChange={(event) => setCategory(event.target.value)}
          className={`${fieldClass} w-auto`}
          aria-label="Catégorie"
        >
          <option value="">Toutes les catégories</option>
          {categories.map((value) => (
            <option key={value} value={value}>
              {value}
            </option>
          ))}
        </select>
        <select
          value={kind}
          onChange={(event) => setKind(event.target.value as typeof kind)}
          className={`${fieldClass} w-auto`}
          aria-label="Type de réservation"
        >
          <option value="">Tous les types</option>
          <option value="slot">Par créneaux</option>
          <option value="loan">Emprunts continus</option>
        </select>
        <select
          value={availability}
          onChange={(event) => setAvailability(event.target.value as typeof availability)}
          className={`${fieldClass} w-auto`}
          aria-label="Disponibilité"
        >
          <option value="">Toutes les disponibilités</option>
          {(Object.keys(AVAILABILITY_LABELS) as ItemAvailability[]).map((value) => (
            <option key={value} value={value}>
              {AVAILABILITY_LABELS[value]}
            </option>
          ))}
        </select>
        {filtersActive && (
          <button
            type="button"
            onClick={() => {
              setSearch("");
              setCategory("");
              setKind("");
              setAvailability("");
            }}
            className="text-sm text-blue-600 hover:underline"
          >
            Réinitialiser
          </button>
        )}
        <button
          type="button"
          onClick={onCreate}
          className="ml-auto inline-flex items-center gap-2 rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
        >
          <Plus className="h-4 w-4" />
          Nouvel élément
        </button>
      </div>

      <div className="overflow-hidden rounded-lg bg-white shadow">
        {isLoading ? (
          <SkeletonRows rows={8} columns={5} />
        ) : isError ? (
          <div className="p-8 text-center text-sm text-red-600">
            {apiErrorMessage(error, "Impossible de charger le catalogue")}{" "}
            <button type="button" onClick={() => refetch()} className="underline">
              Réessayer
            </button>
          </div>
        ) : items.length === 0 ? (
          <EmptyState
            icon={<PackageSearch className="h-10 w-10" />}
            title="Aucun élément réservable"
            description="Créez une catégorie dans l'onglet Catégories, puis ajoutez les salles et le matériel."
          />
        ) : visible.length === 0 ? (
          <EmptyState
            icon={<Search className="h-10 w-10" />}
            title="Aucun résultat"
            description="Aucun élément ne correspond à ces filtres."
          />
        ) : (
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200 text-sm">
              <thead className="bg-gray-50 text-left text-xs font-medium text-gray-500">
                <tr>
                  <th className="px-4 py-3">
                    <SortHeader id="name">Élément</SortHeader>
                  </th>
                  <th className="px-4 py-3 uppercase">Type</th>
                  <th className="px-4 py-3">
                    <SortHeader id="availability">État</SortHeader>
                  </th>
                  <th className="px-4 py-3 uppercase">Détail</th>
                  <th className="px-4 py-3">
                    <SortHeader id="waitlist">File</SortHeader>
                  </th>
                  <th className="px-4 py-3" />
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-100">
                {visible.map((item) => (
                  <tr key={item.id} className="hover:bg-gray-50">
                    <td className="px-4 py-3">
                      <button
                        type="button"
                        onClick={() => onOpen(item)}
                        className="text-left font-medium text-gray-900 hover:text-blue-700 hover:underline"
                      >
                        {item.name}
                      </button>
                      <div className="text-xs text-gray-500">
                        {item.category_path || `${item.club_name} (racine)`}
                        {item.location ? ` · ${item.location}` : ""}
                      </div>
                    </td>
                    <td className="px-4 py-3">
                      <Badge tone={item.slot ? "blue" : "gray"}>
                        {item.slot ? (
                          <CalendarClock className="h-3 w-3" />
                        ) : (
                          <Package className="h-3 w-3" />
                        )}
                        {item.slot ? "Créneaux" : "Emprunt"}
                      </Badge>
                    </td>
                    <td className="px-4 py-3">
                      <AvailabilityBadge value={item.availability} />
                    </td>
                    <td className="px-4 py-3">
                      <StatusCell item={item} />
                    </td>
                    <td className="px-4 py-3">
                      {item.waitlist_count > 0 ? (
                        <Badge tone="blue">{item.waitlist_count}</Badge>
                      ) : (
                        <span className="text-gray-300">—</span>
                      )}
                    </td>
                    <td className="px-4 py-3 text-right">
                      <button
                        type="button"
                        onClick={() => onEdit(item)}
                        className="rounded p-1.5 text-indigo-600 hover:bg-indigo-50"
                        title="Modifier"
                        aria-label={`Modifier ${item.name}`}
                      >
                        <Pencil className="h-4 w-4" />
                      </button>
                      <button
                        type="button"
                        onClick={() => setDeleting(item)}
                        className="rounded p-1.5 text-red-600 hover:bg-red-50"
                        title="Supprimer"
                        aria-label={`Supprimer ${item.name}`}
                      >
                        <Trash2 className="h-4 w-4" />
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
        {!isLoading && items.length > 0 && (
          <div className="border-t border-gray-100 px-4 py-2 text-xs text-gray-500">
            {visible.length} élément{visible.length > 1 ? "s" : ""} sur {items.length}
          </div>
        )}
      </div>

      <Modal
        isOpen={deleting !== null}
        onClose={() => setDeleting(null)}
        title="Supprimer l'élément"
      >
        {deleting && (
          <div className="space-y-4 text-sm text-gray-700">
            <p>
              Supprimer <strong>{deleting.name}</strong> ? Son historique de réservations et ses
              files d'attente seront supprimés. Impossible tant qu'une réservation est en cours ou à
              venir.
            </p>
            <div className="flex justify-end gap-2">
              <LoadingButton
                onClick={() => setDeleting(null)}
                className="border border-gray-300 text-gray-700 hover:bg-gray-50"
              >
                Annuler
              </LoadingButton>
              <LoadingButton
                onClick={confirmDelete}
                loading={deleteItem.isPending}
                className="bg-red-600 text-white hover:bg-red-700"
              >
                Supprimer
              </LoadingButton>
            </div>
          </div>
        )}
      </Modal>
    </div>
  );
}
