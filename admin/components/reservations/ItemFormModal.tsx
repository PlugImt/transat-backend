"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { CalendarClock, Package } from "lucide-react";
import { useEffect, useMemo } from "react";
import { useForm } from "react-hook-form";
import toast from "react-hot-toast";
import { z } from "zod";
import Modal from "@/components/Modal";
import type { ReservationTreeItem } from "@/lib/api";
import { useCreateReservationItem, useUpdateReservationItem } from "@/lib/hooks";
import type { AdminItem, ItemStatus } from "@/lib/reservationTypes";
import type { Club } from "@/lib/types";
import { apiErrorMessage, fieldClass, ITEM_STATUS_LABELS, LoadingButton } from "./ui";

// Limits of the columns of reservation_element.
const schema = z.object({
  name: z.string().trim().min(1, "Le nom est requis").max(100, "100 caractères maximum"),
  slot: z.boolean(),
  description: z.string().max(100, "100 caractères maximum"),
  location: z.string().max(50, "50 caractères maximum"),
  status: z.enum(["active", "maintenance", "inactive"]),
  status_note: z.string().max(200, "200 caractères maximum"),
  warning_message: z.string().max(500, "500 caractères maximum"),
  confirmation_message: z.string().max(500, "500 caractères maximum"),
  club_id: z.string(),
  category_id: z.string(),
});

type FormValues = z.infer<typeof schema>;

const EMPTY: FormValues = {
  name: "",
  slot: false,
  description: "",
  location: "",
  status: "active",
  status_note: "",
  warning_message: "",
  confirmation_message: "",
  club_id: "",
  category_id: "",
};

interface ItemFormModalProps {
  isOpen: boolean;
  onClose: () => void;
  /** The item being edited; none to create one. */
  item?: AdminItem | null;
  defaultClubId?: number | null;
  defaultCategoryId?: number | null;
  clubs: Club[];
  tree: ReservationTreeItem[];
}

const flattenCategories = (nodes: ReservationTreeItem[], level = 0) => {
  const out: { id: number; name: string; clubId?: number }[] = [];
  for (const node of nodes) {
    if (node.type === "category" && node.id) {
      out.push({
        id: node.id,
        name: `${"— ".repeat(level)}${node.name} (${node.club_name})`,
        clubId: node.club_id,
      });
      if (node.children) out.push(...flattenCategories(node.children, level + 1));
    }
  }
  return out;
};

function FieldError({ message }: { message?: string }) {
  return message ? <p className="mt-1 text-xs text-red-600">{message}</p> : null;
}

export function ItemFormModal({
  isOpen,
  onClose,
  item,
  defaultClubId,
  defaultCategoryId,
  clubs,
  tree,
}: ItemFormModalProps) {
  const createItem = useCreateReservationItem();
  const updateItem = useUpdateReservationItem();
  const categories = useMemo(() => flattenCategories(tree), [tree]);

  const {
    register,
    handleSubmit,
    reset,
    watch,
    setValue,
    setError,
    formState: { errors, isSubmitting },
  } = useForm<FormValues>({ resolver: zodResolver(schema), defaultValues: EMPTY });

  const slot = watch("slot");
  const status = watch("status");
  const categoryId = watch("category_id");
  const pending = isSubmitting || createItem.isPending || updateItem.isPending;

  // Switching between slots and loans under reservations that are in progress would make them meaningless.
  const typeLocked = !!item && (!!item.current || item.upcoming_count > 0);

  useEffect(() => {
    if (!isOpen) return;
    if (item) {
      reset({
        name: item.name,
        slot: item.slot,
        description: item.description ?? "",
        location: item.location ?? "",
        status: item.status,
        status_note: item.status_note ?? "",
        warning_message: item.warning_message ?? "",
        confirmation_message: item.confirmation_message ?? "",
        club_id: String(item.club_id),
        category_id: item.category_id ? String(item.category_id) : "",
      });
    } else {
      reset({
        ...EMPTY,
        club_id: defaultClubId ? String(defaultClubId) : "",
        category_id: defaultCategoryId ? String(defaultCategoryId) : "",
      });
    }
  }, [isOpen, item, defaultClubId, defaultCategoryId, reset]);

  // A category fixes the club.
  useEffect(() => {
    if (item || !categoryId) return;
    const category = categories.find((c) => String(c.id) === categoryId);
    if (category?.clubId) setValue("club_id", String(category.clubId));
  }, [categoryId, categories, item, setValue]);

  const submit = handleSubmit(async (values) => {
    if (!item && !values.category_id && !values.club_id) {
      setError("club_id", { message: "Choisissez un club ou une catégorie" });
      return;
    }
    try {
      if (item) {
        await updateItem.mutateAsync({
          id: item.id,
          data: {
            name: values.name.trim(),
            slot: typeLocked ? undefined : values.slot,
            description: values.description.trim(),
            location: values.location.trim(),
            status: values.status,
            status_note: values.status === "active" ? "" : values.status_note.trim(),
            warning_message: values.warning_message.trim(),
            confirmation_message: values.confirmation_message.trim(),
          },
        });
        toast.success("Élément modifié");
      } else {
        const created = await createItem.mutateAsync({
          name: values.name.trim(),
          slot: values.slot,
          description: values.description.trim() || undefined,
          location: values.location.trim() || undefined,
          status: values.status,
          id_club_parent: values.club_id ? Number(values.club_id) : undefined,
          id_category_parent: values.category_id ? Number(values.category_id) : undefined,
        });
        const id: number | undefined = created?.item?.id;
        const extras =
          values.status_note.trim() ||
          values.warning_message.trim() ||
          values.confirmation_message.trim();
        if (id && extras) {
          await updateItem.mutateAsync({
            id,
            data: {
              status_note: values.status === "active" ? "" : values.status_note.trim(),
              warning_message: values.warning_message.trim(),
              confirmation_message: values.confirmation_message.trim(),
            },
          });
        }
        toast.success("Élément créé");
      }
      onClose();
    } catch (error) {
      toast.error(apiErrorMessage(error, "Erreur lors de l'enregistrement"));
    }
  });

  return (
    <Modal
      isOpen={isOpen}
      onClose={pending ? () => {} : onClose}
      title={item ? "Modifier l'élément" : "Nouvel élément"}
      maxWidth="xl"
    >
      <form onSubmit={submit} className="space-y-4" noValidate>
        <div>
          <label htmlFor="item-name" className="mb-1 block text-sm font-medium text-gray-700">
            Nom *
          </label>
          <input
            id="item-name"
            className={fieldClass}
            placeholder="Salle de répétition, Câble HDMI…"
            {...register("name")}
          />
          <FieldError message={errors.name?.message} />
        </div>

        <fieldset>
          <legend className="mb-1 text-sm font-medium text-gray-700">Type de réservation</legend>
          <div className="grid gap-2 sm:grid-cols-2">
            {[
              {
                value: true,
                icon: <CalendarClock className="h-5 w-5" />,
                title: "Par créneaux",
                text: "Réservé par tranches d'1 h, de 8 h à minuit (heure de Paris). Pour les salles et le matériel d'activité.",
              },
              {
                value: false,
                icon: <Package className="h-5 w-5" />,
                title: "Emprunt continu",
                text: "L'utilisateur prend l'objet et le rend quand il a fini. Des rappels partent après 1, 5 et 7 jours.",
              },
            ].map((option) => (
              <label
                key={option.title}
                className={`flex cursor-pointer gap-3 rounded-lg border p-3 text-sm ${
                  slot === option.value
                    ? "border-blue-500 bg-blue-50"
                    : "border-gray-200 hover:bg-gray-50"
                } ${typeLocked ? "cursor-not-allowed opacity-60" : ""}`}
              >
                <input
                  type="radio"
                  className="mt-1"
                  checked={slot === option.value}
                  disabled={typeLocked}
                  onChange={() => setValue("slot", option.value)}
                />
                <span>
                  <span className="flex items-center gap-1 font-medium text-gray-900">
                    {option.icon}
                    {option.title}
                  </span>
                  <span className="mt-1 block text-gray-600">{option.text}</span>
                </span>
              </label>
            ))}
          </div>
          {typeLocked && (
            <p className="mt-1 text-xs text-amber-700">
              Le type ne peut pas changer tant que des réservations sont en cours ou à venir.
            </p>
          )}
        </fieldset>

        <div className="grid gap-4 sm:grid-cols-2">
          <div>
            <label htmlFor="item-status" className="mb-1 block text-sm font-medium text-gray-700">
              État
            </label>
            <select id="item-status" className={fieldClass} {...register("status")}>
              {(Object.keys(ITEM_STATUS_LABELS) as ItemStatus[]).map((value) => (
                <option key={value} value={value}>
                  {ITEM_STATUS_LABELS[value]}
                </option>
              ))}
            </select>
            <p className="mt-1 text-xs text-gray-500">
              {status === "active"
                ? "Visible dans l'application et réservable."
                : "Masqué dans l'application et non réservable. Les réservations existantes restent."}
            </p>
          </div>
          {status !== "active" && (
            <div>
              <label
                htmlFor="item-status-note"
                className="mb-1 block text-sm font-medium text-gray-700"
              >
                Motif (visible ici uniquement)
              </label>
              <input
                id="item-status-note"
                className={fieldClass}
                placeholder="Écran cassé, en réparation…"
                {...register("status_note")}
              />
              <FieldError message={errors.status_note?.message} />
            </div>
          )}
        </div>

        {!item && (
          <div className="grid gap-4 sm:grid-cols-2">
            <div>
              <label
                htmlFor="item-category"
                className="mb-1 block text-sm font-medium text-gray-700"
              >
                Catégorie
              </label>
              <select id="item-category" className={fieldClass} {...register("category_id")}>
                <option value="">Aucune (à la racine du club)</option>
                {categories.map((category) => (
                  <option key={category.id} value={category.id}>
                    {category.name}
                  </option>
                ))}
              </select>
            </div>
            <div>
              <label htmlFor="item-club" className="mb-1 block text-sm font-medium text-gray-700">
                Club {!categoryId && "*"}
              </label>
              <select
                id="item-club"
                className={fieldClass}
                disabled={!!categoryId}
                {...register("club_id")}
              >
                <option value="">Sélectionner un club</option>
                {clubs.map((club) => (
                  <option key={club.id_clubs} value={club.id_clubs}>
                    {club.name}
                  </option>
                ))}
              </select>
              <FieldError message={errors.club_id?.message} />
              {categoryId && <p className="mt-1 text-xs text-gray-500">Hérité de la catégorie.</p>}
            </div>
          </div>
        )}
        {item && (
          <p className="text-xs text-gray-500">
            {item.club_name}
            {item.category_path ? ` · ${item.category_path}` : ""} — l'emplacement ne se modifie pas
            ici.
          </p>
        )}

        <div className="grid gap-4 sm:grid-cols-2">
          <div>
            <label
              htmlFor="item-description"
              className="mb-1 block text-sm font-medium text-gray-700"
            >
              Description
            </label>
            <input id="item-description" className={fieldClass} {...register("description")} />
            <FieldError message={errors.description?.message} />
          </div>
          <div>
            <label htmlFor="item-location" className="mb-1 block text-sm font-medium text-gray-700">
              Localisation
            </label>
            <input
              id="item-location"
              className={fieldClass}
              placeholder="Bâtiment, étage…"
              {...register("location")}
            />
            <FieldError message={errors.location?.message} />
          </div>
        </div>

        <details
          className="rounded-lg border border-gray-200 p-3"
          open={!!(item?.warning_message || item?.confirmation_message)}
        >
          <summary className="cursor-pointer text-sm font-medium text-gray-700">
            Messages affichés dans l'application
          </summary>
          <div className="mt-3 space-y-3">
            <div>
              <label htmlFor="item-warning" className="mb-1 block text-sm text-gray-700">
                Avertissement avant la réservation
              </label>
              <textarea
                id="item-warning"
                rows={2}
                className={fieldClass}
                {...register("warning_message")}
              />
              <FieldError message={errors.warning_message?.message} />
            </div>
            <div>
              <label htmlFor="item-confirmation" className="mb-1 block text-sm text-gray-700">
                Message de confirmation
              </label>
              <textarea
                id="item-confirmation"
                rows={2}
                className={fieldClass}
                {...register("confirmation_message")}
              />
              <FieldError message={errors.confirmation_message?.message} />
            </div>
          </div>
        </details>

        <div className="flex justify-end gap-2 border-t pt-4">
          <LoadingButton
            onClick={onClose}
            disabled={pending}
            className="border border-gray-300 text-gray-700 hover:bg-gray-50"
          >
            Annuler
          </LoadingButton>
          <LoadingButton
            type="submit"
            loading={pending}
            className="bg-blue-600 text-white hover:bg-blue-700"
          >
            {item ? "Enregistrer" : "Créer"}
          </LoadingButton>
        </div>
      </form>
    </Modal>
  );
}
